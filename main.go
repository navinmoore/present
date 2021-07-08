package main

import (
	"present/router"

	"present/middleware"

	"github.com/kataras/iris/v12"
)

func newApp() *iris.Application {
	app := iris.New()

	middleware.InitLog()

	app.Use(middleware.LoggerHandler)

	router.InitRouter(app)
	return app
}

func main() {
	app := newApp()
	_ = app.Run(iris.Addr(":8080"))
}
