package emby

// SubtitleDeliveryMethod
type SubtitleDeliveryMethod string

const (
	Embed                               SubtitleDeliveryMethod = "Embed"
	Encode                              SubtitleDeliveryMethod = "Encode"
	External                            SubtitleDeliveryMethod = "External"
	HLS                                 SubtitleDeliveryMethod = "Hls"
	SubtitleDeliveryMethodVideoSideData SubtitleDeliveryMethod = "VideoSideData"
)
