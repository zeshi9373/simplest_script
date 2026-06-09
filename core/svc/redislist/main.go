package redislist

import (
	"context"
	"simplest_script/core"
	"simplest_script/core/svc"
	"simplest_script/core/svc/kafkaclient"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func RedisListConsumer(ctx context.Context, key string, handler kafkaclient.SyncConsumer) {
	for {
		msg, err := svc.NewRedis(core.RDSDefault).BRPop(5*time.Second, key).Result()

		if err != nil || len(msg) < 2 {
			if err != nil {
				hlog.Errorf("RedisListConsumer error: %v, key: %s", err, key)
			}
			continue
		}

		var status bool
		handler.Consume(msg[1], &status)

		if !status {
			svc.NewRedis(core.RDSDefault).LPush(key, msg[1])
		}

		select {
		case <-ctx.Done():
			return
		default:
		}

		if svc.KillSignal {
			hlog.Info("程序退出，redis list退出， key: " + key)
			return
		}
	}
}
