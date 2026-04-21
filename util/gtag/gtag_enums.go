// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtag

import (
	"fmt"

	"github.com/gogf/gf/v2/internal/json"
)

// Type name => enums json.
var enumsMap = make(map[string]json.RawMessage)

// EnumItem describes one enum item with optional comment.
type EnumItem struct {
	Value   any
	Comment string
}

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
	items, err := GetEnumItemsByType(typeName)
	if err != nil {
		return nil, err
	}
	values := make([]any, 0, len(items))
	for _, item := range items {
		values = append(values, item.Value)
	}
	return values, nil
}

// GetEnumItemsByType retrieves and returns enum items by type name.
// It supports both legacy format:
//
//	["a", "b"]
//
// and the structured format:
//
//	[{"value":"a","comment":"..."},{"value":"b","comment":"..."}]
func GetEnumItemsByType(typeName string) ([]EnumItem, error) {
	enums := enumsMap[typeName]
	if len(enums) == 0 {
		return nil, nil
	}
	var rawItems []json.RawMessage
	if err := json.Unmarshal(enums, &rawItems); err != nil {
		return nil, err
	}
	items := make([]EnumItem, 0, len(rawItems))
	for _, rawItem := range rawItems {
		itemMap := make(map[string]json.RawMessage)
		if err := json.Unmarshal(rawItem, &itemMap); err == nil && len(itemMap) > 0 {
			var (
				valueRaw   json.RawMessage
				commentRaw json.RawMessage
			)
			valueRaw = itemMap["value"]
			commentRaw = itemMap["comment"]
			if len(valueRaw) > 0 {
				var enumValue any
				if err := json.Unmarshal(valueRaw, &enumValue); err != nil {
					return nil, err
				}
				comment := ""
				if len(commentRaw) > 0 {
					if err := json.Unmarshal(commentRaw, &comment); err != nil {
						var commentAny any
						if err2 := json.Unmarshal(commentRaw, &commentAny); err2 == nil {
							comment = anyToString(commentAny)
						}
					}
				}
				items = append(items, EnumItem{
					Value:   enumValue,
					Comment: comment,
				})
				continue
			}
		}
		var enumValue any
		if err := json.Unmarshal(rawItem, &enumValue); err != nil {
			return nil, err
		}
		items = append(items, EnumItem{Value: enumValue})
	}
	return items, nil
}

func anyToString(v any) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprint(v)
}
