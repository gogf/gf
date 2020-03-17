// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/gogf/gf/encoding/gparser"
)

// Json converts <r> to JSON format content.
func (r Result) Json() string {
	content, _ := gparser.VarToJson(r.List())
	return string(content)
}

// Xml converts <r> to XML format content.
func (r Result) Xml(rootTag ...string) string {
	content, _ := gparser.VarToXml(r.List(), rootTag...)
	return string(content)
}

// List converts <r> to a List.
func (r Result) List() List {
	list := make(List, len(r))
	for k, v := range r {
		list[k] = v.Map()
	}
	return list
}

// Array retrieves and returns specified column values as slice.
// The parameter <field> is optional is the column field is only one.
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

// MapKeyStr converts <r> to a map[string]Map of which key is specified by <key>.
func (r Result) MapKeyStr(key string) map[string]Map {
	m := make(map[string]Map)
	for _, item := range r {
		if v, ok := item[key]; ok {
			m[v.String()] = item.Map()
		}
	}
	return m
}

// MapKeyInt converts <r> to a map[int]Map of which key is specified by <key>.
func (r Result) MapKeyInt(key string) map[int]Map {
	m := make(map[int]Map)
	for _, item := range r {
		if v, ok := item[key]; ok {
			m[v.Int()] = item.Map()
		}
	}
	return m
}

// MapKeyUint converts <r> to a map[uint]Map of which key is specified by <key>.
func (r Result) MapKeyUint(key string) map[uint]Map {
	m := make(map[uint]Map)
	for _, item := range r {
		if v, ok := item[key]; ok {
			m[v.Uint()] = item.Map()
		}
	}
	return m
}

// RecordKeyInt converts <r> to a map[int]Record of which key is specified by <key>.
func (r Result) RecordKeyStr(key string) map[string]Record {
	m := make(map[string]Record)
	for _, item := range r {
		if v, ok := item[key]; ok {
			m[v.String()] = item
		}
	}
	return m
}

// RecordKeyInt converts <r> to a map[int]Record of which key is specified by <key>.
func (r Result) RecordKeyInt(key string) map[int]Record {
	m := make(map[int]Record)
	for _, item := range r {
		if v, ok := item[key]; ok {
			m[v.Int()] = item
		}
	}
	return m
}

// RecordKeyUint converts <r> to a map[uint]Record of which key is specified by <key>.
func (r Result) RecordKeyUint(key string) map[uint]Record {
	m := make(map[uint]Record)
	for _, item := range r {
		if v, ok := item[key]; ok {
			m[v.Uint()] = item
		}
	}
	return m
}

// Structs converts <r> to struct slice.
// Note that the parameter <pointer> should be type of *[]struct/*[]*struct.
func (r Result) Structs(pointer interface{}) (err error) {
	l := len(r)
	if l == 0 {
		return sql.ErrNoRows
	}
	t := reflect.TypeOf(pointer)
	if t.Kind() != reflect.Ptr {
		return fmt.Errorf("pointer should be type of pointer, but got: %v", t.Kind())
	}
	array := reflect.MakeSlice(t.Elem(), l, l)
	itemType := array.Index(0).Type()
	for i := 0; i < l; i++ {
		if itemType.Kind() == reflect.Ptr {
			e := reflect.New(itemType.Elem()).Elem()
			if err = r[i].Struct(e); err != nil {
				return err
			}
			array.Index(i).Set(e.Addr())
		} else {
			e := reflect.New(itemType).Elem()
			if err = r[i].Struct(e); err != nil {
				return err
			}
			array.Index(i).Set(e)
		}
	}
	reflect.ValueOf(pointer).Elem().Set(array)
	return nil
}

// IsEmpty checks and returns whether <r> is empty.
func (r Result) IsEmpty() bool {
	return len(r) == 0
}
