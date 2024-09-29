package log

import (
	"io"

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

	loggerManager.ServiceLogger.SetReportCaller(true) // 设置报告调用方

	// 设置样式
	loggerManager.AccessLogger.SetFormatter(aLS)
	loggerManager.ServiceLogger.SetFormatter(sLS)

	if !cfg.Logger.AccessLogger.Console { // 访问日志不输出到终端
		loggerManager.AccessLogger.Out = io.Discard
	}

	if !cfg.Logger.ServiceLogger.Console { // 服务日志不输出到终端
		loggerManager.ServiceLogger.Out = io.Discard
	}

	if cfg.Logger.AccessLogger.File {
		loggerManager.AccessLogger.AddHook(aLS)
	}

	if cfg.Logger.ServiceLogger.File {
		loggerManager.ServiceLogger.AddHook(sLS)
	}
}
