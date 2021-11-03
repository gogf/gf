package main

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

func main() {
	ghttp.PostContent("http://127.0.0.1:8199/", "array[]=1&array[]=2")
}
