package cmd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func GetNowUnix() int64 {
	return time.Now().UTC().Unix()
}

func GetNowUnixMilli() int64 {
	return time.Now().UTC().UnixMilli()
}

func ToUnixSlash(s string) string {
	// for windows
	return strings.ReplaceAll(s, "\\", "/")
}

func GetEnv(k string, defaultVal string) string {
	ev := strings.Trim(os.Getenv(k), " ")
	if ev == "" {
		return defaultVal
	}
	return ev
}

func IsAnyEmpty(args ...string) bool {
	for _, arg := range args {
		if arg == "" {
			return true
		}
	}
	return false
}

func MakeDirs(dpath string) error {
	_, err := os.Stat(dpath)
	if err != nil {
		DebugInfo("MakeDirs", dpath)
		err = os.MkdirAll(dpath, os.ModePerm)
		PrintError("MakeDirs:MkdirAll", err)
	}
	return nil
}

func Str2Float64(n string, defaultValue float64) float64 {
	s, err := strconv.ParseFloat(n, 64)
	if err != nil {
		PrintError("Str2Float64", err)
		return defaultValue
	}
	return s
}

func Str2Int(n string, defaultValue int) int {
	s, err := strconv.Atoi(n)
	if err != nil {
		PrintError("Str2Int", err)
		return defaultValue
	}
	return s
}

func Str2Int64(n string, defaultValue int64) int64 {
	s, err := strconv.ParseInt(n, 10, 64)
	if err != nil {
		PrintError("Str2Int64", err)
		return defaultValue
	}
	return s
}

func Int2Int64(n int, defaultValue int64) int64 {
	s := strconv.Itoa(n)
	m, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		PrintError("Int2Int64", err)
		return defaultValue
	}
	return m
}

func Int2Str(n int) string {
	return strconv.Itoa(n)
}

func Int64ToString(n int64) string {
	return strconv.FormatInt(n, 10)
}

func Json2Bson(s string) (b any, err error) {
	err = bson.UnmarshalExtJSON([]byte(s), false, &b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Json2Map(j string) map[string]any {
	var m map[string]any
	err := json.Unmarshal([]byte(j), &m)
	if err != nil {
		DebugWarn("Json2Map", err)
		return nil
	}
	return m
}

func Map2Json(m map[string]any) string {
	j, err := json.Marshal(m)
	if err != nil {
		DebugWarn("Map2Json", err)
		return ""
	}
	return string(j)
}

func Map2Bson(m map[string]any) []byte {
	b, err := bson.Marshal(m)
	if err != nil {
		DebugWarn("Map2Bson", err)
		return nil
	}
	return b
}

func JsonMarshal(s any) string {
	bf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(bf)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.Encode(s)
	return bf.String()
}

func Base64Dec(s string) string {
	dec, err := base64.StdEncoding.DecodeString(s)
	PrintError("Base64Dec", err)
	return string(dec)
}

func AnyInt2Int64(v any) (int64, error) {
	err := NewError("invalid value")
	switch v.(type) {
	case int, uint, int8, uint8, int16, uint16, int32, uint32:
		vs := fmt.Sprintf("%v", v)
		vs64, err := strconv.ParseInt(vs, 10, 64)
		if err != nil {
			return 0, err
		}
		return vs64, nil
	case int64:
		return v.(int64), nil
	case uint64:
		u64, ok := v.(uint64)
		if !ok {
			return 0, err
		}
		vs := fmt.Sprintf("%v", u64)
		vs64, err := strconv.ParseInt(vs, 10, 64)
		//var min64 int64 = -9223372036854775807
		//var max64 int64 = 9223372036854775807
		if err != nil {
			return 0, err
		}

		return vs64, nil

	default:
		return 0, err
	}
}

func AnyFloat2Float64(v any) (float64, error) {
	err := NewError("invalid value")
	switch v.(type) {
	case float32:
		vs := fmt.Sprintf("%f", v)
		f64, err := strconv.ParseFloat(vs, 64)
		if err != nil {
			return 0, err
		}
		return f64, nil
	case float64:
		return v.(float64), nil

	default:
		return 0, err
	}
}

func AnyNumber2String(v any) (string, error) {
	err := NewError("invalid value")
	vStr := ""
	switch v.(type) {
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64, float32, float64, string:
		vStr = fmt.Sprintf("%v", v)
	default:
		return "", err
	}
	return vStr, nil
}
