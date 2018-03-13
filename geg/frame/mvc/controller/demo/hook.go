package demo

import "gitee.com/johng/gf/g/net/ghttp"

func init() {
    ghttp.GetServer().BindHandler("/hook", func(r *ghttp.Request){
        r.Response.WriteString("This is hook content!\n")
    })
    ghttp.GetServer().BindHookHandlerInit("/hook", func(r *ghttp.Request){
        r.Response.WriteString("Init hook 1!\n")
    })
    ghttp.GetServer().BindHookHandlerInit("/hook", func(r *ghttp.Request){
        r.Response.WriteString("Init hook 2!\n")
    })
    ghttp.GetServer().BindHookHandlerShut("/hook", func(r *ghttp.Request){
        r.Response.WriteString("Shut hook 1!\n")
    })
    ghttp.GetServer().BindHookHandlerShut("/hook", func(r *ghttp.Request){
        r.Response.WriteString("Shut hook 2!\n")
    })
}