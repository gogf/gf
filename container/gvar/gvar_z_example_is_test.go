// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
package gvar_test

import (
	"github.com/gogf/gf/v2/frame/g"
)

// IsNil
func ExampleVar_IsNil() {
	g.Dump(g.NewVar(0).IsNil())
	g.Dump(g.NewVar(0.1).IsNil())
	// true
	g.Dump(g.NewVar(nil).IsNil())
	g.Dump(g.NewVar("").IsNil())
	g.Dump(g.NewVar(g.Map{}).IsNil())
	g.Dump(g.NewVar(g.Map{"k": "v"}).IsNil())
	g.Dump(g.NewVar(g.Slice{}).IsNil())
	g.Dump(g.NewVar(g.Slice{0}).IsNil())

	// Output:
	// false
	// false
	// true
	// false
	// false
	// false
	// false
	// false
}

// IsEmpty
func ExampleVar_IsEmpty() {
	g.Dump(g.NewVar(0).IsEmpty())
	g.Dump(g.NewVar(0.1).IsEmpty())
	g.Dump(g.NewVar(nil).IsEmpty())
	g.Dump(g.NewVar("").IsEmpty())
	g.Dump(g.NewVar(g.Map{}).IsEmpty())
	g.Dump(g.NewVar(g.Map{"k": "v"}).IsEmpty())
	g.Dump(g.NewVar(g.Slice{}).IsEmpty())
	g.Dump(g.NewVar(g.Slice{0}).IsEmpty())

	// Output:
	// true
	// false
	// true
	// true
	// true
	// false
	// true
	// false
}

// IsInt
func ExampleVar_IsInt() {
	g.Dump(g.NewVar(0).IsInt())
	g.Dump(g.NewVar(-1).IsInt())
	g.Dump(g.NewVar(uint8(8)).IsInt())
	g.Dump(g.NewVar(0.1).IsInt())
	g.Dump(g.NewVar(nil).IsInt())
	g.Dump(g.NewVar("").IsInt())
	g.Dump(g.NewVar(g.Map{}).IsInt())
	g.Dump(g.NewVar(g.Map{"k": "v"}).IsInt())
	g.Dump(g.NewVar(g.Slice{}).IsInt())
	g.Dump(g.NewVar(g.Slice{0}).IsInt())

	// Output:
	// true
	// true
	// false
	// false
	// false
	// false
	// false
	// false
	// false
	// false
}

// IsUint
func ExampleVar_IsUint() {
	g.Dump(g.NewVar(0).IsUint())
	g.Dump(g.NewVar(-1).IsUint())
	g.Dump(g.NewVar(uint8(8)).IsUint())
	g.Dump(g.NewVar(0.1).IsUint())
	g.Dump(g.NewVar(nil).IsUint())
	g.Dump(g.NewVar("").IsUint())
	g.Dump(g.NewVar(g.Map{}).IsUint())
	g.Dump(g.NewVar(g.Map{"k": "v"}).IsUint())
	g.Dump(g.NewVar(g.Slice{}).IsUint())
	g.Dump(g.NewVar(g.Slice{0}).IsUint())

	// Output:
	// false
	// false
	// true
	// false
	// false
	// false
	// false
	// false
	// false
	// false
}

// IsFloat
func ExampleVar_IsFloat() {
	g.Dump(g.NewVar(0).IsFloat())
	g.Dump(g.NewVar(-1).IsFloat())
	g.Dump(g.NewVar(uint8(8)).IsFloat())
	g.Dump(g.NewVar(float64(8)).IsFloat())
	g.Dump(g.NewVar(0.1).IsFloat())
	g.Dump(g.NewVar(nil).IsFloat())
	g.Dump(g.NewVar("").IsFloat())
	g.Dump(g.NewVar(g.Map{}).IsFloat())
	g.Dump(g.NewVar(g.Map{"k": "v"}).IsFloat())
	g.Dump(g.NewVar(g.Slice{}).IsFloat())
	g.Dump(g.NewVar(g.Slice{0}).IsFloat())

	// Output:
	// false
	// false
	// false
	// true
	// true
	// false
	// false
	// false
	// false
	// false
	// false
}

// IsSlice
func ExampleVar_IsSlice() {
	g.Dump(g.NewVar(0).IsSlice())
	g.Dump(g.NewVar(-1).IsSlice())
	g.Dump(g.NewVar(uint8(8)).IsSlice())
	g.Dump(g.NewVar(float64(8)).IsSlice())
	g.Dump(g.NewVar(0.1).IsSlice())
	g.Dump(g.NewVar(nil).IsSlice())
	g.Dump(g.NewVar("").IsSlice())
	g.Dump(g.NewVar(g.Map{}).IsSlice())
	g.Dump(g.NewVar(g.Map{"k": "v"}).IsSlice())
	g.Dump(g.NewVar(g.Slice{}).IsSlice())
	g.Dump(g.NewVar(g.Slice{0}).IsSlice())

	// Output:
	// false
	// false
	// false
	// false
	// false
	// false
	// false
	// false
	// false
	// true
	// true
}

// IsMap
func ExampleVar_IsMap() {
	g.Dump(g.NewVar(0).IsMap())
	g.Dump(g.NewVar(-1).IsMap())
	g.Dump(g.NewVar(uint8(8)).IsMap())
	g.Dump(g.NewVar(float64(8)).IsMap())
	g.Dump(g.NewVar(0.1).IsMap())
	g.Dump(g.NewVar(nil).IsMap())
	g.Dump(g.NewVar("").IsMap())
	g.Dump(g.NewVar(g.Map{}).IsMap())
	g.Dump(g.NewVar(g.Map{"k": "v"}).IsMap())
	g.Dump(g.NewVar(g.Slice{}).IsMap())
	g.Dump(g.NewVar(g.Slice{0}).IsMap())

	// Output:
	// false
	// false
	// false
	// false
	// false
	// false
	// false
	// true
	// true
	// false
	// false
}

// IsStruct
func ExampleVar_IsStruct() {
	g.Dump(g.NewVar(0).IsStruct())
	g.Dump(g.NewVar(-1).IsStruct())
	g.Dump(g.NewVar(uint8(8)).IsStruct())
	g.Dump(g.NewVar(float64(8)).IsStruct())
	g.Dump(g.NewVar(0.1).IsStruct())
	g.Dump(g.NewVar(nil).IsStruct())
	g.Dump(g.NewVar("").IsStruct())
	g.Dump(g.NewVar(g.Map{}).IsStruct())
	g.Dump(g.NewVar(g.Map{"k": "v"}).IsStruct())
	g.Dump(g.NewVar(g.Slice{}).IsStruct())
	g.Dump(g.NewVar(g.Slice{0}).IsStruct())
	a := &struct {
	}{}
	g.Dump(g.NewVar(a).IsStruct())
	g.Dump(g.NewVar(*a).IsStruct())
	g.Dump(g.NewVar(&a).IsStruct())

	// Output:
	// false
	// false
	// false
	// false
	// false
	// false
	// false
	// false
	// false
	// false
	// false
	// true
	// true
	// true
}
