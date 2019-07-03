package main

import (
	"github.com/gogf/gf/g/net/ghttp"
	"net/http"
)

func main() {
	s := ghttp.GetServer()
	s.BindHandler("/log/handler", func(r *ghttp.Request) {
		r.Response.WriteStatus(http.StatusNotFound, "文件找不到了")
	})
	s.SetAccessLogEnabled(true)
	s.SetErrorLogEnabled(true)
	//s.SetLogHandler(func(r *ghttp.Request, error ...interface{}) {
	//    if len(error) > 0 {
	//        // 如果是错误日志
	//        fmt.Println("错误产生了：", error[0])
	//    }
	//    // 这里是请求日志
	//    fmt.Println("请求处理完成，请求地址:", r.URL.String(), "请求结果:", r.Response.Status)
	//})
	s.SetPort(8199)
	s.Run()
}
