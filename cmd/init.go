package cmd

import (
	"flag"
)

var (
	isDebug bool
)

func init() {
	flag.BoolVar(&isDebug, "debug", false, "是否启用调试模式")
}
