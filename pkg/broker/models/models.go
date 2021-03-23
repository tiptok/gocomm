package models

import "github.com/Shopify/sarama"

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

type ConsumeRetryOption struct {
	// Enable  enable consume try retry ,true:enable false:disable
	Enable bool
	// 最大重试次数
	MaxRetryTime int
	// 下一次重试间隔 单位:second
	NextRetryTimeSpan int
	// 消息仓库
	Store MessageStore
}

type RetryMessage struct {
	Message       *sarama.ConsumerMessage
	RetryTime     int
	NextRetryTime int64
	MaxRetryTime  int
}
