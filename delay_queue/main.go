package delayqueue

import (
	"simplest_script/core"
)

type HandlerFunc interface {
	Handler(params string) core.DelayQueueResult
}

var HandlerEntry = make(map[string]HandlerFunc)

// 延时队列
func InitEntry() {

}
