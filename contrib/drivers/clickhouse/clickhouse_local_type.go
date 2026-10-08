// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
)

// localTypeMap maps every type family the server knows, as listed in
// `system.data_type_families`, to the local Go type of its value.
//
// ClickHouse type names are composed, like `Array(Nullable(Int64))` or `Decimal(10, 2)`,
// and the arguments are stripped before the lookup, so the family name is what is mapped.
//
// The mapping is exhaustive by design and Test_LocalTypeCoverage fails when the server
// gains a type family that is missing here. Matching the name exactly, rather than looking
// for keywords in it, is what keeps `Point` and the `Interval` families from being taken
// for integers because their names happen to contain "int", which used to convert their
// values to 0.
var localTypeMap = map[string]gdb.LocalType{
	"aggregatefunction":               gdb.LocalTypeString,
	"array":                           gdb.LocalTypeString,
	"bfloat16":                        gdb.LocalTypeFloat64,
	"bigint":                          gdb.LocalTypeInt64,
	"bigint signed":                   gdb.LocalTypeInt64,
	"bigint unsigned":                 gdb.LocalTypeUint64,
	"binary":                          gdb.LocalTypeBytes,
	"binary large object":             gdb.LocalTypeBytes,
	"binary varying":                  gdb.LocalTypeBytes,
	"bit":                             gdb.LocalTypeInt64Bytes, // bit(1) is a boolean; see CheckLocalTypeForField
	"blob":                            gdb.LocalTypeBytes,
	"bool":                            gdb.LocalTypeBool,
	"boolean":                         gdb.LocalTypeBool,
	"byte":                            gdb.LocalTypeString,
	"bytea":                           gdb.LocalTypeString,
	"char":                            gdb.LocalTypeString,
	"char large object":               gdb.LocalTypeString,
	"char varying":                    gdb.LocalTypeString,
	"character":                       gdb.LocalTypeString,
	"character large object":          gdb.LocalTypeString,
	"character varying":               gdb.LocalTypeString,
	"clob":                            gdb.LocalTypeString,
	"date":                            gdb.LocalTypeDate,
	"date32":                          gdb.LocalTypeDatetime,
	"datetime":                        gdb.LocalTypeDatetime,
	"datetime32":                      gdb.LocalTypeDatetime,
	"datetime64":                      gdb.LocalTypeDatetime,
	"dec":                             gdb.LocalTypeString,
	"decimal":                         gdb.LocalTypeString,
	"decimal128":                      gdb.LocalTypeString,
	"decimal256":                      gdb.LocalTypeString,
	"decimal32":                       gdb.LocalTypeString,
	"decimal64":                       gdb.LocalTypeString,
	"double":                          gdb.LocalTypeFloat64,
	"double precision":                gdb.LocalTypeFloat64,
	"dynamic":                         gdb.LocalTypeString,
	"enum":                            gdb.LocalTypeString,
	"enum16":                          gdb.LocalTypeString,
	"enum8":                           gdb.LocalTypeString,
	"fixed":                           gdb.LocalTypeString,
	"fixedstring":                     gdb.LocalTypeString,
	"float":                           gdb.LocalTypeFloat64,
	"float32":                         gdb.LocalTypeFloat64,
	"float64":                         gdb.LocalTypeFloat64,
	"geometry":                        gdb.LocalTypeString,
	"inet4":                           gdb.LocalTypeString,
	"inet6":                           gdb.LocalTypeString,
	"int":                             gdb.LocalTypeInt,
	"int signed":                      gdb.LocalTypeInt,
	"int unsigned":                    gdb.LocalTypeUint,
	"int1":                            gdb.LocalTypeInt,
	"int1 signed":                     gdb.LocalTypeInt,
	"int1 unsigned":                   gdb.LocalTypeUint,
	"int128":                          gdb.LocalTypeBigInt,
	"int16":                           gdb.LocalTypeInt,
	"int256":                          gdb.LocalTypeBigInt,
	"int32":                           gdb.LocalTypeInt,
	"int64":                           gdb.LocalTypeInt,
	"int8":                            gdb.LocalTypeInt,
	"integer":                         gdb.LocalTypeInt,
	"integer signed":                  gdb.LocalTypeInt,
	"integer unsigned":                gdb.LocalTypeUint,
	"intervalday":                     gdb.LocalTypeString, // was int
	"intervalhour":                    gdb.LocalTypeString, // was int
	"intervalmicrosecond":             gdb.LocalTypeString, // was int
	"intervalmillisecond":             gdb.LocalTypeString, // was int
	"intervalminute":                  gdb.LocalTypeString, // was int
	"intervalmonth":                   gdb.LocalTypeString, // was int
	"intervalnanosecond":              gdb.LocalTypeString, // was int
	"intervalquarter":                 gdb.LocalTypeString, // was int
	"intervalsecond":                  gdb.LocalTypeString, // was int
	"intervalweek":                    gdb.LocalTypeString, // was int
	"intervalyear":                    gdb.LocalTypeString, // was int
	"ipv4":                            gdb.LocalTypeString,
	"ipv6":                            gdb.LocalTypeString,
	"json":                            gdb.LocalTypeJson,
	"linestring":                      gdb.LocalTypeString,
	"longblob":                        gdb.LocalTypeBytes,
	"longtext":                        gdb.LocalTypeString,
	"lowcardinality":                  gdb.LocalTypeString,
	"map":                             gdb.LocalTypeString,
	"mediumblob":                      gdb.LocalTypeBytes,
	"mediumint":                       gdb.LocalTypeInt,
	"mediumint signed":                gdb.LocalTypeInt,
	"mediumint unsigned":              gdb.LocalTypeUint,
	"mediumtext":                      gdb.LocalTypeString,
	"multilinestring":                 gdb.LocalTypeString,
	"multipolygon":                    gdb.LocalTypeString,
	"national char":                   gdb.LocalTypeString,
	"national char varying":           gdb.LocalTypeString,
	"national character":              gdb.LocalTypeString,
	"national character large object": gdb.LocalTypeString,
	"national character varying":      gdb.LocalTypeString,
	"nchar":                           gdb.LocalTypeString,
	"nchar large object":              gdb.LocalTypeString,
	"nchar varying":                   gdb.LocalTypeString,
	"nested":                          gdb.LocalTypeString,
	"nothing":                         gdb.LocalTypeString,
	"nullable":                        gdb.LocalTypeString,
	"numeric":                         gdb.LocalTypeString,
	"nvarchar":                        gdb.LocalTypeString,
	"object":                          gdb.LocalTypeString,
	"point":                           gdb.LocalTypeString, // was int
	"polygon":                         gdb.LocalTypeString,
	"real":                            gdb.LocalTypeFloat32,
	"ring":                            gdb.LocalTypeString,
	"set":                             gdb.LocalTypeString,
	"signed":                          gdb.LocalTypeString,
	"simpleaggregatefunction":         gdb.LocalTypeString,
	"single":                          gdb.LocalTypeString,
	"smallint":                        gdb.LocalTypeInt,
	"smallint signed":                 gdb.LocalTypeInt,
	"smallint unsigned":               gdb.LocalTypeUint,
	"string":                          gdb.LocalTypeString,
	"text":                            gdb.LocalTypeString,
	"time":                            gdb.LocalTypeTime,
	"timestamp":                       gdb.LocalTypeDatetime,
	"tinyblob":                        gdb.LocalTypeBytes,
	"tinyint":                         gdb.LocalTypeInt,
	"tinyint signed":                  gdb.LocalTypeInt,
	"tinyint unsigned":                gdb.LocalTypeUint,
	"tinytext":                        gdb.LocalTypeString,
	"tuple":                           gdb.LocalTypeString,
	"uint128":                         gdb.LocalTypeBigInt,
	"uint16":                          gdb.LocalTypeInt,
	"uint256":                         gdb.LocalTypeBigInt,
	"uint32":                          gdb.LocalTypeInt,
	"uint64":                          gdb.LocalTypeInt,
	"uint8":                           gdb.LocalTypeInt,
	"unsigned":                        gdb.LocalTypeString,
	"uuid":                            gdb.LocalTypeString,
	"varbinary":                       gdb.LocalTypeBytes,
	"varchar":                         gdb.LocalTypeString,
	"varchar2":                        gdb.LocalTypeString,
	"variant":                         gdb.LocalTypeString,
	"year":                            gdb.LocalTypeString,
}

// CheckLocalTypeForField checks and returns corresponding local golang type for given db type.
// The parameter `fieldType` is the type name reported by the driver, like `Int64`,
// `Point`, `IntervalDay` or `Array(Nullable(Int64))`.
//
// `bit` is left to the core, as its local type depends on its precision rather than its
// name: the core takes bit(1) as a boolean.
func (d *Driver) CheckLocalTypeForField(ctx context.Context, fieldType string, fieldValue any) (gdb.LocalType, error) {
	var typeName string
	match, _ := gregex.MatchString(`(.+?)\((.+)\)`, fieldType)
	if len(match) == 3 {
		typeName = gstr.Trim(match[1])
	} else {
		typeName = fieldType
	}
	typeName = strings.ToLower(typeName)
	if typeName == "bit" {
		return d.Core.CheckLocalTypeForField(ctx, fieldType, fieldValue)
	}
	if localType, ok := localTypeMap[typeName]; ok {
		return localType, nil
	}
	return d.Core.CheckLocalTypeForField(ctx, fieldType, fieldValue)
}
