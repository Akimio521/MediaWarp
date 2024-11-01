package handler

import (
	"MediaWarp/constants"
	"MediaWarp/internal/config"
	"net/http"
)

// 媒体服务器处理接口
type MediaServerHandler interface {
	Init()                                           // 初始化函数
	ReverseProxy(http.ResponseWriter, *http.Request) // 转发请求至上游服务器
	GetRegexpRouteRules() []RegexpRouteRule          // 获取正则路由表
}

var mediaServerHandler MediaServerHandler

func init() {
	switch config.MediaServer.Type {
	case constants.EMBY:
		mediaServerHandler = &EmbyServerHandler{}
	default:
		panic("错误媒体服务器类型")
	}
	mediaServerHandler.Init()
}

// 获取媒体服务器接口
func GetMediaServer() MediaServerHandler {
	return mediaServerHandler
}
