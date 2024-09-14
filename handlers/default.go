package handlers

import (
	"MediaWarp/core"
	"MediaWarp/pkg"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	config             = core.GetConfig()
	remoteHTTP, _      = url.Parse(config.Server.GetHTTPEndpoint())
	remoteWebSocket, _ = url.Parse(config.Server.GetWebSocketEndpoint())
	upGrader           = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// 默认路由（直接转发请求到后端）
func DefaultHandler(ctx *gin.Context) {
	if websocket.IsWebSocketUpgrade(ctx.Request) {
		reverseProxyWebSocketHandler(ctx)
	} else {
		reverseProxyHTTPHandler(ctx)
	}
}

// 处理HTTP请求
func reverseProxyHTTPHandler(ctx *gin.Context) {
	hostName, _ := pkg.SplitHostPort(ctx.Request.Host)
	proxy := httputil.NewSingleHostReverseProxy(remoteHTTP)
	proxy.Director = func(req *http.Request) {
		req.Header = ctx.Request.Header
		req.Host = hostName //remote.Host
		req.URL.Scheme = remoteHTTP.Scheme
		req.URL.Host = remoteHTTP.Host
		req.URL.Path = ctx.Request.URL.Path
	}
	proxy.ServeHTTP(ctx.Writer, ctx.Request)
}

// 处理WebSocket请求
func reverseProxyWebSocketHandler(ctx *gin.Context) {
	clientConn, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer clientConn.Close()

	wsURL := remoteWebSocket.String() + ctx.Request.URL.Path + "?" + ctx.Request.URL.RawQuery
	serverConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer serverConn.Close()

	// 转发消息
	go func() {
		for {
			messageType, message, err := clientConn.ReadMessage()
			if err != nil {
				break
			}
			err = serverConn.WriteMessage(messageType, message)
			if err != nil {
				break
			}
		}
		serverConn.Close()
	}()

	for {
		messageType, message, err := serverConn.ReadMessage()
		if err != nil {
			break
		}
		err = clientConn.WriteMessage(messageType, message)
		if err != nil {
			break
		}
	}
	clientConn.Close()
}
