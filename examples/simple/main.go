package main

import (
	"os"
	"runtime"

	"github.com/brad-jones/goerr/v2"
	"github.com/brad-jones/goexec/v2"
)

func main() {
	if err := goexec.Run("ping", pingArg(), "4", "127.0.0.1"); err != nil {
		goerr.PrintTrace(err)
		os.Exit(1)
	}
}

func pingArg() string {
	if runtime.GOOS == "windows" {
		return "-n"
	}
	return "-c"
}
