// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

var (
	rawTypeGeometryOnce sync.Once
	rawTypeGeometryErr  error
)

// rawTypeInitGeometry initializes the DMGEO package of the spatial types and functions once, as a
// DM instance does not have it until SP_INIT_GEO_SYS creates it.
func rawTypeInitGeometry(t *testing.T) {
	rawTypeGeometryOnce.Do(func() {
		value, err := db.GetValue(ctx, "SELECT SF_CHECK_GEO_SYS() FROM DUAL")
		if err != nil {
			rawTypeGeometryErr = err
			return
		}
		if value.Int() == 0 {
			_, rawTypeGeometryErr = db.Exec(ctx, "CALL SP_INIT_GEO_SYS(1)")
		}
	})
	if rawTypeGeometryErr != nil {
		t.Fatal(rawTypeGeometryErr)
	}
}

func Test_Raw_Insert(t *testing.T) {
	table := createUserTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		user := db.Model(table)
		// Note: DM rejects a column in VALUES (-2120), where MySQL reads the default 0 of id, so the
		// raw expression adds to 0.
		result, err := user.Data(g.Map{
			"id":          gdb.Raw("0+2"),
			"passport":    "port_1",
			"password":    "pass_1",
			"nickname":    "name_1",
			"create_time": gdb.Raw("now()"),
		}).Insert()
		t.AssertNil(err)
		n, _ := result.LastInsertId()
		t.Assert(n, 2)
	})
}

func Test_Raw_BatchInsert(t *testing.T) {
	table := createUserTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		user := db.Model(table)
		// Note: DM rejects a column in VALUES (-2120), where MySQL reads the default 0 of id, so the
		// raw expressions add to 0.
		result, err := user.Data(
			g.List{
				g.Map{
					"id":          gdb.Raw("0+2"),
					"passport":    "port_2",
					"password":    "pass_2",
					"nickname":    "name_2",
					"create_time": gdb.Raw("now()"),
				},
				g.Map{
					"id":          gdb.Raw("0+4"),
					"passport":    "port_4",
					"password":    "pass_4",
					"nickname":    "name_4",
					"create_time": gdb.Raw("now()"),
				},
			},
		).Insert()
		t.AssertNil(err)
		n, _ := result.LastInsertId()
		t.Assert(n, 4)
	})
}

func Test_Raw_Update(t *testing.T) {
	table := createInitUserTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		user := db.Model(table)
		result, err := user.Data(g.Map{
			"id":          gdb.Raw("id+100"),
			"create_time": gdb.Raw("now()"),
		}).Where("id", 1).Update()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
	})
	gtest.C(t, func(t *gtest.T) {
		user := db.Model(table)
		n, err := user.Where("id", 101).Count()
		t.AssertNil(err)
		t.Assert(n, 1)
	})
}

