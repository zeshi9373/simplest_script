package resident

import (
	"os"
	"simplest_script/core"
)

type ResidentHandler interface {
	Handler()
}

var Entry = make(map[string]ResidentHandler)

// 需要区分不同机器
func InitEntry() {
	paritition := os.Getenv("CPA_PARTITION")

	switch paritition {
	case core.MachineScript1:

	default:

	}
}
