//go:build !dev && !nomin

package assets

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed staticmin
var embededFiles embed.FS

var assetFS, _ = fs.Sub(embededFiles, "staticmin")

var Assets = http.FS(assetFS)
