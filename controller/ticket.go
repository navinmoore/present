package controller

import (
	"math/rand"
	"present/common"
	"time"
)

type TickController interface {
	//即开即得
	Get() (res common.Result)
	// 双色球
	GetPrice() (res common.Result)
}

type tickController struct {
	BaseController
}

func NewTickController() TickController {
	return &tickController{}
}

func (t *tickController) Get() (res common.Result) {
	seed := time.Now().UnixNano()
	code := rand.New(rand.NewSource(seed)).Intn(10)
	var prize string
	switch {
	case code == 1:
		prize = "一等奖"
	case code >= 2 && code <= 3:
		prize = "二等奖"
	case code >= 4 && code <= 6:
		prize = "三等奖"
	default:
		prize = "没中奖"
	}
	res.Code = common.SuccessCode
	res.Msg = common.SuccessMsg
	res.Data = prize
	return res
}

func (t *tickController) GetPrice() (res common.Result) {
	seed := time.Now().UnixNano
	r := rand.New(rand.NewSource(seed()))
	var prize [7]int
	for i := 0; i < 6; i++ {
		prize[i] = r.Intn(33) + 1
	}
	prize[6] = r.Intn(16) + 1
	res.Code = common.SuccessCode
	res.Msg = common.SuccessMsg
	res.Data = prize
	return res
}
