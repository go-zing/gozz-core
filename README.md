<p align="center">
    <img width="200" src="https://raw.githubusercontent.com/go-zing/gozz-doc/main/docs/.vuepress/public/logo.png" alt="logo">
</p>

<div align=center>

[![Go Report Card](https://goreportcard.com/badge/github.com/go-zing/gozz-core)](https://goreportcard.com/report/github.com/go-zing/gozz-core)
[![Go Reference](https://pkg.go.dev/badge/github.com/go-zing/gozz-core.svg)](https://pkg.go.dev/github.com/go-zing/gozz-core)

[![License: MIT](https://img.shields.io/github/license/go-zing/gozz-core)](https://github.com/go-zing/gozz-core/blob/master/LICENSE)
[![Last Commit](https://img.shields.io/github/last-commit/go-zing/gozz-core)](https://github.com/go-zing/gozz-core/commits)
[![codecov](https://codecov.io/gh/go-zing/gozz-core/branch/main/graph/badge.svg)](https://codecov.io/gh/go-zing/gozz-core)

</div>

## Documentation

[English](https://go-zing.github.io/gozz) | [简体中文](https://go-zing.github.io/gozz/zh)

## Introduction

gozz-core provides core packages separated from [gozz](https://github.com/go-zing/gozz).
contains core typing and code-generate utils for better independent plugin package referenced.

### Why they independent

In [Golang plugin](https://pkg.go.dev/plugin),
module with same name loaded should be compiled in same version.

So we have great reason to reduce core dependencies of `gozz` in `go.mod` then
provide greater independence.

## License

[Apache-2.0](https://github.com/go-zing/gozz-core/blob/main/LICENSE)
