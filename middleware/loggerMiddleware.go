package middleware

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/kataras/iris/v12"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var (
	Log = logrus.New()
)

func InitLog() {
	Log.Out = os.Stdout
	Log.SetLevel(logrus.InfoLevel)

	path := "./present.log"
	writer, _ := rotatelogs.New(
		path+".%Y%m%d",
		rotatelogs.WithLinkName(path),             // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(60*24*time.Hour),    // 文件最大保存时间
		rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
	)

	pathMap := lfshook.WriterMap{
		logrus.InfoLevel:  writer,
		logrus.PanicLevel: writer,
	}
	Log.Hooks.Add(lfshook.NewHook(
		pathMap,
		&logrus.JSONFormatter{},
	))
}

func LoggerHandler(ctx iris.Context) {
	p := ctx.Request().URL.Path
	method := ctx.Request().Method
	start := time.Now()
	fields := make(map[string]interface{})
	// fields 中的
	fields["user_agent"] = ctx.Request().UserAgent()

	// 如果是POST/PUT请求，并且内容类型为JSON，则读取内容体
	if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch || method == http.MethodGet {
		body, err := ioutil.ReadAll(ctx.Request().Body)
		if err == nil {
			//这有意思，为什么要close Body
			defer ctx.Request().Body.Close()
			buf := bytes.NewBuffer(body)
			ctx.Request().Body = ioutil.NopCloser(buf)
			fields["body"] = string(body)
		}
	}
	// 执行具体的
	ctx.Next()

	//下面是返回日志
	timeConsuming := time.Since(start).Nanoseconds() / 1e6
	Log.WithFields(fields).Infof("[http] [%s] %s %d (%dms)",
		ctx.Request().Method, p, ctx.ResponseWriter().StatusCode(), timeConsuming)

}
