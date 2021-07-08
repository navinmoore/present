package controller

import (
	"math/rand"
	"present/common"
	"time"
)

type AnnualUserController interface {
	Get() (res common.Result)
	PostImport() (res common.Result)
	GetLuck() (res common.Result)
}

type annualUserController struct {
	BaseController
}

func NewAnnualUserController() AnnualUserController {
	return &annualUserController{}
}

func (u *annualUserController) Get() (res common.Result) {
	res.Code = common.SuccessCode
	res.Msg = common.SuccessMsg
	res.Data = annualUserList.len
	return
}

type ImportRequest struct {
	Name string `json:"name" validate:"required"`
}

func (u *annualUserController) PostImport() common.Result {
	var importArgs ImportRequest
	err := u.Ctx.ReadJSON(&importArgs)
	if err != nil {
		return common.ParseParamsErrorResult
	}

	err = validate.Struct(&importArgs)
	if err != nil {
		return common.ParseParamsErrorResult
	}

	annualUserList.Append(importArgs.Name)
	var res common.Result
	res.Code = common.SuccessCode
	res.Msg = common.SuccessMsg
	// res.Data = userList.userList
	res.Data = nil
	return res
}

func (u *annualUserController) GetLuck() (res common.Result) {
	if annualUserList.IsEmpty() {
		res.Code = common.UserLstEmptyErrorCode
		res.Msg = common.UserLstEmptyErrorMsg
		res.Data = nil
		return
	}

	rand.Seed(time.Now().Unix())
	index := rand.Intn(annualUserList.len)
	name := annualUserList.GetValue(index)
	annualUserList.Remove(index)
	res.Code = common.SuccessCode
	res.Msg = common.SuccessMsg
	m := make(map[string]string)
	m["name"] = name
	res.Data = m
	return
}
