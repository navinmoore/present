package controller

import (
	"fmt"
	"math/rand"
	"present/common"
	"sync"
	"time"
)

type WechatShakeController interface {
	Get() (res common.Result)
	GetLuck() (res common.Result)
}

type wechatShakeController struct {
	BaseController
}

func NewWechatShakeController() WechatShakeController {
	return &wechatShakeController{}
}

const (
	giftTypeCoin      = iota // 虚拟币
	giftTypeCoupon           //劵
	giftTypeCouponFix        //相同的劵
	giftTypeRealSmall        // 小奖
	giftTypeRealLarge        //大将
)

type Gift struct {
	Id       int      `json:"id"`
	Name     string   `json:"name"`
	Pic      string   `json:"pic"`
	Link     string   `json:"link"`
	Gtype    int      `json:"gtype"`
	Data     string   `json:"data"`
	Datalist []string `json:"datalist"`
	Total    int      `json:"total"` // 总数 0 不限量
	Left     int      `json:"left"`
	Inuse    bool     `json:"inuse"` // 是否使用中
	Rate     int      `json:"rate"`  // 中奖概率 万分之N
	RateMin  int      `json:"rate_min"`
	RateMax  int      `json:"rate_max"`
}

const RateMax = 10000

type GiftList struct {
	giftList []*Gift
	sync.Mutex
}

var gitListStruct GiftList

func initGift() {

	gitListStruct.giftList = make([]*Gift, 3)
	g1 := Gift{
		Id:       1,
		Name:     "g1",
		Pic:      "",
		Link:     "",
		Gtype:    giftTypeRealSmall,
		Data:     "",
		Datalist: nil,
		Total:    5,
		Left:     5,
		Inuse:    true,
		Rate:     1000,
		RateMin:  0,
		RateMax:  0,
	}
	gitListStruct.giftList[0] = &g1
	g2 := Gift{
		Id:       2,
		Name:     "g2",
		Pic:      "",
		Link:     "",
		Gtype:    giftTypeCoupon,
		Data:     "",
		Datalist: nil,
		Total:    5,
		Left:     5,
		Inuse:    true,
		Rate:     3000,
		RateMin:  0,
		RateMax:  0,
	}
	gitListStruct.giftList[1] = &g2

	g3 := Gift{
		Id:       3,
		Name:     "g3",
		Pic:      "",
		Link:     "",
		Gtype:    giftTypeCoupon,
		Data:     "",
		Datalist: nil,
		Total:    5,
		Left:     5,
		Inuse:    true,
		Rate:     6000,
		RateMin:  0,
		RateMax:  0,
	}
	gitListStruct.giftList[2] = &g3

	//数据整理
	rateStart := 0
	for _, data := range gitListStruct.giftList {
		if !data.Inuse {
			continue
		}
		data.RateMin = rateStart
		data.RateMax = rateStart + data.Rate
		if data.RateMax > RateMax {
			data.RateMax = RateMax
			rateStart = 0
		} else {
			rateStart += data.Rate
		}
	}

}

func (w *wechatShakeController) Get() (res common.Result) {
	count := 0
	total := 0
	for _, data := range gitListStruct.giftList {
		if data.Inuse && (data.Total == 0 || (data.Total > 0 && data.Left > 0)) {
			count++
			total = total + data.Left
		}
	}
	m := make(map[string]int)
	m["count"] = count
	m["total"] = total
	res.Code = common.SuccessCode
	res.Msg = common.SuccessMsg
	res.Data = m
	fmt.Printf("%#v\n", gitListStruct.giftList[0])
	fmt.Printf("%#v\n", gitListStruct.giftList[1])
	fmt.Printf("%#v\n", gitListStruct.giftList[2])
	//fmt.Printf("%+v\n", giftList)
	//fmt.Printf("%v\n", giftList)
	return
}

func (w *wechatShakeController) GetLuck() (res common.Result) {
	code := luckCode()
	fmt.Println(code, "here")
	res.Data = fmt.Sprintln("没有中奖")

	gitListStruct.Lock()
	for _, data := range gitListStruct.giftList {
		if !data.Inuse || (data.Total > 0 && data.Left <= 0) {
			continue
		}
		if data.RateMin <= code && data.RateMax >= code {
			name, ok := sendGift(data)
			if ok {
				res.Data = fmt.Sprintf("恭喜你获得:%s", name)
			} else {
				res.Data = fmt.Sprintln("你来完了，奖品没了")
			}
		}
	}
	defer gitListStruct.Unlock()

	res.Code = common.SuccessCode
	res.Msg = common.SuccessMsg

	return
}

func luckCode() int {
	seed := time.Now().UnixNano()
	code := rand.New(rand.NewSource(seed)).Intn(RateMax)
	return code
}

func sendGift(gift *Gift) (string, bool) {
	if gift.Total == 0 {
		// 不限量的
		return gift.Name, true
	} else {
		if gift.Left > 0 {
			gift.Left = gift.Left - 1
			return gift.Name, true
		}
		return "", false
	}
}
