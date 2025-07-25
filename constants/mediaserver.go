package constants

import (
	"fmt"
	"strings"
)

type MediaServerType uint8 // 媒体服务器类型

const (
	EMBY     MediaServerType = iota // 媒体服务器类型：EmbyServer
	JELLYFIN                        // 媒体服务器类型：Jellyfin
	PLEX                            // 媒体服务器类型：Plex
)

func (m *MediaServerType) UnMarshalJSON(data []byte) error {
	fmt.Println("UnMarshalJSON", string(data))
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

func (m *MediaServerType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	fmt.Println("UnMarshalYAML", s)
	switch strings.ToLower(s) {
	case "emby":
		*m = EMBY
	case "jellyfin":
		*m = JELLYFIN
	case "plex":
		*m = PLEX
	default:
		return fmt.Errorf("invalid MediaServerType: %s", s)
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
