package demo

import (
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/net/ghttp"
)

type ObjectMethod struct{}

func init() {
	obj := &ObjectMethod{}
	g.Server().BindObject("/object-method", obj, "Show1, Show2, Show3")
	g.Server().BindObjectMethod("/object-method-show1", obj, "Show1")
	g.Server().Domain("localhost").BindObject("/object-method", obj, "Show4")
}

func (o *ObjectMethod) Show1(r *ghttp.Request) {
	r.Response.Write("show 1")
}

func (o *ObjectMethod) Show2(r *ghttp.Request) {
	r.Response.Write("show 2")
}

func (o *ObjectMethod) Show3(r *ghttp.Request) {
	r.Response.Write("show 3")
}

func (o *ObjectMethod) Show4(r *ghttp.Request) {
	r.Response.Write("show 4")
}
