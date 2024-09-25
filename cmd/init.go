package cmd

import (
	"MediaWarp/internal/config"
	"flag"
)

var (
	cfg     = config.GetConfig()
	isDebug bool
)

func init() {
	flag.BoolVar(&isDebug, "debug", false, "是否启用调试模式")
}
