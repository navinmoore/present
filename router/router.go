package router

import (
	"present/controller"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func InitRouter(app *iris.Application) {
	bathPathV1 := "/api/v1"
	mvc.Configure(app.Party(bathPathV1+"/annual"), func(m *mvc.Application) {
		// m.Router.Use(middleware.JwtHandler().Serve)
		m.Handle(controller.NewAnnualUserController())
	})
	//NewTickController
	mvc.Configure(app.Party(bathPathV1+"/ticket"), func(m *mvc.Application) {
		m.Handle(controller.NewTickController())
	})

	//wechat
	mvc.Configure(app.Party(bathPathV1+"/wechat"), func(m *mvc.Application) {
		m.Handle(controller.NewWechatShakeController())
	})

	// weibo
	mvc.Configure(app.Party(bathPathV1+"/weibo"), func(m *mvc.Application) {
		m.Handle(controller.NewWeiboRedPackageController())
	})
	//wheel
	mvc.Configure(app.Party(bathPathV1+"/wheel"), func(m *mvc.Application) {
		m.Handle(controller.NewWheelController())
	})
}
