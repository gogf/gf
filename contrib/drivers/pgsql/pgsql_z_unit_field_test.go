// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package pgsql_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

// Test_TableFields tests the TableFields method for retrieving table field information
func Test_TableFields(t *testing.T) {
	table := createAllTypesTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		fields, err := db.TableFields(ctx, table)
		t.AssertNil(err)
		t.Assert(len(fields) > 0, true)

		// Test primary key field
		t.Assert(fields["id"].Name, "id")
		t.Assert(fields["id"].Key, "pri")

		// Test integer types
		t.Assert(fields["col_int2"].Name, "col_int2")
		t.Assert(fields["col_int4"].Name, "col_int4")
		t.Assert(fields["col_int8"].Name, "col_int8")

		// Test float types
		t.Assert(fields["col_float4"].Name, "col_float4")
		t.Assert(fields["col_float8"].Name, "col_float8")
		t.Assert(fields["col_numeric"].Name, "col_numeric")

		// Test character types
		t.Assert(fields["col_char"].Name, "col_char")
		t.Assert(fields["col_varchar"].Name, "col_varchar")
		t.Assert(fields["col_text"].Name, "col_text")

		// Test boolean type
		t.Assert(fields["col_bool"].Name, "col_bool")

		// Test date/time types
		t.Assert(fields["col_date"].Name, "col_date")
		t.Assert(fields["col_timestamp"].Name, "col_timestamp")

		// Test JSON types
		t.Assert(fields["col_json"].Name, "col_json")
		t.Assert(fields["col_jsonb"].Name, "col_jsonb")

		// Test array types
		t.Assert(fields["col_int2_arr"].Name, "col_int2_arr")
		t.Assert(fields["col_int4_arr"].Name, "col_int4_arr")
		t.Assert(fields["col_varchar_arr"].Name, "col_varchar_arr")
	})
}

// Test_TableFields_Types tests field type information
func Test_TableFields_Types(t *testing.T) {
	table := createAllTypesTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		fields, err := db.TableFields(ctx, table)
		t.AssertNil(err)

		// Test integer type names
		t.Assert(fields["col_int2"].Type, "int2")
		t.Assert(fields["col_int4"].Type, "int4")
		t.Assert(fields["col_int8"].Type, "int8")

		// Test float type names
		t.Assert(fields["col_float4"].Type, "float4")
		t.Assert(fields["col_float8"].Type, "float8")
		t.Assert(fields["col_numeric"].Type, "numeric")

		// Test character type names
		t.Assert(fields["col_char"].Type, "bpchar")
		t.Assert(fields["col_varchar"].Type, "varchar")
		t.Assert(fields["col_text"].Type, "text")

		// Test boolean type name
		t.Assert(fields["col_bool"].Type, "bool")

		// Test date/time type names
		t.Assert(fields["col_date"].Type, "date")
		t.Assert(fields["col_timestamp"].Type, "timestamp")
		t.Assert(fields["col_timestamptz"].Type, "timestamptz")

		// Test JSON type names
		t.Assert(fields["col_json"].Type, "json")
		t.Assert(fields["col_jsonb"].Type, "jsonb")

		// Test array type names (PostgreSQL uses _ prefix for array types)
		t.Assert(fields["col_int2_arr"].Type, "_int2")
		t.Assert(fields["col_int4_arr"].Type, "_int4")
		t.Assert(fields["col_int8_arr"].Type, "_int8")
		t.Assert(fields["col_float4_arr"].Type, "_float4")
		t.Assert(fields["col_float8_arr"].Type, "_float8")
		t.Assert(fields["col_numeric_arr"].Type, "_numeric")
		t.Assert(fields["col_varchar_arr"].Type, "_varchar")
		t.Assert(fields["col_text_arr"].Type, "_text")
		t.Assert(fields["col_bool_arr"].Type, "_bool")
	})
}

