package main

import (
	"MediaWarp/config"
	"MediaWarp/router"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("MediaWarp启动中...")
	cfg := config.GetConfig()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	router := router.InitRouter() // 路由初始化

	server := &http.Server{ // http.Server 实例化
		Addr:    cfg.ListenAddr(), // 服务监听地址
		Handler: router,           // 服务处理器
	}
	ln, err := net.Listen("tcp", cfg.ListenAddr())
	if err != nil {
		fmt.Println("出现错误：", err)
		return
	}

	result := server.Serve(ln)
	if result != nil {
		fmt.Println("服务启动错误详情：", result)
	}

	<-sigChan
	fmt.Println("接收到中断信号，正在退出...")
}
