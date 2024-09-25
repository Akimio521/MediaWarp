package log

import "MediaWarp/internal/config"

var (
	cfg    *config.ConfigManager
	logger *LoggerManager = &LoggerManager{} // 全局日志实例
)

func init() {
	cfg = config.GetConfig()
	logger.Init()
}

func GetLogger() *LoggerManager {
	return logger
}
