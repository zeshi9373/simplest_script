package console

import (
	"simplest_script/core"
	"simplest_script/core/svc"
	"time"

	"gorm.io/gorm"
)

type ExportLogConfig struct {
	Db    string
	Table string
}

type ExportLog struct {
	Id           int64     `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	Title        string    `gorm:"column:title;NOT NULL"`
	Header       string    `gorm:"column:header;NOT NULL"`
	Query        string    `gorm:"column:query;NOT NULL"`
	Enums        string    `gorm:"column:enums;NOT NULL"`
	FileName     string    `gorm:"column:file_name;NOT NULL"`
	Status       int       `gorm:"column:status;default:1;NOT NULL"`
	Token        string    `gorm:"column:token;NOT NULL"`
	CreateUserId int       `gorm:"column:create_user_id;default:0;NOT NULL"`
	FilePath     string    `gorm:"column:file_path;NOT NULL"`
	FinishTime   time.Time `gorm:"column:finish_time"`
	ErrorMsg     string    `gorm:"column:error_msg;NOT NULL"`
	CreateTime   time.Time `gorm:"column:create_time"`
	UpdateTime   time.Time `gorm:"column:update_time"`
}

func GetExportLogConfig() ExportLogConfig {
	return ExportLogConfig{
		Db:    core.DBConsole,
		Table: "export_log",
	}
}

func NewExportLogModel() *gorm.DB {
	return svc.NewDb(GetExportLogConfig().Db).Table(GetExportLogConfig().Table)
}
