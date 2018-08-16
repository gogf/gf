package main

import "gitee.com/johng/gf/g"

func main() {
    s := g.Server()
    s.SetDenyRoutes([]string{
        "/config*",
    })
    s.SetPort(8299)
    s.Run()
}
