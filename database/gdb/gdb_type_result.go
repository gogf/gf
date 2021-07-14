// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/encoding/gparser"
	"github.com/gogf/gf/util/gconv"
	"math"
)

// IsEmpty checks and returns whether `r` is empty.
func (r Result) IsEmpty() bool {
	return r.Len() == 0
}

// Len returns the length of result list.
func (r Result) Len() int {
	return len(r)
}

// Size is alias of function Len.
func (r Result) Size() int {
	return r.Len()
}

// Chunk splits an Result into multiple Results,
// the size of each array is determined by `size`.
// The last chunk may contain less than size elements.
func (r Result) Chunk(size int) []Result {
	if size < 1 {
		return nil
	}
	length := len(r)
	chunks := int(math.Ceil(float64(length) / float64(size)))
	var n []Result
	for i, end := 0, 0; chunks > 0; chunks-- {
		end = (i + 1) * size
		if end > length {
			end = length
		}
		n = append(n, r[i*size:end])
		i++
	}
	return n
}

// Json converts `r` to JSON format content.
func (r Result) Json() string {
	content, _ := gparser.VarToJson(r.List())
	return string(content)
}

// Xml converts `r` to XML format content.
func (r Result) Xml(rootTag ...string) string {
	content, _ := gparser.VarToXml(r.List(), rootTag...)
	return string(content)
}

// List converts `r` to a List.
func (r Result) List() List {
	list := make(List, len(r))
	for k, v := range r {
		list[k] = v.Map()
	}
	return list
}

// Array retrieves and returns specified column values as slice.
// The parameter `field` is optional is the column field is only one.
func (r Result) Array(field ...string) []Value {
	array := make([]Value, len(r))
	if len(r) == 0 {
		return array
	}
	key := ""
	if len(field) > 0 && field[0] != "" {
		key = field[0]
	} else {
		for k, _ := range r[0] {
			key = k
			break
		}
	}
	for k, v := range r {
		array[k] = v[key]
	}
	return array
}

// MapKeyValue converts `r` to a map[string]Value of which key is specified by `key`.
// Note that the item value may be type of slice.
func (r Result) MapKeyValue(key string) map[string]Value {
	var (
		s              = ""
		m              = make(map[string]Value)
		tempMap        = make(map[string][]interface{})
		hasMultiValues bool
	)
	for _, item := range r {
		if k, ok := item[key]; ok {
			s = k.String()
			tempMap[s] = append(tempMap[s], item)
			if len(tempMap[s]) > 1 {
				hasMultiValues = true
			}
		}
	}
	for k, v := range tempMap {
		if hasMultiValues {
			m[k] = gvar.New(v)
		} else {
			m[k] = gvar.New(v[0])
		}
	}
	return m
}

// MapKeyStr converts `r` to a map[string]Map of which key is specified by `key`.
func (r Result) MapKeyStr(key string) map[string]Map {
	m := make(map[string]Map)
	for _, item := range r {
		if v, ok := item[key]; ok {
			m[v.String()] = item.Map()
		}
	}
	return m
}

// MapKeyInt converts `r` to a map[int]Map of which key is specified by `key`.
func (r Result) MapKeyInt(key string) map[int]Map {
	m := make(map[int]Map)
	for _, item := range r {
		if v, ok := item[key]; ok {
			m[v.Int()] = item.Map()
		}
	}
	return m
}

// MapKeyUint converts `r` to a map[uint]Map of which key is specified by `key`.
func (r Result) MapKeyUint(key string) map[uint]Map {
	m := make(map[uint]Map)
	for _, item := range r {
		if v, ok := item[key]; ok {
			m[v.Uint()] = item.Map()
		}
	}
	return m
}

// RecordKeyStr converts `r` to a map[string]Record of which key is specified by `key`.
func (r Result) RecordKeyStr(key string) map[string]Record {
	m := make(map[string]Record)
	for _, item := range r {
		if v, ok := item[key]; ok {
			m[v.String()] = item
		}
	}
	return m
}

// RecordKeyInt converts `r` to a map[int]Record of which key is specified by `key`.
func (r Result) RecordKeyInt(key string) map[int]Record {
	m := make(map[int]Record)
	for _, item := range r {
		if v, ok := item[key]; ok {
			m[v.Int()] = item
		}
	}
	return m
}

// RecordKeyUint converts `r` to a map[uint]Record of which key is specified by `key`.
func (r Result) RecordKeyUint(key string) map[uint]Record {
	m := make(map[uint]Record)
	for _, item := range r {
		if v, ok := item[key]; ok {
			m[v.Uint()] = item
		}
	}
	return m
}

// Structs converts `r` to struct slice.
// Note that the parameter `pointer` should be type of *[]struct/*[]*struct.
func (r Result) Structs(pointer interface{}) (err error) {
	return gconv.StructsTag(r, pointer, OrmTagForStruct)
}