// Test_TableFields_Nullable tests field nullable information
func Test_TableFields_Nullable(t *testing.T) {
	table := createAllTypesTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		fields, err := db.TableFields(ctx, table)
		t.AssertNil(err)

		// NOT NULL fields should have Null = false
		t.Assert(fields["col_int2"].Null, false)
		t.Assert(fields["col_int4"].Null, false)
		t.Assert(fields["col_numeric"].Null, false)
		t.Assert(fields["col_varchar"].Null, false)
		t.Assert(fields["col_bool"].Null, false)
		t.Assert(fields["col_varchar_arr"].Null, false)

		// Nullable fields should have Null = true
		t.Assert(fields["col_int8"].Null, true)
		t.Assert(fields["col_text"].Null, true)
		t.Assert(fields["col_json"].Null, true)
	})
}

// Test_TableFields_Comments tests field comment information
func Test_TableFields_Comments(t *testing.T) {
	table := createAllTypesTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		fields, err := db.TableFields(ctx, table)
		t.AssertNil(err)

		// Test fields with comments
		t.Assert(fields["id"].Comment, "Primary key ID")
		t.Assert(fields["col_int2"].Comment, "int2 type (smallint)")
		t.Assert(fields["col_int4"].Comment, "int4 type (integer)")
		t.Assert(fields["col_int8"].Comment, "int8 type (bigint)")
		t.Assert(fields["col_numeric"].Comment, "numeric type with precision")
		t.Assert(fields["col_varchar"].Comment, "varchar type")
		t.Assert(fields["col_bool"].Comment, "boolean type")
		t.Assert(fields["col_timestamp"].Comment, "timestamp type")
		t.Assert(fields["col_json"].Comment, "json type")
		t.Assert(fields["col_jsonb"].Comment, "jsonb type")

		// Test array field comments
		t.Assert(fields["col_int2_arr"].Comment, "int2 array type (_int2)")
		t.Assert(fields["col_int4_arr"].Comment, "int4 array type (_int4)")
		t.Assert(fields["col_int8_arr"].Comment, "int8 array type (_int8)")
		t.Assert(fields["col_numeric_arr"].Comment, "numeric array type (_numeric)")
		t.Assert(fields["col_varchar_arr"].Comment, "varchar array type (_varchar)")
		t.Assert(fields["col_text_arr"].Comment, "text array type (_text)")
	})
}

// Test_Field_Type_Conversion tests type conversion for various PostgreSQL types
func Test_Field_Type_Conversion(t *testing.T) {
	table := createInitAllTypesTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Query a single record
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one.IsEmpty(), false)

		// Test integer type conversions
		t.Assert(one["col_int2"].Int(), 1)
		t.Assert(one["col_int4"].Int(), 10)
		t.Assert(one["col_int8"].Int64(), int64(100))

		// Test float type conversions
		t.Assert(one["col_float4"].Float32() > 0, true)
		t.Assert(one["col_float8"].Float64() > 0, true)

		// Test string type conversions
		t.AssertNE(one["col_varchar"].String(), "")
		t.AssertNE(one["col_text"].String(), "")

		// Test boolean type conversion
		t.Assert(one["col_bool"].Bool(), false) // i=1, 1%2==0 is false
	})
}

// Test_Field_Array_Type_Conversion tests array type conversion
func Test_Field_Array_Type_Conversion(t *testing.T) {
	table := createInitAllTypesTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Query a single record
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one.IsEmpty(), false)

		// Test integer array type conversions
		int2Arr := one["col_int2_arr"].Ints()
		t.Assert(len(int2Arr), 3)
		t.Assert(int2Arr[0], 1)
		t.Assert(int2Arr[1], 2)
		t.Assert(int2Arr[2], 1)

		int4Arr := one["col_int4_arr"].Ints()
		t.Assert(len(int4Arr), 3)
		t.Assert(int4Arr[0], 10)
		t.Assert(int4Arr[1], 20)
		t.Assert(int4Arr[2], 1)

		int8Arr := one["col_int8_arr"].Int64s()
		t.Assert(len(int8Arr), 3)
		t.Assert(int8Arr[0], int64(100))
		t.Assert(int8Arr[1], int64(200))
		t.Assert(int8Arr[2], int64(1))

		// Test string array type conversions
		varcharArr := one["col_varchar_arr"].Strings()
		t.Assert(len(varcharArr), 3)
		t.Assert(varcharArr[0], "a")
		t.Assert(varcharArr[1], "b")
		t.Assert(varcharArr[2], "c1")

		textArr := one["col_text_arr"].Strings()
		t.Assert(len(textArr), 3)
		t.Assert(textArr[0], "x")
		t.Assert(textArr[1], "y")
		t.Assert(textArr[2], "z1")

		// Test boolean array type conversions
		// col_bool_arr is '{true, false, %t}' where %t = i%2==0, for i=1 it's false
		boolArr := one["col_bool_arr"].Bools()
		t.Assert(len(boolArr), 3)
		t.Assert(boolArr[0], true)  // literal true
		t.Assert(boolArr[1], false) // literal false
		t.Assert(boolArr[2], false) // i=1, 1%2==0 is false
	})
}

