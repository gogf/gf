package main

import (
	"fmt"

	"github.com/gogf/gf/encoding/gcompress"
)

func main() {
	err := gcompress.ZipPath(
		`/Users/john/Workspace/Go/GOPATH/src/github.com/gogf/gf/test`,
		`/Users/john/Workspace/Go/GOPATH/src/github.com/gogf/gf/test.zip`,
	)
	fmt.Println(err)
}
