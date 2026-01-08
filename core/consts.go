package core

import "time"

const (
	MsgSuccess = "成功"
	MsgFail    = "失败"

	CodeOK   = 0
	CodeFail = -1

	DBMain    = "test_main"
	DBConsole = "test_console"

	RDSDefault = "default"
	RDSData    = "data"

	MQExchangeDefault = ""

	MQQueueDefault        = ""
	MQQueueData           = "app-data"
	MQQueueTrack          = "log-track"
	MQQueueMonitor        = "log-monitor"
	MQQueueOppoCallBack   = "oppo-call-back"
	MQQueueSourceCallback = "source-callback"

	ExpireTimeSecond10 = 10 * time.Second
	ExpireTimeSecond30 = 30 * time.Second
	ExpireTimeMinute   = 60 * time.Second
	ExpireTimeHour     = 60 * 60 * time.Second
	ExpireTimeHour3    = 3 * 60 * 60 * time.Second
	ExpireTimeDay      = 24 * 60 * 60 * time.Second
	ExpireTimeDay3     = 3 * 24 * 60 * 60 * time.Second
)
