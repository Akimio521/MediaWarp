package main

import (
	"MediaWarp/constants"
	"MediaWarp/internal/config"
	"MediaWarp/internal/logger"
	"MediaWarp/internal/router"
	"MediaWarp/pkg"
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	isDebug bool
)

func init() {
	printLOGO()
	flag.BoolVar(&isDebug, "debug", false, "是否启用调试模式")
	flag.Parse()
}

func main() {
	if isDebug {
		logger.SetLevel(logrus.DebugLevel)
		logger.Warning("已启用调试模式")
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	logger.Info("MediaWarp 监听端口：", config.Port)
	ginR := router.InitRouter() // 路由初始化
	logger.Info("MediaWarp 启动成功")
	ginR.Run(config.ListenAddr()) // 启动服务
}

// 打印LOGO
func printLOGO() {
	fmt.Print(
		constants.LOGO,
		pkg.Center(
			fmt.Sprintf(" MediaWarp %s 启动中 ", config.Version()),
			75,
			"=",
		),
		"\n\n",
	)
}
