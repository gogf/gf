// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtag

import (
	"github.com/gogf/gf/v2/internal/json"
)

// Type name => enums json.
var enumsMap = make(map[string]json.RawMessage)

// SetGlobalEnums sets the global enums into package.
// Note that this operation is not concurrent safety.
func SetGlobalEnums(enumsJSON string) error {
	return json.Unmarshal([]byte(enumsJSON), &enumsMap)
}

// GetGlobalEnums retrieves and returns the global enums.
func GetGlobalEnums() (string, error) {
	enumsBytes, err := json.Marshal(enumsMap)
	if err != nil {
		return "", err
	}
	return string(enumsBytes), nil
}

// GetEnumsByType retrieves and returns the stored enums json by type name.
// The type name is like: github.com/gogf/gf/v2/encoding/gjson.ContentType
func GetEnumsByType(typeName string) string {
	return string(enumsMap[typeName])
}

// GetEnumValuesByType retrieves and returns enum values by type name.
// It supports both legacy format:
//
//	["a", "b"]
//
// and the structured format:
//
//	[{"value":"a","comment":"..."},{"value":"b","comment":"..."}]
func GetEnumValuesByType(typeName string) ([]any, error) {
	enums := enumsMap[typeName]
	if len(enums) == 0 {
		return nil, nil
	}
	var items []any
	if err := json.Unmarshal(enums, &items); err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return items, nil
	}
	firstMap, ok := items[0].(map[string]any)
	if !ok {
		return items, nil
	}
	if _, ok = firstMap["value"]; !ok {
		return items, nil
	}
	values := make([]any, 0, len(items))
	for _, item := range items {
		itemMap, ok := item.(map[string]any)
		if !ok {
			values = append(values, item)
			continue
		}
		if value, ok := itemMap["value"]; ok {
			values = append(values, value)
			continue
		}
		values = append(values, item)
	}
	return values, nil
}
