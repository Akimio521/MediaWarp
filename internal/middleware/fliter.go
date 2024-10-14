package middleware

import (
	"MediaWarp/constants"
	"strings"

	"github.com/gin-gonic/gin"
)

// 客户端过滤器
func ClientFilter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userAgent := ctx.Request.UserAgent()
		var allowed bool
		if userAgent == "" { // 开启了客户端过滤器后禁止所有未提供User-Agent的链接
			allowed = false
		} else {
			if cfg.ClientFilter.Mode == constants.WHITELIST { // 白名单模式
				allowed = false
				for _, ua := range cfg.ClientFilter.ClientList {
					if strings.Contains(userAgent, ua) {
						allowed = true
						break
					}
				}
			} else if cfg.ClientFilter.Mode == constants.BLACKLIST { // 黑名单模式
				allowed = true
				for _, ua := range cfg.ClientFilter.ClientList {
					if strings.Contains(userAgent, ua) {
						allowed = false
						break
					}
				}
			} else {
				logger.ServiceLogger.Error("未知的客户端过滤器模式，已关闭客户端过滤器")
				cfg.ClientFilter.Enable = false
				allowed = true
			}
		}

		if !allowed {
			ctx.AbortWithStatus(403) // 禁止访问
			logger.ServiceLogger.Info("客户端过滤器拦截了请求，User-Agent: ", userAgent)
			return
		}
		logger.ServiceLogger.Debug("客户端过滤器放行了请求，User-Agent: ", userAgent)
		ctx.Next()

	}
}
