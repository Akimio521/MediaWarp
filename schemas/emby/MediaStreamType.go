package emby

// MediaStreamType
type MediaStreamType string

const (
	Attachment    MediaStreamType = "Attachment"
	Audio         MediaStreamType = "Audio"
	Data          MediaStreamType = "Data"
	EmbeddedImage MediaStreamType = "EmbeddedImage"
	Subtitle      MediaStreamType = "Subtitle"
	Unknown       MediaStreamType = "Unknown"
	Video         MediaStreamType = "Video"
)
