package demo

import (
	"github.com/gogf/gf/frame/gmvc"
	"github.com/gogf/gf/net/ghttp"
)

type ControllerTemplate struct {
	gmvc.Controller
}

func (c *ControllerTemplate) Info() {
	c.View.Assign("name", "john")
	c.View.Assigns(map[string]interface{}{
		"age":   18,
		"score": 100,
	})
	c.View.Display("view/user/index.tpl")
}

func init() {
	ghttp.GetServer().BindController("/template", &ControllerTemplate{})
}
