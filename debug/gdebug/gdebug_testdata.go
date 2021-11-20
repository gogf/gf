// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdebug

import (
	"io/ioutil"
	"path/filepath"
)

// TestDataPath retrieves and returns the testdata path of current package,
// which is used for unit testing cases only.
// The optional parameter `names` specifies the sub-folders/sub-files,
// which will be joined with current system separator and returned with the path.
func TestDataPath(names ...string) string {
	path := CallerDirectory() + string(filepath.Separator) + "testdata"
	for _, name := range names {
		path += string(filepath.Separator) + name
	}
	return path
}

// TestDataContent retrieves and returns the file content for specified testdata path of current package
func TestDataContent(names ...string) string {
	path := TestDataPath(names...)
	if path != "" {
		data, err := ioutil.ReadFile(path)
		if err == nil {
			return string(data)
		}
	}
	return ""
}
