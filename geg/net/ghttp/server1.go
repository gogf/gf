package main

import "gitee.com/johng/gf/g/net/ghttp"

func main() {
    s := ghttp.GetServer()
    s.SetIndexFolder(true)
    s.SetServerRoot("/home/www/")
    s.Run()
}
