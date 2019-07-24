package main

import (
	"github.com/gogf/gf/g/os/glog"
)

func PrintLog(content string) {
	glog.Skip(0).Line().Println("line number with skip:", content)
	glog.Line(true).Println("line number without skip:", content)
}

func main() {
	PrintLog("just test")
}
