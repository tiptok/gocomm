package cache

const (
	DefaultObjectExpire = 60 * 2 //60 * 60 * 24 *7
)

type (
	CachedRepository struct {
		mlCache *MultiLevelCache
	}
	QueryOptions struct {
		ObjectToExpire int
		NoCacheFlag    bool
	}
	QueryOption func(options *QueryOptions) *QueryOptions
)

func NewCachedRepository(c *MultiLevelCache) *CachedRepository {
	return &CachedRepository{
		mlCache: c,
	}
}

func NewDefaultCachedRepository() *CachedRepository {
	return NewCachedRepository(mlCache)
}

func (c *CachedRepository) QueryCache(key string, v interface{}, queryFunc LoadFunc, options ...QueryOption) error {
	option := NewQueryOptions(options...)
	if option.NoCacheFlag || len(key) == 0 {
		if object, err := queryFunc(); err != nil {
			return err
		} else {
			Clone(object, v)
		}
		return nil
	}
	return c.mlCache.GetObject(key, v, DefaultObjectExpire, queryFunc)
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

func WithNoCacheFlag() QueryOption {
	return func(options *QueryOptions) *QueryOptions {
		options.NoCacheFlag = true
		return options
	}
}
func NewQueryOptions(options ...QueryOption) *QueryOptions {
	option := new(QueryOptions)
	for i := range options {
		options[i](option)
	}
	return option
}
