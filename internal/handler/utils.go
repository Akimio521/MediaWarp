package handler

import (
	"MediaWarp/constants"
	"MediaWarp/internal/config"
	"MediaWarp/internal/logging"
	"bytes"
	"io"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// 响应修改创建器
//
// 将需要修改上游响应的处理器包装成一个 gin.HandlerFunc 处理器
func responseModifyCreater(proxy *httputil.ReverseProxy, modifyResponseFN func(rw *http.Response) error) gin.HandlerFunc {
	proxy.ModifyResponse = modifyResponseFN
	return func(ctx *gin.Context) {
		proxy.ServeHTTP(ctx.Writer, ctx.Request)
	}
}

// 根据 Strm 文件路径识别 Strm 文件类型
//
// 返回 Strm 文件类型和一个可选配置
func recgonizeStrmFileType(strmFilePath string) (constants.StrmFileType, any) {
	if config.HTTPStrm.Enable {
		for _, prefix := range config.HTTPStrm.PrefixList {
			if strings.HasPrefix(strmFilePath, prefix) {
				logging.Debug(strmFilePath + " 成功匹配路径：" + prefix + "，Strm 类型：" + string(constants.HTTPStrm))
				return constants.HTTPStrm, nil
			}
		}
	}
	if config.AlistStrm.Enable {
		for _, alistStrmConfig := range config.AlistStrm.List {
			for _, prefix := range alistStrmConfig.PrefixList {
				if strings.HasPrefix(strmFilePath, prefix) {
					logging.Debug(strmFilePath + " 成功匹配路径：" + prefix + "，Strm 类型：" + string(constants.AlistStrm) + "，AlistServer 地址：" + alistStrmConfig.ADDR)
					return constants.AlistStrm, alistStrmConfig.ADDR
				}
			}
		}
	}
	logging.Debug(strmFilePath + " 未匹配任何路径，Strm 类型：" + string(constants.UnknownStrm))
	return constants.UnknownStrm, nil
}

// 更新响应体
//
// 修改响应体、更新Content-Length
func updateBody(rw *http.Response, content []byte) {
	rw.Body = io.NopCloser(bytes.NewBuffer(content)) // 重置响应体

	// 更新 Content-Length 头
	rw.ContentLength = int64(len(content))
	rw.Header.Set("Content-Length", strconv.Itoa(len(content)))
}