func Test_Raw_Where(t *testing.T) {
	table1 := createUserTable("Test_Raw_Where_Table1")
	table2 := createUserTable("Test_Raw_Where_Table2")
	defer dropTable(table1)
	defer dropTable(table2)

	// https://github.com/gogf/gf/issues/3922
	// Note: the DM driver does not quote identifiers.
	gtest.C(t, func(t *gtest.T) {
		expectSql := "SELECT * FROM Test_Raw_Where_Table1 AS A WHERE NOT EXISTS (SELECT B.id FROM Test_Raw_Where_Table2 AS B WHERE B.id=A.id) LIMIT 1"
		sql, err := gdb.ToSQL(ctx, func(ctx context.Context) error {
			s := db.Model(table2).As("B").Ctx(ctx).Fields("B.id").Where("B.id", gdb.Raw("A.id"))
			m := db.Model(table1).As("A").Ctx(ctx).Where("NOT EXISTS ?", s).Limit(1)
			_, err := m.All()
			return err
		})
		t.AssertNil(err)
		t.Assert(expectSql, sql)
	})
	gtest.C(t, func(t *gtest.T) {
		expectSql := "SELECT * FROM Test_Raw_Where_Table1 AS A WHERE NOT EXISTS (SELECT B.id FROM Test_Raw_Where_Table2 AS B WHERE B.id=A.id) LIMIT 1"
		sql, err := gdb.ToSQL(ctx, func(ctx context.Context) error {
			s := db.Model(table2).As("B").Ctx(ctx).Fields("B.id").Where(gdb.Raw("B.id=A.id"))
			m := db.Model(table1).As("A").Ctx(ctx).Where("NOT EXISTS ?", s).Limit(1)
			_, err := m.All()
			return err
		})
		t.AssertNil(err)
		t.Assert(expectSql, sql)
	})
	// https://github.com/gogf/gf/issues/3915
	gtest.C(t, func(t *gtest.T) {
		expectSql := "SELECT * FROM Test_Raw_Where_Table1 WHERE passport < nickname"
		sql, err := gdb.ToSQL(ctx, func(ctx context.Context) error {
			m := db.Model(table1).Ctx(ctx).WhereLT("passport", gdb.Raw("nickname"))
			_, err := m.All()
			return err
		})
		t.AssertNil(err)
		t.Assert(expectSql, sql)
	})
}

// Test_DataType_JSON_Insert tests JSON data insertion
func Test_DataType_JSON_Insert(t *testing.T) {
	table := "test_json_" + gtime.TimestampMicroStr()
	// Note: DM has no JSON column type, JSON documents are stored in CLOB columns.
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INT PRIMARY KEY AUTO_INCREMENT, data CLOB)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert simple JSON object
		result, err := db.Model(table).Data(g.Map{
			"data": `{"name":"John","age":30}`,
		}).Insert()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		// Verify data
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		expected := map[string]interface{}{"name": "John", "age": float64(30)}
		var actual map[string]interface{}
		err = json.Unmarshal([]byte(one["DATA"].String()), &actual)
		t.AssertNil(err)
		t.Assert(actual, expected)
	})
}

// Test_DataType_JSON_Extract tests JSON_EXTRACT function
func Test_DataType_JSON_Extract(t *testing.T) {
	table := "test_json_extract_" + gtime.TimestampMicroStr()
	// Note: DM has no JSON column type, JSON documents are stored in CLOB columns.
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INT PRIMARY KEY AUTO_INCREMENT, data CLOB)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert test data
		_, err := db.Model(table).Data(g.Map{
			"data": `{"name":"Alice","age":25,"city":"Beijing"}`,
		}).Insert()
		t.AssertNil(err)

		// Extract name using JSON_EXTRACT
		// Note: DM has no JSON_EXTRACT; JSON_VALUE returns the scalar itself, not its JSON text.
		one, err := db.Model(table).Fields("JSON_VALUE(data, '$.name') as name").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["NAME"].String(), `Alice`)

		// Extract age
		one, err = db.Model(table).Fields("JSON_VALUE(data, '$.age') as age").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["AGE"].Int(), 25)
	})
}

// Test_DataType_JSON_Array tests JSON array operations
func Test_DataType_JSON_Array(t *testing.T) {
	table := "test_json_array_" + gtime.TimestampMicroStr()
	// Note: DM has no JSON column type, JSON documents are stored in CLOB columns.
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INT PRIMARY KEY AUTO_INCREMENT, data CLOB)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert JSON array
		_, err := db.Model(table).Data(g.Map{
			"data": `["apple","banana","cherry"]`,
		}).Insert()
		t.AssertNil(err)

		// Extract array element
		// Note: DM has no JSON_EXTRACT; JSON_VALUE returns the scalar itself, not its JSON text.
		// FIRST is reserved in DM, so the alias is quoted.
		one, err := db.Model(table).Fields(`JSON_VALUE(data, '$[0]') as "FIRST"`).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["FIRST"].String(), `apple`)
	})
}

