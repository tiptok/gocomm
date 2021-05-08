package handler

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"github.com/tiptok/gocomm/common"
	"github.com/tiptok/gocomm/pkg/cache/model"
	"github.com/tiptok/gocomm/pkg/log"
	"net/http"
	"strings"
)

const (
	apqExtension = "apq"

	defaultExpire = 60 // 60s
)

// AtomicPersistenceHandler  if routers match , atomic persistence response data to cache store,cache will be used in future lookups
// links:https://gqlgen.com/reference/apq/
// links:https://github.com/apollographql/apollo-link-persisted-queries
func AtomicPersistenceQueryHandler(options ...option) func(http.Handler) http.Handler {
	option := NewOptions(options...)
	option.ValidAPQ()
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var (
				queryHash string
				err       error
			)
			if !(r.Method == http.MethodPost || r.Method == http.MethodGet) {
				next.ServeHTTP(w, r)
				return
			}
			if !checkRouter(r, option.routers) {
				next.ServeHTTP(w, r)
				return
			}
			if queryHash, err = ComputeHttpRequestQueryHash(r); err != nil {
				log.Error(err)
				next.ServeHTTP(w, r)
				return
			}
			var item string
			// if cache is miss , store the newest data to cache
			if v, err := option.cache.Get(redisKey(option.serviceName, queryHash), &item); err != nil || v == nil {
				if err != nil {
					log.Error(err)
				}
				responseBuf := bytes.NewBuffer(nil)
				crw := newCacheResponseWrite(w, responseBuf)
				next.ServeHTTP(crw, r)
				if err := option.cache.Set(redisKey(option.serviceName, queryHash), model.NewItem(responseBuf.String(), option.expire)); err != nil {
					log.Error(err)
				}
				return
			}
			// 此处不能提前设置状态，否则beego内部框架会识别response已被处理,导致content-type:text-plain(一直是)
			// 详见 :https://blog.csdn.net/yes169yes123/article/details/103126655
			// w.WriteHeader(http.StatusAccepted)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Write([]byte(item))
			return
		})
	}
}

// computeQueryHash compute hash key
func ComputeQueryHash(query string) string {
	b := sha256.Sum256([]byte(query))
	return hex.EncodeToString(b[:])
}

//ComputeHttpRequestQueryHash  compute request query hash
func ComputeHttpRequestQueryHash(r *http.Request) (string, error) {
	var queryHash string
	if r.Method == http.MethodGet {
		queryHash = ComputeQueryHash(r.URL.String())
	} else if r.Method == http.MethodPost {
		body, err := common.DumpReadCloser(r.Body)
		if err != nil {
			return "", err
		}
		queryHash = ComputeQueryHash(r.URL.String() + string(body))
	}
	return queryHash, nil
}

func checkRouter(r *http.Request, routers []string) bool {
	for i := range routers {
		if common.KeyMatch3(r.URL.Path, routers[i]) {
			return true
		}
	}
	return false
}

func redisKey(serviceName, hash string) string {
	return strings.Join([]string{serviceName, apqExtension, hash}, ":")
}

// cacheResponseWrite buffer response data in future use
type cacheResponseWrite struct {
	writer http.ResponseWriter
	buf    *bytes.Buffer
}

func (w *cacheResponseWrite) Header() http.Header {
	return w.writer.Header()
}

func (w *cacheResponseWrite) Write(bs []byte) (int, error) {
	w.buf.Write(bs)
	return w.writer.Write(bs)
}

func (w *cacheResponseWrite) WriteHeader(code int) {
	w.writer.WriteHeader(code)
}

func newCacheResponseWrite(writer http.ResponseWriter, buf *bytes.Buffer) *cacheResponseWrite {
	return &cacheResponseWrite{
		writer: writer,
		buf:    buf,
	}
}
