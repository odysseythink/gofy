package main

import (
	"embed"
	"io/fs"
)

//go:embed all:frontend/dist
var frontendFS embed.FS

// FrontendAssets returns the frontend dist directory as an fs.FS
func FrontendAssets() (fs.FS, error) {
	return fs.Sub(frontendFS, "frontend/dist")
}
