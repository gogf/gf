package main

import (
	"fmt"

	"github.com/gogf/gf/g/encoding/gbase64"
	"github.com/gogf/gf/g/net/ghttp"
)

type Order struct{}

func (order *Order) Get(r *ghttp.Request) {
	r.Response.Write("GET")
}

func main() {
	s := `BgsnyD6IBEzExNDUzNjEzNDg4MzYxOTMzNjQSShAJGgzns7vnu5/pgJrnn6UiOGh0dHA6Ly9wdWItbWVkLWxvZ28uaW1ncy5tZWRsaW5rZXIubmV0L25ldy1zeXN0ZW1AM3gucG5`
	b, err := gbase64.DecodeString(s)
	fmt.Println(err)
	fmt.Println(string(b))
}
