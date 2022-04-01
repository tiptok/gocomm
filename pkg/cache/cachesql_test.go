package cache

import (
	"github.com/tiptok/gocomm/pkg/cache/gzcache"
	"testing"
)

func TestQueryUniqueIndexCache(t *testing.T) {
	mlc := InitMultiLevelCache()
	mlc.RegisterCache(gzcache.NewNodeCache("127.0.0.1:6379", ""))
	csql := NewCachedRepository(mlCache)
	cacheKeyFunc := func() string {
		return "user:simnum:18860"
	}
	queryPrimaryKeyFunc := func(o interface{}) string {
		return "user:1001"
	}
	queryFunc := func() (interface{}, error) {
		return user{
			Id:     "1001",
			Name:   "tiptok",
			Simnum: "18860",
		}, nil
	}
	var v user
	err := csql.QueryUniqueIndexCache(cacheKeyFunc, &v, queryPrimaryKeyFunc, queryFunc)
	if err != nil {
		t.Fatal(err)
	}
	if v.Name != "tiptok" {
		t.Fatal("except not equal")
	}
}

type user struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Simnum string `json:"simnum"`
}
