package handler

import (
	"regexp"

	"github.com/gin-gonic/gin"
)

// 正则表达式路由规则
type RegexpRouteRule struct {
	Regexp  *regexp.Regexp
	Handler gin.HandlerFunc
}
