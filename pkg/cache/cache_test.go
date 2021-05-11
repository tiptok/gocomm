package cache

import (
	"github.com/stretchr/testify/assert"
	"github.com/tiptok/gocomm/config"
	"github.com/tiptok/gocomm/pkg/cache/gzcache"
	memory "github.com/tiptok/gocomm/pkg/cache/memory"
	redis_cache "github.com/tiptok/gocomm/pkg/cache/redis_cache"
	"github.com/tiptok/gocomm/pkg/redis"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

/*script
go test -bench=. -v pkg/cache/cache_test.go pkg/cache/cache.go pkg/cache/option.go pkg/cache/utils.go pkg/cache/multi_level_cache.go
*/
const (
	MemoryRedisCache = iota + 1
	RedisCacheFlag
	MemoryCache
	GzCache
)

const CacheType = GzCache

var (
	count       int32
	redisHost   = "127.0.0.1:6379"
	redisPasswd = ""
)

func initInstance() {
	switch CacheType {
	case MemoryRedisCache:
		redis.Init(config.Redis{Addr: redisHost, Password: redisPasswd, MaxIdle: 4})
		InitDefault(WithDefaultRedisPool(redis.GetRedisPool()))
	case RedisCacheFlag:
		redis.Init(config.Redis{Addr: redisHost, Password: redisPasswd, MaxIdle: 4})
		mlCache = NewMultiLevelCacheNew(WithDefaultRedisPool(redis.GetRedisPool()))
		mlCache.RegisterCache(
			redis_cache.NewRedisCache(mlCache.Options.DefaultRedisPool),
		)
	case MemoryCache:
		redis.Init(config.Redis{Addr: redisHost, Password: redisPasswd, MaxIdle: 4})
		mlCache = NewMultiLevelCacheNew(WithDefaultRedisPool(redis.GetRedisPool()))
		mlCache.RegisterCache(
			memory.NewMemCache(mlCache.Options.CleanInterval),
		)
	case GzCache:
		InitMultiLevelCache().
			RegisterCache(gzcache.NewNodeCache(redisHost, redisPasswd))
	default:
	}
	count = 0
}

// 测试并行情况
func TestCacheGetConcurrent(t *testing.T) {
	initInstance()
	for i := 1; i <= 10; i++ {
		doConcurrent(t)
	}
}

func doConcurrent(t *testing.T) {
	key := GetKey("partner", "order", "statics", 1)
	var wg sync.WaitGroup
	for i := 1; i <= 1000; i++ {
		go func(index int) {
			wg.Add(1)
			defer wg.Done()
			var target *TestStruct = &TestStruct{}
			GetObject(key, target, 30, newTestStruct)
			//fmt.Println(fmt.Sprintf("%v get : %v  error:%v", index, target, err))
			assert.Equal(t, target.Name, "tiptok")
		}(i)
	}
	wg.Wait()
}

// 测试过期
func TestCacheExpire(t *testing.T) {
	initInstance()
	key := GetKey("partner", "order", "statics", time.Now().Unix())
	var target *TestStruct = &TestStruct{}
	var ttl = 1
	var count int32
	createInstance := func() (interface{}, error) {
		value := &TestStruct{
			Id:   int(time.Now().Unix()),
			Name: "tiptok",
		}
		atomic.AddInt32(&count, 1)
		return value, nil
	}
	if err := GetObject(key, target, ttl, createInstance); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second * 1)
	if err := GetObject(key, target, ttl, createInstance); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, int32(2), count)
}

// 测试默认 内存+redis缓存 获取缓存对象
func TestCacheGetExpire(t *testing.T) {
	initInstance()
	key := GetKey("partner", "order", "statics", time.Now().Unix())
	var target *TestStruct = &TestStruct{}
	if err := GetObject(key, target, 1, newTestStruct); err != nil {
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
	time.Sleep(1 * time.Second)
	if !item.Expire() {
		t.Fatal("get object has expire")
	}
}

// 测试默认 内存+redis缓存 删除缓存对象
func TestCacheDelete(t *testing.T) {
	initInstance()
	key := GetKey("partner", "order", "statics", time.Now().Unix())
	var target *TestStruct = &TestStruct{}
	if err := GetObject(key, target, 100, newTestStruct); err != nil {
		t.Fatal(err)
	}
	if err := Delete(key); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Millisecond * 50)
	GetObject(key, target, 100, newTestStruct)
	assert.Equal(t, int32(2), count)
}

// 测试redis缓存
func TestCacheGet(t *testing.T) {
	initInstance()
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

// Benchmark MultiLevelCache  Memory + Redis
func BenchmarkCacheGet(b *testing.B) {
	initInstance()
	key := GetKey("partner", "order", "statics", time.Now().Unix())
	var target *TestStruct = &TestStruct{}
	var ttl = 10
	for i := 0; i < b.N; i++ {
		if err := GetObject(key, target, ttl, newTestStruct); err != nil {
			assert.Error(b, err)
		}
		//assert.Equal(b,"tiptok",target.Name)
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
		return nil, err
	}
	return &v, nil
}
func newTestStruct() (interface{}, error) {
	value := &TestStruct{
		Id:   int(time.Now().Unix()),
		Name: "tiptok",
	}
	atomic.AddInt32(&count, 1)
	//fmt.Println("create instance...", value)
	return value, nil
}
