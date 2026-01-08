package exec

import (
	"simplest_script/core/svc/kafkaclient"
	"simplest_script/exec/test"
)

var Entry = make(map[string]kafkaclient.SyncConsumer)

func InitEntry() {
	Entry["exec_test"] = &test.TestExec{}
}
