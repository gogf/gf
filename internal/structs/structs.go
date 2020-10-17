// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package structs provides functions for struct conversion.
package structs

import "github.com/gqcn/structs"

// Field is alias of structs.Field.
type Field struct {
	*structs.Field
	// Retrieved tag name. There might be more than one tags in the field,
	// but only one can be retrieved according to calling function rules.
	Tag string
}
