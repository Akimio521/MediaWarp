package constants

type ServerType string
type FliterMode string

const (
	EMBY     ServerType = "Emby"
	JELLYFIN ServerType = "Jellyfin"
	PLEX     ServerType = "Plex"
)

const (
	WHITELIST FliterMode = "WhiteList"
	BLACKLIST FliterMode = "BlackList"
)
