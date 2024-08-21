package controllers

import (
	"MediaWarp/config"
	"MediaWarp/pkg"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

var cfg *config.ConfigManager = config.GetConfig()

// 默认路由（直接转发请求到后端）
func DefaultHandler(ctx *gin.Context) {
	fmt.Println("默认规则，\t请求方法：", ctx.Request.Method, "\n请求路径：", ctx.Request.URL.Path, "\n请求参数：", ctx.Request.URL.RawQuery)
	hostName, _ := pkg.SplitHostPort(ctx.Request.Host)

	remote, _ := url.Parse(cfg.Origin)

	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		req.Header = ctx.Request.Header
		req.Host = hostName //remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
		req.URL.Path = ctx.Request.URL.Path
	}
	proxy.ServeHTTP(ctx.Writer, ctx.Request)

}
