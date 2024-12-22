package constants_test

import (
	"MediaWarp/constants"
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
			"",
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
			"PlaybackInfoHandler",
		},
	}
	for caseName, testCase := range embyRouteTestCases {
		t.Run(caseName, func(t *testing.T) {
			for result, reg := range constants.EmbyRegexp["router"] {
				if reg.MatchString(testCase.URI) && result != testCase.Target { // 匹配但不相等
					t.Errorf("%s 路由错误。期望: %s, 实际: %s", caseName, testCase.Target, result)
				}
			}
		})
	}
}
