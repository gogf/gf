// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package utils

// MapPossibleItemByKey tries to find the possible key-value pair for given key ignoring cases and symbols.
//
// Note that this function might be of low performance.
func MapPossibleItemByKey(data map[string]any, key string) (foundKey string, foundValue any) {
	if len(data) == 0 {
		return
	}
	if v, ok := data[key]; ok {
		return key, v
	}
	// Loop checking.
	for k, v := range data {
		if EqualFoldWithoutChars(k, key) {
			return k, v
		}
	}
	return "", nil
}

// MapContainsPossibleKey checks if the given `key` is contained in given map `data`.
// It checks the key ignoring cases and symbols.
//
// Note that this function might be of low performance.
func MapContainsPossibleKey(data map[string]any, key string) bool {
	if k, _ := MapPossibleItemByKey(data, key); k != "" {
		return true
	}
	return false
}
