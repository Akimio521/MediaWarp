package controllers

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// 拦截basehtmlplayer.js实现Web跨域播放Strm
func ModifyBaseHtmlPlayerHandler(ctx *gin.Context) {
	fmt.Println("匹配路由：/web/modules/htmlvideoplayer/basehtmlplayer.js\t请求方法：", ctx.Request.Method)
	version := ctx.Query("v")
	fmt.Println("请求basehtmlplayer.js版本：", version)
	resp, err := http.Get(cfg.Origin + ctx.Request.URL.Path + "?" + ctx.Request.URL.RawQuery)
	if err != nil {
		fmt.Println("请求失败，使用回源策略，错误信息：", err)
		DefaultHandler(ctx)
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应体失败，使用回源策略，错误信息：", err)
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
