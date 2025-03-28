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
	encoding := rw.Header.Get("Content-Encoding")

	var reader io.Reader
	switch encoding {
	case "gzip":
		logging.Debug("解码 GZIP 数据")
		gr, err := gzip.NewReader(rw.Body)
		if err != nil {
			return nil, fmt.Errorf("gzip reader error: %w", err)
		}
		defer gr.Close()
		reader = gr

	case "br":
		logging.Debug("解码 Brotli 数据")
		reader = brotli.NewReader(rw.Body)

	case "": // 无压缩
		logging.Debug("无压缩数据")
		reader = rw.Body

	default:
		return nil, fmt.Errorf("unsupported Content-Encoding: %s", encoding)
	}
	return io.ReadAll(reader)
}

// 更新响应体
//
// 修改响应体、更新Content-Length
func updateBody(rw *http.Response, content []byte) error {
	encoding := rw.Header.Get("Content-Encoding")
	var (
		compressed bytes.Buffer
		writer     io.Writer
	)

	// 根据原始编码选择压缩方式
	switch encoding {
	case "gzip":
		logging.Debug("使用 GZIP 重新编码数据")
		gw := gzip.NewWriter(&compressed)
		defer gw.Close()
		writer = gw

	case "br":
		logging.Debug("使用 Brotli 重新编码数据")
		bw := brotli.NewWriter(&compressed)
		defer bw.Close()
		writer = bw

	case "": // 无压缩
		logging.Debug("无压缩数据")
		writer = &compressed

	default:
		logging.Warningf("不支持的重新编码：%s，将不对数据进行压缩编码", encoding)
		rw.Header.Del("Content-Encoding")
	}

	if _, err := writer.Write(content); err != nil {
		return fmt.Errorf("compression write error: %w", err)
	}

	// Brotli 需要显式 Flush
	if bw, ok := writer.(*brotli.Writer); ok {
		if err := bw.Flush(); err != nil {
			return err
		}
	}

	// 设置新 Body
	rw.Body = io.NopCloser(bytes.NewReader(compressed.Bytes()))
	rw.ContentLength = int64(compressed.Len())
	rw.Header.Set("Content-Length", strconv.Itoa(compressed.Len())) // 更新响应头

	return nil
}
