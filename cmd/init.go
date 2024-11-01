package cmd

import (
	"MediaWarp/internal/log"
	"flag"
)

var (
	logger  = log.GetLogger()
	isDebug bool
)

func init() {
	flag.BoolVar(&isDebug, "debug", false, "是否启用调试模式")
}
