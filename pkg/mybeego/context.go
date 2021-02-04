package mybeego

import (
	"encoding/json"
	"github.com/astaxie/beego/context"
	"github.com/tiptok/gocomm/common"
	"strconv"
)

type ContextController struct {
}

func (controller ContextController) JsonUnmarshal(ctx *context.Context, v interface{}) error {
	body := ctx.Input.RequestBody
	if len(body) == 0 {
		body = []byte("{}")
	}
	return common.Unmarshal(body, v)
}

func (controller ContextController) BodyKeys(ctx *context.Context, firstCaseToUpper bool) []string {
	var bodyKV map[string]json.RawMessage
	controller.JsonUnmarshal(ctx, &bodyKV)
	if len(bodyKV) == 0 {
		return []string{}
	}
	var list []string
	for k, _ := range bodyKV {
		list = append(list, common.CamelCase(k, true))
	}
	return list
}

func (controller *ContextController) Resp(ctx *context.Context, msg interface{}) {
	ctx.Output.JSON(msg, false, false)
	ctx.Input.SetData("outputData", msg)
}

func (controller ContextController) GetLimitInfo(ctx *context.Context) (offset int, limit int) {
	offset, _ = strconv.Atoi(ctx.Input.Query("pageNumber"))
	limit, _ = strconv.Atoi(ctx.Input.Query("limit"))
	if offset > 0 {
		offset = (offset - 1) * limit
	}
	return
}
