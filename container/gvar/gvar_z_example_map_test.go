// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
package gvar_test

import (
	"fmt"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"reflect"
)

// Map
func ExampleVar_Map() {
	var m = g.Map{"id": 1, "price": 100.00}
	var a1 = gvar.New(m)
	var b1 = a1.Map()
	fmt.Println(b1["id"])

	// Output:
	// 1
}

// MapStrAny
func ExampleVar_MapStrAny() {
	var m1 = g.Map{"id": 1, "price": 100}
	var a1 = gvar.New(m1)
	var a2 = a1.MapStrAny()
	fmt.Println(a2["price"])

	// Output:
	// 100
}

// MapStrStr
func ExampleVar_MapStrStr() {
	var m1 = g.Map{"id": 1, "price": 100}
	var a1 = gvar.New(m1)
	var a2 = a1.MapStrStr()
	fmt.Println(a2["price"] + "元")

	// Output:
	// 100元
}

// MapStrVar
func ExampleVar_MapStrVar() {
	var m1 = g.Map{"id": 1, "price": 100}
	var a1 = gvar.New(m1)
	var a2 = a1.MapStrVar()
	fmt.Println(a2["price"].Float64() * 100)

	// Output:
	// 10000
}

// MapDeep
func ExampleVar_MapDeep() {
	var m1 = g.Map{"id": 1, "price": 100}
	var m2 = g.Map{"product": m1}
	var a1 = gvar.New(m2)
	var a2 = a1.MapDeep()
	fmt.Println(a2["product"])

	// Output:
	// map[id:1 price:100]
}

// MapStrStrDeep
func ExampleVar_MapStrStrDeep() {
	var m1 = g.Map{"id": 1, "price": 100}
	var m2 = g.Map{"product": m1}
	var a1 = gvar.New(m2)
	var a2 = a1.MapStrStrDeep()
	fmt.Println(a2["product"])

	// Output:
	// {"id":1,"price":100}
}

// MapStrVarDeep
func ExampleVar_MapStrVarDeep() {
	var m1 = g.Map{"id": 1, "price": 100}
	var m2 = g.Map{"product": m1}
	var a1 = gvar.New(m2)
	var a2 = a1.MapStrVarDeep()
	fmt.Println(a2["product"])

	// Output:
	// {"id":1,"price":100}
}

// Maps
func ExampleVar_Maps() {
	var m = gvar.New(g.ListIntInt{g.MapIntInt{0: 100, 1: 200}, g.MapIntInt{0: 300, 1: 400}}).Maps()
	fmt.Println(m[0], reflect.TypeOf(m[0]))

	// Output:
	// map[0:100 1:200] map[string]interface {}
}

// MapsDeep
func ExampleVar_MapsDeep() {
	var p1 = g.MapStrAny{"product": g.Map{"id": 1, "price": 100}}
	var p2 = g.MapStrAny{"product": g.Map{"id": 2, "price": 200}}
	var a1 = gvar.New(g.ListStrAny{p1, p2})
	var a2 = a1.MapsDeep()
	fmt.Println(a2[0], reflect.TypeOf(a2[0]))

	// Output:
	// map[product:map[id:1 price:100]] map[string]interface {}
}

// MapToMap
func ExampleVar_MapToMap() {
	var m1 = gvar.New(g.MapIntInt{0: 100, 1: 200})
	var m2 = g.MapStrStr{}
	var m3 = g.MapStrAny{}
	var m4 = g.MapAnyInt{}
	m1.MapToMap(&m2)
	m1.MapToMap(&m3)
	m1.MapToMap(&m4)
	fmt.Println(m2, reflect.TypeOf(m2))
	fmt.Println(m3, reflect.TypeOf(m3))
	fmt.Println(m4, reflect.TypeOf(m4))

	// Output:
	// map[0:100 1:200] map[string]string
	// map[0:100 1:200] map[string]interface {}
	// map[0:100 1:200] map[interface {}]int
}

// MapToMaps
func ExampleVar_MapToMaps() {
	var p1 = g.MapStrAny{"product": g.Map{"id": 1, "price": 100}}
	var p2 = g.MapStrAny{"product": g.Map{"id": 2, "price": 200}}
	var a1 = gvar.New(g.ListStrAny{p1, p2})
	var a2 []g.MapStrStr
	a1.MapToMaps(&a2)
	fmt.Println(a2, reflect.TypeOf(a2))

	// Output:
	// [map[product:{"id":1,"price":100}] map[product:{"id":2,"price":200}]] []map[string]string
}

// MapToMapDeep
func ExampleVar_MapToMapDeep() {
	var p1 = gvar.New(g.MapStrAny{"product": g.Map{"id": 1, "price": 100}})
	var p2 = g.MapStrAny{}
	p1.MapToMap(&p2)
	fmt.Println(p2, reflect.TypeOf(p2))

	// Output:
	// map[product:map[id:1 price:100]] map[string]interface {}
}
