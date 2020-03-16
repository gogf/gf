// Copyright 2019-2020 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdebug

import (
	"path/filepath"
)

// TestDataPath retrieves and returns the testdata path of current package.
// It is used for unit testing cases only.
func TestDataPath() string {
	return CallerDirectory() + string(filepath.Separator) + "testdata"
}
