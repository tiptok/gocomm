package kafkax

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/tiptok/gocomm/pkg/broker/models"
	"log"
	"strings"
	"time"
)

// sarame kafka 消息生产
type KafkaMessageProducer struct {
	KafkaHosts string
	LogInfo    models.LogInfo
}

// 同步发送
func (engine *KafkaMessageProducer) Publish(messages []*models.Message, option map[string]interface{}) (*models.MessagePublishResult, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Retry.Max = 10
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Version = sarama.V0_11_0_0
	brokerList := strings.Split(engine.KafkaHosts, ",")
	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Println(err)
		}
	}()
	var successMessageIds []int64
	var errMessageIds []int64
	for _, message := range messages {
		if value, err := json.Marshal(message); err == nil {
			msg := &sarama.ProducerMessage{
				Topic:     message.Topic,
				Value:     sarama.StringEncoder(value),
				Timestamp: time.Now(),
			}
			partition, offset, err := producer.SendMessage(msg)
			if err != nil {
				errMessageIds = append(errMessageIds, message.Id)
				log.Println(err)
			} else {
				successMessageIds = append(successMessageIds, message.Id)
				var append = make(map[string]interface{})
				append["topic"] = message.Topic
				append["partition"] = partition
				append["offset"] = offset
				log.Println("kafka消息发送", append)
			}
		}
	}
	return &models.MessagePublishResult{SuccessMessageIds: successMessageIds, ErrorMessageIds: errMessageIds}, nil
}

// 消息调度器
type MessageDispatcher struct {
	notifications     chan struct{}
	messageChan       chan *models.Message
	dispatchTicker    *time.Ticker
	messageRepository models.MessageRepository
	producer          models.MessageProducer
}

func (dispatcher *MessageDispatcher) MessagePublishedNotice() error {
	time.Sleep(time.Second * 2)
	dispatcher.notifications <- struct{}{}
	return nil
}

func (dispatcher *MessageDispatcher) MessagePublish(messages []*models.Message) error {
	for i := range messages {
		dispatcher.messageChan <- messages[i]
	}
	return nil
}

// go dispatcher.Dispatch() 启动一个独立协程
func (dispatcher *MessageDispatcher) Dispatch() {
	for {
		select {
		case <-dispatcher.dispatchTicker.C:
			go func(dispatcher *MessageDispatcher) {
				dispatcher.notifications <- struct{}{}
			}(dispatcher)
		case <-dispatcher.notifications:
			if dispatcher.messageRepository == nil {
				continue
			}
			messages, _ := dispatcher.messageRepository.FindNoPublishedStoredMessages()
			var messagesInProcessIds []int64
			for i := range messages {
				messagesInProcessIds = append(messagesInProcessIds, messages[i].Id)
			}
			if messages != nil && len(messages) > 0 {
				dispatcher.messageRepository.FinishMessagesStatus(messagesInProcessIds, int(models.InProcess))

				reuslt, err := dispatcher.producer.Publish(messages, nil)
				if err == nil && len(reuslt.SuccessMessageIds) > 0 {
					dispatcher.messageRepository.FinishMessagesStatus(reuslt.SuccessMessageIds, int(models.Finished))
				}
				//发送失败的消息ID列表 更新状态 进行中->未开始
				if len(reuslt.ErrorMessageIds) > 0 {
					dispatcher.messageRepository.FinishMessagesStatus(reuslt.ErrorMessageIds, int(models.UnFinished))
				}
			}
		case msg := <-dispatcher.messageChan:
			dispatcher.producer.Publish([]*models.Message{msg}, nil)
		}
	}
}

type MessageDirector struct {
	messageRepository models.MessageRepository
	dispatcher        *MessageDispatcher
}

func (d *MessageDirector) PublishMessages(messages []*models.Message) error {
	if d.dispatcher == nil {
		return fmt.Errorf("dispatcher还没有启动")
	}
	if d.messageRepository == nil {
		d.dispatcher.MessagePublish(messages)
		return nil
	}
	for _, message := range messages {
		if err := d.messageRepository.SaveMessage(message); err != nil {
			return err
		}
	}
	if err := d.dispatcher.MessagePublishedNotice(); err != nil {
		return err
	}
	return nil
}

// 消息发布器
// options["kafkaHosts"]="localhost:9092"
// options["timeInterval"]=time.Second*60*5
func NewMessageDirector(messageRepository models.MessageRepository, options map[string]interface{}) *MessageDirector {
	dispatcher := &MessageDispatcher{
		notifications:     make(chan struct{}),
		messageRepository: messageRepository,
		messageChan:       make(chan *models.Message, 100),
	}

	var hosts string
	if kafkaHosts, ok := options["kafkaHosts"]; ok {
		hosts = kafkaHosts.(string)
	} else {
		hosts = "localhost:9092"
	}
	dispatcher.producer = &KafkaMessageProducer{KafkaHosts: hosts, LogInfo: models.DefaultLog}

	if interval, ok := options["timeInterval"]; ok {
		dispatcher.dispatchTicker = time.NewTicker(interval.(time.Duration))
	} else {
		dispatcher.dispatchTicker = time.NewTicker(time.Second * 60 * 5)
	}
	go dispatcher.Dispatch()

	return &MessageDirector{
		messageRepository: messageRepository,
		dispatcher:        dispatcher,
	}
}
