package main

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/net/ghttp"
)

func main() {
    s := g.Server()

    s.BindHookHandlerByMap("/*", map[string]ghttp.HandlerFunc{
        ghttp.HOOK_BEFORE_SERVE: func(r *ghttp.Request) {
            if r.Router.Uri == "/logout" {
                r.Response.Write("HOOK_SERV")
                r.Exit()
            }
        },
    })

    s.BindStatusHandler(404, func(r *ghttp.Request) {
        r.Response.Writeln("This is customized 404 page")
    })
    s.BindStatusHandler(500, func(r *ghttp.Request) {
        r.Response.Writeln("This is customized 500 page")
    })
    s.BindStatusHandler(403, func(r *ghttp.Request) {
        r.Response.Writeln("This is customized 403 page")
    })

    s.BindHandler("/", func(r *ghttp.Request) {
        r.Response.Write("Hello World")
    })

    p := &P{}
    s.BindHandler("/login", p.Login)
    s.BindHandler("/logout", p.Logout)

    s.BindHandler("/api/getuser", p.GetUser)
    s.BindHandler("/api/:name", p.AnyName)

    s.SetPort(6655)
    s.Run()
}

type P struct {
}

func (p *P) Login(c *ghttp.Request) {
    c.Cookie.SetCookie("username", "sdf", "", "/", 300)
    c.Response.Write("this is login")
}

func (p *P) Logout(c *ghttp.Request) {
    c.Response.Write("this is logout")
}

func (p *P) GetUser(c *ghttp.Request) {
    c.Cookie.Remove("username", "", "/")
    c.Response.Write("this is GetUser")
}

func (p *P) AnyName(c *ghttp.Request) {
    c.Response.Write("this is AnyName")
}