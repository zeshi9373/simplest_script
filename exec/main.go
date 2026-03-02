package exec

import (
	"context"
	"fmt"
	"math"
	"os"
	"simplest_script/core"
	"simplest_script/core/svc/kafkaclient"
	"simplest_script/core/svc/redislist"
	"simplest_script/internal/model/console"
	"sync"
	"time"
)

var Scripts = make(map[string]ScriptConfig, 0)
var TopicNum = make(map[string]int, 0)
var mx sync.Mutex

type Paritition struct {
	Paritition int32
	NewOffset  int64
	Offset     int64
}

type ScriptConfig struct {
	Flag                string `json:"flag"`
	Name                string `json:"name"`
	ExecCmd             string `json:"exec_cmd"`
	Key                 string `json:"key"`
	Topic               string `json:"topic"`
	Progress            int    `json:"progress"`
	GroupId             string `json:"group_id"`
	MaxProgress         int    `json:"max_progress"`
	ProgressLagLimit    int    `json:"progress_lag_limit"`
	ProgressAvgMsgcount int    `json:"progress_avg_msgcount"`
	Status              int    `json:"status"`
}

func Init() {
	paritition := os.Getenv("SCRIPT_PARTITION")
	scripts := make([]console.ScriptConfig, 0)
	console.NewScriptConfigModel().Where("type = 1 and machine = ? and status = 1", paritition).Find(&scripts)

	for _, v := range scripts {
		if len(v.Topic) > 0 {
			Scripts[v.Topic] = ScriptConfig{
				Flag:                v.Flag,
				Name:                v.Name,
				ExecCmd:             v.ExecCmd,
				Key:                 v.Key,
				Topic:               v.Topic,
				Progress:            v.Progress,
				GroupId:             v.GroupId,
				MaxProgress:         v.MaxProgress,
				ProgressLagLimit:    v.ProgressLagLimit,
				ProgressAvgMsgcount: v.ProgressAvgMsgcount,
				Status:              v.Status,
			}
		} else {
			Scripts[v.Key] = ScriptConfig{
				Flag:                v.Flag,
				Name:                v.Name,
				ExecCmd:             v.ExecCmd,
				Key:                 v.Key,
				Topic:               v.Topic,
				Progress:            v.Progress,
				GroupId:             v.GroupId,
				MaxProgress:         v.MaxProgress,
				ProgressLagLimit:    v.ProgressLagLimit,
				ProgressAvgMsgcount: v.ProgressAvgMsgcount,
				Status:              v.Status,
			}
		}
	}

	InitEntry()

	for _, v := range scripts {
		if core.StatusIsEnv(v.Status) {
			if _, ok := Entry[v.ExecCmd]; ok {
				if len(v.Topic) == 0 && len(v.Key) == 0 {
					continue
				}

				fmt.Printf("exec script %v \n", v)

				if len(v.Topic) > 0 {
					for i := 0; i < v.Progress; i++ {
						go func() {
							kafkaclient.NewConsumer(v.Flag).ConsumerHandlerMessage(context.Background(), v.Topic, v.GroupId, Entry[v.ExecCmd])
						}()
					}
				} else if len(v.Key) > 0 {
					for i := 0; i < v.Progress; i++ {
						go func() {
							redislist.RedisListConsumer(context.Background(), v.Key, Entry[v.ExecCmd])
						}()
					}
				}
			}
		}
	}

	// 启动定时检查消息堆积调度消费者数量
	go progressCronCheck()
}

func ProgressAdd(flag, topic string, msgCount int64) {
	defer mx.Unlock()
	mx.Lock()

	if _, ok := Scripts[topic]; !ok {
		return
	}

	if Scripts[topic].MaxProgress == Scripts[topic].Progress {
		return
	}

	var progressNum int
	progressNumPlan := int(math.Floor(float64(msgCount) / float64(Scripts[topic].ProgressAvgMsgcount)))
	progressNumAllow := Scripts[topic].MaxProgress - TopicNum[topic]

	if progressNumPlan > progressNumAllow {
		progressNum = progressNumAllow
	} else {
		progressNum = progressNumPlan
	}

	if progressNum > 0 {
		if _, ok := Entry[Scripts[topic].ExecCmd]; ok {
			for i := 0; i < progressNum; i++ {
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3) // 3分钟

				go func() {
					defer TopicNumDone(topic)
					defer cancel()
					TopicNumAdd(topic)

					if len(Scripts[topic].Topic) > 0 {
						kafkaclient.NewConsumer(flag).ConsumerHandlerMessage(ctx, Scripts[topic].Topic, Scripts[topic].GroupId, Entry[Scripts[topic].ExecCmd])
					} else {
						redislist.RedisListConsumer(ctx, Scripts[topic].Key, Entry[Scripts[topic].ExecCmd])
					}
				}()
			}
		}
	}
}

func TopicNumAdd(topic string) {
	mx.Lock()
	defer mx.Unlock()

	TopicNum[topic]++
}

func TopicNumDone(topic string) {
	mx.Lock()
	defer mx.Unlock()

	TopicNum[topic]--
}

func progressCronCheck() {
	ticker := time.NewTicker(3 * time.Minute)
	defer kafkaClientClose()
	defer ticker.Stop()

	for range ticker.C {
		// 调用检查Kafka消息余量的函数
		CheckKafkaProgress()
	}
}
