package core

import "os"

const (
	EnvRelease      = "release"
	EnvTest         = "test"
	EnvDev          = "dev"
	EnvPre          = "pre"
	EnvPreValue     = 8
	EnvReleaseValue = 4
	EnvTestValue    = 2
	EnvDevValue     = 1
)

func StatusIsEnv(status int) bool {
	switch os.Getenv("SCRIPT_ENV") {
	case EnvRelease:
		return status&EnvReleaseValue == EnvReleaseValue
	case EnvTest:
		return status&EnvTestValue == EnvTestValue
	case EnvDev:
		return status&EnvDevValue == EnvDevValue
	case EnvPre:
		return status&EnvPreValue == EnvPreValue
	default:
		return status&EnvDevValue == EnvDevValue
	}
}
