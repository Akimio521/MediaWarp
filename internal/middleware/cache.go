package middleware

import (
	"MediaWarp/internal/cache"
	"MediaWarp/internal/logger"
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 缓存规则
//
// 根据请求路径自定义缓存时间
var cacheRules = []CacheRule{
	{Regexp: regexp.MustCompile(`(?i)^/.*embywebsocket\??`), Time: 0},                         // WebSocket不进入缓存处理
	{Regexp: regexp.MustCompile(`(?i)^/.*items/.*/playbackinfo\??`), Time: 1 * time.Hour},     // 播放信息
	{Regexp: regexp.MustCompile(`(?i)^/.*videos/.*/subtitles`), Time: 30 * 24 * time.Hour},    // 字幕
	{Regexp: regexp.MustCompile(`(?i)^/.*(\.html|\.css|\.js|\.woff)`), Time: 1 * time.Minute}, // 静态资源
}

// 缓存中间件
//
// 根据URI判断数据是否需要缓存
// 缓存Key为
func Cache() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			cacheDuration     time.Duration      = getCacheTime(ctx) // 缓存时间
			responseCacheData *ResponseCacheData = nil               // 缓存数据
		)

		// 判断是否需要缓存
		if cacheDuration == 0 {
			ctx.Next()
			return
		}

		// 缓存逻辑
		cacheKey := getCacheKey(ctx)                      // 计算缓存Key
		cacheData, ok := cache.Get("GIN-Cache", cacheKey) // 查询缓存记录
		logger.Debug("GIN-Cache：", "OK:", ok, ",KEY:", cacheKey)

		if ok { // 缓存存在，直接返回缓存数据
			responseCacheData = cacheData.(*ResponseCacheData)
			logger.Debug("GIN-Cache：缓存存在，直接返回缓存数据，HttpStatus:", responseCacheData.StatusCode)
			responseWithCache(ctx, responseCacheData)
		} else { // 缓存不存在，继续处理请求
			logger.Debug("GIN-Cache：缓存不存在，继续处理请求")
			customWirter := &ResponseWriterWarp{
				ResponseWriter: ctx.Writer,      // 使用原始的响应器初始化
				body:           &bytes.Buffer{}, // 初始化响应体缓存
			}
			ctx.Writer = customWirter // 使用自定义的响应器

			ctx.Next() // 处理请求
			code := ctx.Writer.Status()
			if code >= http.StatusOK && code < http.StatusMultipleChoices { // 响应是2xx的成功响应，更新缓存记录
				logger.Debug("GIN-Cache：创建缓存")
				responseCacheData = &ResponseCacheData{ // 创建缓存数据
					StatusCode: code, //ctx.Request.Response.StatusCode,
					Header:     ctx.Writer.Header().Clone(),
					Body:       customWirter.body.Bytes(),
				}
				go cache.Update("GIN-Cache", cacheKey, responseCacheData, cacheDuration)
			}
		}
	}
}

// 获取缓存时间
//
// 根据默认缓存时间为10分钟，无需缓存的请求返回0
func getCacheTime(ctx *gin.Context) time.Duration {
	if ctx.Writer.Header().Get("Expired") == "-1" {
		logger.Debug("Expired=-1 => 不缓存")
		return 0
	}

	if ctx.Request.Body != nil && ctx.Request.ContentLength != 0 { // 请求体不为空
		logger.Debug("请求体不为空 => 不缓存")
		return 0

	}

	for _, cacheRule := range cacheRules {
		if cacheRule.Regexp.MatchString(ctx.Request.RequestURI) {
			return cacheRule.Time
		}
	}
	// return 0
	if ctx.Request.Method == http.MethodGet {
		return 10 * time.Minute // 默认缓存时间为10分钟
	} else {
		logger.Debug("请求方法不为GET => 不缓存")
		return 0
	}
}

// 计算Key时忽略的查询参数
var CacheKeyIgnoreQuery = []string{
	// Fileball
	"starttimeticks",
	"x-playback-session-id",

	// Emby
	"playsessionid",
}

