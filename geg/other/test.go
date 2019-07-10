package main

import (
	"encoding/base64"
	"fmt"
	"github.com/gogf/gf/g/encoding/gbase64"
)

func main() {
	data := "HwHsGhXMaGc==="
	datab, err := gbase64.Decode([]byte(data))
	fmt.Println(err)
	fmt.Println(datab)
	fmt.Println(string(datab))

	s, e := base64.StdEncoding.DecodeString(data)
	fmt.Println(e)
	fmt.Println(string(s))
}
