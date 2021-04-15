package mybeego

import (
	"fmt"
	"github.com/tiptok/gocomm/pkg/log"
	"time"
)

const (
	TOTAL_SWITCH_ON  int    = 0 // 通行
	TOTAL_SWITCH_OFF int    = 1 // 关闭，系统停止受理
	SWITCH_INFO_KEY  string = "switch_info"
)

type Message struct {
	Errno   int         `json:"errno"`
	Errmsg  string      `json:"errmsg"`
	SysTime int64       `json:"sys_time"`
	Data    interface{} `json:"data"`
}

var ErrnoMsg map[int]string

//var MessageMap map[int]*Message

func NewMessage(code int) *Message {
	return &Message{
		Errno:   code,
		Errmsg:  ErrnoMsg[code],
		SysTime: time.Now().Unix(),
	}
}

func NewErrMessage(code int, errMsg ...interface{}) *Message {
	defer func() {
		if p := recover(); p != nil {
			log.Error(p)
		}
	}()
	msg := NewMessage(code)
	if len(errMsg) > 1 {
		msg.Data = fmt.Sprintf(errMsg[0].(string), errMsg[1:]...)
	} else if len(errMsg) == 1 {
		msg.Data = errMsg[0].(string)
	} else {
		msg.Data = nil
	}
	return msg
}

func init() {
	// 注：错误码9999消息文本可自定义
	ErrnoMsg = make(map[int]string)
	ErrnoMsg[0] = "成功"

	ErrnoMsg[1] = "系统错误"
	ErrnoMsg[2] = "参数错误"
	ErrnoMsg[3] = "系统升级中"
	ErrnoMsg[4] = "您目前使用的版本过低，无法显示最新的相关内容，请使用响单单APP最新版本。"
	ErrnoMsg[5] = "描述包含敏感词，请重新编辑"
	ErrnoMsg[6] = "重复提交，请稍后再试"
}
