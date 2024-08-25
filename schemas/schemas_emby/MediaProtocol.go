package schemas_emby

// MediaProtocol
type MediaProtocol string

const (
	FTP  MediaProtocol = "Ftp"
	File MediaProtocol = "File"
	HTTP MediaProtocol = "Http"
	Mms  MediaProtocol = "Mms"
	RTP  MediaProtocol = "Rtp"
	RTSP MediaProtocol = "Rtsp"
	Rtmp MediaProtocol = "Rtmp"
	UDP  MediaProtocol = "Udp"
)
