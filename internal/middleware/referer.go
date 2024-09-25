package middleware

import (
	"MediaWarp/constants"

	"github.com/gin-gonic/gin"
)

// 设置Referer策略
func SetRefererPolicy(value constants.REFERER_VALUE) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Referrer-Policy", string(value))
	}
}
