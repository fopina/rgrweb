//go:build tools

package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

var assetPath = "static/"
var cacheFile = "rgrweb.appcache"
var ignoredFiles = map[string]interface{}{
	cacheFile:              nil,
	"apple-touch-icon.png": nil,
	"manifest.json":        nil,
	"favicon.ico":          nil,
}

func main() {
	root := filepath.Dir(assetPath) + string(os.PathSeparator) + "."
	output := filepath.Dir(assetPath) + string(os.PathSeparator) + cacheFile
	log.Printf("Re-generating %s\n", output)

	hash := sha256.New()
	w, err := os.OpenFile(output, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteString("CACHE MANIFEST\n")

	walkFn := func(path string, d fs.DirEntry, err error) error {
		if d.Type().IsRegular() {
			if _, ok := ignoredFiles[d.Name()]; ok {
				return nil
			}

			rel, err := filepath.Rel(root, path)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(rel)
			c, err := os.OpenFile(path, os.O_RDONLY, 0)
			if err != nil {
				log.Fatal(err)
			}
			_, err = io.Copy(hash, c)
			if err != nil {
				log.Fatal(err)
			}
			w.WriteString(rel)
			w.WriteString("\n")
		}
		return nil
	}
	filepath.WalkDir(root, walkFn)

	sum := hash.Sum(nil)
	w.WriteString(fmt.Sprintf("\n# %x\n", sum))
	w.Close()
}
