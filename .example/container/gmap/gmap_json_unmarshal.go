package main

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/container/gmap"
)

func main() {
	m := gmap.Map{}
	s := []byte(`{"name":"john","score":100}`)
	json.UnmarshalUseNumber(s, &m)
	fmt.Println(m.Map())
}
