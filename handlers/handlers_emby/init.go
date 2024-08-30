package handlers_emby

import (
	"MediaWarp/api"
	"MediaWarp/core"
)

var (
	config     = core.GetConfig()
	logger     = core.GetLogger()
	embyServer = config.Server.(*api.EmbyServer)
)
