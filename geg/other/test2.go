package main

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/os/gfile"
)

func main() {
    g.Dump(gfile.ScanDir("/var/log", "*.log, *.gz", true))
}