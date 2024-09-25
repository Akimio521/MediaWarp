package alist

type AlistResponse[T any] struct {
	Code    int64  `json:"code"`    // 状态码
	Data    T      `json:"data"`    // data
	Message string `json:"message"` // 信息
}

type AuthLoginData struct {
	Token string `json:"token"` // token
}

type FsGetData struct {
	Created  string      `json:"created"` // 创建时间
	HashInfo interface{} `json:"hash_info"`
	Hashinfo string      `json:"hashinfo"`
	Header   string      `json:"header"`
	IsDir    bool        `json:"is_dir"`   // 是否是文件夹
	Modified string      `json:"modified"` // 修改时间
	Name     string      `json:"name"`     // 文件名
	Provider string      `json:"provider"`
	RawURL   string      `json:"raw_url"` // 原始url
	Readme   string      `json:"readme"`  // 说明
	Related  interface{} `json:"related"`
	Sign     string      `json:"sign"`  // 签名
	Size     int64       `json:"size"`  // 大小
	Thumb    string      `json:"thumb"` // 缩略图
	Type     int64       `json:"type"`  // 类型
}
