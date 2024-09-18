package middleware

import (
	"MediaWarp/constants"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 记录访问日志
func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()
		method := ctx.Request.Method
		path := ctx.Request.URL.Path
		query := ctx.Request.URL.RawQuery
		if query != "" {
			path = path + "?" + query
		}

		ctx.Next()

		wasteTime := time.Since(startTime)
		clientIP := ctx.ClientIP()
		statusCode := ctx.Writer.Status()

		statusColor, methodColor := getColor(statusCode, method)

		logger.AccessLogger.Infof(
			`【AcessLog】 %s|%s| %s |%s|%s "%s"`,
			fmt.Sprintf("%-20s", startTime.Format("2006-1-2 15:04:05")),
			fmt.Sprintf("\033[4%dm %d \033[0m", statusColor, statusCode),
			fmt.Sprintf("%-10s", wasteTime),
			fmt.Sprintf("\033[4%dm %-7s\033[0m", methodColor, method),
			clientIP,
			path,
		)
	}
}

// 根据Http状态码和Http请求方法获取颜色
func getColor(statusCode int, method string) (uint8, uint8) {
	var statusColor, methodColor uint8
	switch {
	case statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices:
		statusColor = constants.COLOR_STATUS_CODE200
	case statusCode >= http.StatusMultipleChoices && statusCode < http.StatusBadRequest:
		statusColor = constants.COLOR_STATUS_CODE300
	case statusCode >= http.StatusBadRequest && statusCode < http.StatusInternalServerError:
		statusColor = constants.COLOR_STATUS_CODE400
	case statusCode >= http.StatusInternalServerError:
		statusColor = constants.COLOR_STATUS_CODE500
	default:
		statusColor = constants.COLOR_BACK
	}
	switch method {
	case http.MethodGet:
		methodColor = constants.COLOR_METHOD_GET
	case http.MethodPost:
		methodColor = constants.COLOR_METHOD_POST
	case http.MethodPut:
		methodColor = constants.COLOR_METHOD_PUT
	case http.MethodPatch:
		methodColor = constants.COLOR_METHOD_PATCH
	case http.MethodDelete:
		methodColor = constants.COLOR_METHOD_DELETE
	case http.MethodHead:
		methodColor = constants.COLOR_METHOD_HEAD
	case http.MethodOptions:
		methodColor = constants.COLOR_METHOD_OPTIONS
	default:
		methodColor = constants.COLOR_BACK
	}
	return statusColor, methodColor
}
