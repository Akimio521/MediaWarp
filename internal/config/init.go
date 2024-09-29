package config

var cfg *ConfigManager // 全局配置实例

// 实例化一个ConfigManager对象
func NewConfig() *ConfigManager {
	cm := ConfigManager{}
	cm.Init()
	return &cm
}

// 获取全局ConfigManager对象
func GetConfig() *ConfigManager {
	if cfg == nil {
		cfg = NewConfig()
	}
	return cfg
}
