package config

var cfg ConfigManager // 全局配置实例

func init() {
	cfg.Init()
}

// 获取配置
func GetConfig() *ConfigManager {
	return &cfg
}
