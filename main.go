package main

import (
	"MediaWarp/core"
	"MediaWarp/router"
	"fmt"
)

var config = core.GetConfig()

func main() {
	fmt.Printf("MediaWarp(%s)启动中...\n", config.Version())
	defer fmt.Println("MediaWarp正在关闭...")

	ginR := router.InitRouter()   // 路由初始化
	ginR.Run(config.ListenAddr()) // 启动服务
}
