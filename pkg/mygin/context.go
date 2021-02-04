package mygin

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/tiptok/gocomm/common"
	"net/http"
	"strconv"
)

type ContextController struct {
}

func (controller ContextController) JsonUnmarshal(ctx *gin.Context, v interface{}) error {
	return ctx.ShouldBind(v)
}

func (controller ContextController) BodyKeys(ctx *gin.Context, firstCaseToUpper bool) []string {
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

func (controller *ContextController) Resp(ctx *gin.Context, msg interface{}) {
	ctx.Set("outputData", msg)
	ctx.JSON(http.StatusOK, msg)
}

func (controller ContextController) GetLimitInfo(ctx *gin.Context) (offset int, limit int) {
	offset, _ = strconv.Atoi(ctx.Query("pageNumber"))
	limit, _ = strconv.Atoi(ctx.Query("limit"))
	if offset > 0 {
		offset = (offset - 1) * limit
	}
	return
}