// Test_Field_Array_Insert tests inserting array data
func Test_Field_Array_Insert(t *testing.T) {
	table := createAllTypesTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert with array values
		_, err := db.Model(table).Data(g.Map{
			"col_int2":        1,
			"col_int4":        10,
			"col_numeric":     99.99,
			"col_varchar":     "test",
			"col_bool":        true,
			"col_int2_arr":    []int{1, 2, 3},
			"col_int4_arr":    []int{10, 20, 30},
			"col_varchar_arr": []string{"a", "b", "c"},
		}).Insert()
		t.AssertNil(err)

		// Query and verify
		one, err := db.Model(table).OrderDesc("id").One()
		t.AssertNil(err)

		t.Assert(one["col_int2"].Int(), 1)
		t.Assert(one["col_varchar"].String(), "test")
		t.Assert(one["col_bool"].Bool(), true)

		int2Arr := one["col_int2_arr"].Ints()
		t.Assert(len(int2Arr), 3)
		t.Assert(int2Arr[0], 1)
		t.Assert(int2Arr[1], 2)
		t.Assert(int2Arr[2], 3)

		varcharArr := one["col_varchar_arr"].Strings()
		t.Assert(len(varcharArr), 3)
		t.Assert(varcharArr[0], "a")
		t.Assert(varcharArr[1], "b")
		t.Assert(varcharArr[2], "c")
	})
}

// Test_Field_Array_Update tests updating array data
func Test_Field_Array_Update(t *testing.T) {
	table := createInitAllTypesTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Update array values
		_, err := db.Model(table).Where("id", 1).Data(g.Map{
			"col_int2_arr":    []int{100, 200, 300},
			"col_varchar_arr": []string{"x", "y", "z"},
		}).Update()
		t.AssertNil(err)

		// Query and verify
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)

		int2Arr := one["col_int2_arr"].Ints()
		t.Assert(len(int2Arr), 3)
		t.Assert(int2Arr[0], 100)
		t.Assert(int2Arr[1], 200)
		t.Assert(int2Arr[2], 300)

		varcharArr := one["col_varchar_arr"].Strings()
		t.Assert(len(varcharArr), 3)
		t.Assert(varcharArr[0], "x")
		t.Assert(varcharArr[1], "y")
		t.Assert(varcharArr[2], "z")
	})
}

// Test_Field_JSON_Type tests JSON/JSONB type handling
func Test_Field_JSON_Type(t *testing.T) {
	table := createAllTypesTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert with JSON values
		testData := g.Map{
			"name":  "test",
			"value": 123,
			"items": []string{"a", "b", "c"},
		}
		_, err := db.Model(table).Data(g.Map{
			"col_int2":    1,
			"col_int4":    10,
			"col_numeric": 99.99,
			"col_varchar": "test",
			"col_bool":    true,
			"col_json":    testData,
			"col_jsonb":   testData,
		}).Insert()
		t.AssertNil(err)

		// Query and verify
		one, err := db.Model(table).OrderDesc("id").One()
		t.AssertNil(err)

		// Test JSON field
		jsonMap := one["col_json"].Map()
		t.Assert(jsonMap["name"], "test")
		t.Assert(jsonMap["value"], 123)

		// Test JSONB field
		jsonbMap := one["col_jsonb"].Map()
		t.Assert(jsonbMap["name"], "test")
		t.Assert(jsonbMap["value"], 123)
	})
}

