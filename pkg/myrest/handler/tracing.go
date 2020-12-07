package handler

import (
	"github.com/tiptok/gocomm/pkg/log"
	"github.com/tiptok/gocomm/pkg/trace"
	xsys "github.com/tiptok/gocomm/xsys"
	"net/http"
)

func TracingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		carrier, err := trace.Extract(trace.HttpFormat, r.Header)
		// ErrInvalidCarrier means no trace id was set in http header
		if err != nil && err != trace.ErrInvalidCarrier {
			log.Error(err)
		}

		ctx, span := trace.StartServerSpan(r.Context(), carrier, xsys.Hostname(), r.RequestURI)
		defer span.Finish()
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
