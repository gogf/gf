// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mssql

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"
	mssqldriver "github.com/microsoft/go-mssqldb"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gregex"
)

// characterFieldTypes are the SQL Server types storing text, into which a value is bound as a string.
var characterFieldTypes = map[string]bool{
	"char":     true,
	"varchar":  true,
	"nchar":    true,
	"nvarchar": true,
	"text":     true,
	"ntext":    true,
	"xml":      true,
}

// binaryFieldTypes are the SQL Server types storing bytes, into which NULL is bound as binary.
var binaryFieldTypes = map[string]bool{
	"binary":    true,
	"varbinary": true,
	"image":     true,
}

// dateTimeOffsetLayout is the layout of a datetimeoffset value with all its fractional digits,
// which SQL Server converts from a string whatever its DATEFORMAT setting.
const dateTimeOffsetLayout = "2006-01-02 15:04:05.9999999 -07:00"

// zonelessTimeFieldTypes are the SQL Server types storing a date and time without time zone, into
// which a time is written as its wall clock in the time zone of the driver.
var zonelessTimeFieldTypes = map[string]bool{
	"datetime":      true,
	"datetime2":     true,
	"smalldatetime": true,
}

// CheckLocalTypeForField checks and returns corresponding local Golang type for given db type.
//
// SQL Server type mapping (only types not handled by Core are listed):
//
//	| SQL Server Type   | Local Go Type |
//	|-------------------|---------------|
//	| uniqueidentifier  | uuid.UUID     |
//	| bit               | bool          |
func (d *Driver) CheckLocalTypeForField(ctx context.Context, fieldType string, fieldValue any) (gdb.LocalType, error) {
	typeName, _ := gregex.ReplaceString(`\(.+\)`, "", fieldType)
	typeName = strings.ToLower(typeName)

	switch typeName {
	case "uniqueidentifier":
		return gdb.LocalTypeUUID, nil

	case "bit":
		return gdb.LocalTypeBool, nil

	default:
		return d.Core.CheckLocalTypeForField(ctx, fieldType, fieldValue)
	}
}

// ConvertValueForField converts value to the type of the record field `fieldType`.
// A value that Core encodes as JSON bytes is bound as a string into a text column, as SQL Server
// converts bytes to text by reinterpreting them; NULL is bound as binary into a binary column, as
// SQL Server refuses to convert an untyped NULL parameter to binary; a time is converted to the
// time zone of the driver for a date and time column without time zone, as SQL Server stores the
// wall clock of the time in its own time zone, and is bound as a string with its offset for a
// datetimeoffset column, so that it keeps its own time zone, which DoFilter does not convert then.
// A date or a time column gets the date or time of day of the time itself, as Core writes it.
func (d *Driver) ConvertValueForField(ctx context.Context, fieldType string, fieldValue any) (any, error) {
	typeName := strings.ToLower(fieldType)
	if zonelessTimeFieldTypes[typeName] {
		fieldValue, _ = d.timeInLocation(fieldValue)
	}
	convertedValue, err := d.Core.ConvertValueForField(ctx, fieldType, fieldValue)
	if err != nil {
		return nil, err
	}
	if b, ok := convertedValue.([]byte); ok && characterFieldTypes[typeName] {
		return string(b), nil
	}
	if convertedValue == nil && binaryFieldTypes[typeName] {
		return []byte(nil), nil
	}
	if typeName == "datetimeoffset" {
		switch v := convertedValue.(type) {
		case time.Time:
			return v.Format(dateTimeOffsetLayout), nil
		case *time.Time:
			if v != nil {
				return v.Format(dateTimeOffsetLayout), nil
			}
		}
	}
	return convertedValue, nil
}

// ConvertValueForLocal converts value to local Golang type of value according field type name from database.
//
// SQL Server stores UNIQUEIDENTIFIER on the wire as a 16-byte binary blob whose first 8 bytes
// follow the little-endian COM/Win32 GUID layout, while the remaining 8 bytes are big-endian.
// Reading the raw bytes as a string yields garbage; go-mssqldb's UniqueIdentifier.Scan
// performs the required byte-order swap so the returned [uuid.UUID] matches the canonical
// RFC 4122 form (and what tools like SSMS display).
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

// timeInLocation returns `value` converted to the time zone of the driver, the local time zone if
// it has none, and true, if `value` is a time, a gtime or a valid sql.NullTime in another time
// zone, a pointer to one, or a sql.NamedArg of any of them, or else `value` itself and false.
// Other driver.Valuer values are left as they are, as the types of go-mssqldb are among them.
func (d *Driver) timeInLocation(value any) (any, bool) {
	location := d.location
	if location == nil {
		location = time.Local
	}
	switch v := value.(type) {
	case time.Time:
		if v.Location() != location {
			return v.In(location), true
		}
	case *time.Time:
		if v != nil && v.Location() != location {
			t := v.In(location)
			return &t, true
		}
	case gtime.Time:
		if v.Location() != location {
			return v.In(location), true
		}
	case *gtime.Time:
		if v != nil && v.Location() != location {
			return v.In(location), true
		}
	case sql.NullTime:
		if v.Valid && v.Time.Location() != location {
			v.Time = v.Time.In(location)
			return v, true
		}
	case *sql.NullTime:
		if v != nil && v.Valid && v.Time.Location() != location {
			t := sql.NullTime{Time: v.Time.In(location), Valid: true}
			return &t, true
		}
	case sql.NamedArg:
		if converted, ok := d.timeInLocation(v.Value); ok {
			v.Value = converted
			return v, true
		}
	}
	return value, false
}
