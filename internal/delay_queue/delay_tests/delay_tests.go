package delay_tests

import (
	"encoding/json"
	"simplest_script/core"
)

type DelayTests struct {
}
type TestMsg struct {
	Id int `json:"id"`
}

func (l *DelayTests) Handler(params string) core.DelayQueueResult {
	var TestMsg TestMsg
	err := json.Unmarshal([]byte(params), &TestMsg)
	if err != nil {
		return core.DelayQueueResult{
			Status: true,
			Data:   "参数解析错误",
		}
	}

	// 业务逻辑

	return core.DelayQueueResult{
		Status: true,
		Data:   "success",
	}
}
