package router

import (
	"MediaWarp/handlers/handlers_emby"
	"regexp"
)

// Emby路由初始化
func initEmbyRouterRule() []RegexpRouterPair {
	embyRouterRules := []RegexpRouterPair{
		{
			Regexp:  regexp.MustCompile(`(?i)^/.*videos/.*/(stream|original)`),
			Handler: handlers_emby.VideosHandler,
		},
	}

	if config.Web.Enable {
		if config.Web.Index || config.Web.Head || config.Web.ExternalPlayerUrl || config.Web.BeautifyCSS {
			embyRouterRules = append(embyRouterRules,
				RegexpRouterPair{
					Regexp:  regexp.MustCompile(`^/web/index.html$`),
					Handler: handlers_emby.IndexHandler,
				},
			)
		}
	}
	embyRouterRules = append(embyRouterRules,
		RegexpRouterPair{
			Regexp:  regexp.MustCompile(`^/web/modules/htmlvideoplayer/basehtmlplayer.js$`),
			Handler: handlers_emby.ModifyBaseHtmlPlayerHandler,
		},
	)

	return embyRouterRules
}
