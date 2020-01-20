package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
)

type TokenRequest struct {
	Scope     string
	Watermark bool
	Policy    *g.Var
}

func main() {
	s := "123456"
	fmt.Println(s[0:2])
	fmt.Println(s[1:3])
	return
	//	s := `
	//{
	//  "policy": {"name":"john"},
	//  "scope": "pub-med-panel",
	//  "watermark": true
	//}
	//`
	var t *TokenRequest
	m := g.Map{
		"policy":    g.Map{"name": "john"},
		"scope":     "pub-med-panel",
		"watermark": true,
	}
	err := gconv.Struct(m, &t)
	fmt.Println(err)
	fmt.Println(t.Policy)
}
