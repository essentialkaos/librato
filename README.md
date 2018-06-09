<p align="center"><a href="#readme"><img src="https://gh.kaos.st/mdtoc.svg"/></a></p>

<p align="center"><a href="#installation">Installation</a> • <a href="#examples">Examples</a> • <a href="#license">License</a></p>

<p align="center">
  <a href="https://godoc.org/pkg.re/essentialkaos/librato.v7"><img src="https://godoc.org/pkg.re/essentialkaos/librato.v7?status.svg"></a>
  <a href="https://travis-ci.org/essentialkaos/librato"><img src="https://travis-ci.org/essentialkaos/librato.svg"></a>
  <a href="https://goreportcard.com/report/github.com/essentialkaos/librato"><img src="https://goreportcard.com/badge/github.com/essentialkaos/librato"></a>
  <a href="https://codebeat.co/projects/github-com-essentialkaos-librato"><img alt="codebeat badge" src="https://codebeat.co/badges/f82e704d-67a7-4c6f-9e5d-1acf058c937b" /></a>
  <a href="https://essentialkaos.com/ekol"><img src="https://gh.kaos.st/ekol.svg"></a>
</p>

Package for working with [Librato Metrics](https://www.librato.com) API from Go code.

### Installation

Before the initial install allows git to use redirects for [pkg.re](https://github.com/essentialkaos/pkgre) service (reason why you should do this described [here](https://github.com/essentialkaos/pkgre#git-support)):

```
git config --global http.https://pkg.re.followRedirects true
```

Make sure you have a working Go 1.8+ workspace ([instructions](https://golang.org/doc/install)), then

```
go get pkg.re/essentialkaos/librato.v7
```

For update to latest stable release, do:

```
go get -u pkg.re/essentialkaos/librato.v7
```

### Examples

* [Basic Usage](examples/basic_example.go)
* [Metrics Collector](examples/collector_example.go)
* [Async Sending](examples/async_example.go)
* [Annotations](examples/annotations_example.go)

### License

[EKOL](https://essentialkaos.com/ekol)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
