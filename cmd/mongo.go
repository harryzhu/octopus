package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/kataras/iris/v12"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/writeconcern"
)

type mgoGetOption struct {
	Skip       int64
	Limit      int64
	Sort       any
	Projection any
}

type mgoGetFilter struct {
	Filter any
}

var (
	mgoClient        *mongo.Client
	mgoDB            *mongo.Database
	defaultGetFilter mgoGetFilter
	defaultGetOption mgoGetOption
)

func initMongo() {

	mgoconn := GetEnv("MONGOCONN", "")
	mgodatabase := GetEnv("MONGODATABASE", "")

	if IsAnyEmpty(mgoconn, mgodatabase) {
		DebugWarn("mongodb init", "cannot get env vars: MONGOCONN / MONGODATABASE, mongodb service is not available")
	}
	DebugInfo("MongoDB Database", mgodatabase)

	mgoClient = mgoGetClient(mgoconn)

	mgoDB = mgoClient.Database(mgodatabase)

	defaultGetFilter = mgoGetFilter{
		Filter: bson.D{},
	}
	defaultGetOption = mgoGetOption{
		Skip:       0,
		Limit:      1000,
		Sort:       bson.D{{"_id", 1}},
		Projection: nil}
}

func mgoGetClient(conn string) *mongo.Client {
	m, err := mongo.Connect(options.Client().SetConnectTimeout(3 * time.Second).ApplyURI(conn))
	if err != nil {
		PrintError("mgoGetClient", err)
	}

	return m
}

func mgoPing() error {
	err := mgoClient.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}
	return nil
}

func mgoKV2BsonD(kv map[string]any) bson.D {
	var bd bson.D
	for k, v := range kv {
		be := bson.E{Key: k, Value: v}
		bd = append(bd, be)
	}
	return bd
}

func mgoSave(c string, _id string, kvs iris.Map) error {
	if IsAnyEmpty(c, _id) {
		return NewError("collection/_id cannot be empty")
	}

	ts := GetNowUnix()

	kvs["_id"] = _id
	kvs["create_at"] = ts
	kvs["update_at"] = ts

	coll := mgoDB.Collection(c)
	res, err := coll.InsertOne(context.TODO(), kvs)
	if err != nil {
		DebugWarn("mgoSave.10", err.Error())
		return err
	}

	DebugInfo("mgoSave.20", res.InsertedID)

	return nil
}

func mgoUpdate(c string, _id string, kvs iris.Map) (int64, error) {
	if IsAnyEmpty(c, _id) {
		return 0, NewError("collection/_id cannot be empty")
	}
	DebugInfo("mgoUpdate.10", c, ":", _id)

	kvs["update_at"] = GetNowUnix()
	filter := bson.D{{"_id", _id}}
	opts := options.UpdateOne().SetUpsert(true)

	coll := mgoDB.Collection(c)

	var modifiedRows int64
	for k, v := range kvs {
		if k == "_id" {
			continue
		}
		if len(k) > 0 {
			fmt.Printf("--------- %T\n", v)
			update := bson.D{{"$set", bson.D{{k, v}}}}
			res, err := coll.UpdateOne(context.TODO(), filter, update, opts)
			if err != nil {
				DebugWarn("mgoUpdate.20 Error", err.Error())
				return 0, err
			}
			DebugInfo("mgoUpdate.30 OK", k, " ==> ", v)
			modifiedRows += res.ModifiedCount

		}
	}

	return modifiedRows, nil
}

func mgoUpdateWithTransaction(c string, _id string, kvs iris.Map) (int64, error) {
	if IsAnyEmpty(c, _id) {
		return 0, NewError("collection/_id cannot be empty")
	}

	// start-session
	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	session, err := mgoClient.StartSession()
	if err != nil {
		PrintError("mgoUpdate", err)
		return 0, err
	}
	// Defers ending the session after the transaction is committed or ended
	defer session.EndSession(context.TODO())

	filter := bson.D{{"_id", _id}}
	opts := options.UpdateOne().SetUpsert(true)
	coll := mgoDB.Collection(c)

	kvs["update_at"] = GetNowUnix()

	result, err := session.WithTransaction(context.TODO(), func(ctx context.Context) (any, error) {
		var modifiedRows int64
		for k, v := range kvs {
			if k == "_id" {
				continue
			}
			if len(k) > 0 {

				fmt.Printf("---%s--%v--|-- %T-----\n", k, v, v)
				update := bson.D{{"$set", bson.D{{k, v}}}}

				res, err := coll.UpdateOne(context.TODO(), filter, update, opts)
				if err != nil {
					break
					return 0, err
				}

				DebugInfo("mgoUpdateWithTransaction.2 OK", k, " ==> ", v)
				modifiedRows += res.ModifiedCount
			}
		}
		return modifiedRows, err
	}, txnOptions)

	if err != nil {
		DebugWarn("mgoUpdate.2", err.Error())
		session.AbortTransaction(context.TODO())
		session.EndSession(context.TODO())
		return 0, err
	}
	session.CommitTransaction(context.TODO())
	session.EndSession(context.TODO())

	return result.(int64), nil
}

