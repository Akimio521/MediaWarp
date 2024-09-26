package handler

import (
	"MediaWarp/internal/service"
	"MediaWarp/internal/service/alist"
	"MediaWarp/internal/service/emby"
	"MediaWarp/pkg"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"regexp"
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
		if cfg.Web.Index || cfg.Web.Head || cfg.Web.ExternalPlayerUrl || cfg.Web.BeautifyCSS {
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

// 拦截basehtmlplayer.js实现Web跨域播放Strm
func (embyServerHandler *EmbyServerHandler) ModifyBaseHtmlPlayerHandler(ctx *gin.Context) {
	version := ctx.Query("v")
	logger.ServiceLogger.Info("请求basehtmlplayer.js版本：", version)
	resp, err := http.Get(embyServerHandler.server.GetEndpoint() + ctx.Request.URL.Path + "?" + ctx.Request.URL.RawQuery)
	if err != nil {
		logger.ServiceLogger.Warning("请求失败，使用回源策略，错误信息：", err)
		embyServerHandler.ReverseProxy(ctx.Writer, ctx.Request)
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.ServiceLogger.Warning("读取响应体失败，使用回源策略，错误信息：", err)
		embyServerHandler.ReverseProxy(ctx.Writer, ctx.Request)
		return
	}
	modifiedBody := strings.ReplaceAll(string(body), `mediaSource.IsRemote&&"DirectPlay"===playMethod?null:"anonymous"`, "null")

	for key, values := range resp.Header {
		for _, value := range values {
			if key != "Content-Length" {
				ctx.Writer.Header().Add(key, value)
			} else {
				ctx.Header("Content-Length", fmt.Sprintf("%d", len(modifiedBody)))
			}
		}
	}

	// 返回修改后的内容
	ctx.Data(resp.StatusCode, resp.Header.Get("Content-Type"), []byte(modifiedBody))

}

// 首页处理方法
func (embyServerHandler *EmbyServerHandler) IndexHandler(ctx *gin.Context) {
	var (
		htmlFilePath   string = filepath.Join(cfg.StaticDir(), "index.html")
		headFilePath   string = filepath.Join(cfg.StaticDir(), "head")
		isFile         bool
		err            error
		htmlContent    []byte
		headContent    []byte
		retHtmlContent string
	)

	if cfg.Web.Index { // 自定义首页
		isFile, err = pkg.IsFile(htmlFilePath)
		if err != nil {
			logger.ServiceLogger.Warning("判断路径是否为文件出错，错误信息：", err)
			isFile = false
		}

		if isFile {
			logger.ServiceLogger.Debug(htmlFilePath, "存在并且是文件")
			htmlContent, err = pkg.GetFileContent(htmlFilePath)
			if err != nil {
				logger.ServiceLogger.Warning("读取文件内容出错，使用回源策略，错误信息：", err)
			}
		}
	}
	if len(htmlContent) == 0 { // 未启用自定义首页或自定义首页文件内容读取失败
		logger.ServiceLogger.Debug("请求上游EmbyServer获取HTML内容")
		htmlContent, err = embyServerHandler.server.GetIndexHtml()
		if err != nil {
			logger.ServiceLogger.Warning("请求上游EmbyServer获取HTML内容失败，使用回源策略，错误信息：", err)
		}
	}

	if len(htmlContent) == 0 {
		embyServerHandler.ReverseProxy(ctx.Writer, ctx.Request)
		return
	}

	if cfg.Web.Head { // 自定义HEAD
		isFile, err = pkg.IsFile(headFilePath)
		if err != nil {
			logger.ServiceLogger.Warning("判断路径是否为文件出错，错误信息：", err)
			isFile = false
		}

		if isFile {
			headContent, err = pkg.GetFileContent(headFilePath)
			if err != nil {
				logger.ServiceLogger.Warning("读取文件内容出错，不添加额外HEAD，错误信息：", err)
			}
		} else {
			logger.ServiceLogger.Debug(headFilePath, "不存在或不是文件")
		}
	}

	if len(headContent) == 0 {
		retHtmlContent = string(htmlContent)
	} else {
		retHtmlContent = strings.Replace(string(htmlContent), "</head>", string(headContent)+"\n"+"</head>", 1)
	}

	if cfg.Web.ExternalPlayerUrl { // 外部播放器
		retHtmlContent = strings.Replace(retHtmlContent, "</head>", `<script src="/MediaWarp/static/embedded/js/ExternalPlayerUrl.js"></script>`+"\n"+"</head>", 1)
	}
	if cfg.Web.ActorPlus { // 过滤没有头像的演员和制作人员
		retHtmlContent = strings.Replace(retHtmlContent, "</head>", `<script src="/MediaWarp/static/embedded/js/ActorPlus.js"></script>`+"\n"+"</head>", 1)
	}
	if cfg.Web.FanartShow { // 显示同人图（fanart图）
		retHtmlContent = strings.Replace(retHtmlContent, "</head>", `<script src="/MediaWarp/static/embedded/js/FanartShow.js"></script>`+"\n"+"</head>", 1)
	}
	if cfg.Web.Danmaku { // 弹幕
		retHtmlContent = strings.Replace(retHtmlContent, "</body>", `<script src="https://cdn.jsdelivr.net/gh/RyoLee/emby-danmaku@gh-pages/ede.user.js" defer></script>`+"\n"+"</body>", 1)
	}

	if cfg.Web.BeautifyCSS { // 美化CSS
		retHtmlContent = strings.Replace(retHtmlContent, "</head>", `<link rel="stylesheet" href="/MediaWarp/static/embedded/css/Beautify.css" type="text/css" media="all" />`+"\n"+"</head>", 1)
	}

	ctx.Header("Content-Type", "text/html; charset=UTF-8")
	ctx.Header("Content-Length", fmt.Sprintf("%d", len(retHtmlContent)))
	ctx.Header("expires", "-1")
	ctx.String(http.StatusOK, retHtmlContent)
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
