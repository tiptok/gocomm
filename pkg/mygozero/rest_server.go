package mygozero

import (
	"github.com/tal-tech/go-zero/rest"
	"net/http"
)

type ServerRouters struct {
	RestServer *rest.Server
	Routers    []rest.Route
}

func (server *ServerRouters) GET(relativePath string, handler http.HandlerFunc) {
	server.AddRoute(rest.Route{Method: http.MethodGet, Path: relativePath, Handler: handler})
}

func (server *ServerRouters) POST(relativePath string, handler http.HandlerFunc) {
	server.AddRoute(rest.Route{Method: http.MethodPost, Path: relativePath, Handler: handler})
}

func (server *ServerRouters) PUT(relativePath string, handler http.HandlerFunc) {
	server.AddRoute(rest.Route{Method: http.MethodPut, Path: relativePath, Handler: handler})
}

func (server *ServerRouters) DELETE(relativePath string, handler http.HandlerFunc) {
	server.AddRoute(rest.Route{Method: http.MethodDelete, Path: relativePath, Handler: handler})
}

func (server *ServerRouters) AddRoute(router rest.Route) {
	server.Routers = append(server.Routers, router)
}
