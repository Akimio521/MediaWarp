package emby

import (
	"MediaWarp/constants"
	"MediaWarp/utils"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

type EmbyServer struct {
	endpoint string
	apiKey   string // 认证方式：APIKey；获取方式：Emby控制台 -> 高级 -> API密钥
	proxy    *httputil.ReverseProxy
}

// 初始化函数
func (embyServer *EmbyServer) Init() {
	embyServer.initProxy()
}

// 初始化proxy
func (embyServer *EmbyServer) initProxy() {
	target, _ := url.Parse(embyServer.endpoint)
	embyServer.proxy = httputil.NewSingleHostReverseProxy(target)
}

// 获取媒体服务器类型
func (embyServer *EmbyServer) GetType() constants.MediaServerType {
	return constants.EMBY
}

// 获取EmbyServer连接地址
//
// 包含协议、服务器域名（IP）、端口号
// 示例：return "http://emby.example.com:8096"
func (embyServer *EmbyServer) GetEndpoint() string {
	return embyServer.endpoint
}

// 获取EmbyServer的API Key
func (embyServer *EmbyServer) GetAPIKey() string {
	return embyServer.apiKey
}

// 获取反代服务器
//
// 根据EmbyServer的proxy创建一个新的反代服务器用于处理请求
// 对此httputil.ReverseProxy进行修改不影响EmbyServer的ReverseProxy()方法的行为
func (embyServer *EmbyServer) GetReverseProxy() *httputil.ReverseProxy {
	return &httputil.ReverseProxy{Director: embyServer.proxy.Director}
}

// 反代上游响应
func (embyServer *EmbyServer) ReverseProxy(rw http.ResponseWriter, req *http.Request) {
	if embyServer.proxy != nil {
		embyServer.proxy.ServeHTTP(rw, req)
	} else {
		panic("反代服务器未初始化")
	}
}

// ItemsService
// /Items
func (embyServer *EmbyServer) ItemsServiceQueryItem(ids string, limit int, fields string) (*EmbyResponse, error) {
	var (
		params       = url.Values{}
		itemResponse = &EmbyResponse{}
	)
	params.Add("Ids", ids)
	params.Add("Limit", strconv.Itoa(limit))
	params.Add("Fields", fields)
	params.Add("api_key", embyServer.GetAPIKey())
	api := embyServer.GetEndpoint() + "/Items?" + params.Encode()
	resp, err := http.Get(api)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, itemResponse)
	if err != nil {
		return nil, err
	}
	return itemResponse, nil
}

// 获取index.html内容 API：/web/index.html
func (embyServer *EmbyServer) GetIndexHtml() ([]byte, error) {
	resp, err := http.Get(embyServer.GetEndpoint() + "/web/index.html")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	htmlContent, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return htmlContent, nil
}

// 获取EmbyServer实例
func New(addr string, apiKey string) *EmbyServer {
	emby := &EmbyServer{
		endpoint: utils.GetEndpoint(addr),
		apiKey:   apiKey,
	}
	emby.Init()
	return emby
}
