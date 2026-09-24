// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package sqlitecgo_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

func Test_Raw_Insert(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		user := db.Model(table)
		result, err := user.Data(g.Map{
			"id":          gdb.Raw("2"),
			"passport":    "port_1",
			"password":    "pass_1",
			"nickname":    "name_1",
			"create_time": gdb.Raw("datetime('now')"),
		}).Insert()
		t.AssertNil(err)
		n, _ := result.LastInsertId()
		t.Assert(n, 2)

		one, err := db.Model(table).Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(one["passport"].String(), "port_1")
		t.AssertNE(one["create_time"].String(), "")
	})
}

func Test_Raw_BatchInsert(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		user := db.Model(table)
		result, err := user.Data(
			g.List{
				g.Map{
					"id":          gdb.Raw("2"),
					"passport":    "port_2",
					"password":    "pass_2",
					"nickname":    "name_2",
					"create_time": gdb.Raw("datetime('now')"),
				},
				g.Map{
					"id":          gdb.Raw("4"),
					"passport":    "port_4",
					"password":    "pass_4",
					"nickname":    "name_4",
					"create_time": gdb.Raw("datetime('now')"),
				},
			},
		).Insert()
		t.AssertNil(err)
		n, _ := result.LastInsertId()
		t.Assert(n, 4)

		all, err := db.Model(table).Order("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["id"].Int(), 2)
		t.Assert(all[1]["id"].Int(), 4)
	})
}

func Test_Raw_Update(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		user := db.Model(table)
		result, err := user.Data(g.Map{
			"id":          gdb.Raw("id+100"),
			"create_time": gdb.Raw("datetime('now')"),
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
	table1 := createTable("Test_Raw_Where_Table1")
	table2 := createTable("Test_Raw_Where_Table2")
	defer dropTable(table1)
	defer dropTable(table2)

	// https://github.com/gogf/gf/issues/3922
	gtest.C(t, func(t *gtest.T) {
		expectSql := "SELECT * FROM `Test_Raw_Where_Table1` AS A WHERE NOT EXISTS (SELECT B.id FROM `Test_Raw_Where_Table2` AS B WHERE `B`.`id`=A.id) LIMIT 1"
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
		expectSql := "SELECT * FROM `Test_Raw_Where_Table1` AS A WHERE NOT EXISTS (SELECT B.id FROM `Test_Raw_Where_Table2` AS B WHERE B.id=A.id) LIMIT 1"
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
		expectSql := "SELECT * FROM `Test_Raw_Where_Table1` WHERE `passport` < `nickname`"
		sql, err := gdb.ToSQL(ctx, func(ctx context.Context) error {
			m := db.Model(table1).Ctx(ctx).WhereLT("passport", gdb.Raw("`nickname`"))
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
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INTEGER PRIMARY KEY AUTOINCREMENT, data JSON)")
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

		valid, err := db.Model(table).Fields("json_valid(data) as ok").Where("id", 1).Value()
		t.AssertNil(err)
		t.Assert(valid.Int(), 1)
	})
}

// Test_DataType_JSON_Extract tests json_extract function
func Test_DataType_JSON_Extract(t *testing.T) {
	table := "test_json_extract_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INTEGER PRIMARY KEY AUTOINCREMENT, data JSON)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"data": `{"name":"Alice","age":25,"city":"Beijing"}`,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Fields("json_extract(data, '$.name') as name").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["name"].String(), "Alice")

		one, err = db.Model(table).Fields("json_quote(json_extract(data, '$.name')) as name").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["name"].String(), `"Alice"`)

		one, err = db.Model(table).Fields("json_extract(data, '$.age') as age").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["age"].Int(), 25)
	})
}

// Test_DataType_JSON_Set tests json_set function
func Test_DataType_JSON_Set(t *testing.T) {
	table := "test_json_set_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INTEGER PRIMARY KEY AUTOINCREMENT, data JSON)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"data": `{"name":"Bob"}`,
		}).Insert()
		t.AssertNil(err)

		_, err = db.Exec(ctx, fmt.Sprintf("UPDATE %s SET data = json_set(data, '$.age', 30) WHERE id = 1", table))
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

