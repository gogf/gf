package main

import (
	"fmt"
	"github.com/gogf/gf/encoding/gparser"
)

type User struct {
	Name string `xml:"name" json:"name"`
	Age  int    `xml:"bb" json:"dd" gconv:"aa"`
	Addr string `xml:"cc"`
}

func main() {
	user := User{
		Name: "sss",
		Age:  22,
		Addr: "kaldsj",
	}

	xmlStr, err := gparser.VarToXmlIndent(user, "user")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(xmlStr))
}
