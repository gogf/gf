package main

import (
	"github.com/jin502437344/gf/os/gfile"
	"github.com/jin502437344/gf/util/gutil"
)

func main() {
	gutil.Dump(gfile.ScanDir("/Users/john/Documents", "*.*"))
	gutil.Dump(gfile.ScanDir("/home/john/temp/newproject", "*", true))
}
