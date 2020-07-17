package main

import (
	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/net/ghttp"
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
		{"ALL", "*", HookHandler, ghttp.HOOK_BEFORE_SERVE},
		{"ALL", "/handler", Handler},
		{"ALL", "/obj", obj},
		{"GET", "/obj/show", obj, "Show"},
		{"REST", "/obj/rest", obj},
	})
	s.SetPort(8199)
	s.Run()
}
