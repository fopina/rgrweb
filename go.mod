module github.com/fopina/rgrweb

go 1.18

require (
	github.com/shurcooL/vfsgen v0.0.0-20200824052919-0d455de96546
	github.com/spf13/pflag v1.0.5
	github.com/warthog618/gpio v1.0.0
)

require (
	github.com/shurcooL/httpfs v0.0.0-20190707220628-8d4bc4ba7749 // indirect
	golang.org/x/sys v0.1.0 // indirect
	golang.org/x/tools v0.1.12 // indirect
)

replace github.com/warthog618/gpio v1.0.0 => github.com/fopina/gpio v0.0.0-20221017234359-903628097509
