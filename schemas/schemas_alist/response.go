package schemas_alist

type AlistResponse[T any] struct {
	Code    int64  `json:"code"`    // 状态码
	Data    T      `json:"data"`    // data
	Message string `json:"message"` // 信息
}
