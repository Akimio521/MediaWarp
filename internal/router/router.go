package router

import (
	"MediaWarp/constants"
	"MediaWarp/internal/assets"
	"MediaWarp/internal/config"
	"MediaWarp/internal/handler"
	"MediaWarp/internal/logging"
	"MediaWarp/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	ginR := gin.New()
	ginR.Use(middleware.QueryCaseInsensitive())
	ginR.Use(middleware.SetRefererPolicy(constants.SAME_ORIGIN))
	ginR.Use(middleware.Logger())

	if config.ClientFilter.Enable {
		ginR.Use(middleware.ClientFilter())
		logging.Info("客户端过滤中间件已启用")
	} else {
		logging.Info("客户端过滤中间件未启用")
	}

	if config.Cache.WebCache {
		ginR.Use(middleware.Cache())
		logging.Info("Web 缓存中间件已启用")
	} else {
		logging.Info("Web 缓存中间件未启用")
	}

	mediawarpRouter := ginR.Group("/MediaWarp")
	{
		mediawarpRouter.Any("/version", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, config.Version())
		})
		static := mediawarpRouter.Group("/static")
		{
			if config.Web.Enable {
				static.StaticFS("/embedded", http.FS(assets.EmbeddedStaticAssets))
				if config.Web.Custom {
					static.Static("/custom", config.StaticDir())
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
		if rule.Regexp.MatchString(ctx.Request.RequestURI) { // 带有查询参数的字符串：/emby/Items/54/Images/Primary?maxWidth=600&tag=f66addf8af207bdc39cdb4dd56db0d0b&quality=90
			rule.Handler(ctx)
			return
		}
	}

	// 未匹配路由
	mediaServerHandler.ReverseProxy(ctx.Writer, ctx.Request)
}