// Test_DataType_JSON_Array tests JSON array operations
func Test_DataType_JSON_Array(t *testing.T) {
	table := "test_json_array_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INTEGER PRIMARY KEY AUTOINCREMENT, data JSON)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"data": `["apple","banana","cherry"]`,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Fields("json_extract(data, '$[0]') as first").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["first"].String(), "apple")

		one, err = db.Model(table).Fields("json_array_length(data) as n").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["n"].Int(), 3)
	})
}

// Test_DataType_JSON_Null tests JSON NULL handling
func Test_DataType_JSON_Null(t *testing.T) {
	table := "test_json_null_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INTEGER PRIMARY KEY AUTOINCREMENT, data JSON)")
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
		t.Assert(one["data"].Val(), nil)
	})
}

// Test_DataType_JSON_Complex tests complex nested JSON
func Test_DataType_JSON_Complex(t *testing.T) {
	table := "test_json_complex_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INTEGER PRIMARY KEY AUTOINCREMENT, data JSON)")
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

		one, err := db.Model(table).Fields("json_extract(data, '$.user.contacts.email') as email").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["email"].String(), "charlie@example.com")

		one, err = db.Model(table).Fields("json_extract(data, '$.user.tags[1]') as tag").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["tag"].String(), "gopher")
	})
}

// Test_DataType_JSON_Query tests JSON query with WHERE clause
func Test_DataType_JSON_Query(t *testing.T) {
	table := "test_json_query_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INTEGER PRIMARY KEY AUTOINCREMENT, data JSON)")
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

		count, err := db.Model(table).Where("json_extract(data, '$.age') > ?", 25).Count()
		t.AssertNil(err)
		t.Assert(count, 1)

		count, err = db.Model(table).Where("json_extract(data, '$.name') = ?", "Eve").Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}

// Test_DataType_JSON_Update tests updating JSON data
func Test_DataType_JSON_Update(t *testing.T) {
	table := "test_json_update_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INTEGER PRIMARY KEY AUTOINCREMENT, data JSON)")
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

// Test_DataType_Binary_Small tests small binary data
func Test_DataType_Binary_Small(t *testing.T) {
	table := "test_binary_small_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INTEGER PRIMARY KEY AUTOINCREMENT, data BLOB)")
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

		storage, err := db.Model(table).Fields("typeof(data) as tp").Where("id", 1).Value()
		t.AssertNil(err)
		t.Assert(storage.String(), "blob")
	})
}

// Test_DataType_Binary_Large tests large binary data (1MB+)
func Test_DataType_Binary_Large(t *testing.T) {
	table := "test_binary_large_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INTEGER PRIMARY KEY AUTOINCREMENT, data BLOB)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		size := 1024 * 1024
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

// Test_DataType_Binary_Integrity tests binary data integrity with checksum
func Test_DataType_Binary_Integrity(t *testing.T) {
	table := "test_binary_integrity_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INTEGER PRIMARY KEY AUTOINCREMENT, data BLOB, checksum VARCHAR(64))")
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
		t.Assert(one["checksum"].String(), checksum)
	})
}

// Test_DataType_Binary_Empty tests empty and NULL binary
func Test_DataType_Binary_Empty(t *testing.T) {
	table := "test_binary_empty_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INTEGER PRIMARY KEY AUTOINCREMENT, data BLOB)")
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
		t.Assert(one["data"].Val(), nil)

		all, err := db.Model(table).Fields("typeof(data) as tp").Order("id").All()
		t.AssertNil(err)
		t.Assert(all[0]["tp"].String(), "blob")
		t.Assert(all[1]["tp"].String(), "null")
	})
}

// Test_DataType_Decimal_HighPrecision tests high precision decimal stored in a NUMERIC affinity column
func Test_DataType_Decimal_HighPrecision(t *testing.T) {
	table := "test_decimal_precision_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INTEGER PRIMARY KEY AUTOINCREMENT, amount DECIMAL(65,30), amount_text TEXT)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		value := "12345678901234567890123456789012345.123456789012345678901234567890"
		_, err := db.Model(table).Data(g.Map{
			"amount":      value,
			"amount_text": value,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Fields("amount, amount_text, typeof(amount) as tp").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["tp"].String(), "real")
		t.AssertNE(one["amount"].String(), value)
		t.AssertGT(one["amount"].Float64(), 1.2e34)
		t.AssertLT(one["amount"].Float64(), 1.3e34)
		t.Assert(one["amount_text"].String(), value)
	})
}

