package main

import (
	"fmt"

	"github.com/jin502437344/gf/os/gfile"
	"github.com/jin502437344/gf/util/gutil"
)

var dirpath1 = "/home/john/Workspace/temp/"
var dirpath2 = "/home/john/Workspace/temp/1"
var filepath1 = "/home/john/Workspace/temp/test.php"
var filepath2 = "/tmp/tmp.test"

type BinData struct {
	name string
	age  int
}

func info() {
	fmt.Println(gfile.Info(dirpath1))
}

func scanDir() {
	gutil.Dump(gfile.ScanDir(dirpath1, "*"))
}

func getContents() {
	fmt.Printf("%s\n", gfile.GetContents(filepath1))
}

func putContents() {
	fmt.Println(gfile.PutContentsAppend(filepath2, "123"))
}

func putBinContents() {
	fmt.Println(gfile.PutBytes(filepath2, []byte("abc")))
}

func main() {
	//info()
	//getContents()
	//putContents()
	putBinContents()
	//scanDir()
}
