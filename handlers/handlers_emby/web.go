package handlers_emby

import (
	"MediaWarp/handlers"
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
	resp, err := http.Get(config.Server.GetHTTPEndpoint() + ctx.Request.URL.Path + "?" + ctx.Request.URL.RawQuery)
	if err != nil {
		logger.ServerLogger.Warning("请求失败，使用回源策略，错误信息：", err)
		handlers.DefaultHandler(ctx)
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.ServerLogger.Warning("读取响应体失败，使用回源策略，错误信息：", err)
		handlers.DefaultHandler(ctx)
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

	if config.Web.Index { // 自定义首页
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
		}
	}
	if len(htmlContent) == 0 { // 未启用自定义首页或自定义首页文件内容读取失败
		logger.ServerLogger.Debug("请求上游EmbyServer获取HTML内容")
		htmlContent, err = embyServer.GetIndexHtml()
		if err != nil {
			logger.ServerLogger.Warning("请求上游EmbyServer获取HTML内容失败，使用回源策略，错误信息：", err)
		}
	}

	if len(htmlContent) == 0 {
		handlers.DefaultHandler(ctx)
		return
	}

	if config.Web.Head { // 自定义HEAD
		isFile, err = pkg.IsFile(headFilePath)
		if err != nil {
			logger.ServerLogger.Warning("判断路径是否为文件出错，错误信息：", err)
			isFile = false
		}

		if isFile {
			headContent, err = pkg.GetFileContent(headFilePath)
			if err != nil {
				logger.ServerLogger.Warning("读取文件内容出错，不添加额外HEAD，错误信息：", err)
			}
		} else {
			logger.ServerLogger.Debug(headFilePath, "不存在或不是文件")
		}
	}

	if len(headContent) == 0 {
		retHtmlContent = string(htmlContent)
	} else {
		retHtmlContent = strings.Replace(string(htmlContent), "</head>", string(headContent)+"\n"+"</head>", 1)
	}

	if config.Web.ExternalPlayerUrl {
		retHtmlContent = strings.Replace(retHtmlContent, "</head>", `<script src="/MediaWarp/Resources/js/ExternalPlayerUrl.js"></script>`+"\n"+"</head>", 1)
	}
	if config.Web.ActorPlus {
		retHtmlContent = strings.Replace(retHtmlContent, "</head>", `<script src="/MediaWarp/Resources/js/ActorPlus.js"></script>`+"\n"+"</head>", 1)
	}
	if config.Web.FanartShow {
		retHtmlContent = strings.Replace(retHtmlContent, "</head>", `<script src="/MediaWarp/Resources/js/FanartShow.js"></script>`+"\n"+"</head>", 1)
	}
	if config.Web.Danmaku {
		retHtmlContent = strings.Replace(retHtmlContent, "</body>", `<script src="https://cdn.jsdelivr.net/gh/RyoLee/emby-danmaku@gh-pages/ede.user.js" defer></script>`+"\n"+"</body>", 1)
	}

	if config.Web.BeautifyCSS {
		retHtmlContent = strings.Replace(retHtmlContent, "</head>", `<link rel="stylesheet" href="/MediaWarp/Resources/css/Beautify.css" type="text/css" media="all" />`+"\n"+"</head>", 1)
	}

	ctx.Header("Content-Type", "text/html; charset=UTF-8")
	ctx.Header("Content-Length", fmt.Sprintf("%d", len(retHtmlContent)))
	ctx.Header("expires", "-1")
	ctx.String(http.StatusOK, retHtmlContent)
}
