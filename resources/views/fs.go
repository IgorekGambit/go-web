package views

import "embed"

//go:embed layouts/*.html pages/*.html
var Files embed.FS
