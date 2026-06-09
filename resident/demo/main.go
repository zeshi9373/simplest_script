package demo

import (
	"simplest_script/core/svc"
)

type Demo struct{}

func (h *Demo) Handler() {
	for {
		if svc.KillSignal {
			break
		}
		// 业务处理逻辑
	}

}
