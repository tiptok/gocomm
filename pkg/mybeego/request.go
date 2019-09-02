package mybeego

import (
	"fmt"
	"sync/atomic"
)

type RequestHead struct {
	Token    string // 登录令牌
	Uid      int64  // 用户id
	AppId    int    // APP唯一标志
	Version  string // 客户端版本
	Os       string // 手机系统版本
	From     string // 请求来源
	Screen   string // 屏幕尺寸
	Model    string // 机型信息
	Channel  string // 渠道信息
	Net      string // 当前网络状态
	DeviceId string // 设备Id
	LoginIp  string // 登录IP
	Jwt      string // jwt

	requestId string //请求编号 md5
	reqIndex int64   //请求链序号
	//lastOpTime int64 //保存上一次操作请求时间戳，暂时未使用(计算链路耗时)
}
func (reqHead *RequestHead)SetRequestId(addString ...string){
	if (len(addString)==0){
		return
	}
	reqHead.requestId = addString[0]
}
func(reqHead *RequestHead)GetRequestId()string{
	atomic.AddInt64(&reqHead.reqIndex,1)
	return fmt.Sprintf("%s.%d",reqHead.requestId,reqHead.reqIndex)
}


type TotalSwitchStr struct {
	TotalSwitch int    `json:"total_switch"` // 总开关:0-on; 1-off
	MessageBody string `json:"message_body"` // 消息提示信息
}