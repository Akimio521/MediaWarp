package middleware

import (
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// 将请求查询参数的键转换为小写
func QueryCaseInsensitive() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		queryParams := make(url.Values)
		for key, values := range ctx.Request.URL.Query() {
			queryParams.Add(strings.ToLower(key), strings.Join(values, ","))
		}
		ctx.Request.URL.RawQuery = queryParams.Encode()
	}
}
