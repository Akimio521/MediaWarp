package cmd

import (
	"MediaWarp/internal/config"
	"MediaWarp/internal/logger"
	"MediaWarp/internal/router"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// MidaWarp启动函数
func SetUP() {
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
