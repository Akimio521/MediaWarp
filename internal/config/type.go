package config

import "MediaWarp/constants"

// 程序版本信息
type VersionInfo struct {
	AppVersion string // 程序版本号
	CommitHash string // GIt Commit Hash
	BuildData  string // 编译时间
	GoVersion  string // 编译 Golang 版本
	OS         string // 操作系统
	Arch       string //  架构
}

// 上游媒体服务器相关设置
type MediaServerSetting struct {
	Type constants.MediaServerType // 媒体服务器类型
	ADDR string                    // 地址
	AUTH string                    // 认证授权KEY
}

// 缓存相关设置
type CacheSetting struct {
	Type     constants.CacheType // 缓存类型
	WebCache bool                // 是否启用Web缓存中间件
}

// 日志设置
type LoggerSetting struct {
	AccessLogger  BaseLoggerSetting // 访问日志相关配置
	ServiceLogger BaseLoggerSetting // 服务日志相关配置
}

// 基础日志配置字段
type BaseLoggerSetting struct {
	Console bool // 是否将日志输出到终端中
	File    bool // 是否将日志输出到文件中
}

// Web前端自定义设置
type WebSetting struct {
	Enable            bool   // 启用自定义前端设置
	Custom            bool   // 启用用户自定义静态资源
	Index             bool   // 是否从static目录读取index.html文件作为首页
	Head              string // 添加到index.html的HEAD中
	ExternalPlayerUrl bool   // 是否开启外置播放器
	ActorPlus         bool   // 过滤没有头像的演员和制作人员
	FanartShow        bool   // 显示同人图（fanart图）
	Danmaku           bool   // Web显示弹幕
	VideoTogether     bool   // VideoTogether
}

// 客户端User-Agent过滤设置
type ClientFilterSetting struct {
	Enable     bool
	Mode       constants.FliterMode
	ClientList []string
}

// HTTPStrm播放设置
type HTTPStrmSetting struct {
	Enable     bool
	TransCode  bool // false->强制关闭转码 true->保持原有转码设置
	PrefixList []string
}

// AlistStrm具体设置
type AlistSetting struct {
	ADDR       string
	Username   string
	Password   string
	Token      *string
	PrefixList []string
}

// AlistStrm播放设置
type AlistStrmSetting struct {
	Enable    bool
	TransCode bool // false->强制关闭转码 true->保持原有转码设置
	RawURL    bool // 是否使用原始 URL
	List      []AlistSetting
}

// 字幕设置
type SubtitleSetting struct {
	Enable   bool
	SRT2ASS  bool // SRT 字幕转 ASS 字幕
	ASSStyle []string
	SubSet   bool // ASS 字幕字体子集化
}
