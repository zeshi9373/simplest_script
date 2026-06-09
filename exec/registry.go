package exec

import (
	"simplest_script/core/svc/kafkaclient"
)

var Entry = make(map[string]kafkaclient.SyncConsumer)

func InitEntry() {

}
