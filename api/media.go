package api

type MediaServer interface {
	GetType() string
	GetADDR() string
	GetToken() string
}
