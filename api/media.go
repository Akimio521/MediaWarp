package api

import (
	"MediaWarp/constants"
	"errors"
)

type MediaServer interface {
	Init()
	GetType() constants.MediaServerType
	GetEndpoint() string // 包含协议、服务器域名（IP）、端口号
	GetToken() string
}

type NewMediaServer func(string, string) MediaServer

// 通过名字获取MediaServer
func GetMediaServer(serverType constants.MediaServerType, addr string, token string) (MediaServer, error) {
	newFunc, ok := MediaServerMap[serverType]
	if !ok {
		return nil, errors.New("未知媒体服务器类型：" + string(serverType))
	}
	mediaserver := newFunc(addr, token)
	mediaserver.Init()
	return mediaserver, nil
}

// 注册可用媒体服务器
func registerMediaServerMap(serverType constants.MediaServerType, newFunc NewMediaServer) {
	MediaServerMap[serverType] = newFunc
}
