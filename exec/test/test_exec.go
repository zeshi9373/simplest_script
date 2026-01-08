package test

import "encoding/json"

type testData struct {
	Id int `json:"id"`
}

type TestExec struct{}

func (h *TestExec) Consume(msg string, status *bool) {
	data := testData{}
	json.Unmarshal([]byte(msg), &data)

	if data.Id == 0 {
		*status = true
		return
	}

	// 业务逻辑

	*status = true
}
