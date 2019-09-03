package main

import (
	"fmt"

	"github.com/gogf/gf/encoding/gcompress"
	"github.com/gogf/gf/os/gfile"
)

func main() {
	fmt.Println(gfile.Basename("/dir/*"))
	return
	err := gcompress.ZipPath(
		"/Users/john/Workspace/Go/GOPATH/src/github.com/gogf/gf/.example/other",
		"/Users/john/Temp/test.zip",
	)
	if err != nil {
		panic(err)
	}
}
