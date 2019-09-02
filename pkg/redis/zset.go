package redis

import (
	"errors"
	"github.com/garyburd/redigo/redis"
)

//设置集合
func Zadd(key string, score float64, member interface{}) error {
	if len(key) < 1 {
		return errors.New("Invalid argument")
	}
	c := redisPool.Get()
	defer c.Close()
	var err error
	_, err = c.Do("ZADD", key, score, member)
	if err != nil {
		return err
	}
	return nil
}

//有序集合中对指定成员的分数加上增量 increment
func Zincrby(key string, increment int64, member interface{}) error {
	if len(key) < 1 {
		return errors.New("Invalid argument")
	}
	c := redisPool.Get()
	defer c.Close()
	var err error
	_, err = c.Do("ZINCRBY", key, increment, member)
	if err != nil {
		return err
	}
	return nil
}

// 有序集合中对指定成员的分数加上增量 increment
// 注意：次函数只能获取member是整形的情况，如遇到member不少整形的情况需要另外函数
func Zrevrange(key string, start, stop int64) ([]string, error) {
	if len(key) < 1 {
		return nil, errors.New("Invalid argument")
	}
	c := redisPool.Get()
	defer c.Close()
	datas, err := redis.Strings(c.Do("ZREVRANGE", key, start, stop, "WITHSCORES"))
	if err != nil {
		return nil, err
	}
	return datas, nil
}
