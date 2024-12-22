package constants_test

import (
	"MediaWarp/constants"
	"testing"
)

type RouteTestSet struct {
	URI    string
	Target string
}

var EmbyRouteTestSets = []RouteTestSet{
	{ // 字幕
		URI:    "/Videos/88697/21ed6a9972693ffa82571197cb406b64/Subtitles/3/0/Stream",
		Target: "",
	},
	{ // 视频1
		URI:    "/Videos/88697/stream?mediasourceid=21ed6a9972693ffa82571197cb406b64&static=true",
		Target: "VideosHandler",
	},
	{ // 视频2（增加前缀，修改大小写）
		URI:    "/emby/videos/88697/stream?mediasourceid=21ed6a9972693ffa82571197cb406b64&static=true",
		Target: "VideosHandler",
	},
	{ // PlaybackInfo
		URI:    "/Items/88697/PlaybackInfo?userid=9d882dc8ec514b2ca14652262df0afad",
		Target: "PlaybackInfoHandler",
	},
}

func TestEmbyRoute(t *testing.T) {
	for _, testSet := range EmbyRouteTestSets {
		for result, reg := range constants.EmbyRegexp["router"] {
			if reg.MatchString(testSet.URI) {
				if result != testSet.Target {
					t.Errorf("%s 路由错误。期望: %s, 实际: %s", testSet.URI, testSet.Target, result)
				}
			}
		}
	}
}
