package router

import (
	"MediaWarp/controllers"
	"MediaWarp/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	ginR := gin.New()
	ginR.Use(middleware.QueryCaseInsensitive())
	ginR.Use(middleware.LogMiddleware())

	// UserLibraryService
	registerRoutes(ginR, "/Users/:userId/Items", controllers.DefaultHandler, http.MethodGet)
	registerRoutes(ginR, "/Users/:userId/Items/Resume", controllers.DefaultHandler, http.MethodGet)
	registerRoutes(ginR, "/Users/:userId/Items/Latest", controllers.DefaultHandler, http.MethodGet)
	registerRoutes(ginR, "/Users/:userId/Items/Views", controllers.DefaultHandler, http.MethodGet)
	registerRoutes(ginR, "/Users/:userId/Items/:itemId", controllers.ItemsHandler, http.MethodGet)
	// ItemsService
	registerRoutes(ginR, "/Items/:itemId/PlaybackInfo", controllers.PlaybackInfoHandler, http.MethodGet)
	registerRoutes(ginR, "/Items/:itemId/PlaybackInfo", controllers.PlaybackInfoHandler, http.MethodPost)
	// VideoService
	registerRoutes(ginR, "/Videos/:itemId/:name", controllers.VideosHandler, http.MethodGet)

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
