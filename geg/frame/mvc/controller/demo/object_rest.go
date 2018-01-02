package demo

import "gitee.com/johng/gf/g/net/ghttp"

// 测试绑定对象
type ObjectRest struct {}

func init() {
    ghttp.GetServer().BindObjectRest("/object-rest", &ObjectRest{})
}

func (o *ObjectRest) Get(r *ghttp.Request) {
    r.Response.WriteString("It's show time bibi!")
}

