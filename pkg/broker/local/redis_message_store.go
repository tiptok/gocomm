package local

import (
	"github.com/tal-tech/go-zero/core/stores/redis"
	"github.com/tiptok/gocomm/common"
	"github.com/tiptok/gocomm/pkg/broker/models"
	"sync/atomic"
	"time"
)

type RedisMessageStore struct {
	start      int64
	redis      *redis.Redis
	normalKey  string
	errorKey   string
	serverFlag string // 服务标识
}

func (store *RedisMessageStore) GetMessage() ([]*models.RetryMessage, error) {
	var result []*models.RetryMessage
	if store.start == 0 {
		store.start = time.Now().Unix()
		return result, nil
	}
	val, err := store.redis.ZrangebyscoreWithScores(store.normalKey, 0, store.start)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(val); i++ {
		var item *models.RetryMessage
		err = common.UnmarshalFromString(val[i].Key, &item)
		if err != nil {
			continue
		}
		result = append(result, item)
	}
	store.redis.Zremrangebyscore(store.normalKey, 0, store.start)
	atomic.CompareAndSwapInt64(&store.start, store.start, time.Now().Unix())
	return result, nil
}

func (store *RedisMessageStore) StoreMessage(msg *models.RetryMessage) error {
	if msg.MaxRetryTime <= msg.RetryTime {
		_, err := store.redis.Zadd(store.errorKey, msg.NextRetryTime, common.JsonAssertString(msg))
		return err
	}
	_, err := store.redis.Zadd(store.normalKey, msg.NextRetryTime, common.JsonAssertString(msg))
	return err
}

func NewRedisMessageStore(serverFlag string, redisAddr, redisPass string) *RedisMessageStore {
	return &RedisMessageStore{
		redis:      redis.NewRedis(redisAddr, redis.NodeType, redisPass),
		serverFlag: serverFlag,
		normalKey:  "consume:message_retry:" + serverFlag,
		errorKey:   "consume:message_retry_error:" + serverFlag,
	}
}
