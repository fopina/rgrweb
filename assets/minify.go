//go:build tools

package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"

	min "github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/json"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/tdewolff/minify/v2/xml"
)

var assetPath = "static/"
var minifiedPath = "staticmin/"

var filetypeMime = map[string]string{
	"css":  "text/css",
	"htm":  "text/html",
	"html": "text/html",
	"js":   "application/javascript",
	"mjs":  "application/javascript",
	"json": "application/json",
	"svg":  "image/svg+xml",
	"xml":  "text/xml",
}

func openInputFile(input string) (io.ReadCloser, error) {
	r, err := os.Open(input)
	if err != nil {
		return nil, fmt.Errorf("open input file %q: %w", input, err)
	}
	return r, nil
}

func openOutputFile(output string) (*os.File, error) {
	dir := filepath.Dir(output)
	if err := os.MkdirAll(dir, 0777); err != nil {
		return nil, fmt.Errorf("creating directory %q: %w", dir, err)
	}

	w, err := os.OpenFile(output, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return nil, fmt.Errorf("open output file %q: %w", output, err)
	}
	return w, nil
}

func minifyAssets() bool {
	output := minifiedPath
	root := filepath.Dir(assetPath) + string(os.PathSeparator) + "."

	cssMinifier := &css.Minifier{}
	htmlMinifier := &html.Minifier{}
	jsMinifier := &js.Minifier{}
	jsonMinifier := &json.Minifier{}
	svgMinifier := &svg.Minifier{}
	xmlMinifier := &xml.Minifier{}
	m := min.New()
	m.Add("text/css", cssMinifier)
	m.Add("text/html", htmlMinifier)
	m.Add("image/svg+xml", svgMinifier)
	m.AddRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma|j|live)script(1\\.[0-5])?$|^module$"), jsMinifier)
	m.AddRegexp(regexp.MustCompile("[/+]json$"), jsonMinifier)
	m.AddRegexp(regexp.MustCompile("[/+]xml$"), xmlMinifier)

	success := true

	walkFn := func(path string, d fs.DirEntry, err error) error {
		if d.Type().IsRegular() {
			rel, err := filepath.Rel(root, path)
			if err != nil {
				log.Fatal(err)
			}
			outputPath := filepath.Join(output, rel)
			r, err := openInputFile(path)
			if err != nil {
				log.Fatal(err)
			}
			w, err := openOutputFile(outputPath)
			if err != nil {
				log.Fatal(err)
			}
			ext := filepath.Ext(path)
			if 0 < len(ext) {
				ext = ext[1:]
			}
			srcMimetype, ok := filetypeMime[ext]
			if ok {
				if err = m.Minify(srcMimetype, w, r); err != nil {
					log.Fatalf("cannot minify %s: %v\n", path, err)
				} else {
					info1, err := os.Stat(path)
					if err != nil {
						log.Fatal(err)
					}
					info2, err := os.Stat(outputPath)
					if err != nil {
						log.Fatal(err)
					}
					rate := float64(info2.Size()) * 100 / float64(info1.Size())
					log.Printf("Minified %v to %v: %.2f%% (%d > %d) \n", path, outputPath, rate, info1.Size(), info2.Size())
				}
			} else {
				io.Copy(w, r)
				log.Printf("Copied %v to %v\n", path, outputPath)
			}
			r.Close()
			w.Close()
		}
		return nil
	}
	os.RemoveAll(output)
	filepath.WalkDir(root, walkFn)

	return success
}

func main() {
	if !minifyAssets() {
		os.Exit(1)
	}
}
