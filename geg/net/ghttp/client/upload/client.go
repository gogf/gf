package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/net/ghttp"
)

func main() {
    path := "/home/john/Workspace/Go/gitee.com/johng/gf/version.go"
    r, e := ghttp.Post("http://127.0.0.1:8199/upload", "name=john&age=18&upload-file=@file:" + path)
    if e != nil {
        glog.Error(e)
    } else {
        fmt.Println(string(r.ReadAll()))
        r.Close()
    }
}