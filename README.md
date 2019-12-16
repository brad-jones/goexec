# goexec

[![GoReport](https://goreportcard.com/badge/brad-jones/goexec)](https://goreportcard.com/report/brad-jones/goexec)
[![GoLang](https://img.shields.io/badge/golang-%3E%3D%201.13.4-lightblue.svg)](https://golang.org)
[![GoDoc](https://godoc.org/github.com/brad-jones/goexec?status.svg)](https://godoc.org/github.com/brad-jones/goexec)
[![License](https://img.shields.io/github/license/brad-jones/goexec.svg)](https://github.com/brad-jones/goexec/blob/master/LICENSE)

Package goexec is a fluent decorator based API for os/exec.

## Usage

`go get -u github.com/brad-jones/goexec`

```go
package main

import (
	"github.com/brad-jones/goexec"
)

func main() {
	goexec.Run("ping", "-c", "4", "8.8.8.8")
}
```
