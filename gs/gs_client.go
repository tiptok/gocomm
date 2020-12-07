package gs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/tiptok/gocomm/common"
	"strconv"
	"strings"
)

type (
	GatewayClient struct {
		BaseUrl      string
		GlobalHeader map[string]interface{}
		//httplib.BeegoHTTPRequest
		MapApi map[string]ApiImplement
	}
	ApiImplement interface {
		Path(...interface{}) string
		Method() string
	}
	ClientOptions struct {
		PathParam  []interface{}
		PathQuery  map[string]interface{}
		JsonObject interface{}
	}
	HttpRequest struct {
		*httplib.BeegoHTTPRequest
	}
	MapData map[string]interface{}
	Router  struct {
		Key        string
		Path       string
		HttpMethod string
	}
)

func NewGateWayClient(baseUrl string) *GatewayClient {
	return &GatewayClient{
		BaseUrl: baseUrl,
		MapApi:  make(map[string]ApiImplement),
	}
}
func (c *GatewayClient) WithGlobalHeader(header map[string]interface{}) *GatewayClient {
	c.GlobalHeader = header
	return c
}
func (c *GatewayClient) NewRequest(key string, option ...ClientOption) *HttpRequest {
	api := c.MapApi[key]
	options := NewClientOptions(option...)
	rawUrl := c.BaseUrl + api.Path()
	if len(options.PathParam) > 0 {
		rawUrl = c.BaseUrl + api.Path(options.PathParam...)
	}
	request := httplib.NewBeegoRequest(rawUrl, api.Method())
	if options.JsonObject != nil {
		request.JSONBody(options.JsonObject)
	}
	if len(options.PathQuery) > 0 {
		for k, v := range options.PathQuery {
			request.Param(k, common.AssertString(v))
		}
	}
	if len(c.GlobalHeader) != 0 {
		for k, v := range c.GlobalHeader {
			request.Header(k, fmt.Sprintf("%v", v))
		}
	}
	return &HttpRequest{request}
}
func (c *GatewayClient) AddApi(key, path, method string) {
	c.MapApi[key] = apiStruct{
		PathFormat: path,
		HttpMethod: method,
	}
}

type apiStruct struct {
	PathFormat string
	HttpMethod string
}

func (api apiStruct) Path(args ...interface{}) string {
	if len(args) == 0 {
		return api.PathFormat
	}
	return fmt.Sprintf(api.PathFormat, args...)
}
func (api apiStruct) Method() string {
	return api.HttpMethod
}

type ClientOption func(options *ClientOptions)

func WithPathParam(params ...interface{}) func(options *ClientOptions) {
	return func(options *ClientOptions) {
		options.PathParam = params
	}
}
func WithJsonObject(object interface{}) func(options *ClientOptions) {
	return func(options *ClientOptions) {
		options.JsonObject = object
	}
}
func WithPathQuery(pathQuery map[string]interface{}) func(options *ClientOptions) {
	return func(options *ClientOptions) {
		options.PathQuery = pathQuery
	}
}
func NewClientOptions(options ...ClientOption) *ClientOptions {
	option := &ClientOptions{}
	for i := range options {
		options[i](option)
	}
	return option
}

func (b *HttpRequest) ToJSON(v interface{}) error {
	data, err := b.Bytes()
	if err != nil {
		return err
	}
	parse := json.NewDecoder(bytes.NewBuffer(data))
	parse.UseNumber()
	return parse.Decode(v)
}

// FindField find key value from MapData
// eq:key.key1.key2  will find map[key][key1][key2]
func (m MapData) FindField(field string) (interface{}, bool) {
	var result = false
	if len(field) == 0 {
		return m, result
	}
	fieldChains := strings.Split(field, ".")
	var cur interface{} = m[fieldChains[0]]
	for i := 1; i < len(fieldChains); i++ {
		mapFiled, ok := cur.(map[string]interface{})
		if !ok {
			return nil, result
		}
		cur = mapFiled[fieldChains[i]]
	}
	return cur, true
}

// 打印 map 的结构
func (m MapData) PrintMapStruct() string {
	if m == nil {
		return ""
	}
	return common.JsonAssertString(m)
}

func (m MapData) Int(field string) int {
	v, _ := m.FindField(field)
	vInt, _ := strconv.Atoi(string(v.(json.Number)))
	return vInt
}
func (m MapData) String(field string) string {
	v, _ := m.FindField(field)
	return common.AssertString(v)
}
