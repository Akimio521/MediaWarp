package config

var cfg *ConfigManager // 全局配置实例

// 获取全局ConfigManager对象
func GetConfig() *ConfigManager {
	if cfg == nil {
		cfg = &ConfigManager{}
		cfg.Init()
	}
	return cfg
}
