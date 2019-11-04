package main

import (
	"fmt"
	"github.com/gogf/gf/container/gmap"
)

func Test(data *gmap.Map) {
	data = gmap.New()
	fmt.Println(data)
}
func main() {
	var m *gmap.Map
	fmt.Println(m)
	Test(m)
	fmt.Println(m)
}
