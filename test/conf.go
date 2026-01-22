package test

import (
	"flag"
	"simplest_script/core/conf"
)

var configFile *string
var c conf.Config

func Init() {
	configFile = flag.String("f", "../../etc/dev.yaml", "the config file")
	flag.Parse()
	conf.MustLoad(*configFile, &c)
}
