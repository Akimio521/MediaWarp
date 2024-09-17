package router

import (
	"MediaWarp/core"
	"regexp"

	"github.com/gin-gonic/gin"
)

type RegexpRouterPair struct {
	Regexp  *regexp.Regexp
	Handler gin.HandlerFunc
}

var (
	config          = core.GetConfig()
	embyRouterRules = initEmbyRouterRule()
)
