package demo

import (
<<<<<<< HEAD
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g/frame/gmvc"
)

type ControllerExit struct {
    gmvc.Controller
}

func (c *ControllerExit) Init(r *ghttp.Request) {
    c.Controller.Init(r)
    c.Response.Write("exit, it will not print \"show\"")
    c.Request.Exit()
}

func (c *ControllerExit) Show() {
    c.Response.Write("show")
}

func init() {
    ghttp.GetServer().BindController("/exit", &ControllerExit{})
=======
	"github.com/gogf/gf/g/frame/gmvc"
	"github.com/gogf/gf/g/net/ghttp"
)

type ControllerExit struct {
	gmvc.Controller
}

func (c *ControllerExit) Init(r *ghttp.Request) {
	c.Controller.Init(r)
	c.Response.Write("exit, it will not print \"show\"")
	c.Request.Exit()
}

func (c *ControllerExit) Show() {
	c.Response.Write("show")
}

func init() {
	ghttp.GetServer().BindController("/exit", &ControllerExit{})
>>>>>>> upstream/master
}
