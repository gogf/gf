package main

import (
    "fmt"
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g"
)

func main() {
    s := g.Server()

    // 多事件回调示例，事件1
    pattern1 := "/:name/info/{uid}"
    s.BindHookHandlerByMap(pattern1, map[string]ghttp.HandlerFunc {
        "BeforeServe"  : func(r *ghttp.Request){
            fmt.Println("打印到Server端终端")
        },
    })
    s.BindHandler(pattern1, func(r *ghttp.Request) {
        r.Response.Write("用户:", r.Get("name"), ", uid:", r.Get("uid"))
    })

    // 多事件回调示例，事件2
    pattern2 := "/{object}/list/{page}.java"
    s.BindHookHandlerByMap(pattern2, map[string]ghttp.HandlerFunc {
        "BeforeOutput" : func(r *ghttp.Request){
            r.Response.SetBuffer([]byte(
                fmt.Sprintf("通过事件修改输出内容, object: %s, page: %s",
                    r.Get("object"), r.GetRouterString("page"))),
            )
        },
    })
    s.SetPort(8199)
    s.Run()
}