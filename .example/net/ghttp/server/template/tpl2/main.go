package main

import (
	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/frame/gins"
	"github.com/jin502437344/gf/net/ghttp"
)

func main() {
	s := g.Server()
	s.SetServerRoot("public")
	s.SetNameToUriType(ghttp.URI_TYPE_ALLLOWER)
	s.SetErrorLogEnabled(true)
	s.SetAccessLogEnabled(true)
	s.SetPort(2333)

	s.BindHandler("/", func(r *ghttp.Request) {
		content, _ := gins.View().Parse("test.html", nil)
		r.Response.Write(content)
	})

	s.Run()
}
