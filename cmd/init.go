package cmd

import (
	"MediaWarp/core"
	"flag"
)

var (
	config  = core.GetConfig()
	logger  = core.GetLogger()
	isDebug bool
)

func init() {
	flag.BoolVar(&isDebug, "debug", false, "是否启用调试模式")
}
