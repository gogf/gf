// Copyright 2012-2018 Charles Banning. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file

//	x2j_valuesAt.go: Extract values from an arbitrary XML doc that are at same level as "key". 
//                  Tag path can include wildcard characters.

package x2j

import (
	"strings"

	"github.com/clbanning/mxj"
)

// ------------------- sweep up everything for some point in the node tree ---------------------

// ValuesAtTagPath - deliver all values at the same level of the document as the specified key.
//	See ValuesAtKeyPath().
// If there are no values for the path 'nil' is returned.
// A return value of (nil, nil) means that there were no values and no errors parsing the doc.
//   'doc' is the XML document
//   'path' is a dot-separated path of tag nodes
//   'getAttrs' can be set 'true' to return attribute values for "*"-terminated path
//          If a node is '*', then everything beyond is scanned for values.
//          E.g., "doc.books' might return a single value 'book' of type []interface{}, but
//                "doc.books.*" could return all the 'book' entries as []map[string]interface{}.
//                "doc.books.*.author" might return all the 'author' tag values as []string - or
//            		"doc.books.*.author.lastname" might be required, depending on he schema.
func ValuesAtTagPath(doc, path string, getAttrs ...bool) ([]interface{}, error) {
	var a bool
	if len(getAttrs) == 1 {
		a = getAttrs[0]
	}
	m, err := mxj.NewMapXml([]byte(doc))
	if err != nil {
		return nil, err
	}

	v := ValuesAtKeyPath(m, path, a)
	return v, nil
}

// ValuesAtKeyPath - deliver all values at the same depth in a map[string]interface{} value
//	If v := ValuesAtKeyPath(m,"x.y.z") 
//	then there exists a _,vv := range v
//	such that v.(map[string]interface{})[z] == ValuesFromKeyPath(m,"x.y.z")
// If there are no values for the path 'nil' is returned.
//   'm' is the map to be walked
//   'path' is a dot-separated path of key values
//   'getAttrs' can be set 'true' to return attribute values for "*"-terminated path
//          If a node is '*', then everything beyond is walked.
//          E.g., see ValuesFromTagPath documentation.
func ValuesAtKeyPath(m map[string]interface{}, path string, getAttrs ...bool) []interface{} {
	var a bool
	if len(getAttrs) == 1 {
		a = getAttrs[0]
	}
	keys := strings.Split(path, ".")
	lenKeys := len(keys)
	ret := make([]interface{}, 0)
	if lenKeys > 1 {
		// use function in x2j_valuesFrom.go
		valuesFromKeyPath(&ret, m, keys[:lenKeys-1], a)
		if len(ret) == 0 {
			return nil
		}
	} else {
		ret = append(ret,interface{}(m))
	}

	// scan the value set and see if key occurs
	key := keys[lenKeys-1]
	// wildcard is special
	if key == "*" {
		return ret
	}
	for _, v := range ret {
		switch v.(type) {
		case map[string]interface{}:
			if _, ok := v.(map[string]interface{})[key]; ok {
				return ret
			}
		}
	}

	// no instance of key in penultimate value set
	return nil
}

