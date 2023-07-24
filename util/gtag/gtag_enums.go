// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtag

import (
	"github.com/gogf/gf/v2/internal/json"
)

var (
	// Type name => enums json.
	enumsMap = make(map[string]json.RawMessage)
)

// SetGlobalEnums sets the global enums into package.
// Note that this operation is not concurrent safety.
func SetGlobalEnums(enumsJson string) error {
	return json.Unmarshal([]byte(enumsJson), &enumsMap)
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
