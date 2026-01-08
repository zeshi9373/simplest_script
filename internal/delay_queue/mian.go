package delayqueue

import (
	"simplest_script/crontab"
	delayqueue "simplest_script/delay_queue"
)

type DelayQueue struct {
}

func (l *DelayQueue) Handler(params string) *crontab.Result {
	delayqueue.NewDelayQueue().Handler()

	return &crontab.Result{
		Status: 0,
		Data:   nil,
	}
}
