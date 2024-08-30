package pkg

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 多前缀注册路由
func RegisterRoutesWithPrefixs(router *gin.Engine, rawPath string, handler gin.HandlerFunc, method string, prefixs ...string) {
	for _, prefix := range prefixs {
		path := prefix + rawPath
		switch method {
		case http.MethodGet:
			router.GET(path, handler)
		case http.MethodPost:
			router.POST(path, handler)
		}
	}
}
