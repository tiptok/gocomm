package trace

//import (
//	"fmt"
//	"github.com/tiptok/gocomm/common"
//	"net/http"
//
//	"github.com/opentracing-contrib/go-stdlib/nethttp"
//	"github.com/opentracing/opentracing-go"
//	"github.com/uber/jaeger-client-go"
//	jaegercfg "github.com/uber/jaeger-client-go/config"
//	"github.com/uber/jaeger-lib/metrics"
//	"github.com/uber/jaeger-lib/metrics/metricstest"
//
//	"github.com/tiptok/gocomm/pkg/log"
//)
//
//func Init(serviceName, addr string) (opentracing.Tracer, error) {
//	// Sample configuration for testing. Use constant sampling to sample every trace
//	// and enable LogSpan to log every span via configured Logger.
//	cfg := jaegercfg.Configuration{
//		Sampler: &jaegercfg.SamplerConfig{
//			Type:  jaeger.SamplerTypeConst,
//			Param: 1,
//		},
//		Reporter: &jaegercfg.ReporterConfig{
//			LogSpans: true,
//		},
//	}
//
//	cfg.ServiceName = serviceName
//
//	// Example logger and metrics factory. Use github.com/uber/jaeger-client-go/log
//	// and github.com/uber/jaeger-lib/metrics respectively to bind to real logging and metrics
//	// frameworks.
//	jLogger := &jaegerLogger{}
//	jMetricsFactory := metrics.NullFactory
//
//	metricsFactory := metricstest.NewFactory(0)
//	metrics := jaeger.NewMetrics(metricsFactory, nil)
//
//	sender, err := jaeger.NewUDPTransport(addr, 0)
//	if err != nil {
//		log.Error("could not initialize jaeger sender:", err.Error())
//		return nil, err
//	}
//
//	repoter := jaeger.NewRemoteReporter(sender, jaeger.ReporterOptions.Metrics(metrics))
//
//	tracer, _, err := cfg.NewTracer(
//		jaegercfg.Logger(jLogger),
//		jaegercfg.Metrics(jMetricsFactory),
//		jaegercfg.Reporter(repoter),
//	)
//	if err != nil {
//		return nil, fmt.Errorf("new trace error: %v", err)
//	}
//
//	return tracer, nil
//
//}
//
//type jaegerLogger struct{}
//
//func (l *jaegerLogger) Error(msg string) {
//	log.Error(common.LogF("ERROR: %s", msg))
//}
//
//// Infof logs a message at info priority
//func (l *jaegerLogger) Infof(msg string, args ...interface{}) {
//	log.Info(common.LogF(msg,args))
//}
//
//func TracingMiddleware(handler http.Handler) http.Handler {
//	return nethttp.Middleware(
//		opentracing.GlobalTracer(),
//		handler,
//		nethttp.MWSpanObserver(func(span opentracing.Span, r *http.Request) {
//
//		}),
//		nethttp.OperationNameFunc(func(r *http.Request) string {
//			return "HTTP " + r.Method + " " + r.RequestURI
//		}),
//	)
//}