// Test_DataType_JSON_Null tests JSON NULL handling
func Test_DataType_JSON_Null(t *testing.T) {
	table := "test_json_null_" + gtime.TimestampMicroStr()
	// Note: DM has no JSON column type, JSON documents are stored in CLOB columns.
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INT PRIMARY KEY AUTO_INCREMENT, data CLOB)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert NULL value
		_, err := db.Model(table).Data(g.Map{
			"data": nil,
		}).Insert()
		t.AssertNil(err)

		// Verify NULL
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["DATA"].IsNil(), true)
	})
}

// Test_DataType_JSON_Complex tests complex nested JSON
func Test_DataType_JSON_Complex(t *testing.T) {
	table := "test_json_complex_" + gtime.TimestampMicroStr()
	// Note: DM has no JSON column type, JSON documents are stored in CLOB columns.
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INT PRIMARY KEY AUTO_INCREMENT, data CLOB)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert complex nested JSON
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

		// Extract nested field
		// Note: DM has no JSON_EXTRACT; JSON_VALUE returns the scalar itself, not its JSON text.
		one, err := db.Model(table).Fields("JSON_VALUE(data, '$.user.contacts.email') as email").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["EMAIL"].String(), `charlie@example.com`)
	})
}

// Test_DataType_JSON_Query tests JSON query with WHERE clause
func Test_DataType_JSON_Query(t *testing.T) {
	table := "test_json_query_" + gtime.TimestampMicroStr()
	// Note: DM has no JSON column type, JSON documents are stored in CLOB columns.
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INT PRIMARY KEY AUTO_INCREMENT, data CLOB)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert multiple JSON records
		_, err := db.Model(table).Data(g.List{
			g.Map{"data": `{"name":"David","age":20}`},
			g.Map{"data": `{"name":"Eve","age":30}`},
			g.Map{"data": `{"name":"Frank","age":25}`},
		}).Insert()
		t.AssertNil(err)

		// Query by JSON field value
		// Note: DM has no JSON_EXTRACT; JSON_VALUE extracts the scalar at the path.
		count, err := db.Model(table).Where("JSON_VALUE(data, '$.age') > ?", 25).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}

// Test_DataType_JSON_Update tests updating JSON data
func Test_DataType_JSON_Update(t *testing.T) {
	table := "test_json_update_" + gtime.TimestampMicroStr()
	// Note: DM has no JSON column type, JSON documents are stored in CLOB columns.
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INT PRIMARY KEY AUTO_INCREMENT, data CLOB)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert initial data
		_, err := db.Model(table).Data(g.Map{
			"data": `{"name":"Grace","age":28}`,
		}).Insert()
		t.AssertNil(err)

		// Update entire JSON
		_, err = db.Model(table).Data(g.Map{
			"data": `{"name":"Grace","age":29}`,
		}).Where("id", 1).Update()
		t.AssertNil(err)

		// Verify update
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		expected := map[string]interface{}{"name": "Grace", "age": float64(29)}
		var actual map[string]interface{}
		err = json.Unmarshal([]byte(one["DATA"].String()), &actual)
		t.AssertNil(err)
		t.Assert(actual, expected)
	})
}

// Test_DataType_Binary_Small tests small binary data
func Test_DataType_Binary_Small(t *testing.T) {
	table := "test_binary_small_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INT PRIMARY KEY AUTO_INCREMENT, data BLOB)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert small binary data
		binaryData := []byte{0x00, 0x01, 0x02, 0x03, 0xFF}
		_, err := db.Model(table).Data(g.Map{
			"data": binaryData,
		}).Insert()
		t.AssertNil(err)

		// Verify data
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(bytes.Equal(one["DATA"].Bytes(), binaryData), true)
	})
}

