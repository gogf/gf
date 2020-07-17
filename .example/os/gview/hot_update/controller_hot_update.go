package main

import (
	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/frame/gmvc"
)

func init() {
	g.View().SetPath(`D:\Workspace\Go\GOPATH\src\gitee.com\johng\gf\geg\os\gview`)
}

// 测试控制器注册模板热更新机制
type Controller struct {
	gmvc.Controller
}

// 测试模板热更新机制
func (c *Controller) Test() {
	b, _ := c.View.Parse("gview.tpl")
	c.Response.Write(b)
}

func main() {
	s := g.Server()
	s.BindController("/", &Controller{})
	s.Run()
}
