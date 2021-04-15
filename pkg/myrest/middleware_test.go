package myrest

import (
	"fmt"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/tiptok/gocomm/pkg/myrest/handler"
	"github.com/tiptok/gocomm/pkg/myrest/httpx"
	"log"
	"net/http"
	"testing"
	"time"
)

const maxByte = 1024 * 1024
const maxConn = 4

func TestBeegoMiddleware(t *testing.T) {
	work := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello world")
	}
	beeWork := func(ctx *context.Context) {
		ctx.WriteString("hello world")
		var m map[string]interface{}
		httpx.ParseJsonBody(ctx.Request, &m)
		if len(m) > 0 {
			log.Println(m)
		}
	}
	web.Handler("/work", HandlerFuncUseMiddleware(work))
	web.Get("/work2", HandlerToBeeFunc(HandlerFuncUseMiddleware(work)))
	//beego.Get("/v1/work",BeeUseMiddleware(beeWork,
	//	handler.TracingHandler,
	//	handler.LogHandler,
	//	handler.LimitConnHandler(maxConn),
	//	handler.TimeoutHandler(time.Second*5),
	//	handler.RecoverHandler(),
	//	handler.LimitBytesHandler(maxByte),
	//	beeWorkMiddleware("mid1"),beeWorkMiddleware("mid2")),
	//)
	web.Post("/v1/work", bindMiddleware(beeWork))

	web.BConfig.CopyRequestBody = true
	web.Run(":8080")
}

func beeWorkMiddleware(midName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			//log.Println("detect middleware :",midName,time.Now().UnixNano())
			//time.Sleep(time.Millisecond*10)
			next.ServeHTTP(writer, request)
		})
	}
}

func bindMiddleware(work func(c *context.Context)) func(c *context.Context) {
	return BeeUseMiddleware(work,
		handler.TracingHandler,
		handler.LogHandler,
		handler.LimitConnHandler(maxConn),
		handler.TimeoutHandler(time.Second*5),
		handler.RecoverHandler(),
		handler.LimitBytesHandler(maxByte),
		beeWorkMiddleware("mid1"),
		beeWorkMiddleware("mid2"),
	)
}
