package gs

import (
	"errors"
	"fmt"
	"github.com/tiptok/gocomm/common"
	"github.com/tiptok/gocomm/pkg/log"
	"net/http"
	"net/http/httputil"
)

var ApplicationError = errors.New("application service error")

type GatewayService struct {
	c          *GatewayClient
	debugModel bool
}

func (svr *GatewayService) Invoke(methodKey string, option ...ClientOption) (*ResponseData, error) {
	var (
		response     *http.Response
		err          error
		responseData *ResponseData
		request      *HttpRequest
	)
	if request, err = svr.c.NewRequest(methodKey, option...); err != nil {
		return nil, err
	}
	if response, err = request.Response(); err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if svr.debugModel {
		dump, _ := httputil.DumpRequest(request.GetRequest(), true)
		rspDump, _ := httputil.DumpResponse(response, true)
		log.Debug(fmt.Sprintf("【HttpRequest】 \n%v%v", string(dump), string(rspDump)))
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status error code:%v", response.StatusCode)
	}
	if err = common.UnmarshalFromReader(response.Body, &responseData); err != nil {
		return responseData, err
	}
	if err == nil && len(responseData.Msg) > 0 && responseData.Code > 0 {
		return responseData, ApplicationError
	}
	return responseData, nil
}

// 启用Debug日志  在生产环境上需要关闭！
func (svr *GatewayService) WithDebugModel(debugModel bool) *GatewayService {
	svr.debugModel = true
	return svr
}

func NewManagerService(baseUrl string, routers []Router, globalOption ...ClientOption) *GatewayService {
	svr := &GatewayService{
		c: NewGateWayClient(baseUrl, globalOption...),
	}
	for _, r := range routers {
		svr.c.AddApi(r.Key, r.Path, r.HttpMethod)
	}
	return svr
}
