package main

import (
    "gitee.com/johng/gf/g/util/gutil"
    "gitee.com/johng/gf/g/os/gfile"
)

func main() {
	gutil.Dump(gfile.ScanDir("/tmp", "*test*"))
	gutil.Dump(gfile.Glob("/tmp/*", true))
}