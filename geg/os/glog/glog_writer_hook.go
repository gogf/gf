package main

import (
	"fmt"
	"github.com/gogf/gf/g/net/ghttp"
	"github.com/gogf/gf/g/os/glog"
	"github.com/gogf/gf/g/text/gregex"
)

type MyWriter struct {
	logger *glog.Logger
}

func (w *MyWriter) Write(p []byte) (n int, err error) {
	s := string(p)
	if gregex.IsMatchString(`\[(PANI|FATA)\]`, s) {
		fmt.Println("SERIOUS ISSUE OCCURRED!! I'd better tell monitor in first time!")
		ghttp.PostContent("http://monitor.mydomain.com", s)
	}
	return w.logger.Write(p)
}

func main() {
	glog.SetWriter(&MyWriter{
		logger : glog.New(),
	})
	glog.Fatal("FATAL ERROR")
}
