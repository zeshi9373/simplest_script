package resident

import (
	"os"
	"simplest_script/core"
	"simplest_script/resident/demo"
)

type ResidentHandler interface {
	Handler()
}

var Entry = make(map[string]ResidentHandler)

func InitEntry() {
	paritition := os.Getenv("SCRIPT_PARTITION")

	switch paritition {
	case core.MachineScript1:
		Entry["demo"] = &demo.Demo{}
	case core.MachineScript2:
	default:

	}

}
