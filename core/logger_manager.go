package core

import (
	"github.com/sirupsen/logrus"
)

type loggerManager struct {
	AccessLogger *logrus.Logger
	ServerLogger *logrus.Logger
}

func (l *loggerManager) Init() {
	var sLS = &serviceLoggerSetting{}
	var aLS = &accessLoggerSetting{}

	l.ServerLogger.SetReportCaller(true) // 设置报告调用方

	l.AccessLogger.SetFormatter(aLS)
	l.ServerLogger.SetFormatter(sLS)

	if config.LoggerSetting.Enable { // 是否启用日志文件记录
		if config.LoggerSetting.AccessLogger.Enable {
			l.ServerLogger.Debug("将访问日志记录到", config.AccessLogPath())
			l.AccessLogger.AddHook(aLS)
		} else {
			l.ServerLogger.Debug("不启用访问日志")
		}
		if config.LoggerSetting.ServiceLogger.Enable {
			l.ServerLogger.Debug("将服务日志记录到", config.ServiceLogPath())
			l.ServerLogger.AddHook(sLS)
		} else {
			l.ServerLogger.Debug("不启用服务日志")
		}
	} else {
		l.ServerLogger.Debug("不启用日志记录")
	}
}

func init() {
	logger.Init()
}

// -----------------外部引用部分----------------- //
var logger = loggerManager{logrus.New(), logrus.New()}

func GetLogger() *loggerManager {
	return &logger
}
