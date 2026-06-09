package console

import (
	"simplest_script/core"
	"simplest_script/core/svc"

	"gorm.io/gorm"
)

// CrontabLog crontab_log表的Model
type CrontabLog struct {
	Id         int64  `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL" json:"id"`
	Pid        int    `gorm:"column:pid;NOT NULL;default:0" json:"pid"`               // 进程id
	Name       string `gorm:"column:name;NOT NULL" json:"name"`                       // 名称
	ExecCmd    string `gorm:"column:exec_cmd;NOT NULL" json:"exec_cmd"`               // 运行命令
	Params     string `gorm:"column:params;NOT NULL" json:"params"`                   // 参数
	Result     string `gorm:"column:result;NOT NULL" json:"result"`                   // 结果
	Status     string `gorm:"column:status;NOT NULL" json:"status"`                   // 状态
	StartTime  int64  `gorm:"column:start_time;NOT NULL;default:0" json:"start_time"` // 开始时间（毫秒）
	EndTime    int64  `gorm:"column:end_time;NOT NULL;default:0" json:"end_time"`     // 结束时间（毫秒）
	CostTime   int    `gorm:"column:cost_time;NOT NULL;default:0" json:"cost_time"`   // 运行时间（毫秒）
	Partition  string `gorm:"column:partition;NOT NULL" json:"partition"`             // 运行机器
	Uk         string `gorm:"column:uk;NOT NULL" json:"uk"`                           // 标识
	CreateTime string `gorm:"column:create_time" json:"create_time"`                  // 创建时间
	UpdateTime string `gorm:"column:update_time" json:"update_time"`                  // 更新时间
}

// 配置信息
type CrontabLogConfig struct {
	Db    string
	Table string
}

// 获取配置
func GetCrontabLogConfig() CrontabLogConfig {
	return CrontabLogConfig{
		Db:    core.DBConsole,
		Table: "crontab_log",
	}
}

// 创建新的Model实例
func NewCrontabLogModel() *gorm.DB {
	return svc.NewDb(GetCrontabLogConfig().Db).Table(GetCrontabLogConfig().Table)
}
