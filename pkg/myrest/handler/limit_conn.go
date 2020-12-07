package handler

import (
	"fmt"
	"github.com/tiptok/gocomm/pkg/log"
	"github.com/tiptok/gocomm/pkg/myrest/httpx"
	"github.com/tiptok/gocomm/pkg/myrest/internal"
	"github.com/tiptok/gocomm/sync/limit"
	"net/http"
)

func LimitConnHandler(n int) func(http.Handler) http.Handler {
	if n <= 0 {
		return func(next http.Handler) http.Handler {
			return next
		}
	}

	return func(next http.Handler) http.Handler {
		latchLimiter := limit.NewLimit(n)

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if latchLimiter.TryBorrow() {
				defer func() {
					if err := latchLimiter.Return(); err != nil {
						log.Error(err)
					}
				}()

				next.ServeHTTP(w, r)
			} else {
				internal.Errorf(r, fmt.Sprintf("(%s - %s) Concurrent connections over %d, rejected with code %d",
					r.RequestURI, httpx.GetRemoteAddr(r), n, http.StatusServiceUnavailable))
				w.WriteHeader(http.StatusServiceUnavailable)
			}
		})
	}
}
