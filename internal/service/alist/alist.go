package alist

import (
	"MediaWarp/internal/config"
	"MediaWarp/internal/logging"
	"MediaWarp/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/allegro/bigcache"
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
	client   *http.Client
	cache    *bigcache.BigCache
}

// 获得AlistServer实例
func New(addr string, username string, password string, token *string) *AlistServer {
	s := AlistServer{
		endpoint: utils.GetEndpoint(addr),
		username: username,
		password: password,
		client:   utils.GetHTTPClient(),
	}
	if token != nil {
		s.token = alistToken{
			value:    *token,
			expireAt: time.Time{},
		}
	}
	if config.Cache.Enable && config.Cache.AlistTTL > 0 {
		cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(config.Cache.AlistTTL))
		if err == nil {
			s.cache = cache
		} else {
			logging.Warning("创建 Alist API 缓存失败: ", err)
		}
	}
	return &s
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

	loginData, err := alistServer.authLogin() // 重新生成一个token
	alistServer.token.mutex.RUnlock()
	if err != nil {
		return "", err
	}

	alistServer.token.mutex.Lock()
	defer alistServer.token.mutex.Unlock()
	alistServer.token.value = loginData.Token
	alistServer.token.expireAt = time.Now().Add(tokenDuration) // Token 有效期为30分钟

	return loginData.Token, nil
}

func (alistServer *AlistServer) doRequest(method string, path string, reqBody io.Reader, result any, needToken bool, needCache bool) error {
	var resp AlistResponse[any]

	cacheKey := method + "|" + path
	if needCache && alistServer.cache != nil {
		if data, err := alistServer.cache.Get(cacheKey); err == nil {
			return json.Unmarshal(data, result)
		}
	}
	var url = alistServer.GetEndpoint() + path
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	if needToken {
		token, err := alistServer.getToken()
		if err != nil {
			return err
		}
		req.Header.Add("Authorization", token)
	}

	res, err := alistServer.client.Do(req)
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		return fmt.Errorf("解析响应体失败: %w", err)
	}

	if resp.Code != http.StatusOK {
		return fmt.Errorf("请求失败，HTTP 状态码: %d, 相应状态码: %d, 相应信息: %s", res.StatusCode, resp.Code, resp.Message)
	}

	data, err := json.Marshal(resp.Data)
	if err != nil {
		return fmt.Errorf("序列化响应数据失败: %w", err)
	}
	err = json.Unmarshal(data, result)
	if err != nil {
		return fmt.Errorf("反序列化到结果类型失败: %w", err)
	}

	if needCache && alistServer.cache != nil {
		err = alistServer.cache.Set(cacheKey, data)
		if err != nil {
			return fmt.Errorf("缓存响应体失败: %w", err)
		}
	}

	return nil
}

// ==========Alist API(v3) 相关操作==========

// 登录Alist（获取一个新的Token）
func (alistServer *AlistServer) authLogin() (*AuthLoginData, error) {
	var (
		payload = strings.NewReader(fmt.Sprintf(`{"username": "%s","password": "%s"}`, alistServer.GetUsername(), alistServer.password))
		data    AuthLoginData
	)

	err := alistServer.doRequest(
		http.MethodPost,
		"/api/auth/login",
		payload,
		&data,
		false,
		false,
	)
	if err != nil {
		return nil, fmt.Errorf("登录失败: %w", err)
	}

	return &data, nil
}

// 获取某个文件/目录信息
func (alistServer *AlistServer) FsGet(path string) (*FsGetData, error) {
	var (
		data    FsGetData
		payload = strings.NewReader(fmt.Sprintf(`{"path": "%s","password": "","page": 1,"per_page": 0,"refresh": false}`, path))
	)
	err := alistServer.doRequest(
		http.MethodPost,
		"/api/fs/get",
		payload,
		&data,
		true,
		true,
	)
	if err != nil {
		return nil, fmt.Errorf("获取文件/目录信息失败: %w", err)
	}
	return &data, nil
}

func (alistServer *AlistServer) Me() (*UserInfoData, error) {
	var data UserInfoData
	err := alistServer.doRequest(
		http.MethodGet,
		"/api/me",
		nil,
		&data,
		true,
		true,
	)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %w", err)
	}
	return &data, nil
}

// GetFileURL 获取文件的可访问 URL
func (alistServer *AlistServer) GetFileURL(p string, isRawURL bool) (string, error) {
	fileData, err := alistServer.FsGet(p)
	if err != nil {
		return "", fmt.Errorf("获取文件信息失败：%w", err)
	}
	if isRawURL {
		return fileData.RawURL, nil
	}
	userInfo, err := alistServer.Me()
	if err != nil {
		return "", fmt.Errorf("获取用户当前信息失败：%w", err)
	}
	var sign string
	if fileData.Sign != "" {
		sign = "?sign=" + fileData.Sign
	}
	return alistServer.GetEndpoint() + path.Join("/d", userInfo.BasePath, p) + sign, nil
}
