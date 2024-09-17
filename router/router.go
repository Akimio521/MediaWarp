package router

import (
	"MediaWarp/constants"
	"MediaWarp/handlers/handlers_emby"
	"MediaWarp/middleware"
	"MediaWarp/resources"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	ginR := gin.New()
	ginR.Use(middleware.QueryCaseInsensitive())
	ginR.Use(middleware.Logger())
	ginR.Use(middleware.ClientFilter())

	mediawarpRouter := ginR.Group("/MediaWarp")
	{
		if config.Web.Enable && config.Web.Static {
			mediawarpRouter.Static("/Static", config.StaticDir())
		}
		mediawarpRouter.StaticFS("/Resources", resources.ResourcesFS)
	}

	ginR.NoRoute(RegexpRouterHandler)
	return ginR
}

// 正则表达式路由处理器
//
// 依次匹配路由规则, 找到对应的处理器
func RegexpRouterHandler(ctx *gin.Context) {
	var (
		rules      []RegexpRouterPair
		serverType = config.Server.GetType()
	)
	switch serverType {
	case constants.EMBY:
		rules = embyRouterRules
	default:
		panic("未识别服务器类型：" + serverType)
	}
	for _, rule := range rules {
		if rule.Regexp.MatchString(ctx.Request.RequestURI) {
			rule.Handler(ctx)
			return
		}
	}
	// 未匹配路由
	handlers_emby.DefaultHandler(ctx)
}
