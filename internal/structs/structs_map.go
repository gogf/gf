// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package structs

// MapField retrieves struct field as map[name/tag]*Field from <pointer>, and returns the map.
//
// The parameter <pointer> should be type of struct/*struct.
//
// The parameter <priority> specifies the priority tag array for retrieving from high to low.
//
// The parameter <recursive> specifies whether retrieving the struct field recursively.
//
// Note that it only retrieves the exported attributes with first letter up-case from struct.
func MapField(pointer interface{}, priority []string) (map[string]*Field, error) {
	tagFields, err := getFieldValuesByTagPriority(pointer, priority, map[string]struct{}{})
	if err != nil {
		return nil, err
	}
	tagFieldMap := make(map[string]*Field, len(tagFields))
	for _, field := range tagFields {
		tagField := field
		tagFieldMap[field.Name()] = tagField
		if tagField.TagValue != "" {
			tagFieldMap[tagField.TagValue] = tagField
		}
	}
	return tagFieldMap, nil
}
