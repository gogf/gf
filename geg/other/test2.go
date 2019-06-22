package main

import (
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/net/ghttp"
)

type Schedule struct{}

type Task struct{}

func (c *Schedule) ListDir(r *ghttp.Request) {
	r.Response.Writeln("ListDir")
}

func (c *Task) Add(r *ghttp.Request) {
	r.Response.Writeln("Add")
}

func (c *Task) Task(r *ghttp.Request) {
	r.Response.Writeln("Task")
}

// 实现权限校验
// 通过事件回调，类似于中间件机制，但是可控制的粒度更细，可以精准注册到路由规则
func AuthHookHandler(r *ghttp.Request) {
	// 如果权限校验失败，调用 r.ExitAll() 退出执行流程
}

func main() {
	s := g.Server()
	s.Group("/schedule").Bind([]ghttp.GroupItem{
		{"ALL", "*", AuthHookHandler, ghttp.HOOK_BEFORE_SERVE},
		{"POST", "/schedule", new(Schedule)},
		{"POST", "/task", new(Task)},
	})
	s.SetPort(8199)
	s.Run()
}
