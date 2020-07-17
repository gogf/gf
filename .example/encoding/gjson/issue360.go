package main

import (
	"fmt"
	"github.com/jin502437344/gf/encoding/gjson"
)

func main() {
	s := `
{"apiVersion":"v1","kind":"Service","metadata":{"labels":{"name":"http-daemon"},"name":"http-daemon","namespace":"default"},"spec":{"ports":[{"name":"http-daemon","port":8080,"protocol":"TCP","targetPort":9212}],"selector":{"app":"http-daemon","version":"v0930-082326"}}}
`
	js, err := gjson.DecodeToJson(s)
	if err != nil {
		panic(err)
	}
	//g.Dump(js.ToMap())
	y, _ := js.ToYamlString()
	fmt.Println(y)
}
