package demo

import (
	"github.com/gogf/gf/g/net/ghttp"
)

type ControllerDomain struct{}

// 初始化控制器对象，并绑定操作到Web Server
func init() {
	// 只有localhost域名下才能访问该对象，
	// 对应URL为：http://localhost:8199/test/show
	// 通过该地址将无法访问到内容：http://127.0.0.1:8199/test/show
	ghttp.GetServer().Domain("localhost").BindObject("/domain", &ControllerDomain{})
}

// 用于对象映射
func (d *ControllerDomain) Show(r *ghttp.Request) {
	r.Response.Write("It's show time bibi!")
}
