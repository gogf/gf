package demo

import (
<<<<<<< HEAD
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g/frame/gmvc"
)

type ControllerTemplate struct {
    gmvc.Controller
}

func (c *ControllerTemplate) Info() {
    c.View.Assign("name", "john")
    c.View.Assigns(map[string]interface{}{
        "age"   : 18,
        "score" : 100,
    })
    c.View.Display("user/index.tpl")
}

func init() {
    ghttp.GetServer().BindController("/template", &ControllerTemplate{})
}





=======
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
>>>>>>> upstream/master
