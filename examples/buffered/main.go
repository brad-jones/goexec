package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/brad-jones/goerr/v2"
	"github.com/brad-jones/goexec/v2"
)

func main() {
	out, err := goexec.RunBuffered("ping", pingArg(), "4", "127.0.0.1")
	if err != nil {
		fmt.Println(out.StdErr)
		goerr.PrintTrace(err)
		os.Exit(1)
	}

	fmt.Println(strings.ReplaceAll(out.StdOut, "127.0.0.1", "1.2.3.4"))
}

func pingArg() string {
	if runtime.GOOS == "windows" {
		return "-n"
	}
	return "-c"
}
