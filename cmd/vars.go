package cmd

import "fmt"

type Item struct {
	Action string         `json: "action" msgpack:"action"`
	Bucket string         `json: "bucket" msgpack:"bucket"`
	Data   map[string]any `json: "data"   msgpack:"data"`
}

type Info struct {
	Message string           `json: "message" msgpack:"message"`
	Status  string           `json: "status"  msgpack:"status"`
	Bucket  string           `json: "bucket"  msgpack:"bucket"`
	Data    []map[string]any `json: "data"    msgpack:"data"`
}

func NewInfo(b string) Info {
	info := Info{
		Message: "",
		Status:  "",
		Bucket:  b,
		Data:    []map[string]any{},
	}

	return info
}

func ParseItemData(itemData map[string]any) map[string]string {
	fid := ""
	r2Path := ""
	r2MIME := ""
	mcValue := ""
	for k, v := range itemData {
		//fmt.Printf("==== v's type ==== %T\n", v)
		if k == "_id" {
			vs, ok := v.(string)
			if !ok {

				DebugWarn("", "_id type shoule be string")
			} else {
				fid = vs
			}
		}

		if k == "path" {
			vs, ok := v.(string)
			if !ok {
				DebugWarn("", "path type shoule be string")
			} else {
				r2Path = vs
			}
		}
		if k == "mime" {
			vs, ok := v.(string)
			if !ok {
				DebugWarn("", "mime type shoule be string")
			} else {
				r2MIME = vs
			}
		}
		if k == "mc_value" {
			vs, err := AnyNumber2String(v)
			if err != nil {
				DebugWarn("", "value type shoule be string")
			} else {
				mcValue = vs
			}

		}
	}

	res := map[string]string{
		"_id":      fid,
		"path":     r2Path,
		"mime":     r2MIME,
		"mc_value": mcValue,
	}

	DebugInfo("ParseItemData", fmt.Sprintf("%+v", res))

	return res
}
