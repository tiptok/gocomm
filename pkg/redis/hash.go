package redis

import (
	"errors"
	"github.com/garyburd/redigo/redis"
)

//设置指定hash指定key的值
func Hset(key string, field string, value interface{}, timeout int64) error {
	if len(key) < 1 || len(field) < 1 {
		return errors.New("Invalid argument")
	}
	c := redisPool.Get()
	defer c.Close()
	var err error
	_, err = c.Do("HSET", key, field, value)
	if err != nil {
		return err
	}
	//设置有效时间
	if timeout > 0 {
		length := Hlen(key)
		if 1 == length {
			_, err := c.Do("EXPIRE", key, timeout)
			if err != nil {
				Del(key)
				return err
			}
		}
	}
	return nil
}

//获取指定hash的所有key
func Hkeys(key string) ([][]byte, error) {
	if len(key) < 1 {
		return nil, errors.New("Invalid argument")
	}
	c := redisPool.Get()
	defer c.Close()
	v, err := redis.ByteSlices(c.Do("HKEYS", key))
	if err != nil {
		return nil, err
	}
	return v, nil
}

//获取指定hash指定key的value
func Hget(key string, field string) (string, error) {
	if len(key) < 1 {
		return "", errors.New("invalid argument")
	}
	c := redisPool.Get()
	defer c.Close()
	v, err := redis.String(c.Do("HGET", key, field))
	if err != nil {
		return "", err
	}
	return v, nil
}

//获取指定hash的key和value
func Hgetall(key string) (map[string]string, error) {
	if len(key) < 1 {
		return nil, errors.New("Invalid argument")
	}
	c := redisPool.Get()
	defer c.Close()
	v, err := redis.StringMap(c.Do("HGETALL", key))
	if err != nil {
		return nil, err
	}
	return v, nil
}

//获取hash字段数量
func Hlen(key string) int {
	if len(key) < 1 {
		return 0
	}
	c := redisPool.Get()
	defer c.Close()
	v, err := redis.Int(c.Do("HLEN", key))
	if err != nil || v <= 0 {
		return 0
	}
	return v
}

// 删除哈希指定字段
func Hdel(key string, field string) error {
	if len(key) < 1 {
		return errors.New("Invalid argument")
	}
	c := redisPool.Get()
	defer c.Close()
	var err error
	_, err = c.Do("HDEL", key, field)
	if err != nil {
		return err
	}
	return nil
}

// 查看哈希表 key 中，指定的字段是否存在
func Hexists(key string, field string) bool {
	if len(key) < 1 || len(field) < 1 {
		return false
	}
	cli := redisPool.Get()
	defer cli.Close()
	v, err := redis.Int(cli.Do("HEXISTS", key, field))
	if err != nil || v == 0 {
		return false
	}
	return true
}
