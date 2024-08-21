package controllers

import (
	"MediaWarp/schemas/emby"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// /Users/:userId/Items/:itemId 处理
func ItemsHandler(ctx *gin.Context) {
	fmt.Println("匹配路由：ItemsHandler", "\t请求方法:", ctx.Request.Method, "\n路径参数:", ctx.Request.URL.Path, "\n查询参数:", ctx.Request.URL.RawQuery)

	var api string
	if strings.Contains(ctx.Request.URL.RawQuery, "api_key=") {
		api = cfg.Origin + ctx.Request.URL.Path + "?" + ctx.Request.URL.RawQuery
	} else {
		fmt.Println("未找到api_key参数")
		if ctx.Request.URL.RawQuery == "" {
			api = cfg.Origin + ctx.Request.URL.Path + "?api_key=" + cfg.ApiKey
		} else {
			api = cfg.Origin + ctx.Request.URL.Path + "?" + ctx.Request.URL.RawQuery + "&api_key=" + cfg.ApiKey
		}
	}

	client := http.Client{}
	req, err := http.NewRequest(ctx.Request.Method, api, ctx.Request.Body)
	if err != nil {
		fmt.Println("创建Request出错：", err)
		return
	}

	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败：", err)
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应体失败：", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "读取响应体失败"})
		return
	}

	var userItemsDetailResponse emby.UsersItemsDetailResponse
	err = json.Unmarshal(body, &userItemsDetailResponse)
	if err != nil {
		fmt.Println("解析Json错误：", err)
		fmt.Printf("%s", body)
		return
	}

	for index, mediasource := range userItemsDetailResponse.MediaSources {
		if *mediasource.Protocol == emby.HTTP {
			trueValue := true
			falseValue := false
			userItemsDetailResponse.MediaSources[index].SupportsDirectPlay = &trueValue
			userItemsDetailResponse.MediaSources[index].SupportsDirectStream = &trueValue
			userItemsDetailResponse.MediaSources[index].SupportsTranscoding = &falseValue
		}
	}

	jsonData, err := json.Marshal(userItemsDetailResponse)
	if err != nil {
		fmt.Println("序列化播放信息失败：", err)
		ctx.JSON(500, gin.H{"error": "序列化播放信息失败"})
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
