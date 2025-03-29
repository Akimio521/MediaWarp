package constants // 包名保持简洁，若范围明确可保留

// ReferrerPolicy 定义 HTTP Referer 策略类型
type ReferrerPolicy string

const (
	NoReferrer                  ReferrerPolicy = "no-referrer"                     // 不发送 Referer 头部
	NoReferrerWhenDowngrade     ReferrerPolicy = "no-referrer-when-downgrade"      // 从 HTTPS 到 HTTP 时不发送
	Origin                      ReferrerPolicy = "origin"                          // 跨域发送域名，同域发送完整 URL
	OriginWhenCrossOrigin       ReferrerPolicy = "origin-when-cross-origin"        // 跨域仅发送域名
	SameOrigin                  ReferrerPolicy = "same-origin"                     // 仅同域发送
	StrictOrigin                ReferrerPolicy = "strict-origin"                   // HTTPS 到 HTTPS 时发送域名
	StrictOriginWhenCrossOrigin ReferrerPolicy = "strict-origin-when-cross-origin" // 跨域发送域名，同域完整 URL（仅 HTTPS）
	UnsafeURL                   ReferrerPolicy = "unsafe-url"                      // 始终发送完整 URL
)
