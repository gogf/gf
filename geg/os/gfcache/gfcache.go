package main

import (
	"fmt"
	"github.com/gogf/gf/g/os/gfcache"
	"github.com/gogf/gf/g/os/gfile"
	"time"
)

func main() {
	s := 0
	r := ""
	path := gfile.TempDir() + gfile.Separator + "temp"
	gfile.PutContents(path, "hello")

	s = gfcache.GetSize()
	r = gfcache.GetContents(path)
	fmt.Println(s, r)

	gfile.PutContentsAppend(path, " john")

	// 等待1秒以便gfsnotify回调能处理完成
	time.Sleep(time.Second)

	s = gfcache.GetSize()
	r = gfcache.GetContents(path)
	fmt.Println(s, r)
}
