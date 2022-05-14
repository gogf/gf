// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package goai

import (
	"github.com/gogf/gf/v2/text/gstr"
)

// XExtensions stores the `x-` custom extensions.
type XExtensions map[string]string

func (oai *OpenApiV3) tagMapToXExtensions(tagMap map[string]string, extensions XExtensions) {
	for k, v := range tagMap {
		if gstr.HasPrefix(k, "x-") || gstr.HasPrefix(k, "X-") {
			extensions[k] = v
		}
	}
}
