package handler

import (
	"github.com/tiptok/gocomm/pkg/myrest/internal"
	"net/http"
	"runtime/debug"
)

func RecoverHandler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if p := recover(); p != nil {
					internal.Errorf(r, "%v\n%s", p, debug.Stack())
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
