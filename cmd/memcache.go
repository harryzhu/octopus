package cmd

import (
	"strings"

	"github.com/kataras/iris/v12"
)

func actionMemcache(ctx iris.Context) {

	item := ctx.Values().Get("requestBody").(Item)
	result := NewInfo(item.Bucket)
	result.Data = []iris.Map{}

	memcacheItemData := ParseItemData(item.Data)

	fid := memcacheItemData["_id"]
	mcvalue := memcacheItemData["mc_value"]

	tStart := GetNowUnixMilli()

	if item.Action == "SAVE" {
		result.Message = "_id / Data.mc_value cannot be empty"
		result.Status = "error"

		if len(mcvalue) > maxEntrySize {
			result.Message = strings.Join([]string{"Data.mc_value length is oversized(max: ", Int2Str(maxEntrySize), " bytes)"}, "")
		}

		if fid != "" && mcvalue != "" && len(mcvalue) <= maxEntrySize {
			DebugInfo("actionMemcache: SAVE", fid)

			bkey := bcacheKeyJoin(item.Bucket, fid)
			bval := []byte(mcvalue)

			err := bcacheSet(bkey, bval)
			if err != nil {
				result.Message = err.Error()
				result.Status = "error"
			} else {
				result.Message = "saved successfully"
				result.Status = "ok"
				result.Data = append(result.Data, iris.Map{
					"_id":       fid,
					"elapse_ms": GetNowUnixMilli() - tStart,
				})
			}
		}
	}

	if item.Action == "GET" {
		DebugInfo("actionMemcache: GET", fid)

		result.Message = "_id cannot be empty"
		result.Status = "error"

		if fid != "" {
			bkey := bcacheKeyJoin(item.Bucket, fid)
			bval := bcacheGet(bkey)

			if bval != nil {
				result.Message = "get successfully"
				result.Status = "ok"
				result.Data = append(result.Data, iris.Map{
					"_id":       fid,
					"mc_value":  string(bval),
					"elapse_ms": GetNowUnixMilli() - tStart,
				})
			} else {
				result.Message = "get failed"
				result.Status = "error"
				result.Data = append(result.Data, iris.Map{
					"_id":       fid,
					"mc_value":  "",
					"elapse_ms": GetNowUnixMilli() - tStart,
				})
			}

		}
	}

	if item.Action == "DELETE" {
		DebugInfo("actionMemcache: DELETE", fid)

		result.Message = "_id cannot be empty"
		result.Status = "error"

		if fid != "" {
			bkey := bcacheKeyJoin(item.Bucket, fid)
			err := bcacheDelete(bkey)

			if err != nil {
				result.Message = err.Error()
				result.Status = "error"
				result.Data = append(result.Data, iris.Map{
					"_id":       fid,
					"elapse_ms": GetNowUnixMilli() - tStart,
				})
			} else {
				result.Message = "deleted successfully"
				result.Status = "ok"
				result.Data = append(result.Data, iris.Map{
					"_id":       fid,
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

	ctx.JSON(result)

}
