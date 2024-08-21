package emby

// MediaSourceInfo
type MediaSourceInfo struct {
	AddAPIKeyToDirectStreamURL *bool                     `json:"AddApiKeyToDirectStreamUrl,omitempty"`
	AnalyzeDurationMS          *int64                    `json:"AnalyzeDurationMs"`
	Bitrate                    *int64                    `json:"Bitrate"`
	BufferMS                   *int64                    `json:"BufferMs"`
	Container                  *string                   `json:"Container,omitempty"`
	ContainerStartTimeTicks    *int64                    `json:"ContainerStartTimeTicks"`
	DefaultAudioStreamIndex    *int64                    `json:"DefaultAudioStreamIndex"`
	DefaultSubtitleStreamIndex *int64                    `json:"DefaultSubtitleStreamIndex"`
	DirectStreamURL            *string                   `json:"DirectStreamUrl,omitempty"`
	EncoderPath                *string                   `json:"EncoderPath,omitempty"`
	EncoderProtocol            *MediaProtocol            `json:"EncoderProtocol,omitempty"`
	Formats                    []string                  `json:"Formats,omitempty"`
	HasMixedProtocols          *bool                     `json:"HasMixedProtocols,omitempty"`
	ID                         *string                   `json:"Id,omitempty"`
	IsInfiniteStream           *bool                     `json:"IsInfiniteStream,omitempty"`
	IsRemote                   *bool                     `json:"IsRemote,omitempty"`
	ItemID                     *string                   `json:"ItemId,omitempty"`
	LiveStreamID               *string                   `json:"LiveStreamId,omitempty"`
	MediaStreams               []MediaStream             `json:"MediaStreams,omitempty"`
	Name                       *string                   `json:"Name,omitempty"`
	OpenToken                  *string                   `json:"OpenToken,omitempty"`
	Path                       *string                   `json:"Path,omitempty"`
	ProbePath                  *string                   `json:"ProbePath,omitempty"`
	ProbeProtocol              *MediaProtocol            `json:"ProbeProtocol,omitempty"`
	Protocol                   *MediaProtocol            `json:"Protocol,omitempty"`
	ReadAtNativeFramerate      *bool                     `json:"ReadAtNativeFramerate,omitempty"`
	RequiredHTTPHeaders        map[string]string         `json:"RequiredHttpHeaders,omitempty"`
	RequiresClosing            *bool                     `json:"RequiresClosing,omitempty"`
	RequiresLooping            *bool                     `json:"RequiresLooping,omitempty"`
	RequiresOpening            *bool                     `json:"RequiresOpening,omitempty"`
	RunTimeTicks               *int64                    `json:"RunTimeTicks"`
	ServerID                   *string                   `json:"ServerId,omitempty"`
	Size                       *int64                    `json:"Size"`
	SortName                   *string                   `json:"SortName,omitempty"`
	SupportsDirectPlay         *bool                     `json:"SupportsDirectPlay,omitempty"`
	SupportsDirectStream       *bool                     `json:"SupportsDirectStream,omitempty"`
	SupportsProbing            *bool                     `json:"SupportsProbing,omitempty"`
	SupportsTranscoding        *bool                     `json:"SupportsTranscoding,omitempty"`
	Timestamp                  *TransportStreamTimestamp `json:"Timestamp,omitempty"`
	TrancodeLiveStartIndex     *int64                    `json:"TrancodeLiveStartIndex"`
	TranscodingContainer       *string                   `json:"TranscodingContainer,omitempty"`
	TranscodingSubProtocol     *string                   `json:"TranscodingSubProtocol,omitempty"`
	TranscodingURL             *string                   `json:"TranscodingUrl,omitempty"`
	Type                       *MediaSourceType          `json:"Type,omitempty"`
	Video3DFormat              *Video3DFormat            `json:"Video3DFormat,omitempty"`
	WallClockStart             *string                   `json:"WallClockStart"`
}
