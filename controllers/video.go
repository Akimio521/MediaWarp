package controllers

import (
	"MediaWarp/api"
	"MediaWarp/core"
	"MediaWarp/schemas/schemas_emby"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// /Videos/:itemId/:name，302重定向播放Strm
func VideosHandler(ctx *gin.Context) {
	params := ctx.Request.URL.Query()
	mediaSourceID := params.Get("mediasourceid")
	var apiKey string
	apiKey = params.Get("api_key")
	if apiKey == "" {
		apiKey = params.Get("x-emby-token")
	}
	if apiKey == "" {
		apiKey = config.ApiKey
	}

	embyServer := api.EmbyServer{
		ServerURL: config.Origin,
		ApiKey:    apiKey,
	}
	ItemResponse, err := embyServer.ItemsServiceQueryItem(mediaSourceID, 1, "Path,MediaSources")
	if err != nil {
		logger.ServerLogger.Warning("请求ItemsServiceQueryItem失败：", err)
		return
	}
	item := ItemResponse.Items[0]
	for _, mediasource := range item.MediaSources {
		if *mediasource.ID == mediaSourceID {

			if strings.ToUpper(*mediasource.Container) == "STRM" { // 判断是否为Strm文件
				if *mediasource.Protocol == schemas_emby.HTTP {
					httpStrmRedirect(ctx, &mediasource)
					return
				}
				if config.AlistStrm.Enable { // 判断是否启用AlistStrm
					if alistStrmRedirect(ctx, &mediasource, &item) {
						return
					}
					logger.ServerLogger.Warning("未匹配AlistStrm路径：", *item.Path)
				}

			}
			logger.ServerLogger.Info("本地视频：", *mediasource.Path)
			DefaultHandler(ctx)
			return
		}
	}
}

// Strm文件内部为标准http协议时，302重定向播放
func httpStrmRedirect(ctx *gin.Context, mediasource *schemas_emby.MediaSourceInfo) {
	logger := core.GetLogger()
	logger.ServerLogger.Info("Http协议Strm重定向：", *mediasource.Path)
	ctx.Redirect(http.StatusFound, *mediasource.Path)
}

// 判断是否注册为AlistStrm，实现302重定向播放
func alistStrmRedirect(ctx *gin.Context, mediasource *schemas_emby.MediaSourceInfo, item *schemas_emby.BaseItemDto) (mathed bool) {
	mathed = false
	logger := core.GetLogger()
	var (
		alistPath   string
		alistServer api.AlistServer
	)
	for _, alistStrmConfig := range config.AlistStrm.List {
		if strings.HasPrefix(*item.Path, alistStrmConfig.Prefix) {
			alistPath = *mediasource.Path
			alistServer = alistStrmConfig.AlistServer
			logger.ServerLogger.Info("匹配AlistStrm路径：", *item.Path)
			break
		}
	}
	if alistPath != "" { // 匹配成功
		err := alistServer.AuthLogin()
		if err != nil {
			logger.ServerLogger.Warning("Alist登录失败：", err)
			return
		}
		fsGet, err := alistServer.FsGet(alistPath)
		if err != nil {
			logger.ServerLogger.Warning("请求GetFile失败：", err)
			return
		}
		logger.ServerLogger.Info("AlistStrm重定向：", fsGet.RawURL)
		ctx.Redirect(http.StatusFound, fsGet.RawURL)
		return true
	}
	return
}
