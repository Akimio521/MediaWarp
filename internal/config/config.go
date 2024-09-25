package config

import (
	"MediaWarp/constants"
	"MediaWarp/internal/service"
	"MediaWarp/internal/utils/cacheutils"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// 上游媒体服务器相关设置
type MeidaServerSetting struct {
	Type constants.MediaServerType
	ADDR string
	AUTH string
}

// 日志设置
type LoggerSetting struct {
	Enable        bool              // 是否启用日志
	AccessLogger  BaseLoggerSetting // 访问日志相关配置
	ServiceLogger BaseLoggerSetting // 服务日志相关配置
}

// 基础日志配置字段
type BaseLoggerSetting struct {
	Enable bool // 是否启用日志
	File   bool // 是否将日志输出到文件中
}

// Web前端自定义设置
type WebSetting struct {
	Enable            bool
	Static            bool
	Index             bool
	Head              bool
	ExternalPlayerUrl bool
	ActorPlus         bool
	FanartShow        bool
	Danmaku           bool
	BeautifyCSS       bool
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
	Port         int
	CacheType    constants.CacheType
	Cache        cacheutils.Cache
	MeidaServer  MeidaServerSetting
	Logger       LoggerSetting
	Web          WebSetting
	ClientFilter ClientFilterSetting
	HTTPStrm     HTTPStrmSetting
	AlistStrm    AlistStrmSetting
}

// 读取并解析配置文件
func (configManager *ConfigManager) loadConfig() {
	var (
		vip = viper.New()
		err error
	)

	vip.SetConfigFile(configManager.ConfigPath())
	vip.SetConfigType("yaml")

	if err = vip.ReadInConfig(); err != nil {
		panic(err)
	}

	configManager.MeidaServer.Type = constants.MediaServerType(vip.GetString("Server.Type"))
	configManager.CacheType = constants.CacheType(vip.GetString("Cache.Type"))
	configManager.Cache = cacheutils.GetCache(configManager.CacheType)

	err = vip.Unmarshal(configManager)
	if err != nil {
		panic(err)
	}
}

// 创建文件夹
func (*ConfigManager) createDir() {
	if err := os.MkdirAll(cfg.ConfigDir(), os.ModePerm); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(cfg.LogDir(), os.ModePerm); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(cfg.StaticDir(), os.ModePerm); err != nil {
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

// 获取项目根目录
func (*ConfigManager) RootDir() string {
	executablePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(executablePath)
}

// 获取配置文件目录
func (configManager *ConfigManager) ConfigDir() string {
	return filepath.Join(configManager.RootDir(), "config")
}

// 获取配置文件路径
func (configManager *ConfigManager) ConfigPath() string {
	return filepath.Join(configManager.ConfigDir(), "config.yaml")
}

// 获取日志目录
func (configManager *ConfigManager) LogDir() string {
	return filepath.Join(configManager.RootDir(), "logs")
}

// 获取访问日志文件路径
func (configManager *ConfigManager) AccessLogPath() string {
	return filepath.Join(configManager.LogDir(), "access.log")
}

// 获取服务日志文件路径
func (configManager *ConfigManager) ServiceLogPath() string {
	return filepath.Join(configManager.LogDir(), "service.log")
}

// 获取静态资源文件目录
func (configManager *ConfigManager) StaticDir() string {
	return filepath.Join(configManager.RootDir(), "static")
}

// MediaWarp监听地址
//
// 监听所有网卡
func (configManager *ConfigManager) ListenAddr() string {
	return fmt.Sprintf(":%d", configManager.Port)
}
