package kafkaclient

import (
	"simplest_script/core/conf"
	"strings"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaWriterConnection struct {
	LastUsedTime int64
	Writer       *kafka.Writer
}

var KafkaWriterClient = make(map[string]chan *KafkaWriterConnection, 0)
var mxWriter sync.Mutex

type Producer struct {
	Flag    string
	Brokers string
	MaxIdle int
}

func NewProducer(flag string) *Producer {
	switch flag {
	default:
		return &Producer{
			Brokers: conf.Conf.Kafka.Default.Brokers,
			MaxIdle: conf.Conf.Kafka.Default.MaxIdle,
		}
	}
}

func (l *Producer) GetKafkaWriterClient(flag, topic string) *kafka.Writer {
	defer mxWriter.Unlock()
	mxWriter.Lock()

CreateProducer:

	if len(KafkaWriterClient[l.Flag+topic]) == 0 {
		return l.createWriter(topic)
	}

	client := <-KafkaWriterClient[l.Flag+topic]

	if time.Now().Unix()-client.LastUsedTime > 600 { // 关闭超过10分钟的连接
		client.Writer.Close()
		goto CreateProducer
	}

	return client.Writer
}

func (l *Producer) createWriter(topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(strings.Split(l.Brokers, ",")...),
		Topic:    topic,
		Balancer: &kafka.Hash{},
	}
}
func (l *Producer) PutKafkaWriterClient(topic string, client *kafka.Writer) {
	defer mxWriter.Unlock()
	mxWriter.Lock()

	if len(KafkaWriterClient[l.Flag+topic]) < l.MaxIdle {
		KafkaWriterClient[l.Flag+topic] <- &KafkaWriterConnection{
			Writer:       client,
			LastUsedTime: time.Now().Unix(),
		}
	} else {
		client.Close()
	}
}

func CloseKafkaWriterClient() {
	if len(KafkaWriterClient) == 0 {
		return
	}

	for _, v := range KafkaWriterClient {
		for {
			if len(v) == 0 {
				break
			}

			select {
			case conn := <-v:
				conn.Writer.Close()
			default:
				return
			}
		}
	}
}
