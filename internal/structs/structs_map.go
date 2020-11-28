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
// Note that it only retrieves the exported attributes with first letter up-case from struct.
func MapField(pointer interface{}, priority []string) (map[string]*Field, error) {
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
				m, err := MapField(field.value, priority)
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
