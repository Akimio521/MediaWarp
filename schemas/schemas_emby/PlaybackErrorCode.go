package schemas_emby

// PlaybackErrorCode
type PlaybackErrorCode string

const (
	NoCompatibleStream PlaybackErrorCode = "NoCompatibleStream"
	NotAllowed         PlaybackErrorCode = "NotAllowed"
	RateLimitExceeded  PlaybackErrorCode = "RateLimitExceeded"
)
