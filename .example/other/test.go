package main

import (
	"fmt"
	"github.com/gogf/gf/os/gres"
)

func main() {
	//buffer := bytes.NewBuffer(nil)
	//buffer.WriteString("\x00")
	//hex.Decode()
	//if v, e := strconv.ParseInt(s[2:], 16, 64); e == nil {
	//	return v
	//}
	//s := "\x00"
	//fmt.Println([]byte(s))
	//return
	err := gres.PackToGoFile(
		"/Users/john/Workspace/Go/GOPATH/src/github.com/gogf/gf-cli/public",
		"/Users/john/Workspace/Go/GOPATH/src/github.com/gogf/gf/.example/other/config.go",
		"main",
	)
	fmt.Println(err)
}
