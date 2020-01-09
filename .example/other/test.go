package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gfile"
)

func main() {
	g.Dump(gfile.ScanDirFile("/Users/john/Temp/test", "Dockerfile", true))
	//if err := gfile.ReplaceDir("gf-empty", "app", "/Users/john/Temp/test", "*.*", true); err != nil {
	//	glog.Fatal("content replacing failed,", err.Error())
	//}

}
