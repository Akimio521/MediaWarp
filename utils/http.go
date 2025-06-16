package utils

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"time"
)

// 判断 url 是否经过编码
func IsURLEncoded(u string) bool {
	unescaped, err := url.QueryUnescape(u)
	if err != nil {
		return false // 非法编码
	}
	return url.QueryEscape(unescaped) == u
}

// 创建优化配置的 HTTP 客户端
func createOptimizedClient() *http.Client {
	transport := &http.Transport{
		// 连接池配置
		MaxIdleConns:        runtime.NumCPU() * 80, // 全局最大空闲连接
		MaxIdleConnsPerHost: runtime.NumCPU() * 5,  // 每个主机最大空闲连接
		IdleConnTimeout:     90 * time.Second,      // 空闲连接超时时间

		// 连接复用优化
		DisableKeepAlives: false, // 启用 Keep-Alive
		ForceAttemptHTTP2: true,  // 启用 HTTP/2

		// 连接建立优化
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,  // 连接超时
			KeepAlive: 30 * time.Second, // Keep-Alive 周期
		}).DialContext,

		// TLS 配置
		TLSHandshakeTimeout: 5 * time.Second, // TLS 握手超时
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false, // 生产环境应为 false
			MinVersion:         tls.VersionTLS12,
		},
	}

	return &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second, // 整个请求超时
	}
}

var httpClient *http.Client = createOptimizedClient() // 全局客户端单例

// 获取全局 HTTP 客户端（线程安全）
func GetHTTPClient() *http.Client {
	return httpClient
}
