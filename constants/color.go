package constants

const (
	COLOR_BACK   uint8 = iota // 黑色
	COLOR_RED                 // 红色
	COLOR_GREEN               // 绿色
	COLOR_YELLOW              // 黄色
	COLOR_BLUE                // 蓝色
	COLOR_PURPLE              // 紫色
	COLOR_CYAN                // 青色
	COLOR_GRAY                // 灰色
)

const (
	COLOR_STATUS_CODE200 = COLOR_GREEN  // HTTP 请求成功状态颜色
	COLOR_STATUS_CODE300 = COLOR_GRAY   // HTTP 重定向状态颜色
	COLOR_STATUS_CODE400 = COLOR_YELLOW // HTTP 客户端错误颜色
	COLOR_STATUS_CODE500 = COLOR_RED    // HTTP 服务器错误颜色
)

const (
	COLOR_METHOD_GET     = COLOR_BLUE   // HTTP GET 请求方法颜色
	COLOR_METHOD_POST    = COLOR_CYAN   // HTTP POST 请求方法颜色
	COLOR_METHOD_PUT     = COLOR_YELLOW // HTTP PUT 请求方法颜色
	COLOR_METHOD_PATCH   = COLOR_GREEN  // HTTP PATCH 请求方法颜色
	COLOR_METHOD_DELETE  = COLOR_RED    // HTTP DELETE 请求方法颜色
	COLOR_METHOD_HEAD    = COLOR_PURPLE // HTTP HEAD 请求方法颜色
	COLOR_METHOD_OPTIONS = COLOR_GRAY   // HTTP OPTIONS 请求方法颜色
)
