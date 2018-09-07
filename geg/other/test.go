package main

import (
    "gitee.com/johng/gf/g/os/gfile"
    "strings"
    "gitee.com/johng/gf/g/os/gfcache"
    "fmt"
)

func main() {
	files := gfile.ScanDir("/home/john/Workspace/med3-svr", true)
	for _, file := range files {
	    if strings.Index(gfcache.GetContents(file), "ENV") != -1 {
            fmt.Println(file)
        }
    }
}