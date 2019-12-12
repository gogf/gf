package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"reflect"
)

// 结构体
type Data struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	data := Data{Name: "abcdefg"}
	data1 := Data{}
	data2 := Data{}

	g.Redis().Do("SET", "goods:id", data)
	v, _ := g.Redis().DoVar("GET", "goods:id")

	v.Struct(&data1)
	gconv.Struct(v, &data2)

	fmt.Println(v, data1, data2)
	fmt.Println(reflect.TypeOf(v), reflect.TypeOf(data1))
}
