package cmd

import (
	//"bytes"
	"context"
	"encoding/json"

	"fmt"
	"strings"
	"time"

	"github.com/allegro/bigcache/v3"
)

var bcache *bigcache.BigCache

const maxEntrySize int = 512

var MaxCacheSize int = MemcacheSizeMB << 20

func initBigcache() {
	var err error
	if MaxCacheSize < 16<<20 {
		MaxCacheSize = 16 << 20
	}
	config := bigcache.Config{
		Shards:             256,
		LifeWindow:         10 * time.Minute,
		CleanWindow:        5 * time.Minute,
		MaxEntriesInWindow: 1000 * 10 * 60,
		MaxEntrySize:       maxEntrySize,
		Verbose:            true,
		HardMaxCacheSize:   MaxCacheSize,
		OnRemove:           nil,
		OnRemoveWithReason: nil,
	}

	if IsDebug {
		config.LifeWindow = 30 * time.Second
		config.CleanWindow = 15 * time.Second
	}

	DebugInfo("memcacheInit:LifeWindow", config.LifeWindow)
	DebugInfo("memcacheInit:CleanWindow", config.CleanWindow)
	DebugInfo("memcacheInit:HardMaxCacheSize", config.HardMaxCacheSize)

	bcache, err = bigcache.New(context.Background(), config)
	PrintError("memcacheInit", err)
}

func bcachePing() error {
	errSet := bcacheSet("test", []byte("ok"))
	if errSet != nil {
		return errSet
	}
	b := bcacheGet("test")
	if b == nil {
		return NewError("bcache get error")
	}
	if string(b) != "ok" {
		return NewError("bcache set/get value error")
	}
	return nil
}

func bcacheSet(k string, v []byte) error {
	return bcache.Set(k, v)
}

func bcacheGet(k string) []byte {
	v, err := bcache.Get(k)
	if err != nil {
		return nil
	}
	//DebugInfo("bcacheGet from cache", k)
	return v
}

func bcacheDelete(k string) error {

	err := bcache.Delete(k)
	if err != nil {
		if err == bigcache.ErrEntryNotFound {
			return nil
		}
	} else {
		return err
	}

	return nil
}

func jsonEnc(data any) []byte {
	b, err := json.Marshal(data)
	if err != nil {
		PrintError("jsonEnc", err)
		return nil
	}
	return b
}

func jsonDec(data []byte, dataStruct any) error {
	err := json.Unmarshal(data, &dataStruct)
	PrintError("jsonDec", err)
	return err
}

func bcacheKeyJoin(args ...any) string {
	var info []string
	for _, arg := range args {
		info = append(info, fmt.Sprintf("%v", arg))
	}
	return strings.Join(info, "::")
}

func bcacheScan(prefix string) (data map[string]string) {
	if prefix == "" {
		return data
	}

	data = make(map[string]string, 100)
	iterator := bcache.Iterator()
	count := 0
	valSafe := ""

	for iterator.SetNext() {
		if count > 1000 {
			break
		}
		current, err := iterator.Value()
		PrintError("bcacheScan", err)
		k := current.Key()
		var val []byte

		if strings.HasPrefix(k, prefix) {
			val = current.Value()
			if val == nil {
				continue
			}
			valSafe = string(val)
			if len(valSafe) > 512 {
				valSafe = strings.Join([]string{valSafe[0:512], "..."}, " ")
			}
			data[k] = valSafe
			count++

		}

	}
	return data
}
