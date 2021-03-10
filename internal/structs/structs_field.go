// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package structs

import "reflect"

// Tag returns the value associated with key in the tag string. If there is no
// such key in the tag, Tag returns the empty string.
func (f *Field) Tag(key string) string {
	return f.Field.Tag.Get(key)
}

// IsEmbedded returns true if the given field is an anonymous field (embedded)
func (f *Field) IsEmbedded() bool {
	return f.Field.Anonymous
}

// IsExported returns true if the given field is exported.
func (f *Field) IsExported() bool {
	return f.Field.PkgPath == ""
}

// Name returns the name of the given field
func (f *Field) Name() string {
	return f.Field.Name
}

// Type returns the type of the given field
func (f *Field) Type() Type {
	return Type{
		Type: f.Field.Type,
	}
}

// Kind returns the reflect.Kind for Value of Field `f`.
func (f *Field) Kind() reflect.Kind {
	return f.Value.Kind()
}

// OriginalKind retrieves and returns the original reflect.Kind for Value of Field `f`.
func (f *Field) OriginalKind() reflect.Kind {
	var (
		kind  = f.Value.Kind()
		value = f.Value
	)
	for kind == reflect.Ptr {
		value = value.Elem()
		kind = value.Kind()
	}
	return kind
}

// FieldMap retrieves and returns struct field as map[name/tag]*Field from `pointer`.
//
// The parameter `pointer` should be type of struct/*struct.
//
// The parameter `priority` specifies the priority tag array for retrieving from high to low.
// If it's given `nil`, it returns map[name]*Field, of which the `name` is attribute name.
//
// Note that it only retrieves the exported attributes with first letter up-case from struct.
func FieldMap(pointer interface{}, priority []string) (map[string]*Field, error) {
	fields, err := getFieldValues(pointer)
	if err != nil {
		return nil, err
	}
	var (
		tagValue = ""
		mapField = make(map[string]*Field)
	)
	for _, field := range fields {
		// Only retrieve exported attributes.
		if !field.IsExported() {
			continue
		}
		tagValue = ""
		for _, p := range priority {
			tagValue = field.Tag(p)
			if tagValue != "" && tagValue != "-" {
				break
			}
		}
		tempField := field
		tempField.TagValue = tagValue
		if tagValue != "" {
			mapField[tagValue] = tempField
		} else {
			if field.IsEmbedded() {
				m, err := FieldMap(field.Value, priority)
				if err != nil {
					return nil, err
				}
				for k, v := range m {
					if _, ok := mapField[k]; !ok {
						tempV := v
						mapField[k] = tempV
					}
				}
			} else {
				mapField[field.Name()] = tempField
			}
		}
	}
	return mapField, nil
}
