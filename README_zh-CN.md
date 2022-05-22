# GoKoa

![GitHub Workflow Status](https://img.shields.io/github/workflow/status/xiaojianzhong/gokoa/Golang)
![GitHub repo size](https://img.shields.io/github/repo-size/xiaojianzhong/gokoa)
![GitHub last commit](https://img.shields.io/github/last-commit/xiaojianzhong/gokoa)
![GitHub](https://img.shields.io/github/license/xiaojianzhong/gokoa)

[Koa](koa) 风格的，基于 Go 的 web 框架。

阅读其他语言的版本：[English](./README.md) | 简体中文

## 目录

- [介绍](#introduction)
- [环境要求](#prerequisites)
- [安装](#installation)
- [快速上手](#quick-start)
- [哪些功能没有被实现？](#what-are-not-implemented)
- [文档](#documentation)
  - [`Application`](#application)
  - [`Context`](#context)
  - [`Request`](#request)
  - [`Response`](#response)
- [有关开发](#development)
  - [如何执行测试？](#how-to-test)

## <a name="introduction"></a> 介绍

[GoKoa](gokoa) 是一个 Koa 风格的 web 框架，致力于减少 Node.js 开发者迁移到 Go web 开发过程中的学习成本。

GoKoa 用 Go 编写，这有利于提高处理 HTTP 请求的性能。

## <a name="prerequisites"></a> 环境要求

1. [Go](go) >= 1.13

## <a name="installation"></a> 安装

通过执行 `go get`，你可以轻松获得 GoKoa 的最新版本：

```bash
$ go get github.com/xiaojianzhong/gokoa
```

## <a name="quick-start"></a> 快速上手

```go
package main

import (
	"github.com/xiaojianzhong/gokoa"
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

## <a name="what-are-not-implemented"></a> 哪些功能没有被实现？

1. 缺少事件机制（`app.OnError()` 除外）
2. 无法自定义 HTTP 响应状态描述信息（这是违背最佳实践的行为）
3. 无法绕过 GoKoa 的响应处理器（事实上，Koa 也已经废弃了这个功能）
4. 无法访问与 HTTP 连接挂钩的 socket
5. 缺少 `headerSent` 属性

## <a name="documentation"></a> 文档

### <a name="application"></a> `Application`

### <a name="context"></a> `Context`

### <a name="request"></a> `Request`

### <a name="response"></a> `Response`

## <a name="development"></a> 有关开发

### <a name="how-to-test"></a> 如何执行测试？

通过执行 `go test`，你可以轻松运行 GoKoa 的所有单元测试：

```bash
$ go test github.com/xiaojianzhong/gokoa
```

[koa]: https://github.com/koajs/koa
[go]: https://golang.org/
[gokoa]: https://github.com/xiaojianzhong/gokoa

