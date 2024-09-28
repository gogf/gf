// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gjson

import "github.com/gogf/gf/v2/os/gfile"

// Load loads content from specified file `path`, and creates a Json object from its content.
// Deprecated: use LoadPath instead.
func Load(path string, safe ...bool) (*Json, error) {
	var isSafe bool
	if len(safe) > 0 {
		isSafe = safe[0]
	}
	return LoadPath(path, Options{
		Safe: isSafe,
	})
}

// LoadPath loads content from specified file `path`, and creates a Json object from its content.
func LoadPath(path string, options Options) (*Json, error) {
	if p, err := gfile.Search(path); err != nil {
		return nil, err
	} else {
		path = p
	}
	if options.Type == "" {
		options.Type = ContentType(gfile.Ext(path))
	}
	return loadContentWithOptions(gfile.GetBytesWithCache(path), options)
}
