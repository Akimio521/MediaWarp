package api

import (
	"MediaWarp/schemas/schemas_alist"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type CacheItem struct {
	Data   interface{}
	Expiry time.Time
}

type AlistServer struct {
	URL         string
	Username    string
	Password    string
	token       string
	tokenExpiry time.Time
	cachePool   map[string]CacheItem
	cacheMutex  sync.Mutex
}

// 登录Alist（获取Token）
func (alistServer *AlistServer) authLogin() (err error) {
	funcInfo := "Alist登录"
	url := alistServer.URL + "/api/auth/login"
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(`{
    "username": "%s",
    "password": "%s"
}`, alistServer.Username, alistServer.Password))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		err = fmt.Errorf("创建%s请求失败: %w", funcInfo, err)
		return
	}
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("请求%s失败: %w", funcInfo, err)
		return
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("读取%s响应体失败: %w", funcInfo, err)
		return
	}

	var (
		authLoginData schemas_alist.AuthLoginData
		alistResponse schemas_alist.AlistResponse
	)

	err = json.Unmarshal(body, &alistResponse)
	if err != nil {
		err = fmt.Errorf("解析%s响应体失败: %w", funcInfo, err)
		return
	}
	if alistResponse.Code != 200 {
		err = errors.New(alistResponse.Message)
		return
	}
	dataBytes, err := json.Marshal(alistResponse.Data)
	if err != nil {
		err = fmt.Errorf("序列化%s响应体data失败: %w", funcInfo, err)
		return
	}

	// 将 JSON 字符串解析为 authLoginData 结构体
	err = json.Unmarshal(dataBytes, &authLoginData)
	if err != nil {
		err = fmt.Errorf("反序列化%s响应体data失败: %w", funcInfo, err)
		return
	}

	alistServer.token = authLoginData.Token
	alistServer.tokenExpiry = time.Now().Add(48 * time.Hour) // Token 有效期为2天
	return
}

// 检查alistServer.token是否过期，过期则自动重新登录获取新alistServer.token
func (alistServer *AlistServer) checkToken() (err error) {
	if time.Now().After(alistServer.tokenExpiry) {
		err = alistServer.authLogin()
		if err != nil {
			return err
		}
	}
	return nil
}

// 更新缓存
func (alistServer *AlistServer) updateCache(key string, data interface{}) {
	alistServer.cacheMutex.Lock()
	defer alistServer.cacheMutex.Unlock()
	alistServer.cachePool[key] = CacheItem{
		Data:   data,
		Expiry: time.Now().Add(1 * time.Hour),
	}
}

// 从缓存中获取数据
func (alistServer *AlistServer) getCache(key string) (data interface{}, found bool) {
	alistServer.cacheMutex.Lock()
	defer alistServer.cacheMutex.Unlock()
	item, found := alistServer.cachePool[key]
	return item.Data, found
}

// 获取某个文件/目录信息
func (alistServer *AlistServer) FsGet(path string) (fsGetData schemas_alist.FsGetData, err error) {
	err = alistServer.checkToken()
	if err != nil {
		return
	}
	if alistServer.cachePool == nil {
		alistServer.cachePool = make(map[string]CacheItem)
	}

	var alistResponse schemas_alist.AlistResponse
	cacheKey := fmt.Sprintf("API_FsGet_%s", path)
	cacheData, found := alistServer.getCache(cacheKey)
	if found {
		fsGetData = cacheData.(schemas_alist.FsGetData)
		return
	}

	funcInfo := "Alist获取某个文件/目录信息"
	url := alistServer.URL + "/api/fs/get"
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf(`{
    "path": "%s",
    "password": "",
    "page": 1,
    "per_page": 0,
    "refresh": false
}`, path))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		err = fmt.Errorf("创建%s请求失败: %w", funcInfo, err)
		return
	}
	req.Header.Add("Authorization", alistServer.token)
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://apifox.com)")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("请求%s信息失败: %w", funcInfo, err)
		return
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("读取%s响应体失败: %w", funcInfo, err)
		return
	}

	err = json.Unmarshal(body, &alistResponse)
	if err != nil {
		err = fmt.Errorf("解析%s响应体失败: %w", funcInfo, err)
		return
	}
	if alistResponse.Code != 200 {
		err = errors.New(alistResponse.Message)
		return
	}
	dataBytes, err := json.Marshal(alistResponse.Data)
	if err != nil {
		err = fmt.Errorf("序列化%s响应体data失败: %w", funcInfo, err)
		return
	}

	err = json.Unmarshal(dataBytes, &fsGetData)
	if err != nil {
		err = fmt.Errorf("反序列化%s响应体data失败: %w", funcInfo, err)
		return
	}
	alistServer.updateCache(cacheKey, fsGetData)
	return
}
