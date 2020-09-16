package main

import (
	"runtime"

	"github.com/brad-jones/goasync/v2/await"
	"github.com/brad-jones/goerr/v2"
	"github.com/brad-jones/goexec/v2"
)

func main() {
	defer goerr.Handle(func(err error) {
		goerr.PrintTrace(err)
	})

	await.MustFastAllOrError(
		goexec.RunPrefixedAsync("ip1", "ping", pingArg(), "4", "127.0.0.1"),
		goexec.RunPrefixedAsync("ip2", "ping", pingArg(), "4", "127.0.0.2"),
	)
}

func pingArg() string {
	if runtime.GOOS == "windows" {
		return "-n"
	}
	return "-c"
}
