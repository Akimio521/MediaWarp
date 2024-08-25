package api

import (
	"MediaWarp/schemas/emby"
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
func (embyServer *EmbyServer) ItemsServiceQueryItem(ids string, limit int, fields string) (ItemResponse *emby.ItemResponse, err error) {
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
