// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package util provides util functions for internal usage.
package structtag

import (
	"reflect"
	"strings"

	"github.com/gogf/gf/third/github.com/fatih/structs"
)

// Map recursively retrieves struct tags as map[tag]attribute from <pointer>, and returns it.
func Map(pointer interface{}, priority []string) map[string]string {
	tagMap := make(map[string]string)
	fields := ([]*structs.Field)(nil)
	if v, ok := pointer.(reflect.Value); ok {
		fields = structs.Fields(v.Interface())
	} else {
		fields = structs.Fields(pointer)
	}
	tag := ""
	name := ""
	for _, field := range fields {
		tag = ""
		for _, p := range priority {
			tag = field.Tag(p)
			if tag != "" {
				break
			}
		}
		name = field.Name()
		// Only retrieve exported attributes.
		if name[0] < byte('A') || name[0] > byte('Z') {
			continue
		}
		if tag != "" {
			for _, v := range strings.Split(tag, ",") {
				tagMap[strings.TrimSpace(v)] = name
			}
		} else {
			rv := reflect.ValueOf(field.Value())
			kind := rv.Kind()
			if kind == reflect.Ptr {
				rv = rv.Elem()
				kind = rv.Kind()
			}
			if kind == reflect.Struct {
				for k, v := range Map(rv, priority) {
					if _, ok := tagMap[k]; !ok {
						tagMap[k] = v
					}
				}
			}
		}
	}
	return tagMap
}
