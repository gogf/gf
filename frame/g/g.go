// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package g

import (
	"github.com/gogf/gf/container/gvar"
)

// Var is a universal variable interface, like generics.
type Var = gvar.Var

// Frequently-used map type alias.
type Map = map[string]interface{}
type MapAnyAny = map[interface{}]interface{}
type MapAnyStr = map[interface{}]string
type MapAnyInt = map[interface{}]int
type MapStrAny = map[string]interface{}
type MapStrStr = map[string]string
type MapStrInt = map[string]int
type MapIntAny = map[int]interface{}
type MapIntStr = map[int]string
type MapIntInt = map[int]int
type MapAnyBool = map[interface{}]bool
type MapStrBool = map[string]bool
type MapIntBool = map[int]bool

// Frequently-used slice type alias.
type List = []Map
type ListAnyAny = []Map
type ListAnyStr = []MapAnyStr
type ListAnyInt = []MapAnyInt
type ListStrAny = []MapStrAny
type ListStrStr = []MapStrStr
type ListStrInt = []MapStrInt
type ListIntAny = []MapIntAny
type ListIntStr = []MapIntStr
type ListIntInt = []MapIntInt
type ListAnyBool = []MapAnyBool
type ListStrBool = []MapStrBool
type ListIntBool = []MapIntBool

// Frequently-used slice type alias.
type Slice = []interface{}
type SliceAny = []interface{}
type SliceStr = []string
type SliceInt = []int

// Array is alias of Slice.
type Array = []interface{}
type ArrayAny = []interface{}
type ArrayStr = []string
type ArrayInt = []int
