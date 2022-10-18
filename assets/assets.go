//go:build !dev
// +build !dev

package assets

import (
	"embed"
	"io/fs"
)

//go:embed static
var embededFiles embed.FS

var Assets, _ = fs.Sub(embededFiles, "static")
