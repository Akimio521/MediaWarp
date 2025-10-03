package constants

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type MediaServerType uint8 // 媒体服务器类型

const (
	EMBY     MediaServerType = iota // 媒体服务器类型：EmbyServer
	JELLYFIN                        // 媒体服务器类型：Jellyfin
	PLEX                            // 媒体服务器类型：Plex
)

func (m *MediaServerType) UnMarshalJSON(data []byte) error {
	switch string(data) {
	case `"Emby"`:
		*m = EMBY
	case `"Jellyfin"`:
		*m = JELLYFIN
	case `"Plex"`:
		*m = PLEX
	default:
		return fmt.Errorf("invalid MediaServerType: %s", string(data))
	}
	return nil
}

func (m MediaServerType) String() string {
	switch m {
	case EMBY:
		return "Emby"
	case JELLYFIN:
		return "Jellyfin"
	case PLEX:
		return "Plex"
	default:
		return "Unknown"
	}
}

func (m *MediaServerType) UnmarshalYAML(value *yaml.Node) error {
	var s string
	if err := value.Decode(&s); err != nil {
		return err
	}
	switch s {
	case "Emby":
		*m = EMBY
	case "Jellyfin":
		*m = JELLYFIN
	case "Plex":
		*m = PLEX
	default:
		return fmt.Errorf("invalid MediaServerType: %s", s)
	}
	return nil
}
