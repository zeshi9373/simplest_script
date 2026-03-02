package console

import (
	"simplest_script/core"
	"simplest_script/core/svc"
	"time"

	"gorm.io/gorm"
)

// ScriptConfig script_config表的Model
type ScriptConfig struct {
	Id                  int       `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL" json:"id"`
	Type                int       `gorm:"column:type;NOT NULL;default:1" json:"type"`                                 // 类型 1异步消费 2定时任务
	Cron                string    `gorm:"column:cron;NOT NULL" json:"cron"`                                           // 定时表达式
	Flag                string    `gorm:"column:flag;NOT NULL" json:"flag"`                                           // kafka实例标识
	Name                string    `gorm:"column:name;NOT NULL" json:"name"`                                           // 名称
	ExecCmd             string    `gorm:"column:exec_cmd;NOT NULL" json:"execCmd"`                                    // 执行指令
	Params              string    `gorm:"column:params;NOT NULL" json:"params"`                                       // 参数
	Topic               string    `gorm:"column:topic;NOT NULL" json:"topic"`                                         // 主题
	GroupId             string    `gorm:"column:group_id;NOT NULL" json:"groupId"`                                    // 消费组
	Key                 string    `gorm:"column:key;NOT NULL" json:"key"`                                             // redis队列key
	Progress            int       `gorm:"column:progress;NOT NULL;default:2" json:"progress"`                         // 初始消费者数量
	MaxProgress         int       `gorm:"column:max_progress;NOT NULL;default:0" json:"maxProgress"`                  // 最大消费者数量
	ProgressLagLimit    int       `gorm:"column:progress_lag_limit;NOT NULL;default:0" json:"progressLagLimit"`       // 消息堆积阈值
	ProgressAvgMsgcount int       `gorm:"column:progress_avg_msgcount;NOT NULL;default:0" json:"progressAvgMsgcount"` // 新增单个消费者处理的消息数
	Status              int       `gorm:"column:status;NOT NULL;default:0" json:"status"`                             // 状态 1 执行 0不执行
	Machine             string    `gorm:"column:machine;NOT NULL" json:"machine"`                                     // 运行机器标识
	IsLog               int       `gorm:"column:is_log;NOT NULL;default:0" json:"isLog"`                              // 定时任务是否记日志
	CreateTime          time.Time `gorm:"column:create_time;NOT NULL;default:CURRENT_TIMESTAMP" json:"createTime"`    // 创建时间
	UpdateTime          time.Time `gorm:"column:update_time;NOT NULL;default:0000-00-00 00:00:00" json:"updateTime"`  // 更新时间
}

// 配置信息
type ScriptConfigConfig struct {
	Db    string
	Table string
}

// 获取配置
func GetScriptConfigConfig() ScriptConfigConfig {
	return ScriptConfigConfig{
		Db:    core.DBConsole,
		Table: "script_config",
	}
}

// 创建新的Model实例
func NewScriptConfigModel() *gorm.DB {
	return svc.NewDb(GetScriptConfigConfig().Db).Table(GetScriptConfigConfig().Table)
}
