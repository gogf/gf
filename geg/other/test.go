package main

import (
    "fmt"
<<<<<<< HEAD
    "gitee.com/johng/gf/g/text/gstr"
)

func main() {
    fmt.Println(gstr.TrimLeftStr("gogo我爱gogo", "go"))
    fmt.Println(gstr.TrimRightStr("gogo我爱gogo", "go"))
}
=======
<<<<<<< HEAD
    "gitee.com/johng/gf/g/net/ghttp"
    "strings"
    "time"
)

func main() {
    for {
        time.Sleep(500*time.Millisecond)
        fmt.Println(strings.TrimSpace(ghttp.GetContent("http://127.0.0.1:8881")))
    }
=======
    "gitee.com/johng/gf/g/os/gfile"
)

func main() {
    fmt.Println(gfile.RealPath("config"))
>>>>>>> master
}
>>>>>>> qiangg_reuseport
