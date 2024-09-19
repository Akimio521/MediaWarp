package cmd

import (
	"MediaWarp/router"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// MidaWarp启动函数
func SetUP() {
	if isDebug {
		logger.ServerLogger.SetLevel(logrus.DebugLevel)
		logger.ServerLogger.Warning("已启用调试模式")
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	ginR := router.InitRouter() // 路由初始化
	logger.ServerLogger.Info("MediaWarp启动成功")
	ginR.Run(config.ListenAddr()) // 启动服务
}
