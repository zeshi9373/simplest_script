package exec

import (
	"fmt"
	"os"
	"simplest_script/core"
	"simplest_script/core/conf"
	"simplest_script/core/svc"
	"strings"

	"github.com/IBM/sarama"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// 监测kafka消费者需不需要增加消费者
func CheckKafkaProgress() {
	if os.Getenv("CSCRIPT_ENV") != core.EnvRelease {
		fmt.Println("监测kafka消费者")
	}

	saramaCfg := sarama.NewConfig()
	saramaCfg.ClientID = "cpa-" + os.Getenv("SCRIPT_PARTITION")
	saramaCfg.Version = sarama.V2_8_0_0
	saramaCfg.Consumer.Return.Errors = false

	client, err := sarama.NewClient(strings.Split(conf.Conf.Kafka.Brokers, ","), saramaCfg)
	if err != nil {
		fmt.Println("创建 Kafka 客户端失败:", err)
	}

	defer client.Close()

	for topic, detail := range Scripts {
		var lag int64
		if len(detail.Topic) > 0 {
			lag = getConsumerTopicLag(client, detail.GroupId, topic, detail.MaxProgress)
		} else if len(detail.Key) > 0 {
			lag = getRedisListLag(detail.Key)
		}

		if lag > int64(detail.ProgressLagLimit) {
			ProgressAdd(topic, lag)
		}
	}
}

func getConsumerTopicLag(client sarama.Client, groupId, topic string, partition int) int64 {
	var lag int64
	offsetMgr, err := sarama.NewOffsetManagerFromClient(groupId, client)

	if err != nil {
		fmt.Println("创建偏移量管理器失败:", err)
	}

	defer offsetMgr.Close()

	for i := 0; i < partition; i++ {
		latestOffset, err := client.GetOffset(topic, int32(i), sarama.OffsetNewest)

		if err != nil {
			hlog.Errorf("警告：获取分区 %d 最新偏移量失败: %v，跳过该分区", int32(i), err)
			continue
		}

		// 获取该分区的偏移量
		partitionOffset, err := offsetMgr.ManagePartition(topic, int32(i))
		if err != nil {
			hlog.Errorf("获取分区 %d 偏移量管理器失败: %w", int32(i), err)
			continue
		}

		// 获取已提交的偏移量（若无提交则返回 OffsetOldest）
		committedOffset, _ := partitionOffset.NextOffset()

		partitionOffset.Close()

		lag += latestOffset - committedOffset
	}

	if lag < 0 {
		lag = 0
	}

	return lag
}

func getRedisListLag(key string) int64 {
	var lag int64
	lag, _ = svc.NewRedis(core.RDSDefault).LLen(key).Result()

	return lag
}
