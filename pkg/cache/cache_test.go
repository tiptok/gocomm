package cache

import (
	"fmt"
	"github.com/tiptok/gocomm/config"
	"github.com/tiptok/gocomm/pkg/log"
	"github.com/tiptok/gocomm/pkg/redis"
	"sync"
	"testing"
	"time"
)

type TestStruct struct {
	Id   int
	Name string
}

// this will be called by deepcopy to improves reflect copy performance
func (p TestStruct) DeepCopy() interface{} {
	c := p
	return &c
}

func getStruct(id uint32) (*TestStruct, error) {
	key := GetKey("val", id)
	var v TestStruct
	err := GetObject(key, &v, 100, func() (interface{}, error) {
		return func() (interface{}, error) {
			return searchById(10)
		}()
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &v, nil
}

func TestCache(t *testing.T) {
	redis.Init(config.Redis{Addr: "127.0.0.1:6379", Password: "123456", MaxIdle: 4})
	InitDefault(WithDefaultRedisPool(redis.GetRedisPool()))
	v, e := getStruct(100)
	if e != nil {
		t.Fatal(e)
	}
	log.Info(v)
}

func searchById(id int) (*TestStruct, error) {
	return &TestStruct{
		Id:   id,
		Name: "tip tok",
	}, nil
}

func TestConcurrent(t *testing.T) {
	redis.Init(config.Redis{Addr: "127.0.0.1:6379", Password: "123456", MaxIdle: 4})
	InitDefault(WithDefaultRedisPool(redis.GetRedisPool()))
	t.Log("开始")
	for i := 1; i <= 10; i++ {
		doConcurrent()
		time.Sleep(time.Second * 10)
	}
	t.Log("结束")
}

func doConcurrent() {
	key := GetKey("partner", "order", "statics", 1)
	var wg sync.WaitGroup
	for i := 1; i <= 10; i++ {
		go func(index int) {
			wg.Add(1)
			defer wg.Done()
			var target *TestStruct = &TestStruct{}
			err := GetObject(key, target, 30, newTestStruct)
			fmt.Println(fmt.Sprintf("%v get : %v  error:%v", index, target, err))
		}(i)
	}
	wg.Wait()
}

func newTestStruct() (interface{}, error) {
	value := &TestStruct{
		Id:   int(time.Now().Unix()),
		Name: "hello tip tok",
	}
	fmt.Println("create instance...", value)
	return value, nil
}

func TestExpire(t *testing.T) {
	redis.Init(config.Redis{Addr: "127.0.0.1:6379", Password: "123456", MaxIdle: 4})
	InitDefault(WithDefaultRedisPool(redis.GetRedisPool()))
	key := GetKey("partner", "order", "statics", 1)
	var target *TestStruct = &TestStruct{}
	var ttl = 10
	if err := GetObject(key, target, ttl, newTestStruct); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 11)
	if err := GetObject(key, target, ttl, newTestStruct); err != nil {
		t.Fatal(err)
	}
}

func BenchmarkMultiLevelCacheNew_Get(b *testing.B) {
	redis.Init(config.Redis{Addr: "127.0.0.1:6379", Password: "123456", MaxIdle: 4})
	InitDefault(WithDefaultRedisPool(redis.GetRedisPool()))
	key := GetKey("partner", "order", "statics", 1)
	var target *TestStruct = &TestStruct{}
	var ttl = 10
	for i := 0; i < b.N; i++ {
		if err := GetObject(key, target, ttl, newTestStruct); err != nil {
			b.Fatal(err)
		}
	}
}

func TestMultiLevelCacheNew_Delete(t *testing.T) {
	redis.Init(config.Redis{Addr: "127.0.0.1:6379", Password: "123456", MaxIdle: 4})
	InitDefault(WithDefaultRedisPool(redis.GetRedisPool()))
	key := GetKey("partner", "order", "statics", 1)
	var target *TestStruct = &TestStruct{}
	if err := GetObject(key, target, 100, newTestStruct); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 20)
	if err := Delete(key); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 5)
}
