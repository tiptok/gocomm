package redis

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/tiptok/gocomm/config"
	"time"
)

const (
	INFINITE int64 = 0
	SECOND   int64 = 1
	MINUTE   int64 = 60
	HOUR     int64 = 3600
	DAY      int64 = 24 * HOUR
	WEEK     int64 = 7 * DAY
	MONTH    int64 = 30 * DAY
	YEAR     int64 = 365 * DAY
)

var (
	// 连接池
	redisPool *redis.Pool
)

func InitWithDb(size int, addr, password, db string) error {
	redisPool = &redis.Pool{
		MaxIdle:     size,
		IdleTimeout: 180 * time.Second,
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		Dial: func() (redis.Conn, error) {
			return dialWithDB(addr, password, db)
		},
	}

	_, err := ping()
	return err
}

func Init(conf config.Redis) error {
	redisPool = &redis.Pool{
		MaxIdle:     conf.MaxIdle,
		IdleTimeout: 180 * time.Second,
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		Dial: func() (redis.Conn, error) {
			return dial(conf.Addr, conf.Password)
		},
	}
	_, err := ping()
	return err
}

func ping() (bool, error) {
	c := redisPool.Get()
	defer c.Close()
	data, err := c.Do("PING")
	if err != nil || data == nil {
		return false, err
	}
	return (data == "PONG"), nil
}

func dial(addr, password string) (redis.Conn, error) {
	c, err := redis.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	if password != "" {
		if _, err := c.Do("AUTH", password); err != nil {
			c.Close()
			return nil, err
		}
	}
	return c, err
}

func dialWithDB(addr, password, db string) (redis.Conn, error) {
	c, err := dial(addr, password)
	if err != nil {
		return nil, err
	}
	if _, err := c.Do("SELECT", db); err != nil {
		c.Close()
		return nil, err
	}
	return c, err
}

//判断键是否存在
func Exists(key string) (bool, error) {
	if len(key) <= 0 {
		return false, errors.New("Empty key")
	}
	c := redisPool.Get()
	defer c.Close()
	exists, err := redis.Bool(c.Do("EXISTS", key))
	return exists, err
}

//删除指定键
func Del(key string) (bool, error) {
	if len(key) <= 0 {
		return false, errors.New("Empty key")
	}
	c := redisPool.Get()
	defer c.Close()
	return redis.Bool(c.Do("DEL", key))
}

//批量删除指定键
func DelMulti(key ...interface{}) (bool, error) {
	if len(key) <= 0 {
		return false, errors.New("Invalid argument")
	}
	c := redisPool.Get()
	defer c.Close()
	return redis.Bool(c.Do("DEL", key...))
}

// func LikeDeletes(key string) error {
// 	conn := RedisConn.Get()
// 	defer conn.Close()

// 	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
// 	if err != nil {
// 		return err
// 	}

// 	for _, key := range keys {
// 		_, err = Delete(key)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

//设置制定键的生存周期
func Expire(key string, timeout int64) error {
	if len(key) <= 0 || timeout < 0 {
		return fmt.Errorf("Invalid argument: key=[%s] timeout=[%d]", key, timeout)
	}
	c := redisPool.Get()
	defer c.Close()
	if timeout == 0 {
		return nil
	}
	_, err := c.Do("EXPIRE", key, timeout)
	return err
}

func setExpire(c redis.Conn, key string, timeout int64) error {
	if len(key) <= 0 || timeout < 0 {
		return fmt.Errorf("Invalid argument: key=[%s] timeout=[%d]", key, timeout)
	}
	_, err := c.Do("EXPIRE", key, timeout)
	return err
}

func DelPattern(pattern string) bool {
	if len(pattern) <= 0 {
		return false
	}
	c := redisPool.Get()
	defer c.Close()
	result, err := redis.Strings(c.Do("Keys", pattern))
	if err != nil {
		return false
	}
	// TODO:pipeline
	for i := range result {
		Del(result[i])
	}
	return true
}

func GetKeyPattern(pattern string) ([]string, error) {
	if len(pattern) < 1 {
		return nil, errors.New("Invalid argument")
	}
	c := redisPool.Get()
	defer c.Close()
	result, err := redis.Strings(c.Do("Keys", pattern))
	if err != nil {
		return nil, err
	}
	return result, nil
}
