// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package structs

import (
	"reflect"

	"github.com/gogf/gf/third/github.com/fatih/structs"
)

// TagMapName retrieves struct tags as map[tag]attribute from <pointer>, and returns it.
//
// The parameter <recursive> specifies whether retrieving the struct field recursively.
//
// Note that it only retrieves the exported attributes with first letter up-case from struct.
func TagMapName(pointer interface{}, priority []string, recursive bool) map[string]string {
	tagMap := TagMapField(pointer, priority, recursive)
	if len(tagMap) > 0 {
		m := make(map[string]string, len(tagMap))
		for k, v := range tagMap {
			m[k] = v.Name()
		}
		return m
	}
	return nil
}

// TagMapField retrieves struct tags as map[tag]*Field from <pointer>, and returns it.
//
// The parameter <recursive> specifies whether retrieving the struct field recursively.
//
// Note that it only retrieves the exported attributes with first letter up-case from struct.
func TagMapField(pointer interface{}, priority []string, recursive bool) map[string]*Field {
	tagMap := make(map[string]*Field)
	fields := ([]*structs.Field)(nil)
	if v, ok := pointer.(reflect.Value); ok {
		fields = structs.Fields(v.Interface())
	} else {
		rv := reflect.ValueOf(pointer)
		kind := rv.Kind()
		if kind == reflect.Ptr {
			rv = rv.Elem()
			kind = rv.Kind()
		}
		// If pointer is type of **struct and nil, then automatically create a temporary struct,
		// which is used for structs.Fields.
		if kind == reflect.Ptr && (!rv.IsValid() || rv.IsNil()) {
			fields = structs.Fields(reflect.New(rv.Type().Elem()).Elem().Interface())
		} else {
			fields = structs.Fields(pointer)
		}
	}
	tag := ""
	name := ""
	for _, field := range fields {
		name = field.Name()
		// Only retrieve exported attributes.
		if name[0] < byte('A') || name[0] > byte('Z') {
			continue
		}

		tag = ""
		for _, p := range priority {
			tag = field.Tag(p)
			if tag != "" {
				break
			}
		}
		if tag != "" {
			tagMap[tag] = field
		}
		if recursive {
			rv := reflect.ValueOf(field.Value())
			kind := rv.Kind()
			if kind == reflect.Ptr {
				rv = rv.Elem()
				kind = rv.Kind()
			}
			if kind == reflect.Struct {
				for k, v := range TagMapField(rv, priority, true) {
					if _, ok := tagMap[k]; !ok {
						tagMap[k] = v
					}
				}
			}
		}
	}
	return tagMap
}
