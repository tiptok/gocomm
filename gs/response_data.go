package gs

import (
	"encoding/json"
	"github.com/tiptok/gocomm/common"
	"strconv"
)

type ResponseData struct {
	Code int     `json:"code"`
	Msg  string  `json:"msg"`
	Data MapData `json:"data"`
}

func (data *ResponseData) FindField(field string) (interface{}, bool) {
	if data.Data == nil {
		return nil, false
	}
	return data.Data.FindField(field)
}
func (data *ResponseData) MustFindField(field string) interface{} {
	if data.Data == nil {
		return nil
	}
	v, _ := data.Data.FindField(field)
	return v
}
func (data *ResponseData) Int(field string) int {
	v := data.MustFindField(field)
	vInt, _ := strconv.Atoi(string(v.(json.Number)))
	return vInt
}
func (data *ResponseData) String(field string) string {
	v := data.MustFindField(field)
	return common.AssertString(v)
}
func (data *ResponseData) PrintMapDataStruct() string {
	return data.Data.PrintMapStruct()
}
