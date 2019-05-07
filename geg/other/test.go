package main

import (
	"github.com/gogf/gf/g/net/ghttp"
	"github.com/gogf/gf/g/os/gfile"
)

func main() {
if r, e := ghttp.Get("https://goframe.org/cover.png"); e != nil {
	panic(e)
} else {
	defer r.Close()
	gfile.PutBinContents("/Users/john/Temp/cover.png", r.ReadAll())
}
}