// Test_Field_Scan_To_Struct tests scanning results to struct
func Test_Field_Scan_To_Struct(t *testing.T) {
	table := createInitAllTypesTable()
	defer dropTable(table)

	type TestRecord struct {
		Id         int64    `json:"id"`
		ColInt2    int16    `json:"col_int2"`
		ColInt4    int32    `json:"col_int4"`
		ColInt8    int64    `json:"col_int8"`
		ColVarchar string   `json:"col_varchar"`
		ColBool    bool     `json:"col_bool"`
		ColInt2Arr []int    `json:"col_int2_arr"`
		ColInt4Arr []int    `json:"col_int4_arr"`
		ColInt8Arr []int64  `json:"col_int8_arr"`
		ColTextArr []string `json:"col_text_arr"`
	}

	gtest.C(t, func(t *gtest.T) {
		var record TestRecord
		err := db.Model(table).Where("id", 1).Scan(&record)
		t.AssertNil(err)

		t.Assert(record.Id, int64(1))
		t.Assert(record.ColInt2, int16(1))
		t.Assert(record.ColInt4, int32(10))
		t.Assert(record.ColInt8, int64(100))
		t.AssertNE(record.ColVarchar, "")
		t.Assert(record.ColBool, false)

		// Test array fields scanned to struct
		t.Assert(len(record.ColInt2Arr), 3)
		t.Assert(record.ColInt2Arr[0], 1)
		t.Assert(record.ColInt2Arr[1], 2)
		t.Assert(record.ColInt2Arr[2], 1)

		t.Assert(len(record.ColTextArr), 3)
		t.Assert(record.ColTextArr[0], "x")
		t.Assert(record.ColTextArr[1], "y")
		t.Assert(record.ColTextArr[2], "z1")
	})
}

// Test_Field_Scan_To_Struct_Slice tests scanning multiple results to struct slice
func Test_Field_Scan_To_Struct_Slice(t *testing.T) {
	table := createInitAllTypesTable()
	defer dropTable(table)

	type TestRecord struct {
		Id         int64    `json:"id"`
		ColInt2    int16    `json:"col_int2"`
		ColVarchar string   `json:"col_varchar"`
		ColInt2Arr []int    `json:"col_int2_arr"`
		ColTextArr []string `json:"col_text_arr"`
	}

	gtest.C(t, func(t *gtest.T) {
		var records []TestRecord
		err := db.Model(table).OrderAsc("id").Limit(5).Scan(&records)
		t.AssertNil(err)

		t.Assert(len(records), 5)

		// Verify first record
		t.Assert(records[0].Id, int64(1))
		t.Assert(records[0].ColInt2, int16(1))
		t.Assert(len(records[0].ColInt2Arr), 3)

		// Verify last record
		t.Assert(records[4].Id, int64(5))
		t.Assert(records[4].ColInt2, int16(5))
	})
}

// Test_Field_Empty_Array tests handling empty arrays
func Test_Field_Empty_Array(t *testing.T) {
	table := createAllTypesTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert with empty array values (using default)
		_, err := db.Model(table).Data(g.Map{
			"col_int2":    1,
			"col_int4":    10,
			"col_numeric": 99.99,
			"col_varchar": "test",
			"col_bool":    true,
		}).Insert()
		t.AssertNil(err)

		// Query and verify empty arrays
		one, err := db.Model(table).OrderDesc("id").One()
		t.AssertNil(err)

		// Default empty arrays
		int2Arr := one["col_int2_arr"].Ints()
		t.Assert(len(int2Arr), 0)

		varcharArr := one["col_varchar_arr"].Strings()
		t.Assert(len(varcharArr), 0)
	})
}

