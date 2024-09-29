package cmd

import (
	"MediaWarp/internal/router"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// MidaWarp启动函数
func SetUP() {
	if isDebug {
		logger.ServiceLogger.SetLevel(logrus.DebugLevel)
		fmt.Println("已启用调试模式")
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	ginR := router.InitRouter() // 路由初始化
	fmt.Println("MediaWarp启动成功")
	ginR.Run(cfg.ListenAddr()) // 启动服务
}
