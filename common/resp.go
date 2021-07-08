package common

const (
	// 成功
	SuccessCode int    = 20000
	SuccessMsg  string = "SUCCESS"
	//参数错误
	ParamsErrorCode      int    = 40000
	ParamsErrorMsg       string = "参数错误"
	ParseParamsErrorCode int    = 40001
	ParseParamsErrorMsg  string = "解析参数错误"

	UserLstEmptyErrorCode int    = 40002
	UserLstEmptyErrorMsg  string = "人员列表不能为空"
	// 系统错误
	SystemErrorCode int    = 50000
	SystemErrorMsg  string = "系统错误，请联系管理员"
)

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

var ParseParamsErrorResult = Result{
	Code: ParseParamsErrorCode,
	Msg:  ParamsErrorMsg,
	Data: nil,
}
