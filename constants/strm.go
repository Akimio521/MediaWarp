package constants

type StrmFileType uint8 // Strm 文件类型

const (
	UnknownStrm StrmFileType = iota
	HTTPStrm
	AlistStrm
)

func (s StrmFileType) String() string {
	switch s {
	case HTTPStrm:
		return "HTTPStrm"
	case AlistStrm:
		return "AlistStrm"
	default:
		return "UnknownStrm"
	}
}