// Test_Field_Null_Values tests handling NULL values
func Test_Field_Null_Values(t *testing.T) {
	table := createAllTypesTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert minimal required fields, leaving nullable fields as NULL
		_, err := db.Model(table).Data(g.Map{
			"col_int2":        1,
			"col_int4":        10,
			"col_numeric":     99.99,
			"col_varchar":     "test",
			"col_bool":        true,
			"col_varchar_arr": []string{},
		}).Insert()
		t.AssertNil(err)

		// Query and verify NULL handling
		one, err := db.Model(table).OrderDesc("id").One()
		t.AssertNil(err)

		// Nullable fields should return appropriate zero values
		t.Assert(one["col_text"].IsNil() || one["col_text"].IsEmpty(), true)
		t.Assert(one["col_int8_arr"].IsNil() || one["col_int8_arr"].IsEmpty(), true)
	})
}

// Test_Field_Float_Array_Type_Conversion tests float array type conversion (_float4, _float8)
func Test_Field_Float_Array_Type_Conversion(t *testing.T) {
	table := createInitAllTypesTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Query a single record
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one.IsEmpty(), false)

		// Test float4 array type conversions
		float4Arr := one["col_float4_arr"].Float32s()
		t.Assert(len(float4Arr), 3)
		t.Assert(float4Arr[0] > 0, true)
		t.Assert(float4Arr[1] > 0, true)

		// Test float8 array type conversions
		float8Arr := one["col_float8_arr"].Float64s()
		t.Assert(len(float8Arr), 3)
		t.Assert(float8Arr[0] > 0, true)
		t.Assert(float8Arr[1] > 0, true)
	})
}

// Test_Field_Numeric_Array_Type_Conversion tests numeric/decimal array type conversion
func Test_Field_Numeric_Array_Type_Conversion(t *testing.T) {
	table := createInitAllTypesTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Query a single record
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one.IsEmpty(), false)

		// Test numeric array type conversions
		numericArr := one["col_numeric_arr"].Float64s()
		t.Assert(len(numericArr), 3)
		t.Assert(numericArr[0] > 0, true)
		t.Assert(numericArr[1] > 0, true)

		// Test decimal array type conversions
		decimalArr := one["col_decimal_arr"].Float64s()
		if !one["col_decimal_arr"].IsNil() {
			t.Assert(len(decimalArr) > 0, true)
		}
	})
}

// Test_Field_Bool_Array_Type_Conversion tests bool array type conversion more thoroughly
func Test_Field_Bool_Array_Type_Conversion(t *testing.T) {
	table := createAllTypesTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert with specific bool array values
		_, err := db.Model(table).Data(g.Map{
			"col_int2":     1,
			"col_int4":     10,
			"col_numeric":  99.99,
			"col_varchar":  "test",
			"col_bool":     true,
			"col_bool_arr": []bool{true, false, true},
		}).Insert()
		t.AssertNil(err)

		// Query and verify
		one, err := db.Model(table).OrderDesc("id").One()
		t.AssertNil(err)

		// Test bool array
		boolArr := one["col_bool_arr"].Bools()
		t.Assert(len(boolArr), 3)
		t.Assert(boolArr[0], true)
		t.Assert(boolArr[1], false)
		t.Assert(boolArr[2], true)
	})
}

// Test_Field_Char_Array_Type tests char array type (_char)
func Test_Field_Char_Array_Type(t *testing.T) {
	table := createAllTypesTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert with char array values
		_, err := db.Model(table).Data(g.Map{
			"col_int2":        1,
			"col_int4":        10,
			"col_numeric":     99.99,
			"col_varchar":     "test",
			"col_bool":        true,
			"col_char_arr":    []string{"a", "b", "c"},
			"col_varchar_arr": []string{},
		}).Insert()
		t.AssertNil(err)

		// Query and verify
		one, err := db.Model(table).OrderDesc("id").One()
		t.AssertNil(err)

		// Test char array
		charArr := one["col_char_arr"].Strings()
		t.Assert(len(charArr), 3)
	})
}

