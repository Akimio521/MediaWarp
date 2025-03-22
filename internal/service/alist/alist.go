package alist

import (
	"MediaWarp/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

type alistToken struct {
	value    string       // 令牌 Token
	expireAt time.Time    // 令牌过期时间
	mutex    sync.RWMutex // 令牌锁
}
type AlistServer struct {
	endpoint string // 服务器入口 URL
	username string // 用户名
	password string // 密码
	token    alistToken
}

// 得到服务器入口
//
// 避免直接访问 endpoint 字段
func (alistServer *AlistServer) GetEndpoint() string {
	return alistServer.endpoint
}

// 得到用户名
//
// 避免直接访问 username 字段
func (alistServer *AlistServer) GetUsername() string {
	return alistServer.username
}

// 得到一个可用的 Token
//
// 先从缓存池中读取，若过期或者未找到则重新生成
func (alistServer *AlistServer) getToken() (string, error) {
	var tokenDuration = 2*24*time.Hour - 5*time.Minute // Token 有效期为 2 天，提前 5 分钟刷新

	alistServer.token.mutex.RLock()
	if alistServer.token.value != "" && (alistServer.token.expireAt.IsZero() || time.Now().Before(alistServer.token.expireAt)) {
		// 零值表示永不过期
		defer alistServer.token.mutex.RUnlock()
		return alistServer.token.value, nil
	}

	token, err := alistServer.authLogin() // 重新生成一个token
	alistServer.token.mutex.RUnlock()
	if err != nil {
		return "", err
	}

	alistServer.token.mutex.Lock()
	defer alistServer.token.mutex.Unlock()
	alistServer.token.value = token
	alistServer.token.expireAt = time.Now().Add(tokenDuration) // Token 有效期为30分钟

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
		err = fmt.Errorf("创建 %s 请求失败: %w", funcInfo, err)
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("请求 %s 失败: %w", funcInfo, err)
		return "", err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("读取 %s 响应体失败: %w", funcInfo, err)
		return "", err
	}

	err = json.Unmarshal(body, &authLoginResponse)
	if err != nil {
		err = fmt.Errorf("解析 %s 响应体失败: %w", funcInfo, err)
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
		funcInfo          = "Alist获取某个文件/目录信息"
		url               = alistServer.GetEndpoint() + "/api/fs/get"
		method            = "POST"
		payload           = strings.NewReader(fmt.Sprintf(`{"path": "%s","password": "","page": 1,"per_page": 0,"refresh": false}`, path))
	)

	// 未从缓存池中读取到数据
	token, err := alistServer.getToken()
	if err != nil {
		return fsGetDataResponse.Data, nil
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		err = fmt.Errorf("创建 %s 请求失败: %w", funcInfo, err)
		return fsGetDataResponse.Data, err
	}
	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("请求 %s 信息失败: %w", funcInfo, err)
		return fsGetDataResponse.Data, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("读取 %s 响应体失败: %w", funcInfo, err)
		return fsGetDataResponse.Data, err
	}

	err = json.Unmarshal(body, &fsGetDataResponse)
	if err != nil {
		err = fmt.Errorf("解析 %s 响应体失败: %w", funcInfo, err)
		return fsGetDataResponse.Data, err
	}
	if fsGetDataResponse.Code != 200 {
		err = errors.New(fsGetDataResponse.Message)
		return fsGetDataResponse.Data, err
	}

	return fsGetDataResponse.Data, nil
}

// 获得AlistServer实例
func New(addr string, username string, password string, token *string) *AlistServer {
	s := AlistServer{
		endpoint: utils.GetEndpoint(addr),
		username: username,
		password: password,
	}
	if token != nil {
		s.token = alistToken{
			value:    *token,
			expireAt: time.Time{},
		}
	}
	return &s
}
