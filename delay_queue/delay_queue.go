package delayqueue

import (
	"encoding/json"
	"simplest_script/core"
	"simplest_script/core/tool"
	"simplest_script/internal/model/console"
	"sync"
	"time"
)

type DelayQueue struct {
}

func NewDelayQueue() *DelayQueue {
	return &DelayQueue{}
}

func (l *DelayQueue) Handler() {
	minStartTime := time.Now().Truncate(time.Minute).Unix()
	endTime := minStartTime + 60
	list := make([]console.DelayQueueLog, 0)
	console.NewDelayQueueLogModel().Where("status = 1 and exec_time < ?", endTime).Find(&list)

	var gr sync.WaitGroup

	for _, log := range list {
		gr.Add(1)
		go l.delayQueue(&gr, log)
	}

	gr.Wait()

	time.Sleep(2 * time.Second) // 延迟2秒,保证日志写入成功
}

func (l *DelayQueue) delayQueue(gr *sync.WaitGroup, data console.DelayQueueLog) {
	defer gr.Done()
	// 延迟毫秒数
	sleepMilli := int64(data.ExecTime*1000) - time.Now().UnixMilli()

	if sleepMilli > 0 {
		time.Sleep(time.Duration(sleepMilli) * time.Millisecond)
	}

	console.NewDelayQueueLogModel().Where("id = ?", data.Id).Updates(map[string]any{
		"status":      2,
		"update_time": time.Now().Format(tool.DatetimeLayout),
	})

	InitEntry()
	if _, ok := HandlerEntry[data.ExecCmd]; !ok {
		console.NewDelayQueueLogModel().Where("id = ?", data.Id).Updates(map[string]any{
			"status":      4,
			"result":      "请输入脚本方法不存在 " + data.ExecCmd,
			"update_time": time.Now().Format(tool.DatetimeLayout),
		})
	}

	res := HandlerEntry[data.ExecCmd].Handler(data.Params)
	resData, _ := json.Marshal(res)

	if res.Status {
		console.NewDelayQueueLogModel().Where("id = ?", data.Id).Updates(map[string]any{
			"status":      3,
			"result":      string(resData),
			"update_time": time.Now().Format(tool.DatetimeLayout),
		})
	} else {
		console.NewDelayQueueLogModel().Where("id = ?", data.Id).Updates(map[string]any{
			"status":      4,
			"result":      string(resData),
			"update_time": time.Now().Format(tool.DatetimeLayout),
		})
	}
}

func (l *DelayQueue) Push(params ...core.DelayQueuePushParams) {
	for _, param := range params {
		var execTime int

		if param.ExecTime > 0 {
			execTime = param.ExecTime
		} else {
			execTime = int(time.Now().Unix() + param.DelayTime)
		}

		log := console.DelayQueueLog{
			Name:       param.Name,
			ExecCmd:    param.ExecCmd,
			ExecTime:   execTime,
			Params:     param.Params,
			Status:     1,
			CreateTime: time.Now().Format(tool.DatetimeLayout),
			UpdateTime: time.Now().Format(tool.DatetimeLayout),
		}

		console.NewDelayQueueLogModel().Create(&log)
	}
}
