package handler

import (
	"MediaWarp/constants"
	"MediaWarp/internal/config"
	"MediaWarp/internal/logging"
	"MediaWarp/internal/service"
	"MediaWarp/internal/service/emby"
	"MediaWarp/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"path"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Emby服务器处理器
type EmbyServerHandler struct {
	server         *emby.EmbyServer
	modifyProxyMap map[string]*httputil.ReverseProxy
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
			Regexp:  constants.EmbyRegexp["router"]["VideosHandler"],
			Handler: embyServerHandler.VideosHandler,
		},
		{
			Regexp:  constants.EmbyRegexp["router"]["PlaybackInfoHandler"],
			Handler: embyServerHandler.PlaybackInfoHandler,
		},
		{
			Regexp:  constants.EmbyRegexp["router"]["ModifyBaseHtmlPlayerHandler"],
			Handler: embyServerHandler.ModifyBaseHtmlPlayerHandler,
		},
	}

	if config.Web.Enable {
		if config.Web.Index || config.Web.Head != "" || config.Web.ExternalPlayerUrl || config.Web.BeautifyCSS {
			embyRouterRules = append(embyRouterRules,
				RegexpRouteRule{
					Regexp:  constants.EmbyRegexp["router"]["WebIndex"],
					Handler: embyServerHandler.IndexHandler,
				},
			)
		}
	}
	return embyRouterRules
}

// 根据 Strm 文件路径识别 Strm 文件类型
//
// 返回 Strm 文件类型和一个可选配置
func (embyServerHandler *EmbyServerHandler) RecgonizeStrmFileType(strmFilePath string) (constants.StrmFileType, any) {
	if config.HTTPStrm.Enable {
		for _, prefix := range config.HTTPStrm.PrefixList {
			if strings.HasPrefix(strmFilePath, prefix) {
				logging.Debug(strmFilePath + " 成功匹配路径：" + prefix + "，Strm 类型：" + string(constants.HTTPStrm))
				return constants.HTTPStrm, nil
			}
		}
	}
	if config.AlistStrm.Enable {
		for _, alistStrmConfig := range config.AlistStrm.List {
			for _, prefix := range alistStrmConfig.PrefixList {
				if strings.HasPrefix(strmFilePath, prefix) {
					logging.Debug(strmFilePath + " 成功匹配路径：" + prefix + "，Strm 类型：" + string(constants.AlistStrm) + "，AlistServer 地址：" + alistStrmConfig.ADDR)
					return constants.AlistStrm, alistStrmConfig.ADDR
				}
			}
		}
	}
	logging.Debug(strmFilePath + " 未匹配任何路径，Strm 类型：" + string(constants.UnknownStrm))
	return constants.UnknownStrm, nil
}

