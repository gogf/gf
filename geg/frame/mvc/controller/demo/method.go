package demo

import (
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/frame/gmvc"
)

type Method struct {
	gmvc.Controller
}

func init() {
	// 第三个参数指定主要注册的方法，其他方法不注册，方法名称会自动追加到给定路由后面，构成新路由
	// 以下注册会中注册两个新路由: /method/name, /method/age
	g.Server().BindController("/method", new(Method), "Name, Age")
	// 绑定路由到指定的方法执行，以下注册只会注册一个路由: /method-name
	g.Server().BindControllerMethod("/method-name", new(Method), "Name")
}

func (c *Method) Name() {
	c.Response.Write("John")
}

func (c *Method) Age() {
	c.Response.Write("18")
}

func (c *Method) Info() {
	c.Response.Write("Info")
}
