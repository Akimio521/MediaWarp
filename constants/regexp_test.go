package constants_test

import (
	"regexp"
	"testing"
)

func TestEmbyRoute(t *testing.T) {
	type RouteTestCase struct {
		URI    string
		Target string
	}
	var embyRouteTestCases = map[string]RouteTestCase{
		"字幕": {
			"/Videos/88697/21ed6a9972693ffa82571197cb406b64/Subtitles/3/0/Stream",
			"ModifySubtitles",
		},
		"4.9+字幕": {
			"/emby/Videos/45/mediasource_45/Subtitles/0/0/Stream.subrip?api_key=e12acc0815f74e9da6a86c9e8c2d45d8",
			"ModifySubtitles",
		},
		"4.9+字幕2": {
			"/emby/Videos/146/mediasource_146/Subtitles/3/0/Stream.srt?api_key=4b988503e747491ca53ff22527a13f08",
			"ModifySubtitles",
		},
		"视频1": {
			"/Videos/88697/stream?mediasourceid=21ed6a9972693ffa82571197cb406b64&static=true",
			"VideosHandler",
		},
		"视频2（增加前缀，修改大小写）": {
			"/emby/videos/88697/stream?mediasourceid=21ed6a9972693ffa82571197cb406b64&static=true",
			"VideosHandler",
		},
		"PlaybackInfo": {
			"/Items/88697/PlaybackInfo?userid=9d882dc8ec514b2ca14652262df0afad",
			"ModifyPlaybackInfo",
		},
		"WEB JavaScript": {
			"/web/videos/videos.js?v=4.8.10.0",
			"",
		},
	}
	var allReg = map[string]*regexp.Regexp{
		"VideosHandler":        regexp.MustCompile(`(?i)^(/emby)?/Videos/\d+/(stream|original)(\.\w+)?$`),
		"ModifyBaseHtmlPlayer": regexp.MustCompile(`(?i)^/web/modules/htmlvideoplayer/basehtmlplayer.js$`),
		"ModifyIndex":          regexp.MustCompile(`^/web/index.html$`),
		"ModifyPlaybackInfo":   regexp.MustCompile(`(?i)^(/emby)?/Items/\d+/PlaybackInfo$`),
		"ModifySubtitles":      regexp.MustCompile(`(?i)^(/emby)?/Videos/\d+/\w+/subtitles$`),
	}
	for caseName, testCase := range embyRouteTestCases {
		t.Run(caseName, func(t *testing.T) {
			for result, reg := range allReg {
				if reg.MatchString(testCase.URI) && result != testCase.Target { // 匹配但不相等
					t.Errorf("%s 路由错误。期望: %s, 实际: %s", caseName, testCase.Target, result)
				}
			}
		})
	}
}

func TestJellyfinRoute(t *testing.T) {
	type RouteTestCase struct {
		URI    string
		Target string
	}
	var embyRouteTestCases = map[string]RouteTestCase{
		"字幕": {
			"/Videos/88697/21ed6a9972693ffa82571197cb406b64/Subtitles/3/0/Stream",
			"ModifySubtitles",
		},
		"4.9+字幕": {
			"/emby/Videos/45/mediasource_45/Subtitles/0/0/Stream.subrip?api_key=e12acc0815f74e9da6a86c9e8c2d45d8",
			"ModifySubtitles",
		},
		"4.9+字幕2": {
			"/emby/Videos/146/mediasource_146/Subtitles/3/0/Stream.srt?api_key=4b988503e747491ca53ff22527a13f08",
			"ModifySubtitles",
		},
		"视频1": {
			"/Videos/99ef1b42a4dbb5c7695211cca3e49007/stream?Static=true&mediaSourceId=99ef1b42a4dbb5c7695211cca3e49007&Tag=e9d6bd3c4ec166ece59f25dff839eb9c&api_key=13d4f3d5c13542749a95e01974931792",
			"VideosHandler",
		},
		"视频2（增加前缀，修改大小写）": {
			"/emby/videos/88697/stream?mediasourceid=21ed6a9972693ffa82571197cb406b64&static=true",
			"VideosHandler",
		},
		"PlaybackInfo": {
			"/emby/Items/99ef1b42a4dbb5c7695211cca3e49007/PlaybackInfo?AutoOpenLiveStream=true&IsPlayback=true&UserId=f5b381091c3b430a85b66ebfe5cb5aa1",
			"ModifyPlaybackInfo",
		},
		"WEB JavaScript": {
			"/web/videos/videos.js?v=4.8.10.0",
			"",
		},
	}

	var allReg = map[string]*regexp.Regexp{
		"VideosHandler":      regexp.MustCompile(`(?i)^(/emby)?/Videos/\w+/(stream|original)(\.\w+)?$`), // /Videos/813a630bcf9c3f693a2ec8c498f868d2/stream /Videos/205953b114bb8c9dc2c7ba7e44b8024c/stream.mp4
		"ModifyIndex":        regexp.MustCompile(`(?i)^(/emby)?/web/$`),
		"ModifyPlaybackInfo": regexp.MustCompile(`(?i)^(/emby)?/Items/\w+/PlaybackInfo$`),
		"ModifySubtitles":    regexp.MustCompile(`(?i)^(/emby)?/Videos/\d+/\w+/subtitles$`),
	}
	for caseName, testCase := range embyRouteTestCases {
		t.Run(caseName, func(t *testing.T) {
			for result, reg := range allReg {
				if reg.MatchString(testCase.URI) && result != testCase.Target { // 匹配但不相等
					t.Errorf("%s 路由错误。期望: %s, 实际: %s", caseName, testCase.Target, result)
				}
			}
		})
	}
}
