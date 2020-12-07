package gs

import (
	"net/http"
	"testing"
)

func TestNewManagerService(t *testing.T) {
	GlobalRouter := []Router{
		{UserPost, "/user", http.MethodPost},
		{UserPut, "/user/%v", http.MethodPut},
		{UserGet, "/user/%v", http.MethodGet},
		{UserDelete, "/user/%v", http.MethodDelete},
		{UserList, "/user", http.MethodGet},

		{AuthLogin, "/auth/login", http.MethodPost},
	}

	svr := NewManagerService("http://127.0.0.1:8080/v1", GlobalRouter)
	data, err := svr.Invoke(UserPut, WithPathParam(map[string]interface{}{"id": 1}))
	if err != nil {
		t.Log(data.PrintMapDataStruct())
	}
}
