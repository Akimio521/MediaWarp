package core

import "github.com/sirupsen/logrus"

var (
	logger = loggerManager{logrus.New(), logrus.New()}
	config = configManager{}
)

func init() {
	config.Init()
	logger.Init()
}

// -----------------外部引用部分----------------- //

func GetConfig() *configManager {
	return &config
}

func GetLogger() *loggerManager {
	return &logger
}
