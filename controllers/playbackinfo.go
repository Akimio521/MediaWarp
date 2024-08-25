package controllers

import (
	"MediaWarp/constants"
	"MediaWarp/pkg"
	"MediaWarp/schemas/schemas_emby"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// /Items/:itemId/PlaybackInfo 处理播放信息请求
func PlaybackInfoHandler(ctx *gin.Context) {
	body, err := pkg.GetRespBody(ctx, config.Origin, config.ApiKey)
	if err != nil {
		logger.ServerLogger.Warning("获取Body出错：", err)
		return
	}

	var playbackInfoResponse schemas_emby.PlaybackInfoResponse
	err = json.Unmarshal(body, &playbackInfoResponse)
	if err != nil {
		logger.ServerLogger.Warning("解析Json错误：", err)
		return
	}

	for index, mediasource := range playbackInfoResponse.MediaSources {
		if *mediasource.Protocol == schemas_emby.HTTP {
			playbackInfoResponse.MediaSources[index].SupportsDirectPlay = &constants.BOOL_TRUE
			playbackInfoResponse.MediaSources[index].SupportsDirectStream = &constants.BOOL_TRUE
			playbackInfoResponse.MediaSources[index].SupportsTranscoding = &constants.BOOL_FALSE
		}
	}

	ctx.JSON(http.StatusOK, playbackInfoResponse)
}
