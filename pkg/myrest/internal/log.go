package internal

import (
	"fmt"
	"github.com/tiptok/gocomm/pkg/log"
	"github.com/tiptok/gocomm/pkg/myrest/httpx"
	"net/http"
)

func Errorf(r *http.Request, format string, args ...interface{}) {
	log.Info(formatWithReq(r, fmt.Sprintf(format, args...)))
}

func formatWithReq(r *http.Request, v string) string {
	return fmt.Sprintf("[http] (%s - %s) %s", r.RequestURI, httpx.GetRemoteAddr(r), v)
}
