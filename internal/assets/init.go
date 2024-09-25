package assets

import (
	"embed"
)

//go:embed js/* css/*
var EmbeddedStaticAssets embed.FS
