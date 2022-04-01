package cache

const (
	DefaultObjectExpire = 60 * 60 * 24
	//cacheSafeGapBetweenIndexAndPrimary = 5 // sec
)

type (
	CachedRepository struct {
		mlCache *MultiLevelCache
		option  *QueryOptions
	}
	QueryOptions struct {
		ObjectToExpire int
		NoCacheFlag    bool
	}
	QueryOption    func(options *QueryOptions) *QueryOptions
	cacheKeyFunc   func() string
	primaryKeyFunc func(obj interface{}) string
)

func NewCachedRepository(c *MultiLevelCache, options ...QueryOption) *CachedRepository {
	option := NewQueryOptions(options...)
	return &CachedRepository{
		mlCache: c,
		option:  option,
	}
}

func NewDefaultCachedRepository() *CachedRepository {
	return NewCachedRepository(mlCache)
}

func (c *CachedRepository) QueryCache(keyFunc cacheKeyFunc, v interface{}, queryFunc LoadFunc, options ...QueryOption) error {
	option := NewQueryOptions(options...)
	key := keyFunc()
	if option.NoCacheFlag || len(key) == 0 {
		if object, err := queryFunc(); err != nil {
			return err
		} else {
			Clone(object, v)
		}
		return nil
	}
	return c.mlCache.GetObject(key, v, option.ObjectToExpire, queryFunc)
}

func (c *CachedRepository) Query(queryFunc LoadFunc, deleteKeys ...string) (interface{}, error) {
	var ret interface{}
	var err error
	if ret, err = queryFunc(); err != nil {
		return ret, err
	}
	for _, key := range deleteKeys {
		if len(key) == 0 {
			continue
		}
		if err = c.mlCache.Delete(key); err != nil {
			return ret, err
		}
	}
	return ret, err
}

// QueryUniqueIndexCache 根据唯一索引查询缓存数据
//
// uniqueIndexKeyFunc 唯一索引键值函数
// v                  查询对象
// primaryKeyFunc     主键索引键值
// queryFunc          查询函数
func (c *CachedRepository) QueryUniqueIndexCache(uniqueIndexKeyFunc cacheKeyFunc, v interface{}, primaryKeyFunc keyFunc, queryFunc LoadFunc, options ...QueryOption) error {
	option := NewQueryOptions(options...)
	key := uniqueIndexKeyFunc()
	if option.NoCacheFlag || len(key) == 0 {
		if object, err := queryFunc(); err != nil {
			return err
		} else {
			Clone(object, v)
		}
		return nil
	}

	var primaryKey string
	ok, err := c.mlCache.GetCacheWithoutLoad(key, &primaryKey)
	if err != nil {
		return err
	}
	if ok && len(primaryKey) > 0 {
		return c.mlCache.GetObject(primaryKey, v, -1, queryFunc)
	}

	var primaryObj interface{}
	queryPrimaryKeyFunc := func() (interface{}, error) {
		primaryObj, err = queryFunc()
		if err != nil {
			return nil, err
		}
		return primaryKeyFunc(primaryObj), err
	}
	err = c.mlCache.Load(key, &primaryKey, -1, queryPrimaryKeyFunc, -1) // ttl可改为可配置的，默认不过期，如果有ttl 需要注意primary 、 unique index 之前的过期时间
	if err != nil {
		return err
	}
	return c.mlCache.GetObject(primaryKeyFunc(primaryObj), v, option.ObjectToExpire, func() (interface{}, error) {
		return primaryObj, nil
	})
}

func WithNoCacheFlag() QueryOption {
	return func(options *QueryOptions) *QueryOptions {
		options.NoCacheFlag = true
		return options
	}
}
func WithObjectToExpire(expire int) QueryOption {
	return func(options *QueryOptions) *QueryOptions {
		options.ObjectToExpire = expire
		return options
	}
}
func NewQueryOptions(options ...QueryOption) *QueryOptions {
	option := new(QueryOptions)
	option.ObjectToExpire = DefaultObjectExpire
	for i := range options {
		options[i](option)
	}
	return option
}
