package demo

import "gitee.com/johng/gf/g/net/ghttp"

type Object struct {}

func init() {
    ghttp.GetServer().BindObject("/object", &Object{})
}

func (o *Object) Show(s *ghttp.Server, r *ghttp.ClientRequest, w *ghttp.ServerResponse) {
    w.WriteString("It's show time bibi!")
}

