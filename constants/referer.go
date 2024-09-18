package constants

type REFERER_VALUE string // 常用 referer 策略值

const (
	NO_REFERRER                     REFERER_VALUE = "no-referrer"                     // 不发送 Referer 头部
	NO_REFERRER_WHEN_DOWNGRADE      REFERER_VALUE = "no-referrer-when-downgrade"      // 仅在从 HTTPS 站点向 HTTP 站点发送请求时不发送 Referer 头部
	ORIGIN                          REFERER_VALUE = "origin"                          // 在跨域请求时仅发送来源的域名，在同域请求时发送完整的 URL
	ORIGIN_WHEN_CROSS_ORIGIN        REFERER_VALUE = "origin-when-cross-origin"        // 仅发送来源的域名，不包含路径和查询参数
	SAME_ORIGIN                     REFERER_VALUE = "same-origin"                     // 仅在同域请求时发送 Referer 头部，在跨域请求时不发送
	STRICT_ORIGIN                   REFERER_VALUE = "strict-origin"                   // 仅发送来源的域名，并且仅在从 HTTPS 站点向 HTTPS 站点发送请求时发送
	STRICT_ORIGIN_WHEN_CROSS_ORIGIN REFERER_VALUE = "strict-origin-when-cross-origin" // 在跨域请求时仅发送来源的域名，在同域请求时发送完整的 URL，并且仅在从 HTTPS 站点向 HTTPS 站点发送请求时发送
	UNSAFE_URL                      REFERER_VALUE = "unsafe-url"                      // 发送完整的 URL，无论是同域还是跨域请求
)
