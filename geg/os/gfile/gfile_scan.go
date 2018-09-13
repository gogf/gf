package main

import (
    "gitee.com/johng/gf/g/util/gutil"
    "gitee.com/johng/gf/g/os/gfile"
)

func main() {
    gutil.Dump(gfile.ScanDir("/home/john/Documents", "*"))
    gutil.Dump(gfile.ScanDir("/home/john/temp/newproject", "*", true))
}