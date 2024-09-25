package router

import "MediaWarp/internal/config"

var (
	cfg *config.ConfigManager
)

func init() {
	cfg = config.GetConfig()

}
