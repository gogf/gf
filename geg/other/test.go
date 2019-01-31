package main

import (
    "gitee.com/johng/gf/g/os/gspath"
    "gitee.com/johng/gf/g/util/gutil"
)

func main() {
    gutil.Dump(gspath.Get("/Users/john/Temp/config").AllPaths())
}