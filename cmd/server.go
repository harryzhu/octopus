package cmd

import (
	stdContext "context"
	"fmt"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/basicauth"
)

const ActionAllowed string = "SAVE|UPDATE|DELETE|GET"

var app *iris.Application

func StartHTTPServer() {
	app = iris.New()

	homeAPI := app.Party("/")
	{
		homeAPI.Use(iris.Compression)

		homeAPI.Get("/", homeIndex)
	}

	r2API := app.Party("/r2/")
	{
		if AdminUser != "" && AdminPassword != "" {
			auth := basicauth.Default(map[string]string{
				AdminUser: AdminPassword,
			})
			r2API.Use(auth)
		}
		r2API.Use(mustPost)
		r2API.Use(decodeBody)

		r2API.Use(iris.Compression)
		r2API.Post("action", actionR2)
	}

	mongoAPI := app.Party("/mongo/")
	{
		if AdminUser != "" && AdminPassword != "" {
			auth := basicauth.Default(map[string]string{
				AdminUser: AdminPassword,
			})
			mongoAPI.Use(auth)
		}
		mongoAPI.Use(mustPost)
		mongoAPI.Use(decodeBody)

		mongoAPI.Use(iris.Compression)
		mongoAPI.Post("action", actionMongo)
	}

	quickAPI := app.Party("/quick/")
	{
		quickAPI.Get("uuid", uuidQuick)
		quickAPI.Get("snowflake", snowflakeQuick)
		quickAPI.Get("memcache/list/{prefix:path}", memcacheListQuick)
		quickAPI.Get("memcache/get/{key:path}", memcacheGetQuick)
	}

	memcacheAPI := app.Party("/memcache/")
	{
		if AdminUser != "" && AdminPassword != "" {
			auth := basicauth.Default(map[string]string{
				AdminUser: AdminPassword,
			})
			memcacheAPI.Use(auth)
		}
		memcacheAPI.Use(mustPost)
		memcacheAPI.Use(decodeBody)

		memcacheAPI.Post("action", actionMemcache)
	}

	statusAPI := app.Party("/status/")
	{
		statusAPI.Get("/", serviceStatus)
	}

	app.Listen(fmt.Sprintf("%s:%d", Host, Port))

}

func StopHTTPServer() {
	fmt.Println("stopping the server ...")
	appShutdown()
}

func appShutdown() {
	ctx, cancel := stdContext.WithTimeout(stdContext.Background(), 3*time.Second)
	defer cancel()
	// close all hosts
	err := app.Shutdown(ctx)
	if err != nil {
		PrintError("appShutdown:", err)
	}
}

func notFound(ctx iris.Context) {
	ctx.StatusCode(iris.StatusNotFound)
	ctx.JSON(iris.Map{})
}

func homeIndex(ctx iris.Context) {
	resp := iris.Map{}
	ctx.JSON(resp)
}

func mustPost(ctx iris.Context) {

	if ctx.Method() != "POST" {
		DebugInfo("ERROR: mustPost", "pls use POST")
		ctx.StopWithError(iris.StatusMethodNotAllowed, NewError("pls use POST"))
		ctx.StopExecution()
		return
	}

	DebugInfo("mustPost", "Next")

	ctx.Next()
}

func decodeBody(ctx iris.Context) {
	item := Item{}

	err := ctx.ReadMsgPack(&item)

	if err != nil {
		DebugWarn("decodeBody", err.Error())
		ctx.StopWithError(iris.StatusBadRequest, NewError(err.Error()))
		ctx.StopExecution()
		return
	}

	if IsAnyEmpty(item.Action, item.Bucket) {
		errStr := "action/bucket cannot be empty"
		DebugWarn("decodeBody", errStr)
		ctx.StopWithError(iris.StatusBadRequest, NewError(errStr))
		ctx.StopExecution()
		return
	}

	if strings.Contains(ActionAllowed, item.Action) != true {
		errStr := strings.Join([]string{"action must be one of", ActionAllowed}, ": ")
		DebugWarn("decodeBody", errStr)
		ctx.StopWithError(iris.StatusBadRequest, NewError(errStr))
		ctx.StopExecution()
		return
	}

	ctx.Values().Set("requestBody", item)

	ctx.Next()
}

func actionR2(ctx iris.Context) {

	item := ctx.Values().Get("requestBody").(Item)

	DebugInfo("actionR2", item.Data)

	result := NewInfo(item.Bucket)
	result.Data = []iris.Map{}

	tStart := GetNowUnixMilli()

	r2ItemData := ParseItemData(item.Data)

	fid := r2ItemData["_id"]
	fpath := r2ItemData["path"]
	fmime := r2ItemData["mime"]

	if item.Action == "SAVE" {
		DebugInfo("actionR2: SAVE", fpath, ", mime: ", fmime)

		result.Message = "file _id/path/mime cannot be empty"
		result.Status = "error"

		if fid != "" && fpath != "" && fmime != "" {
			err := R2Save(item.Bucket, fid, fpath, fmime)
			if err != nil {
				result.Message = err.Error()
				result.Status = "error"
			} else {
				result.Status = "ok"
				result.Message = "saved successfully"
				result.Data = append(result.Data, iris.Map{
					"_id":       fid,
					"action":    "r2.SAVE",
					"elapse_ms": GetNowUnixMilli() - tStart,
				})
			}
		}

	}

	if item.Action == "DELETE" {
		DebugInfo("actionR2: DELETE", fid)

		result.Message = "file _id cannot be empty"
		result.Status = "error"

		if fid != "" {
			err := R2Delete(item.Bucket, fid)
			if err != nil {
				result.Message = err.Error()
				result.Status = "error"
			} else {
				result.Message = "deleted successfully"
				result.Status = "ok"
				result.Data = append(result.Data, iris.Map{
					"_id":       fid,
					"action":    "r2.DELETE",
					"elapse_ms": GetNowUnixMilli() - tStart,
				})
			}
		}

	}

	if result.Status == "ok" {
		ctx.StatusCode(iris.StatusOK)
	} else {
		ctx.StatusCode(iris.StatusBadRequest)
	}

	ctx.JSON(result, iris.JSON{
		UnescapeHTML: true,
	})
}

