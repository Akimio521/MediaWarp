package schemas_alist

type AlistResponse struct {
	Code    int64       `json:"code"`    // 状态码
	Data    interface{} `json:"data"`    // data
	Message string      `json:"message"` // 信息
}
