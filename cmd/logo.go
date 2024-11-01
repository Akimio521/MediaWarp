package cmd

import (
	"MediaWarp/constants"
	"MediaWarp/internal/config"
	"MediaWarp/pkg"
	"fmt"
)

// 打印LOGO
func PrintLOGO() {
	fmt.Print(
		constants.LOGO,
		pkg.Center(
			fmt.Sprintf(" MediaWarp %s 启动中... ", config.Version()),
			75,
			"=",
		),
		"\n\n",
	)
}
