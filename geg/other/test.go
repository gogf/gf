package main

import (
	"github.com/gogf/gf/g/os/glog"

	"github.com/gogf/gf/g/os/gcache"
)

func localCache() {
	result := gcache.GetOrSetFunc("test.key.1", func() interface{} {
		return nil
	}, 1000*60*2)
	if result == nil {
		glog.Error("未获取到值")
	} else {
		glog.Infofln("result is $v", result)
	}
}

func TestCache() {
	for i := 0; i < 100; i++ {
		localCache()
	}
}

func main() {
	TestCache()
}
