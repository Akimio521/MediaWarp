package constants

import "regexp"

type EmbyRegexps struct {
	Router RouterRegexps
	Others OthersRegexps
}

type RouterRegexps struct {
	VideosHandler        *regexp.Regexp // 普通视频处理接口匹配
	ModifyBaseHtmlPlayer *regexp.Regexp // 修改 Web 的 basehtmlplayer.js
	ModifyIndex          *regexp.Regexp // Web 首页
	ModifyPlaybackInfo   *regexp.Regexp // 播放信息处理接口
	ModifySubtitles      *regexp.Regexp // 字幕处理接口
}

type OthersRegexps struct {
	VideoRedirectReg *regexp.Regexp // 视频重定向匹配，统一视频请求格式
}

var EmbyRegexp = &EmbyRegexps{
	Router: RouterRegexps{
		VideosHandler:        regexp.MustCompile(`(?i)^(/emby)?/videos/\d+/(stream|original)$`),
		ModifyBaseHtmlPlayer: regexp.MustCompile(`(?i)^/web/modules/htmlvideoplayer/basehtmlplayer.js$`),
		ModifyIndex:          regexp.MustCompile(`^/web/index.html$`),
		ModifyPlaybackInfo:   regexp.MustCompile(`(?i)^(/emby)?/Items/\d+/PlaybackInfo$`),
		ModifySubtitles:      regexp.MustCompile(`(?i)^(/emby)?/Videos/\d+/\w+/subtitles$`),
	},
	Others: OthersRegexps{
		VideoRedirectReg: regexp.MustCompile(`(?i)^(/emby)?/videos/(.*)/stream/(.*)`),
	},
}

type JellyfinRouterRegexps struct {
	VideosHandler      *regexp.Regexp // 普通视频处理接口匹配
	ModifyIndex        *regexp.Regexp // Web 首页
	ModifyPlaybackInfo *regexp.Regexp // 播放信息处理接口
	ModifySubtitles    *regexp.Regexp // 字幕处理接口
}
type JellyfinRegexps struct {
	Router JellyfinRouterRegexps
}

var JellyfinRegexp = &JellyfinRegexps{
	Router: JellyfinRouterRegexps{
		VideosHandler:      regexp.MustCompile(`/Videos/\w+/(stream|original)$`), // /Videos/813a630bcf9c3f693a2ec8c498f868d2/stream?
		ModifyIndex:        regexp.MustCompile(`^/web/$`),
		ModifyPlaybackInfo: regexp.MustCompile(`^/Items/\w+$`),
		ModifySubtitles:    regexp.MustCompile(`/Videos/\d+/\w+/subtitles$`),
	},
}
