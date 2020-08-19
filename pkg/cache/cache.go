package cache

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	. "github.com/tiptok/gocomm/pkg/cache/model"
	. "github.com/tiptok/gocomm/pkg/cache/redis_cache"
	"sync"
)

const (
	//缓存不需要重新加载
	noReload = iota
	//缓存部分重新加载
	someReload
	//缓存全部重新加载
	allReload
)

type Cache interface {
	// get cached value by key.
	Get(key string, obj interface{}) (*Item, error)
	// set cached item
	Set(k string, it *Item) error
	// delete cached value by key.
	Delete(key string) error
}

type LoadFunc func() (interface{}, error)

type MultiLevelCache struct {
	Options *Options

	pool *redis.Pool

	mlCache *MLCache
	// RWMutex map for each cache key
	muxm sync.Map
}

type MLCache struct {
	Current Cache
	Next    *MLCache
}

func NewMultiLevelCacheNew(option ...Option) *MultiLevelCache {
	o := NewOptions(option...)
	if o.DefaultRedisPool == nil {
		panic("redis pool is nil")
	}
	c := &MultiLevelCache{
		Options: o,
	}
	c.pool = o.DefaultRedisPool

	// subscribe key deletion
	go c.subscribe(o.DeleteChannel)
	return c
}

//GetObject  get cache from multilevelCache
func (c *MultiLevelCache) GetObject(key string, obj interface{}, ttl int, f LoadFunc) error {
	return c.getObjectWithExpiration(key, obj, ttl, f)
}
func (c *MultiLevelCache) getObjectWithExpiration(key string, obj interface{}, ttl int, f LoadFunc) error {
	var (
		item *Item
		err  error
	)
	if c.mlCache == nil {
		return fmt.Errorf("mlCache is nil")
	}
	if c.mlCache.Current == nil {
		return fmt.Errorf("mlCache is nil")
	}
	cacheLink := c.mlCache
	deep := -1
	reload := noReload //0:不需要重新加载 1.重新加载部分 2.重新加载全部
	for {
		if cacheLink == nil || cacheLink.Current == nil {
			reload = allReload
			break
		}
		cache := cacheLink.Current
		if item, err = cache.Get(key, obj); err != nil || item == nil {
			deep++
			reload = someReload
			cacheLink = cacheLink.Next
			continue
		}
		if item.Expire() {
			deep++
			reload = someReload
			cacheLink = cacheLink.Next
			continue
		} else {
			break
		}
	}

	switch reload {
	case someReload:
		if item == nil {
			break
		}
		c.TraverseCache(deep, func(c Cache) error {
			return c.Set(key, item)
		})
		break
	case allReload:
		err = c.Load(key, obj, ttl, f, deep)
		return err
	}

	return Clone(item.Object, obj)
}

//Load  get cache if expire or not exist , create new cache to multiLevelCache
func (c *MultiLevelCache) Load(key string, obj interface{}, ttl int, f LoadFunc, deep int) error {
	mux := c.GetMutex(key)
	mux.Lock()
	defer func() {
		mux.Unlock()
		c.ReleaseMutex(key)
	}()
	if v, err := c.mlCache.Current.Get(key, obj); err != nil || !itemNeedReload(v) {
		if v != nil {
			err = Clone(v.Object, obj)
		}
		return err
	}
	o, err := f()
	if err != nil {
		return err
	}

	it := NewItem(o, ttl)
	it.MarshData, _ = json.Marshal(it)
	if err = Clone(o, obj); err != nil {
		return err
	}

	return c.TraverseCache(deep, func(c Cache) error {
		return c.Set(key, it)
	})
}

//GetMutex get a mutex , one key one mutex
func (c *MultiLevelCache) GetMutex(key string) *sync.RWMutex {
	var mux *sync.RWMutex
	nMux := new(sync.RWMutex)
	if oMux, ok := c.muxm.LoadOrStore(key, nMux); ok {
		mux = oMux.(*sync.RWMutex)
		nMux = nil
	} else {
		mux = nMux
	}
	return mux
}

//ReleaseMutex release mutex from map
func (c *MultiLevelCache) ReleaseMutex(key string) {
	var mux *sync.RWMutex
	if oMux, ok := c.muxm.Load(key); ok {
		mux = oMux.(*sync.RWMutex)
	} else {
		return
	}
	mux.Lock()
	defer mux.Unlock()
	c.muxm.Delete(key)
	return
}

// notify all cache nodes to delete key
func (c *MultiLevelCache) Delete(key string) error {
	return RedisPublish(c.Options.DeleteChannel, key, c.pool)
}

// redis subscriber for key deletion
func (c *MultiLevelCache) subscribe(key string) error {
	conn := c.pool.Get()
	defer conn.Close()

	psc := redis.PubSubConn{Conn: conn}
	if err := psc.Subscribe(key); err != nil {
		return err
	}

	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			key := string(v.Data)
			c.delete(key)
		case error:
			return v
		}
	}
}
func (c *MultiLevelCache) delete(key string) error {
	return c.TraverseCache(-1, func(c Cache) error {
		return c.Delete(key)
	})
}

//遍历缓存
//@untilIndex 遍历直到序号 -1：遍历所有 n:遍历指定数量
func (c *MultiLevelCache) TraverseCache(untilIndex int, do func(c Cache) error) error {
	cacheLink := c.mlCache
	if untilIndex < 0 {
		untilIndex = 999
	}
	for i := 0; i <= untilIndex || i < 0; i++ {
		cache := cacheLink.Current
		err := do(cache)
		if err != nil {
			return err
		}
		if cacheLink.Next == nil {
			break
		}
		cacheLink = cacheLink.Next
	}
	return nil
}

//注册缓存
func (c *MultiLevelCache) RegisterCache(cache ...Cache) {
	var start = 0
	if len(cache) == 0 {
		return
	}
	if c.mlCache == nil {
		c.mlCache = &MLCache{
			Current: cache[0],
		}
		start += 1
	}
	var cacheLink *MLCache = c.mlCache
	for i := start; i < len(cache); i++ {
		c.registerCache(cacheLink, cache[i])
	}
}
func (c *MultiLevelCache) registerCache(cacheLink *MLCache, cache Cache) error {
	if cacheLink.Next == nil {
		cacheLink.Next = &MLCache{
			Current: cache,
		}
		return nil
	}
	return c.registerCache(cacheLink.Next, cache)
}

//缓存项是否需要重载
func itemNeedReload(item *Item) bool {
	if item == nil {
		return true
	}
	return item.Expire()
}
