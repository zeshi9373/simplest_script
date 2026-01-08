package console

import (
	"simplest_script/core"
	"simplest_script/core/svc"
	"time"

	"gorm.io/gorm"
)

type CrontabLogConfig struct {
	Db    string
	Table string
}

type CrontabLog struct {
	Id         int64     `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	Pid        int       `gorm:"column:pid;default:0;NOT NULL"`
	Name       string    `gorm:"column:name;NOT NULL"`
	ExecCmd    string    `gorm:"column:exec_cmd;NOT NULL"`
	Params     string    `gorm:"column:params;NOT NULL"`
	Result     string    `gorm:"column:result;NOT NULL"`
	Status     string    `gorm:"column:status;NOT NULL"`
	StartTime  int       `gorm:"column:start_time;default:0;NOT NULL"`
	EndTime    int       `gorm:"column:end_time;default:0;NOT NULL"`
	CostTime   int       `gorm:"column:cost_time;default:0;NOT NULL"`
	Partition  string    `gorm:"column:partition;NOT NULL"`
	Uk         string    `gorm:"column:uk;NOT NULL"`
	CreateTime time.Time `gorm:"column:create_time"`
	UpdateTime time.Time `gorm:"column:update_time"`
}

func GetCrontabLogConfig() CrontabLogConfig {
	return CrontabLogConfig{
		Db:    core.DBConsole,
		Table: "crontab_log",
	}
}
func NewCrontabLogModel() *gorm.DB {
	return svc.NewDb(GetCrontabLogConfig().Db).Table(GetCrontabLogConfig().Table)
}
