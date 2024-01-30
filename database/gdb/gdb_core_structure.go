// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"database/sql/driver"
	"reflect"
	"strings"
	"time"

	"github.com/gogf/gf/v2/encoding/gbinary"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
)

// GetFieldTypeStr retrieves and returns the field type string for certain field by name.
func (c *Core) GetFieldTypeStr(ctx context.Context, fieldName, table, schema string) string {
	field := c.GetFieldType(ctx, fieldName, table, schema)
	if field != nil {
		return field.Type
	}
	return ""
}

// GetFieldType retrieves and returns the field type object for certain field by name.
func (c *Core) GetFieldType(ctx context.Context, fieldName, table, schema string) *TableField {
	fieldsMap, err := c.db.TableFields(ctx, table, schema)
	if err != nil {
		intlog.Errorf(
			ctx,
			`TableFields failed for table "%s", schema "%s": %+v`,
			table, schema, err,
		)
		return nil
	}
	for tableFieldName, tableField := range fieldsMap {
		if tableFieldName == fieldName {
			return tableField
		}
	}
	return nil
}

// ConvertDataForRecord is a very important function, which does converting for any data that
// will be inserted into table/collection as a record.
//
// The parameter `value` should be type of *map/map/*struct/struct.
// It supports embedded struct definition for struct.
func (c *Core) ConvertDataForRecord(ctx context.Context, value interface{}, table string) (map[string]interface{}, error) {
	var (
		err  error
		data = MapOrStructToMapDeep(value, true)
	)
	for fieldName, fieldValue := range data {
		data[fieldName], err = c.db.ConvertValueForField(
			ctx,
			c.GetFieldTypeStr(ctx, fieldName, table, c.GetSchema()),
			fieldValue,
		)
		if err != nil {
			return nil, gerror.Wrapf(err, `ConvertDataForRecord failed for value: %#v`, fieldValue)
		}
	}
	return data, nil
}

