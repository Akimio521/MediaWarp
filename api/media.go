package api

import "MediaWarp/constants"

type MediaServer interface {
	GetType() constants.MediaServerType
	GetEndpoint() string // 包含协议、服务器域名（IP）、端口号
	GetToken() string
}
