package main

import (
	"github.com/jin502437344/gf/net/ghttp"
)

func main() {
	ghttp.PostContent("http://127.0.0.1:8199/", "array[]=1&array[]=2")
}
