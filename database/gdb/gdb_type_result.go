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

// 将结果集转换为JSON字符串
func (r Result) Json() string {
	content, _ := gparser.VarToJson(r.List())
	return string(content)
}

// 将结果集转换为XML字符串
func (r Result) Xml(rootTag ...string) string {
	content, _ := gparser.VarToXml(r.List(), rootTag...)
	return string(content)
}

// 将结果集转换为List类型返回，便于json处理
func (r Result) List() List {
	l := make(List, len(r))
	for k, v := range r {
		l[k] = v.Map()
	}
	return l
}

// 将结果列表按照指定的字段值做map[string]Map
func (r Result) MapKeyStr(key string) map[string]Map {
	m := make(map[string]Map)
	for _, item := range r {
		if v, ok := item[key]; ok {
			m[v.String()] = item.Map()
		}
	}
	return m
}

// 将结果列表按照指定的字段值做map[int]Map
func (r Result) MapKeyInt(key string) map[int]Map {
	m := make(map[int]Map)
	for _, item := range r {
		if v, ok := item[key]; ok {
			m[v.Int()] = item.Map()
		}
	}
	return m
}

// 将结果列表按照指定的字段值做map[uint]Map
func (r Result) MapKeyUint(key string) map[uint]Map {
	m := make(map[uint]Map)
	for _, item := range r {
		if v, ok := item[key]; ok {
			m[v.Uint()] = item.Map()
		}
	}
	return m
}

// 将结果列表按照指定的字段值做map[string]Record
func (r Result) RecordKeyStr(key string) map[string]Record {
	m := make(map[string]Record)
	for _, item := range r {
		if v, ok := item[key]; ok {
			m[v.String()] = item
		}
	}
	return m
}

// 将结果列表按照指定的字段值做map[int]Record
func (r Result) RecordKeyInt(key string) map[int]Record {
	m := make(map[int]Record)
	for _, item := range r {
		if v, ok := item[key]; ok {
			m[v.Int()] = item
		}
	}
	return m
}

// 将结果列表按照指定的字段值做map[uint]Record
func (r Result) RecordKeyUint(key string) map[uint]Record {
	m := make(map[uint]Record)
	for _, item := range r {
		if v, ok := item[key]; ok {
			m[v.Uint()] = item
		}
	}
	return m
}

// 将结果列表转换为指定对象的slice。
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
