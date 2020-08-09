package cache

import (
	"github.com/garyburd/redigo/redis"
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
