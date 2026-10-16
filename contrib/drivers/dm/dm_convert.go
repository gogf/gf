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
// numeric column as the decimal number it is with all its digits, as MySQL writes them. A time is passed to the underlying driver as it is, which writes it to a column
// without time zone as its wall clock in the local time zone, as MySQL does with `loc=Local`.
// A zero time is written as NULL. The field type is passed to the core in lower case, which formats
// a time for a DATE or TIME column by the type name in lower case.
func (d *Driver) ConvertValueForField(ctx context.Context, fieldType string, fieldValue any) (any, error) {
	convertedValue, err := d.Core.ConvertValueForField(ctx, strings.ToLower(fieldType), fieldValue)
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
// scale of the column, as MySQL returns it.
func (d *Driver) ConvertColumnValueForLocal(ctx context.Context, columnType *sql.ColumnType, fieldValue any) (any, error) {
	if _, scale, ok := columnType.DecimalSize(); ok {
		if text, ok := fieldValue.(string); ok {
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
