package core

import (
	"MediaWarp/api"
	"MediaWarp/constants"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type baseConfig struct {
	Enable bool
}

type loggerConfig struct {
	Enable        bool
	AccessLogger  baseConfig
	ServiceLogger baseConfig
}

type webConfig struct {
	Enable            bool
	Static            bool
	Index             bool
	Head              bool
	ExternalPlayerUrl bool
	BeautifyCSS       bool
}

type clientFilterConfig struct {
	Enable     bool
	Mode       constants.FliterMode
	ClientList []string
}

type alistStrmConfig struct {
	Enable bool
	List   []baseAlistStrmConfig
}

type baseAlistStrmConfig struct {
	AlistServer api.AlistServer
	PrefixList  []string
}

type httpStrmConfig struct {
	Enable     bool
	PrefixList []string
}

type configManager struct {
	Host          string
	Port          int
	Server        api.MediaServer
	LoggerSetting loggerConfig
	Web           webConfig
	ClientFilter  clientFilterConfig
	HttpStrm      httpStrmConfig
	AlistStrm     alistStrmConfig
}

// 读取并解析配置文件
func (config *configManager) LoadConfig() {
	vip := viper.New()
	vip.SetConfigFile(config.ConfigPath())
	vip.SetConfigType("yaml")

	if err := vip.ReadInConfig(); err != nil {
		panic(err)
	}
	serverType := constants.ServerType(vip.GetString("Server.TYPE"))
	switch serverType {
	case constants.EMBY:
		config.Server = &api.EmbyServer{
			ADDR:  vip.GetString("Server.ADDR"),
			TOKEN: vip.GetString("Server.TOKEN"),
		}
	default:
		panic("未知的媒体服务器类型：" + serverType)
	}

	err := vip.Unmarshal(config)
	if err != nil {
		panic(err)
	}
}

// 创建文件夹
func (config *configManager) CreateDir() {
	if err := os.MkdirAll(config.ConfigDir(), os.ModePerm); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(config.LogDir(), os.ModePerm); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(config.StaticDir(), os.ModePerm); err != nil {
		panic(err)
	}
}

// 初始化configManager
func (config *configManager) Init() {
	config.LoadConfig()
	config.CreateDir()
}

// 获取版本号
func (config *configManager) Version() string {
	return APP_VERSION
}

// 获取项目根目录
func (c *configManager) RootDir() string {
	executablePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(executablePath)
}

// 获取配置文件目录
func (config *configManager) ConfigDir() string {
	return filepath.Join(config.RootDir(), "config")
}

// 获取配置文件路径
func (config *configManager) ConfigPath() string {
	return filepath.Join(config.ConfigDir(), "config.yaml")
}

// 获取日志目录
func (config *configManager) LogDir() string {
	return filepath.Join(config.RootDir(), "logs")
}

// 获取访问日志文件路径
func (config *configManager) AccessLogPath() string {
	return filepath.Join(config.LogDir(), "access.log")
}

// 获取服务日志文件路径
func (config *configManager) ServiceLogPath() string {
	return filepath.Join(config.LogDir(), "service.log")
}

// 获取静态资源文件目录
func (config *configManager) StaticDir() string {
	return filepath.Join(config.RootDir(), "static")
}

// MediaWarp监听地址
func (config *configManager) ListenAddr() string {
	return fmt.Sprintf("%s:%d", config.Host, config.Port)
}

// -----------------外部引用部分----------------- //
var config configManager

func GetConfig() *configManager {
	return &config
}

func init() {
	config.Init()
}
