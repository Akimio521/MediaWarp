// UserItemDataDto.go

package schemas_emby

// UserItemDataDto
type UserItemDataDto struct {
	IsFavorite            *bool    `json:"IsFavorite,omitempty"`
	ItemID                *string  `json:"ItemId,omitempty"`
	Key                   *string  `json:"Key,omitempty"`
	LastPlayedDate        *string  `json:"LastPlayedDate"`
	PlaybackPositionTicks *int64   `json:"PlaybackPositionTicks,omitempty"`
	PlayCount             *int64   `json:"PlayCount"`
	Played                *bool    `json:"Played,omitempty"`
	PlayedPercentage      *float64 `json:"PlayedPercentage"`
	Rating                *float64 `json:"Rating"`
	ServerID              *string  `json:"ServerId,omitempty"`
	UnplayedItemCount     *int64   `json:"UnplayedItemCount"`
}