// Test_DataType_Decimal_Calculation tests decimal arithmetic
func Test_DataType_Decimal_Calculation(t *testing.T) {
	table := "test_decimal_calc_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INTEGER PRIMARY KEY AUTOINCREMENT, price DECIMAL(10,2), quantity DECIMAL(10,2))")
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

		one, err := db.Model(table).Fields("price * quantity as total").Where("id", 1).One()
		t.AssertNil(err)
		t.AssertLT(math.Abs(one["total"].Float64()-69.965), 1e-9)

		one, err = db.Model(table).Fields("printf('%.4f', price * quantity) as total").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["total"].String(), "69.9650")
	})
}

// Test_DataType_Decimal_Boundary tests decimal boundary values
func Test_DataType_Decimal_Boundary(t *testing.T) {
	table := "test_decimal_boundary_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INTEGER PRIMARY KEY AUTOINCREMENT, value DECIMAL(10,2))")
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
		t.Assert(all[2]["value"].String(), "0")
		t.Assert(all[2]["value"].Float64(), 0)

		types, err := db.Model(table).Fields("typeof(value) as tp").Order("id").All()
		t.AssertNil(err)
		t.Assert(types[0]["tp"].String(), "real")
		t.Assert(types[2]["tp"].String(), "integer")
	})
}

// Test_DataType_Decimal_Null tests NULL decimal values
func Test_DataType_Decimal_Null(t *testing.T) {
	table := "test_decimal_null_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INTEGER PRIMARY KEY AUTOINCREMENT, value DECIMAL(10,2))")
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
		t.Assert(one["value"].Val(), nil)
	})
}

// Test_DataType_Datetime_Timezone tests datetime with timezone handling
func Test_DataType_Datetime_Timezone(t *testing.T) {
	table := "test_datetime_tz_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INTEGER PRIMARY KEY AUTOINCREMENT, created_at DATETIME)")
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

		stored, err := db.Model(table).Fields("typeof(created_at) as tp").Where("id", 1).Value()
		t.AssertNil(err)
		t.Assert(stored.String(), "text")
	})
}

// Test_DataType_Datetime_Precision tests datetime with microsecond precision
func Test_DataType_Datetime_Precision(t *testing.T) {
	table := "test_datetime_precision_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INTEGER PRIMARY KEY AUTOINCREMENT, created_at DATETIME)")
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

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		expected := "2024-01-15 12:30:45"
		actual := one["created_at"].String()[:19]
		t.Assert(actual, expected)
		t.Assert(one["created_at"].GTime().Microsecond(), 123456)

		raw, err := db.GetValue(ctx, "SELECT CAST(created_at AS TEXT) FROM "+table+" WHERE id = 1")
		t.AssertNil(err)
		t.Assert(raw.String(), dt)
	})
}

// Test_DataType_Datetime_Boundary tests datetime boundary values
func Test_DataType_Datetime_Boundary(t *testing.T) {
	table := "test_datetime_boundary_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INTEGER PRIMARY KEY AUTOINCREMENT, dt DATETIME)")
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

// Test_DataType_Datetime_Null tests NULL datetime
func Test_DataType_Datetime_Null(t *testing.T) {
	table := "test_datetime_null_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INTEGER PRIMARY KEY AUTOINCREMENT, dt DATETIME)")
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
		t.Assert(one["dt"].Val(), nil)
	})
}

// Test_DataType_Datetime_Update tests datetime updates
func Test_DataType_Datetime_Update(t *testing.T) {
	table := "test_datetime_update_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INTEGER PRIMARY KEY AUTOINCREMENT, dt DATETIME)")
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

// Test_DataType_Enum_Valid tests valid enumerated values constrained by CHECK
func Test_DataType_Enum_Valid(t *testing.T) {
	table := "test_enum_valid_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INTEGER PRIMARY KEY AUTOINCREMENT, status TEXT CHECK (status IN ('pending','approved','rejected')))")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.List{
			g.Map{"status": "pending"},
			g.Map{"status": "approved"},
			g.Map{"status": "rejected"},
		}).Insert()
		t.AssertNil(err)

		all, err := db.Model(table).Order("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(all[0]["status"].String(), "pending")
		t.Assert(all[1]["status"].String(), "approved")
		t.Assert(all[2]["status"].String(), "rejected")
	})
}

