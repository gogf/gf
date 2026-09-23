// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
)

// localTypeMap maps every type name the driver can report to the local Go type of its value.
//
// DM reports a type name in two shapes, and both are keys here because
// GetFormattedDBTypeNameForField reduces them differently:
//
//   - TableFields reads ALL_TAB_COLUMNS and yields the name with its precision, like
//     `INTERVAL DAY TO SECOND(24)`, which is reduced to `interval day to second`.
//   - ColumnTypes yields the name alone, like `INTERVAL DAY TO SECOND`, which is reduced
//     to its first word, `interval`.
//
// Matching the name exactly, rather than looking for keywords in it, is what keeps the
// INTERVAL types from being taken for integers because their names contain "int", which
// used to convert their values to 0.
var localTypeMap = map[string]gdb.LocalType{
	// Character types. CHARACTER, LONG and RAW are aliases that the server reports under
	// their base name, and are listed for the case it reports them as declared.
	"char": gdb.LocalTypeString, "character": gdb.LocalTypeString,
	"varchar": gdb.LocalTypeString, "varchar2": gdb.LocalTypeString,
	"nchar": gdb.LocalTypeString, "nvarchar": gdb.LocalTypeString,
	"nvarchar2": gdb.LocalTypeString,
	"text":      gdb.LocalTypeString, "long": gdb.LocalTypeString,
	"longvarchar": gdb.LocalTypeString, "clob": gdb.LocalTypeString,
	"rowid": gdb.LocalTypeString,

	// Exact and approximate numerics.
	"numeric": gdb.LocalTypeString, "decimal": gdb.LocalTypeString,
	"dec": gdb.LocalTypeString, "number": gdb.LocalTypeString,
	"int": gdb.LocalTypeInt, "integer": gdb.LocalTypeInt,
	"bigint": gdb.LocalTypeInt64, "smallint": gdb.LocalTypeInt,
	"tinyint": gdb.LocalTypeInt, "byte": gdb.LocalTypeString,
	"pls_integer": gdb.LocalTypeInt, "money": gdb.LocalTypeString,
	"float": gdb.LocalTypeFloat64, "double": gdb.LocalTypeFloat64,
	"double precision": gdb.LocalTypeFloat64, "real": gdb.LocalTypeFloat32,

	// Bit and binary types. BOOLEAN is not a column type of DM, whose boolean column type
	// is BIT, but the driver reads it as a result column type name, so it is mapped too.
	"bit": gdb.LocalTypeBool, "boolean": gdb.LocalTypeBool, "bool": gdb.LocalTypeBool,
	"binary": gdb.LocalTypeBytes, "varbinary": gdb.LocalTypeBytes, "raw": gdb.LocalTypeBytes,
	"longvarbinary": gdb.LocalTypeBytes, "blob": gdb.LocalTypeBytes,
	"image": gdb.LocalTypeString, "bfile": gdb.LocalTypeString,
	"json": gdb.LocalTypeJson, "jsonb": gdb.LocalTypeJsonb,

	// Date and time types.
	"date": gdb.LocalTypeDate, "time": gdb.LocalTypeTime,
	"timestamp": gdb.LocalTypeDatetime, "datetime": gdb.LocalTypeDatetime,
	"time with time zone":            gdb.LocalTypeDatetime,
	"datetime with time zone":        gdb.LocalTypeDatetime,
	"timestamp with time zone":       gdb.LocalTypeDatetime,
	"timestamp with local time zone": gdb.LocalTypeDatetime,

	// Interval types. Their names contain "int" while holding no integer at all, which is
	// what this issue is about. Both shapes of the name are listed.
	"interval":                  gdb.LocalTypeString,
	"interval year":             gdb.LocalTypeString,
	"interval year to month":    gdb.LocalTypeString,
	"interval month":            gdb.LocalTypeString,
	"interval day":              gdb.LocalTypeString,
	"interval day to hour":      gdb.LocalTypeString,
	"interval day to minute":    gdb.LocalTypeString,
	"interval day to second":    gdb.LocalTypeString,
	"interval hour":             gdb.LocalTypeString,
	"interval hour to minute":   gdb.LocalTypeString,
	"interval hour to second":   gdb.LocalTypeString,
	"interval minute":           gdb.LocalTypeString,
	"interval minute to second": gdb.LocalTypeString,
	"interval second":           gdb.LocalTypeString,
}

// CheckLocalTypeForField checks and returns corresponding local golang type for given db type.
// The parameter `fieldType` is the type name reported by the driver, like `VARCHAR(10)`,
// `INTERVAL DAY TO SECOND(24)` or `INTERVAL DAY TO SECOND`.
func (d *Driver) CheckLocalTypeForField(ctx context.Context, fieldType string, fieldValue any) (gdb.LocalType, error) {
	var typeName string
	match, _ := gregex.MatchString(`(.+?)\((.+)\)`, fieldType)
	if len(match) == 3 {
		typeName = gstr.Trim(match[1])
	} else {
		typeName = fieldType
	}
	typeName = strings.ToLower(typeName)
	if localType, ok := localTypeMap[typeName]; ok {
		return localType, nil
	}
	return d.Core.CheckLocalTypeForField(ctx, fieldType, fieldValue)
}