func mgoDelete(c string, _id string) (int64, error) {
	if IsAnyEmpty(c, _id) {
		return 0, NewError("collection/_id cannot be empty")
	}
	DebugInfo("mgoDelete.10", c, ":", _id)

	filter := bson.D{{"_id", _id}}
	opts := options.DeleteOne().SetCollation(&options.Collation{
		Locale:    "en_US",
		Strength:  1,
		CaseLevel: false,
	})

	coll := mgoDB.Collection(c)
	res, err := coll.DeleteOne(context.TODO(), filter, opts)
	if err != nil {
		DebugWarn("mgoDelete.20", err.Error())
		return 0, err
	}
	DebugInfo("mgoDelete.30", res.DeletedCount)

	return res.DeletedCount, nil
}

func mgoGet(c string, mgfilter mgoGetFilter, mgopt mgoGetOption) ([]bson.M, error) {
	var bmResult []bson.M

	coll := mgoDB.Collection(c)
	opts := options.Find()
	opts.SetSort(mgopt.Sort)
	opts.SetProjection(mgopt.Projection)
	opts.SetSkip(mgopt.Skip)
	opts.SetLimit(mgopt.Limit)

	filter := mgfilter.Filter

	fmt.Printf("mgoGet:filter: %v\n", filter)
	fmt.Printf("mgoGet:opts:%+v\n", opts)

	cursor, err := coll.Find(context.TODO(), filter, opts)

	if err != nil {
		PrintError("mgoGet: find", err)
		return bmResult, err
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		PrintError("mgoGet: cursor.All", err)
		return bmResult, err
	}

	for _, result := range results {
		bmResult = append(bmResult, result)
	}

	return bmResult, nil
}

func mgoParseGetFilter(itemDataFilter map[string]any) mgoGetFilter {
	filter := mgoGetFilter{
		Filter: bson.D{},
	}
	b := Map2Bson(itemDataFilter)
	if b != nil {
		filter.Filter = b
		DebugInfo("mgoParseGetFilter:", filter)
	}
	return filter
}

func mgoParseGetOption(itemDataOption map[string]any) mgoGetOption {
	option := mgoGetOption{
		Skip:       0,
		Limit:      1000,
		Sort:       bson.D{{"_id", 1}},
		Projection: nil,
	}

	for k, v := range itemDataOption {
		DebugInfo("mgoParseGetOption", k, "=>", fmt.Sprintf("%T", v), ": ", v)
		if k == "skip" {
			vskip, err := AnyInt2Int64(v)
			if err != nil {
				PrintError("mgoParseGetOption: skip", err)
			} else {
				option.Skip = vskip
			}

		}

		if k == "limit" {
			vlimit, err := AnyInt2Int64(v)
			if err != nil {
				PrintError("mgoParseGetOption: limit", err)
			} else {
				option.Skip = vlimit
			}
		}

		if k == "sort" {
			mi, ok := v.(map[string]any)
			if ok {
				b := Map2Bson(mi)
				if b != nil {
					option.Sort = b
				}
			} else {
				PrintError("mgoParseGetOption: sort", NewError("sort is invalid"))
			}

		}

		if k == "projection" {

			switch v.(type) {
			case []any:
				option.Projection = nil
			case map[string]any:
				b := Map2Bson(v.(map[string]any))
				if b != nil {
					option.Projection = b
				}

			default:
				PrintError("mgoParseGetOption: projection", NewError("projection is invalid"))
			}

		}

	}

	return option

}
