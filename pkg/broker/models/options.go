package models

import "time"

var defaultKafkaHost = "localhost:9092"
var defaultInterval = time.Second * 60 * 5

type MessageOptions struct {
	KafkaHost         string
	Interval          time.Duration
	MessageRepository MessagePublisherRepository
	Version           string
	//默认false: 格式化成models.Message对象，
	//true:需要原始sarama.ConsumeMessage对象
	HandlerOriginalMessageFlag bool
	// Enable  enable consume try retry ,true:enable false:disable
	EnableConsumeRetry bool
	ConsumeRetryOption *ConsumeRetryOption
}

type MessageOption func(options *MessageOptions)

func WithMessageProduceRepository(repository MessagePublisherRepository) MessageOption {
	return func(options *MessageOptions) {
		options.MessageRepository = repository
	}
}

func WithInterval(interval time.Duration) MessageOption {
	return func(options *MessageOptions) {
		options.Interval = interval
	}
}

func WithKafkaHost(kafkaHost string) MessageOption {
	return func(options *MessageOptions) {
		options.KafkaHost = kafkaHost
	}
}

func WithVersion(version string) MessageOption {
	return func(options *MessageOptions) {
		options.Version = version
	}
}

func WithHandlerOriginalMessageFlag(flag bool) MessageOption {
	return func(options *MessageOptions) {
		options.HandlerOriginalMessageFlag = flag
	}
}

// WithConsumeRetryOption set message retry config
// maxRetryTime is max try time,if over limit time ,will abort retry
// retryDuration is interval of timer to work
// store is message persistent container
func WithConsumeRetryOption(maxRetryTime int, retryDuration int, store MessageStore) MessageOption {
	return func(options *MessageOptions) {
		options.ConsumeRetryOption = &ConsumeRetryOption{
			MaxRetryTime:      maxRetryTime,
			NextRetryTimeSpan: retryDuration,
			Store:             store,
		}
		if maxRetryTime > 0 {
			options.EnableConsumeRetry = true
		}
	}
}

func NewMessageOptions(options ...MessageOption) *MessageOptions {
	option := &MessageOptions{
		KafkaHost:                  defaultKafkaHost,
		Interval:                   defaultInterval,
		HandlerOriginalMessageFlag: false,
	}
	for i := range options {
		options[i](option)
	}
	return option
}
