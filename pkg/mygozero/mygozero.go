package mygozero

import (
	"github.com/zeromicro/go-zero/rest"
	"net/http"
)

var ServerRouter = &ServerRouters{}

func GET(relativePath string, handler http.HandlerFunc) {
	ServerRouter.AddRoute(rest.Route{Method: http.MethodGet, Path: relativePath, Handler: handler})
}

func POST(relativePath string, handler http.HandlerFunc) {
	ServerRouter.AddRoute(rest.Route{Method: http.MethodPost, Path: relativePath, Handler: handler})
}

func PUT(relativePath string, handler http.HandlerFunc) {
	ServerRouter.AddRoute(rest.Route{Method: http.MethodPut, Path: relativePath, Handler: handler})
}

func DELETE(relativePath string, handler http.HandlerFunc) {
	ServerRouter.AddRoute(rest.Route{Method: http.MethodDelete, Path: relativePath, Handler: handler})
}

func AddRoute(router rest.Route) {
	ServerRouter.Routers = append(ServerRouter.Routers, router)
}
