package middleware

import (
	"MediaWarp/constants"
	"MediaWarp/internal/config"
	"MediaWarp/internal/logger"
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
			if config.ClientFilter.Mode == constants.WHITELIST { // 白名单模式
				allowed = false
				for _, ua := range config.ClientFilter.ClientList {
					if strings.Contains(userAgent, ua) {
						allowed = true
						break
					}
				}
			} else if config.ClientFilter.Mode == constants.BLACKLIST { // 黑名单模式
				allowed = true
				for _, ua := range config.ClientFilter.ClientList {
					if strings.Contains(userAgent, ua) {
						allowed = false
						break
					}
				}
			} else {
				logger.Error("未知的客户端过滤器模式，已关闭客户端过滤器")
				config.ClientFilter.Enable = false
				allowed = true
			}
		}

		if !allowed {
			ctx.AbortWithStatus(403) // 禁止访问
			logger.Info("客户端过滤器拦截了请求，User-Agent: ", userAgent)
			return
		}
		logger.Debug("客户端过滤器放行了请求，User-Agent: ", userAgent)
		ctx.Next()

	}
}
