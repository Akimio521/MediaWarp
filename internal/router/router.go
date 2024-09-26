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
	ginR.Use(middleware.Cache())

	mediawarpRouter := ginR.Group("/MediaWarp")
	{
		static := mediawarpRouter.Group("/static")
		{
			if cfg.Web.Enable && cfg.Web.Static {
				static.Static("/custom ", cfg.StaticDir())
			}
			static.StaticFS("/embedded", http.FS(assets.EmbeddedStaticAssets))
		}
	}

	ginR.NoRoute(RegexpRouterHandler)
	return ginR
}

// 正则表达式路由处理器
//
// 依次匹配路由规则, 找到对应的处理器
func RegexpRouterHandler(ctx *gin.Context) {
	mediaServer := handler.GetMediaServer()

	for _, rule := range mediaServer.GetRegexpRouteRules() {
		if rule.Regexp.MatchString(ctx.Request.RequestURI) {
			rule.Handler(ctx)
			return
		}
	}

	// 未匹配路由
	mediaServer.ReverseProxy(ctx.Writer, ctx.Request)
}
