//go:build dev && nomin

package assets

import "net/http"

// Assets contains project assets.
var Assets http.FileSystem = http.Dir("assets/static")
