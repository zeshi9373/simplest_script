package kafkaclient

import (
	"context"
	"simplest_script/core/conf"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/segmentio/kafka-go"
)

type SyncConsumer interface {
	Consume(msg string, status *bool)
}

func GetConsumerClient(topic string, groupId string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:        strings.Split(conf.Conf.Kafka.Brokers, ","),
		Topic:          topic,
		GroupID:        groupId,
		MaxBytes:       10e6, // 10MB
		CommitInterval: time.Second,
	})
}

func ConsumerHandlerMessage(ctx context.Context, topic string, groupId string, handler SyncConsumer) {
	reader := GetConsumerClient(topic, groupId)
	defer reader.Close()

	for {
		m, err := reader.ReadMessage(ctx)

		if err != nil {
			hlog.Error("Kafka消费出错 topic: " + topic + " groupId: " + groupId + " error: " + err.Error())
			continue
		}

		var status bool
		handler.Consume(string(m.Value), &status)

		if status {
			reader.CommitMessages(ctx, m)
		}
	}
}
