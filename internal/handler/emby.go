package handler

import (
	"MediaWarp/internal/service"
	"MediaWarp/internal/service/alist"
	"MediaWarp/internal/service/emby"
	"MediaWarp/pkg"
	"bytes"
	"io"
	"net/http"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Emby服务器处理器
type EmbyServerHandler struct {
	server *emby.EmbyServer
}

// 初始化
func (embyServerHandler *EmbyServerHandler) Init() {
	embyServerHandler.server = emby.New(cfg.MeidaServer.ADDR, cfg.MeidaServer.AUTH)
}

// 转发请求至上游服务器
func (embyServerHandler *EmbyServerHandler) ReverseProxy(rw http.ResponseWriter, req *http.Request) {
	embyServerHandler.server.ReverseProxy(rw, req)
}

// 正则路由表
func (embyServerHandler *EmbyServerHandler) GetRegexpRouteRules() []RegexpRouteRule {
	embyRouterRules := []RegexpRouteRule{
		{
			Regexp:  regexp.MustCompile(`(?i)^/.*videos/.*/(stream|original)`),
			Handler: embyServerHandler.VideosHandler,
		},
		{
			Regexp:  regexp.MustCompile(`^/web/modules/htmlvideoplayer/basehtmlplayer.js$`),
			Handler: embyServerHandler.ModifyBaseHtmlPlayerHandler,
		},
	}

	if cfg.Web.Enable {
		if cfg.Web.Index || cfg.Web.Head != "" || cfg.Web.ExternalPlayerUrl || cfg.Web.BeautifyCSS {
			embyRouterRules = append(embyRouterRules,
				RegexpRouteRule{
					Regexp:  regexp.MustCompile(`^/web/index.html$`),
					Handler: embyServerHandler.IndexHandler,
				},
			)
		}
	}
	return embyRouterRules
}

// 视频流处理器
//
// 支持播放本地视频、重定向HttpStrm、AlistStrm
func (embyServerHandler *EmbyServerHandler) VideosHandler(ctx *gin.Context) {
	// EmbyServer <= 4.8 ====> mediaSourceID = 343121
	// EmbyServer >= 4.9 ====> mediaSourceID = mediasource_31
	mediaSourceID := ctx.Query("mediasourceid")

	logger.ServiceLogger.Debug("请求ItemsServiceQueryItem：", mediaSourceID)
	itemResponse, err := embyServerHandler.server.ItemsServiceQueryItem(strings.Replace(mediaSourceID, "mediasource_", "", 1), 1, "Path,MediaSources") // 查询item需要去除前缀仅保留数字部分

	if err != nil {
		logger.ServiceLogger.Warning("请求ItemsServiceQueryItem失败：", err)
		return
	}
	item := itemResponse.Items[0]
	for _, mediasource := range item.MediaSources {
		if *mediasource.ID == mediaSourceID { // EmbyServer >= 4.9 返回的ID带有前缀mediasource_
			redirectURL := getRedirctURL(&mediasource, &item)
			if redirectURL != "" {
				ctx.Redirect(http.StatusFound, redirectURL)
				return
			}

			logger.ServiceLogger.Info("本地视频：", *mediasource.Path)
			embyServerHandler.server.ReverseProxy(ctx.Writer, ctx.Request)
			return
		}
	}
}

// 修改basehtmlplayer.js
//
// 用于修改播放器JS，实现跨域播放Strm
func (embyServerHandler *EmbyServerHandler) ModifyBaseHtmlPlayerHandler(ctx *gin.Context) {
	version := ctx.Query("v")
	logger.ServiceLogger.Debug("请求basehtmlplayer.js版本：", version)

	proxy := embyServerHandler.server.GetReverseProxy()
	proxy.ModifyResponse = func(rw *http.Response) error {
		defer rw.Body.Close()
		body, err := io.ReadAll(rw.Body)
		if err != nil {
			return err
		}

		modifiedBodyStr := strings.ReplaceAll(string(body), `mediaSource.IsRemote&&"DirectPlay"===playMethod?null:"anonymous"`, "null") // 修改响应体
		updateBody(rw, modifiedBodyStr)
		return nil
	}

	proxy.ServeHTTP(ctx.Writer, ctx.Request)
}

// 首页处理方法
func (embyServerHandler *EmbyServerHandler) IndexHandler(ctx *gin.Context) {
	var (
		htmlFilePath    string = path.Join(cfg.StaticDir(), "index.html")
		modifiedBodyStr string
		addHEAD         string
	)
	if !cfg.Web.Enable { //直接转发相关请求
		embyServerHandler.server.ReverseProxy(ctx.Writer, ctx.Request)
	} else { // 修改请求
		proxy := embyServerHandler.server.GetReverseProxy()
		proxy.ModifyResponse = func(rw *http.Response) error {

			if !cfg.Web.Index { // 从上游获取响应体
				body, err := io.ReadAll(rw.Body)
				defer rw.Body.Close()
				if err != nil {
					return err
				}
				modifiedBodyStr = string(body)
			} else { // 从本地文件读取index.html
				htmlContent, err := pkg.GetFileContent(htmlFilePath)
				if err != nil {
					logger.ServiceLogger.Warning("读取文件内容出错，错误信息：", err)
					return err
				} else {
					modifiedBodyStr = string(htmlContent)
				}
			}

			if cfg.Web.Head != "" { // 用户自定义HEAD
				addHEAD += cfg.Web.Head + "\n"
			}
			if cfg.Web.ExternalPlayerUrl { // 外部播放器
				addHEAD += `<script src="/MediaWarp/static/embedded/js/ExternalPlayerUrl.js"></script>` + "\n"
			}
			if cfg.Web.ActorPlus { // 过滤没有头像的演员和制作人员
				addHEAD += `<script src="/MediaWarp/static/embedded/js/ActorPlus.js"></script>` + "\n"
			}
			if cfg.Web.FanartShow { // 显示同人图（fanart图）
				addHEAD += `<script src="/MediaWarp/static/embedded/js/FanartShow.js"></script>` + "\n"
			}
			if cfg.Web.Danmaku { // 弹幕
				addHEAD += `<script src="https://cdn.jsdelivr.net/gh/RyoLee/emby-danmaku@gh-pages/ede.user.js" defer></script>` + "\n"
			}

			if cfg.Web.BeautifyCSS { // 美化CSS
				addHEAD += `<link rel="stylesheet" href="/MediaWarp/static/embedded/css/Beautify.css" type="text/css" media="all" />` + "\n"
			}
			modifiedBodyStr = strings.Replace(modifiedBodyStr, "</head>", addHEAD+"</head>", 1) // 将添加HEAD
			updateBody(rw, modifiedBodyStr)
			return nil
		}
		proxy.ServeHTTP(ctx.Writer, ctx.Request)
	}
}

// 获取302重定向URL
func getRedirctURL(mediasource *emby.MediaSourceInfo, item *emby.BaseItemDto) (redirectURL string) {
	if *mediasource.Protocol == emby.HTTP && cfg.HTTPStrm.Enable { // 判断是否为http协议Strm
		httpStrmRedirectURL := getHttpStrmRedirectURL(mediasource, item)
		if httpStrmRedirectURL != "" {
			return httpStrmRedirectURL
		}
	}
	if strings.ToUpper(*mediasource.Container) == "STRM" { // 判断是否为Strm文件
		if cfg.AlistStrm.Enable { // 判断是否启用AlistStrm
			alistStrmRedirectURL := getAlistStrmRedirect(mediasource, item)
			if alistStrmRedirectURL != "" {
				return alistStrmRedirectURL
			}
		}
	}
	return ""
}

// Strm文件内部为标准http协议时，获取302重定向URL
func getHttpStrmRedirectURL(mediasource *emby.MediaSourceInfo, item *emby.BaseItemDto) (url string) {
	for _, prefix := range cfg.HTTPStrm.PrefixList {
		if strings.HasPrefix(*item.Path, prefix) {
			logger.ServiceLogger.Debug(*item.Path, "匹配HttpStrm路径：", prefix, "成功")
			logger.ServiceLogger.Info("Http协议Strm重定向：", *mediasource.Path)
			return *mediasource.Path
		}
	}
	logger.ServiceLogger.Info("未匹配HttpStrm路径：", *item.Path)
	return ""
}

// 判断是否注册为AlistStrm，获取302重定向URL
func getAlistStrmRedirect(mediasource *emby.MediaSourceInfo, item *emby.BaseItemDto) (url string) {
	var (
		alistPath   string
		alistServer *alist.AlistServer
	)
	for _, alistStrmConfig := range cfg.AlistStrm.List {
		for _, perfix := range alistStrmConfig.PrefixList {
			if strings.HasPrefix(*item.Path, perfix) {
				alistPath = *mediasource.Path                              // 获取Strm文件中的AlistPath
				alistServer = service.GetAlistServer(alistStrmConfig.ADDR) // 获取AlistServer，无需重新生成实例
				logger.ServiceLogger.Debug(*item.Path, "匹配AlistStrm路径：", perfix, "成功")
				break
			}
		}

	}
	if alistPath != "" { // 匹配成功
		fsGet, err := alistServer.FsGet(alistPath)
		if err != nil {
			logger.ServiceLogger.Warning("请求GetFile失败：", err)
			return ""
		}
		logger.ServiceLogger.Info("AlistStrm重定向：", fsGet.RawURL)
		return fsGet.RawURL
	}
	logger.ServiceLogger.Info("未匹配AlistStrm路径：", *item.Path)
	return ""
}

// 更新响应体
//
// 修改响应体、更新Content-Length
func updateBody(rw *http.Response, s string) {
	rw.Body = io.NopCloser(bytes.NewBuffer([]byte(s))) // 重置响应体

	// 更新 Content-Length 头
	rw.ContentLength = int64(len(s))
	rw.Header.Set("Content-Length", strconv.Itoa(len(s)))

}
