package script

import (
	delayqueue "simplest_script/internal/delay_queue"
	"simplest_script/internal/handler/export"
)

func InitEntry() {
	HandlerEntry["export"] = &export.Export{}
	HandlerEntry["delay_queue"] = &delayqueue.DelayQueue{}

}
