package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				// 忽略 http.ErrAbortHandler 的 panic
				if r == http.ErrAbortHandler {
					return
				}

				// 处理其他 panic
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Internal Server Error",
					"msg":   r,
				})
			}
		}()
		c.Next()
	}
}
