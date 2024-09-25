package cmd

import (
	"MediaWarp/constants"
	"MediaWarp/pkg"
	"fmt"
)

// 打印LOGO
func PrintLOGO() {
	fmt.Print(
		constants.LOGO,
		pkg.Center(
			fmt.Sprintf(" MediaWarp %s 启动中... ", cfg.Version()),
			75,
			"=",
		),
		"\n\n",
	)
}
