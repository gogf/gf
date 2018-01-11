package main

import (
    "fmt"
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g/os/glog"
)

func main() {
    path := "/home/john/Workspace/Go/GOPATH/src/gitee.com/johng/gf/version.go"
    r, e := ghttp.Post("http://127.0.0.1:8199/upload?page=1", "name=john&upload-file=@file:" + path)
    if e != nil {
        glog.Error(e)
    } else {
        fmt.Println(string(r.ReadAll()))
        r.Close()
    }
}