// 处理播放信息请求处理接口
//
// /Items/:itemId/PlaybackInfo
func (embyServerHandler *EmbyServerHandler) PlaybackInfoHandler(ctx *gin.Context) {
	key := "PlaybackInfoHandler"
	if embyServerHandler.modifyProxyMap == nil {
		embyServerHandler.modifyProxyMap = make(map[string]*httputil.ReverseProxy)
	}

	if _, ok := embyServerHandler.modifyProxyMap[key]; !ok {
		proxy := embyServerHandler.server.GetReverseProxy()
		proxy.ModifyResponse = func(rw *http.Response) error {
			defer rw.Body.Close()
			body, err := io.ReadAll(rw.Body)
			if err != nil {
				logging.Warning("读取 Body 出错：", err)
				return err
			}

			var playbackInfoResponse emby.PlaybackInfoResponse
			err = json.Unmarshal(body, &playbackInfoResponse)
			if err != nil {
				logging.Warning("解析 emby.PlaybackInfoResponse Json 错误：", err)
				return err
			}
			for index, mediasource := range playbackInfoResponse.MediaSources {
				logging.Debug("请求 ItemsServiceQueryItem：" + *mediasource.ID)
				itemResponse, err := embyServerHandler.server.ItemsServiceQueryItem(strings.Replace(*mediasource.ID, "mediasource_", "", 1), 1, "Path,MediaSources") // 查询 item 需要去除前缀仅保留数字部分
				if err != nil {
					logging.Warning("请求 ItemsServiceQueryItem 失败：", err)
					return err
				}
				item := itemResponse.Items[0]
				strmFileType, opt := embyServerHandler.RecgonizeStrmFileType(*item.Path)
				switch strmFileType {
				case constants.HTTPStrm: // HTTPStrm 设置支持直链播放并且支持转码
					*playbackInfoResponse.MediaSources[index].SupportsDirectPlay = true
					*playbackInfoResponse.MediaSources[index].SupportsDirectStream = true
					if mediasource.DirectStreamURL != nil {
						apikeypair, err := utils.ResolveEmbyAPIKVPairs(*mediasource.DirectStreamURL)
						if err != nil {
							logging.Warning("解析API键值对失败：", err)
						}
						directStreamURL := fmt.Sprintf("/videos/%s/stream?MediaSourceId=%s&Static=true&%s", *mediasource.ItemID, *mediasource.ID, apikeypair)
						playbackInfoResponse.MediaSources[index].DirectStreamURL = &directStreamURL
						logging.Debug("设置直链播放链接为: " + directStreamURL)
					}

				case constants.AlistStrm: // AlistStm 设置支持直链播放并且禁止转码
					*playbackInfoResponse.MediaSources[index].SupportsDirectPlay = true
					*playbackInfoResponse.MediaSources[index].SupportsDirectStream = true
					*playbackInfoResponse.MediaSources[index].SupportsTranscoding = false
					playbackInfoResponse.MediaSources[index].TranscodingURL = nil
					playbackInfoResponse.MediaSources[index].TranscodingSubProtocol = nil
					playbackInfoResponse.MediaSources[index].TranscodingContainer = nil
					apikeypair, err := utils.ResolveEmbyAPIKVPairs(*mediasource.DirectStreamURL)
					if err != nil {
						logging.Warning("解析API键值对失败：", err)
						continue
					}
					directStreamURL := fmt.Sprintf("/videos/%s/stream?MediaSourceId=%s&Static=true&%s", *mediasource.ItemID, *mediasource.ID, apikeypair)
					playbackInfoResponse.MediaSources[index].DirectStreamURL = &directStreamURL
					container := strings.TrimPrefix(path.Ext(*mediasource.Path), ".")
					playbackInfoResponse.MediaSources[index].Container = &container
					logging.Debug("设置直链播放链接为: " + directStreamURL + "，容器为: " + container)
					alistServer := service.GetAlistServer(opt.(string))
					fsGetData, err := alistServer.FsGet(*mediasource.Path)
					if err != nil {
						logging.Warning("请求 FsGet 失败：", err)
						continue
					}
					playbackInfoResponse.MediaSources[index].Size = &fsGetData.Size
					logging.Debug("设置文件大小为: " + strconv.FormatInt(fsGetData.Size, 10))
				}
			}

			body, err = json.Marshal(playbackInfoResponse)
			if err != nil {
				logging.Warning("序列化 emby.PlaybackInfoResponse Json 错误：", err)
				return err
			}

			rw.Body = io.NopCloser(bytes.NewBuffer(body)) // 重置响应体
			// 更新 Content-Length 头
			rw.ContentLength = int64(len(body))
			rw.Header.Set("Content-Length", strconv.Itoa(len(body)))
			// 更新 Content-Type 头
			rw.Header.Set("Content-Type", "application/json")

			return nil

		}
		embyServerHandler.modifyProxyMap[key] = proxy
	}
	embyServerHandler.modifyProxyMap[key].ServeHTTP(ctx.Writer, ctx.Request)
}

