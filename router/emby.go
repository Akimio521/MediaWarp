package router

import (
	"MediaWarp/handlers/handlers_emby"
	"MediaWarp/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
)

func initEmbyRouter(router *gin.Engine) {
	// VideoService
	pkg.RegisterRoutesWithPrefixs(router, "/Videos/:itemId/:name", handlers_emby.VideosHandler, http.MethodGet, "/emby")
	pkg.RegisterRoutesWithPrefixs(router, "/videos/:itemId/:name", handlers_emby.VideosHandler, http.MethodGet, "/emby")

	if config.Web.Enable {
		router.GET("/web/index.html", handlers_emby.IndexHandler)
		if config.Web.Static { // 自定义静态资源
			router.Static("/MediaWarp/Static", config.StaticDir())
		}
	}
	router.GET("/web/modules/htmlvideoplayer/basehtmlplayer.js", handlers_emby.ModifyBaseHtmlPlayerHandler)
}
