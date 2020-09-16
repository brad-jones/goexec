# goexec

[![PkgGoDev](https://pkg.go.dev/badge/github.com/brad-jones/goexec/v2)](https://pkg.go.dev/github.com/brad-jones/goexec/v2)
[![GoReport](https://goreportcard.com/badge/github.com/brad-jones/goexec/v2)](https://goreportcard.com/report/github.com/brad-jones/goexec/v2)
[![GoLang](https://img.shields.io/badge/golang-%3E%3D%201.15.1-lightblue.svg)](https://golang.org)
![.github/workflows/main.yml](https://github.com/brad-jones/goexec/workflows/.github/workflows/main.yml/badge.svg?branch=v2)
[![semantic-release](https://img.shields.io/badge/%20%20%F0%9F%93%A6%F0%9F%9A%80-semantic--release-e10079.svg)](https://github.com/semantic-release/semantic-release)
[![Conventional Commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg)](https://conventionalcommits.org)
[![KeepAChangelog](https://img.shields.io/badge/Keep%20A%20Changelog-1.0.0-%23E05735)](https://keepachangelog.com/)
[![License](https://img.shields.io/github/license/brad-jones/goexec.svg)](https://github.com/brad-jones/goexec/blob/v2/LICENSE)

Package goexec is a fluent decorator based API for os/exec.

_Looking for v1, see the [master branch](https://github.com/brad-jones/goexec/tree/master)_

## Quick Start

`go get -u github.com/brad-jones/goexec/v2`

```go
package main

import (
	"github.com/brad-jones/goexec/v2"
)

func main() {
	goexec.Run("ping", "-c", "4", "8.8.8.8")
}
```

_Also see further working examples under: <https://github.com/brad-jones/goexec/tree/v2/examples>_
