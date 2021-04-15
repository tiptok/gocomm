package cache

import (
	"github.com/garyburd/redigo/redis"
	"github.com/tiptok/gocomm/pkg/log"
	"time"
)

const (
	delKeyChannel = "delkey"
	cleanInterval = time.Second * 10
)

type Options struct {
	CleanInterval    time.Duration
	DefaultRedisPool *redis.Pool
	DeleteChannel    string
	DebugMode        bool
	Log              func() log.Log
}

type Option func(o *Options)

func NewOptions(option ...Option) *Options {
	o := &Options{
		CleanInterval: cleanInterval,
		DeleteChannel: delKeyChannel,
	}
	for i := range option {
		option[i](o)
	}
	return o
}

func WithCleanInterval(i time.Duration) Option {
	return func(o *Options) {
		o.CleanInterval = i
	}
}

func WithDefaultRedisPool(p *redis.Pool) Option {
	return func(o *Options) {
		o.DefaultRedisPool = p
	}
}

func WithDeleteChannel(s string) Option {
	return func(o *Options) {
		o.DeleteChannel = s
	}
}

func WithDebugLog(DebugModule bool, log func() log.Log) Option {
	return func(o *Options) {
		o.DebugMode = DebugModule
		o.Log = log
	}
}
