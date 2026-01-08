package conf

import (
	"os"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	yaml "gopkg.in/yaml.v2"
)

var Conf *Config

func MustLoad(filepath string, v *Config) {
	byteStream, err := os.ReadFile(filepath)

	if err != nil {
		hlog.Info("read config file error")
	}

	if err := yaml.Unmarshal(byteStream, v); err != nil {
		hlog.Info("unmarshal config file error")
	}

	Conf = v
}
