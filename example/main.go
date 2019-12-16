package main

import (
	"github.com/brad-jones/goexec"
)

func main() {
	goexec.MustRunPrefixed("foo", "ping", "-c", "4", "1.1.1.1")
}
