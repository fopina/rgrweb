//go:build !dev
// +build !dev

package assets

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed static
var embededFiles embed.FS

var assetFS, _ = fs.Sub(embededFiles, "static")

var Assets = http.FS(assetFS)
