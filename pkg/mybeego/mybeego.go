package mybeego

import (
	"github.com/astaxie/beego"
	"github.com/tiptok/gocomm/pkg/myrest"
	"net/http"
)

var Handlers []func(http.Handler) http.Handler

func GET(rootpath string, f beego.FilterFunc) *beego.App {
	if len(Handlers) == 0 {
		return beego.Get(rootpath, f)
	}
	return beego.Get(rootpath, myrest.BeeUseMiddleware(f, Handlers...))
}

func POST(rootpath string, f beego.FilterFunc) *beego.App {
	if len(Handlers) == 0 {
		return beego.Post(rootpath, f)
	}
	return beego.Post(rootpath, myrest.BeeUseMiddleware(f, Handlers...))
}

func DELETE(rootpath string, f beego.FilterFunc) *beego.App {
	if len(Handlers) == 0 {
		return beego.Delete(rootpath, f)
	}
	return beego.Delete(rootpath, myrest.BeeUseMiddleware(f, Handlers...))
}

func PUT(rootpath string, f beego.FilterFunc) *beego.App {
	if len(Handlers) == 0 {
		return beego.Put(rootpath, f)
	}
	return beego.Put(rootpath, myrest.BeeUseMiddleware(f, Handlers...))
}

func Use(middle ...func(http.Handler) http.Handler) {
	Handlers = append(Handlers, middle...)
}
