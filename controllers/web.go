package controllers

import (
	"MediaWarp/pkg"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// 拦截basehtmlplayer.js实现Web跨域播放Strm
func ModifyBaseHtmlPlayerHandler(ctx *gin.Context) {
	version := ctx.Query("v")
	logger.ServerLogger.Info("请求basehtmlplayer.js版本：", version)
	resp, err := http.Get(config.Origin + ctx.Request.URL.Path + "?" + ctx.Request.URL.RawQuery)
	if err != nil {
		logger.ServerLogger.Warning("请求失败，使用回源策略，错误信息：", err)
		DefaultHandler(ctx)
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.ServerLogger.Warning("读取响应体失败，使用回源策略，错误信息：", err)
		DefaultHandler(ctx)
		return
	}
	modifiedBody := strings.ReplaceAll(string(body), `mediaSource.IsRemote&&"DirectPlay"===playMethod?null:"anonymous"`, "null")

	for key, values := range resp.Header {
		for _, value := range values {
			if key != "Content-Length" {
				ctx.Writer.Header().Add(key, value)
			} else {
				ctx.Header("Content-Length", fmt.Sprintf("%d", len(modifiedBody)))
			}
		}
	}

	// 返回修改后的内容
	ctx.Data(resp.StatusCode, resp.Header.Get("Content-Type"), []byte(modifiedBody))

}

// 首页处理函数
func IndexHandler(ctx *gin.Context) {
	htmlFilePath := filepath.Join(config.StaticDir(), "index.html")
	headFilePath := filepath.Join(config.StaticDir(), "head")
	var (
		htmlContent    []byte
		retHtmlContent string
	)

	if pkg.PathExists(htmlFilePath) && pkg.IsFile(htmlFilePath) { // 检查htmlFilePath是否存在并且是文件
		logger.ServerLogger.Debug(htmlFilePath, "存在并且是文件")
		indexFile, err := os.OpenFile(htmlFilePath, os.O_RDONLY, 0666)
		if err != nil {
			logger.ServerLogger.Warning("打开", htmlFilePath, "失败，错误信息：", err)
		}
		if indexFile != nil {
			defer indexFile.Close()
			htmlContent, err = io.ReadAll(indexFile)
			if err != nil {
				logger.ServerLogger.Warning("读取", htmlFilePath, "内容失败，错误信息：", err)
			}
		}

	} else { // 请求上游EmbyServer获取HTML内容
		logger.ServerLogger.Debug("请求上游EmbyServer获取HTML内容")
		resp, err := http.Get(config.Origin + ctx.Request.URL.Path + "?" + ctx.Request.URL.RawQuery)
		if err != nil {
			logger.ServerLogger.Warning("请求失败，使用回源策略，错误信息：", err)
		}
		if resp != nil {
			defer resp.Body.Close()
			htmlContent, err = io.ReadAll(resp.Body)
			if err != nil {
				logger.ServerLogger.Warning("读取响应体失败，使用回源策略，错误信息：", err)
			}
		}
	}

	if htmlContent == nil {
		DefaultHandler(ctx)
		return
	}

	if pkg.PathExists(headFilePath) && pkg.IsFile(headFilePath) { // 检查 headFilePath 是否存在并且是文件
		logger.ServerLogger.Debug(headFilePath, "存在并且是文件")
		headFile, err := os.OpenFile(headFilePath, os.O_RDONLY, 0666)
		if err != nil {
			logger.ServerLogger.Warning("打开", headFilePath, "失败，错误信息：", err)
			retHtmlContent = string(htmlContent)
		} else {
			defer headFile.Close()
			headContent, err := io.ReadAll(headFile)
			if err != nil {
				logger.ServerLogger.Warning("读取", headFilePath, "内容失败，错误信息：", err)
				retHtmlContent = string(htmlContent)
			} else if len(headContent) > 0 {
				// 将 headContent 插入到 htmlContent 的 </head> 标签之前
				retHtmlContent = strings.Replace(string(htmlContent), "</head>", string(headContent)+"\n"+"</head>", 1)
			} else {
				logger.ServerLogger.Warning(headFilePath, "文件内容为空")
				retHtmlContent = string(htmlContent)
			}
		}
	} else {
		logger.ServerLogger.Debug(headFilePath, "不存在或不是文件")
		retHtmlContent = string(htmlContent)
	}

	ctx.Header("Content-Type", "text/html; charset=UTF-8")
	ctx.Header("Content-Length", fmt.Sprintf("%d", len(retHtmlContent)))
	ctx.Header("expires", "-1")
	ctx.String(http.StatusOK, retHtmlContent)
}
