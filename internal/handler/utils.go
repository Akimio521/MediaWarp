package handler

import (
	"MediaWarp/constants"
	"MediaWarp/internal/config"
	"MediaWarp/internal/logging"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"

	"github.com/andybalholm/brotli"
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

// 读取响应体
//
// 读取响应体，解压缩 GZIP、Brotli 数据（若响应体被压缩）
func readBody(rw *http.Response) ([]byte, error) {
	var data []byte
	var err error
	encodingType := rw.Header.Get("Content-Encoding")
	switch encodingType {
	case "": // 无 Content-Encoding 头
		logging.Debug("无 Content-Encoding 头")
		if data, err = io.ReadAll(rw.Body); err != nil { // 不是压缩格式，直接读取数据
			return nil, fmt.Errorf("读取 Body 出错：%w", err)
		}
	case "gzip":
		logging.Debug("解压 GZIP 数据")
		// 解压缩 GZIP 数据
		gzipReader, err := gzip.NewReader(rw.Body)
		if err != nil {
			return nil, fmt.Errorf("创建 GZIP 解压器失败：%w", err)
		}
		defer gzipReader.Close()

		if data, err = io.ReadAll(gzipReader); err != nil { // 读取解压后的数据
			return nil, fmt.Errorf("读取解压后的数据失败：%w", err)
		}
	case "br":
		logging.Debug("解压 Brotli 数据")
		// 解压 Brotli 数据
		brotliReader := brotli.NewReader(rw.Body)
		if data, err = io.ReadAll(brotliReader); err != nil {
			return nil, fmt.Errorf("读取 Brotli 解压后的数据失败：%w", err)
		}
	default:
		return nil, fmt.Errorf("未知的 Content-Encoding：%s", encodingType)
	}
	return data, nil
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
