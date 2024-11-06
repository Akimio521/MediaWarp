package logging

import (
	"MediaWarp/internal/config"
	"io"

	"github.com/sirupsen/logrus"
)

var (
	accessLogger  *logrus.Logger // 访问日志
	serviceLogger *logrus.Logger // 服务日志
)

func init() {
	var (
		aLS = &accessLoggerSetting{}  // 访问日志logrus相关设置
		sLS = &serviceLoggerSetting{} // 服务日志logrus相关设置
	)
	accessLogger = logrus.New()
	serviceLogger = logrus.New()

	serviceLogger.SetReportCaller(true) // 设置报告调用方

	// 设置样式
	accessLogger.SetFormatter(aLS)
	serviceLogger.SetFormatter(sLS)

	if !config.Logger.AccessLogger.Console { // 访问日志不输出到终端
		accessLogger.Out = io.Discard
	}

	if !config.Logger.ServiceLogger.Console { // 服务日志不输出到终端
		serviceLogger.Out = io.Discard
	}

	if config.Logger.AccessLogger.File {
		accessLogger.AddHook(aLS)
	}

	if config.Logger.ServiceLogger.File {
		serviceLogger.AddHook(sLS)
	}

}

// 访问日志
//
// 默认日志级别为 Info
func AccessLog(format string, args ...interface{}) {
	accessLogger.Infof(format, args...)
}

// 服务日志
//
// Debug 级别日志
func Debug(args ...interface{}) {
	serviceLogger.Debug(args...)
}

// 服务日志
//
// Info 级别日志
func Info(args ...interface{}) {
	serviceLogger.Info(args...)
}

// 服务日志
//
// Warning 级别日志
func Warning(args ...interface{}) {
	serviceLogger.Warning(args...)
}

// 服务日志
//
// Error 级别日志
func Error(args ...interface{}) {
	serviceLogger.Error(args...)
}

// 服务日志
//
// Fatal 级别日志
func Fatal(args ...interface{}) {
	serviceLogger.Fatal(args...)
}

// 服务日志
//
// Panic 级别日志
func Panic(args ...interface{}) {
	serviceLogger.Panic(args...)
}

// 服务日志
//
// 设置日志级别
func SetLevel(level logrus.Level) {
	serviceLogger.SetLevel(level)
}
