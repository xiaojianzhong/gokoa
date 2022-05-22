/*
Package gokoa is a Koa styled web framework written in Go.

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
*/
package gokoa
