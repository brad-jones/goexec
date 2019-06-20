# goexec

[![GoReport](https://goreportcard.com/badge/brad-jones/goexec)](https://goreportcard.com/report/brad-jones/goexec)
[![GoLang](https://img.shields.io/badge/golang-%3E%3D%201.12.6-lightblue.svg)](https://golang.org)
[![GoDoc](https://godoc.org/github.com/brad-jones/goexec?status.svg)](https://godoc.org/github.com/brad-jones/goexec)
[![License](https://img.shields.io/github/license/brad-jones/goexec.svg)](https://github.com/brad-jones/goexec/blob/master/LICENSE)

Package goexec is a fluent decorator based API for os/exec.

## Usage

`go get -u github.com/brad-jones/goexec`

```go
package main

import (
	"github.com/brad-jones/goasync/await"
	"github.com/brad-jones/goexec"
)

func main() {
	done1, err1 := goexec.RunAsync("ping", "-c", "4", "8.8.8.8")
	done2, err2 := goexec.RunAsync("ping", "-c", "4", "8.8.4.4")
	select {
	case <-await.AllAsync(done1, done2):
	case e := <-await.AnyErrorAsync(err1, err2):
		panic(e)
	}
}
```
