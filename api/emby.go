package api

import (
	"MediaWarp/constants"
	"MediaWarp/schemas/schemas_emby"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
)

type EmbyServer struct {
	ADDR     string
	TOKEN    string
	endpoint string
	proxy    *httputil.ReverseProxy
}

// 初始化函数
func (embyServer *EmbyServer) Init() {
	embyServer.initEndpoint()
	embyServer.initProxy()
}

// 初始化endpoint
func (embyServer *EmbyServer) initEndpoint() {
	endpoint := embyServer.ADDR
	if !strings.HasPrefix(endpoint, "http") {
		endpoint = "http://" + endpoint
	}
	embyServer.endpoint = strings.TrimSuffix(endpoint, "/")
}

// 初始化proxy
func (embyServer *EmbyServer) initProxy() {
	target, _ := url.Parse(embyServer.GetEndpoint())
	embyServer.proxy = httputil.NewSingleHostReverseProxy(target)
}

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

// 获取EmbyServer的APIKey
func (embyServer *EmbyServer) GetToken() string {
	return embyServer.TOKEN
}

// 反代代理处理函数
func (embyServer *EmbyServer) ReverseProxy(rw http.ResponseWriter, req *http.Request) {
	embyServer.proxy.ServeHTTP(rw, req)
}

// ItemsService
// /Items
func (embyServer *EmbyServer) ItemsServiceQueryItem(ids string, limit int, fields string) (ItemResponse *schemas_emby.ItemResponse, err error) {
	params := url.Values{}
	params.Add("Ids", ids)
	params.Add("Limit", strconv.Itoa(limit))
	params.Add("Fields", fields)
	params.Add("api_key", embyServer.GetToken())
	api := embyServer.GetEndpoint() + "/Items?" + params.Encode()
	resp, err := http.Get(api)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &ItemResponse)
	if err != nil {
		return
	}
	return
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

// 实例化EmbyServer
func NewEmbyServer(addr string, token string) *EmbyServer {
	emby := &EmbyServer{
		ADDR:  addr,
		TOKEN: token,
	}
	emby.Init()
	return emby
}

// 实例化EmbyServer（返回一个符合MediaServer的接口）
func registerEmbyServer(addr string, token string) MediaServer {
	return NewEmbyServer(addr, token)
}
