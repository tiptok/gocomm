package mybeego

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"testing"
)

func TestExample(t *testing.T) {
	initSentinel(t)
	Example()
}

func Example() {
	web.InsertFilter("/*", web.BeforeExec, SentinelMiddleware(
		WithBlockFallback(
			func(ctx *context.Context) {
				ctx.Output.SetStatus(400)
				ctx.Output.Body([]byte("too many request; the quota used up"))
				//ctx.Abort(400,"too many request; the quota used up")
			})),
	)
	web.Get("/test", func(ctx *context.Context) {
		ctx.WriteString("test")
	})
	web.Get("/work", func(ctx *context.Context) {
		ctx.WriteString("work")
	})
	web.Run(":8089")
}

func initSentinel(t *testing.T) {
	err := sentinel.InitDefault()
	if err != nil {
		t.Fatalf("Unexpected error: %+v", err)
	}

	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource: "GET:/test",
			//Type:      flow.QPS,
			//Count:           100,
			Threshold:       100,
			ControlBehavior: flow.Reject,
		},
		{
			Resource: "GET:/work",
			//MetricType:      flow.QPS,
			//Count:           0,
			Threshold:       100,
			ControlBehavior: flow.Reject,
		},
	})
	if err != nil {
		t.Fatalf("Unexpected error: %+v", err)
		return
	}
}
