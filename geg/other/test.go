package main

import (
    "gitee.com/johng/gf/g/util/gutil"
    "gitee.com/johng/gf/g/os/genv"
)

func main() {
    gutil.Dump(genv.All())
}