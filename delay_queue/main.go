package delayqueue

import (
	"simplest_script/core"
	delayTests "simplest_script/internal/delay_queue/delay_tests"
)

type HandlerFunc interface {
	Handler(params string) core.DelayQueueResult
}

var HandlerEntry = make(map[string]HandlerFunc)

// 延时队列
func InitEntry() {
	HandlerEntry["delay_test"] = &delayTests.DelayTests{}
}
