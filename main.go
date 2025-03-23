package main

import (
	"MediaWarp/constants"
	"MediaWarp/internal/config"
	"MediaWarp/internal/handler"
	"MediaWarp/internal/logging"
	"MediaWarp/internal/router"
	"MediaWarp/internal/service"
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
	fmt.Print(constants.LOGO)
	fmt.Println(utils.Center(fmt.Sprintf(" MediaWarp %s ", config.Version().AppVersion), 71, "="))
	flag.BoolVar(&showVersion, "version", false, "显示版本信息")
	flag.BoolVar(&isDebug, "debug", false, "是否启用调试模式")
}

func main() {
	flag.Parse()

	if showVersion {
		ver, _ := json.MarshalIndent(config.Version(), "", "  ")
		fmt.Println(string(ver))
		return
	}

	if err := config.Init(); err != nil { // 初始化配置
		logging.Error("配置初始化失败：", err)
		return
	}
	logging.Init()           // 初始化日志
	service.InitAlistSerer() // 初始化Alist服务器
	handler.Init()           // 初始化媒体服务器处理器

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
