package cache

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	. "github.com/tiptok/gocomm/pkg/cache/model"
)

type RedisCache struct {
	pool *redis.Pool
}

func NewRedisCache(p *redis.Pool) *RedisCache {
	c := &RedisCache{
		pool: p,
	}
	return c
}

// read item from redis
func (c *RedisCache) Get(key string, obj interface{}) (*Item, error) {
	body, err := RedisGetString(key, c.pool)
	if err != nil && err != redis.ErrNil {
		return nil, err
	}
	if body == "" {
		return nil, nil //fmt.Errorf("redis cache body is empty'")
	}
	var it Item
	it.Object = obj
	err = json.Unmarshal([]byte(body), &it)
	if err != nil {
		return nil, err
	}
	//fmt.Println("rediscache get cache item:",key,obj)
	return &it, nil
}

func (c *RedisCache) Set(k string, it *Item) error {
	//fmt.Println("rediscache set cache item:",it.String())
	return RedisSetString(k, string(it.Data()), it.TTL*4, c.pool)
}

func (c *RedisCache) Delete(keyPattern string) error {
	//fmt.Println("rediscache delete key:",keyPattern)
	return RedisDelKey(keyPattern, c.pool)
}
