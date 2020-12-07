package trace

//import (
//	"github.com/opentracing/opentracing-go"
//	"github.com/opentracing/opentracing-go/ext"
//	"net/http"
//)
//
//func TracingHTTPRequest(tracer opentracing.Tracer,tracerName string,tagValue interface{}) (func(next http.Handler) http.Handler) {
//	return func(next http.Handler) (http.Handler) {
//		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
//			// Try to join to a trace propagated in `req`.
//			//步骤1 解客户端span
//			wireContext, err := tracer.Extract(
//				opentracing.TextMap,
//				opentracing.HTTPHeadersCarrier(req.Header),
//			)
//			if err!=nil{
//				panic(err)
//			}
//			//步骤2 启动服务端span
//			span := tracer.StartSpan(tracerName, ext.RPCServerOption(wireContext))
//			span.SetTag("server", tagValue)
//			//部署4 关闭span
//			defer span.Finish()
//			// 部署3 store span in context
//			ctx := opentracing.ContextWithSpan(req.Context(), span)
//			// update request context to include our new span
//			req = req.WithContext(ctx)
//			// next middleware or actual request handler
//			next.ServeHTTP(w, req)
//		})
//	}
//}
