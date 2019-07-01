atomicfs
========

A library for atomic filesystem operations in Go.

[![Go Report Card](https://goreportcard.com/badge/go.gophers.dev/pkgs/atomicfs)](https://goreportcard.com/report/go.gophers.dev/pkgs/atomicfs)
[![Build Status](https://travis-ci.com/shoenig/atomicfs.svg?branch=master)](https://travis-ci.com/shoenig/atomicfs)
[![GoDoc](https://godoc.org/go.gophers.dev/pkgs/atomicfs?status.svg)](https://godoc.org/go.gophers.dev/pkgs/atomicfs)
[![NetflixOSS Lifecycle](https://img.shields.io/osslifecycle/shoenig/atomicfs.svg)](OSSMETADATA)
[![GitHub](https://img.shields.io/github/license/shoenig/atomicfs.svg)](LICENSE)

# Project Overview

The `go.gophers.dev/pkgs/atomicfs` module provides a package for performing atomic
filesystem operations.

#### Reading material
- https://rcrowley.org/2010/01/06/things-unix-can-do-atomically.html

# Getting Started

The `atomicfs` package can be installed by running
```bash
$ go get go.gophers.dev/pkgs/atomicfs
```

#### Example usage
```golang
writer := atomicfs.NewFileWriter(atomicfs.Options{
    TmpDirectory: "/tmp",
    Mode:         0600,
})

_ = writer.Write(input, output)
```

# Contributing

The `go.gophers.dev/pkgs/atomicfs` module is always improving with new features
and error corrections. For contributing bug fixes and new features please file an issue.

# License

The `go.gophers.dev/pkgs/atomicfs` module is open source under the [BSD-3-Clause](LICENSE) license.
