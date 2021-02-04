package mybeego

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/tiptok/gocomm/common"
)

type BaseController struct {
	beego.Controller
}

func (controller BaseController) JsonUnmarshal(v interface{}) error {
	body := controller.Ctx.Input.RequestBody
	if len(body) == 0 {
		body = []byte("{}")
	}
	return common.Unmarshal(body, v)
}

func (controller BaseController) BodyKeys(firstCaseToUpper bool) []string {
	var bodyKV map[string]json.RawMessage
	controller.JsonUnmarshal(&bodyKV)
	if len(bodyKV) == 0 {
		return []string{}
	}
	var list []string
	for k, _ := range bodyKV {
		list = append(list, common.CamelCase(k, true))
	}
	return list
}

func (controller *BaseController) Resp(msg interface{}) {
	controller.Data["json"] = msg
	controller.Ctx.Input.SetData("outputData", msg)
	controller.ServeJSON()
}

func (controller BaseController) GetLimitInfo() (offset int, limit int) {
	offset, _ = controller.GetInt("pageNumber")
	limit, _ = controller.GetInt("limit")
	if offset > 0 {
		offset = (offset - 1) * limit
	}
	return
}

//获取请求头信息
//func (controller *BaseController) GetRequestHeader(ctx *context.Context) *protocol.RequestHeader {
//	h := &protocol.RequestHeader{}
//
//	if v := ctx.Input.GetData("x-mmm-id"); v != nil {
//		h.UserId = int64(v.(int))
//	}
//	if v := ctx.Input.GetData("x-mmm-uname"); v != nil {
//		h.UserName = v.(string)
//	}
//	h.Token = ctx.Input.Header("Authorization")
//	if len(h.Token) > 0 && len(strings.Split(h.Token, " ")) > 1 {
//		h.Token = strings.Split(h.Token, " ")[1]
//	}
//	h.BodyKeys = controller.BodyKeys(true)
//	if v := ctx.Request.URL.Query(); len(v) > 0 {
//		for k, _ := range v {
//			h.BodyKeys = append(h.BodyKeys, common.CamelCase(k, true))
//		}
//	}
//	return h
//}
