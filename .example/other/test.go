package main

import (
	"github.com/gogf/gf/os/glog"
	"runtime/debug"
)

func doPrintLog(content string) {
	debug.PrintStack()
	glog.Skip(1).Line().Println("line number with skip:", content)
	glog.Line().Println("line number without skip:", content)
}

func PrintLog(content string) {
	doPrintLog(content)
}

func main() {
	PrintLog("just test")
}
