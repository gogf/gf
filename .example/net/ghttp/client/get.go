package main

import (
	"fmt"
	"github.com/jin502437344/gf/net/ghttp"
)

func main() {
	r, err := ghttp.Get("http://127.0.0.1:8199/11111/11122")
	fmt.Println(err)
	fmt.Println(r.Header)
}
