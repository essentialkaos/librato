# Librato [![GoDoc](https://godoc.org/pkg.re/essentialkaos/librato.v7?status.svg)](https://godoc.org/pkg.re/essentialkaos/librato.v7) [![Build Status](https://travis-ci.org/essentialkaos/librato.svg?branch=master)](https://travis-ci.org/essentialkaos/librato) [![Go Report Card](https://goreportcard.com/badge/github.com/essentialkaos/librato)](https://goreportcard.com/report/github.com/essentialkaos/librato) [![codebeat badge](https://codebeat.co/badges/f82e704d-67a7-4c6f-9e5d-1acf058c937b)](https://codebeat.co/projects/github-com-essentialkaos-librato) [![License](https://gh.kaos.io/ekol.svg)](https://essentialkaos.com/ekol)

Package for working with [Librato Metrics](https://www.librato.com) API from Go code.

## Installation

Before the initial install allows git to use redirects for [pkg.re](https://github.com/essentialkaos/pkgre) service (reason why you should do this described [here](https://github.com/essentialkaos/pkgre#git-support)):

```
git config --global http.https://pkg.re.followRedirects true
```

Make sure you have a working Go 1.6+ workspace ([instructions](https://golang.org/doc/install)), then

```
go get pkg.re/essentialkaos/librato.v7
```

For update to latest stable release, do:

```
go get -u pkg.re/essentialkaos/librato.v7
```

## Examples

* [Basic Usage](examples/basic_example.go)
* [Metrics Collector](examples/collector_example.go)
* [Async Sending](examples/async_example.go)
* [Annotations](examples/annotations_example.go)

## License

[EKOL](https://essentialkaos.com/ekol)
