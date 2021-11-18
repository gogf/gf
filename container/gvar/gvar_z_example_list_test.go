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
)

// ListItemValues
func ExampleVar_ListItemValues() {
	// Map
	var goods1 = g.List{
		g.Map{"id": 1, "price": 100.00},
		g.Map{"id": 2, "price": 0},
		g.Map{"id": 3, "price": nil},
	}
	var a1 = gvar.New(goods1)
	fmt.Println(a1.ListItemValues("id"))
	fmt.Println(a1.ListItemValues("price"))

	// Struct
	type Good struct {
		Id    int
		Price float64
		Cost  interface{}
	}
	var goods2 = g.Slice{
		Good{1, 100.00, 50},
		&Good{2, 0, 0},
		Good{3, 200, nil},
	}
	var a2 = gvar.New(goods2)
	fmt.Println(a2.ListItemValues("Id"))
	fmt.Println(a2.ListItemValues("Price"))
	fmt.Println(a2.ListItemValues("Cost"))

	// Output:
	// [1 2 3]
	// [100 0 <nil>]
	// [1 2 3]
	// [100 0 200]
	// [50 0 <nil>]
}

// ListItemValuesUnique
func ExampleVar_ListItemValuesUnique() {
	// Map
	var goods1 = g.List{
		g.Map{"id": 1, "price": 100.00},
		g.Map{"id": 2, "price": 100.00},
		g.Map{"id": 3, "price": nil},
	}
	var a1 = gvar.New(goods1)
	fmt.Println(a1.ListItemValuesUnique("id"))
	fmt.Println(a1.ListItemValuesUnique("price"))

	// Struct
	type Good struct {
		Id    int
		Price float64
		Cost  interface{}
	}
	var goods2 = g.Slice{
		Good{1, 100.00, 50},
		Good{2, 100, 0},
		&Good{3, 100.00, 0},
		Good{4, 200, nil},
	}
	var a2 = gvar.New(goods2)
	fmt.Println(a2.ListItemValuesUnique("Id"))
	fmt.Println(a2.ListItemValuesUnique("Price"))
	fmt.Println(a2.ListItemValuesUnique("Cost"))

	// Output:
	// [1 2 3]
	// [100 <nil>]
	// [1 2 3 4]
	// [100 200]
	// [50 0 <nil>]
}
