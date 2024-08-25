package controllers

import (
	"MediaWarp/api"
	"MediaWarp/schemas/schemas_emby"
	"net/http"

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
			if *mediasource.Protocol == schemas_emby.HTTP {
				logger.ServerLogger.Info("302重定向：", *mediasource.Path)
				ctx.Redirect(http.StatusFound, *mediasource.Path)
			} else {
				logger.ServerLogger.Info("本地视频：", *mediasource.Path)
				DefaultHandler(ctx)
				return
			}
		}
	}
}
