// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/wangyougui/gf.

package gins

import (
	"github.com/wangyougui/gf/v2/os/gres"
)

// Resource returns an instance of Resource.
// The parameter `name` is the name for the instance.
func Resource(name ...string) *gres.Resource {
	return gres.Instance(name...)
}
