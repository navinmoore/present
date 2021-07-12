package controller

import (
	"fmt"
	"math/rand"
	"present/common"
	"present/middleware"
	"strconv"
	"sync"
	"time"
)

type WeiboRedPackageController interface {
	Get() (res common.Result)
	Post() (res common.Result)
	GetRecv() (res common.Result)
}

type weiboRedPackageController struct {
	BaseController
}

func NewWeiboRedPackageController() WeiboRedPackageController {
	return &weiboRedPackageController{}
}

type Task struct {
	id       uint32
	callback chan uint
}

// sync.Map 线程安全
//var packageList *sync.Map = new(sync.Map)
var packageLists [common.TaskNum]*sync.Map

//var chTasks  = make(chan task)
var ChTaskList []chan Task = make([]chan Task, common.TaskNum)

func init() {
	for i := 0; i < common.TaskNum; i++ {
		packageLists[i] = new(sync.Map)
		ChTaskList[i] = make(chan Task)
	}
}

// 全部的红包地址
func (c *weiboRedPackageController) Get() (res common.Result) {

	rs := make(map[uint32][2]int)

	for i := 0; i < common.TaskNum; i++ {
		packageLists[i].Range(func(key, value interface{}) bool {
			id := key.(uint32)
			list := value.([]uint)
			var money int
			for _, v := range list {
				money += int(v)
			}
			rs[id] = [2]int{len(list), money}
			return true
		})
	}

	res.Code = common.SuccessCode
	res.Msg = common.SuccessMsg
	res.Data = rs
	return
}

type WeiBoRedPackageParams struct {
	Uid   int `json:"uid" validate:"required"`
	Money int `json:"money" validate:"required"` // 总共有多少钱 分为单位
	Num   int `json:"num" validate:"required"`   // 多少个红包
}

// 设置红包
func (c *weiboRedPackageController) Post() (res common.Result) {
	var weiBoRedPackageParams WeiBoRedPackageParams
	_ = c.Ctx.ReadJSON(&weiBoRedPackageParams)

	uid := weiBoRedPackageParams.Uid
	money := weiBoRedPackageParams.Money
	num := weiBoRedPackageParams.Num

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	rMax := 0.55 //随机分配最大值
	list := make([]uint, num)
	leftMoney := money
	leftNum := num

	for leftNum > 0 {
		if leftNum == 1 {
			list[num-1] = uint(leftMoney)
			break
		}

		if leftMoney == leftNum {
			for i := num - leftNum; i < num; i++ {
				list[i] = 1
			}
			break
		}

		rMoney := int(float64(leftMoney-leftNum) * rMax)
		m := r.Intn(rMoney)
		if m < 1 {
			m = 1
		}
		list[num-leftNum] = uint(m)
		leftMoney = leftMoney - m
		leftNum--
	}

	id := r.Uint32()
	//packageList.packageList[id] = list
	packageLists[id%common.TaskNum].Store(id, list)
	fmt.Println(packageLists[id%common.TaskNum].Load(id))
	res.Code = common.SuccessCode
	res.Msg = common.SuccessMsg
	res.Data = fmt.Sprintf("/api/v1/weibo/recv?id=%d&uid=%d&num=%d", id, uid, num)
	return res
}

func (c *weiboRedPackageController) GetRecv() (res common.Result) {
	//id := c.Ctx.URLParamInt32Default("id", 0)

	var qid string
	qid = c.Ctx.Request().URL.Query().Get("id")
	var id int
	id, _ = strconv.Atoi(qid)
	if id == 0 {
		res = common.ParseParamsErrorResult
		return
	}
	rlist, ok := packageLists[id%common.TaskNum].Load(uint32(id))
	if ok && rlist != nil {
		// 构造一个task
		callback := make(chan uint)
		t := Task{id: uint32(id), callback: callback}
		fmt.Println(ChTaskList, "here", len(ChTaskList), id%common.TaskNum)
		//发送
		ChTaskList[id%common.TaskNum] <- t
		fmt.Println(ChTaskList, "here##")
		//等待回调
		data := <-t.callback
		fmt.Println(ChTaskList, "here", len(ChTaskList))

		if data <= 0 {
			res.Code = common.SuccessCode
			res.Msg = common.SuccessMsg
			res.Data = fmt.Sprintln("没有红包了")
		} else {
			res.Code = common.SuccessCode
			res.Msg = common.SuccessMsg
			res.Data = fmt.Sprintf("恭喜你，抢得%.2f", float64(data)*0.01)
			middleware.Log.Infof("恭喜你，抢得%.2f", float64(data)*0.01)
		}
	} else {
		res.Code = common.SuccessCode
		res.Msg = common.SuccessMsg
		res.Data = fmt.Sprintln("没有红包了")
	}
	return
}

func FetchRedPackage(chTasks chan Task) {
	for {
		fmt.Println("FetchRedPackage")
		t := <-chTasks
		fmt.Println("FetchRedPackage ttt")
		fmt.Println("haha")
		id := t.id
		rlist, ok := packageLists[id%common.TaskNum].Load(uint32(id))
		list := rlist.([]uint)
		if ok {
			if len(list) == 0 {
				packageLists[id%common.TaskNum].Delete(uint32(id))
				t.callback <- 0
			} else {
				data := list[0]
				packageLists[id%common.TaskNum].Store(uint32(id), list[1:])
				t.callback <- data
			}
		}
	}
}
