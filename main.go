package main

import (
	"present/controller"
	"present/router"

	"present/middleware"

	"github.com/kataras/iris/v12"
)

func newApp() *iris.Application {
	app := iris.New()

	middleware.InitLog()

	app.Use(middleware.LoggerHandler)

	//for i:=0; i<common.TaskNum; i++{
	//	go func(chan controller.Task){
	//		controller.FetchRedPackage(controller.ChTaskList[i])
	//	}(controller.ChTaskList[i])
	//	fmt.Println(i)
	//}

	//fmt.Println(controller.ChTaskList)
	for _, v := range controller.ChTaskList {
		go func() {
			controller.FetchRedPackage(v)
		}()
	}

	router.InitRouter(app)

	return app
}

func main() {
	app := newApp()
	_ = app.Run(iris.Addr(":8080"))
}
