package logger

import (
	"simplest_script/core/logger"
	"time"
)

func Info() {
	logger.NewLogger("test_log").Info("测试", logger.Fields{
		"test": "测试",
	})
	logger.NewLogger("test_log_2").Info("测试", logger.Fields{
		"test": "测试2",
	})
	logger.NewLogger("test_log_3").Info("测试", logger.Fields{
		"test": "测试3",
	})
	logger.NewLogger("test_log_4").Info("测试", logger.Fields{
		"test": "测试4",
	})
	time.Sleep(2 * time.Second)
}