// 视频流处理器
//
// 支持播放本地视频、重定向HttpStrm、AlistStrm
func (embyServerHandler *EmbyServerHandler) VideosHandler(ctx *gin.Context) {
	orginalPath := ctx.Request.URL.Path
	matches := constants.EmbyRegexp["others"]["VideoRedirectReg"].FindStringSubmatch(orginalPath)
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
	itemResponse, err := embyServerHandler.server.ItemsServiceQueryItem(strings.Replace(mediaSourceID, "mediasource_", "", 1), 1, "Path,MediaSources") // 查询 item 需要去除前缀仅保留数字部分
	if err != nil {
		logging.Warning("请求 ItemsServiceQueryItem 失败：", err)
		embyServerHandler.server.ReverseProxy(ctx.Writer, ctx.Request)
		return
	}

	item := itemResponse.Items[0]

	if !strings.HasSuffix(strings.ToLower(*item.Path), ".strm") { // 不是 Strm 文件
		logging.Debug("播放本地视频：" + *item.Path + "，不进行处理")
		embyServerHandler.server.ReverseProxy(ctx.Writer, ctx.Request)
		return
	}

	strmFileType, opt := embyServerHandler.RecgonizeStrmFileType(*item.Path)
	for _, mediasource := range item.MediaSources {
		if *mediasource.ID == mediaSourceID { // EmbyServer >= 4.9 返回的ID带有前缀mediasource_
			switch strmFileType {
			case constants.HTTPStrm:
				if *mediasource.Protocol == emby.HTTP {
					logging.Info("HTTPStrm 重定向至：", *mediasource.Path)
					ctx.Redirect(http.StatusFound, *mediasource.Path)
				}
				return
			case constants.AlistStrm:
				if strings.ToUpper(*mediasource.Container) == "STRM" { // 判断是否为Strm文件
					alistServerAddr := opt.(string)
					alistServer := service.GetAlistServer(alistServerAddr)
					fsGetData, err := alistServer.FsGet(*mediasource.Path)
					if err != nil {
						logging.Warning("请求 FsGet 失败：", err)
						return
					}
					logging.Info("AlistStrm 重定向至：", fsGetData.RawURL)
					ctx.Redirect(http.StatusFound, fsGetData.RawURL)
					return
				}
				return
			case constants.UnknownStrm:
				embyServerHandler.server.ReverseProxy(ctx.Writer, ctx.Request)
				return
			}
		}
	}
}

// 修改basehtmlplayer.js
//
// 用于修改播放器JS，实现跨域播放Strm
func (embyServerHandler *EmbyServerHandler) ModifyBaseHtmlPlayerHandler(ctx *gin.Context) {
	key := "ModifyBaseHtmlPlayerHandler"

	version := ctx.Query("v")
	logging.Debug("请求 basehtmlplayer.js 版本：", version)

	if embyServerHandler.modifyProxyMap == nil {
		embyServerHandler.modifyProxyMap = make(map[string]*httputil.ReverseProxy)
	}

	if _, ok := embyServerHandler.modifyProxyMap[key]; !ok {
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
		embyServerHandler.modifyProxyMap[key] = proxy
	}
	embyServerHandler.modifyProxyMap[key].ServeHTTP(ctx.Writer, ctx.Request)
}

// 首页处理方法
func (embyServerHandler *EmbyServerHandler) IndexHandler(ctx *gin.Context) {
	var (
		key             string = "IndexHandler"
		htmlFilePath    string = path.Join(config.StaticDir(), "index.html")
		modifiedBodyStr string
		addHEAD         string
	)
	if embyServerHandler.modifyProxyMap == nil {
		embyServerHandler.modifyProxyMap = make(map[string]*httputil.ReverseProxy)
	}

	if _, ok := embyServerHandler.modifyProxyMap[key]; !ok {
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
		embyServerHandler.modifyProxyMap[key] = proxy
	}
	embyServerHandler.modifyProxyMap[key].ServeHTTP(ctx.Writer, ctx.Request)
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
