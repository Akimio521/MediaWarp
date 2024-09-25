package handler

import (
	"MediaWarp/constants"
	"fmt"
	"net/http"
)

type MediaServer interface {
	Init()                                           // 初始化函数
	ReverseProxy(http.ResponseWriter, *http.Request) // 转发请求至上游服务器
	GetRegexpRouteRules() []RegexpRouteRule          // 正则路由表
}

var mediaServer MediaServer

func init() {
	switch cfg.MeidaServer.Type {
	case constants.EMBY:
		mediaServer = &EmbyServer{}
	default:
		panic("错误媒体服务器类型")
	}
	fmt.Println("初始化媒体服务器接口")
	mediaServer.Init()
}

// 获取媒体服务器接口
func GetMediaServer() MediaServer {
	return mediaServer
}
