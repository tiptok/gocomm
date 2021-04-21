package gs

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/client/httplib"
	"github.com/tiptok/gocomm/common"
	"strings"
	"time"
)

const (
	defaultConnectTimeOut   = time.Second * 5
	defaultReadWriteTimeOut = time.Second * 5
)

type (
	GatewayClient struct {
		BaseUrl       string
		MapApi        map[string]ApiImplement
		GlobalOptions []ClientOption
	}
	ApiImplement interface {
		Path(...interface{}) string
		Method() string
	}
	ClientOptions struct {
		PathParam        []interface{}
		PathQuery        map[string]interface{}
		JsonObject       interface{}
		Header           map[string]interface{}
		TlsConfig        *tls.Config
		connectTimeout   time.Duration
		readWriteTimeout time.Duration
	}
	HttpRequest struct {
		*httplib.BeegoHTTPRequest
	}
	Router struct {
		Key        string
		Path       string
		HttpMethod string
	}
	ClientOption func(options *ClientOptions)
)

// NewGateWayClient make a gateway client
// hostUrl 服务 host 地址 eg:http://example.com/
// globalOptions 全局配置
func NewGateWayClient(hostUrl string, globalOptions ...ClientOption) *GatewayClient {
	return &GatewayClient{
		BaseUrl:       hostUrl,
		MapApi:        make(map[string]ApiImplement),
		GlobalOptions: globalOptions,
	}
}

// NewRequest 新建http请求
func (c *GatewayClient) NewRequest(key string, option ...ClientOption) (*HttpRequest, error) {
	api, ok := c.MapApi[key]
	if !ok {
		return nil, fmt.Errorf("routers %v is unregistered", key)
	}
	if len(c.GlobalOptions) > 0 {
		option = append(option, c.GlobalOptions...)
	}
	options := NewClientOptions(option...)
	rawUrl := c.rawUrl(api, options)
	request := httplib.NewBeegoRequest(rawUrl, api.Method())
	request.SetTimeout(options.connectTimeout, options.readWriteTimeout)
	if strings.HasPrefix(rawUrl, "https") {
		request.SetTLSClientConfig(options.TlsConfig)
	}
	if len(options.PathQuery) > 0 {
		for k, v := range options.PathQuery {
			request.Param(k, common.AssertString(v))
		}
	}
	if len(options.Header) > 0 {
		for k, v := range options.Header {
			request.Header(k, fmt.Sprintf("%v", v))
		}
	}
	if options.JsonObject != nil {
		request.JSONBody(options.JsonObject)
	}
	return &HttpRequest{request}, nil
}

// AddApi 添加路由
func (c *GatewayClient) AddApi(key, path, method string) {
	c.MapApi[key] = apiStruct{
		PathFormat: path,
		HttpMethod: method,
	}
}

// rawUrl
func (c *GatewayClient) rawUrl(api ApiImplement, options *ClientOptions) string {
	rawUrl := c.BaseUrl + api.Path()
	if len(options.PathParam) > 0 {
		rawUrl = c.BaseUrl + api.Path(options.PathParam...)
	}
	return rawUrl
}

type apiStruct struct {
	PathFormat string
	HttpMethod string
}

func (api apiStruct) Path(args ...interface{}) string {
	if len(args) == 0 {
		return api.PathFormat
	}
	if mapArgs, ok := args[0].(map[string]interface{}); ok {
		mapApi := api.PathFormat
		for k, v := range mapArgs {
			old := ":" + k
			if strings.Contains(mapApi, old) {
				mapApi = strings.Replace(mapApi, old, common.AssertString(v), 1)
			}
		}
		return mapApi
	}
	return fmt.Sprintf(api.PathFormat, args...)
}
func (api apiStruct) Method() string {
	return api.HttpMethod
}

func NewClientOptions(options ...ClientOption) *ClientOptions {
	option := &ClientOptions{
		TlsConfig:        &tls.Config{InsecureSkipVerify: true},
		connectTimeout:   defaultConnectTimeOut,
		readWriteTimeout: defaultReadWriteTimeOut,
	}
	for i := range options {
		options[i](option)
	}
	return option
}

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
		if options.Header == nil {
			options.Header = make(map[string]interface{})
		}
		for k, v := range header {
			options.Header[k] = v
		}
	}
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
