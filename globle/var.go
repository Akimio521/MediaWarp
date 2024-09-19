package globle

import (
	"MediaWarp/constants"
	"MediaWarp/pkg/cachepkg"
)

// 全局变量存放地
var (
	GlobleCache cachepkg.Cache
)

func init() {
	GlobleCache = cachepkg.AvailableCaches[constants.MEMORY_CACHE]
}
