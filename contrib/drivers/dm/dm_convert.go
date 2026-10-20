// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"math/big"
	"strconv"
	"strings"
	"time"

	"gitee.com/chunanyong/dm"
)

// characterFieldTypes are the field types whose values are bound as text. The underlying driver
// binds []byte as binary data, which DM converts to its hexadecimal text for these types.
var characterFieldTypes = map[string]bool{
	"CHAR":        true,
	"CHARACTER":   true,
	"VARCHAR":     true,
	"VARCHAR2":    true,
	"NCHAR":       true,
	"NVARCHAR":    true,
	"NVARCHAR2":   true,
	"TEXT":        true,
	"LONG":        true,
	"LONGVARCHAR": true,
	"CLOB":        true,
	"NCLOB":       true,
}

// binaryFieldTypes are the field types whose values are bound as binary data. The underlying
// driver binds a string as text, which DM converts to binary data as hexadecimal digits.
var binaryFieldTypes = map[string]bool{
	"BINARY":        true,
	"VARBINARY":     true,
	"RAW":           true,
	"BLOB":          true,
	"IMAGE":         true,
	"LONGVARBINARY": true,
}

// exactNumericFieldTypes are the field types of exact numbers. The underlying driver binds text
// for them as a binary float with a 64-bit mantissa, which keeps about 20 significant digits.
var exactNumericFieldTypes = map[string]bool{
	"DECIMAL": true,
	"DEC":     true,
	"NUMERIC": true,
	"NUMBER":  true,
}

// ConvertValueForField converts value to the type of the record field.
// The parameter `fieldType` is the field type of the table, like `VARCHAR`, `TIMESTAMP` or `CLOB`.
// The parameter `fieldValue` is the value that to be committed to record field.
//
// Text is written to a character column as text, to a binary column as its bytes, and to an exact
// numeric column as the decimal number it is with all its digits. A time is written:
//   - to a column without a time zone as its wall clock in the configured Timezone or else the
//     local time zone, see sessionTimes; a DATE or a TIME column keeps its date or its time of day;
//   - to a column with a time zone as its instant, which the column keeps in the offset of the
//     session;
//   - to a column with the local time zone as the text of its wall clock in the time zone of the
//     session, see localTimeFieldTypes;
//   - to a character column as its text.
//
// A zero time is written as NULL.
func (d *Driver) ConvertValueForField(ctx context.Context, fieldType string, fieldValue any) (any, error) {
	if t, ok := timeOfValue(fieldValue); ok {
		if value, ok := convertTimeForField(fieldType, t); ok {
			return value, nil
		}
	}
	convertedValue, err := d.Core.ConvertValueForField(ctx, fieldType, fieldValue)
	if err != nil {
		return nil, err
	}
	switch value := convertedValue.(type) {
	case []byte:
		if characterFieldTypes[strings.ToUpper(fieldType)] {
			return string(value), nil
		}
	case string:
		if binaryFieldTypes[strings.ToUpper(fieldType)] {
			return []byte(value), nil
		}
		if exactNumericFieldTypes[strings.ToUpper(fieldType)] {
			if decimal, ok := parseDecimal(value); ok {
				return decimal, nil
			}
		}
	case *time.Time:
		if value.IsZero() {
			return nil, nil
		}
	}
	return convertedValue, nil
}

// ConvertColumnValueForLocal converts a value scanned from a result set column to local Golang type.
// The underlying driver scans the value of an exact numeric column, being DECIMAL, DEC, NUMERIC or
// NUMBER, as its decimal text without the zero before the decimal point and without trailing
// zeros, like ".5" for 0.50 of DECIMAL(10,2). That text is formatted with formatDecimal to the
// scale of the column.
//
// The value of a NUMBER column without fractional digits is an integer, as on Oracle: int64, or
// uint64 beyond the range of int64, or else the text. The result set reports NUMBER without
// precision with a scale of 0 too, so a value with fractional digits is kept as its text.
//
// The value of a DATETIME or a TIMESTAMP column, which the underlying driver returns in the fixed
// time zone of its session, is relabelled as its wall clock in the configured Timezone or else the
// local time zone, as ConvertValueForField writes it to such a column. The wall clock of the value
// is kept as it is if relabelling it would not: a wall clock that a daylight saving time skips, in
// the time zone of the session, has no equal in the other time zone. The core formats the value of
// a DATE or a TIME column as the text of its wall clock.
func (d *Driver) ConvertColumnValueForLocal(ctx context.Context, columnType *sql.ColumnType, fieldValue any) (any, error) {
	if t, ok := fieldValue.(time.Time); ok {
		switch columnType.DatabaseTypeName() {
		case "DATETIME", "TIMESTAMP":
			fieldValue = relabelWallClock(t, d.wallClockLocation())
		}
	}
	if _, scale, ok := columnType.DecimalSize(); ok {
		if text, ok := fieldValue.(string); ok {
			if scale == 0 && columnType.DatabaseTypeName() == "NUMBER" {
				if i, err := strconv.ParseInt(text, 10, 64); err == nil {
					return i, nil
				}
				if u, err := strconv.ParseUint(text, 10, 64); err == nil {
					return u, nil
				}
			}
			fieldValue = formatDecimal(text, scale)
		}
	}
	return d.Core.ConvertColumnValueForLocal(ctx, columnType, fieldValue)
}

// formatDecimal returns the decimal text `text` with a zero before its decimal point if it has
// no integer digit, and with at least `scale` fractional digits, like "0.50" for ".5" and
// "-1500.00" for "-1500" with a scale of 2.
func formatDecimal(text string, scale int64) string {
	var sign string
	if strings.HasPrefix(text, "-") {
		sign, text = "-", text[1:]
	}
	integer, fraction, _ := strings.Cut(text, ".")
	if integer == "" {
		integer = "0"
	}
	for int64(len(fraction)) < scale {
		fraction += "0"
	}
	if fraction == "" {
		return sign + integer
	}
	return sign + integer + "." + fraction
}

// decimalText is the text of a decimal number that binds the number with all its digits. It is a
// string, which the core formats as the text in the sql it logs.
type decimalText string

// Value returns the decimal number for the underlying driver to bind.
func (t decimalText) Value() (driver.Value, error) {
	text := strings.TrimSpace(string(t))
	value, _, err := big.ParseFloat(text, 10, uint(len(text))*4+64, big.ToNearestEven)
	if err != nil {
		return nil, err
	}
	decimal, err := dm.NewDecimalFromBigFloat(value)
	if err != nil {
		return nil, err
	}
	return *decimal, nil
}

// parseDecimal returns `text` as a decimalText, or false if it is not the text of a finite
// decimal number.
func parseDecimal(text string) (decimalText, bool) {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return "", false
	}
	value, _, err := big.ParseFloat(trimmed, 10, 64, big.ToNearestEven)
	if err != nil || value.IsInf() {
		return "", false
	}
	return decimalText(text), true
}
