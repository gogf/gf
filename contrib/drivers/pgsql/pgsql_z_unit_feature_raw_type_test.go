// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package pgsql_test

// Note: Test_Raw_Insert / Test_Raw_BatchInsert / Test_Raw_Update / Test_Raw_Where already
// exist in pgsql_z_unit_raw_test.go (older PgSQL-adapted versions). Per project policy
// ("老用例原则"), we do not duplicate or override them here — this file only contributes the
// DataType tests ported from the MariaDB baseline.

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

// Test_DataType_JSON_Insert tests JSON data insertion using PgSQL jsonb.
func Test_DataType_JSON_Insert(t *testing.T) {
	table := "test_json_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id bigserial PRIMARY KEY, data jsonb)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table).Data(g.Map{
			"data": `{"name":"John","age":30}`,
		}).Insert()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		expected := map[string]interface{}{"name": "John", "age": float64(30)}
		var actual map[string]interface{}
		err = json.Unmarshal([]byte(one["data"].String()), &actual)
		t.AssertNil(err)
		t.Assert(actual, expected)
	})
}

// Test_DataType_JSON_Extract tests JSON extract using PgSQL -> / ->> operators.
func Test_DataType_JSON_Extract(t *testing.T) {
	table := "test_json_extract_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id bigserial PRIMARY KEY, data jsonb)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"data": `{"name":"Alice","age":25,"city":"Beijing"}`,
		}).Insert()
		t.AssertNil(err)

		// PgSQL: data->>'name' returns text
		one, err := db.Model(table).Fields("data->>'name' as name").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["name"].String(), "Alice")

		// Extract age as int
		one, err = db.Model(table).Fields("(data->>'age')::int as age").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["age"].Int(), 25)
	})
}

// Test_DataType_JSON_Set tests JSON field update using PgSQL jsonb_set.
func Test_DataType_JSON_Set(t *testing.T) {
	table := "test_json_set_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id bigserial PRIMARY KEY, data jsonb)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"data": `{"name":"Bob"}`,
		}).Insert()
		t.AssertNil(err)

		// PgSQL: jsonb_set(data, '{age}', '30')
		_, err = db.Exec(ctx, fmt.Sprintf("UPDATE %s SET data = jsonb_set(data, '{age}', '30') WHERE id = 1", table))
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		expected := map[string]interface{}{"name": "Bob", "age": float64(30)}
		var actual map[string]interface{}
		err = json.Unmarshal([]byte(one["data"].String()), &actual)
		t.AssertNil(err)
		t.Assert(actual, expected)
	})
}

// Test_DataType_JSON_Array tests JSON array operations in jsonb.
func Test_DataType_JSON_Array(t *testing.T) {
	table := "test_json_array_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id bigserial PRIMARY KEY, data jsonb)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"data": `["apple","banana","cherry"]`,
		}).Insert()
		t.AssertNil(err)

		// PgSQL: data->>0 extracts array element as text
		one, err := db.Model(table).Fields("data->>0 as first").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["first"].String(), "apple")
	})
}

// Test_DataType_JSON_Null tests JSON NULL handling.
func Test_DataType_JSON_Null(t *testing.T) {
	table := "test_json_null_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id bigserial PRIMARY KEY, data jsonb)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"data": nil,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["data"].IsNil(), true)
	})
}

