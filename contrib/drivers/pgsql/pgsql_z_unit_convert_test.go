// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package pgsql_test

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/test/gtest"

	"github.com/gogf/gf/contrib/drivers/pgsql/v2"
)

// Test_CheckLocalTypeForField tests the CheckLocalTypeForField method
// for various PostgreSQL types
func Test_CheckLocalTypeForField(t *testing.T) {
	var (
		ctx    = context.Background()
		driver = pgsql.Driver{}
	)

	gtest.C(t, func(t *gtest.T) {
		// Test basic integer types
		localType, err := driver.CheckLocalTypeForField(ctx, "int2", nil)
		t.AssertNil(err)
		t.Assert(localType, gdb.LocalTypeInt)

		localType, err = driver.CheckLocalTypeForField(ctx, "int4", nil)
		t.AssertNil(err)
		t.Assert(localType, gdb.LocalTypeInt)

		localType, err = driver.CheckLocalTypeForField(ctx, "int8", nil)
		t.AssertNil(err)
		t.Assert(localType, gdb.LocalTypeInt64)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test integer array types
		localType, err := driver.CheckLocalTypeForField(ctx, "_int2", nil)
		t.AssertNil(err)
		t.Assert(localType, gdb.LocalTypeInt32Slice)

		localType, err = driver.CheckLocalTypeForField(ctx, "_int4", nil)
		t.AssertNil(err)
		t.Assert(localType, gdb.LocalTypeInt32Slice)

		localType, err = driver.CheckLocalTypeForField(ctx, "_int8", nil)
		t.AssertNil(err)
		t.Assert(localType, gdb.LocalTypeInt64Slice)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test float array types
		localType, err := driver.CheckLocalTypeForField(ctx, "_float4", nil)
		t.AssertNil(err)
		t.Assert(localType, gdb.LocalTypeFloat32Slice)

		localType, err = driver.CheckLocalTypeForField(ctx, "_float8", nil)
		t.AssertNil(err)
		t.Assert(localType, gdb.LocalTypeFloat64Slice)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test boolean array type
		localType, err := driver.CheckLocalTypeForField(ctx, "_bool", nil)
		t.AssertNil(err)
		t.Assert(localType, gdb.LocalTypeBoolSlice)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test string array types
		localType, err := driver.CheckLocalTypeForField(ctx, "_varchar", nil)
		t.AssertNil(err)
		t.Assert(localType, gdb.LocalTypeStringSlice)

		localType, err = driver.CheckLocalTypeForField(ctx, "_text", nil)
		t.AssertNil(err)
		t.Assert(localType, gdb.LocalTypeStringSlice)

		localType, err = driver.CheckLocalTypeForField(ctx, "_char", nil)
		t.AssertNil(err)
		t.Assert(localType, gdb.LocalTypeStringSlice)

		localType, err = driver.CheckLocalTypeForField(ctx, "_bpchar", nil)
		t.AssertNil(err)
		t.Assert(localType, gdb.LocalTypeStringSlice)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test numeric array types
		localType, err := driver.CheckLocalTypeForField(ctx, "_numeric", nil)
		t.AssertNil(err)
		t.Assert(localType, gdb.LocalTypeFloat64Slice)

		localType, err = driver.CheckLocalTypeForField(ctx, "_decimal", nil)
		t.AssertNil(err)
		t.Assert(localType, gdb.LocalTypeFloat64Slice)

		localType, err = driver.CheckLocalTypeForField(ctx, "_money", nil)
		t.AssertNil(err)
		t.Assert(localType, gdb.LocalTypeFloat64Slice)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test bytea array type
		localType, err := driver.CheckLocalTypeForField(ctx, "_bytea", nil)
		t.AssertNil(err)
		t.Assert(localType, gdb.LocalTypeBytesSlice)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test uuid type
		localType, err := driver.CheckLocalTypeForField(ctx, "uuid", nil)
		t.AssertNil(err)
		t.Assert(localType, gdb.LocalTypeUUID)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test uuid array type
		localType, err := driver.CheckLocalTypeForField(ctx, "_uuid", nil)
		t.AssertNil(err)
		t.Assert(localType, gdb.LocalTypeUUIDSlice)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test type with precision, e.g., "numeric(10,2)"
		localType, err := driver.CheckLocalTypeForField(ctx, "int2(5)", nil)
		t.AssertNil(err)
		t.Assert(localType, gdb.LocalTypeInt)

		localType, err = driver.CheckLocalTypeForField(ctx, "int4(10)", nil)
		t.AssertNil(err)
		t.Assert(localType, gdb.LocalTypeInt)

		localType, err = driver.CheckLocalTypeForField(ctx, "INT8(20)", nil)
		t.AssertNil(err)
		t.Assert(localType, gdb.LocalTypeInt64)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test uppercase type names
		localType, err := driver.CheckLocalTypeForField(ctx, "INT2", nil)
		t.AssertNil(err)
		t.Assert(localType, gdb.LocalTypeInt)

		localType, err = driver.CheckLocalTypeForField(ctx, "_INT4", nil)
		t.AssertNil(err)
		t.Assert(localType, gdb.LocalTypeInt32Slice)
	})
}

