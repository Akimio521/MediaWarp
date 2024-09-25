package middleware

import (
	"MediaWarp/internal/config"
	"MediaWarp/internal/log"
)

var (
	cfg    = config.GetConfig()
	logger = log.GetLogger()
)
