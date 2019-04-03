package demo

import "github.com/gogf/gf/g/net/ghttp"

// 测试绑定对象
type ObjectRest struct{}

func init() {
	ghttp.GetServer().BindObjectRest("/object-rest", &ObjectRest{})
}

// RESTFul - GET
func (o *ObjectRest) Get(r *ghttp.Request) {
	r.Response.Write("RESTFul HTTP Method GET")
}

// RESTFul - POST
func (c *ObjectRest) Post(r *ghttp.Request) {
	r.Response.Write("RESTFul HTTP Method POST")
}

// RESTFul - DELETE
func (c *ObjectRest) Delete(r *ghttp.Request) {
	r.Response.Write("RESTFul HTTP Method DELETE")
}

// 该方法无法映射，将会无法访问到
func (c *ObjectRest) Hello(r *ghttp.Request) {
	r.Response.Write("Hello")
}
