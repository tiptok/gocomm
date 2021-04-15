package myrest

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/gin-gonic/gin"
	"github.com/justinas/alice"
	"net/http"
)

type BeegoEngine struct {
}

// HandlerFuncUseMiddleware http.HandlerFunc use middleware to intercept http request
func HandlerFuncUseMiddleware(work http.HandlerFunc, middle ...func(http.Handler) http.Handler) http.Handler {
	chain := midChain(middle...)
	return chain.ThenFunc(work)
}

// BeeUseMiddleware beego.HandlerFilter use middleware o intercept http request
func BeeUseMiddleware(work func(c *context.Context), middle ...func(http.Handler) http.Handler) web.FilterFunc {
	chain := midChain(middle...)
	return func(c *context.Context) {
		svr := chain.ThenFunc(BeeFuncToHandlerFunc(c, work))
		svr.ServeHTTP(c.ResponseWriter, c.Request)
	}
}

// midChain return middleware chains
func midChain(middle ...func(http.Handler) http.Handler) alice.Chain {
	var mid []alice.Constructor
	mid = append(mid)
	for _, v := range middle {
		mid = append(mid, v)
	}
	chain := alice.New(mid...)
	return chain
}

// HandlerToBeeFunc  http.Handler convert to beego.FilterFunc
func HandlerToBeeFunc(h http.Handler) web.FilterFunc {
	return func(c *context.Context) {
		h.ServeHTTP(c.ResponseWriter, c.Request)
	}
}

// HandlerFuncToBeeFunc http.HandlerFunc convert to beego.FilterFunc
func HandlerFuncToBeeFunc(h http.HandlerFunc) web.FilterFunc {
	return func(c *context.Context) {
		h.ServeHTTP(c.ResponseWriter, c.Request)
	}
}

// BeeFuncToHandlerFunc   beego.FilterFunc  convert to  http.HandlerFunc
func BeeFuncToHandlerFunc(c *context.Context, work func(c *context.Context)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.Request = r
		c.ResponseWriter = &context.Response{
			ResponseWriter: w,
		}
		work(c)
	}
}

func GinMiddlewareChains(middle ...func(http.Handler) http.Handler) gin.HandlerFunc {
	chain := midChain(middle...)
	return func(c *gin.Context) {
		svr := chain.ThenFunc(func(http.ResponseWriter, *http.Request) {})
		svr.ServeHTTP(c.Writer, c.Request)
	}
}
