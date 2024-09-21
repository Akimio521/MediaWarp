package handlers_emby

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
	// EmbyServer <= 4.8 ====> mediaSourceID = 343121
	// EmbyServer >= 4.9 ====> mediaSourceID = mediasource_31
	mediaSourceID := ctx.Query("mediasourceid")

	logger.ServerLogger.Debug("请求ItemsServiceQueryItem：", mediaSourceID)
	itemResponse, err := embyServer.ItemsServiceQueryItem(strings.Replace(mediaSourceID, "mediasource_", "", 1), 1, "Path,MediaSources") // 查询item需要去除前缀仅保留数字部分

	if err != nil {
		logger.ServerLogger.Warning("请求ItemsServiceQueryItem失败：", err)
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

			logger.ServerLogger.Info("本地视频：", *mediasource.Path)
			embyServer.ReverseProxy().ServeHTTP(ctx.Writer, ctx.Request)
			return
		}
	}
}

// 获取302重定向URL
func getRedirctURL(mediasource *schemas_emby.MediaSourceInfo, item *schemas_emby.BaseItemDto) (redirectURL string) {
	if *mediasource.Protocol == schemas_emby.HTTP && config.HttpStrm.Enable { // 判断是否为http协议Strm
		httpStrmRedirectURL := getHttpStrmRedirectURL(mediasource, item)
		if httpStrmRedirectURL != "" {
			return httpStrmRedirectURL
		}
	}
	if strings.ToUpper(*mediasource.Container) == "STRM" { // 判断是否为Strm文件
		if config.AlistStrm.Enable { // 判断是否启用AlistStrm
			alistStrmRedirectURL := getAlistStrmRedirect(mediasource, item)
			if alistStrmRedirectURL != "" {
				return alistStrmRedirectURL
			}
		}
	}
	return ""
}

// Strm文件内部为标准http协议时，获取302重定向URL
func getHttpStrmRedirectURL(mediasource *schemas_emby.MediaSourceInfo, item *schemas_emby.BaseItemDto) (url string) {
	logger := core.GetLogger()
	for _, prefix := range config.HttpStrm.PrefixList {
		if strings.HasPrefix(*item.Path, prefix) {
			logger.ServerLogger.Debug(*item.Path, "匹配HttpStrm路径：", prefix, "成功")
			logger.ServerLogger.Info("Http协议Strm重定向：", *mediasource.Path)
			return *mediasource.Path
		}
	}
	logger.ServerLogger.Info("未匹配HttpStrm路径：", *item.Path)
	return ""
}

// 判断是否注册为AlistStrm，获取302重定向URL
func getAlistStrmRedirect(mediasource *schemas_emby.MediaSourceInfo, item *schemas_emby.BaseItemDto) (url string) {
	logger := core.GetLogger()
	var (
		alistPath   string
		alistServer *api.AlistServer
	)
	for index := range config.AlistStrm.List {
		alistStrmConfig := &config.AlistStrm.List[index]
		for _, perfix := range alistStrmConfig.PrefixList {
			if strings.HasPrefix(*item.Path, perfix) {
				alistPath = *mediasource.Path
				alistServer = &alistStrmConfig.AlistServer // 获取AlistServer，无需重新生成实例
				logger.ServerLogger.Debug(*item.Path, "匹配AlistStrm路径：", perfix, "成功")
				break
			}
		}

	}
	if alistPath != "" { // 匹配成功
		fsGet, err := alistServer.FsGet(alistPath)
		if err != nil {
			logger.ServerLogger.Warning("请求GetFile失败：", err)
			return ""
		}
		logger.ServerLogger.Info("AlistStrm重定向：", fsGet.RawURL)
		return fsGet.RawURL
	}
	logger.ServerLogger.Info("未匹配AlistStrm路径：", *item.Path)
	return ""
}