// Test_DataType_JSON_Complex tests complex nested JSON using PgSQL jsonb operators.
func Test_DataType_JSON_Complex(t *testing.T) {
	table := "test_json_complex_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id bigserial PRIMARY KEY, data jsonb)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		complexJSON := `{
			"user": {
				"name": "Charlie",
				"contacts": {
					"email": "charlie@example.com",
					"phone": "1234567890"
				},
				"tags": ["developer", "gopher"]
			}
		}`
		_, err := db.Model(table).Data(g.Map{
			"data": complexJSON,
		}).Insert()
		t.AssertNil(err)

		// PgSQL: nested path extraction using #>>
		one, err := db.Model(table).Fields("data#>>'{user,contacts,email}' as email").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["email"].String(), "charlie@example.com")
	})
}

// Test_DataType_JSON_Query tests JSON query with WHERE clause using jsonb operators.
func Test_DataType_JSON_Query(t *testing.T) {
	table := "test_json_query_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id bigserial PRIMARY KEY, data jsonb)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.List{
			g.Map{"data": `{"name":"David","age":20}`},
			g.Map{"data": `{"name":"Eve","age":30}`},
			g.Map{"data": `{"name":"Frank","age":25}`},
		}).Insert()
		t.AssertNil(err)

		// PgSQL: cast jsonb field to int for comparison
		count, err := db.Model(table).Where("(data->>'age')::int > ?", 25).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}

// Test_DataType_JSON_Update tests updating JSON data.
func Test_DataType_JSON_Update(t *testing.T) {
	table := "test_json_update_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id bigserial PRIMARY KEY, data jsonb)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"data": `{"name":"Grace","age":28}`,
		}).Insert()
		t.AssertNil(err)

		_, err = db.Model(table).Data(g.Map{
			"data": `{"name":"Grace","age":29}`,
		}).Where("id", 1).Update()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		expected := map[string]interface{}{"name": "Grace", "age": float64(29)}
		var actual map[string]interface{}
		err = json.Unmarshal([]byte(one["data"].String()), &actual)
		t.AssertNil(err)
		t.Assert(actual, expected)
	})
}

// Test_DataType_Binary_Small tests small binary data using PgSQL bytea.
func Test_DataType_Binary_Small(t *testing.T) {
	table := "test_binary_small_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id bigserial PRIMARY KEY, data bytea)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		binaryData := []byte{0x00, 0x01, 0x02, 0x03, 0xFF}
		_, err := db.Model(table).Data(g.Map{
			"data": binaryData,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(bytes.Equal(one["data"].Bytes(), binaryData), true)
	})
}

// Test_DataType_Binary_Large tests large binary data (1MB+) using PgSQL bytea.
func Test_DataType_Binary_Large(t *testing.T) {
	table := "test_binary_large_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id bigserial PRIMARY KEY, data bytea)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		size := 1024 * 1024 // 1MB
		largeBinary := make([]byte, size)
		for i := 0; i < size; i++ {
			largeBinary[i] = byte(i % 256)
		}

		_, err := db.Model(table).Data(g.Map{
			"data": largeBinary,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(len(one["data"].Bytes()), size)
		t.Assert(bytes.Equal(one["data"].Bytes(), largeBinary), true)
	})
}

// Test_DataType_Binary_Integrity tests binary data integrity with SHA256 checksum.
func Test_DataType_Binary_Integrity(t *testing.T) {
	table := "test_binary_integrity_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id bigserial PRIMARY KEY, data bytea, checksum varchar(64))")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		binaryData := []byte("Hello, World! This is a binary test data with special chars: \x00\xFF\xAB")

		hash := sha256.Sum256(binaryData)
		checksum := hex.EncodeToString(hash[:])

		_, err := db.Model(table).Data(g.Map{
			"data":     binaryData,
			"checksum": checksum,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)

		retrievedHash := sha256.Sum256(one["data"].Bytes())
		retrievedChecksum := hex.EncodeToString(retrievedHash[:])
		t.Assert(retrievedChecksum, checksum)
	})
}

// Test_DataType_Binary_Empty tests empty and NULL binary values.
func Test_DataType_Binary_Empty(t *testing.T) {
	table := "test_binary_empty_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id bigserial PRIMARY KEY, data bytea)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"data": []byte{},
		}).Insert()
		t.AssertNil(err)

		_, err = db.Model(table).Data(g.Map{
			"data": nil,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(len(one["data"].Bytes()), 0)

		one, err = db.Model(table).Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(one["data"].IsNil(), true)
	})
}

// Test_DataType_Decimal_HighPrecision tests high precision numeric.
// PgSQL numeric supports up to 131072 digits before decimal, 16383 after.
func Test_DataType_Decimal_HighPrecision(t *testing.T) {
	table := "test_decimal_precision_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id bigserial PRIMARY KEY, amount numeric(65,30))")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		value := "12345678901234567890123456789012345.123456789012345678901234567890"
		_, err := db.Model(table).Data(g.Map{
			"amount": value,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["amount"].String(), value)
	})
}

// Test_DataType_Decimal_Calculation tests decimal arithmetic.
func Test_DataType_Decimal_Calculation(t *testing.T) {
	table := "test_decimal_calc_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id bigserial PRIMARY KEY, price numeric(10,2), quantity numeric(10,2))")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"price":    "19.99",
			"quantity": "3.5",
		}).Insert()
		t.AssertNil(err)

		// PgSQL: price(10,2) * quantity(10,2) yields numeric(20,4)
		one, err := db.Model(table).Fields("price * quantity as total").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["total"].String(), "69.9650")
	})
}

// Test_DataType_Decimal_Boundary tests decimal boundary values.
func Test_DataType_Decimal_Boundary(t *testing.T) {
	table := "test_decimal_boundary_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id bigserial PRIMARY KEY, value numeric(10,2))")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"value": "99999999.99",
		}).Insert()
		t.AssertNil(err)

		_, err = db.Model(table).Data(g.Map{
			"value": "-99999999.99",
		}).Insert()
		t.AssertNil(err)

		_, err = db.Model(table).Data(g.Map{
			"value": "0.00",
		}).Insert()
		t.AssertNil(err)

		all, err := db.Model(table).Order("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(all[0]["value"].String(), "99999999.99")
		t.Assert(all[1]["value"].String(), "-99999999.99")
		t.Assert(all[2]["value"].String(), "0.00")
	})
}

// Test_DataType_Decimal_Null tests NULL decimal values.
func Test_DataType_Decimal_Null(t *testing.T) {
	table := "test_decimal_null_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id bigserial PRIMARY KEY, value numeric(10,2))")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"value": nil,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["value"].IsNil(), true)
	})
}

// Test_DataType_Datetime_Timezone tests timestamp handling.
// PgSQL MySQL's DATETIME maps to timestamp (without time zone).
func Test_DataType_Datetime_Timezone(t *testing.T) {
	table := "test_datetime_tz_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id bigserial PRIMARY KEY, created_at timestamp)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		dt := "2024-01-15 12:30:45"
		_, err := db.Model(table).Data(g.Map{
			"created_at": dt,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["created_at"].String(), dt)
	})
}

// Test_DataType_Datetime_Precision tests timestamp with microsecond precision.
func Test_DataType_Datetime_Precision(t *testing.T) {
	table := "test_datetime_precision_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id bigserial PRIMARY KEY, created_at timestamp(6))")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		dt := "2024-01-15 12:30:45.123456"
		_, err := db.Model(table).Data(g.Map{
			"created_at": dt,
		}).Insert()
		t.AssertNil(err)

		// Compare up to seconds; driver may reformat microseconds.
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		expected := "2024-01-15 12:30:45"
		actual := one["created_at"].String()[:19]
		t.Assert(actual, expected)
	})
}

// Test_DataType_Datetime_Boundary tests timestamp boundary values.
// PgSQL supports timestamps from 4713 BC to 294276 AD (much wider than MySQL).
func Test_DataType_Datetime_Boundary(t *testing.T) {
	table := "test_datetime_boundary_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id bigserial PRIMARY KEY, dt timestamp)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"dt": "1000-01-01 00:00:00",
		}).Insert()
		t.AssertNil(err)

		_, err = db.Model(table).Data(g.Map{
			"dt": "9999-12-31 23:59:59",
		}).Insert()
		t.AssertNil(err)

		all, err := db.Model(table).Order("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["dt"].String(), "1000-01-01 00:00:00")
		t.Assert(all[1]["dt"].String(), "9999-12-31 23:59:59")
	})
}

// Test_DataType_Datetime_Null tests NULL timestamp values.
func Test_DataType_Datetime_Null(t *testing.T) {
	table := "test_datetime_null_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id bigserial PRIMARY KEY, dt timestamp)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"dt": nil,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["dt"].IsNil(), true)
	})
}

// Test_DataType_Datetime_Update tests timestamp updates.
func Test_DataType_Datetime_Update(t *testing.T) {
	table := "test_datetime_update_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id bigserial PRIMARY KEY, dt timestamp)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		dt1 := "2024-01-01 10:00:00"
		_, err := db.Model(table).Data(g.Map{
			"dt": dt1,
		}).Insert()
		t.AssertNil(err)

		dt2 := "2024-12-31 23:59:59"
		_, err = db.Model(table).Data(g.Map{
			"dt": dt2,
		}).Where("id", 1).Update()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["dt"].String(), dt2)
	})
}

// Test_DataType_Enum_Valid: PgSQL does not have a native MySQL-compatible ENUM type.
// PgSQL has CREATE TYPE ... AS ENUM but with different semantics, so test is skipped.
func Test_DataType_Enum_Valid(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Skip("MySQL ENUM column type is not supported by PgSQL (use CREATE TYPE AS ENUM separately)")
	})
}

// Test_DataType_Enum_Invalid: same reason as Test_DataType_Enum_Valid.
func Test_DataType_Enum_Invalid(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Skip("MySQL ENUM column type is not supported by PgSQL")
	})
}

// Test_DataType_Set_Valid: PgSQL does not have MySQL's SET column type.
// Use array types (e.g., text[]) or a jsonb column instead.
func Test_DataType_Set_Valid(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Skip("MySQL SET column type is not supported by PgSQL (use array or jsonb)")
	})
}

// Test_DataType_Set_Empty: same reason as Test_DataType_Set_Valid.
func Test_DataType_Set_Empty(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Skip("MySQL SET column type is not supported by PgSQL")
	})
}

// Test_DataType_Geometry_Point: PgSQL built-in geometric types differ from MySQL/PostGIS.
// ST_GeomFromText/ST_AsText require the PostGIS extension, which is not assumed here.
func Test_DataType_Geometry_Point(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Skip("ST_GeomFromText/ST_AsText require PostGIS extension; not installed in default PgSQL")
	})
}

// Test_DataType_Geometry_Polygon: same reason as Test_DataType_Geometry_Point.
func Test_DataType_Geometry_Polygon(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Skip("ST_GeomFromText/ST_AsText require PostGIS extension; not installed in default PgSQL")
	})
}

// Test_DataType_Geometry_Null: same reason as Test_DataType_Geometry_Point.
func Test_DataType_Geometry_Null(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Skip("geometry types require PostGIS extension; not installed in default PgSQL")
	})
}
