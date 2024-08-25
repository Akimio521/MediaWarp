// ItemResponse.go

package emby

// QueryResult_BaseItemDto
type ItemResponse struct {
	Items            []BaseItemDto `json:"Items,omitempty"`
	TotalRecordCount *int64        `json:"TotalRecordCount,omitempty"`
}
