// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

// AttributeMap contains the attribute key and value as map for easy filtering.
type AttributeMap map[string]any

// Sets adds given attribute map to current map.
func (m AttributeMap) Sets(attrMap map[string]any) {
	for k, v := range attrMap {
		m[k] = v
	}
}

// Pick picks and returns attributes by given attribute keys.
func (m AttributeMap) Pick(keys ...string) Attributes {
	var attrs = make(Attributes, 0)
	for _, key := range keys {
		value, ok := m[key]
		if !ok {
			continue
		}
		attrs = append(attrs, NewAttribute(key, value))
	}
	return attrs
}

// PickEx picks and returns attributes of which the given attribute keys does not in given `keys`.
func (m AttributeMap) PickEx(keys ...string) Attributes {
	var (
		exKeyMap = make(map[string]struct{})
		attrs    = make(Attributes, 0)
	)
	for _, key := range keys {
		exKeyMap[key] = struct{}{}
	}
	for k, v := range m {
		_, ok := exKeyMap[k]
		if ok {
			continue
		}
		attrs = append(attrs, NewAttribute(k, v))
	}
	return attrs
}
