package kafkaclient

import (
	"simplest_script/core"
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

var KafkaWriterClient = make(map[string]chan KafkaWriterConnection, 0)
var mxWriter sync.Mutex

type Producer struct {
	Flag    string
	Brokers string
	MaxIdle int
}

func NewProducer(flag string) *Producer {
	switch flag {
	case core.FlagData:
		return &Producer{
			Flag:    flag,
			Brokers: conf.Conf.Kafka.Data.Brokers,
			MaxIdle: conf.Conf.Kafka.Data.MaxIdle,
		}
	default:
		return &Producer{
			Flag:    flag,
			Brokers: conf.Conf.Kafka.Default.Brokers,
			MaxIdle: conf.Conf.Kafka.Default.MaxIdle,
		}
	}
}

func (l *Producer) GetKafkaWriterClient(topic string) *kafka.Writer {
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

	if KafkaWriterClient[l.Flag+topic] == nil {
		KafkaWriterClient[l.Flag+topic] = make(chan KafkaWriterConnection, l.MaxIdle)
	}

	if len(KafkaWriterClient[l.Flag+topic]) >= l.MaxIdle {
		client.Close()
	} else {
		KafkaWriterClient[l.Flag+topic] <- KafkaWriterConnection{
			LastUsedTime: time.Now().Unix(),
			Writer:       client,
		}
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
