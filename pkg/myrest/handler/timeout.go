package handler

import (
	"net/http"
	"time"
)

const reason = "Request Timeout"

// TODO:gin上面使用有问题，http.ContentType 写入失败（待完善）
func TimeoutHandler(duration time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if duration > 0 {
			return http.TimeoutHandler(next, duration, reason)
		} else {
			return next
		}
	}
}
