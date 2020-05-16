// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid

import (
	"strings"
)

// checkRequired checks <value> using required rules.
func checkRequired(value, ruleKey, ruleVal string, params map[string]string) bool {
	required := false
	switch ruleKey {
	// Required.
	case "required":
		required = true

	// Required unless all given field and its value are equal.
	// Example: required-if: id,1,age,18
	case "required-if":
		required = false
		array := strings.Split(ruleVal, ",")
		// It supports multiple field and value pairs.
		if len(array)%2 == 0 {
			for i := 0; i < len(array); {
				tk := array[i]
				tv := array[i+1]
				if v, ok := params[tk]; ok {
					if strings.Compare(tv, v) == 0 {
						required = true
						break
					}
				}
				i += 2
			}
		}

	// Required unless all given field and its value are not equal.
	// Example: required-unless: id,1,age,18
	case "required-unless":
		required = true
		array := strings.Split(ruleVal, ",")
		// It supports multiple field and value pairs.
		if len(array)%2 == 0 {
			for i := 0; i < len(array); {
				tk := array[i]
				tv := array[i+1]
				if v, ok := params[tk]; ok {
					if strings.Compare(tv, v) == 0 {
						required = false
						break
					}
				}
				i += 2
			}
		}

	// Required if any of given fields are not empty.
	// Example: required-with:id,name
	case "required-with":
		required = false
		array := strings.Split(ruleVal, ",")
		for i := 0; i < len(array); i++ {
			if params[array[i]] != "" {
				required = true
				break
			}
		}

	// Required if all of given fields are not empty.
	// Example: required-with:id,name
	case "required-with-all":
		required = true
		array := strings.Split(ruleVal, ",")
		for i := 0; i < len(array); i++ {
			if params[array[i]] == "" {
				required = false
				break
			}
		}

	// Required if any of given fields are empty.
	// Example: required-with:id,name
	case "required-without":
		required = false
		array := strings.Split(ruleVal, ",")
		for i := 0; i < len(array); i++ {
			if params[array[i]] == "" {
				required = true
				break
			}
		}

	// Required if all of given fields are empty.
	// Example: required-with:id,name
	case "required-without-all":
		required = true
		array := strings.Split(ruleVal, ",")
		for i := 0; i < len(array); i++ {
			if params[array[i]] != "" {
				required = false
				break
			}
		}
	}
	if required {
		return !(value == "")
	} else {
		return true
	}
}
