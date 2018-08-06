package main

import (
	"gitee.com/johng/gf/g/os/gfile"
)

func main() {
	path := "/home/john/Documents/temp"
	flags1 := gfile.IsFile(path)
	if flags1 == true {
		println("有")
	} else {
		println("无")
	}
}
