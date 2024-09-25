package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// 将请求查询参数的键转换为小写
func QueryCaseInsensitive() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取所有查询参数
		queryParams := ctx.Request.URL.Query()

		// 创建一个新的 map 来存储大小写不敏感的查询参数
		caseInsensitiveParams := make(map[string][]string)
		for key, values := range queryParams {
			// 	将查询参数的键转换为小写，并存储值
			caseInsensitiveParams[strings.ToLower(key)] = values
		}

		// 清空原始查询参数并设置新的大小写不敏感的查询参数
		ctx.Request.URL.RawQuery = ""
		q := ctx.Request.URL.Query()
		for key, values := range caseInsensitiveParams {
			for _, value := range values {
				q.Add(key, value)
			}
		}
		ctx.Request.URL.RawQuery = q.Encode()
	}
}