// Test_DataType_Enum_Invalid tests values rejected by the CHECK constraint
func Test_DataType_Enum_Invalid(t *testing.T) {
	table := "test_enum_invalid_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INTEGER PRIMARY KEY AUTOINCREMENT, status TEXT CHECK (status IN ('pending','approved','rejected')))")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"status": "invalid_status",
		}).Insert()
		t.AssertNE(err, nil)
		t.Assert(gstr.Contains(err.Error(), "CHECK constraint failed"), true)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}

// Test_DataType_Set_Valid tests multi-value flags stored as a delimited TEXT column
func Test_DataType_Set_Valid(t *testing.T) {
	table := "test_set_valid_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INTEGER PRIMARY KEY AUTOINCREMENT, permissions TEXT)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"permissions": "read",
		}).Insert()
		t.AssertNil(err)

		_, err = db.Model(table).Data(g.Map{
			"permissions": "read,write",
		}).Insert()
		t.AssertNil(err)

		_, err = db.Model(table).Data(g.Map{
			"permissions": "read,write,execute",
		}).Insert()
		t.AssertNil(err)

		all, err := db.Model(table).Order("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(all[0]["permissions"].String(), "read")
		t.Assert(all[1]["permissions"].String(), "read,write")
		t.Assert(all[2]["permissions"].String(), "read,write,execute")

		count, err := db.Model(table).Where("permissions LIKE ?", "%write%").Count()
		t.AssertNil(err)
		t.Assert(count, 2)
	})
}

// Test_DataType_Set_Empty tests empty multi-value flags
func Test_DataType_Set_Empty(t *testing.T) {
	table := "test_set_empty_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INTEGER PRIMARY KEY AUTOINCREMENT, permissions TEXT)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"permissions": "",
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["permissions"].String(), "")
		t.Assert(one["permissions"].IsNil(), false)
	})
}

// Test_DataType_Geometry_Point tests WKT point text stored in a POINT declared column
func Test_DataType_Geometry_Point(t *testing.T) {
	table := "test_geo_point_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INTEGER PRIMARY KEY AUTOINCREMENT, location POINT)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		point := "POINT(116.4074 39.9042)"
		_, err := db.Model(table).Data(g.Map{
			"location": point,
		}).Insert()
		t.AssertNil(err)

		storage, err := db.Model(table).Fields("typeof(location) as tp").Where("id", 1).Value()
		t.AssertNil(err)
		t.Assert(storage.String(), "text")

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		if got := one["location"].String(); got != point {
			t.Errorf(
				`Known bug, see https://github.com/gogf/gf/issues/4842: declared type POINT contains "int", so gdb.Core.CheckLocalTypeForField reports `+
					`LocalTypeInt and the TEXT storage class value is converted away; got %q, want %q`,
				got, point,
			)
		}
	})
}

// Test_DataType_Geometry_Polygon tests WKT polygon text stored in a POLYGON declared column
func Test_DataType_Geometry_Polygon(t *testing.T) {
	table := "test_geo_polygon_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INTEGER PRIMARY KEY AUTOINCREMENT, area POLYGON)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		polygon := "POLYGON((0 0,10 0,10 10,0 10,0 0))"
		_, err := db.Model(table).Data(g.Map{
			"area": polygon,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["area"].String(), polygon)

		storage, err := db.Model(table).Fields("typeof(area) as tp").Where("id", 1).Value()
		t.AssertNil(err)
		t.Assert(storage.String(), "text")
	})
}

// Test_DataType_Geometry_Null tests NULL geometry values
func Test_DataType_Geometry_Null(t *testing.T) {
	table := "test_geo_null_" + gtime.TimestampMicroStr()
	_, err := db.Exec(ctx, "CREATE TABLE "+table+" (id INTEGER PRIMARY KEY AUTOINCREMENT, location POINT)")
	if err != nil {
		t.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"location": nil,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["location"].IsNil(), true)
		t.Assert(one["location"].Val(), nil)
	})
}
