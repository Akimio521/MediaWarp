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
	URL        string
	Username   string
	Password   string
	cachePool  map[string]CacheItem
	cacheMutex sync.Mutex
}

// 更新缓存中的数据
func (alistServer *AlistServer) updateCache(key string, data interface{}, cacheDuration time.Duration) {
	alistServer.cacheMutex.Lock()
	defer alistServer.cacheMutex.Unlock()
	alistServer.cachePool[key] = CacheItem{
		Data:   data,
		Expiry: time.Now().Add(cacheDuration),
	}
}

// 从缓存中获取数据
func (alistServer *AlistServer) getCache(key string) (data interface{}, found bool) {
	alistServer.cacheMutex.Lock()
	defer alistServer.cacheMutex.Unlock()
	item, found := alistServer.cachePool[key]
	if found && time.Now().After(item.Expiry) { // 找到缓存并且缓存已过期
		delete(alistServer.cachePool, key) // 删除缓存
		found = false                      // 修改为未找到缓存
	}
	return item.Data, found
}

// 得到一个可用的Token
//
// 先从缓存池中读取，若过期或者未找到则重新生成
func (alistServer *AlistServer) getToken() (string, error) {
	var (
		cacheKey      = "API_TOKEN"
		token         string
		cacheDuration = 48*time.Hour - 5*time.Minute // Alist v3 Token默认有效期为2天，5分钟作为误差
	)
	cacheToken, found := alistServer.getCache(cacheKey)
	if found { // 找到token
		token = cacheToken.(string)
		return token, nil
	}

	// 未找到已缓存的token
	token, err := alistServer.authLogin() // 重新生成一个token
	if err != nil {
		return "", err
	}
	alistServer.updateCache(cacheKey, token, cacheDuration) // 将新生成的token添加到缓存池中

	return token, nil
}

// ==========Alist API(v3) 相关操作==========

// 登录Alist（获取一个新的Token）
func (alistServer *AlistServer) authLogin() (string, error) {
	var (
		funcInfo          = "Alist登录"
		url               = alistServer.URL + "/api/auth/login"
		method            = "POST"
		payload           = strings.NewReader(fmt.Sprintf(`{"username": "%s","password": "%s"}`, alistServer.Username, alistServer.Password))
		authLoginResponse schemas_alist.AlistResponse[schemas_alist.AuthLoginData]
	)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		err = fmt.Errorf("创建%s请求失败: %w", funcInfo, err)
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("请求%s失败: %w", funcInfo, err)
		return "", err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("读取%s响应体失败: %w", funcInfo, err)
		return "", err
	}

	err = json.Unmarshal(body, &authLoginResponse)
	if err != nil {
		err = fmt.Errorf("解析%s响应体失败: %w", funcInfo, err)
		return "", err
	}
	if authLoginResponse.Code != 200 {
		err = errors.New(authLoginResponse.Message)
		return "", err
	}

	return authLoginResponse.Data.Token, nil
}

// 获取某个文件/目录信息
func (alistServer *AlistServer) FsGet(path string) (schemas_alist.FsGetData, error) {
	var (
		fsGetDataResponse schemas_alist.AlistResponse[schemas_alist.FsGetData]
		token             string
		cacheKey          = fmt.Sprintf("API_FsGet_%s", path)
		cacheDuration     = 30 * time.Minute // 缓存时间为30分钟
		funcInfo          = "Alist获取某个文件/目录信息"
		url               = alistServer.URL + "/api/fs/get"
		method            = "POST"
		payload           = strings.NewReader(fmt.Sprintf(`{"path": "%s","password": "","page": 1,"per_page": 0,"refresh": false}`, path))
	)

	if alistServer.cachePool == nil {
		alistServer.cachePool = make(map[string]CacheItem)
	}

	cacheData, found := alistServer.getCache(cacheKey)
	if found { // 在缓存中查询到数据
		fsGetData := cacheData.(schemas_alist.FsGetData)
		return fsGetData, nil
	}

	// 未从缓存池中读取到数据
	token, err := alistServer.getToken()
	if err != nil {
		return fsGetDataResponse.Data, nil
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		err = fmt.Errorf("创建%s请求失败: %w", funcInfo, err)
		return fsGetDataResponse.Data, err
	}
	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("请求%s信息失败: %w", funcInfo, err)
		return fsGetDataResponse.Data, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("读取%s响应体失败: %w", funcInfo, err)
		return fsGetDataResponse.Data, err
	}

	err = json.Unmarshal(body, &fsGetDataResponse)
	if err != nil {
		err = fmt.Errorf("解析%s响应体失败: %w", funcInfo, err)
		return fsGetDataResponse.Data, err
	}
	if fsGetDataResponse.Code != 200 {
		err = errors.New(fsGetDataResponse.Message)
		return fsGetDataResponse.Data, err
	}

	alistServer.updateCache(cacheKey, fsGetDataResponse.Data, cacheDuration)
	return fsGetDataResponse.Data, nil
}