// Test_DataType_Binary_Large tests large binary data (1MB+)
func Test_DataType_Binary_Large(t *testing.T) {
	table := "test_binary_large_" + gtime.TimestampMicroStr()
	// Note: DM has no MEDIUMBLOB, its BLOB holds the 1MB value.
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INT PRIMARY KEY AUTO_INCREMENT, data BLOB)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Create 1MB binary data
		size := 1024 * 1024 // 1MB
		largeBinary := make([]byte, size)
		for i := 0; i < size; i++ {
			largeBinary[i] = byte(i % 256)
		}

		// Insert large binary data
		_, err := db.Model(table).Data(g.Map{
			"data": largeBinary,
		}).Insert()
		t.AssertNil(err)

		// Verify data
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(len(one["DATA"].Bytes()), size)
		t.Assert(bytes.Equal(one["DATA"].Bytes(), largeBinary), true)
	})
}

// Test_DataType_Binary_Integrity tests binary data integrity with checksum
func Test_DataType_Binary_Integrity(t *testing.T) {
	table := "test_binary_integrity_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INT PRIMARY KEY AUTO_INCREMENT, data BLOB, checksum VARCHAR(64))")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Create random binary data
		binaryData := []byte("Hello, World! This is a binary test data with special chars: \x00\xFF\xAB")

		// Calculate SHA256 checksum
		hash := sha256.Sum256(binaryData)
		checksum := hex.EncodeToString(hash[:])

		// Insert with checksum
		_, err := db.Model(table).Data(g.Map{
			"data":     binaryData,
			"checksum": checksum,
		}).Insert()
		t.AssertNil(err)

		// Verify integrity
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)

		retrievedHash := sha256.Sum256(one["DATA"].Bytes())
		retrievedChecksum := hex.EncodeToString(retrievedHash[:])
		t.Assert(retrievedChecksum, checksum)
	})
}

// Test_DataType_Binary_Empty tests empty and NULL binary
func Test_DataType_Binary_Empty(t *testing.T) {
	table := "test_binary_empty_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INT PRIMARY KEY AUTO_INCREMENT, data BLOB)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert empty binary
		_, err := db.Model(table).Data(g.Map{
			"data": []byte{},
		}).Insert()
		t.AssertNil(err)

		// Insert NULL
		_, err = db.Model(table).Data(g.Map{
			"data": nil,
		}).Insert()
		t.AssertNil(err)

		// Verify empty
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(len(one["DATA"].Bytes()), 0)

		// Verify NULL
		one, err = db.Model(table).Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(one["DATA"].IsNil(), true)
	})
}

// Test_DataType_Decimal_HighPrecision tests high precision decimal (65,30)
func Test_DataType_Decimal_HighPrecision(t *testing.T) {
	table := "test_decimal_precision_" + gtime.TimestampMicroStr()
	// Note: DM allows a precision of at most 38, so the column is DECIMAL(38,30) and the value keeps
	// the 30 fractional digits with 8 integer digits.
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INT PRIMARY KEY AUTO_INCREMENT, amount DECIMAL(38,30))")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert high precision decimal
		value := "12345678.123456789012345678901234567890"
		_, err := db.Model(table).Data(g.Map{
			"amount": value,
		}).Insert()
		t.AssertNil(err)

		// Verify precision
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["AMOUNT"].String(), value)
	})
}

// Test_DataType_Decimal_Calculation tests decimal arithmetic
func Test_DataType_Decimal_Calculation(t *testing.T) {
	table := "test_decimal_calc_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INT PRIMARY KEY AUTO_INCREMENT, price DECIMAL(10,2), quantity DECIMAL(10,2))")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert test data
		_, err := db.Model(table).Data(g.Map{
			"price":    "19.99",
			"quantity": "3.5",
		}).Insert()
		t.AssertNil(err)

		// Calculate total using SQL
		// Note: DM gives the product of two DECIMAL(10,2) no fixed scale (DEC(38,0) in the result
		// metadata), where MySQL types it DECIMAL(21,4), so there is no trailing zero.
		one, err := db.Model(table).Fields("price * quantity as total").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["TOTAL"].String(), "69.965")
	})
}

