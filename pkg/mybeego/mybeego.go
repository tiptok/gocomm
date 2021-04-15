package mybeego

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/tiptok/gocomm/pkg/myrest"
	"net/http"
)

var Handlers []func(http.Handler) http.Handler

func GET(rootpath string, f web.FilterFunc) *web.HttpServer {
	if len(Handlers) == 0 {
		return web.Get(rootpath, f)
	}
	return web.Get(rootpath, myrest.BeeUseMiddleware(f, Handlers...))
}

func POST(rootpath string, f web.FilterFunc) *web.HttpServer {
	if len(Handlers) == 0 {
		return web.Post(rootpath, f)
	}
	return web.Post(rootpath, myrest.BeeUseMiddleware(f, Handlers...))
}

func DELETE(rootpath string, f web.FilterFunc) *web.HttpServer {
	if len(Handlers) == 0 {
		return web.Delete(rootpath, f)
	}
	return web.Delete(rootpath, myrest.BeeUseMiddleware(f, Handlers...))
}

func PUT(rootpath string, f web.FilterFunc) *web.HttpServer {
	if len(Handlers) == 0 {
		return web.Put(rootpath, f)
	}
	return web.Put(rootpath, myrest.BeeUseMiddleware(f, Handlers...))
}

func Use(middle ...func(http.Handler) http.Handler) {
	Handlers = append(Handlers, middle...)
}
