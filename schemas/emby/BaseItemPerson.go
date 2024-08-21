// BaseItemPerson.go

package emby

// BaseItemPerson
type BaseItemPerson struct {
	ID              *string     `json:"Id,omitempty"`
	Name            *string     `json:"Name,omitempty"`
	PrimaryImageTag *string     `json:"PrimaryImageTag,omitempty"`
	Role            *string     `json:"Role,omitempty"`
	Type            *PersonType `json:"Type,omitempty"`
}