// Test_Field_Bytea_Type tests bytea (binary) type conversion
func Test_Field_Bytea_Type(t *testing.T) {
	table := createAllTypesTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert with binary data
		binaryData := []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f} // "Hello" in hex
		_, err := db.Model(table).Data(g.Map{
			"col_int2":        1,
			"col_int4":        10,
			"col_numeric":     99.99,
			"col_varchar":     "test",
			"col_bool":        true,
			"col_bytea":       binaryData,
			"col_varchar_arr": []string{},
		}).Insert()
		t.AssertNil(err)

		// Query and verify
		one, err := db.Model(table).OrderDesc("id").One()
		t.AssertNil(err)

		// Test bytea field
		result := one["col_bytea"].Bytes()
		t.Assert(len(result), 5)
		t.Assert(result[0], 0x48) // 'H'
	})
}

// Test_Field_Bytea_Array_Type tests bytea array type (_bytea)
func Test_Field_Bytea_Array_Type(t *testing.T) {
	table := createAllTypesTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert with bytea array values using raw SQL
		// PostgreSQL bytea array literal format: ARRAY[E'\\x010203', E'\\x040506']::bytea[]
		_, err := db.Exec(ctx, fmt.Sprintf(`
			INSERT INTO %s (col_int2, col_int4, col_numeric, col_varchar, col_bool, col_varchar_arr, col_bytea_arr)
			VALUES (1, 10, 99.99, 'test', true, '{}', ARRAY[E'\\x010203', E'\\x040506']::bytea[])
		`, table))
		t.AssertNil(err)

		// Query and verify bytea array
		one, err := db.Model(table).OrderDesc("id").One()
		t.AssertNil(err)

		// Test bytea array field - should be converted to [][]byte
		byteaArrVal := one["col_bytea_arr"]
		t.Assert(byteaArrVal.IsNil(), false)

		// Verify the array contains the expected data
		byteaArr := byteaArrVal.Interfaces()
		t.Assert(len(byteaArr), 2)
	})
}

// Test_Field_Date_Array_Type tests date array type (_date)
func Test_Field_Date_Array_Type(t *testing.T) {
	table := createAllTypesTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Note: PostgreSQL _date array is not yet mapped in the driver
		// This test documents the limitation but can be extended when support is added

		_, err := db.Model(table).Data(g.Map{
			"col_int2":        1,
			"col_int4":        10,
			"col_numeric":     99.99,
			"col_varchar":     "test",
			"col_bool":        true,
			"col_varchar_arr": []string{},
		}).Insert()
		t.AssertNil(err)

		// Query and verify NULL date array is handled gracefully
		one, err := db.Model(table).OrderDesc("id").One()
		t.AssertNil(err)
		// date array should be nil or empty
		t.Assert(one["col_date_arr"].IsNil() || one["col_date_arr"].IsEmpty(), true)
	})
}

// Test_Field_Timestamp_Array_Type tests timestamp array type (_timestamp)
func Test_Field_Timestamp_Array_Type(t *testing.T) {
	table := createAllTypesTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Note: PostgreSQL _timestamp array is not yet mapped in the driver
		// This test documents the limitation but can be extended when support is added

		_, err := db.Model(table).Data(g.Map{
			"col_int2":        1,
			"col_int4":        10,
			"col_numeric":     99.99,
			"col_varchar":     "test",
			"col_bool":        true,
			"col_varchar_arr": []string{},
		}).Insert()
		t.AssertNil(err)

		// Query and verify NULL timestamp array is handled gracefully
		one, err := db.Model(table).OrderDesc("id").One()
		t.AssertNil(err)
		// timestamp array should be nil or empty
		t.Assert(one["col_timestamp_arr"].IsNil() || one["col_timestamp_arr"].IsEmpty(), true)
	})
}

// Test_Field_JSONB_Array_Type tests JSONB array type (_jsonb)
func Test_Field_JSONB_Array_Type(t *testing.T) {
	table := createAllTypesTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Note: PostgreSQL _jsonb array is not yet mapped in the driver
		// This test documents the limitation but can be extended when support is added

		_, err := db.Model(table).Data(g.Map{
			"col_int2":        1,
			"col_int4":        10,
			"col_numeric":     99.99,
			"col_varchar":     "test",
			"col_bool":        true,
			"col_varchar_arr": []string{},
		}).Insert()
		t.AssertNil(err)

		// Query and verify NULL jsonb array is handled gracefully
		one, err := db.Model(table).OrderDesc("id").One()
		t.AssertNil(err)
		// jsonb array should be nil or empty
		t.Assert(one["col_jsonb_arr"].IsNil() || one["col_jsonb_arr"].IsEmpty(), true)
	})
}

