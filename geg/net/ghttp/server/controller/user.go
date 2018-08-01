package main

import (
    "gitee.com/johng/gf/g/frame/gmvc"
    "gitee.com/johng/gf/g"
)

type User struct {
    gmvc.Controller
}

func (c *User) Index() {
    c.View.Display("index.html")
}

func main() {
    g.View().SetPath("C:/www/static")
    s := g.Server()
    s.BindController("/user", new(User))
    s.SetPort(8199)
    s.Run()
}





