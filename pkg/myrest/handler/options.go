package handler

import (
	"github.com/tiptok/gocomm/pkg/cache"
	"net/http"
)

type Options struct {
	hashFunc         func(query string) string
	requestQueryHash func(r *http.Request) (string, error)
	serviceName      string
	cache            cache.Cache
	expire           int
	routers          []string
}

type option func(options *Options)

func WithRequestQueryHashFunc(requestQueryHash func(r *http.Request) (string, error)) option {
	return func(options *Options) {
		options.requestQueryHash = requestQueryHash
	}
}

func WithServiceName(serviceName string) option {
	return func(options *Options) {
		options.serviceName = serviceName
	}
}

func WithCache(cache cache.Cache) option {
	return func(options *Options) {
		options.cache = cache
	}
}

// WithExpire set cache expire duration (unit:second)
func WithExpire(expire int) option {
	return func(options *Options) {
		options.expire = expire
	}
}

func WithRouters(routers []string) option {
	return func(options *Options) {
		options.routers = routers
	}
}

func NewOptions(options ...option) *Options {
	option := &Options{
		hashFunc:         ComputeQueryHash,
		requestQueryHash: ComputeHttpRequestQueryHash,
		expire:           defaultExpire,
	}
	for i := range options {
		options[i](option)
	}
	return option
}

func (o *Options) ValidAPQ() error {
	if o.cache == nil {
		panic("Options cache is null")
	}
	if len(o.serviceName) == 0 {
		panic("Options serviceName is empty")
	}
	return nil
}
