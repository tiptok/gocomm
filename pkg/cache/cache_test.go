package cache

import (
	"fmt"
	"github.com/tiptok/gocomm/config"
	memory "github.com/tiptok/gocomm/pkg/cache/memory"
	redis_cache "github.com/tiptok/gocomm/pkg/cache/redis_cache"
	"github.com/tiptok/gocomm/pkg/log"
	"github.com/tiptok/gocomm/pkg/redis"
	"sync"
	"testing"
	"time"
)

// 测试缓存
func TestCache(t *testing.T) {
	redis.Init(config.Redis{Addr: "127.0.0.1:6379", Password: "123456", MaxIdle: 4})
	InitDefault(WithDefaultRedisPool(redis.GetRedisPool()))
	v, e := getStruct(100)
	if e != nil {
		t.Fatal(e)
	}
	log.Info(v)
}

// 测试并行情况
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

// 测试过期
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

// 测试默认 内存+redis缓存 获取缓存对象
func TestMultiLevelCache_GetObject(t *testing.T) {
	redis.Init(config.Redis{Addr: "127.0.0.1:6379", Password: "123456", MaxIdle: 4})
	InitDefault(WithDefaultRedisPool(redis.GetRedisPool()))
	key := GetKey("partner", "order", "statics", 1)
	var target *TestStruct = &TestStruct{}
	if err := GetObject(key, target, 2, newTestStruct); err != nil {
		t.Fatal(err)
	}
	var target2 = &TestStruct{}
	item, err := mlCache.mlCache.Current.Get(key, target2)
	if err != nil {
		t.Fatal("get object error")
	}
	if item.Expire() {
		t.Fatal("get object expire")
	}
	time.Sleep(2 * time.Second)
	if !item.Expire() {
		t.Fatal("get object has expire")
	}
}

// 测试默认 内存+redis缓存 删除缓存对象
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

// 测试redis缓存
func TestMultiLevelCache_Redis(t *testing.T) {
	redis.Init(config.Redis{Addr: "127.0.0.1:6379", Password: "123456", MaxIdle: 4})
	mlCache = NewMultiLevelCacheNew(WithDefaultRedisPool(redis.GetRedisPool()))
	mlCache.RegisterCache(
		redis_cache.NewRedisCache(mlCache.Options.DefaultRedisPool),
	)
	key := GetKey("partner", "order", "statics", 1)
	var target *TestStruct = &TestStruct{}
	if err := mlCache.GetObject(key, target, 100, newTestStruct); err != nil {
		t.Fatal(err)
	}
	var target2 = &TestStruct{}
	item, err := mlCache.mlCache.Current.Get(key, target2)
	if err != nil {
		t.Fatal("redis get object error")
	}
	if item.Expire() {
		t.Fatal("redis get object expire")
	}
}

// 测试内存缓存
func TestMultiLevelCache_Memory(t *testing.T) {
	redis.Init(config.Redis{Addr: "127.0.0.1:6379", Password: "123456", MaxIdle: 4})
	mlCache = NewMultiLevelCacheNew(WithDefaultRedisPool(redis.GetRedisPool()))
	mlCache.RegisterCache(
		memory.NewMemCache(mlCache.Options.CleanInterval),
	)
	key := GetKey("partner", "order", "statics", 1)
	var target *TestStruct = &TestStruct{}
	if err := mlCache.GetObject(key, target, 100, newTestStruct); err != nil {
		t.Fatal(err)
	}
	var target2 = &TestStruct{}
	item, err := mlCache.mlCache.Current.Get(key, target2)
	if err != nil {
		t.Fatal("memory get object error")
	}
	if item.Expire() {
		t.Fatal("memory get object expire")
	}
}

// 批量测试
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

// 缓存对象
type TestStruct struct {
	Id   int
	Name string
}

// this will be called by deepcopy to improves reflect copy performance
func (p TestStruct) DeepCopy() interface{} {
	c := p
	return &c
}
func searchById(id int) (*TestStruct, error) {
	return &TestStruct{
		Id:   id,
		Name: "tip tok",
	}, nil
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
func newTestStruct() (interface{}, error) {
	value := &TestStruct{
		Id:   int(time.Now().Unix()),
		Name: "hello tip tok",
	}
	fmt.Println("create instance...", value)
	return value, nil
}
