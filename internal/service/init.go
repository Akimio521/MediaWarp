package service

import (
	"MediaWarp/internal/service/alist"
	"MediaWarp/internal/utils/cache"
	"MediaWarp/pkg"
	"sync"
)

var (
	alistSeverMap sync.Map
)

// 注册Alist服务器
//
// 将Alist服务器注册到全局Map中
func RegisterAlistServer(addr string, username string, password string, cache cache.Cache) {
	alistServer := alist.New(addr, username, password, cache)
	alistServer.Init()
	alistSeverMap.Store(alistServer.GetEndpoint(), alistServer)
}

// 获取Alist服务器
//
// 从全局Map中获取Alist服务器
// 若未找到则抛出panic
func GetAlistServer(addr string) *alist.AlistServer {
	endpoint := pkg.GetEndpoint(addr)
	if server, ok := alistSeverMap.Load(endpoint); ok {
		return server.(*alist.AlistServer)
	}
	panic("Alist服务器：" + endpoint + "未注册")
}
