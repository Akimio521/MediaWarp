package handler

import (
	"MediaWarp/internal/log"
)

var (
	logger *log.LoggerManager
)

func init() {
	logger = log.GetLogger()
}
