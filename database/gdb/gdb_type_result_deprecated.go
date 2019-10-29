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

// Deprecated.
func (r Result) ToJson() string {
	content, _ := gparser.VarToJson(r.List())
	return string(content)
}

// Deprecated.
func (r Result) ToXml(rootTag ...string) string {
	content, _ := gparser.VarToXml(r.List(), rootTag...)
	return string(content)
}

// Deprecated.
func (r Result) ToList() List {
	l := make(List, len(r))
	for k, v := range r {
		l[k] = v.Map()
	}
	return l
}

// Deprecated.
func (r Result) ToStringMap(key string) map[string]Map {
	m := make(map[string]Map)
	for _, item := range r {
		if v, ok := item[key]; ok {
			m[v.String()] = item.Map()
		}
	}
	return m
}

// Deprecated.
func (r Result) ToIntMap(key string) map[int]Map {
	m := make(map[int]Map)
	for _, item := range r {
		if v, ok := item[key]; ok {
			m[v.Int()] = item.Map()
		}
	}
	return m
}

// Deprecated.
func (r Result) ToUintMap(key string) map[uint]Map {
	m := make(map[uint]Map)
	for _, item := range r {
		if v, ok := item[key]; ok {
			m[v.Uint()] = item.Map()
		}
	}
	return m
}

// Deprecated.
func (r Result) ToStringRecord(key string) map[string]Record {
	m := make(map[string]Record)
	for _, item := range r {
		if v, ok := item[key]; ok {
			m[v.String()] = item
		}
	}
	return m
}

// Deprecated.
func (r Result) ToIntRecord(key string) map[int]Record {
	m := make(map[int]Record)
	for _, item := range r {
		if v, ok := item[key]; ok {
			m[v.Int()] = item
		}
	}
	return m
}

// Deprecated.
func (r Result) ToUintRecord(key string) map[uint]Record {
	m := make(map[uint]Record)
	for _, item := range r {
		if v, ok := item[key]; ok {
			m[v.Uint()] = item
		}
	}
	return m
}

// Deprecated.
func (r Result) ToStructs(pointer interface{}) (err error) {
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
