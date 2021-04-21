package gs

import (
	"encoding/json"
	"github.com/tiptok/gocomm/common"
	"reflect"
	"strconv"
	"strings"
)

const (
	splitChar = "."
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
func (data *ResponseData) Int64(field string) int64 {
	return data.Data.Int64(field)
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

func NewMapData() MapData {
	m := make(map[string]interface{})
	return m
}
func (data MapData) AddFiled(field string, value interface{}) MapData {
	fields := strings.Split(field, splitChar)
	var cur MapData
	cur = data
	for index, f := range fields {
		if index != (len(fields) - 1) {
			if _, ok := cur[f]; !ok {
				cur[f] = make(map[string]interface{})
			}
			cur = cur[f].(map[string]interface{})
			continue
		}
		if _, ok := cur[f]; !ok {
			cur[f] = value
		}
	}
	return data
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
func (data MapData) Int64(field string) int64 {
	v := data.MustFindField(field)
	vInt, _ := strconv.Atoi(string(v.(json.Number)))
	return int64(vInt)
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
func (data MapData) Bool(field string) bool {
	v := data.MustFindField(field)
	vb, _ := strconv.ParseBool(string(v.(json.Number)))
	return vb
}

// PrintMapStruct 打印 map 的结构
func (data MapData) PrintMapStruct() string {
	if data == nil {
		return ""
	}
	return common.JsonAssertString(data)
}

func (data MapData) GetFiledMap(field string) map[string]interface{} {
	fields := strings.Split(field, splitChar)
	cur := data
	for _, f := range fields {
		if _, ok := cur[f]; !ok {
			cur[f] = make(map[string]interface{})
		}
		cur = cur[f].(map[string]interface{})
	}
	return cur
}
func (data MapData) SetFieldMap(fieldMap map[string]interface{}, field string, value interface{}) MapData {
	if value == nil {
		return data
	}
	v := reflect.ValueOf(value)
	if !v.IsValid() {
		return data
	}
	if v.IsZero() {
		return data
	}
	fieldMap[field] = value
	return data
}

// FindField find key value from MapData
// eq:key.key1.key2  will find map[key][key1][key2]
func (data MapData) findField(field string) (interface{}, bool) {
	if len(data) == 0 {
		return nil, false
	}
	if len(field) == 0 {
		return nil, false
	}
	fieldChains := strings.Split(field, ".")
	var cur interface{} = data[fieldChains[0]]
	for i := 1; i < len(fieldChains); i++ {
		mapFiled, ok := cur.(map[string]interface{})
		if !ok {
			return nil, false
		}
		cur = mapFiled[fieldChains[i]]
	}
	return cur, true
}
