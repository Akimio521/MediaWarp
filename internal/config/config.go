package config

import (
	"MediaWarp/constants"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/viper"
)

var (
	version = VersionInfo{
		AppVersion: appVersion,
		CommitHash: commitHash,
		BuildData:  buildDate,
		GoVersion:  runtime.Version(),
		OS:         runtime.GOOS,
		Arch:       runtime.GOARCH,
	}

	Port         int                 // MediaWarp开放端口
	MediaServer  MediaServerSetting  // 上游媒体服务器设置
	Logger       LoggerSetting       // 日志设置
	Web          WebSetting          // Web服务器设置
	ClientFilter ClientFilterSetting // 客户端过滤设置
	HTTPStrm     HTTPStrmSetting     // HTTPSTRM设置
	AlistStrm    AlistStrmSetting    // AlistStrm设置
	Subtitle     SubtitleSetting     // 字幕设置
)

// 获取版本信息
func Version() *VersionInfo {
	return &version
}

// 二进制文件目录
func RootDir() string {
	executablePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(executablePath)
}

// 配置文件目录
func ConfigDir() string {
	return filepath.Join(RootDir(), "config")
}

// 配置文件路径
func ConfigPath() string {
	return filepath.Join(ConfigDir(), "config.yaml")
}

// 获取日志目录
//
// 总日志目录
// ./logs
func LogDir() string {
	return filepath.Join(RootDir(), "logs")
}

// 获取日志目录
//
// 带有日期
// ./logs/2024-9-29
func LogDirWithDate() string {
	return filepath.Join(LogDir(), time.Now().Format("2006-1-2"))
}

// 访问日志文件路径
func AccessLogPath() string {
	return filepath.Join(LogDirWithDate(), "access.log")
}

// 服务日志文件路径
func ServiceLogPath() string {
	return filepath.Join(LogDirWithDate(), "service.log")
}

// 静态资源文件目录
//
// 用户自定义静态文件存放地址
func CostomDir() string {
	return filepath.Join(RootDir(), "static")
}

// MediaWarp监听地址
//
// 监听所有网卡
func ListenAddr() string {
	return fmt.Sprintf(":%d", Port)
}

// ====================== 私有函数 ======================

// 初始化configManager
func Init() error {
	if err := loadConfig(); err != nil {
		return err
	}
	if err := createDir(); err != nil {
		return err
	}
	return nil
}

// 读取并解析配置文件
func loadConfig() error {
	viper.AddConfigPath(ConfigDir())
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	Port = viper.GetInt("Port")
	MediaServer.Type = constants.MediaServerType(viper.GetString("MediaServer.Type"))
	MediaServer.ADDR = viper.GetString("MediaServer.ADDR")
	MediaServer.AUTH = viper.GetString("MediaServer.AUTH")

	if err := viper.UnmarshalKey("Logger", &Logger); err != nil {
		return fmt.Errorf("LoggerSetting 解析失败, %v", err)
	}
	if err := viper.UnmarshalKey("Web", &Web); err != nil {
		return fmt.Errorf("WebSetting  解析失败, %v", err)
	}
	if err := viper.UnmarshalKey("ClientFilter", &ClientFilter); err != nil {
		return fmt.Errorf("ClientFilterSetting  解析失败, %v", err)
	}
	if err := viper.UnmarshalKey("HTTPStrm", &HTTPStrm); err != nil {
		return fmt.Errorf("HTTPStrmSetting  解析失败, %v", err)
	}
	if err := viper.UnmarshalKey("AlistStrm", &AlistStrm); err != nil {
		return fmt.Errorf("AlistStrmSetting  解析失败, %v", err)
	}
	if err := viper.UnmarshalKey("Subtitle", &Subtitle); err != nil {
		return fmt.Errorf("SubtitleSetting  解析失败, %v", err)
	}
	return nil
}

// 创建文件夹
func createDir() error {
	if err := os.MkdirAll(ConfigDir(), os.ModePerm); err != nil {
		return fmt.Errorf("创建配置文件夹失败: %v", err)
	}
	if err := os.MkdirAll(LogDir(), os.ModePerm); err != nil {
		return fmt.Errorf("创建日志文件夹失败: %v", err)
	}
	if err := os.MkdirAll(CostomDir(), os.ModePerm); err != nil {
		panic(err)
	}
	return nil
}
