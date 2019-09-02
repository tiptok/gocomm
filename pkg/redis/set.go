package redis

import (
	"errors"

	"github.com/garyburd/redigo/redis"
)

//设置集合
func Sadd(key string, value string, timeout int64) error {
	if len(key) < 1 || len(value) < 1 {
		return errors.New("Invalid argument")
	}
	c := redisPool.Get()
	defer c.Close()
	var err error
	_, err = c.Do("SADD", key, value)
	if err != nil {
		return err
	}
	//设置有效时间
	if 1 == Scard(key) && timeout > 0 {
		_, err := c.Do("EXPIRE", key, timeout)
		if err != nil {
			Del(key)
			return err
		}
	}
	return nil
}

//删除集合中一个元素
func Srem(key string, value string) error {
	if len(key) < 1 || len(value) < 1 {
		return errors.New("Invalid argument")
	}
	c := redisPool.Get()
	defer c.Close()
	var err error
	_, err = c.Do("SREM", key, value)
	if err != nil {
		return err
	}
	return nil
}

//随机获取集合中的1个
func Srandmember(key string) (string, error) {
	if len(key) < 1 {
		return "", errors.New("Invalid argument")
	}
	c := redisPool.Get()
	defer c.Close()
	v, err := redis.String(c.Do("SRANDMEMBER", key))
	if err != nil {
		return "", err
	}
	return v, nil
}

//获取set集合成员数量
func Scard(key string) int {
	if len(key) < 1 {
		return 0
	}
	c := redisPool.Get()
	defer c.Close()
	v, err := redis.Int(c.Do("SCARD", key))
	if err != nil || v <= 0 {
		return 0
	}
	return v
}

//获取集合中所有元素
func Smembers(key string) ([]string, error) {
	if len(key) < 1 {
		return nil, errors.New("Invalid argument")
	}
	c := redisPool.Get()
	defer c.Close()
	v, err := redis.Strings(c.Do("SMEMBERS", key))
	if err != nil {
		return nil, err
	}
	return v, nil
}

//获取集合中所有元素
func SmembersInt(key string) ([]int, error) {
	if len(key) < 1 {
		return nil, errors.New("Invalid argument")
	}
	c := redisPool.Get()
	defer c.Close()
	v, err := redis.Ints(c.Do("SMEMBERS", key))
	if err != nil {
		return nil, err
	}
	return v, nil
}
