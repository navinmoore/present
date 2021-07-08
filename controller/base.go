package controller

import (
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
)

type BaseController struct {
	Ctx iris.Context
}

var (
	validate *validator.Validate
	userList UserList
)

type UserList struct {
	userList []string
	sync.Mutex
	len int
}

func (u *UserList) Append(name string) {
	u.Lock()
	defer u.Unlock()
	u.userList = append(u.userList, name)
	u.len = u.len + 1
}

func (u *UserList) Remove(index int) bool {
	u.Lock()
	defer u.Unlock()
	if u.len == 0 {
		return false
	}
	if index == 0 {
		u.userList = u.userList[1:]
	} else if index == u.len {
		u.userList = u.userList[:u.len-1]
	} else {
		u.userList = append(u.userList[:index], u.userList[index+1:]...)
	}
	u.len = u.len - 1
	return true
}

func (u *UserList) IsEmpty() bool {
	return u.len == 0
}

func (u *UserList) GetValue(index int) string {
	return u.userList[index]
}

func init() {
	validate = validator.New()

	// validate.RegisterValidation("gender", func(fl validator.FieldLevel) bool {
	// 	// 获取 Field 的值
	// 	gender := fl.Field().String()
	// 	if gender != "girl" && gender != "boy" {
	// 		return false
	// 	}
	// 	return true
	// })
}
