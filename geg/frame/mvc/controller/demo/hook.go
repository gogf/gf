package demo

import "gitee.com/johng/gf/g/net/ghttp"

func init() {
    ghttp.GetServer().BindHandler("/hook", func(r *ghttp.Request){
        r.Response.Write("This is hook content!\n")
    })
    ghttp.GetServer().BindHookHandlerInit("/hook", func(r *ghttp.Request){
        r.Response.Write("Init hook 1!\n")
    })
    ghttp.GetServer().BindHookHandlerInit("/hook", func(r *ghttp.Request){
        r.Response.Write("Init hook 2!\n")
    })
    ghttp.GetServer().BindHookHandlerShut("/hook", func(r *ghttp.Request){
        r.Response.Write("Shut hook 1!\n")
    })
    ghttp.GetServer().BindHookHandlerShut("/hook", func(r *ghttp.Request){
        r.Response.Write("Shut hook 2!\n")
    })
}