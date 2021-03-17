package gs

import (
	"encoding/json"
	"github.com/tiptok/gocomm/common"
	"strconv"
	"strings"
)

type (
	ResponseData struct {
		Code int     `json:"code"`
		Msg  string  `json:"msg"`
		Data MapData `json:"data"`
	}
	MapData map[string]interface{}
)

func (data *ResponseData) FindField(field string) (interface{}, bool) {
	return data.Data.FindField(field)
}
func (data *ResponseData) MustFindField(field string) interface{} {
	v, _ := data.Data.FindField(field)
	return v
}
func (data *ResponseData) Int(field string) int {
	return data.Data.Int(field)
}
func (data *ResponseData) String(field string) string {
	return data.Data.String(field)
}
func (data *ResponseData) Float64(field string) float64 {
	return data.Data.Float64(field)
}
func (data *ResponseData) PrintMapDataStruct() string {
	return data.Data.PrintMapStruct()
}

func (data MapData) FindField(field string) (interface{}, bool) {
	return data.findField(field)
}
func (data MapData) MustFindField(field string) interface{} {
	v, _ := data.findField(field)
	return v
}
func (data MapData) Int(field string) int {
	v := data.MustFindField(field)
	vInt, _ := strconv.Atoi(string(v.(json.Number)))
	return vInt
}
func (data MapData) String(field string) string {
	v := data.MustFindField(field)
	return common.AssertString(v)
}
func (data MapData) Float64(field string) float64 {
	v := data.MustFindField(field)
	vFloat, _ := strconv.ParseFloat(string(v.(json.Number)), 10)
	return vFloat
}

// PrintMapStruct 打印 map 的结构
func (m MapData) PrintMapStruct() string {
	if m == nil {
		return ""
	}
	return common.JsonAssertString(m)
}

// FindField find key value from MapData
// eq:key.key1.key2  will find map[key][key1][key2]
func (m MapData) findField(field string) (interface{}, bool) {
	if len(m) == 0 {
		return nil, false
	}
	if len(field) == 0 {
		return nil, false
	}
	fieldChains := strings.Split(field, ".")
	var cur interface{} = m[fieldChains[0]]
	for i := 1; i < len(fieldChains); i++ {
		mapFiled, ok := cur.(map[string]interface{})
		if !ok {
			return nil, false
		}
		cur = mapFiled[fieldChains[i]]
	}
	return cur, true
}
