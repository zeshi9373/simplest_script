package core

type DelayQueueResult struct {
	Data   any  `json:"data"`
	Status bool `json:"status"`
}

type DelayQueuePushParams struct {
	Name      string `json:"name"`
	ExecCmd   string `json:"exec_cmd"`
	Params    string `json:"params"`
	DelayTime int64  `json:"delay_time"` // 延迟秒数
	ExecTime  int    `json:"exec_time"`
}
