package mybeego

import (
	"github.com/beego/beego/v2/server/web"
	"net/http"
)

type BeegoEngine struct {
	routers    []Router
	middleWare []func(http.Handler) http.Handler
}

type Router struct {
	Method  string
	Path    string
	Handler web.FilterFunc
}

func (engine *BeegoEngine) Run(addr string) {
	server := web.NewHttpSever()
	server.Run(addr)
}

func (engine *BeegoEngine) InitRouters() {

}

func (engine *BeegoEngine) Use(middle ...func(http.Handler) http.Handler) {
	Handlers = append(Handlers, middle...)
}

func (engine *BeegoEngine) Routers() []Router {
	return engine.routers
}

func NewBeegoEngine() *BeegoEngine {
	return &BeegoEngine{}
}
