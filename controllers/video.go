package controllers

import (
	"MediaWarp/schemas/emby"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// /Videos/:itemId/:name，302重定向播放Strm
func VideosHandler(ctx *gin.Context) {
	itemId := ctx.Param("itemId")
	mediaSourceID := ctx.Query("mediasourceid")

	resp, err := http.Get(config.Origin + "/Items/" + itemId + "/PlaybackInfo?mediaSourceId=" + mediaSourceID + "&api_key=" + config.ApiKey)
	if err != nil {
		logger.ServerLogger.Warning("请求失败：", err)
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.ServerLogger.Warning("读取响应体失败：", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "读取响应体失败"})
		return
	}

	var playbackInfoResponse emby.PlaybackInfoResponse
	err = json.Unmarshal(body, &playbackInfoResponse)
	if err != nil {
		logger.ServerLogger.Warning("解析Json错误：", err)
		fmt.Printf("%s", body)
		return
	}
	for _, mediasource := range playbackInfoResponse.MediaSources {
		if *mediasource.ID == mediaSourceID {
			if *mediasource.Protocol == emby.HTTP {
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
