package constants

type FliterMode uint8

const (
	WHITELIST FliterMode = iota
	BLACKLIST
)

func (f FliterMode) String() string {
	switch f {
	case WHITELIST:
		return "WhiteList"
	case BLACKLIST:
		return "BlackList"
	default:
		return "Unknown"
	}
}
