package router

import (
	"MediaWarp/controllers"
	"MediaWarp/core"
	"MediaWarp/middleware"
	"MediaWarp/pkg"
	"MediaWarp/resources"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	config := core.GetConfig()

	ginR := gin.New()
	ginR.Use(middleware.QueryCaseInsensitive())
	ginR.Use(middleware.LogMiddleware())
	ginR.Use(middleware.ClientFilter())

	// VideoService
	pkg.RegisterRoutesWithPrefixs(ginR, "/Videos/:itemId/:name", controllers.VideosHandler, http.MethodGet, "/emby")
	pkg.RegisterRoutesWithPrefixs(ginR, "/videos/:itemId/:name", controllers.VideosHandler, http.MethodGet, "/emby")

	if config.Web.Enable {
		ginR.GET("/web/index.html", controllers.IndexHandler)
		if config.Web.Static { // 自定义静态资源
			ginR.Static("/MediaWarp/Static", config.StaticDir())
		}
	}
	ginR.StaticFS("/MediaWarp/Resources", resources.ResourcesFS)
	ginR.GET("/web/modules/htmlvideoplayer/basehtmlplayer.js", controllers.ModifyBaseHtmlPlayerHandler)
	ginR.NoRoute(controllers.DefaultHandler)

	return ginR
}