func actionMongo(ctx iris.Context) {
	ctx.Header("Content-Type", "text/json;charset=UTF-8")
	item := ctx.Values().Get("requestBody").(Item)

	result := NewInfo(item.Bucket)
	result.Data = []iris.Map{}

	DebugInfo("actionMongo.1", item)

	mongoItemData := ParseItemData(item.Data)

	fid := mongoItemData["_id"]

	tStart := GetNowUnixMilli()

	if item.Action == "DELETE" {
		DebugInfo("actionMongo: DELETE", fid)

		result.Message = "mongo _id cannot be empty"
		result.Status = "error"

		if fid != "" {
			deletedCount, err := mgoDelete(item.Bucket, fid)
			if err != nil {
				result.Message = err.Error()
				result.Status = "error"
				DebugWarn("actionMongo.DELETE", result.Message)
			} else {
				result.Message = "deleted successfully"
				result.Status = "ok"
				result.Data = append(result.Data, iris.Map{
					"_id":           fid,
					"action":        "mongo.DELETE",
					"deleted_count": deletedCount,
					"elapse_ms":     GetNowUnixMilli() - tStart,
				})
			}
		}

	}

	if item.Action == "SAVE" {
		DebugInfo("actionMongo: SAVE", fid)

		result.Message = "mongo _id cannot be empty"
		result.Status = "error"

		if fid != "" {
			err := mgoSave(item.Bucket, fid, item.Data)
			if err != nil {
				result.Message = err.Error()
				result.Status = "error"
				DebugWarn("actionMongo.SAVE", result.Message)
			} else {
				result.Message = "saved successfully"
				result.Status = "ok"
				result.Data = append(result.Data, iris.Map{
					"_id":       fid,
					"action":    "mongo.SAVE",
					"elapse_ms": GetNowUnixMilli() - tStart,
				})
			}
		}
	}

	if item.Action == "UPDATE" {
		DebugInfo("actionMongo: UPDATE", fid)

		result.Message = "mongo _id cannot be empty"
		result.Status = "error"

		var modifiedRows int64
		var err error
		if fid != "" {
			if len(item.Data) < 5 {
				modifiedRows, err = mgoUpdate(item.Bucket, fid, item.Data)
			} else {
				modifiedRows, err = mgoUpdateWithTransaction(item.Bucket, fid, item.Data)
			}
			PrintError("actionMongo.UPDATE", err)

			if err != nil {
				result.Message = err.Error()
				result.Status = "error"
			} else {
				result.Message = "updated successfully"
				result.Status = "ok"
				result.Data = append(result.Data, iris.Map{
					"_id":           fid,
					"action":        "mongo.UPDATE",
					"modified_rows": modifiedRows,
					"elapse_ms":     GetNowUnixMilli() - tStart,
				})
			}
		}

	}

	if item.Action == "GET" {
		result.Message = "cannot get result"
		result.Status = "error"

		DebugInfo("actionMongo: GET", item.Data)
		mgfilter := defaultGetFilter
		mgoption := defaultGetOption

		for k, v := range item.Data {
			if k == "filter" {
				vf, ok := v.(map[string]any)
				if ok {
					mgfilter = mgoParseGetFilter(vf)
				} else {
					fmt.Printf("item.Data.filter Type--- %T : %v \n", v, v)
					DebugWarn("actionMongo: GET", "filter is invalid")
				}

			}

			if k == "option" {
				vo, ok := v.(map[string]any)
				if ok {
					mgoption = mgoParseGetOption(vo)
				} else {
					fmt.Printf("---option Type--- %T \n", v)
					DebugWarn("actionMongo: GET", "option is invalid")
				}

			}

		}

		fmt.Printf("---FINAL: filter--- %+v \n", mgfilter)
		fmt.Printf("---FINAL: option--- %+v \n", mgoption)

		rows, err := mgoGet(item.Bucket, mgfilter, mgoption)
		if err != nil {
			result.Message = err.Error()
			result.Status = "error"
		} else {
			result.Message = "GET successfully"
			result.Status = "ok"
			for _, row := range rows {
				result.Data = append(result.Data, row)
			}
		}

	}

	if result.Status == "ok" {
		ctx.StatusCode(iris.StatusOK)
	} else {
		ctx.StatusCode(iris.StatusBadRequest)
	}

	ctx.JSON(result, iris.JSON{
		UnescapeHTML: true,
	})
}
