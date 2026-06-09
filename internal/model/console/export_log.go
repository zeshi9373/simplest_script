package console

import (
	"simplest_script/core"
	"simplest_script/core/svc"

	"gorm.io/gorm"
)

// ExportLog export_log表的Model
type ExportLog struct {
	Id           int64  `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL" json:"id"`
	Title        string `gorm:"column:title;NOT NULL" json:"title"`                             // 标题
	Header       string `gorm:"column:header;NOT NULL" json:"header"`                           // 表头
	Query        string `gorm:"column:query;NOT NULL" json:"query"`                             // 请求参数
	Enums        string `gorm:"column:enums;NOT NULL" json:"enums"`                             // 枚举值
	FileName     string `gorm:"column:file_name;NOT NULL" json:"file_name"`                     // 文件名
	Status       int    `gorm:"column:status;NOT NULL;default:1" json:"status"`                 // 状态 1 等待导出 2进行中 3已完成 4失败
	Token        string `gorm:"column:token;NOT NULL" json:"token"`                             // token
	CreateUserId int    `gorm:"column:create_user_id;NOT NULL;default:0" json:"create_user_id"` // 用户id
	FilePath     string `gorm:"column:file_path;NOT NULL" json:"file_path"`                     // 文件地址
	FinishTime   string `gorm:"column:finish_time" json:"finish_time"`
	ErrorMsg     string `gorm:"column:error_msg;NOT NULL" json:"error_msg"` // 错误信息
	CreateTime   string `gorm:"column:create_time" json:"create_time"`
	UpdateTime   string `gorm:"column:update_time" json:"update_time"`
}

// 配置信息
type ExportLogConfig struct {
	Db    string
	Table string
}

// 获取配置
func GetExportLogConfig() ExportLogConfig {
	return ExportLogConfig{
		Db:    core.DBConsole,
		Table: "export_log",
	}
}

// 创建新的Model实例
func NewExportLogModel() *gorm.DB {
	return svc.NewDb(GetExportLogConfig().Db).Table(GetExportLogConfig().Table)
}
