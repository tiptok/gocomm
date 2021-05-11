package cache

import (
	memory "github.com/tiptok/gocomm/pkg/cache/memory"
	redis_cache "github.com/tiptok/gocomm/pkg/cache/redis_cache"
)

var mlCache *MultiLevelCache

func InitDefault(option ...Option) {
	mlCache = NewMultiLevelCacheNew(option...)
	mlCache.RegisterCache(
		memory.NewMemCache(mlCache.Options.CleanInterval),
		redis_cache.NewRedisCache(mlCache.Options.DefaultRedisPool),
	)
	return
}

func GetObject(key string, obj interface{}, ttl int, f LoadFunc) error {
	return mlCache.GetObject(key, obj, ttl, f)
}

func Delete(key string) error {
	return mlCache.Delete(key)
}

func InitMultiLevelCache(option ...Option) *MultiLevelCache {
	if mlCache == nil {
		mlCache = NewMultiLevelCacheNew(option...)
	}
	return mlCache
}

func RegisterCache(cache ...Cache) {
	if mlCache == nil {
		mlCache = NewMultiLevelCacheNew()
	}
	mlCache.RegisterCache(cache...)
}
