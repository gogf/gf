// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gstructs

import (
	"reflect"
	"strconv"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gtag"
)

// ParseTag parses tag string into map.
// For example:
// ParseTag(`v:"required" p:"id" d:"1"`) => map[v:required p:id d:1].
func ParseTag(tag string) map[string]string {
	var (
		key  string
		data = make(map[string]string)
	)
	for tag != "" {
		// Skip leading space.
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}
		// Scan to colon. A space, a quote or a control character is a syntax error.
		// Strictly speaking, control chars include the range [0x7f, 0x9f], not just
		// [0x00, 0x1f], but in practice, we ignore the multi-byte control characters
		// as it is simpler to inspect the tag's bytes than the tag's runes.
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != ':' && tag[i] != '"' && tag[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		key = tag[:i]
		tag = tag[i+1:]

		// Scan quoted string to find value.
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		quotedValue := tag[:i+1]
		tag = tag[i+1:]
		value, err := strconv.Unquote(quotedValue)
		if err != nil {
			panic(gerror.WrapCodef(gcode.CodeInvalidParameter, err, `error parsing tag "%s"`, tag))
		}
		data[key] = gtag.Parse(value)
	}
	return data
}

// TagFields retrieves and returns struct tags as []Field from `pointer`.
//
// The parameter `pointer` should be type of struct/*struct.
//
// Note that,
// 1. It only retrieves the exported attributes with first letter upper-case from struct.
// 2. The parameter `priority` should be given, it only retrieves fields that has given tag.
func TagFields(pointer any, priority []string) ([]Field, error) {
	return getFieldValuesByTagPriority(pointer, priority, map[string]struct{}{})
}

// TagMapName retrieves and returns struct tags as map[tag]attribute from `pointer`.
//
// The parameter `pointer` should be type of struct/*struct.
//
// Note that,
// 1. It only retrieves the exported attributes with first letter upper-case from struct.
// 2. The parameter `priority` should be given, it only retrieves fields that has given tag.
// 3. If one field has no specified tag, it uses its field name as result map key.
func TagMapName(pointer any, priority []string) (map[string]string, error) {
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

// TagMapField retrieves struct tags as map[tag]Field from `pointer`, and returns it.
// The parameter `object` should be either type of struct/*struct/[]struct/[]*struct.
//
// Note that,
// 1. It only retrieves the exported attributes with first letter upper-case from struct.
// 2. The parameter `priority` should be given, it only retrieves fields that has given tag.
// 3. If one field has no specified tag, it uses its field name as result map key.
func TagMapField(object any, priority []string) (map[string]Field, error) {
	fields, err := TagFields(object, priority)
	if err != nil {
		return nil, err
	}
	tagMap := make(map[string]Field, len(fields))
	for _, field := range fields {
		tagField := field
		tagMap[field.TagValue] = tagField
	}
	return tagMap, nil
}

func getFieldValues(structObject any) ([]Field, error) {
	var (
		reflectValue reflect.Value
		reflectKind  reflect.Kind
	)
	if v, ok := structObject.(reflect.Value); ok {
		reflectValue = v
		reflectKind = reflectValue.Kind()
	} else {
		reflectValue = reflect.ValueOf(structObject)
		reflectKind = reflectValue.Kind()
	}
	for {
		switch reflectKind {
		case reflect.Pointer:
			if !reflectValue.IsValid() || reflectValue.IsNil() {
				// If pointer is type of *struct and nil, then automatically create a temporary struct.
				reflectValue = reflect.New(reflectValue.Type().Elem()).Elem()
				reflectKind = reflectValue.Kind()
			} else {
				reflectValue = reflectValue.Elem()
				reflectKind = reflectValue.Kind()
			}
		case reflect.Array, reflect.Slice:
			reflectValue = reflect.New(reflectValue.Type().Elem()).Elem()
			reflectKind = reflectValue.Kind()
		default:
			goto exitLoop
		}
	}

exitLoop:
	for reflectKind == reflect.Pointer {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Kind()
	}
	if reflectKind != reflect.Struct {
		return nil, gerror.NewCode(
			gcode.CodeInvalidParameter,
			"given value should be either type of struct/*struct/[]struct/[]*struct",
		)
	}
	var (
		structType = reflectValue.Type()
		length     = reflectValue.NumField()
		fields     = make([]Field, length)
	)
	for i := 0; i < length; i++ {
		fields[i] = Field{
			Value: reflectValue.Field(i),
			Field: structType.Field(i),
		}
	}
	return fields, nil
}

func getFieldValuesByTagPriority(
	pointer any, priority []string, repeatedTagFilteringMap map[string]struct{},
) ([]Field, error) {
	fields, err := getFieldValues(pointer)
	if err != nil {
		return nil, err
	}
	var (
		tagName   string
		tagValue  string
		tagFields = make([]Field, 0)
	)
	for _, field := range fields {
		// Only retrieve exported attributes.
		if !field.IsExported() {
			continue
		}
		tagValue = ""
		for _, p := range priority {
			tagName = p
			tagValue = field.Tag(p)
			if tagValue != "" && tagValue != "-" {
				break
			}
		}
		if tagValue != "" {
			// Filter repeated tag.
			if _, ok := repeatedTagFilteringMap[tagValue]; ok {
				continue
			}
			tagField := field
			tagField.TagName = tagName
			tagField.TagValue = tagValue
			tagFields = append(tagFields, tagField)
		}
		// If this is an embedded attribute, it retrieves the tags recursively.
		if field.IsEmbedded() && field.OriginalKind() == reflect.Struct {
			subTagFields, err := getFieldValuesByTagPriority(field.Value, priority, repeatedTagFilteringMap)
			if err != nil {
				return nil, err
			} else {
				tagFields = append(tagFields, subTagFields...)
			}
		}
	}
	return tagFields, nil
}
