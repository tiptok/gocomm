package trace

import (
	"context"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/reporter"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/tiptok/gocomm/pkg/trace/tracespec"
	"net/http"
)

func StartClientSpan(ctx context.Context, serviceName, operationName string) (context.Context, tracespec.Trace) {
	if span, ok := ctx.Value(tracespec.TracingKey).(*Span); ok {
		return span.Fork(ctx, serviceName, operationName)
	}

	return ctx, emptyNoopSpan
}

func StartServerSpan(ctx context.Context, carrier Carrier, serviceName, operationName string) (
	context.Context, tracespec.Trace) {
	span := newServerSpan(carrier, serviceName, operationName)
	return context.WithValue(ctx, tracespec.TracingKey, span), span
}

func Extract(format, carrier interface{}) (Carrier, error) {
	switch v := format.(type) {
	case int:
		if v == HttpFormat {
			return emptyHttpPropagator.Extract(carrier)
		} else if v == GrpcFormat {
			return emptyGrpcPropagator.Extract(carrier)
		}
	}

	return nil, ErrInvalidCarrier
}

func Inject(format, carrier interface{}) (Carrier, error) {
	switch v := format.(type) {
	case int:
		if v == HttpFormat {
			return emptyHttpPropagator.Inject(carrier)
		} else if v == GrpcFormat {
			return emptyGrpcPropagator.Inject(carrier)
		}
	}

	return nil, ErrInvalidCarrier
}

func HttpInject(ctx context.Context, request *http.Request) {
	if span, ok := ctx.Value(tracespec.TracingKey).(*Span); ok {
		request.Header.Add(traceIdKey, span.TraceId())
		request.Header.Add(spanIdKey, span.SpanId())
	}
	return
}

/*********************自定义 tracing*************************/

var globalReport reporter.Reporter
var globalLocalEndpoint *model.Endpoint
var globalSpanFinish func(*Span)

func BindReporter(report reporter.Reporter) {
	globalReport = report
}

// BindZipkinReporter 通过zipkin http 进行上报
// eg: url http://127.0.0.1:9411/api/v2/spans
// TODO:可配置上报策略，上报周期等
func BindZipkinReporter(url string) {
	globalReport = zipkinhttp.NewReporter(url)
	globalSpanFinish = zipkinOnFinish
}

func BindEndpoint(serviceName string, hostPort string) error {
	ep, err := zipkin.NewEndpoint(serviceName, hostPort)
	if err == nil {
		globalLocalEndpoint = ep
	}
	return err
}

func BindSpanFinish(onFinish func(span *Span)) {
	globalSpanFinish = onFinish
}
