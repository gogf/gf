// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package sqlite

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
)

const (
	fieldTypeDate       = "date"
	fieldTypeDatetime   = "datetime"
	fieldTypeTimestamp  = "timestamp"
	fieldTypeTimestampz = "timestamptz"
)

func (d *Driver) CheckLocalTypeForField(ctx context.Context, fieldType string, fieldValue interface{}) (gdb.LocalType, error) {
	typeName := strings.ToLower(fieldType)
	switch typeName {
	case fieldTypeDate:
		return gdb.LocalTypeDate, nil
	case fieldTypeDatetime, fieldTypeTimestamp, fieldTypeTimestampz:
		return gdb.LocalTypeDatetime, nil
	default:
		if strings.Contains(typeName, "time") {
			return gdb.LocalTypeString, nil
		}
		return d.Core.CheckLocalTypeForField(ctx, fieldType, fieldValue)
	}
}
