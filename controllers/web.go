package controllers

import (
	"MediaWarp/pkg"
	"fmt"
	"io"
	"net/http"
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
	var (
		htmlFilePath   string = filepath.Join(config.StaticDir(), "index.html")
		headFilePath   string = filepath.Join(config.StaticDir(), "head")
		isFile         bool
		err            error
		htmlContent    []byte
		headContent    []byte
		retHtmlContent string
	)

	isFile, err = pkg.IsFile(htmlFilePath)
	if err != nil {
		logger.ServerLogger.Warning("判断路径是否为文件出错，错误信息：", err)
		isFile = false
	}

	if isFile {
		logger.ServerLogger.Debug(htmlFilePath, "存在并且是文件")
		htmlContent, err = pkg.GetFileContent(htmlFilePath)
		if err != nil {
			logger.ServerLogger.Warning("读取文件内容出错，使用回源策略，错误信息：", err)
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

	isFile, err = pkg.IsFile(headFilePath)
	if err != nil {
		logger.ServerLogger.Warning("判断路径是否为文件出错，错误信息：", err)
		isFile = false
	}

	if isFile {
		headContent, err = pkg.GetFileContent(htmlFilePath)
		if err != nil {
			logger.ServerLogger.Warning("读取文件内容出错，使用回源策略，错误信息：", err)
			retHtmlContent = string(htmlContent)
		} else {
			retHtmlContent = strings.Replace(string(htmlContent), "</head>", string(headContent)+"\n"+"</head>", 1)
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
