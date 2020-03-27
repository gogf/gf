package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/frame/gins"
	"github.com/gogf/gf/os/glog"
)

func getNameLogger(name, path string) *glog.Logger {
	inst := gins.Get(name)
	if inst == nil {
		logger := g.Log(name)
		logConf := map[string]interface{}{
			"Path":  path,
			"Level": "ALL",
		}
		if err := logger.SetConfigWithMap(logConf); err != nil {
			panic(err)
		}
		return logger
	}
	return inst.(*glog.Logger)
}

func main() {
	alog := getNameLogger("logger.日志1", "c:/logger1")
	alog.Print(1)

	blog := getNameLogger("logger.日志2", "c:/logger2")
	blog.Print(2)
}
