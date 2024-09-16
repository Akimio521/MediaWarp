package handlers

import (
	"MediaWarp/core"
	"net/url"
)

var (
	config                 = core.GetConfig()
	mediaServerEndpoint, _ = url.Parse(config.Server.GetEndpoint())
)
