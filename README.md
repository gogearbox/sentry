<p align="center">
	<a href="https://gogearbox.com">
    	<img src="https://raw.githubusercontent.com/gogearbox/gearbox/master/assets/gearbox-512.png"/>
	</a>
    <br />
    <a href="https://godoc.org/github.com/gogearbox/netadaptor">
      <img src="https://godoc.org/github.com/gogearbox/netadaptor?status.png" />
    </a>
    <img src="https://github.com/gogearbox/netadaptor/workflows/Test%20&%20Build/badge.svg?branch=master" />
    <a href="https://goreportcard.com/report/github.com/gogearbox/gearbox">
      <img src="https://goreportcard.com/badge/github.com/gogearbox/netadaptor" />
    </a>
	<a href="https://discord.com/invite/CT8my4R">
      <img src="https://img.shields.io/discord/716724372642988064?label=Discord&logo=discord">
  	</a>
    <a href="https://deepsource.io/gh/gogearbox/netadaptor/?ref=repository-badge" target="_blank">
      <img alt="DeepSource" title="DeepSource" src="https://static.deepsource.io/deepsource-badge-light-mini.svg">
    </a>
</p>

**sentry** middleware


### Supported Go versions & installation

:gear: gearbox requires version `1.11` or higher of Go ([Download Go](https://golang.org/dl/))

Just use [go get](https://golang.org/cmd/go/#hdr-Add_dependencies_to_current_module_and_install_them) to download and install gearbox

```bash
go get -u github.com/gogearbox/gearbox
go get u- github.com/gogearbox/sentry
```


### Examples

```go
package main

import (
	  "github.com/getsentry/sentry-go"
	  "github.com/gogearbox/gearbox"
    sentrymiddleware "github.com/gogearbox/sentry"
)

func main() {
	// Setup gearbox
	gb := gearbox.New()

	// Initialize sentry
	_ = sentry.Init(sentry.ClientOptions{
		Dsn:              PROJECT_DSN,
	})

	// Register the sentry middleware for all requests
	gb.Use(sentrymiddleware.New())

	// Define your handler
	gb.Post("/hello", func(ctx gearbox.Context) {
		panic("There is an issue")
	})

	// Start service
	gb.Start(":3000")
}

```
