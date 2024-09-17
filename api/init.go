package api

import "MediaWarp/constants"

var (
	MediaServerMap = map[constants.MediaServerType]NewMediaServer{}
)

func init() {
	registerMediaServerMap(constants.EMBY, registerEmbyServer)
}
