package emby

// /Items/:itemID/PlaybackInfo的响应
type PlaybackInfoResponse struct {
	ErrorCode     *PlaybackErrorCode `json:"ErrorCode,omitempty"`
	MediaSources  []MediaSourceInfo  `json:"MediaSources,omitempty"`
	PlaySessionID *string            `json:"PlaySessionId,omitempty"`
}