// Test_Field_UUID_Array_Type tests UUID array type (_uuid)
func Test_Field_UUID_Array_Type(t *testing.T) {
	table := createAllTypesTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert with UUID array values using raw SQL
		// PostgreSQL uuid array literal format: ARRAY['uuid1', 'uuid2']::uuid[]
		uuid1 := "550e8400-e29b-41d4-a716-446655440000"
		uuid2 := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
		uuid3 := "6ba7b811-9dad-11d1-80b4-00c04fd430c8"
		_, err := db.Exec(ctx, fmt.Sprintf(`
			INSERT INTO %s (col_int2, col_int4, col_numeric, col_varchar, col_bool, col_varchar_arr, col_uuid_arr)
			VALUES (1, 10, 99.99, 'test', true, '{}', ARRAY['%s', '%s', '%s']::uuid[])
		`, table, uuid1, uuid2, uuid3))
		t.AssertNil(err)

		// Query and verify UUID array
		one, err := db.Model(table).OrderDesc("id").One()
		t.AssertNil(err)

		// Test UUID array field - should be converted to []uuid.UUID
		uuidArrVal := one["col_uuid_arr"]
		t.Assert(uuidArrVal.IsNil(), false)

		// Verify the array contains the expected data as []uuid.UUID
		uuidArr := uuidArrVal.Interfaces()
		t.Assert(len(uuidArr), 3)

		// Verify each element is uuid.UUID type
		u1, ok := uuidArr[0].(uuid.UUID)
		t.Assert(ok, true)
		t.Assert(u1.String(), uuid1)

		u2, ok := uuidArr[1].(uuid.UUID)
		t.Assert(ok, true)
		t.Assert(u2.String(), uuid2)

		u3, ok := uuidArr[2].(uuid.UUID)
		t.Assert(ok, true)
		t.Assert(u3.String(), uuid3)
	})
}

// Test_Field_UUID_Type tests UUID type
func Test_Field_UUID_Type(t *testing.T) {
	table := createInitAllTypesTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Query and verify UUID field
		one, err := db.Model(table).OrderAsc("id").One()
		t.AssertNil(err)

		// Test UUID field - should be converted to uuid.UUID
		uuidVal := one["col_uuid"]
		t.Assert(uuidVal.IsNil(), false)

		// Verify the value is uuid.UUID type
		uuidObj, ok := uuidVal.Val().(uuid.UUID)
		t.Assert(ok, true)

		// Verify the UUID format
		uuidStr := uuidObj.String()
		t.Assert(len(uuidStr) > 0, true)
		// UUID should contain the pattern from insert: 550e8400-e29b-41d4-a716-44665544000X
		t.Assert(uuidStr, "550e8400-e29b-41d4-a716-446655440001")

		// Also verify we can still get string representation via .String()
		t.Assert(uuidVal.String(), "550e8400-e29b-41d4-a716-446655440001")
	})
}

// Test_Field_Bytea_Array_Type_Scan tests bytea array type and scanning
func Test_Field_Bytea_Array_Type_Scan(t *testing.T) {
	table := createInitAllTypesTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Query and verify bytea array field
		one, err := db.Model(table).OrderAsc("id").One()
		t.AssertNil(err)

		// Test bytea array field
		byteaArrVal := one["col_bytea_arr"]
		// bytea array should not be nil since we inserted data
		t.Assert(byteaArrVal.IsNil(), false)
	})
}

// Test_Field_Date_Array_Type_Scan tests date array type and scanning
func Test_Field_Date_Array_Type_Scan(t *testing.T) {
	table := createInitAllTypesTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Query and verify date array field
		one, err := db.Model(table).OrderAsc("id").One()
		t.AssertNil(err)

		// Test date array field
		dateArrVal := one["col_date_arr"]
		t.Assert(dateArrVal.IsNil(), false)

		// Verify the array contains the expected data
		dateArr := dateArrVal.Strings()
		t.Assert(len(dateArr) > 0, true)
	})
}