// ConvertValueForField converts value to the type of the record field.
// The parameter `fieldType` is the target record field.
// The parameter `fieldValue` is the value that to be committed to record field.
func (c *Core) ConvertValueForField(ctx context.Context, fieldType string, fieldValue interface{}) (interface{}, error) {
	var (
		err            error
		convertedValue = fieldValue
	)
	// If `value` implements interface `driver.Valuer`, it then uses the interface for value converting.
	if valuer, ok := fieldValue.(driver.Valuer); ok {
		if convertedValue, err = valuer.Value(); err != nil {
			if err != nil {
				return nil, err
			}
		}
		return convertedValue, nil
	}
	// Default value converting.
	var (
		rvValue = reflect.ValueOf(fieldValue)
		rvKind  = rvValue.Kind()
	)
	for rvKind == reflect.Ptr {
		rvValue = rvValue.Elem()
		rvKind = rvValue.Kind()
	}
	switch rvKind {
	case reflect.Slice, reflect.Array, reflect.Map:
		// It should ignore the bytes type.
		if _, ok := fieldValue.([]byte); !ok {
			// Convert the value to JSON.
			convertedValue, err = json.Marshal(fieldValue)
			if err != nil {
				return nil, err
			}
		}

	case reflect.Struct:
		switch r := fieldValue.(type) {
		// If the time is zero, it then updates it to nil,
		// which will insert/update the value to database as "null".
		case time.Time:
			if r.IsZero() {
				convertedValue = nil
			}

		case gtime.Time:
			if r.IsZero() {
				convertedValue = nil
			} else {
				convertedValue = r.Time
			}

		case *gtime.Time:
			if r.IsZero() {
				convertedValue = nil
			} else {
				convertedValue = r.Time
			}

		case *time.Time:
			// Nothing to do.

		case Counter, *Counter:
			// Nothing to do.

		default:
			// If `value` implements interface iNil,
			// check its IsNil() function, if got ture,
			// which will insert/update the value to database as "null".
			if v, ok := fieldValue.(iNil); ok && v.IsNil() {
				convertedValue = nil
			} else if s, ok := fieldValue.(iString); ok {
				// Use string conversion in default.
				convertedValue = s.String()
			} else {
				// Convert the value to JSON.
				convertedValue, err = json.Marshal(fieldValue)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return convertedValue, nil
}

// CheckLocalTypeForField checks and returns corresponding type for given db type.
func (c *Core) CheckLocalTypeForField(ctx context.Context, fieldType string, fieldValue interface{}) (LocalType, error) {
	var (
		typeName    string
		typePattern string
	)
	match, _ := gregex.MatchString(`(.+?)\((.+)\)`, fieldType)
	if len(match) == 3 {
		typeName = gstr.Trim(match[1])
		typePattern = gstr.Trim(match[2])
	} else {
		typeName = gstr.Split(fieldType, " ")[0]
	}

	typeName = strings.ToLower(typeName)

	switch typeName {
	case
		fieldTypeBinary,
		fieldTypeVarbinary,
		fieldTypeBlob,
		fieldTypeTinyblob,
		fieldTypeMediumblob,
		fieldTypeLongblob:
		return LocalTypeBytes, nil

	case
		fieldTypeInt,
		fieldTypeTinyint,
		fieldTypeSmallInt,
		fieldTypeSmallint,
		fieldTypeMediumInt,
		fieldTypeMediumint,
		fieldTypeSerial:
		if gstr.ContainsI(fieldType, "unsigned") {
			return LocalTypeUint, nil
		}
		return LocalTypeInt, nil

	case
		fieldTypeBigInt,
		fieldTypeBigint,
		fieldTypeBigserial:
		if gstr.ContainsI(fieldType, "unsigned") {
			return LocalTypeUint64, nil
		}
		return LocalTypeInt64, nil

	case
		fieldTypeReal:
		return LocalTypeFloat32, nil

	case
		fieldTypeDecimal,
		fieldTypeMoney,
		fieldTypeNumeric,
		fieldTypeSmallmoney:
		return LocalTypeString, nil
	case
		fieldTypeFloat,
		fieldTypeDouble:
		return LocalTypeFloat64, nil

	case
		fieldTypeBit:
		// It is suggested using bit(1) as boolean.
		if typePattern == "1" {
			return LocalTypeBool, nil
		}
		s := gconv.String(fieldValue)
		// mssql is true|false string.
		if strings.EqualFold(s, "true") || strings.EqualFold(s, "false") {
			return LocalTypeBool, nil
		}
		if gstr.ContainsI(fieldType, "unsigned") {
			return LocalTypeUint64Bytes, nil
		}
		return LocalTypeInt64Bytes, nil

	case
		fieldTypeBool:
		return LocalTypeBool, nil

	case
		fieldTypeDate:
		return LocalTypeDate, nil

	case
		fieldTypeDatetime,
		fieldTypeTimestamp,
		fieldTypeTimestampz:
		return LocalTypeDatetime, nil

	case
		fieldTypeJson:
		return LocalTypeJson, nil

	case
		fieldTypeJsonb:
		return LocalTypeJsonb, nil

	default:
		// Auto-detect field type, using key match.
		switch {
		case strings.Contains(typeName, "text") || strings.Contains(typeName, "char") || strings.Contains(typeName, "character"):
			return LocalTypeString, nil

		case strings.Contains(typeName, "float") || strings.Contains(typeName, "double") || strings.Contains(typeName, "numeric"):
			return LocalTypeFloat64, nil

		case strings.Contains(typeName, "bool"):
			return LocalTypeBool, nil

		case strings.Contains(typeName, "binary") || strings.Contains(typeName, "blob"):
			return LocalTypeBytes, nil

		case strings.Contains(typeName, "int"):
			if gstr.ContainsI(fieldType, "unsigned") {
				return LocalTypeUint, nil
			}
			return LocalTypeInt, nil

		case strings.Contains(typeName, "time"):
			return LocalTypeDatetime, nil

		case strings.Contains(typeName, "date"):
			return LocalTypeDatetime, nil

		default:
			return LocalTypeString, nil
		}
	}
}

// ConvertValueForLocal converts value to local Golang type of value according field type name from database.
// The parameter `fieldType` is in lower case, like:
// `float(5,2)`, `unsigned double(5,2)`, `decimal(10,2)`, `char(45)`, `varchar(100)`, etc.
func (c *Core) ConvertValueForLocal(
	ctx context.Context, fieldType string, fieldValue interface{},
) (interface{}, error) {
	// If there's no type retrieved, it returns the `fieldValue` directly
	// to use its original data type, as `fieldValue` is type of interface{}.
	if fieldType == "" {
		return fieldValue, nil
	}
	typeName, err := c.db.CheckLocalTypeForField(ctx, fieldType, fieldValue)
	if err != nil {
		return nil, err
	}
	switch typeName {
	case LocalTypeBytes:
		var typeNameStr = string(typeName)
		if strings.Contains(typeNameStr, "binary") || strings.Contains(typeNameStr, "blob") {
			return fieldValue, nil
		}
		return gconv.Bytes(fieldValue), nil

	case LocalTypeInt:
		return gconv.Int(gconv.String(fieldValue)), nil

	case LocalTypeUint:
		return gconv.Uint(gconv.String(fieldValue)), nil

	case LocalTypeInt64:
		return gconv.Int64(gconv.String(fieldValue)), nil

	case LocalTypeUint64:
		return gconv.Uint64(gconv.String(fieldValue)), nil

	case LocalTypeInt64Bytes:
		return gbinary.BeDecodeToInt64(gconv.Bytes(fieldValue)), nil

	case LocalTypeUint64Bytes:
		return gbinary.BeDecodeToUint64(gconv.Bytes(fieldValue)), nil

	case LocalTypeFloat32:
		return gconv.Float32(gconv.String(fieldValue)), nil

	case LocalTypeFloat64:
		return gconv.Float64(gconv.String(fieldValue)), nil

	case LocalTypeBool:
		s := gconv.String(fieldValue)
		// mssql is true|false string.
		if strings.EqualFold(s, "true") {
			return 1, nil
		}
		if strings.EqualFold(s, "false") {
			return 0, nil
		}
		return gconv.Bool(fieldValue), nil

	case LocalTypeDate:
		// Date without time.
		if t, ok := fieldValue.(time.Time); ok {
			return gtime.NewFromTime(t).Format("Y-m-d"), nil
		}
		t, _ := gtime.StrToTime(gconv.String(fieldValue))
		return t.Format("Y-m-d"), nil

	case LocalTypeDatetime:
		if t, ok := fieldValue.(time.Time); ok {
			return gtime.NewFromTime(t), nil
		}
		t, _ := gtime.StrToTime(gconv.String(fieldValue))
		return t, nil

	default:
		return gconv.String(fieldValue), nil
	}
}

// mappingAndFilterData automatically mappings the map key to table field and removes
// all key-value pairs that are not the field of given table.
func (c *Core) mappingAndFilterData(ctx context.Context, schema, table string, data map[string]interface{}, filter bool) (map[string]interface{}, error) {
	fieldsMap, err := c.db.TableFields(ctx, c.guessPrimaryTableName(table), schema)
	if err != nil {
		return nil, err
	}
	fieldsKeyMap := make(map[string]interface{}, len(fieldsMap))
	for k := range fieldsMap {
		fieldsKeyMap[k] = nil
	}
	// Automatic data key to table field name mapping.
	var foundKey string
	for dataKey, dataValue := range data {
		if _, ok := fieldsKeyMap[dataKey]; !ok {
			foundKey, _ = gutil.MapPossibleItemByKey(fieldsKeyMap, dataKey)
			if foundKey != "" {
				if _, ok = data[foundKey]; !ok {
					data[foundKey] = dataValue
				}
				delete(data, dataKey)
			}
		}
	}
	// Data filtering.
	// It deletes all key-value pairs that has incorrect field name.
	if filter {
		for dataKey := range data {
			if _, ok := fieldsMap[dataKey]; !ok {
				delete(data, dataKey)
			}
		}
	}
	return data, nil
}
