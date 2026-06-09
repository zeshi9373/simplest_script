package core

import "time"

const (
	MsgSuccess = "成功"
	MsgFail    = "失败"

	CodeOK        = 0
	CodeFail      = -1
	CodeLoginFail = -10002
	FlagDefault   = "default"
	FlagData      = "data"

	DBBusiness = "business"
	DBConsole  = "console"

	RDSDefault = "default"
	RDSData    = "data"

	MQExchangeDefault = ""

	MQQueueDefault = ""
	MQQueueData    = "app-data"
	MQQueueTrack   = "log-track"
	MQQueueMonitor = "log-monitor"

	ExpireTimeSecond10 = 10 * time.Second
	ExpireTimeSecond30 = 30 * time.Second
	ExpireTimeMinute   = time.Minute
	ExpireTimeHour     = time.Hour
	ExpireTimeHour3    = 3 * time.Hour
	ExpireTimeDay      = 24 * time.Hour
	ExpireTimeDay3     = 3 * 24 * time.Hour

	MachineScript1 = "script1"
	MachineScript2 = "script2"
)
