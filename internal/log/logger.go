package log

import (
	"github.com/sirupsen/logrus"
)

type LoggerManager struct {
	AccessLogger  *logrus.Logger
	ServiceLogger *logrus.Logger
}

// 初始化loggerManager
func (loggerManager *LoggerManager) Init() {
	var (
		aLS = &accessLoggerSetting{}  // 访问日志logrus相关设置
		sLS = &serviceLoggerSetting{} // 服务日志logrus相关设置
	)
	loggerManager.AccessLogger = logrus.New()
	loggerManager.ServiceLogger = logrus.New()

	if loggerManager.ServiceLogger == nil {
		panic("服务日志对象为nil")
	}

	loggerManager.ServiceLogger.SetReportCaller(true) // 设置报告调用方

	loggerManager.AccessLogger.SetFormatter(aLS)
	loggerManager.ServiceLogger.SetFormatter(sLS)

	if cfg.Logger.Enable { // 是否启用日志文件记录
		if cfg.Logger.AccessLogger.Enable {
			loggerManager.ServiceLogger.Debug("将访问日志记录到", cfg.AccessLogPath())
			loggerManager.AccessLogger.AddHook(aLS)
		} else {
			loggerManager.ServiceLogger.Debug("不启用访问日志")
		}
		if cfg.Logger.ServiceLogger.Enable {
			loggerManager.ServiceLogger.Debug("将服务日志记录到", cfg.ServiceLogPath())
			loggerManager.ServiceLogger.AddHook(sLS)
		} else {
			loggerManager.ServiceLogger.Debug("不启用服务日志")
		}
	} else {
		loggerManager.ServiceLogger.Warning("不启用日志记录")
	}
}
