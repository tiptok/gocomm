package kafkax

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/tiptok/gocomm/common"
	"github.com/tiptok/gocomm/identity/idgen"
	"github.com/tiptok/gocomm/pkg/broker/models"
	"github.com/tiptok/gocomm/pkg/log"
	"strings"
	"sync"
	"time"
)

type SaramaConsumer struct {
	ready             chan bool
	messageHandlerMap map[string]func(message interface{}) error
	//Logger            log.Logger
	kafkaHosts string
	groupId    string
	topicMiss  map[string]string //记录未被消费的topic
	receiver   models.MessageReceiverRepository
}

func (consumer *SaramaConsumer) Setup(sarama.ConsumerGroupSession) error {
	close(consumer.ready)
	return nil
}
func (consumer *SaramaConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}
func (consumer *SaramaConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	var err error
	for message := range claim.Messages() {
		log.Debug(fmt.Sprintf("【kafka】 Receive Message claimed:  timestamp = %v, topic = %s offset = %v value = %v", message.Timestamp, message.Topic, message.Offset, string(message.Value)))
		handler, ok := consumer.messageHandlerMap[message.Topic]
		if e := consumer.messageReceiveBefore(message); e != nil {
			ok = false
			log.Error(e)
		}
		if !ok {
			continue
		}
		var msg = &models.Message{}
		common.JsonUnmarshal(string(message.Value), msg)
		if err = handler(msg); err == nil {
			session.MarkMessage(message, "")
		} else {
			log.Error("Message claimed: kafka消息处理错误 topic =", message.Topic, message.Offset, err)
		}
		if err != nil {
			continue
		}
		consumer.messageReceiveAfter(message)
	}
	return err
}

func (consumer *SaramaConsumer) messageReceiveBefore(message *sarama.ConsumerMessage) error {
	if consumer.receiver == nil {
		return nil
	}

	var params = make(map[string]interface{})
	var err error
	_, ok := consumer.messageHandlerMap[message.Topic]
	if !ok {
		params["status"] = models.Ignore
		_, topicMiss := consumer.topicMiss[message.Topic]
		if !topicMiss {
			fmt.Printf("topic:[%v] has not consumer handler", message.Topic)
		}
		return nil
	}

	_, err = consumer.storeMessage(params, message)
	if err != nil {
		ok = false
		//log.Println("ConsumeClaim:", err)
	}
	return err
}
func (consumer *SaramaConsumer) messageReceiveAfter(message *sarama.ConsumerMessage) {
	if consumer.receiver == nil {
		return
	}
	consumer.finishMessage(map[string]interface{}{"offset": message.Offset, "topic": message.Topic})
}

func (consumer *SaramaConsumer) storeMessage(params map[string]interface{}, message *sarama.ConsumerMessage) (id int64, err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Error(e)
		}
	}()
	id = idgen.Next()
	params = make(map[string]interface{})
	params["id"] = message.Offset
	params["topic"] = message.Topic
	params["partition"] = message.Partition
	params["offset"] = message.Offset
	params["key"] = string(message.Key)
	params["value"] = string(message.Value)
	params["msg_time"] = message.Timestamp.Unix()
	params["create_at"] = time.Now().Unix()
	params["status"] = models.UnFinished //0:未完成 1:已完成 2：未命中
	err = consumer.receiver.ReceiveMessage(params)
	return
}
func (consumer *SaramaConsumer) finishMessage(params map[string]interface{}) error {
	defer func() {
		if e := recover(); e != nil {
			log.Error(e)
		}
	}()
	consumer.receiver.ConfirmReceive(params)
	return nil
}

func (consumer *SaramaConsumer) StartConsume() error {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Version = sarama.V0_11_0_0
	brokerList := strings.Split(consumer.kafkaHosts, ",")
	consumerGroup, err := sarama.NewConsumerGroup(brokerList, consumer.groupId, config)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			var topics []string
			for key := range consumer.messageHandlerMap {
				topics = append(topics, key)
			}
			if err := consumerGroup.Consume(ctx, topics, consumer); err != nil {
				log.Error(err.Error())
				return
			}
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()
	<-consumer.ready
	log.Info("Sarama consumer up and running!...")
	select {
	case <-ctx.Done():
		log.Info("Sarama consumer : context cancelled")
	}
	cancel()
	wg.Wait()
	if err := consumerGroup.Close(); err != nil {
		return err
	}
	return nil
}
func (consumer *SaramaConsumer) WithTopicHandler(topic string, handler func(message interface{}) error) { //*sarama.ConsumerMessage
	consumer.messageHandlerMap[topic] = handler
}
func (consumer *SaramaConsumer) WithMessageReceiver(receiver models.MessageReceiverRepository) {
	consumer.receiver = receiver
}

func NewSaramaConsumer(kafkaHosts string, groupId string) models.Consumer {
	return &SaramaConsumer{
		kafkaHosts:        kafkaHosts,
		groupId:           groupId,
		topicMiss:         make(map[string]string),
		messageHandlerMap: make(map[string]func(message interface{}) error),
		ready:             make(chan bool),
	}
}
