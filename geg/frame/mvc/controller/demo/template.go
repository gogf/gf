package demo

import (
	"github.com/gogf/gf/g/frame/gmvc"
	"github.com/gogf/gf/g/net/ghttp"
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