// Test_ConvertValueForLocal tests the ConvertValueForLocal method
func Test_ConvertValueForLocal(t *testing.T) {
	var (
		ctx    = context.Background()
		driver = pgsql.Driver{}
	)

	gtest.C(t, func(t *gtest.T) {
		// Test _int2 array conversion
		result, err := driver.ConvertValueForLocal(ctx, "_int2", []byte(`{1,2,3}`))
		t.AssertNil(err)
		t.Assert(result, []int32{1, 2, 3})
	})

	gtest.C(t, func(t *gtest.T) {
		// Test _int4 array conversion
		result, err := driver.ConvertValueForLocal(ctx, "_int4", []byte(`{10,20,30}`))
		t.AssertNil(err)
		t.Assert(result, []int32{10, 20, 30})
	})

	gtest.C(t, func(t *gtest.T) {
		// Test _int8 array conversion
		result, err := driver.ConvertValueForLocal(ctx, "_int8", []byte(`{100,200,300}`))
		t.AssertNil(err)
		t.Assert(result, []int64{100, 200, 300})
	})

	gtest.C(t, func(t *gtest.T) {
		// Test _float4 array conversion
		result, err := driver.ConvertValueForLocal(ctx, "_float4", []byte(`{1.1,2.2,3.3}`))
		t.AssertNil(err)
		resultArr := result.([]float32)
		t.Assert(len(resultArr), 3)
		t.Assert(resultArr[0] > 1.0 && resultArr[0] < 1.2, true)
		t.Assert(resultArr[1] > 2.1 && resultArr[1] < 2.3, true)
		t.Assert(resultArr[2] > 3.2 && resultArr[2] < 3.4, true)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test _float8 array conversion
		result, err := driver.ConvertValueForLocal(ctx, "_float8", []byte(`{1.11,2.22,3.33}`))
		t.AssertNil(err)
		resultArr := result.([]float64)
		t.Assert(len(resultArr), 3)
		t.Assert(resultArr[0] > 1.1 && resultArr[0] < 1.12, true)
		t.Assert(resultArr[1] > 2.21 && resultArr[1] < 2.23, true)
		t.Assert(resultArr[2] > 3.32 && resultArr[2] < 3.34, true)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test _bool array conversion
		result, err := driver.ConvertValueForLocal(ctx, "_bool", []byte(`{t,f,t}`))
		t.AssertNil(err)
		t.Assert(result, []bool{true, false, true})
	})

	gtest.C(t, func(t *gtest.T) {
		// Test _varchar array conversion
		result, err := driver.ConvertValueForLocal(ctx, "_varchar", []byte(`{a,b,c}`))
		t.AssertNil(err)
		t.Assert(result, []string{"a", "b", "c"})
	})

	gtest.C(t, func(t *gtest.T) {
		// Test _text array conversion
		result, err := driver.ConvertValueForLocal(ctx, "_text", []byte(`{hello,world}`))
		t.AssertNil(err)
		t.Assert(result, []string{"hello", "world"})
	})

	gtest.C(t, func(t *gtest.T) {
		// Test _char array conversion
		result, err := driver.ConvertValueForLocal(ctx, "_char", []byte(`{x,y,z}`))
		t.AssertNil(err)
		t.Assert(result, []string{"x", "y", "z"})
	})

	gtest.C(t, func(t *gtest.T) {
		// Test _bpchar array conversion
		result, err := driver.ConvertValueForLocal(ctx, "_bpchar", []byte(`{a,b}`))
		t.AssertNil(err)
		t.Assert(result, []string{"a", "b"})
	})

	gtest.C(t, func(t *gtest.T) {
		// Test _numeric array conversion
		result, err := driver.ConvertValueForLocal(ctx, "_numeric", []byte(`{1.11,2.22}`))
		t.AssertNil(err)
		resultArr := result.([]float64)
		t.Assert(len(resultArr), 2)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test _decimal array conversion
		result, err := driver.ConvertValueForLocal(ctx, "_decimal", []byte(`{3.33,4.44}`))
		t.AssertNil(err)
		resultArr := result.([]float64)
		t.Assert(len(resultArr), 2)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test _money array conversion
		result, err := driver.ConvertValueForLocal(ctx, "_money", []byte(`{5.55,6.66}`))
		t.AssertNil(err)
		resultArr := result.([]float64)
		t.Assert(len(resultArr), 2)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test _bytea array conversion
		result, err := driver.ConvertValueForLocal(ctx, "_bytea", []byte(`{"\\x68656c6c6f","\\x776f726c64"}`))
		t.AssertNil(err)
		resultArr := result.([][]byte)
		t.Assert(len(resultArr), 2)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test uuid conversion from []byte
		result, err := driver.ConvertValueForLocal(ctx, "uuid", []byte(`550e8400-e29b-41d4-a716-446655440000`))
		t.AssertNil(err)
		t.Assert(result.(uuid.UUID).String(), "550e8400-e29b-41d4-a716-446655440000")
	})

	gtest.C(t, func(t *gtest.T) {
		// Test uuid conversion from string
		result, err := driver.ConvertValueForLocal(ctx, "uuid", "550e8400-e29b-41d4-a716-446655440000")
		t.AssertNil(err)
		t.Assert(result.(uuid.UUID).String(), "550e8400-e29b-41d4-a716-446655440000")
	})

	gtest.C(t, func(t *gtest.T) {
		// Test uuid conversion error case with invalid uuid
		_, err := driver.ConvertValueForLocal(ctx, "uuid", "invalid-uuid")
		t.AssertNE(err, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test _uuid array conversion
		result, err := driver.ConvertValueForLocal(ctx, "_uuid", []byte(`{550e8400-e29b-41d4-a716-446655440000,6ba7b810-9dad-11d1-80b4-00c04fd430c8}`))
		t.AssertNil(err)
		resultArr := result.([]uuid.UUID)
		t.Assert(len(resultArr), 2)
		t.Assert(resultArr[0].String(), "550e8400-e29b-41d4-a716-446655440000")
		t.Assert(resultArr[1].String(), "6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	})

	gtest.C(t, func(t *gtest.T) {
		// Test _uuid array conversion error case
		_, err := driver.ConvertValueForLocal(ctx, "_uuid", []byte(`{invalid-uuid}`))
		t.AssertNE(err, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test error case with invalid data for _int2
		_, err := driver.ConvertValueForLocal(ctx, "_int2", "invalid")
		t.AssertNE(err, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test error case with invalid data for _int4
		_, err := driver.ConvertValueForLocal(ctx, "_int4", "invalid")
		t.AssertNE(err, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test error case with invalid data for _int8
		_, err := driver.ConvertValueForLocal(ctx, "_int8", "invalid")
		t.AssertNE(err, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test error case with invalid data for _float4
		_, err := driver.ConvertValueForLocal(ctx, "_float4", "invalid")
		t.AssertNE(err, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test error case with invalid data for _float8
		_, err := driver.ConvertValueForLocal(ctx, "_float8", "invalid")
		t.AssertNE(err, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test error case with invalid data for _bool
		_, err := driver.ConvertValueForLocal(ctx, "_bool", "invalid")
		t.AssertNE(err, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test error case with invalid data for _varchar
		_, err := driver.ConvertValueForLocal(ctx, "_varchar", 12345)
		t.AssertNE(err, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test error case with invalid data for _numeric
		_, err := driver.ConvertValueForLocal(ctx, "_numeric", "invalid")
		t.AssertNE(err, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test error case with invalid data for _bytea
		_, err := driver.ConvertValueForLocal(ctx, "_bytea", "invalid")
		t.AssertNE(err, nil)
	})
}

// Test_ConvertValueForField tests the ConvertValueForField method
func Test_ConvertValueForField(t *testing.T) {
	var (
		ctx    = context.Background()
		driver = pgsql.Driver{}
	)

	gtest.C(t, func(t *gtest.T) {
		// Test nil value
		result, err := driver.ConvertValueForField(ctx, "varchar", nil)
		t.AssertNil(err)
		t.Assert(result, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test slice value for non-json type (should convert [] to {})
		result, err := driver.ConvertValueForField(ctx, "int4[]", []int{1, 2, 3})
		t.AssertNil(err)
		t.Assert(result, "{1,2,3}")
	})

	gtest.C(t, func(t *gtest.T) {
		// Test slice value for non-json type with strings
		// Note: gconv.String for []string{"a","b","c"} produces ["a","b","c"] which then gets converted to {"a","b","c"}
		result, err := driver.ConvertValueForField(ctx, "varchar[]", []string{"a", "b", "c"})
		t.AssertNil(err)
		t.Assert(result, `{"a","b","c"}`)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test slice value for json type (should keep [] as is)
		result, err := driver.ConvertValueForField(ctx, "json", []int{1, 2, 3})
		t.AssertNil(err)
		t.Assert(result, "[1,2,3]")
	})

	gtest.C(t, func(t *gtest.T) {
		// Test slice value for jsonb type (should keep [] as is)
		result, err := driver.ConvertValueForField(ctx, "jsonb", []string{"a", "b"})
		t.AssertNil(err)
		t.Assert(result, `["a","b"]`)
	})
}
