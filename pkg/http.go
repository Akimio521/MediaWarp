package pkg

import (
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// 分离主机域名（IP）和端口号
//
// 示例：
// "example.com:8096" 					=> 	"example.com", "8096"
// "[240e:da8:a801:5a47::316]:8096" 	=> 	"240e:da8:a801:5a47::316" "8096"
// "192.168.1.1:8096" 					=> 	"192.168.1.1" "8096"
func SplitHostPort(hostPort string) (host string, port string) {
	host = hostPort

	colon := strings.LastIndexByte(host, ':')
	if colon != -1 && validOptionalPort(host[colon:]) {
		host, port = host[:colon], host[colon+1:]
	}

	// 对地址是IPv6地址进行处理
	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
		host = host[1 : len(host)-1]
	}

	return
}

func validOptionalPort(port string) bool {
	if port == "" {
		return true
	}
	if port[0] != ':' {
		return false
	}
	for _, b := range port[1:] {
		if b < '0' || b > '9' {
			return false
		}
	}
	return true
}

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
