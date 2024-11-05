package utils

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 代理访问上游服务器返回的响应体
func GetRespBody(ctx *gin.Context, target string, api_key string) (body []byte, err error) {
	params := ctx.Request.URL.Query()
	if params.Get("api_key") == "" || params.Get("api_key") == "x-emby-token" {
		params.Set("api_key", api_key)
	}
	api := target + ctx.Request.URL.Path + "?" + params.Encode()
	client := http.Client{}
	req, err := http.NewRequest(ctx.Request.Method, api, ctx.Request.Body)
	if err != nil {
		return
	}

	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	// 清除所有响应头
	// for key := range ctx.Writer.Header() {
	// 	ctx.Writer.Header().Del(key)
	// }

	for key, values := range resp.Header {
		for _, value := range values {
			if key != "Content-Length" {
				ctx.Writer.Header().Set(key, value)
			}
		}
	}
	return
}
