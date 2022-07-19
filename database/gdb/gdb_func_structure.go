// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

const (
	typeString      = "string"
	typeDate        = "date"
	typeDatetime    = "datetime"
	typeInt         = "int"
	typeUint        = "uint"
	typeInt64       = "int64"
	typeUint64      = "uint64"
	typeInt64Slice  = "[]int64"
	typeUint64Slice = "[]uint64"
	typeInt64Bytes  = "int64-bytes"
	typeUint64Bytes = "uint64-bytes"
	typeFloat32     = "float32"
	typeFloat64     = "float64"
	typeBytes       = "[]byte"
	typeBool        = "bool"
)

func CheckValueForLocalType(ctx context.Context, fieldType string, fieldValue interface{}) (string, error) {
	var (
		typeName    string
		typePattern string
	)
	match, _ := gregex.MatchString(`(.+?)\((.+)\)`, fieldType)
	if len(match) == 3 {
		typeName = gstr.Trim(match[1])
		typePattern = gstr.Trim(match[2])
	} else {
		typeName = fieldType
	}
	typeName = strings.ToLower(typeName)
	switch typeName {
	case
		"binary",
		"varbinary",
		"blob",
		"tinyblob",
		"mediumblob",
		"longblob":
		return typeBytes, nil

	case
		"int",
		"tinyint",
		"small_int",
		"smallint",
		"medium_int",
		"mediumint",
		"serial":
		if typePattern == "1" {
			return typeBool, nil
		}
		if gstr.ContainsI(fieldType, "unsigned") {
			return typeUint, nil
		}
		return typeInt, nil

	case "_int4", "_int8":
		if gstr.ContainsI(fieldType, "unsigned") {
			return typeUint64Slice, nil
		}
		return typeInt64Slice, nil

	case
		"big_int",
		"bigint",
		"bigserial":
		if gstr.ContainsI(fieldType, "unsigned") {
			return typeUint64, nil
		}
		return typeInt64, nil

	case "real":
		return typeFloat32, nil

	case
		"float",
		"double",
		"decimal",
		"money",
		"numeric",
		"smallmoney":
		return typeFloat64, nil

	case "bit":
		if typePattern == "1" {
			return typeBool, nil
		}
		s := gconv.String(fieldValue)
		// mssql is true|false string.
		if strings.EqualFold(s, "true") || strings.EqualFold(s, "false") {
			return typeBool, nil
		}
		if gstr.ContainsI(fieldType, "unsigned") {
			return typeUint64Bytes, nil
		}
		return typeInt64Bytes, nil

	case "bool":
		return typeBool, nil

	case "date":
		return typeDate, nil

	case
		"datetime",
		"timestamp",
		"timestamptz":
		return typeDatetime, nil

	default:
		// Auto-detect field type, using key match.
		switch {
		case strings.Contains(typeName, "text") || strings.Contains(typeName, "char") || strings.Contains(typeName, "character"):
			return typeString, nil

		case strings.Contains(typeName, "float") || strings.Contains(typeName, "double") || strings.Contains(typeName, "numeric"):
			return typeFloat64, nil

		case strings.Contains(typeName, "bool"):
			return typeBool, nil

		case strings.Contains(typeName, "binary") || strings.Contains(typeName, "blob"):
			return typeBytes, nil

		case strings.Contains(typeName, "int"):
			return typeInt, nil

		case strings.Contains(typeName, "time"):
			return typeDatetime, nil

		case strings.Contains(typeName, "date"):
			return typeDatetime, nil

		default:
			return typeString, nil
		}
	}
}
