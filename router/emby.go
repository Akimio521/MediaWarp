package router

import (
	"MediaWarp/handlers"
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
		if config.Web.Index || config.Web.Head || config.Web.ExternalPlayerUrl || config.Web.BeautifyCSS {
			pkg.RegisterRoutesWithPrefixs(router, "/web/index.html", handlers_emby.IndexHandler, http.MethodGet)
		}
	}
	router.GET("/web/modules/htmlvideoplayer/basehtmlplayer.js", handlers.DefaultHandler)
}
