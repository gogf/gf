package main

import (
    "github.com/gogf/gf/g/os/glog"
    "github.com/gogf/gf/g/os/gfile"
    "github.com/gogf/gf/g"
)

func main() {
    path := "/tmp/glog-cat"
    glog.SetPath(path)
    glog.StdPrint(false).Cat("cat1").Cat("cat2").Println("test")
    list, err := gfile.ScanDir(path, "*", true)
    g.Dump(err)
    g.Dump(list)
}


