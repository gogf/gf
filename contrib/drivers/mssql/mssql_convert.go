// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mssql

import (
	"context"
	"strings"

	"github.com/google/uuid"
	mssqldriver "github.com/microsoft/go-mssqldb"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gregex"
)

// CheckLocalTypeForField checks and returns corresponding local Golang type for given db type.
//
// SQL Server type mapping (only types not handled by Core are listed):
//
//	| SQL Server Type   | Local Go Type |
//	|-------------------|---------------|
//	| uniqueidentifier  | uuid.UUID     |
func (d *Driver) CheckLocalTypeForField(ctx context.Context, fieldType string, fieldValue any) (gdb.LocalType, error) {
	typeName, _ := gregex.ReplaceString(`\(.+\)`, "", fieldType)
	typeName = strings.ToLower(typeName)

	switch typeName {
	case "uniqueidentifier":
		return gdb.LocalTypeUUID, nil

	default:
		return d.Core.CheckLocalTypeForField(ctx, fieldType, fieldValue)
	}
}

// ConvertValueForLocal converts value to local Golang type of value according field type name from database.
//
// SQL Server stores UNIQUEIDENTIFIER on the wire as a 16-byte binary blob whose first 8 bytes
// follow the little-endian COM/Win32 GUID layout, while the remaining 8 bytes are big-endian.
// Reading the raw bytes as a string yields garbage; [mssql.UniqueIdentifier.Scan] performs the
// required byte-order swap so the returned [uuid.UUID] matches the canonical RFC 4122 form
// (and what tools like SSMS display).
func (d *Driver) ConvertValueForLocal(ctx context.Context, fieldType string, fieldValue any) (any, error) {
	typeName, _ := gregex.ReplaceString(`\(.+\)`, "", fieldType)
	typeName = strings.ToLower(typeName)

	switch typeName {
	case "uniqueidentifier":
		var msID mssqldriver.UniqueIdentifier
		if err := msID.Scan(fieldValue); err != nil {
			return nil, gerror.Wrapf(
				err, "convert uniqueidentifier value %v to uuid.UUID failed", fieldValue,
			)
		}
		return uuid.UUID(msID), nil

	default:
		return d.Core.ConvertValueForLocal(ctx, fieldType, fieldValue)
	}
}
