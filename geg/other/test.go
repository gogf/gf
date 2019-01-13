package main

import (
    "fmt"
    "net/http"

    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/net/ghttp"
)

type Test struct {
    Name string
}

func Handler(r *ghttp.Request) {
    fmt.Println("==========")
    fmt.Println("/v1/*", r.Request.Method)
    if r.Request.Method == "OPTIONS" { // 临时解决方法，但并不完美。而且请求时还是会报错但返回是正常的。
        return						   // 注释掉这行就会跨域失败
    }
    r.Response.WriteJson(Test{Name: "hello"})
    r.Response.WriteHeader(http.StatusOK)
}

func main() {
    s := g.Server()
    s.BindHookHandler("/*any", ghttp.HOOK_BEFORE_SERVE, func(r *ghttp.Request) {
        fmt.Println("/*any", r.Request.Method)
        r.Response.SetAllowCrossDomainRequest("*", "PUT,GET,POST,DELETE,OPTIONS")
        r.Response.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, token")
        if r.Request.Method == "OPTIONS" { // 复杂请求的预处理
            r.Response.WriteHeader(202)
        }
    })
    s.Group("/v1").ALL("*", Handler, ghttp.HOOK_BEFORE_SERVE)
    s.SetPort(6789)
    s.Run()
}