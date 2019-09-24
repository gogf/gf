package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

type User struct{}

func (u *User) ShowList(r *ghttp.Request) {
	r.Response.Write("list")
}

func main() {
	s1 := g.Server(1)
	s2 := g.Server(2)
	s3 := g.Server(3)
	s4 := g.Server(4)

	s1.SetNameToUriType(ghttp.URI_TYPE_DEFAULT)
	s2.SetNameToUriType(ghttp.URI_TYPE_FULLNAME)
	s3.SetNameToUriType(ghttp.URI_TYPE_ALLLOWER)
	s4.SetNameToUriType(ghttp.URI_TYPE_CAMEL)

	s1.BindObject("/{.struct}/{.method}", new(User))
	s2.BindObject("/{.struct}/{.method}", new(User))
	s3.BindObject("/{.struct}/{.method}", new(User))
	s4.BindObject("/{.struct}/{.method}", new(User))

	s1.SetPort(8100)
	s2.SetPort(8200)
	s3.SetPort(8300)
	s4.SetPort(8400)

	s1.Start()
	s2.Start()
	s3.Start()
	s4.Start()

	g.Wait()
}
