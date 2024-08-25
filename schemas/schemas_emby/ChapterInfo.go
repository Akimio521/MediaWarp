// ChapterInfo.go

package schemas_emby

// ChapterInfo
type ChapterInfo struct {
	ChapterIndex       *int64      `json:"ChapterIndex,omitempty"`
	ImageTag           *string     `json:"ImageTag,omitempty"`
	MarkerType         *MarkerType `json:"MarkerType,omitempty"`
	Name               *string     `json:"Name,omitempty"`
	StartPositionTicks *int64      `json:"StartPositionTicks,omitempty"`
}
