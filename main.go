package main

import (
	"MediaWarp/constants"
	"MediaWarp/internal/config"
	"MediaWarp/internal/logging"
	"MediaWarp/internal/router"
	"MediaWarp/utils"
	"flag"
	"fmt"

	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	isDebug     bool
	showVersion bool
)

func init() {
	printLOGO()
	flag.BoolVar(&isDebug, "debug", false, "是否启用调试模式")
	flag.BoolVar(&showVersion, "version", false, "显示版本信息")
}

func main() {
	flag.Parse()

	if showVersion {
		ver, _ := json.MarshalIndent(config.Version(), "", "  ")
		fmt.Println(string(ver))
		return
	}

	if isDebug {
		logging.SetLevel(logrus.DebugLevel)
		logging.Warning("已启用调试模式")
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	logging.Info("MediaWarp 监听端口：", config.Port)
	ginR := router.InitRouter() // 路由初始化
	logging.Info("MediaWarp 启动成功")
	ginR.Run(config.ListenAddr()) // 启动服务
}

// 打印LOGO
func printLOGO() {
	fmt.Print(
		constants.LOGO,
		utils.Center(
			fmt.Sprintf(" MediaWarp %s 启动中 ", config.Version().AppVersion),
			75,
			"=",
		),
		"\n\n",
	)
}
