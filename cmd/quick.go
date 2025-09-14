package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
)

func uuidQuick(ctx iris.Context) {
	count := ctx.URLParamDefault("count", "10")
	var ids []string
	c, err := strconv.Atoi(count)
	if err != nil || c == 0 {
		c = 10
	}

	for range c {
		u := uuid.New()
		ids = append(ids, fmt.Sprintf("%v", u))
	}

	ctx.StatusCode(iris.StatusOK)

	ctx.JSON(iris.Map{
		"algo":  "uuid",
		"count": c,
		"ids":   ids,
	})

}

func snowflakeQuick(ctx iris.Context) {
	nodeID := ctx.URLParamDefault("nodeid", "0")
	count := ctx.URLParamDefault("count", "10")
	if nodeID == "" {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("nodeid cannot be empty")
		return
	}
	nID := Str2Int64(nodeID, 0)
	node, err := snowflake.NewNode(nID)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	c, err := strconv.Atoi(count)
	if err != nil || c == 0 {
		c = 10
	}

	var ids []string

	for range c {
		id := node.Generate()
		ids = append(ids, fmt.Sprintf("%v", id))
	}

	ctx.StatusCode(iris.StatusOK)

	ctx.JSON(iris.Map{
		"algo":    "snowflake",
		"count":   c,
		"node_id": nID,
		"ids":     ids,
	})
}

func memcacheListQuick(ctx iris.Context) {
	prefix := ctx.Params().Get("prefix")
	kvs := bcacheScan(prefix)
	DebugInfo("memcacheListQuick", len(kvs))
	css := `
	<style type="text/css">
		table {
			border-collapse: collapse;
			border: 1px solid #ccc;
		}
		td{
			border: none;
			word-break: break-all;
			border-bottom: 1px solid #ccc;
			border-right: 1px solid #ccc;
			padding: 5px 10px;
		}
	</style>
	`
	s := strings.Join([]string{"<h2>", prefix, "</h2><table>"}, "")
	for k, v := range kvs {
		if k != "" {
			s = strings.Join([]string{s, "<tr><td>", k, "</td><td>", v, "</tr>"}, "")
		}
	}
	s = strings.Join([]string{css, s, "</table>"}, "")
	ctx.Header("Content-Type", "text/html;charset=utf-8")
	ctx.WriteString(s)
}

func memcacheGetQuick(ctx iris.Context) {
	prefix := ctx.Params().Get("key")
	val := bcacheGet(prefix)
	ctx.Header("Content-Type", "text/plain;charset=utf-8")
	if val != nil {
		ctx.StatusCode(iris.StatusOK)
		ctx.Write(val)
	} else {
		ctx.StatusCode(iris.StatusOK)
		ctx.WriteString("")
	}

}

func serviceStatus(ctx iris.Context) {
	r2bucket := ctx.URLParam("r2bucket")
	DebugInfo("", r2bucket)
	r2Status := "ERROR"
	mongoStatus := "ERROR"
	bcacheStatus := "ERROR"

	if WithR2 {
		if r2bucket == "" {
			r2Status = "please provide your r2 bucket name with url parameter `?r2bucket=your-bucket-name`"
		} else {
			r2err := R2Ping(r2bucket)
			if r2err == nil {
				r2Status = "OK"
			}
		}

	} else {
		r2Status = "DISABLED"
	}

	if WithMongo {
		mgoError := mgoPing()
		if mgoError == nil {
			mongoStatus = "OK"
		}
	} else {
		mongoStatus = "DISABLED"
	}

	if WithMemcache {
		bcacheError := bcachePing()
		if bcacheError == nil {
			bcacheStatus = "OK"
		}
	} else {
		bcacheStatus = "DISABLED"
	}

	status := iris.Map{
		"mongodb":  mongoStatus,
		"r2":       r2Status,
		"memcache": bcacheStatus,
	}

	ctx.JSON(status)
}
