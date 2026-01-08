package redislist

import (
	"context"
	"simplest_script/core"
	"simplest_script/core/svc"
	"simplest_script/core/svc/kafkaclient"
	"time"
)

func RedisListConsumer(ctx context.Context, key string, handler kafkaclient.SyncConsumer) {
	for {
		msg, err := svc.NewRedis(core.RDSDefault).BRPop(5*time.Second, key).Result()

		if err != nil || len(msg) < 2 {
			// hlog.Error("RedisListConsumer error: "+err.Error(), " key: "+key+" msg: "+strings.Join(msg, ","))
			continue
		}

		var status bool
		handler.Consume(msg[1], &status)

		if !status {
			svc.NewRedis(core.RDSDefault).LPush(key, msg)
		}

		select {
		case <-ctx.Done():
			return
		default:
		}

		time.Sleep(200 * time.Millisecond)
	}
}