// Test_DataType_Decimal_Boundary tests decimal boundary values
func Test_DataType_Decimal_Boundary(t *testing.T) {
	table := "test_decimal_boundary_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INT PRIMARY KEY AUTO_INCREMENT, value DECIMAL(10,2))")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Test max value (10 digits, 2 decimals: 99999999.99)
		_, err := db.Model(table).Data(g.Map{
			"value": "99999999.99",
		}).Insert()
		t.AssertNil(err)

		// Test min value
		_, err = db.Model(table).Data(g.Map{
			"value": "-99999999.99",
		}).Insert()
		t.AssertNil(err)

		// Test zero
		_, err = db.Model(table).Data(g.Map{
			"value": "0.00",
		}).Insert()
		t.AssertNil(err)

		// Verify all values
		all, err := db.Model(table).Order("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(all[0]["VALUE"].String(), "99999999.99")
		t.Assert(all[1]["VALUE"].String(), "-99999999.99")
		t.Assert(all[2]["VALUE"].String(), "0.00")
	})
}

// Test_DataType_Decimal_Null tests NULL decimal values
func Test_DataType_Decimal_Null(t *testing.T) {
	table := "test_decimal_null_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INT PRIMARY KEY AUTO_INCREMENT, value DECIMAL(10,2))")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert NULL
		_, err := db.Model(table).Data(g.Map{
			"value": nil,
		}).Insert()
		t.AssertNil(err)

		// Verify NULL
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["VALUE"].IsNil(), true)
	})
}

// Test_DataType_Datetime_Timezone tests datetime with timezone handling
func Test_DataType_Datetime_Timezone(t *testing.T) {
	table := "test_datetime_tz_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INT PRIMARY KEY AUTO_INCREMENT, created_at DATETIME)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert datetime
		dt := "2024-01-15 12:30:45"
		_, err := db.Model(table).Data(g.Map{
			"created_at": dt,
		}).Insert()
		t.AssertNil(err)

		// Verify datetime
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["CREATED_AT"].String(), dt)
	})
}

// Test_DataType_Datetime_Precision tests datetime with microsecond precision
func Test_DataType_Datetime_Precision(t *testing.T) {
	table := "test_datetime_precision_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INT PRIMARY KEY AUTO_INCREMENT, created_at DATETIME(6))")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert datetime with microseconds
		dt := "2024-01-15 12:30:45.123456"
		_, err := db.Model(table).Data(g.Map{
			"created_at": dt,
		}).Insert()
		t.AssertNil(err)

		// Verify precision (compare up to seconds, MySQL may format microseconds differently)
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		expected := "2024-01-15 12:30:45"
		actual := one["CREATED_AT"].String()[:19] // Extract first 19 chars (YYYY-MM-DD HH:MM:SS)
		t.Assert(actual, expected)
	})
}

// Test_DataType_Datetime_Boundary tests datetime boundary values
func Test_DataType_Datetime_Boundary(t *testing.T) {
	table := "test_datetime_boundary_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INT PRIMARY KEY AUTO_INCREMENT, dt DATETIME)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Test min datetime (MySQL supports 1000-01-01 00:00:00)
		_, err := db.Model(table).Data(g.Map{
			"dt": "1000-01-01 00:00:00",
		}).Insert()
		t.AssertNil(err)

		// Test max datetime
		_, err = db.Model(table).Data(g.Map{
			"dt": "9999-12-31 23:59:59",
		}).Insert()
		t.AssertNil(err)

		// Verify boundaries
		all, err := db.Model(table).Order("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["DT"].String(), "1000-01-01 00:00:00")
		t.Assert(all[1]["DT"].String(), "9999-12-31 23:59:59")
	})
}

// Test_DataType_Datetime_Null tests NULL datetime
func Test_DataType_Datetime_Null(t *testing.T) {
	table := "test_datetime_null_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INT PRIMARY KEY AUTO_INCREMENT, dt DATETIME)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert NULL
		_, err := db.Model(table).Data(g.Map{
			"dt": nil,
		}).Insert()
		t.AssertNil(err)

		// Verify NULL
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["DT"].IsNil(), true)
	})
}

