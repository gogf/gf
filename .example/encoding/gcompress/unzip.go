package main

import (
	"fmt"

	"github.com/gogf/gf/v2/encoding/gcompress"
)

func main() {
	err := gcompress.UnZipFile(
		`D:\Workspace\Go\GOPATH\src\github.com\gogf\gf\geg\encoding\gcompress\data.zip`,
		`D:\Workspace\Go\GOPATH\src\github.com\gogf\gf\geg`,
	)
	fmt.Println(err)
}
