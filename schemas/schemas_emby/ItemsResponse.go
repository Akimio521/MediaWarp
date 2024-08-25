// ItemResponse.go

package schemas_emby

// QueryResult_BaseItemDto
type ItemResponse struct {
	Items            []BaseItemDto `json:"Items,omitempty"`
	TotalRecordCount *int64        `json:"TotalRecordCount,omitempty"`
}
