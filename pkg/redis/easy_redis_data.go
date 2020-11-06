package redis

type IEasyRedisData interface {
	// 添加
	Add(interface{}) error
	// 获取
	Get(param ...interface{}) (interface{}, error)
	// 移除
	Remove() error
	// 键值
	RedisKey() string
	// 键值
	//Field() string
}

type BaseEasyRedisData struct{}

func (o *BaseEasyRedisData) Add(interface{}) error {
	return nil
}
func (o *BaseEasyRedisData) Get(param ...interface{}) (interface{}, error) {
	return nil, nil
}
func (o *BaseEasyRedisData) Remove(interface{}) error {
	return nil
}
func (o *BaseEasyRedisData) RedisKey() string {
	return ""
}
