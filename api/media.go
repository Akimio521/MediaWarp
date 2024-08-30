package api

import "MediaWarp/constants"

type MediaServer interface {
	GetType() constants.ServerType
	GetADDR() string
	GetToken() string
}
