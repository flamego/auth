# auth

[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/flamego/auth/Go?logo=github&style=for-the-badge)](https://github.com/flamego/auth/actions?query=workflow%3AGo)
[![Codecov](https://img.shields.io/codecov/c/gh/flamego/auth?logo=codecov&style=for-the-badge)](https://app.codecov.io/gh/flamego/auth)
[![GoDoc](https://img.shields.io/badge/GoDoc-Reference-blue?style=for-the-badge&logo=go)](https://pkg.go.dev/github.com/flamego/auth?tab=doc)
[![Sourcegraph](https://img.shields.io/badge/view%20on-Sourcegraph-brightgreen.svg?style=for-the-badge&logo=sourcegraph)](https://sourcegraph.com/github.com/flamego/auth)

Package auth is a middleware that provides request authentication for [Flamego](https://github.com/flamego/flamego).

## Installation

The minimum requirement of Go is **1.16**.

	go get github.com/flamego/auth

## Getting started

```go
package main

import (
	"github.com/flamego/auth"
	"github.com/flamego/flamego"
)

func main() {
	f := flamego.Classic()
	f.Use(auth.Basic("username", "secretpassword"))
	f.Get("/", func(user auth.User) string {
		return "Welcome, " + string(user)
	})
	f.Run()
}
```

## Getting help

- Please [file an issue](https://github.com/flamego/flamego/issues) or [start a discussion](https://github.com/flamego/flamego/discussions) on the [flamego/flamego](https://github.com/flamego/flamego) repository.

## License

This project is under the MIT License. See the [LICENSE](LICENSE) file for the full license text.
