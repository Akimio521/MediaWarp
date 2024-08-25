// UsersItemResponse.go

package schemas_emby

// /Users/:userID/Items的响应
type UserItemsResponse struct {
	Items            []BaseItemDto `json:"Items,omitempty"`
	TotalRecordCount *int64        `json:"TotalRecordCount,omitempty"`
}
