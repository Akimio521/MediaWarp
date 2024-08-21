package core

import "github.com/sirupsen/logrus"

type loggerManager struct {
	AccessLogger *logrus.Logger
	ServerLogger *logrus.Logger
}

func (l *loggerManager) Init() {
	var sLS = &serviceLoggerSetting{}
	var aLS = &accessLoggerSetting{}
	l.AccessLogger.SetFormatter(aLS)
	l.ServerLogger.SetFormatter(sLS)
	l.AccessLogger.AddHook(aLS)
	l.ServerLogger.AddHook(sLS)
}

// -----------------外部引用部分----------------- //
var logger = loggerManager{logrus.New(), logrus.New()}

func GetLogger() *loggerManager {
	return &logger
}

func init() {
	logger.Init()
}
