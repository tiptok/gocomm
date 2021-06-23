package cache

const (
	DefaultObjectExpire = 60 * 60 * 24
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
	QueryOption  func(options *QueryOptions) *QueryOptions
	cacheKeyFunc func() string
	primaryKeyFunc func(obj interface{})string
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

func (c *CachedRepository) QueryUniqueIndexCache(keyFunc cacheKeyFunc,v interface{},queryPrimaryKeyFunc LoadFunc, queryFunc LoadFunc, options ...QueryOption) error {
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
	// 通过 queryPrimaryKeyFunc 先查primaryKey,再用primaryKey查询缓存记录
	var primaryKey string
	if err:=c.mlCache.GetObject(key, &primaryKey, -1, queryPrimaryKeyFunc);err!=nil{
		return err
	}
	return c.mlCache.GetObject(primaryKey, v, option.ObjectToExpire, queryFunc)
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
