// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package g

import "github.com/gogf/gf/container/gvar"

// Var is a universal variable interface, like generics.
type Var = gvar.Var

// Frequently-used map alias.
type (
	Map        = map[string]interface{}
	MapAnyAny  = map[interface{}]interface{}
	MapAnyStr  = map[interface{}]string
	MapAnyInt  = map[interface{}]int
	MapStrAny  = map[string]interface{}
	MapStrStr  = map[string]string
	MapStrInt  = map[string]int
	MapIntAny  = map[int]interface{}
	MapIntStr  = map[int]string
	MapIntInt  = map[int]int
	MapAnyBool = map[interface{}]bool
	MapStrBool = map[string]bool
	MapIntBool = map[int]bool
)

// Frequently-used slice alias.
type (
	List        = []Map
	ListAnyAny  = []MapAnyAny
	ListAnyStr  = []MapAnyStr
	ListAnyInt  = []MapAnyInt
	ListStrAny  = []MapStrAny
	ListStrStr  = []MapStrStr
	ListStrInt  = []MapStrInt
	ListIntAny  = []MapIntAny
	ListIntStr  = []MapIntStr
	ListIntInt  = []MapIntInt
	ListAnyBool = []MapAnyBool
	ListStrBool = []MapStrBool
	ListIntBool = []MapIntBool
)

// Frequently-used slice alias.
type (
	Slice    = []interface{}
	SliceAny = []interface{}
	SliceStr = []string
	SliceInt = []int
)

// Array is alias of Slice.
type (
	Array    = []interface{}
	ArrayAny = []interface{}
	ArrayStr = []string
	ArrayInt = []int
)
