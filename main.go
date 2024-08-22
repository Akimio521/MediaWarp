package main

import (
	"MediaWarp/core"
	"MediaWarp/router"
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var config = core.GetConfig()
var logger = core.GetLogger()
var isDebug bool

func init() {
	flag.BoolVar(&isDebug, "debug", false, "是否启用调试模式")
}

func main() {
	flag.Parse()
	fmt.Printf("MediaWarp(%s)启动中...\n", config.Version())
	if isDebug {
		logger.ServerLogger.SetLevel(logrus.DebugLevel)
		logger.ServerLogger.Warning("已启用调试模式")
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	defer fmt.Println("MediaWarp正在关闭...")

	ginR := router.InitRouter()   // 路由初始化
	ginR.Run(config.ListenAddr()) // 启动服务
}
