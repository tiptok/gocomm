package gs

import (
	"bytes"
	"encoding/json"
	"github.com/tiptok/gocomm/common"
	"net/http"
	"testing"
)

const (
	UserPost   = "UserPost"
	UserPut    = "UserPut"
	UserGet    = "UserGet"
	UserDelete = "UserDelete"
	UserList   = "UserPost"

	AuthLogin = "AuthLogin"
)

var c = NewGateWayClient("http://mmm-godevp-dev.fjmaimaimai.com/v1")

func init() {
	c.WithGlobalHeader(map[string]interface{}{"Authorization": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6IjEiLCJwYXNzd29yZCI6IjdjNGE4ZDA5Y2EzNzYyYWY2MWU1OTUyMDk0M2RjMjY0OTRmODk0MWIiLCJhZGREYXRhIjp7IlVzZXJOYW1lIjoidGlwdG9rIn0sImV4cCI6MTYwNTE1MzQwMiwiaXNzIjoiand0In0.BaGz73D5YSf98jXs-HATO8Ah8Thm415N8UAerlbNt48"})
	c.AddApi(UserPost, "/user", http.MethodPost)
	c.AddApi(UserPut, "/user/%v", http.MethodPut)
	c.AddApi(UserGet, "/user/%v", http.MethodGet)
	c.AddApi(UserDelete, "/user/%v", http.MethodDelete)
	c.AddApi(UserList, "/user", http.MethodGet)

	c.AddApi(AuthLogin, "/auth/login", http.MethodPost)
}

func TestGatewayPost(t *testing.T) {
	request := c.NewRequest(AuthLogin, WithJsonObject(map[string]interface{}{"username": "18860183051", "password": "7c4a8d09ca3762af61e59520943dc26494f8941b"}))
	var responseData *ResponseData
	err := request.ToJSON(&responseData)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("mapstruct:", responseData.PrintMapDataStruct())
	t.Log("expiresIn:", responseData.Int("access.expiresIn"))
	t.Log("accessToken:", responseData.String("access.accessToken"))
}

func TestGatewayGetList(t *testing.T) {
	request := c.NewRequest(UserList, WithPathQuery(map[string]interface{}{"pageNumber": 1, "pageSize": 20}))
	var responseData *ResponseData
	err := request.ToJSON(&responseData)
	if err != nil {
		t.Fatal(err)
	}

	mapRet := map[string]interface{}{
		"total": responseData.MustFindField("gridResult.totalRow"),
		"list":  responseData.MustFindField("gridResult.list"),
	}
	t.Log(common.JsonAssertString(mapRet))
}

func TestGatewayGet(t *testing.T) {
	request := c.NewRequest(UserGet, WithPathParam(map[string]interface{}{"id": 1}))
	var responseData *ResponseData
	err := request.ToJSON(&responseData)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(responseData.PrintMapDataStruct())
	mapRet := map[string]interface{}{
		"user": responseData.MustFindField("user"),
	}
	t.Log(common.JsonAssertString(mapRet))
}

func TestGatewayPut(t *testing.T) {
	request := c.NewRequest(UserPut, WithPathParam(1), WithJsonObject(map[string]interface{}{"name": "tip111"}))
	var responseData *ResponseData
	err := request.ToJSON(&responseData)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(responseData, responseData.PrintMapDataStruct())
}

func TestResponseData(t *testing.T) {
	parse := json.NewDecoder(bytes.NewBuffer([]byte(`{
    "code": 0,
    "msg": "成功",
    "data": {}
}`)))
	parse.UseNumber()
	var responseData *ResponseData
	err := parse.Decode(&responseData)
	if err != nil {
		t.Fatal(err)
	}
}
