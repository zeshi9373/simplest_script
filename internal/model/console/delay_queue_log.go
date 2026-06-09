package console

import (
	"simplest_script/core"
	"simplest_script/core/svc"

	"gorm.io/gorm"
)

// DelayQueueLog delay_queue_log表的Model
type DelayQueueLog struct {
	Id         int64  `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL" json:"id"`
	Name       string `gorm:"column:name;NOT NULL" json:"name"`                     // 名称
	ExecCmd    string `gorm:"column:exec_cmd;NOT NULL" json:"exec_cmd"`             // 执行方法
	Params     string `gorm:"column:params;NOT NULL" json:"params"`                 // 参数
	Status     int    `gorm:"column:status;NOT NULL;default:1" json:"status"`       // 执行状态（1待执行 2执行中 3已完成 4失败）
	Result     string `gorm:"column:result;NOT NULL" json:"result"`                 // 执行结果
	ExecTime   int    `gorm:"column:exec_time;NOT NULL;default:0" json:"exec_time"` // 执行时间（秒）
	CreateTime string `gorm:"column:create_time" json:"create_time"`
	UpdateTime string `gorm:"column:update_time" json:"update_time"`
}

// 配置信息
type DelayQueueLogConfig struct {
	Db    string
	Table string
}

// 获取配置
func GetDelayQueueLogConfig() DelayQueueLogConfig {
	return DelayQueueLogConfig{
		Db:    core.DBConsole,
		Table: "delay_queue_log",
	}
}

// 创建新的Model实例
func NewDelayQueueLogModel() *gorm.DB {
	return svc.NewDb(GetDelayQueueLogConfig().Db).Table(GetDelayQueueLogConfig().Table)
}