// 计算Key时忽略的请求头
var cacheKeyIgnoreHeaders = []string{
	"Range",
	"Host",
	"Referrer",
	"Connection",
	"Accept",
	"Accept-Encoding",
	"Accept-Language",
	"Cache-Control",
	"Upgrade-Insecure-Requests",
	"Referer",
	"Origin",

	// StreamMusic
	"X-Streammusic-Audioid",
	"X-Streammusic-Savepath",

	// IP
	"X-Forwarded-For",
	"X-Real-IP",
	"Forwarded",
	"Client-IP",
	"True-Client-IP",
	"CF-Connecting-IP",
	"X-Cluster-Client-IP",
	"Fastly-Client-IP",
	"X-Client-IP",
	"X-ProxyUser-IP",
	"Via",
	"Forwarded-For",
	"X-From-Cdn",
}

// 得到请求的CacheKey
//
// 计算方式: 取出 请求方法（method）， 请求路径（path），请求头（header），请求体（body）转换成字符串之后字典排序
func getCacheKey(ctx *gin.Context) string {
	var (
		method    string      = ctx.Request.Method      // 请求方法
		path      string      = ctx.Request.URL.Path    // 请求路径
		query     url.Values  = ctx.Request.URL.Query() // 查询参数
		queryStr  string      = ""                      // 查询参数字符串
		header    http.Header = ctx.Request.Header      // 请求头
		headerStr string      = ""                      // 请求头字符串
	)

	// 将查询参数转化为字符串
	for _, key := range CacheKeyIgnoreQuery {
		query.Del(key)
	}
	for key, values := range query { // 对查询参数的值进行排序
		sort.Strings(values)
		query[key] = values
	}
	queryStr = query.Encode()

	// 将请求头转化为字符串
	for _, key := range cacheKeyIgnoreHeaders {
		header.Del(key)
	}
	headKeys := make([]string, 0, len(header))
	for key := range header {
		headKeys = append(headKeys, key)
	}
	sort.Strings(headKeys) // 对请求头的键进行排序
	for _, key := range headKeys {
		sort.Strings(header[key]) // 对请求头的值进行排序
		headerStr += fmt.Sprintf("%s=%s;", key, strings.Join(header[key], "|"))
	}

	return method + path + queryStr + headerStr
}

// 缓存正则表达式和缓存时间对
type CacheRule struct {
	Regexp *regexp.Regexp
	Time   time.Duration
}

// 存放请求的响应信息
type ResponseCacheData struct {
	StatusCode int         // code 响应码
	Header     http.Header // header 响应头信息
	Body       []byte      // body 响应体

}

// 自定义的请求响应器
//
// 用于记录缓存数据
type ResponseWriterWarp struct {
	gin.ResponseWriter // gin 原始的响应器
	// statusCode         int // 响应码
	//header             http.Header   // 响应头
	body *bytes.Buffer // gin 回写响应时, 同步缓存
}

// WriteHeader 重写 WriteHeader 方法
//
// 用于记录响应码
// func (responseWriterWarp *ResponseWriterWarp) WriteHeader(statusCode int) {
// 	responseWriterWarp.statusCode = statusCode
// 	responseWriterWarp.ResponseWriter.WriteHeader(statusCode)
// }

// Write 重写 Write 方法
//
// 用于记录响应体
func (responseWriterWarp *ResponseWriterWarp) Write(data []byte) (int, error) {
	responseWriterWarp.body.Write(data)
	return responseWriterWarp.ResponseWriter.Write(data)
}

// 返回缓存数据
func responseWithCache(ctx *gin.Context, cacheData *ResponseCacheData) {
	ctx.Status(cacheData.StatusCode)            // 设置响应码
	for key, values := range cacheData.Header { // 设置响应头
		for _, value := range values {
			ctx.Writer.Header().Add(key, value)
		}
	}
	ctx.Writer.Write(cacheData.Body) // 设置响应体
	ctx.Abort()
}
