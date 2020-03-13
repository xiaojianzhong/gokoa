# GoKoa

![GitHub Workflow Status](https://img.shields.io/github/workflow/status/azxj/gokoa/Golang)
![GitHub repo size](https://img.shields.io/github/repo-size/azxj/gokoa)
![GitHub last commit](https://img.shields.io/github/last-commit/azxj/gokoa)
![GitHub](https://img.shields.io/github/license/azxj/gokoa)

[Koa][koa] styled web framework written in Go.

Read this in other languages: English | [简体中文](./README_zh-CN.md)

## Table of Contents

- [Introduction](#introduction)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [What Are NOT Implemented](#what-are-not-implemented)
- [Documentation](#documentation)
  - [`Application`](#application)
  - [`Context`](#context)
  - [`Request`](#request)
  - [`Response`](#response)
- [Development](#development)
  - [How to Test](#how-to-test)

## <a name="introduction"></a> Introduction

[GoKoa](gokoa) is a Koa styled web framework, which aims at reducing learning cost in grasping web development in Go for Node.js developers.

GoKoa is written in Go, which is beneficial to improving performance in handling HTTP requests.

## <a name="prerequisites"></a> Prerequisites

1. [Go](go) >= 1.13

## <a name="installation"></a> Installation

You can easily fetch the newest version of GoKoa by executing `go get`:

```bash
$ go get github.com/azxj/gokoa
```

## <a name="quick-start"></a> Quick Start

```go
package main

import (
	"github.com/azxj/gokoa"
)

func main() {
	app := gokoa.NewApplication(nil)

	app.Use(func(ctx *gokoa.Context, fn func() error) error {
		ctx.SetBody("hello gokoa")
		return nil
	})

	app.Listen(8080)
}
```

## <a name="what-are-not-implemented"></a> What Are NOT Implemented

1. lack of event emitting, except for `app.OnError()`
2. not able to customize HTTP reason phrase (don't do that cause it breaks the best practice)
3. not able to bypass GoKoa's response handling (actually it is deprecated by Koa, too)
4. not able to access the socket related to a HTTP connection
5. lack of `headerSent` property

## <a name="documentation"></a> Documentation

### <a name="application"></a> `Application`

### <a name="context"></a> `Context`

### <a name="request"></a> `Request`

### <a name="response"></a> `Response`

## <a name="development"></a> Development

### <a name="how-to-test"></a> How To Test

You can easily run unit tests by executing `go test`:

```bash
$ go test github.com/azxj/gokoa
```

[koa]: https://github.com/koajs/koa
[go]: https://golang.org/
[gokoa]: https://github.com/azxj/gokoa
