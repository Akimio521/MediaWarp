package log

import "MediaWarp/internal/config"

var (
	cfg    *config.ConfigManager
	logger *LoggerManager // 全局日志实例
)

func init() {
	cfg = config.GetConfig()
}

// 实例化一个LoggerManager对象
func NewLogger() *LoggerManager {
	lm := LoggerManager{}
	lm.Init()
	return &lm
}

// 获取全局LoggerManager对象
func GetLogger() *LoggerManager {
	if logger == nil {
		logger = NewLogger()
	}
	return logger
}
