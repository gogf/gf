package main

import (
    "gitee.com/johng/gf/g/net/ghttp"
)

func main() {
    s := ghttp.GetServer()
    s.SetIndexFolder(true)
    s.SetServerRoot("/home/john/Workspace/Go/gf-home/static/plugin/editor.md/css")
    s.SetPort(8199)
    s.Run()
}
