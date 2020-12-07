package handler

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const target = "http://localhost"

func TestRecoverHandler(t *testing.T) {
	Recover := RecoverHandler()
	handler := Recover(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		panic("throw panic:serve error")
	}))
	req := httptest.NewRequest(http.MethodGet, target, nil)
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}
