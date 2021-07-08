package controller

import (
	"math/rand"
	"present/common"
	"time"
)

type UserController interface {
	Get() (res common.Result)
	PostImport() (res common.Result)
	GetLuck() (res common.Result)
}

type userController struct {
	BaseController
}

func NewUserController() UserController {
	return &userController{}
}

func (u *userController) Get() (res common.Result) {
	res.Code = common.SuccessCode
	res.Msg = common.SuccessMsg
	res.Data = nil
	return
}

type ImportRequest struct {
	Name string `json:"name" validate:"required"`
}

func (u *userController) PostImport() common.Result {
	var importArgs ImportRequest
	err := u.Ctx.ReadJSON(&importArgs)
	if err != nil {
		return common.ParseParamsErrorResult
	}

	err = validate.Struct(&importArgs)
	if err != nil {
		return common.ParseParamsErrorResult
	}

	userList.Append(importArgs.Name)
	var res common.Result
	res.Code = common.SuccessCode
	res.Msg = common.SuccessMsg
	// res.Data = userList.userList
	res.Data = nil
	return res
}

func (u *userController) GetLuck() (res common.Result) {
	if userList.IsEmpty() {
		res.Code = common.UserLstEmptyErrorCode
		res.Msg = common.UserLstEmptyErrorMsg
		res.Data = nil
		return
	}

	rand.Seed(time.Now().Unix())
	index := rand.Intn(userList.len)
	name := userList.GetValue(index)
	userList.Remove(index)
	res.Code = common.SuccessCode
	res.Msg = common.SuccessMsg
	m := make(map[string]string)
	m["name"] = name
	res.Data = m
	return
}
