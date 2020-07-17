package main

import (
	"fmt"

	"github.com/jin502437344/gf/encoding/gcompress"
	"github.com/jin502437344/gf/os/gfile"
)

func main() {
	err := gcompress.UnZipContent(
		gfile.GetBytes(`D:\Workspace\Go\GOPATH\src\github.com\jin502437344\gf\geg\encoding\gcompress\data.zip`),
		`D:\Workspace\Go\GOPATH\src\github.com\jin502437344\gf\geg`,
	)
	fmt.Println(err)
}
