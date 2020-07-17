package main

import (
	"fmt"
	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/net/ghttp"
	"io/ioutil"
)

func main() {
	s := g.Server()
	s.SetIndexFolder(true)
	s.BindHandler("/", func(r *ghttp.Request) {
		body1 := r.GetBody()
		body2, _ := ioutil.ReadAll(r.Body)
		fmt.Println(body1)
		fmt.Println(body2)
		r.Response.Write("hello world")
	})
	s.SetPort(8999)
	s.Run()
}
