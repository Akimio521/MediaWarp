package config

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

type configManager struct {
	Server struct {
		Host string
		Port int
	}
	Origin string
	ApiKey string
}

// 获取版本号
func (c *configManager) Version() string {
	return APP_VERSION
}

// 获取项目根目录
func (c *configManager) RootDir() string {
	_, fullFilename, _, _ := runtime.Caller(0)
	return filepath.Dir(filepath.Dir(fullFilename))
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

// MediaWarp监听地址
func (c *configManager) ListenAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// -----------------外部引用部分----------------- //
var config configManager

func init() {
	config.LoadConfig()
}

func GetConfig() *configManager {
	return &config
}
