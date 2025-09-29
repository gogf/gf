// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvar

import (
	"github.com/gogf/gf/v2/util/gconv"
)

// Vars is a slice of *Var.
type Vars []*Var

// Strings converts and returns `vs` as []string.
func (vs Vars) Strings() (s []string) {
	s = make([]string, 0, len(vs))
	for _, v := range vs {
		s = append(s, v.String())
	}
	return s
}

// Bools converts and returns `vs` as []bool.
func (vs Vars) Bools() (s []bool) {
	s = make([]bool, 0, len(vs))
	for _, v := range vs {
		s = append(s, v.Bool())
	}
	return s
}

// Interfaces converts and returns `vs` as []any.
func (vs Vars) Interfaces() (s []any) {
	s = make([]any, 0, len(vs))
	for _, v := range vs {
		s = append(s, v.Val())
	}
	return s
}

// Float32s converts and returns `vs` as []float32.
func (vs Vars) Float32s() (s []float32) {
	s = make([]float32, 0, len(vs))
	for _, v := range vs {
		s = append(s, v.Float32())
	}
	return s
}

// Float64s converts and returns `vs` as []float64.
func (vs Vars) Float64s() (s []float64) {
	s = make([]float64, 0, len(vs))
	for _, v := range vs {
		s = append(s, v.Float64())
	}
	return s
}

// Ints converts and returns `vs` as []Int.
func (vs Vars) Ints() (s []int) {
	s = make([]int, 0, len(vs))
	for _, v := range vs {
		s = append(s, v.Int())
	}
	return s
}

// Int8s converts and returns `vs` as []int8.
func (vs Vars) Int8s() (s []int8) {
	s = make([]int8, 0, len(vs))
	for _, v := range vs {
		s = append(s, v.Int8())
	}
	return s
}

// Int16s converts and returns `vs` as []int16.
func (vs Vars) Int16s() (s []int16) {
	s = make([]int16, 0, len(vs))
	for _, v := range vs {
		s = append(s, v.Int16())
	}
	return s
}

// Int32s converts and returns `vs` as []int32.
func (vs Vars) Int32s() (s []int32) {
	s = make([]int32, 0, len(vs))
	for _, v := range vs {
		s = append(s, v.Int32())
	}
	return s
}

// Int64s converts and returns `vs` as []int64.
func (vs Vars) Int64s() (s []int64) {
	s = make([]int64, 0, len(vs))
	for _, v := range vs {
		s = append(s, v.Int64())
	}
	return s
}

// Uints converts and returns `vs` as []uint.
func (vs Vars) Uints() (s []uint) {
	s = make([]uint, 0, len(vs))
	for _, v := range vs {
		s = append(s, v.Uint())
	}
	return s
}

// Uint8s converts and returns `vs` as []uint8.
func (vs Vars) Uint8s() (s []uint8) {
	s = make([]uint8, 0, len(vs))
	for _, v := range vs {
		s = append(s, v.Uint8())
	}
	return s
}

// Uint16s converts and returns `vs` as []uint16.
func (vs Vars) Uint16s() (s []uint16) {
	s = make([]uint16, 0, len(vs))
	for _, v := range vs {
		s = append(s, v.Uint16())
	}
	return s
}

// Uint32s converts and returns `vs` as []uint32.
func (vs Vars) Uint32s() (s []uint32) {
	s = make([]uint32, 0, len(vs))
	for _, v := range vs {
		s = append(s, v.Uint32())
	}
	return s
}

// Uint64s converts and returns `vs` as []uint64.
func (vs Vars) Uint64s() (s []uint64) {
	s = make([]uint64, 0, len(vs))
	for _, v := range vs {
		s = append(s, v.Uint64())
	}
	return s
}

// Scan converts `vs` to []struct/[]*struct.
func (vs Vars) Scan(pointer any, mapping ...map[string]string) error {
	return gconv.Structs(vs.Interfaces(), pointer, mapping...)
}