// Test_DataType_Datetime_Update tests datetime updates
func Test_DataType_Datetime_Update(t *testing.T) {
	table := "test_datetime_update_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INT PRIMARY KEY AUTO_INCREMENT, dt DATETIME)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert initial datetime
		dt1 := "2024-01-01 10:00:00"
		_, err := db.Model(table).Data(g.Map{
			"dt": dt1,
		}).Insert()
		t.AssertNil(err)

		// Update datetime
		dt2 := "2024-12-31 23:59:59"
		_, err = db.Model(table).Data(g.Map{
			"dt": dt2,
		}).Where("id", 1).Update()
		t.AssertNil(err)

		// Verify update
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["DT"].String(), dt2)
	})
}

// Test_DataType_Enum_Valid tests valid ENUM values
func Test_DataType_Enum_Valid(t *testing.T) {
	table := "test_enum_valid_" + gtime.TimestampMicroStr()
	// Note: DM has no ENUM type; a CHECK constraint restricts the column to the same values.
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INT PRIMARY KEY AUTO_INCREMENT, status VARCHAR(20) CHECK (status IN ('pending','approved','rejected')))")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert all valid values
		_, err := db.Model(table).Data(g.List{
			g.Map{"status": "pending"},
			g.Map{"status": "approved"},
			g.Map{"status": "rejected"},
		}).Insert()
		t.AssertNil(err)

		// Verify all values
		all, err := db.Model(table).Order("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(all[0]["STATUS"].String(), "pending")
		t.Assert(all[1]["STATUS"].String(), "approved")
		t.Assert(all[2]["STATUS"].String(), "rejected")
	})
}

// Test_DataType_Enum_Invalid tests invalid ENUM values (should fail or truncate)
func Test_DataType_Enum_Invalid(t *testing.T) {
	table := "test_enum_invalid_" + gtime.TimestampMicroStr()
	// Note: DM has no ENUM type; a CHECK constraint restricts the column to the same values.
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INT PRIMARY KEY AUTO_INCREMENT, status VARCHAR(20) CHECK (status IN ('pending','approved','rejected')))")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Attempt to insert invalid value (should fail in strict mode)
		_, err := db.Model(table).Data(g.Map{
			"status": "invalid_status",
		}).Insert()
		// In strict SQL mode, this should produce an error
		// In non-strict mode, it might insert empty string
		t.AssertNE(err, nil)
	})
}

// Test_DataType_Set_Valid tests valid SET values
func Test_DataType_Set_Valid(t *testing.T) {
	table := "test_set_valid_" + gtime.TimestampMicroStr()
	// Note: DM has no SET type; the column holds the comma-separated members, which a CHECK
	// constraint restricts to the combinations MySQL stores.
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INT PRIMARY KEY AUTO_INCREMENT, permissions VARCHAR(32) CHECK (permissions IN "+
		"('','read','write','execute','read,write','read,execute','write,execute','read,write,execute')))")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert single value
		_, err := db.Model(table).Data(g.Map{
			"permissions": "read",
		}).Insert()
		t.AssertNil(err)

		// Insert multiple values
		_, err = db.Model(table).Data(g.Map{
			"permissions": "read,write",
		}).Insert()
		t.AssertNil(err)

		// Insert all values
		_, err = db.Model(table).Data(g.Map{
			"permissions": "read,write,execute",
		}).Insert()
		t.AssertNil(err)

		// Verify all values
		all, err := db.Model(table).Order("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(all[0]["PERMISSIONS"].String(), "read")
		t.Assert(all[1]["PERMISSIONS"].String(), "read,write")
		t.Assert(all[2]["PERMISSIONS"].String(), "read,write,execute")
	})
}

