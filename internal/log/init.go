package log

var (
	logger *LoggerManager // 全局日志实例
)

// 获取全局LoggerManager对象
func GetLogger() *LoggerManager {
	if logger == nil {
		logger := &LoggerManager{}
		logger.Init()
	}

	return logger
}
