package main

import (
	"github.com/jin502437344/gf/net/ghttp"
)

func main() {
	s1 := ghttp.GetServer("s1")
	s1.SetAddr(":8080")
	s1.SetIndexFolder(true)
	s1.SetServerRoot("/home/www/static1")
	go s1.Run()

	s2 := ghttp.GetServer("s2")
	s2.SetAddr(":8081")
	s2.SetIndexFolder(true)
	s2.SetServerRoot("/home/www/static2")
	go s2.Run()

	select {}
}
