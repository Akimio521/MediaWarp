package api

import "MediaWarp/constants"

type MediaServer interface {
	GetType() constants.ServerType
	GetHTTPEndpoint() string
	GetWebSocketEndpoint() string
	GetToken() string
}
