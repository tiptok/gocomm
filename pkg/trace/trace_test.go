package trace

import (
	"context"
	"github.com/tiptok/gocomm/pkg/log"
	"github.com/tiptok/gocomm/pkg/trace/tracespec"
	"github.com/tiptok/gocomm/xstring"
	sysx "github.com/tiptok/gocomm/xsys"
	"net/http"
	"testing"
)

func TracingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		carrier, err := Extract(HttpFormat, r.Header)
		// ErrInvalidCarrier means no trace id was set in http header
		if err != nil && err != ErrInvalidCarrier {
			log.Error(err)
		}

		ctx, span := StartServerSpan(r.Context(), carrier, sysx.Hostname(), r.RequestURI)
		defer span.Finish()
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func ContextFolk(c context.Context) tracespec.Trace {
	c1, trace := StartClientSpan(c, xstring.RandId(), "folk")
	trace.Finish()
	c1.Value(tracespec.TracingKey)
	return trace
}
func ContextFlow(ctx context.Context, serviceName, operationName string) (context.Context, tracespec.Trace) {
	if span, ok := ctx.Value(tracespec.TracingKey).(*Span); ok {
		return span.Follow(ctx, serviceName, operationName)
	}

	return ctx, emptyNoopSpan
}

func TestTrace(t *testing.T) {
	printTraceInfo := func(o interface{}) {
		span := o.(*Span)
		t.Log(span.TraceId(), span.SpanId(), span.serviceName, span.operationName, span.startTime.Unix())
	}
	c, trace1 := StartServerSpan(context.Background(), nil, "myservices", "init")
	printTraceInfo(trace1)
	c2, trace2 := ContextFlow(c, "db", "db_init")
	printTraceInfo(trace2)
	printTraceInfo(ContextFolk(c2))
	printTraceInfo(ContextFolk(c2))
	printTraceInfo(ContextFolk(c2))
	c3, trace3 := ContextFlow(c2, "log", "log_init")
	printTraceInfo(trace3)
	printTraceInfo(ContextFolk(c3))
	printTraceInfo(ContextFolk(c3))
}
