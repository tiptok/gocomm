package gzcache

import (
	"database/sql"
	. "github.com/tiptok/gocomm/pkg/cache/model"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/syncx"
	"time"
)

var (
	exclusiveCalls = syncx.NewSingleFlight()
	stats          = cache.NewStat("sqlc")
)

type GZCache struct {
	cache cache.Cache
}

func NewNodeCache(host, pass string, opts ...cache.Option) GZCache {
	return GZCache{cache.New(cache.CacheConf(
		[]cache.NodeConf{
			cache.NodeConf{RedisConf: redis.RedisConf{Host: host, Pass: pass, Type: "node"}, Weight: 100},
		}),
		exclusiveCalls, stats, sql.ErrNoRows, opts...)}
}

// get cached value by key.
func (c GZCache) Get(key string, obj interface{}) (*Item, error) {
	var it Item
	it.Object = obj
	err := c.cache.Get(key, &it)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &it, err
}

// set cached item
func (c GZCache) Set(k string, it *Item) error {
	return c.cache.SetWithExpire(k, it, time.Second*time.Duration(it.TTL))
}

// delete cached value by key.
func (c GZCache) Delete(key string) error {
	return c.cache.Del(key)
}
