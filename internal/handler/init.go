package handler

import (
	"MediaWarp/internal/config"
	"MediaWarp/internal/log"
)

var (
	cfg    *config.ConfigManager
	logger *log.LoggerManager
)

func init() {
	cfg = config.GetConfig()
	logger = log.GetLogger()
}
