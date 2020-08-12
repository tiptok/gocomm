package mybeego

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"net/http"
)

//Sentinel 拦截器（限流,熔断）
func SentinelMiddleware(opts ...Option) beego.FilterFunc {
	options := evaluateOptions(opts)
	return func(ctx *context.Context) {
		resourceName := ctx.Request.Method + ":" + ctx.Request.RequestURI
		if options.resourceExtract != nil {
			resourceName = options.resourceExtract(ctx)
		}
		entry, err := sentinel.Entry(
			resourceName,
			sentinel.WithResourceType(base.ResTypeWeb),
			sentinel.WithTrafficType(base.Inbound),
		)
		if err != nil {
			if options.blockFallback != nil {
				options.blockFallback(ctx)
			} else {
				ctx.Abort(http.StatusTooManyRequests, "")
			}
			return
		}

		defer entry.Exit()
		//ctx
	}
}

type (
	Option  func(*options)
	options struct {
		resourceExtract func(*context.Context) string
		blockFallback   func(*context.Context)
	}
)

func evaluateOptions(opts []Option) *options {
	optCopy := &options{}
	for _, opt := range opts {
		opt(optCopy)
	}

	return optCopy
}

// WithResourceExtractor sets the resource extractor of the web requests.
func WithResourceExtractor(fn func(*context.Context) string) Option {
	return func(opts *options) {
		opts.resourceExtract = fn
	}
}

// WithBlockFallback sets the fallback handler when requests are blocked.
func WithBlockFallback(fn func(ctx *context.Context)) Option {
	return func(opts *options) {
		opts.blockFallback = fn
	}
}
