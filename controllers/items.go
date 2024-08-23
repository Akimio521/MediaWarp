package controllers

import (
	"MediaWarp/constants"
	"MediaWarp/pkg"
	"MediaWarp/schemas/emby"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// /Users/:userId/Items/:itemId 处理
func ItemsHandler(ctx *gin.Context) {
	body, err := pkg.GetRespBody(ctx, config.Origin, config.ApiKey)
	if err != nil {
		logger.ServerLogger.Warning("获取Body出错：", err)
		return
	}
	var userItemsDetailResponse emby.UsersItemsDetailResponse
	err = json.Unmarshal(body, &userItemsDetailResponse)
	if err != nil {
		logger.ServerLogger.Warning("解析Json错误：", err)
		return
	}

	if *userItemsDetailResponse.MediaType == string(emby.Video) {
		logger.ServerLogger.Debug("查找视频类型Item，ItemID：", *userItemsDetailResponse.ID)
		for index, mediasource := range userItemsDetailResponse.MediaSources {
			if *mediasource.Protocol == emby.HTTP {
				userItemsDetailResponse.MediaSources[index].SupportsDirectPlay = &constants.BOOL_TRUE
				userItemsDetailResponse.MediaSources[index].SupportsDirectStream = &constants.BOOL_TRUE
				userItemsDetailResponse.MediaSources[index].SupportsTranscoding = &constants.BOOL_FALSE
			}
		}
	} else {
		logger.ServerLogger.Debug("查找非视频类型Item，ItemID：", *userItemsDetailResponse.ID)
	}

	ctx.JSON(http.StatusOK, userItemsDetailResponse)
}
