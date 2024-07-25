atomicfs
========

A library for atomic filesystem operations in Go.

[![Go Reference](https://pkg.go.dev/badge/noxide.lol/go/atomicfs.svg)](https://pkg.go.dev/noxide.lol/go/atomicfs)
[![License](https://img.shields.io/github/license/noxideproject/atomicfs?color=7C00D8&style=flat-square&label=License)](https://github.com/noxideproject/atomicfs/blob/main/LICENSE)
[![Build](https://img.shields.io/github/actions/workflow/status/noxideproject/atomicfs/ci.yaml?style=flat-square&color=0FAA07&label=Tests)](https://github.com/noxideproject/atomicfs/actions/workflows/ci.yaml)

# Overview

The `noxide.lol/go/atomicfs` module provides a package for performing atomic
filesystem operations.

#### Reading material

- https://rcrowley.org/2010/01/06/things-unix-can-do-atomically.html

# Getting Started

The `atomicfs` package can be installed by running
```bash
$ go get noxide.lol/go/atomicfs
```

#### Examples

```go
writer := atomicfs.NewFileWriter(atomicfs.Options{
    TmpDirectory: "/tmp",
    Mode:         0600,
})

_ = writer.Write(input, output)
```

# Contributing

The `noxide.lol/go/atomicfs` module is always improving with new features and
error corrections. For contributing bug fixes and new features please file an
issue.

# License

The `noxide.lol/go/atomicfs` module is open source under the [BSD-3-Clause](LICENSE)
license.
