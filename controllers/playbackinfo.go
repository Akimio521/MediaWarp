package controllers

import (
	"MediaWarp/schemas/emby"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// /Items/:itemId/PlaybackInfo 处理播放信息请求
func PlaybackInfoHandler(ctx *gin.Context) {
	var api string
	if strings.Contains(ctx.Request.URL.RawQuery, "api_key=") {
		api = config.Origin + ctx.Request.URL.Path + "?" + ctx.Request.URL.RawQuery
	} else {
		logger.ServerLogger.Warning("未找到api_key参数")
		if ctx.Request.URL.RawQuery == "" {
			api = config.Origin + ctx.Request.URL.Path + "?api_key=" + config.ApiKey
		} else {
			api = config.Origin + ctx.Request.URL.Path + "?" + ctx.Request.URL.RawQuery + "&api_key=" + config.ApiKey
		}
	}

	client := http.Client{}
	req, err := http.NewRequest(ctx.Request.Method, api, ctx.Request.Body)
	if err != nil {
		logger.ServerLogger.Warning("创建Request出错：", err)
		return
	}

	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		logger.ServerLogger.Warning("请求失败：", err)
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.ServerLogger.Warning("读取响应体失败：", err)
		return
	}

	var playbackInfoResponse emby.PlaybackInfoResponse
	err = json.Unmarshal(body, &playbackInfoResponse)
	if err != nil {
		logger.ServerLogger.Warning("解析Json错误：", err)
		return
	}

	for index, mediasource := range playbackInfoResponse.MediaSources {
		if *mediasource.Protocol == emby.HTTP {
			trueValue := true
			falseValue := false
			playbackInfoResponse.MediaSources[index].SupportsDirectPlay = &trueValue
			playbackInfoResponse.MediaSources[index].SupportsDirectStream = &trueValue
			playbackInfoResponse.MediaSources[index].SupportsTranscoding = &falseValue
		}
	}

	jsonData, err := json.Marshal(playbackInfoResponse)
	if err != nil {
		logger.ServerLogger.Warning("序列化播放信息失败：", err)
		return
	}

	// 将目标服务的响应返回给客户端
	for key, values := range resp.Header {
		for _, value := range values {
			if key != "Content-Length" {
				ctx.Writer.Header().Add(key, value)
			}
		}
	}
	ctx.Data(resp.StatusCode, resp.Header.Get("Content-Type"), jsonData)
}
