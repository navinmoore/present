package main

import (
	"encoding/json"
	"present/controller"
	"strconv"
	"sync"
	"testing"

	"github.com/kataras/iris/v12/httptest"
	"github.com/stretchr/testify/assert"
	"present/common"
)

func TestMVC(t *testing.T) {
	e := httptest.New(t, newApp())
	var wg sync.WaitGroup
	var tempDate float64 = 0

	result := common.Result{
		Code: common.SuccessCode,
		Msg:  common.SuccessMsg,
		Data: tempDate,
	}

	sth := e.GET("/api/v1/annual").Expect().Status(httptest.StatusOK).Body().Raw()
	var responseResult common.Result
	_ = json.Unmarshal([]byte(sth), &responseResult)

	assert.Equal(t, result, responseResult)

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()
			importRequest := controller.ImportRequest{
				Name: "c" + strconv.Itoa(i),
			}
			e.POST("/api/v1/annual/import").WithJSON(importRequest).Expect().Status(httptest.StatusOK)
		}(i)
	}

	wg.Wait()

	tempDate = 100
	resultN := common.Result{
		Code: common.SuccessCode,
		Msg:  common.SuccessMsg,
		Data: tempDate,
	}

	sth1 := e.GET("/api/v1/annual").Expect().Status(httptest.StatusOK).Body().Raw()
	_ = json.Unmarshal([]byte(sth1), &responseResult)
	assert.Equal(t, resultN, responseResult)

	sth2 := e.GET("/api/v1/annual/luck").Expect().Status(httptest.StatusOK).Body().Raw()
	_ = json.Unmarshal([]byte(sth2), &responseResult)
	assert.Equal(t, common.SuccessCode, responseResult.Code)

	tempDate = 99
	resultM := common.Result{
		Code: common.SuccessCode,
		Msg:  common.SuccessMsg,
		Data: tempDate,
	}
	sth3 := e.GET("/api/v1/annual").Expect().Status(httptest.StatusOK).Body().Raw()
	_ = json.Unmarshal([]byte(sth3), &responseResult)
	assert.Equal(t, resultM, responseResult)

}
