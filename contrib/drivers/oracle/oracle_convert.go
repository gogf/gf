// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package oracle

import (
	"context"
	"database/sql"
	"reflect"
	"strconv"
	"strings"
	"time"

	gora "github.com/sijms/go-ora/v2"

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
	"date":             gdb.LocalTypeDatetime,
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

// maxBindSize is the size in bytes beyond which a character or binary value can no longer be
// bound as VARCHAR2 or RAW, but has to be bound as a LOB.
const maxBindSize = 4000

// characterFieldTypes are the field types whose values are bound as text.
var characterFieldTypes = map[string]bool{
	"CHAR":      true,
	"NCHAR":     true,
	"VARCHAR":   true,
	"VARCHAR2":  true,
	"NVARCHAR2": true,
	"LONG":      true,
	"CLOB":      true,
	"NCLOB":     true,
}

// zonelessTimeTypes are the type names the driver reports for the time values that carry no
// time zone, which the driver returns labeled as UTC.
var zonelessTimeTypes = map[string]bool{
	"DATE":         true,
	"TIMESTAMP":    true,
	"TimeStampDTY": true,
}

// ConvertValueForField converts value to the type of the record field.
// The parameter `fieldType` is the field type of the table, like `VARCHAR2`, `TIMESTAMP` or `CLOB`.
// The parameter `fieldValue` is the value that to be committed to record field.
func (d *Driver) ConvertValueForField(ctx context.Context, fieldType string, fieldValue any) (any, error) {
	convertedValue, err := d.Core.ConvertValueForField(ctx, fieldType, fieldValue)
	if err != nil {
		return nil, err
	}
	convertedValue = convertNamedString(convertedValue)
	if b, ok := convertedValue.([]byte); ok && characterFieldTypes[fieldType] {
		convertedValue = string(b)
	}
	switch fieldType {
	case "CLOB":
		if s, ok := convertedValue.(string); ok && len(s) > maxBindSize {
			return gora.Clob{String: s, Valid: true}, nil
		}

	case "NCLOB":
		if s, ok := convertedValue.(string); ok && len(s) > maxBindSize {
			return gora.NClob{String: s, Valid: true}, nil
		}

	case "BLOB":
		if b, ok := convertedValue.([]byte); ok && len(b) > maxBindSize {
			return gora.Blob{Data: b, Valid: true}, nil
		}

	case "DATE":
		switch t := convertedValue.(type) {
		case time.Time:
			return t.In(time.Local), nil
		case *time.Time:
			return t.In(time.Local), nil
		}

	case "TIMESTAMP":
		switch t := convertedValue.(type) {
		case time.Time:
			return gora.TimeStampTZ(t.In(time.Local)), nil
		case *time.Time:
			return gora.TimeStampTZ(t.In(time.Local)), nil
		}

	case "BINARY_DOUBLE", "BINARY_FLOAT":
		if s, ok := formatBinaryFloat(convertedValue); ok {
			return s, nil
		}
	}
	return convertedValue, nil
}

// ConvertValueForLocal converts value to local Golang type of value according field type name from database.
// The parameter `fieldType` is the type name reported by the driver, like `NUMBER`,
// `TimeStampDTY` or `TimeStampTZ_DTY`.
func (d *Driver) ConvertValueForLocal(ctx context.Context, fieldType string, fieldValue any) (any, error) {
	if t, ok := fieldValue.(time.Time); ok && zonelessTimeTypes[fieldType] {
		fieldValue = time.Date(
			t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.Local,
		)
	}
	return d.Core.ConvertValueForLocal(ctx, fieldType, fieldValue)
}

// ConvertColumnValueForLocal converts a value scanned from a result set column to local Golang type.
// The underlying driver scans a NUMBER value as an int64 or uint64 for an integer and as a string
// for a decimal, but reports float64 as the scan type. The integer is kept as it is, and the
// decimal is converted to float64 only if float64 holds it exactly, or else kept as the string.
func (d *Driver) ConvertColumnValueForLocal(ctx context.Context, columnType *sql.ColumnType, fieldValue any) (any, error) {
	if columnType.DatabaseTypeName() != "NUMBER" {
		return d.Core.ConvertColumnValueForLocal(ctx, columnType, fieldValue)
	}
	if s, ok := fieldValue.(string); ok {
		if f, err := strconv.ParseFloat(s, 64); err == nil && strconv.FormatFloat(f, 'f', -1, 64) == s {
			return f, nil
		}
	}
	return fieldValue, nil
}

// convertNamedString returns the value of a string type defined outside the underlying driver
// as a string, as the underlying driver only binds the string types it defines.
func convertNamedString(value any) any {
	switch value.(type) {
	case string, gdb.Raw, *gdb.Raw, gora.NVarChar, *gora.NVarChar:
		return value
	}
	reflectValue := reflect.ValueOf(value)
	for reflectValue.Kind() == reflect.Pointer && !reflectValue.IsNil() {
		reflectValue = reflectValue.Elem()
	}
	if reflectValue.Kind() == reflect.String {
		return reflectValue.String()
	}
	return value
}

// formatBinaryFloat formats a float value as a numeric literal without decimal separator, like
// `17976931348623157e292`, which Oracle converts to BINARY_DOUBLE or BINARY_FLOAT without loss
// whatever the NLS settings of the session. A float bound as NUMBER does not keep a value
// beyond the exponent range of NUMBER.
func formatBinaryFloat(value any) (string, bool) {
	reflectValue := reflect.ValueOf(value)
	for reflectValue.Kind() == reflect.Pointer && !reflectValue.IsNil() {
		reflectValue = reflectValue.Elem()
	}
	var bitSize int
	switch reflectValue.Kind() {
	case reflect.Float32:
		bitSize = 32
	case reflect.Float64:
		bitSize = 64
	default:
		return "", false
	}
	text := strconv.FormatFloat(reflectValue.Float(), 'e', -1, bitSize)
	mantissa, exponentText, ok := strings.Cut(text, "e")
	if !ok {
		return text, true
	}
	exponent, _ := strconv.Atoi(exponentText)
	if integer, fraction, ok := strings.Cut(mantissa, "."); ok {
		mantissa = integer + fraction
		exponent -= len(fraction)
	}
	return mantissa + "e" + strconv.Itoa(exponent), true
}
