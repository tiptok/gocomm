package redis

import (
	"encoding/json"
	"errors"

	"github.com/garyburd/redigo/redis"
)

func Set(key string, v interface{}, timeout int64) error {
	if len(key) <= 0 || timeout < 0 || v == nil {
		err := errors.New("Invalid argument")
		return err
	}
	c := redisPool.Get()
	defer c.Close()
	switch v := v.(type) {
	case int8, int16, int32, int, int64, uint8, uint16, uint, uint32, uint64, string:
		if timeout == 0 {
			if _, err := c.Do("SET", key, v); err != nil {
				return err
			}
		} else {
			if _, err := c.Do("SETEX", key, timeout, v); err != nil {
				return err
			}
		}
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		if timeout == 0 {
			if _, err = c.Do("SET", key, string(b)); err != nil {
				return err
			}
		} else {
			if _, err = c.Do("SETEX", key, timeout, string(b)); err != nil {
				return err
			}
		}

	}
	return nil
}

func Get(key string) (string, error) {
	if len(key) < 1 {
		return "", errors.New("Invalid argument")
	}
	c := redisPool.Get()
	defer c.Close()
	v, err := redis.String(c.Do("GET", key))
	if err != nil {
		return "", err
	}

	return v, nil
}

//获取指定键的INCR
func Incr(key string, timeout int64) (int64, bool) {
	if len(key) < 1 {
		return 0, false
	}
	c := redisPool.Get()
	defer c.Close()
	var isExpire bool = false
	// timeout大于0并且不存在改key，则需要设置ttl
	exists, err := Exists(key)
	if err != nil {
		return 0, false
	}
	if timeout > INFINITE && !exists {
		isExpire = true
	}
	v, err := redis.Int64(c.Do("INCR", key))
	if err != nil || v == 0 {
		return 0, false
	}
	if isExpire {
		Expire(key, timeout)
	}
	return v, true
}
