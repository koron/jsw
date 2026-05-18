# Jekyll Serve Watcher

[![PkgGoDev](https://pkg.go.dev/badge/github.com/koron/jsw)](https://pkg.go.dev/github.com/koron/jsw)
[![Actions/Go](https://github.com/koron/jsw/workflows/Go/badge.svg)](https://github.com/koron/jsw/actions?query=workflow%3AGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/koron/jsw)](https://goreportcard.com/report/github.com/koron/jsw)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/koron/jsw)

```console
$ jekyll serve --watch
```

`--watch` option waste many CPU time, because it is implemented by stat()
polling in each a second.

**jsw** replaces it by Go.

## Getting started

## How to compile

Install and upgrade.

```console
$ go install github.com/koron/jwc@latest
```

## Execute

Just type in your jekyll project:

```console
$ jsw
```

instead of:

```console
$ jekyll serve --watch
```

## Requirements

*   jekyll (1.0 or above)
