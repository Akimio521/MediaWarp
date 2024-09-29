package alist

import (
	"MediaWarp/internal/utils/cache"
	"MediaWarp/pkg"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type CacheItem struct {
	Data   interface{}
	Expiry time.Time
}

type AlistServer struct {
	endpoint   string
	username   string
	password   string
	cache      cache.Cache
	sapaceName string
}

// 得到缓存SpaceName
func (alistServer *AlistServer) Init() {
	alistServer.sapaceName = alistServer.GetEndpoint() + alistServer.GetUsername() + alistServer.password
}

// 得到服务器入口
//
// 避免直接访问endpoint字段
func (alistServer *AlistServer) GetEndpoint() string {
	return alistServer.endpoint
}

// 得到用户名
//
// 避免直接访问username字段
func (alistServer *AlistServer) GetUsername() string {
	return alistServer.username
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
	cacheToken, found := alistServer.cache.GetCache(alistServer.sapaceName, cacheKey)
	if found { // 找到token
		token = cacheToken.(string)
		return token, nil
	}

	// 未找到已缓存的token
	token, err := alistServer.authLogin() // 重新生成一个token
	if err != nil {
		return "", err
	}
	go alistServer.cache.UpdateCache(alistServer.sapaceName, cacheKey, token, cacheDuration) // 将新生成的token添加到缓存池中

	return token, nil
}

// ==========Alist API(v3) 相关操作==========

// 登录Alist（获取一个新的Token）
func (alistServer *AlistServer) authLogin() (string, error) {
	var (
		funcInfo          = "Alist登录"
		url               = alistServer.GetEndpoint() + "/api/auth/login"
		method            = http.MethodPost
		payload           = strings.NewReader(fmt.Sprintf(`{"username": "%s","password": "%s"}`, alistServer.GetUsername(), alistServer.password))
		authLoginResponse AlistResponse[AuthLoginData]
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
func (alistServer *AlistServer) FsGet(path string) (FsGetData, error) {
	var (
		fsGetDataResponse AlistResponse[FsGetData]
		token             string
		cacheKey          = fmt.Sprintf("API_FsGet_%s", path)
		cacheDuration     = 30 * time.Minute // 缓存时间为30分钟
		funcInfo          = "Alist获取某个文件/目录信息"
		url               = alistServer.GetEndpoint() + "/api/fs/get"
		method            = "POST"
		payload           = strings.NewReader(fmt.Sprintf(`{"path": "%s","password": "","page": 1,"per_page": 0,"refresh": false}`, path))
	)

	cacheData, found := alistServer.cache.GetCache(alistServer.sapaceName, cacheKey)
	if found { // 在缓存中查询到数据
		fsGetData := cacheData.(FsGetData)
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

	go alistServer.cache.UpdateCache(alistServer.sapaceName, cacheKey, fsGetDataResponse.Data, cacheDuration)
	return fsGetDataResponse.Data, nil
}

// 获得AlistServer实例
func New(addr string, username string, password string, cache cache.Cache) *AlistServer {
	return &AlistServer{
		endpoint: pkg.GetEndpoint(addr),
		username: username,
		password: password,
		cache:    cache,
	}
}
