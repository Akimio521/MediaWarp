package router

import (
	"MediaWarp/constants"
	"MediaWarp/handlers"
	"MediaWarp/middleware"
	"MediaWarp/resources"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	ginR := gin.New()
	ginR.Use(middleware.QueryCaseInsensitive())
	ginR.Use(middleware.Logger())
	ginR.Use(middleware.ClientFilter())

	switch config.Server.GetType() {
	case constants.EMBY:
		initEmbyRouter(ginR)
	default:
		panic("未知的媒体服务器类型")
	}

	if config.Web.Enable && config.Web.Static {
		ginR.Static("/MediaWarp/Static", config.StaticDir())
	}
	ginR.StaticFS("/MediaWarp/Resources", resources.ResourcesFS)
	ginR.NoRoute(handlers.DefaultHandler)

	return ginR
}
