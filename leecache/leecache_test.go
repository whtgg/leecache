package leecache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func TestGetter(t *testing.T) {
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key),nil
	})

	expect := []byte("key")
	if v,_ := f.Get("key");!reflect.DeepEqual(v,expect) {
		t.Errorf("callback failed")
	}
}

var db = map[string]string {
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGet(t *testing.T) {
	loadcounts := make(map[string]int,len(db))
	lee := NewGroup("scores",2<<10, GetterFunc(func(key string) ([]byte, error) {
		log.Println("[SlowDB] search key", key)
		if v,ok := db[key];ok {
			if _,ok := loadcounts[key];!ok {
				loadcounts[key]  = 0
 			}
 			loadcounts[key] += 1
 			return []byte(v),nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}))

	for k,v := range db {
		if view,err := lee.Get(k); err != nil || view.String() != v {
			t.Fatal("failed to get value of Tom")
		}
		if _, err := lee.Get(k); err != nil || loadcounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		}
	}

	if view, err := lee.Get("unknown"); err == nil {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}
}