package main

import (
	"fmt"

	"github.com/jin502437344/gf/encoding/gcompress"
)

func main() {
	err := gcompress.UnZipFile(
		`D:\Workspace\Go\GOPATH\src\github.com\jin502437344\gf\geg\encoding\gcompress\data.zip`,
		`D:\Workspace\Go\GOPATH\src\github.com\jin502437344\gf\geg`,
	)
	fmt.Println(err)
}
