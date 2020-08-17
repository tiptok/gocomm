package mybeego

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"testing"
)

func TestExample(t *testing.T) {
	initSentinel(t)
	Example()
}

func Example() {
	beego.InsertFilter("/*", beego.BeforeExec, SentinelMiddleware(
		WithBlockFallback(
			func(ctx *context.Context) {
				ctx.Output.SetStatus(400)
				ctx.Output.Body([]byte("too many request; the quota used up"))
				//ctx.Abort(400,"too many request; the quota used up")
			})),
	)
	beego.Get("/test", func(ctx *context.Context) {
		ctx.WriteString("test")
	})
	beego.Get("/work", func(ctx *context.Context) {
		ctx.WriteString("work")
	})
	beego.Run(":8089")
}

func initSentinel(t *testing.T) {
	err := sentinel.InitDefault()
	if err != nil {
		t.Fatalf("Unexpected error: %+v", err)
	}

	_, err = flow.LoadRules([]*flow.FlowRule{
		{
			Resource:        "GET:/test",
			MetricType:      flow.QPS,
			Count:           100,
			ControlBehavior: flow.Reject,
		},
		{
			Resource:        "GET:/work",
			MetricType:      flow.QPS,
			Count:           0,
			ControlBehavior: flow.Reject,
		},
	})
	if err != nil {
		t.Fatalf("Unexpected error: %+v", err)
		return
	}
}
