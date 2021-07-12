package controller

import (
	"math/rand"
	"present/common"
	"present/middleware"
	"strconv"
	"sync"
	"time"
)

// 抽奖前用户已知全部奖品信息
// 后端设置各个奖品的中奖概率和数量限制
// 更新奖品库存的时候存在并发安全性问题

type Prate struct {
	Rate      int // 万分之N的中奖概率
	Total     int // 总数量限制
	CodeStart int // 中奖概率开始编码
	CodeEnd   int //中奖概率结束编码
	Left      int //剩余奖品数量
	sync.Mutex
}

func (p *Prate) Sub() {
	p.Lock()
	if p.Left > 0 {
		p.Left -= 1
	}
	p.Unlock()
}

func (p *Prate) GetLeft() int {
	p.Lock()
	a := p.Left
	p.Unlock()
	return a
}

var prizeList []string = []string{
	"一等奖, 港澳游",
	"二等奖，MBP",
	"三等奖, Iphone12",
	"谢谢",
}

var rateList []Prate = []Prate{
	{Rate: 1, Total: 1, CodeStart: 0, CodeEnd: 100, Left: 1},
	{Rate: 2, Total: 2, CodeStart: 100, CodeEnd: 300, Left: 2},
	{Rate: 5, Total: 10, CodeStart: 300, CodeEnd: 800, Left: 10},
	{Rate: 100, Total: 0, CodeStart: 0, CodeEnd: 9999, Left: 0},
}

type WheelController interface {
	//
	Get() (res common.Result)
	GetDebug() (res common.Result)
}

type wheelController struct {
	BaseController
}

func NewWheelController() WheelController {
	return &wheelController{}
}

func (w *wheelController) Get() (res common.Result) {
	//w.Ctx.Header("Content-Type", "text/html")
	result := make([]string, len(prizeList))
	for i, value := range prizeList {
		result[i] = value + "剩余:" + strconv.Itoa(rateList[i].Left)
	}
	res.Code = common.SuccessCode
	res.Msg = common.SuccessMsg
	res.Data = result
	return
}

func (w *wheelController) GetDebug() (res common.Result) {
	result := make([]string, len(prizeList))
	for i, value := range prizeList {
		result[i] = value + "剩余:" + strconv.Itoa(rateList[i].Left) + "概率:" + strconv.Itoa(rateList[i].Rate)
	}
	res.Code = common.SuccessCode
	res.Msg = common.SuccessMsg
	res.Data = result
	return
}

func (w *wheelController) GetLuck() (res common.Result) {
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	code := r.Intn(10000)
	var myPrize string
	var prizeRate *Prate
	//从奖品列表匹配是否中奖
	for i, prize := range prizeList {
		rate := &rateList[i]
		if code >= rate.CodeStart && code <= rate.CodeEnd {
			myPrize = prize
			prizeRate = rate
			break
		}
	}
	if myPrize == "谢谢" {
		myPrize = "谢谢"
		res.Data = myPrize
		res.Code = common.SuccessCode
		res.Msg = common.SuccessMsg
		return
	}
	if prizeRate.GetLeft() == 0 {
		res.Data = myPrize
		res.Code = common.SuccessCode
		res.Msg = common.SuccessMsg
		return
	} else if prizeRate.GetLeft() > 0 {
		//prizeRate.Left -= 1
		prizeRate.Sub()
		res.Data = myPrize
		res.Code = common.SuccessCode
		res.Msg = common.SuccessMsg
		middleware.Log.Infof("恭喜你，%s", myPrize)
		return
	} else {
		myPrize = "谢谢"
		res.Data = myPrize
		res.Code = common.SuccessCode
		res.Msg = common.SuccessMsg
		return
	}

}
