package api

import (
	"MediaWarp/schemas/schemas_emby"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type EmbyServer struct {
	ServerURL string
	ApiKey    string
}

// ItemsService
// /Items
func (embyServer *EmbyServer) ItemsServiceQueryItem(ids string, limit int, fields string) (ItemResponse *schemas_emby.ItemResponse, err error) {
	params := url.Values{}
	params.Add("Ids", ids)
	params.Add("Limit", strconv.Itoa(limit))
	params.Add("Fields", fields)
	params.Add("api_key", embyServer.ApiKey)
	api := embyServer.ServerURL + "/Items?" + params.Encode()
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
	resp, err := http.Get(embyServer.ServerURL + "/web/index.html")
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
