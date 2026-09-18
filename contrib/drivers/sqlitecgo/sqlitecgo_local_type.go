// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package sqlitecgo

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
)

// CheckLocalTypeForField checks and returns corresponding local golang type for given db type.
//
// SQLite has no fixed set of type names to map: a declared type is free text, and it only
// gives the column an affinity, as described in https://www.sqlite.org/datatype3.html#affinity.
// The affinity is a preference, not a guarantee. A value that does not fit it keeps its own
// storage class, so a column declared `point` has INTEGER affinity, because its name contains
// "int", while the value `(1.5,2.5)` is still stored as text.
//
// The declared type therefore decides the local type only as far as the value agrees with it.
// When the name suggests a number but the driver hands over text, SQLite has already
// determined the value is not a number, and coercing it anyway would replace it with 0.
func (d *Driver) CheckLocalTypeForField(ctx context.Context, fieldType string, fieldValue any) (gdb.LocalType, error) {
	localType, err := d.Core.CheckLocalTypeForField(ctx, fieldType, fieldValue)
	if err != nil {
		return localType, err
	}
	switch localType {
	case
		gdb.LocalTypeInt, gdb.LocalTypeUint,
		gdb.LocalTypeInt32, gdb.LocalTypeUint32,
		gdb.LocalTypeInt64, gdb.LocalTypeUint64,
		gdb.LocalTypeBigInt,
		gdb.LocalTypeFloat32, gdb.LocalTypeFloat64:
		switch fieldValue.(type) {
		case string, []byte:
			return gdb.LocalTypeString, nil
		default:
		}
	default:
	}
	return localType, nil
}
