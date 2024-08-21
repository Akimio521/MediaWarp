package main

import (
	"MediaWarp/config"
	"MediaWarp/router"
	"fmt"
)

var cfg = config.GetConfig()

func main() {
	fmt.Printf("MediaWarp(%s)启动中...\n", cfg.Version())
	defer fmt.Println("MediaWarp正在关闭...")

	ginR := router.InitRouter() // 路由初始化
	ginR.Run(cfg.ListenAddr())  // 启动服务
}
