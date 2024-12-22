package handler

import (
	"MediaWarp/internal/config"
	"MediaWarp/internal/logging"
	"MediaWarp/internal/service"
	"MediaWarp/internal/service/emby"
	"MediaWarp/utils"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var embyRegexp = map[string]*regexp.Regexp{ // Emby 相关的正则表达式
	"VideosHandler":               regexp.MustCompile(`(?i)^(.*)/videos/(.*)/(stream|original)`),              // 普通视频处理接口匹配
	"ModifyBaseHtmlPlayerHandler": regexp.MustCompile(`(?i)^/web/modules/htmlvideoplayer/basehtmlplayer.js$`), // 修改 Web 的 basehtmlplayer.js
	"videoRedirectReg":            regexp.MustCompile(`(?i)^(.*)/videos/(.*)/stream/(.*)$`),                   // 视频重定向匹配，统一视频请求格式
}

// Emby服务器处理器
type EmbyServerHandler struct {
	server *emby.EmbyServer
}

// 初始化
func (embyServerHandler *EmbyServerHandler) Init() {
	embyServerHandler.server = emby.New(config.MediaServer.ADDR, config.MediaServer.AUTH)
}

// 转发请求至上游服务器
func (embyServerHandler *EmbyServerHandler) ReverseProxy(rw http.ResponseWriter, req *http.Request) {
	embyServerHandler.server.ReverseProxy(rw, req)
}

// 正则路由表
func (embyServerHandler *EmbyServerHandler) GetRegexpRouteRules() []RegexpRouteRule {
	embyRouterRules := []RegexpRouteRule{
		{
			Regexp:  embyRegexp["VideosHandler"],
			Handler: embyServerHandler.VideosHandler,
		},
		{
			Regexp:  embyRegexp["ModifyBaseHtmlPlayerHandler"],
			Handler: embyServerHandler.ModifyBaseHtmlPlayerHandler,
		},
	}

	if config.Web.Enable {
		if config.Web.Index || config.Web.Head != "" || config.Web.ExternalPlayerUrl || config.Web.BeautifyCSS {
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
	orginalPath := ctx.Request.URL.Path
	matches := embyRegexp["videoRedirectReg"].FindStringSubmatch(orginalPath)
	if len(matches) == 3 {
		if strings.ToLower(matches[0]) == "emby" {
			redirectPath := fmt.Sprintf("/videos/%s/stream", matches[1])
			logging.Debug(orginalPath + " 重定向至：" + redirectPath)
			ctx.Redirect(http.StatusFound, redirectPath)
			return
		}
	}
	if len(matches) == 2 {
		redirectPath := fmt.Sprintf("/videos/%s/stream", matches[0])
		logging.Debug(orginalPath + " 重定向至：" + redirectPath)
		ctx.Redirect(http.StatusFound, redirectPath)
		return
	}
	if ctx.Request.Method == http.MethodHead { // 不额外处理 HEAD 请求
		embyServerHandler.ReverseProxy(ctx.Writer, ctx.Request)
		logging.Debug("VideosHandler 不处理 HEAD 请求，转发至上游服务器")
		return
	}

	// EmbyServer <= 4.8 ====> mediaSourceID = 343121
	// EmbyServer >= 4.9 ====> mediaSourceID = mediasource_31
	mediaSourceID := ctx.Query("mediasourceid")

	logging.Debug("请求 ItemsServiceQueryItem：", mediaSourceID)
	itemResponse, err := embyServerHandler.server.ItemsServiceQueryItem(strings.Replace(mediaSourceID, "mediasource_", "", 1), 1, "Path,MediaSources") // 查询item需要去除前缀仅保留数字部分

	if err != nil {
		logging.Warning("请求 ItemsServiceQueryItem 失败：", err)
		embyServerHandler.server.ReverseProxy(ctx.Writer, ctx.Request)
		return
	}
	item := itemResponse.Items[0]
	for _, mediasource := range item.MediaSources {
		if *mediasource.ID == mediaSourceID { // EmbyServer >= 4.9 返回的ID带有前缀mediasource_
			isRedirect := true

			// HTTPStrm 处理
			if *mediasource.Protocol == emby.HTTP && config.HTTPStrm.Enable {
				for _, prefix := range config.HTTPStrm.PrefixList {
					if strings.HasPrefix(*item.Path, prefix) {
						logging.Debug(*item.Path, " 匹配 HTTPStrm 路径：", prefix, " 成功")
						logging.Info("HTTPStrm 重定向至：", *mediasource.Path)
						ctx.Redirect(http.StatusFound, *mediasource.Path)
						return
					}
				}
				logging.Info("未匹配 HTTPStrm 路径：", *item.Path)
				isRedirect = false // 未匹配到 HTTPStrm 路径，但未启用，不进行后续重定向
			}

			// AlistStrm 处理
			if isRedirect && strings.ToUpper(*mediasource.Container) == "STRM" && config.AlistStrm.Enable { // 判断是否为Strm文件
				for _, alistStrmConfig := range config.AlistStrm.List {
					for _, perfix := range alistStrmConfig.PrefixList {
						if strings.HasPrefix(*item.Path, perfix) {
							alistServer := service.GetAlistServer(alistStrmConfig.ADDR)
							fsGetData, err := alistServer.FsGet(*mediasource.Path)
							if err != nil {
								logging.Warning("请求 FsGet 失败：", err)
								return
							}
							logging.Info("AlistStrm 重定向至：", fsGetData.RawURL)
							ctx.Redirect(http.StatusFound, fsGetData.RawURL)
							return
						}
					}
				}
				logging.Info("未匹配 AlistStrm 路径：", *item.Path)
				isRedirect = false // 未匹配到 AlistStrm 路径，但未启用，不进行后续重定向
			}

			logging.Info("本地视频：", *mediasource.Path)
			embyServerHandler.server.ReverseProxy(ctx.Writer, ctx.Request)
			return
		}
	}
	logging.Debug("非视频流请求，转发至上游服务器") // 可能是字幕文件等
	embyServerHandler.server.ReverseProxy(ctx.Writer, ctx.Request)
}

// 修改basehtmlplayer.js
//
// 用于修改播放器JS，实现跨域播放Strm
func (embyServerHandler *EmbyServerHandler) ModifyBaseHtmlPlayerHandler(ctx *gin.Context) {
	version := ctx.Query("v")
	logging.Debug("请求 basehtmlplayer.js 版本：", version)

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
		htmlFilePath    string = path.Join(config.StaticDir(), "index.html")
		modifiedBodyStr string
		addHEAD         string
	)
	if !config.Web.Enable { //直接转发相关请求
		embyServerHandler.server.ReverseProxy(ctx.Writer, ctx.Request)
	} else { // 修改请求
		proxy := embyServerHandler.server.GetReverseProxy()
		proxy.ModifyResponse = func(rw *http.Response) error {

			if !config.Web.Index { // 从上游获取响应体
				body, err := io.ReadAll(rw.Body)
				defer rw.Body.Close()
				if err != nil {
					return err
				}
				modifiedBodyStr = string(body)
			} else { // 从本地文件读取index.html
				htmlContent, err := utils.GetFileContent(htmlFilePath)
				if err != nil {
					logging.Warning("读取文件内容出错，错误信息：", err)
					return err
				} else {
					modifiedBodyStr = string(htmlContent)
				}
			}

			if config.Web.Head != "" { // 用户自定义HEAD
				addHEAD += config.Web.Head + "\n"
			}
			if config.Web.ExternalPlayerUrl { // 外部播放器
				addHEAD += `<script src="/MediaWarp/static/embedded/js/ExternalPlayerUrl.js"></script>` + "\n"
			}
			if config.Web.ActorPlus { // 过滤没有头像的演员和制作人员
				addHEAD += `<script src="/MediaWarp/static/embedded/js/ActorPlus.js"></script>` + "\n"
			}
			if config.Web.FanartShow { // 显示同人图（fanart图）
				addHEAD += `<script src="/MediaWarp/static/embedded/js/FanartShow.js"></script>` + "\n"
			}
			if config.Web.Danmaku { // 弹幕
				addHEAD += `<script src="https://cdn.jsdelivr.net/gh/RyoLee/emby-danmaku@gh-pages/ede.user.js" defer></script>` + "\n"
			}

			if config.Web.BeautifyCSS { // 美化CSS
				addHEAD += `<link rel="stylesheet" href="/MediaWarp/static/embedded/css/Beautify.css" type="text/css" media="all" />` + "\n"
			}
			modifiedBodyStr = strings.Replace(modifiedBodyStr, "</head>", addHEAD+"</head>", 1) // 将添加HEAD
			updateBody(rw, modifiedBodyStr)
			return nil
		}
		proxy.ServeHTTP(ctx.Writer, ctx.Request)
	}
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
