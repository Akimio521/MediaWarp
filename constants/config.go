package constants

const (
	EMBY     string = "Emby"
	JELLYFIN string = "Jellyfin"
	PLEX     string = "Plex"
)

type FliterMode string

const (
	WHITELIST FliterMode = "WhiteList"
	BLACKLIST FliterMode = "BlackList"
)