// Test_Field_Timestamp_Array_Type_Scan tests timestamp array type and scanning
func Test_Field_Timestamp_Array_Type_Scan(t *testing.T) {
	table := createInitAllTypesTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Query and verify timestamp array field
		one, err := db.Model(table).OrderAsc("id").One()
		t.AssertNil(err)

		// Test timestamp array field
		timestampArrVal := one["col_timestamp_arr"]
		t.Assert(timestampArrVal.IsNil(), false)

		// Verify the array contains the expected data
		timestampArr := timestampArrVal.Strings()
		t.Assert(len(timestampArr) > 0, true)
	})
}

// Test_Field_JSONB_Array_Type_Scan tests JSONB array type and scanning
func Test_Field_JSONB_Array_Type_Scan(t *testing.T) {
	table := createInitAllTypesTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Query and verify JSONB array field
		one, err := db.Model(table).OrderAsc("id").One()
		t.AssertNil(err)

		// Test JSONB array field
		jsonbArrVal := one["col_jsonb_arr"]
		t.Assert(jsonbArrVal.IsNil(), false)
	})
}

// Test_Field_UUID_Query tests querying by UUID field
func Test_Field_UUID_Query(t *testing.T) {
	table := createInitAllTypesTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Test 1: Query by UUID string
		uuidStr := "550e8400-e29b-41d4-a716-446655440001"
		one, err := db.Model(table).Where("col_uuid", uuidStr).One()
		t.AssertNil(err)
		t.Assert(one.IsEmpty(), false)
		t.Assert(one["id"].Int(), 1)

		// Verify the returned UUID is correct
		uuidObj, ok := one["col_uuid"].Val().(uuid.UUID)
		t.Assert(ok, true)
		t.Assert(uuidObj.String(), uuidStr)

		// Test 2: Query by uuid.UUID type directly
		uuidVal, err := uuid.Parse("550e8400-e29b-41d4-a716-446655440002")
		t.AssertNil(err)
		one, err = db.Model(table).Where("col_uuid", uuidVal).One()
		t.AssertNil(err)
		t.Assert(one.IsEmpty(), false)
		t.Assert(one["id"].Int(), 2)

		// Test 3: Query by UUID string using g.Map
		one, err = db.Model(table).Where(g.Map{
			"col_uuid": "550e8400-e29b-41d4-a716-446655440003",
		}).One()
		t.AssertNil(err)
		t.Assert(one.IsEmpty(), false)
		t.Assert(one["id"].Int(), 3)

		// Test 4: Query by uuid.UUID type using g.Map
		uuidVal, err = uuid.Parse("550e8400-e29b-41d4-a716-446655440004")
		t.AssertNil(err)
		one, err = db.Model(table).Where(g.Map{
			"col_uuid": uuidVal,
		}).One()
		t.AssertNil(err)
		t.Assert(one.IsEmpty(), false)
		t.Assert(one["id"].Int(), 4)

		// Test 5: Query non-existent UUID
		one, err = db.Model(table).Where("col_uuid", "00000000-0000-0000-0000-000000000000").One()
		t.AssertNil(err)
		t.Assert(one.IsEmpty(), true)

		// Test 6: Query multiple records by UUID IN clause with strings
		all, err := db.Model(table).WhereIn("col_uuid", g.Slice{
			"550e8400-e29b-41d4-a716-446655440001",
			"550e8400-e29b-41d4-a716-446655440002",
		}).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["id"].Int(), 1)
		t.Assert(all[1]["id"].Int(), 2)

		// Test 7: Query multiple records by UUID IN clause with uuid.UUID types
		uuid1, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440003")
		uuid2, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440004")
		all, err = db.Model(table).WhereIn("col_uuid", g.Slice{uuid1, uuid2}).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["id"].Int(), 3)
		t.Assert(all[1]["id"].Int(), 4)
	})
}
