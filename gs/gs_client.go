package gs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/tiptok/gocomm/common"
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
		Header     map[string]interface{}
	}
	HttpRequest struct {
		*httplib.BeegoHTTPRequest
	}
	Router struct {
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
	if len(options.Header) > 0 {
		for k, v := range options.Header {
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

// WithPathParam 请求路径参数填充  eg:/user/:userId/info ,需要填充:userId值
func WithPathParam(params ...interface{}) func(options *ClientOptions) {
	return func(options *ClientOptions) {
		options.PathParam = params
	}
}

// WithJsonObject 参数body
func WithJsonObject(object interface{}) func(options *ClientOptions) {
	return func(options *ClientOptions) {
		options.JsonObject = object
	}
}

// WithPathQuery 请求路径后面的参数 eg:?id=1&&name=22
func WithPathQuery(pathQuery map[string]interface{}) func(options *ClientOptions) {
	return func(options *ClientOptions) {
		options.PathQuery = pathQuery
	}
}

// WithHeader 请求头 eg:x_trace_id=1
func WithHeader(header map[string]interface{}) func(options *ClientOptions) {
	return func(options *ClientOptions) {
		options.Header = header
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
