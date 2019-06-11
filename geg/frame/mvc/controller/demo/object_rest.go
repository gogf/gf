package demo

<<<<<<< HEAD
import "gitee.com/johng/gf/g/net/ghttp"

// 测试绑定对象
type ObjectRest struct {}

func init() {
    ghttp.GetServer().BindObjectRest("/object-rest", &ObjectRest{})
=======
import "github.com/gogf/gf/g/net/ghttp"

// 测试绑定对象
type ObjectRest struct{}

func init() {
	ghttp.GetServer().BindObjectRest("/object-rest", &ObjectRest{})
>>>>>>> upstream/master
}

// RESTFul - GET
func (o *ObjectRest) Get(r *ghttp.Request) {
<<<<<<< HEAD
    r.Response.Write("RESTFul HTTP Method GET")
=======
	r.Response.Write("RESTFul HTTP Method GET")
>>>>>>> upstream/master
}

// RESTFul - POST
func (c *ObjectRest) Post(r *ghttp.Request) {
<<<<<<< HEAD
    r.Response.Write("RESTFul HTTP Method POST")
=======
	r.Response.Write("RESTFul HTTP Method POST")
>>>>>>> upstream/master
}

// RESTFul - DELETE
func (c *ObjectRest) Delete(r *ghttp.Request) {
<<<<<<< HEAD
    r.Response.Write("RESTFul HTTP Method DELETE")
=======
	r.Response.Write("RESTFul HTTP Method DELETE")
>>>>>>> upstream/master
}

// 该方法无法映射，将会无法访问到
func (c *ObjectRest) Hello(r *ghttp.Request) {
<<<<<<< HEAD
    r.Response.Write("Hello")
}
=======
	r.Response.Write("Hello")
}
>>>>>>> upstream/master
