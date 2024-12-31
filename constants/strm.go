package constants

type StrmFileType string // Strm 文件类型

var (
	HTTPStrm    StrmFileType = "HTTPStrm"
	AlistStrm   StrmFileType = "AlistStrm"
	AlistHTTPStrm   StrmFileType = "AlistHTTPStrm"
	UnknownStrm StrmFileType = "UnknownStrm"
)
