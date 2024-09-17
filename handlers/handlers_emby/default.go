package handlers_emby

import (
	"github.com/gin-gonic/gin"
)

// 默认路由（直接转发请求到后端）
func DefaultHandler(ctx *gin.Context) {
	embyServer.ReverseProxy(ctx.Writer, ctx.Request)
}
