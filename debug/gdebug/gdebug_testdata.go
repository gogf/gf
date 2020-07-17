// Copyright 2019-2020 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdebug

import (
	"path/filepath"
)

// TestDataPath retrieves and returns the testdata path of current package,
// which is used for unit testing cases only.
// The optional parameter <names> specifies the its sub-folders/sub-files,
// which will be joined with current system separator and returned with the path.
func TestDataPath(names ...string) string {
	path := CallerDirectory() + string(filepath.Separator) + "testdata"
	for _, name := range names {
		path += string(filepath.Separator) + name
	}
	return path
}
