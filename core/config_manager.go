package core

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type serverConfig struct {
	Host string
	Port int
}

type loggerConfig struct {
	Enable        bool
	AccessLogger  baseLoggerConfig
	ServiceLogger baseLoggerConfig
}

type baseLoggerConfig struct {
	Enable bool
}

type configManager struct {
	Server        serverConfig
	LoggerSetting loggerConfig
	Origin        string
	ApiKey        string
	Static        bool
}

// 读取并解析配置文件
func (c *configManager) LoadConfig() {
	vip := viper.New()
	vip.SetConfigFile(c.ConfigPath())
	vip.SetConfigType("yaml")

	if err := vip.ReadInConfig(); err != nil {
		panic(err)
	}

	err := vip.Unmarshal(c)
	if err != nil {
		panic(err)
	}
}

// 创建文件夹
func (c *configManager) CreateDir() {
	if err := os.MkdirAll(c.ConfigDir(), os.ModePerm); err != nil {
		panic(err)
	}
	if err := os.MkdirAll(c.LogDir(), os.ModePerm); err != nil {
		panic(err)
	}
}

// 初始化configManager
func (c *configManager) Init() {
	c.LoadConfig()
	c.CreateDir()
}

// 获取版本号
func (c *configManager) Version() string {
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
func (c *configManager) ConfigDir() string {
	return filepath.Join(c.RootDir(), "config")
}

// 获取配置文件路径
func (c *configManager) ConfigPath() string {
	return filepath.Join(c.ConfigDir(), "config.yaml")
}

// 获取日志目录
func (c *configManager) LogDir() string {
	return filepath.Join(c.RootDir(), "logs")
}

// 获取访问日志文件路径
func (c *configManager) AccessLogPath() string {
	return filepath.Join(c.LogDir(), "access.log")
}

// 获取服务日志文件路径
func (c *configManager) ServiceLogPath() string {
	return filepath.Join(c.LogDir(), "service.log")
}

// 获取静态资源文件目录
func (c *configManager) StaticDir() string {
	return filepath.Join(c.RootDir(), "static")
}

// MediaWarp监听地址
func (c *configManager) ListenAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// -----------------外部引用部分----------------- //
var config configManager

func GetConfig() *configManager {
	return &config
}

func init() {
	config.Init()
}
