package jellyfin

import (
	"MediaWarp/constants"
	"MediaWarp/utils"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type Jellyfin struct {
	endpoint string
	apiKey   string // 认证方式：APIKey；获取方式：Jellyfin 控制台 -> 高级 -> API密钥
}

// 获取媒体服务器类型
func (jellyfin *Jellyfin) GetType() constants.MediaServerType {
	return constants.JELLYFIN
}

// 获取 Jellyfin 连接地址
//
// 包含协议、服务器域名（IP）、端口号
// 示例：return "http://jellyfin.example.com:8096"
func (jellyfin *Jellyfin) GetEndpoint() string {
	return jellyfin.endpoint
}

// 获取 Jellyfin 的API Key
func (jellyfin *Jellyfin) GetAPIKey() string {
	return jellyfin.apiKey
}

// ItemsService
// /Items
func (jellyfin *Jellyfin) ItemsServiceQueryItem(ids string, limit int, fields string) (*Response, error) {
	var (
		params       = url.Values{}
		itemResponse = &Response{}
	)
	params.Add("Ids", ids)
	params.Add("Limit", strconv.Itoa(limit))
	params.Add("Fields", fields)
	params.Add("api_key", jellyfin.GetAPIKey())

	resp, err := http.Get(jellyfin.GetEndpoint() + "/Items?" + params.Encode())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(body, itemResponse); err != nil {
		return nil, err
	}
	return itemResponse, nil
}

// 获取 Jellyfin 实例
func New(addr string, apiKey string) *Jellyfin {
	jellyfin := &Jellyfin{
		endpoint: utils.GetEndpoint(addr),
		apiKey:   apiKey,
	}
	return jellyfin
}
