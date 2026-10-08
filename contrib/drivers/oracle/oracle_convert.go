// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package oracle

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
)

// localTypeMap maps every type name the underlying driver can report, being the names that
// `TNSType.String` of `github.com/sijms/go-ora` returns, to the local Go type of its value.
//
// The mapping is exhaustive by design and Test_LocalTypeCoverage fails when the driver
// gains a type name that is missing here. Matching the name exactly, rather than looking
// for keywords in it, is what keeps the INTERVAL types from being taken for integers
// because their names happen to contain "int", which used to convert their values to 0.
var localTypeMap = map[string]gdb.LocalType{
	"bdouble":          gdb.LocalTypeFloat64,
	"bfloat":           gdb.LocalTypeFloat64,
	"char":             gdb.LocalTypeString,
	"charz":            gdb.LocalTypeString,
	"date":             gdb.LocalTypeDate,
	"float":            gdb.LocalTypeFloat64,
	"ibdouble":         gdb.LocalTypeFloat64,
	"ibfloat":          gdb.LocalTypeFloat64,
	"intervalds":       gdb.LocalTypeString, // was int
	"intervalds_dty":   gdb.LocalTypeString, // was int
	"intervalym":       gdb.LocalTypeString, // was int
	"intervalym_dty":   gdb.LocalTypeString, // was int
	"long":             gdb.LocalTypeString,
	"longraw":          gdb.LocalTypeString,
	"longvarchar":      gdb.LocalTypeString,
	"longvarraw":       gdb.LocalTypeString,
	"nchar":            gdb.LocalTypeString,
	"nullstr":          gdb.LocalTypeString,
	"number":           gdb.LocalTypeString,
	"ocibloblocator":   gdb.LocalTypeBytes,
	"ocicloblocator":   gdb.LocalTypeString,
	"ocidate":          gdb.LocalTypeDatetime,
	"ocifilelocator":   gdb.LocalTypeString,
	"ociref":           gdb.LocalTypeString,
	"ocistring":        gdb.LocalTypeString,
	"ocixmltype":       gdb.LocalTypeString,
	"raw":              gdb.LocalTypeString,
	"refcursor":        gdb.LocalTypeString,
	"resultset":        gdb.LocalTypeString,
	"rowid":            gdb.LocalTypeString,
	"sb1":              gdb.LocalTypeString,
	"timestamp":        gdb.LocalTypeDatetime,
	"timestampdty":     gdb.LocalTypeDatetime,
	"timestampeltz":    gdb.LocalTypeDatetime,
	"timestampltz_dty": gdb.LocalTypeDatetime,
	"timestamptz":      gdb.LocalTypeDatetime,
	"timestamptz_dty":  gdb.LocalTypeDatetime,
	"timetz":           gdb.LocalTypeDatetime,
	"uint":             gdb.LocalTypeInt,
	"urowid":           gdb.LocalTypeString,
	"varchar":          gdb.LocalTypeString,
	"varnum":           gdb.LocalTypeString,
	"varraw":           gdb.LocalTypeString,
	"xmltype":          gdb.LocalTypeString,
}

// CheckLocalTypeForField checks and returns corresponding local golang type for given db type.
// The parameter `fieldType` is the type name reported by the driver, like `NUMBER`,
// `IntervalDS_DTY` or `TimeStampTZ_DTY`.
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
