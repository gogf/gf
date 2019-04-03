package demo

import (
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/frame/gmvc"
)

type Rest struct {
	gmvc.Controller
}

func init() {
	g.Server().BindControllerRest("/rest", &Rest{})
}

// RESTFul - GET
func (c *Rest) Get() {
	c.Response.Write("RESTFul HTTP Method GET")
}

// RESTFul - POST
func (c *Rest) Post() {
	c.Response.Write("RESTFul HTTP Method POST")
}

// RESTFul - DELETE
func (c *Rest) Delete() {
	c.Response.Write("RESTFul HTTP Method DELETE")
}

// 该方法无法映射，将会无法访问到
func (c *Rest) Hello() {
	c.Response.Write("Hello")
}
