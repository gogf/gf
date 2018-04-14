package demo

import "gitee.com/johng/gf/g/net/ghttp"

func init() {
    ghttp.GetServer().BindHandler("/list", List)
    ghttp.GetServer().Router.SetRule(`\/list\/page\/(\d+)[\/\?]*`, "/list?page=$1&")
}

func List(r *ghttp.Request) {
    r.Response.Write("list page:" + r.GetQueryString("page"))
}