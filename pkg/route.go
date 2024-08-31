package pkg

import (
	"github.com/gin-gonic/gin"
)

// 多前缀注册路由
func RegisterRoutesWithPrefixs(router *gin.Engine, rawPath string, handler gin.HandlerFunc, method string, prefixs ...string) {
	for _, prefix := range prefixs {
		path := prefix + rawPath

		router.Handle(method, rawPath, handler)
		router.Handle(method, path, handler)
	}

}
