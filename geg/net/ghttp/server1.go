package main

import (
    "gitee.com/johng/gf/g/net/ghttp"
)

func main() {
    s := ghttp.GetServer()
    s.SetIndexFolder(true)
    s.SetServerRoot("C:\\Documents and Settings\\Claymore\\桌面\\gf.test")
    s.SetPort(8199)
    s.Run()
}
