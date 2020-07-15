package broker

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/tiptok/gocomm/pkg/log"
	"testing"
	"time"
)

/*
	kafka golang client github.com/Shopify/sarama
	测试
*/
//生产
func ExampleProducer() {
	config := sarama.NewConfig()
	//等待服务器所有副本都保存成功后的响应
	config.Producer.RequiredAcks = sarama.WaitForAll
	//随机的分区类型
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	//是否等待成功和失败后的响应,只有上面的RequireAcks设置不是NoReponse这里才有用.
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	//设置使用的kafka版本,如果低于V0_10_0_0版本,消息中的timestrap没有作用.需要消费和生产同时配置
	config.Version = sarama.V0_11_0_0

	//使用配置,新建一个异步生产者
	producer, e := sarama.NewAsyncProducer([]string{"127.0.0.1:9092"}, config)

	if e != nil {
		panic(e)
	}
	defer producer.AsyncClose()

	//发送的消息,主题,key
	msg := &sarama.ProducerMessage{
		Topic: "ability",
		//Key:   sarama.StringEncoder("test"),
	}

	var value string
	//for {
	value = "this is a message!!"
	//设置发送的真正内容
	//fmt.Scanln(&value)
	//将字符串转化为字节数组
	msg.Value = sarama.ByteEncoder(value)
	fmt.Println(value)

	//使用通道发送
	producer.Input() <- msg

	//循环判断哪个通道发送过来数据.
	select {
	case suc := <-producer.Successes():
		fmt.Println("offset: ", suc.Offset, "timestamp: ", suc.Timestamp.Format("2006-Jan-02 15:04"), "partitions: ", suc.Partition)
	case fail := <-producer.Errors():
		fmt.Println("err: ", fail.Err)
	}
	//}
}

//消费
func ExampleComsumer() {
	config := sarama.NewConfig()
	//接收失败通知
	config.Consumer.Return.Errors = true
	//设置使用的kafka版本,如果低于V0_10_0_0版本,消息中的timestrap没有作用.需要消费和生产同时配置
	config.Version = sarama.V0_11_0_0
	//新建一个消费者
	consumer, e := sarama.NewConsumer([]string{"127.0.0.1:9092"}, config)
	if e != nil {
		panic("error get consumer")
	}
	defer consumer.Close()

	//根据消费者获取指定的主题分区的消费者,Offset这里指定为获取最新的消息.
	partitionConsumer, err := consumer.ConsumePartition("ability", 0, sarama.OffsetNewest)
	if err != nil {
		fmt.Println("error get partition consumer", err)
	}
	timeout := time.After(time.Second * 60 * 5)
	defer partitionConsumer.Close()
	//循环等待接受消息.
	for {
		select {
		//接收消息通道和错误通道的内容.
		case msg := <-partitionConsumer.Messages():
			fmt.Println("key: ", string(msg.Key), "msg offset: ", msg.Offset, " partition: ", msg.Partition, " timestrap: ", msg.Timestamp.Format("2006-Jan-02 15:04"), " value: ", string(msg.Value))
		case err := <-partitionConsumer.Errors():
			fmt.Println(err.Err)
		case <-timeout:
			return
		}
	}
}

func ExampleClient() {
	config := sarama.NewConfig()
	config.Version = sarama.V0_11_0_0
	client, err := sarama.NewClient([]string{"127.0.0.1:9092"}, config)
	if err != nil {
		panic("client create error")
	}
	defer client.Close()
	//获取主题的名称集合
	topics, err := client.Topics()
	if err != nil {
		panic("get topics err")
	}
	for _, e := range topics {
		log.Info(e)
	}
	//获取broker集合
	brokers := client.Brokers()
	//输出每个机器的地址
	for _, broker := range brokers {
		log.Info(broker.Addr())
	}
}

func Test_Broker(t *testing.T) {
	//ExampleComsumer()
}
