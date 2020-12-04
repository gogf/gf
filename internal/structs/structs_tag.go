// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package structs

import (
	"errors"
	"reflect"
)

// TagFields retrieves struct tags as []*Field from <pointer>, and returns it.
//
// The parameter <pointer> should be type of struct/*struct.
//
// Note that it only retrieves the exported attributes with first letter up-case from struct.
func TagFields(pointer interface{}, priority []string) ([]*Field, error) {
	return getFieldValuesByTagPriority(pointer, priority, map[string]struct{}{})
}

// TagMapName retrieves struct tags as map[tag]attribute from <pointer>, and returns it.
//
// The parameter <pointer> should be type of struct/*struct.
//
// Note that it only retrieves the exported attributes with first letter up-case from struct.
func TagMapName(pointer interface{}, priority []string) (map[string]string, error) {
	fields, err := TagFields(pointer, priority)
	if err != nil {
		return nil, err
	}
	tagMap := make(map[string]string, len(fields))
	for _, field := range fields {
		tagMap[field.TagValue] = field.Name()
	}
	return tagMap, nil
}

// TagMapField retrieves struct tags as map[tag]*Field from <pointer>, and returns it.
//
// The parameter <pointer> should be type of struct/*struct.
//
// Note that it only retrieves the exported attributes with first letter up-case from struct.
func TagMapField(pointer interface{}, priority []string) (map[string]*Field, error) {
	fields, err := TagFields(pointer, priority)
	if err != nil {
		return nil, err
	}
	tagMap := make(map[string]*Field, len(fields))
	for _, field := range fields {
		tagField := field
		tagMap[field.TagValue] = tagField
	}
	return tagMap, nil
}

func getFieldValues(value interface{}) ([]*Field, error) {
	var (
		reflectValue reflect.Value
		reflectKind  reflect.Kind
	)
	if v, ok := value.(reflect.Value); ok {
		reflectValue = v
		reflectKind = reflectValue.Kind()
	} else {
		reflectValue = reflect.ValueOf(value)
		reflectKind = reflectValue.Kind()
	}

	for reflectKind == reflect.Ptr {
		if !reflectValue.IsValid() || reflectValue.IsNil() {
			// If pointer is type of *struct and nil, then automatically create a temporary struct.
			reflectValue = reflect.New(reflectValue.Type().Elem()).Elem()
			reflectKind = reflectValue.Kind()
		} else {
			reflectValue = reflectValue.Elem()
			reflectKind = reflectValue.Kind()
		}
	}
	if reflectKind != reflect.Struct {
		return nil, errors.New("given value should be type of struct/*struct")
	}
	var (
		structType = reflectValue.Type()
		length     = reflectValue.NumField()
		fields     = make([]*Field, length)
	)
	for i := 0; i < length; i++ {
		fields[i] = &Field{
			value: reflectValue.Field(i),
			field: structType.Field(i),
		}
	}
	return fields, nil
}

func getFieldValuesByTagPriority(pointer interface{}, priority []string, tagMap map[string]struct{}) ([]*Field, error) {
	fields, err := getFieldValues(pointer)
	if err != nil {
		return nil, err
	}
	var (
		tagValue  = ""
		tagFields = make([]*Field, 0)
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
		if tagValue != "" {
			// Filter repeated tag.
			if _, ok := tagMap[tagValue]; ok {
				continue
			}
			tagField := field
			tagField.TagValue = tagValue
			tagFields = append(tagFields, tagField)
		}
		// If this is an embedded attribute, it retrieves the tags recursively.
		if field.IsEmbedded() {
			if subTagFields, err := getFieldValuesByTagPriority(field.value, priority, tagMap); err != nil {
				return nil, err
			} else {
				tagFields = append(tagFields, subTagFields...)
			}
		}
	}
	return tagFields, nil
}
