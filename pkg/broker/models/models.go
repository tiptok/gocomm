package models

type Message struct {
	Id           int64  `json:"id"`
	Topic        string `json:"topic"`
	Value        string `json:"value"`
	MsgTime      int64  `json:"msg_time"`
	FinishStatus int    `json:"-"` //0:未完成 2：已完成 1：进行中  3：忽略
}

//结束状态
type FinishStatus int

const (
	UnFinished FinishStatus = 0
	InProcess  FinishStatus = 1
	Finished   FinishStatus = 2
	Ignore     FinishStatus = 3
)

type MessagePublishResult struct {
	SuccessMessageIds []int64
	ErrorMessageIds   []int64
}
