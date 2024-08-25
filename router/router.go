package router

import (
	"MediaWarp/controllers"
	"MediaWarp/core"
	"MediaWarp/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	var config = core.GetConfig()
	ginR := gin.New()
	ginR.Use(middleware.QueryCaseInsensitive())
	ginR.Use(middleware.LogMiddleware())
	ginR.Use(middleware.ClientFilter())

	// VideoService
	registerRoutes(ginR, "/Videos/:itemId/:name", controllers.VideosHandler, http.MethodGet)
	registerRoutes(ginR, "/videos/:itemId/:name", controllers.VideosHandler, http.MethodGet)

	ginR.GET("/web/modules/htmlvideoplayer/basehtmlplayer.js", controllers.ModifyBaseHtmlPlayerHandler)
	if config.Static { // 静态资源
		ginR.GET("/web/index.html", controllers.ModifyIndexHandler)
		ginR.Static("/MediaWarp/Static", config.StaticDir())
	}
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
