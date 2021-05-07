package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

type Object struct{}

func (o *Object) Show(r *ghttp.Request) {
	r.Response.Writeln("Show")
}

func (o *Object) Delete(r *ghttp.Request) {
	r.Response.Writeln("REST Delete")
}

func Handler(r *ghttp.Request) {
	r.Response.Writeln("Handler")
}

func HookHandler(r *ghttp.Request) {
	r.Response.Writeln("HOOK Handler")
}

func main() {
	s := g.Server()
	obj := new(Object)
	s.Group("/api").Bind([]ghttp.GroupItem{
		{"ALL", "*", HookHandler, ghttp.HookBeforeServe},
		{"ALL", "/handler", Handler},
		{"ALL", "/obj", obj},
		{"GET", "/obj/show", obj, "Show"},
		{"REST", "/obj/rest", obj},
	})
	s.SetPort(8199)
	s.Run()
}
