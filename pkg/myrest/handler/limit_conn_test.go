package handler

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

const conn = 4

func TestLimitConnHandler(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(conn)
	done := make(chan interface{})
	defer close(done)
	limitConnHandler := LimitConnHandler(conn)
	handler := limitConnHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wg.Done()
		<-done
	}))
	for i := 0; i < conn; i++ {
		go func() {
			request := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
			handler.ServeHTTP(httptest.NewRecorder(), request)
		}()
	}
	wg.Wait()
	request := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, request)
	assert.Equal(t, http.StatusServiceUnavailable, resp.Code)
}
