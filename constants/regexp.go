package constants

import "regexp"

var EmbyRegexp = map[string]map[string]*regexp.Regexp{ // Emby 相关的正则表达式
	"router": {
		"VideosHandler":               regexp.MustCompile(`(?i)^(/emby)?/videos/(.*)/(stream|original)`),          // 普通视频处理接口匹配
		"ModifyBaseHtmlPlayerHandler": regexp.MustCompile(`(?i)^/web/modules/htmlvideoplayer/basehtmlplayer.js$`), // 修改 Web 的 basehtmlplayer.js
		"WebIndex":                    regexp.MustCompile(`^/web/index.html$`),                                    // Web 首页
		"PlaybackInfoHandler":         regexp.MustCompile(`(?i)^(/emby)?/Items/(.*)/PlaybackInfo`),                // 播放信息处理接口
	},
	"others": {
		"VideoRedirectReg": regexp.MustCompile(`(?i)^(/emby)?/videos/(.*)/stream/(.*)`), // 视频重定向匹配，统一视频请求格式
	},
}
