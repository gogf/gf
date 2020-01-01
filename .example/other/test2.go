package main

import (
	"fmt"
	"github.com/gogf/gf/net/ghttp"
)

func main() {
	r := ghttp.PostContent("http://127.0.0.1:8199/test", `<doc><id>1</id><name>john</name><password1>123Abc!@#</password1><password2>123Abc!@#</password2></doc>`)
	fmt.Println(r)
}
