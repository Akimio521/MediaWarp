package controllers

import (
	"MediaWarp/core"
	"MediaWarp/pkg"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

var config = core.GetConfig()
var logger = core.GetLogger()

// 默认路由（直接转发请求到后端）
func DefaultHandler(ctx *gin.Context) {
	hostName, _ := pkg.SplitHostPort(ctx.Request.Host)

	remote, _ := url.Parse(config.Origin)

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
