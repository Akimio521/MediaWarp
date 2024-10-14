package router

import (
	"MediaWarp/constants"
	"MediaWarp/internal/assets"
	"MediaWarp/internal/handler"
	"MediaWarp/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	ginR := gin.New()
	ginR.Use(middleware.QueryCaseInsensitive())
	ginR.Use(middleware.SetRefererPolicy(constants.SAME_ORIGIN))
	ginR.Use(middleware.Logger())
	ginR.Use(middleware.ClientFilter())

	if cfg.Cache.WebCache {
		ginR.Use(middleware.Cache())
		logger.ServiceLogger.Info("Web缓存中间件已启用")
	} else {
		logger.ServiceLogger.Info("Web缓存中间件未启用")
	}

	mediawarpRouter := ginR.Group("/MediaWarp")
	{
		static := mediawarpRouter.Group("/static")
		{
			if cfg.Web.Enable {
				static.StaticFS("/embedded", http.FS(assets.EmbeddedStaticAssets))
				if cfg.Web.Custom {
					static.Static("/custom", cfg.StaticDir())
				}
			}
		}
	}

	ginR.NoRoute(RegexpRouterHandler)
	return ginR
}

// 正则表达式路由处理器
//
// 从媒体服务器处理结构体中获取正则路由规则
// 依次匹配请求, 找到对应的处理器
func RegexpRouterHandler(ctx *gin.Context) {
	mediaServerHandler := handler.GetMediaServer()

	for _, rule := range mediaServerHandler.GetRegexpRouteRules() {
		if rule.Regexp.MatchString(ctx.Request.RequestURI) {
			rule.Handler(ctx)
			return
		}
	}

	// 未匹配路由
	mediaServerHandler.ReverseProxy(ctx.Writer, ctx.Request)
}
