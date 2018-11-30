package main

import (
    "fmt"
    "gitee.com/johng/gf/g/encoding/gparser"
)

func main() {


	type DemoInfo struct {
		Name string
		Age  string
	}


		//r.Response.Write("Hello World")
		l := map[interface{}][]DemoInfo{}

		el := [...]DemoInfo{
			{Name:"Bala", Age:"15"},
			{Name:"CeCe", Age:"18"},
			{Name:"ChenLo", Age:"28"},
			{Name:"Bii", Age:"22"},
			{Name:"Ann", Age:"23"},
			{Name:"Bmx", Age:"88"},
		}

		for _,v := range el{
			l[string(v.Name[0])] = append(l[string(v.Name[0])],v)
		}
		//fmt.Println(l)



		b, err := gparser.VarToJson(l)
		fmt.Println(err)
		fmt.Println(string(b))
}
