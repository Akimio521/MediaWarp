package log

import "MediaWarp/internal/config"

var (
	cfg    *config.ConfigManager
	logger *LoggerManager // 全局日志实例
)

func init() {
	cfg = config.GetConfig()
}

// 获取全局LoggerManager对象
func GetLogger() *LoggerManager {
	if logger == nil {
		logger := &LoggerManager{}
		logger.Init()
	}

	return logger
}
