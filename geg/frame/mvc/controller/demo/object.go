package demo

import "gitee.com/johng/gf/g/net/ghttp"

type Object struct {}

func init() {
    ghttp.GetServer().BindObject("/object", &Object{})
}

func (o *Object) Index(r *ghttp.Request) {
    r.Response.Write("It's index!")
}

func (o *Object) Show(r *ghttp.Request) {
    r.Response.Write("It's show time bibi!")
}

