package trace

import (
	"context"
	"github.com/tiptok/gocomm/pkg/trace/tracespec"
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
