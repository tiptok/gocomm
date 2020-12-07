package gs

type GatewayService struct {
	c *GatewayClient
}

func (svr *GatewayService) Invoke(methodKey string, option ...ClientOption) (responseData *ResponseData, err error) {
	request := svr.c.NewRequest(methodKey, option...)
	err = request.ToJSON(&responseData)
	if err != nil {
		return
	}
	return
}

func NewManagerService(baseUrl string, routers []Router) *GatewayService {
	svr := &GatewayService{
		c: NewGateWayClient(baseUrl),
	}
	for _, r := range routers {
		svr.c.AddApi(r.Key, r.Path, r.HttpMethod)
	}
	return svr
}
