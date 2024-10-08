package config

import (
	"MediaWarp/constants"
	"MediaWarp/internal/service"
	"MediaWarp/internal/utils/cache"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

// 上游媒体服务器相关设置
type MediaServerSetting struct {
	Type constants.MediaServerType // 媒体服务器类型
	ADDR string                    // 地址
	AUTH string                    // 认证授权KEY
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
	BeautifyCSS       bool   // Emby美化CSS样式
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
	PrefixList []string
}

// AlistStrm具体设置
type AlistSetting struct {
	ADDR       string
	Username   string
	Password   string
	PrefixList []string
}

// AlistStrm播放设置
type AlistStrmSetting struct {
	Enable bool
	List   []AlistSetting
}

type ConfigManager struct {
	Port         int                 // MediaWarp开放端口
	CacheType    constants.CacheType // 缓存类型
	Cache        cache.Cache         // 全局缓存接口
	MediaServer  MediaServerSetting  // 上游媒体服务器设置
	Logger       LoggerSetting       // 日志设置
	Web          WebSetting          // Web服务器设置
	ClientFilter ClientFilterSetting // 客户端过滤设置
	HTTPStrm     HTTPStrmSetting     // HTTPSTRM设置
	AlistStrm    AlistStrmSetting    // AlistStrm设置
}

// 读取并解析配置文件
func (configManager *ConfigManager) loadConfig() {
	viper.SetConfigFile(configManager.ConfigPath())
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	configManager.MediaServer.Type = constants.MediaServerType(viper.GetString("Server.Type"))
	configManager.CacheType = constants.CacheType(viper.GetString("Cache.Type"))

	configManager.Cache = cache.GetCache(configManager.CacheType)

	if err := viper.Unmarshal(configManager); err != nil {
		panic(err)
	}
}

// 创建文件夹
func (configManager *ConfigManager) createDir() {
	if err := os.MkdirAll(configManager.ConfigDir(), os.ModePerm); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(configManager.LogDir(), os.ModePerm); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(configManager.StaticDir(), os.ModePerm); err != nil {
		panic(err)
	}
}

// 注册Alist服务器
func (configManager *ConfigManager) registerAlistServer() {
	for _, alistConfig := range configManager.AlistStrm.List {
		service.RegisterAlistServer(alistConfig.ADDR, alistConfig.Username, alistConfig.Password, configManager.Cache)
	}
}

// 初始化configManager
func (configManager *ConfigManager) Init() {
	configManager.loadConfig()
	configManager.createDir()
	configManager.registerAlistServer()
}

// 获取版本号
func (*ConfigManager) Version() string {
	return APP_VERSION
}

// 项目根目录
//
// 二进制可执行文件存放地址
func (*ConfigManager) RootDir() string {
	executablePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(executablePath)
}

// 配置文件目录
func (configManager *ConfigManager) ConfigDir() string {
	return filepath.Join(configManager.RootDir(), "config")
}

// 配置文件路径
func (configManager *ConfigManager) ConfigPath() string {
	return filepath.Join(configManager.ConfigDir(), "config.yaml")
}

// 获取日志目录
//
// 总日志目录
// ./logs
func (configManager *ConfigManager) LogDir() string {
	return filepath.Join(configManager.RootDir(), "logs")
}

// 获取日志目录
//
// 带有日期
// ./logs/2024-9-29
func (configManager *ConfigManager) LogDirWithDate() string {
	return filepath.Join(configManager.LogDir(), time.Now().Format("2006-1-2"))
}

// 访问日志文件路径
func (configManager *ConfigManager) AccessLogPath() string {
	return filepath.Join(configManager.LogDirWithDate(), "access.log")
}

// 服务日志文件路径
func (configManager *ConfigManager) ServiceLogPath() string {
	return filepath.Join(configManager.LogDirWithDate(), "service.log")
}

// 静态资源文件目录
//
// 用户自定义静态文件存放地址
func (configManager *ConfigManager) StaticDir() string {
	return filepath.Join(configManager.RootDir(), "static")
}

// MediaWarp监听地址
//
// 监听所有网卡
func (configManager *ConfigManager) ListenAddr() string {
	return fmt.Sprintf(":%d", configManager.Port)
}
