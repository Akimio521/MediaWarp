package router

import (
	"MediaWarp/controllers"
	"MediaWarp/core"
	"MediaWarp/middleware"
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
	registerRoutes(ginR, "/Videos/:itemId/:name", controllers.VideosHandler, http.MethodGet)
	registerRoutes(ginR, "/videos/:itemId/:name", controllers.VideosHandler, http.MethodGet)

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

// 注册路由
func registerRoutes(router *gin.Engine, path string, handler gin.HandlerFunc, method string) {
	embyPath := "/emby" + path
	switch method {
	case http.MethodGet:
		router.GET(path, handler)
		router.GET(embyPath, handler)
	case http.MethodPost:
		router.POST(path, handler)
		router.POST(embyPath, handler)
	}
}
