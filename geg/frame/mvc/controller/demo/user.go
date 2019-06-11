package demo

import (
<<<<<<< HEAD
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g/frame/gmvc"
)

// 定义业务相关的控制器对象，
// 建议命名规范中控制器统一使用Controller前缀，后期代码维护时便于区分
type ControllerUser struct {
    gmvc.Controller
}

// 初始化控制器对象，并绑定操作到Web Server
func init() {
    // 绑定控制器到指定URI，所有控制器的公开方法将会映射到指定URI末尾
    // 例如该方法执行后，查看效果可访问：
    // http://127.0.0.1:8199/user/name
    // http://127.0.0.1:8199/user/age
    // http://127.0.0.1:8199/user/true-name
    ghttp.GetServer().BindController("/user", &ControllerUser{})
}

// 定义操作逻辑 - 展示姓名
func (c *ControllerUser) Name() {
    c.Response.Write("John")
}

// 定义操作逻辑 - 展示年龄
func (c *ControllerUser) Age() {
    c.Response.Write("18")
}

// 定义操作逻辑 - 展示方法名称如果带多个单词，路由控制器使用英文连接符号"-"进行拼接
func (c *ControllerUser) TrueName() {
    c.View.Assigns(map[string]interface{}{

    })
    c.Response.Write("John Smith")
}



=======
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/frame/gmvc"
)

type User struct {
	gmvc.Controller
}

func init() {
	s := g.Server()
	s.BindController("/user", new(User))
	s.BindController("/user/{.method}/{uid}", new(User), "Info")
	s.BindController("/user/{.method}/{page}.html", new(User), "List")
}

func (u *User) Index() {
	u.Response.Write("User")
}

func (u *User) Info() {
	u.Response.Write("Info - Uid: ", u.Request.Get("uid"))
}

func (u *User) List() {
	u.Response.Write("List - Page: ", u.Request.Get("page"))
}
>>>>>>> upstream/master