// Test_DataType_Set_Empty tests empty SET values
func Test_DataType_Set_Empty(t *testing.T) {
	table := "test_set_empty_" + gtime.TimestampMicroStr()
	// Note: DM has no SET type; the column holds the comma-separated members, which a CHECK
	// constraint restricts to the combinations MySQL stores.
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INT PRIMARY KEY AUTO_INCREMENT, permissions VARCHAR(32) CHECK (permissions IN "+
		"('','read','write','execute','read,write','read,execute','write,execute','read,write,execute')))")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert empty SET
		_, err := db.Model(table).Data(g.Map{
			"permissions": "",
		}).Insert()
		t.AssertNil(err)

		// Verify empty
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["PERMISSIONS"].String(), "")
	})
}

// Test_DataType_Geometry_Point tests POINT geometry type
func Test_DataType_Geometry_Point(t *testing.T) {
	table := "test_geo_point_" + gtime.TimestampMicroStr()
	// Note: DM provides the spatial types and functions in the DMGEO package, initialized once by
	// SP_INIT_GEO_SYS.
	rawTypeInitGeometry(t)
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INT PRIMARY KEY AUTO_INCREMENT, location SYSGEO.ST_POINT)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert POINT using ST_GeomFromText
		_, err := db.Exec(ctx, fmt.Sprintf("INSERT INTO %s (location) VALUES (DMGEO.ST_GeomFromText('POINT(116.4074 39.9042)', 0))", table))
		t.AssertNil(err)

		// Query POINT using ST_AsText
		// Note: DM writes WKT with a space after the type and every coordinate of the double value
		// with 16 decimal places.
		one, err := db.Model(table).Fields("DMGEO.ST_AsText(location) as location_text").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["LOCATION_TEXT"].String(), "POINT (116.4073999999999955 39.9042000000000030)")
	})
}

// Test_DataType_Geometry_Polygon tests POLYGON geometry type
func Test_DataType_Geometry_Polygon(t *testing.T) {
	table := "test_geo_polygon_" + gtime.TimestampMicroStr()
	// Note: DM provides the spatial types and functions in the DMGEO package, initialized once by
	// SP_INIT_GEO_SYS.
	rawTypeInitGeometry(t)
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INT PRIMARY KEY AUTO_INCREMENT, area SYSGEO.ST_POLYGON)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert POLYGON (rectangle)
		polygon := "POLYGON((0 0, 10 0, 10 10, 0 10, 0 0))"
		_, err := db.Exec(ctx, fmt.Sprintf("INSERT INTO %s (area) VALUES (DMGEO.ST_GeomFromText('%s', 0))", table, polygon))
		t.AssertNil(err)

		// Query POLYGON (normalize spaces for comparison)
		one, err := db.Model(table).Fields("DMGEO.ST_AsText(area) as area_text").Where("id", 1).One()
		t.AssertNil(err)
		// Note: DM writes WKT with a space after the type and every coordinate with 16 decimal places.
		expected := "POLYGON ((0.0000000000000000 0.0000000000000000,10.0000000000000000 0.0000000000000000," +
			"10.0000000000000000 10.0000000000000000,0.0000000000000000 10.0000000000000000,0.0000000000000000 0.0000000000000000))"
		actual := strings.ReplaceAll(one["AREA_TEXT"].String(), ", ", ",") // Remove spaces after commas
		t.Assert(actual, expected)
	})
}

// Test_DataType_Geometry_Null tests NULL geometry values
func Test_DataType_Geometry_Null(t *testing.T) {
	table := "test_geo_null_" + gtime.TimestampMicroStr()
	// Note: DM provides the spatial types and functions in the DMGEO package, initialized once by
	// SP_INIT_GEO_SYS.
	rawTypeInitGeometry(t)
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INT PRIMARY KEY AUTO_INCREMENT, location SYSGEO.ST_POINT)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert NULL
		_, err := db.Model(table).Data(g.Map{
			"location": nil,
		}).Insert()
		t.AssertNil(err)

		// Verify NULL
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["LOCATION"].IsNil(), true)
	})
}
