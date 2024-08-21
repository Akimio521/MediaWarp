package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func LogRawRequest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fmt.Println("路径参数:", ctx.Request.URL.Path, "\t查询参数:", ctx.Request.URL.RawQuery, "\t请求方法:", ctx.Request.Method)
		ctx.Next()
	}
}
