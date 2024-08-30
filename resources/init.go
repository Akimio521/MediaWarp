package resources

import (
	"embed"
	"net/http"
)

//go:embed js/* css/*
var resources embed.FS

var ResourcesFS http.FileSystem = http.FS(resources)
