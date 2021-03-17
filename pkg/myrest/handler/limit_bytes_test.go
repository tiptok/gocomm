package handler

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLimitBytesHandler(t *testing.T) {
	maxb := LimitBytesHandler(10)
	handler := maxb(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	req := httptest.NewRequest(http.MethodPost, "http://localhost", bytes.NewBufferString("12345678901"))
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusRequestEntityTooLarge, resp.Code)

	req = httptest.NewRequest(http.MethodPost, "http://localhost", bytes.NewBufferString("1234567890"))
	resp = httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestLimitBytesHandlerNoLimit(t *testing.T) {
	maxb := LimitBytesHandler(-1)
	handler := maxb(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	req := httptest.NewRequest(http.MethodPost, "http://localhost", bytes.NewBufferString("12345678901"))
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}
