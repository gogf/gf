package main

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/net/ghttp"
)

func main() {
    s := g.Server()
    s.BindHookHandler("/*any", ghttp.HOOK_BEFORE_SERVE, func(r *ghttp.Request) {
        r.Response.SetAllowCrossDomainRequest("*", "PUT,GET,POST,DELETE,OPTIONS")
        r.Response.Header().Set("Access-Control-Allow-Credentials", "true")
        r.Response.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, token")
    })
    s.Group("/v1").COMMON("*", func(r *ghttp.Request) {
        r.Response.WriteJson(g.Map{"name" : "john"})
    })
    s.SetPort(6789)
    s.Run()
}