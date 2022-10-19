module github.com/fopina/rgrweb

go 1.18

require (
	github.com/spf13/pflag v1.0.5
	github.com/tdewolff/minify/v2 v2.12.4
	github.com/warthog618/gpio v1.0.0
)

require (
	github.com/tdewolff/parse/v2 v2.6.4 // indirect
	golang.org/x/sys v0.1.0 // indirect
)

replace github.com/warthog618/gpio v1.0.0 => github.com/fopina/gpio v0.0.0-20221018093117-69def52fb60f
