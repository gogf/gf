// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

// chCCreateTypeTable creates a table with an `id UInt64` sorting key followed by the given columns.
func chCCreateTypeTable(prefix string, columns string) string {
	name := fmt.Sprintf(`%s_%d`, prefix, gtime.TimestampNano())
	dropTable(name)
	if _, err := db.Exec(ctx, fmt.Sprintf(
		"CREATE TABLE %s (id UInt64, %s) ENGINE = MergeTree() ORDER BY id", name, columns,
	)); err != nil {
		gtest.Fatal(err)
	}
	return name
}

// Test_Raw_Insert tests inserting a record with raw SQL expressions as values.
// Note: ClickHouse cannot refer to a column in an inserted value as MySQL does with `id+2`,
// and reports no last insert id, so a constant expression is inserted and the record is checked.
func Test_Raw_Insert(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		user := db.Model(table)
		_, err := user.Data(g.Map{
			"id":          gdb.Raw("0+2"),
			"passport":    "port_1",
			"password":    "pass_1",
			"nickname":    "name_1",
			"create_time": gdb.Raw("now()"),
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("passport", "port_1").One()
		t.AssertNil(err)
		t.Assert(one["id"], 2)
		t.AssertGT(one["create_time"].GTime().Year(), 1970)
	})
}

// Test_Raw_BatchInsert tests inserting a batch of records with raw SQL expressions as values.
// Note: see Test_Raw_Insert for the constant expressions inserted in place of `id+2`.
func Test_Raw_BatchInsert(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		user := db.Model(table)
		_, err := user.Data(
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

		array, err := db.Model(table).OrderAsc("id").Array("id")
		t.AssertNil(err)
		t.Assert(array, g.Slice{2, 4})
	})
}

// Test_Raw_Update tests updating a record with raw SQL expressions as values.
// Note: ClickHouse cannot update a column of the sorting key, so `id+100` fails, and a
// mutation reports no affected rows, so the updated record is checked instead.
func Test_Raw_Update(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id":          gdb.Raw("id+100"),
			"create_time": gdb.Raw("now()"),
		}).Where("id", 1).Update()
		t.AssertNE(err, nil)
	})
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"nickname":    gdb.Raw("concat(nickname, '_100')"),
			"create_time": gdb.Raw("now()"),
		}).Where("id", 1).Update()
		t.AssertNil(err)
	})
	gtest.C(t, func(t *gtest.T) {
		n, err := db.Model(table).Where("nickname", "name_1_100").Count()
		t.AssertNil(err)
		t.Assert(n, 1)
	})
}

// Test_Raw_Where tests the SQL built with raw SQL expressions in conditions.
func Test_Raw_Where(t *testing.T) {
	table1 := createTable("Test_Raw_Where_Table1_" + gtime.TimestampNanoStr())
	table2 := createTable("Test_Raw_Where_Table2_" + gtime.TimestampNanoStr())
	defer dropTable(table1)
	defer dropTable(table2)

	// https://github.com/gogf/gf/issues/3922
	gtest.C(t, func(t *gtest.T) {
		expectSql := fmt.Sprintf(
			"SELECT * FROM %s AS A WHERE NOT EXISTS (SELECT B.id FROM %s AS B WHERE B.id=A.id) LIMIT 1",
			table1, table2,
		)
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
		expectSql := fmt.Sprintf(
			"SELECT * FROM %s AS A WHERE NOT EXISTS (SELECT B.id FROM %s AS B WHERE B.id=A.id) LIMIT 1",
			table1, table2,
		)
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
		expectSql := fmt.Sprintf("SELECT * FROM %s WHERE passport < `nickname`", table1)
		sql, err := gdb.ToSQL(ctx, func(ctx context.Context) error {
			m := db.Model(table1).Ctx(ctx).WhereLT("passport", gdb.Raw("`nickname`"))
			_, err := m.All()
			return err
		})
		t.AssertNil(err)
		t.Assert(expectSql, sql)
	})
}

// Test_DataType_JSON_Insert tests JSON data insertion into a String column.
func Test_DataType_JSON_Insert(t *testing.T) {
	table := chCCreateTypeTable("test_json", "data String")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id":   1,
			"data": `{"name":"John","age":30}`,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		expected := map[string]any{"name": "John", "age": float64(30)}
		var actual map[string]any
		err = json.Unmarshal([]byte(one["data"].String()), &actual)
		t.AssertNil(err)
		t.Assert(actual, expected)
	})
}

// Test_DataType_JSON_Extract tests JSONExtractRaw, the counterpart of JSON_EXTRACT.
func Test_DataType_JSON_Extract(t *testing.T) {
	table := chCCreateTypeTable("test_json_extract", "data String")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id":   1,
			"data": `{"name":"Alice","age":25,"city":"Beijing"}`,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Fields("JSONExtractRaw(data, 'name') as name").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["name"].String(), `"Alice"`)

		one, err = db.Model(table).Fields("JSONExtractRaw(data, 'age') as age").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["age"].Int(), 25)
	})
}

// Test_DataType_JSON_Set tests setting a JSON field with JSONMergePatch, the counterpart of JSON_SET.
// Note: ClickHouse 24.11 has no UPDATE statement, so the mutation is written as ALTER TABLE ... UPDATE.
func Test_DataType_JSON_Set(t *testing.T) {
	table := chCCreateTypeTable("test_json_set", "data String")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id":   1,
			"data": `{"name":"Bob"}`,
		}).Insert()
		t.AssertNil(err)

		_, err = db.Exec(ctx, fmt.Sprintf(
			`ALTER TABLE %s UPDATE data = JSONMergePatch(data, '{"age":30}') WHERE id = 1`, table,
		))
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		expected := map[string]any{"name": "Bob", "age": float64(30)}
		var actual map[string]any
		err = json.Unmarshal([]byte(one["data"].String()), &actual)
		t.AssertNil(err)
		t.Assert(actual, expected)
	})
}

// Test_DataType_JSON_Array tests extracting an element of a JSON array.
// Note: the indexes of JSONExtractRaw start at 1, where the paths of JSON_EXTRACT start at 0.
func Test_DataType_JSON_Array(t *testing.T) {
	table := chCCreateTypeTable("test_json_array", "data String")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id":   1,
			"data": `["apple","banana","cherry"]`,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Fields("JSONExtractRaw(data, 1) as first").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["first"].String(), `"apple"`)
	})
}

// Test_DataType_JSON_Null tests NULL in a nullable JSON column.
func Test_DataType_JSON_Null(t *testing.T) {
	table := chCCreateTypeTable("test_json_null", "data Nullable(String)")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id":   1,
			"data": nil,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["data"].IsNil(), true)
	})
}

// Test_DataType_JSON_Complex tests extracting a nested JSON field.
func Test_DataType_JSON_Complex(t *testing.T) {
	table := chCCreateTypeTable("test_json_complex", "data String")
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
			"id":   1,
			"data": complexJSON,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Fields(
			"JSONExtractRaw(data, 'user', 'contacts', 'email') as email",
		).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["email"].String(), `"charlie@example.com"`)
	})
}

// Test_DataType_JSON_Query tests filtering on a JSON field.
func Test_DataType_JSON_Query(t *testing.T) {
	table := chCCreateTypeTable("test_json_query", "data String")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.List{
			g.Map{"id": 1, "data": `{"name":"David","age":20}`},
			g.Map{"id": 2, "data": `{"name":"Eve","age":30}`},
			g.Map{"id": 3, "data": `{"name":"Frank","age":25}`},
		}).Insert()
		t.AssertNil(err)

		count, err := db.Model(table).Where("JSONExtractInt(data, 'age') > ?", 25).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}

// Test_DataType_JSON_Update tests replacing a JSON document.
func Test_DataType_JSON_Update(t *testing.T) {
	table := chCCreateTypeTable("test_json_update", "data String")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id":   1,
			"data": `{"name":"Grace","age":28}`,
		}).Insert()
		t.AssertNil(err)

		_, err = db.Model(table).Data(g.Map{
			"data": `{"name":"Grace","age":29}`,
		}).Where("id", 1).Update()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		expected := map[string]any{"name": "Grace", "age": float64(29)}
		var actual map[string]any
		err = json.Unmarshal([]byte(one["data"].String()), &actual)
		t.AssertNil(err)
		t.Assert(actual, expected)
	})
}

// Test_DataType_Binary_Small tests small binary data in a String column, which holds arbitrary bytes.
func Test_DataType_Binary_Small(t *testing.T) {
	table := chCCreateTypeTable("test_binary_small", "data String")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		binaryData := []byte{0x00, 0x01, 0x02, 0x03, 0xFF}
		_, err := db.Model(table).Data(g.Map{
			"id":   1,
			"data": binaryData,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(bytes.Equal(one["data"].Bytes(), binaryData), true)
	})
}

// Test_DataType_Binary_Large tests large binary data (1MB+)
func Test_DataType_Binary_Large(t *testing.T) {
	table := chCCreateTypeTable("test_binary_large", "data String")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		size := 1024 * 1024
		largeBinary := make([]byte, size)
		for i := 0; i < size; i++ {
			largeBinary[i] = byte(i % 256)
		}

		_, err := db.Model(table).Data(g.Map{
			"id":   1,
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
	table := chCCreateTypeTable("test_binary_integrity", "data String, checksum String")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		binaryData := []byte("Hello, World! This is a binary test data with special chars: \x00\xFF\xAB")
		hash := sha256.Sum256(binaryData)
		checksum := hex.EncodeToString(hash[:])

		_, err := db.Model(table).Data(g.Map{
			"id":       1,
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
	table := chCCreateTypeTable("test_binary_empty", "data Nullable(String)")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id":   1,
			"data": []byte{},
		}).Insert()
		t.AssertNil(err)

		_, err = db.Model(table).Data(g.Map{
			"id":   2,
			"data": nil,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(len(one["data"].Bytes()), 0)
		t.Assert(one["data"].IsNil(), false)

		one, err = db.Model(table).Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(one["data"].IsNil(), true)
	})
}

// Test_DataType_Decimal_HighPrecision tests high precision decimal (65,30)
func Test_DataType_Decimal_HighPrecision(t *testing.T) {
	table := chCCreateTypeTable("test_decimal_precision", "amount Decimal(65, 30)")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		value := "12345678901234567890123456789012345.123456789012345678901234567890"
		_, err := db.Model(table).Data(g.Map{
			"id":     1,
			"amount": value,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["amount"].String(), value)
	})
}

// Test_DataType_Decimal_Calculation tests decimal arithmetic
func Test_DataType_Decimal_Calculation(t *testing.T) {
	table := chCCreateTypeTable("test_decimal_calc", "price Decimal(10, 2), quantity Decimal(10, 2)")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id":       1,
			"price":    "19.99",
			"quantity": "3.5",
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Fields("price * quantity as total").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["total"].String(), "69.9650")
	})
}

// Test_DataType_Decimal_Boundary tests decimal boundary values
func Test_DataType_Decimal_Boundary(t *testing.T) {
	table := chCCreateTypeTable("test_decimal_boundary", "value Decimal(10, 2)")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id":    1,
			"value": "99999999.99",
		}).Insert()
		t.AssertNil(err)

		_, err = db.Model(table).Data(g.Map{
			"id":    2,
			"value": "-99999999.99",
		}).Insert()
		t.AssertNil(err)

		_, err = db.Model(table).Data(g.Map{
			"id":    3,
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

// Test_DataType_Decimal_Null tests NULL decimal values
func Test_DataType_Decimal_Null(t *testing.T) {
	table := chCCreateTypeTable("test_decimal_null", "value Nullable(Decimal(10, 2))")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id":    1,
			"value": nil,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["value"].IsNil(), true)
	})
}

// Test_DataType_Datetime_Timezone tests datetime with timezone handling
func Test_DataType_Datetime_Timezone(t *testing.T) {
	table := chCCreateTypeTable("test_datetime_tz", "created_at DateTime")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		dt := "2024-01-15 12:30:45"
		_, err := db.Model(table).Data(g.Map{
			"id":         1,
			"created_at": dt,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["created_at"].String(), dt)
	})
}

// Test_DataType_Datetime_Precision tests datetime with microsecond precision
func Test_DataType_Datetime_Precision(t *testing.T) {
	table := chCCreateTypeTable("test_datetime_precision", "created_at DateTime64(6)")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		dt := "2024-01-15 12:30:45.123456"
		_, err := db.Model(table).Data(g.Map{
			"id":         1,
			"created_at": dt,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		expected := "2024-01-15 12:30:45"
		actual := one["created_at"].String()[:19]
		t.Assert(actual, expected)
		t.Assert(one["created_at"].Time().Nanosecond(), 123456000)
	})
}

// Test_DataType_Datetime_Boundary tests datetime boundary values
// Note: DateTime64 ranges from 1900-01-01 to 2299-12-31, narrower than the DATETIME of MySQL, and
// clickhouse-go converts its values through int64 nanoseconds, which overflow after 2262-04-11,
// so the bound of 2299 is only checked as text on the server side.
func Test_DataType_Datetime_Boundary(t *testing.T) {
	table := chCCreateTypeTable("test_datetime_boundary", "dt DateTime64(0)")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id": 1,
			"dt": "1900-01-01 00:00:00",
		}).Insert()
		t.AssertNil(err)

		_, err = db.Model(table).Data(g.Map{
			"id": 2,
			"dt": "2262-04-10 23:59:59",
		}).Insert()
		t.AssertNil(err)

		all, err := db.Model(table).Order("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["dt"].String(), "1900-01-01 00:00:00")
		t.Assert(all[1]["dt"].String(), "2262-04-10 23:59:59")
	})
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Exec(ctx, fmt.Sprintf(
			"INSERT INTO %s (id, dt) VALUES (3, '2299-12-31 23:59:59')", table,
		))
		t.AssertNil(err)

		value, err := db.Model(table).Fields("toString(dt)").Where("id", 3).Value()
		t.AssertNil(err)
		t.Assert(value.String(), "2299-12-31 23:59:59")
	})
}

// Test_DataType_Datetime_Null tests NULL datetime
func Test_DataType_Datetime_Null(t *testing.T) {
	table := chCCreateTypeTable("test_datetime_null", "dt Nullable(DateTime)")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id": 1,
			"dt": nil,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["dt"].IsNil(), true)
	})
}

// Test_DataType_Datetime_Update tests datetime updates
func Test_DataType_Datetime_Update(t *testing.T) {
	table := chCCreateTypeTable("test_datetime_update", "dt DateTime")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		dt1 := "2024-01-01 10:00:00"
		_, err := db.Model(table).Data(g.Map{
			"id": 1,
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

// Test_DataType_Enum_Valid tests valid Enum8 values
func Test_DataType_Enum_Valid(t *testing.T) {
	table := chCCreateTypeTable("test_enum_valid", "status Enum8('pending' = 1, 'approved' = 2, 'rejected' = 3)")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.List{
			g.Map{"id": 1, "status": "pending"},
			g.Map{"id": 2, "status": "approved"},
			g.Map{"id": 3, "status": "rejected"},
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

// Test_DataType_Enum_Invalid tests that a value outside of an Enum8 is rejected
func Test_DataType_Enum_Invalid(t *testing.T) {
	table := chCCreateTypeTable("test_enum_invalid", "status Enum8('pending' = 1, 'approved' = 2, 'rejected' = 3)")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id":     1,
			"status": "invalid_status",
		}).Insert()
		t.AssertNE(err, nil)

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}

// Test_DataType_Set_Valid tests sets of values.
// Note: ClickHouse has no SET type; an Array of an Enum8 holds the members instead, read as a JSON array.
func Test_DataType_Set_Valid(t *testing.T) {
	table := chCCreateTypeTable("test_set_valid", "permissions Array(Enum8('read' = 1, 'write' = 2, 'execute' = 3))")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id":          1,
			"permissions": []string{"read"},
		}).Insert()
		t.AssertNil(err)

		_, err = db.Model(table).Data(g.Map{
			"id":          2,
			"permissions": []string{"read", "write"},
		}).Insert()
		t.AssertNil(err)

		_, err = db.Model(table).Data(g.Map{
			"id":          3,
			"permissions": []string{"read", "write", "execute"},
		}).Insert()
		t.AssertNil(err)

		all, err := db.Model(table).Order("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(all[0]["permissions"].Strings(), []string{"read"})
		t.Assert(all[1]["permissions"].Strings(), []string{"read", "write"})
		t.Assert(all[2]["permissions"].Strings(), []string{"read", "write", "execute"})
	})
}

// Test_DataType_Set_Empty tests an empty set of values.
// Note: see Test_DataType_Set_Valid for the Array of an Enum8 used in place of SET.
func Test_DataType_Set_Empty(t *testing.T) {
	table := chCCreateTypeTable("test_set_empty", "permissions Array(Enum8('read' = 1, 'write' = 2, 'execute' = 3))")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id":          1,
			"permissions": []string{},
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(len(one["permissions"].Strings()), 0)
	})
}

// Test_DataType_Geometry_Point tests the Point type with readWKTPoint and wkt, the counterparts
// of ST_GeomFromText and ST_AsText.
// Note: wkt writes the coordinates with 6 significant digits.
func Test_DataType_Geometry_Point(t *testing.T) {
	table := chCCreateTypeTable("test_geo_point", "location Point")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Exec(ctx, fmt.Sprintf(
			"INSERT INTO %s (id, location) VALUES (1, readWKTPoint('POINT(116.4074 39.9042)'))", table,
		))
		t.AssertNil(err)

		one, err := db.Model(table).Fields("wkt(location) as location_text").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["location_text"].String(), "POINT(116.407 39.9042)")

		one, err = db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["location"].String(), "[116.4074,39.9042]")
	})
}

// Test_DataType_Geometry_Polygon tests the Polygon type
func Test_DataType_Geometry_Polygon(t *testing.T) {
	table := chCCreateTypeTable("test_geo_polygon", "area Polygon")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		polygon := "POLYGON((0 0, 10 0, 10 10, 0 10, 0 0))"
		_, err := db.Exec(ctx, fmt.Sprintf("INSERT INTO %s (id, area) VALUES (1, readWKTPolygon('%s'))", table, polygon))
		t.AssertNil(err)

		one, err := db.Model(table).Fields("wkt(area) as area_text").Where("id", 1).One()
		t.AssertNil(err)
		expected := "POLYGON((0 0,10 0,10 10,0 10,0 0))"
		actual := strings.ReplaceAll(one["area_text"].String(), ", ", ",")
		t.Assert(actual, expected)
	})
}

// Test_DataType_Geometry_Null tests NULL geometry values.
// Note: the geometry types of ClickHouse are tuples and arrays, which cannot be Nullable, so a
// nullable Point column is rejected and NULL inserted into a Point column stores (0, 0).
func Test_DataType_Geometry_Null(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		name := fmt.Sprintf(`test_geo_null_%d`, gtime.TimestampNano())
		defer dropTable(name)
		_, err := db.Exec(ctx, fmt.Sprintf(
			"CREATE TABLE %s (id UInt64, location Nullable(Point)) ENGINE = MergeTree() ORDER BY id", name,
		))
		t.AssertNE(err, nil)
	})
	gtest.C(t, func(t *gtest.T) {
		table := chCCreateTypeTable("test_geo_null", "location Point")
		defer dropTable(table)

		_, err := db.Exec(ctx, fmt.Sprintf("INSERT INTO %s (id, location) VALUES (1, NULL)", table))
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["location"].String(), "[0,0]")
	})
}

// Test_ClickHouse_DataType_Integer tests the integer types from Int8 to UInt256 at their bounds.
func Test_ClickHouse_DataType_Integer(t *testing.T) {
	table := chCCreateTypeTable(
		"test_integer",
		"i8 Int8, u8 UInt8, i16 Int16, u16 UInt16, i32 Int32, u32 UInt32, i64 Int64, u64 UInt64, "+
			"i128 Int128, u256 UInt256",
	)
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		i128, _ := new(big.Int).SetString("-170141183460469231731687303715884105728", 10)
		u256, _ := new(big.Int).SetString("115792089237316195423570985008687907853269984665640564039457584007913129639935", 10)
		_, err := db.Model(table).Data(g.Map{
			"id":   1,
			"i8":   int8(-128),
			"u8":   uint8(255),
			"i16":  int16(-32768),
			"u16":  uint16(65535),
			"i32":  int32(-2147483648),
			"u32":  uint32(4294967295),
			"i64":  int64(-9223372036854775808),
			"u64":  uint64(18446744073709551615),
			"i128": i128,
			"u256": u256,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["i8"].Int(), -128)
		t.Assert(one["u8"].Int(), 255)
		t.Assert(one["i16"].Int(), -32768)
		t.Assert(one["u16"].Int(), 65535)
		t.Assert(one["i32"].Int(), -2147483648)
		t.Assert(one["u32"].Int64(), 4294967295)
		t.Assert(one["i64"].Int64(), int64(-9223372036854775808))
		t.AssertEQ(one["u64"].Uint64(), uint64(18446744073709551615))
		t.Assert(one["i128"].String(), i128.String())
		t.Assert(one["u256"].String(), u256.String())
	})
}

// Test_ClickHouse_DataType_Decimal tests the Decimal32, Decimal64, Decimal128 and Decimal256 types.
func Test_ClickHouse_DataType_Decimal(t *testing.T) {
	table := chCCreateTypeTable(
		"test_decimal",
		"d32 Decimal32(4), d64 Decimal64(6), d128 Decimal128(10), d256 Decimal256(20)",
	)
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id":   1,
			"d32":  decimal.RequireFromString("12345.6789"),
			"d64":  decimal.RequireFromString("-123456789012.345678"),
			"d128": decimal.RequireFromString("1234567890123456789012345678.0123456789"),
			"d256": decimal.RequireFromString("-12345678901234567890.12345678901234567890"),
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["d32"].String(), "12345.6789")
		t.Assert(one["d64"].String(), "-123456789012.345678")
		t.Assert(one["d128"].String(), "1234567890123456789012345678.0123456789")
		t.Assert(one["d256"].String(), "-12345678901234567890.12345678901234567890")
		t.Assert(one["d32"].Float64(), 12345.6789)
	})
}

// Test_ClickHouse_DataType_Float tests the Float32, Float64 and Bool types.
func Test_ClickHouse_DataType_Float(t *testing.T) {
	table := chCCreateTypeTable("test_float", "f32 Float32, f64 Float64, flag Bool")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.List{
			g.Map{"id": 1, "f32": float32(1.5), "f64": 3.141592653589793, "flag": true},
			g.Map{"id": 2, "f32": float32(-0.25), "f64": -2.718281828459045, "flag": false},
		}).Insert()
		t.AssertNil(err)

		all, err := db.Model(table).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["f32"].Float32(), float32(1.5))
		t.Assert(all[0]["f64"].Float64(), 3.141592653589793)
		t.Assert(all[0]["flag"].Bool(), true)
		t.Assert(all[1]["f32"].Float32(), float32(-0.25))
		t.Assert(all[1]["f64"].Float64(), -2.718281828459045)
		t.Assert(all[1]["flag"].Bool(), false)

		count, err := db.Model(table).Where("flag", true).Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}

// Test_ClickHouse_DataType_UUID tests the UUID type with uuid.UUID and string values.
func Test_ClickHouse_DataType_UUID(t *testing.T) {
	table := chCCreateTypeTable("test_uuid", "uid UUID")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var (
			uid1 = uuid.New()
			uid2 = uuid.New()
		)
		_, err := db.Model(table).Data(g.List{
			g.Map{"id": 1, "uid": uid1},
			g.Map{"id": 2, "uid": uid2.String()},
		}).Insert()
		t.AssertNil(err)

		all, err := db.Model(table).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["uid"].String(), uid1.String())
		t.Assert(all[1]["uid"].String(), uid2.String())

		one, err := db.Model(table).Where("uid", uid2).One()
		t.AssertNil(err)
		t.Assert(one["id"], 2)

		type Entity struct {
			Id  int
			Uid uuid.UUID
		}
		var entity *Entity
		err = db.Model(table).Where("id", 1).Scan(&entity)
		t.AssertNil(err)
		t.Assert(entity.Uid, uid1)
	})
}

// Test_ClickHouse_DataType_Array tests the Array type with integers, strings and nullable elements.
func Test_ClickHouse_DataType_Array(t *testing.T) {
	table := chCCreateTypeTable("test_array", "ints Array(Int64), strs Array(String), nulls Array(Nullable(Int32))")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var (
			one32  = int32(1)
			three  = int32(3)
			nulls  = []*int32{&one32, nil, &three}
			insert = g.List{
				g.Map{"id": 1, "ints": []int64{1, 2, 3}, "strs": []string{"a", "b"}, "nulls": nulls},
				g.Map{"id": 2, "ints": []int64{}, "strs": []string{}, "nulls": []*int32{}},
			}
		)
		_, err := db.Model(table).Data(insert).Insert()
		t.AssertNil(err)

		all, err := db.Model(table).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["ints"].Int64s(), []int64{1, 2, 3})
		t.Assert(all[0]["strs"].Strings(), []string{"a", "b"})
		t.Assert(all[0]["nulls"].String(), "[1,null,3]")
		t.Assert(len(all[1]["ints"].Int64s()), 0)
		t.Assert(len(all[1]["strs"].Strings()), 0)

		count, err := db.Model(table).Where("has(ints, ?)", 2).Count()
		t.AssertNil(err)
		t.Assert(count, 1)

		type Entity struct {
			Id   int
			Ints []int64
			Strs []string
		}
		var entities []Entity
		err = db.Model(table).OrderAsc("id").Scan(&entities)
		t.AssertNil(err)
		t.Assert(len(entities), 2)
		t.Assert(entities[0].Ints, []int64{1, 2, 3})
		t.Assert(entities[0].Strs, []string{"a", "b"})
		t.Assert(len(entities[1].Ints), 0)
	})
}

// Test_ClickHouse_DataType_Map tests the Map type.
func Test_ClickHouse_DataType_Map(t *testing.T) {
	table := chCCreateTypeTable("test_map", "attrs Map(String, UInt64)")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.List{
			g.Map{"id": 1, "attrs": map[string]uint64{"a": 1, "b": 2}},
			g.Map{"id": 2, "attrs": map[string]uint64{}},
		}).Insert()
		t.AssertNil(err)

		all, err := db.Model(table).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["attrs"].Map(), g.Map{"a": 1, "b": 2})
		t.Assert(len(all[1]["attrs"].Map()), 0)

		value, err := db.Model(table).Fields("attrs['b']").Where("id", 1).Value()
		t.AssertNil(err)
		t.Assert(value.Int(), 2)

		type Entity struct {
			Id    int
			Attrs map[string]uint64
		}
		var entity *Entity
		err = db.Model(table).Where("id", 1).Scan(&entity)
		t.AssertNil(err)
		t.Assert(entity.Attrs, map[string]uint64{"a": 1, "b": 2})
	})
}

// Test_ClickHouse_DataType_Enum16 tests the Enum16 type and filtering on its values.
func Test_ClickHouse_DataType_Enum16(t *testing.T) {
	table := chCCreateTypeTable("test_enum16", "level Enum16('low' = -1000, 'high' = 1000)")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.List{
			g.Map{"id": 1, "level": "low"},
			g.Map{"id": 2, "level": "high"},
		}).Insert()
		t.AssertNil(err)

		all, err := db.Model(table).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["level"].String(), "low")
		t.Assert(all[1]["level"].String(), "high")

		value, err := db.Model(table).Fields("CAST(level, 'Int16')").Where("level", "high").Value()
		t.AssertNil(err)
		t.Assert(value.Int(), 1000)
	})
}

// Test_ClickHouse_DataType_IP tests the IPv4 and IPv6 types with net.IP and string values.
func Test_ClickHouse_DataType_IP(t *testing.T) {
	table := chCCreateTypeTable("test_ip", "v4 IPv4, v6 IPv6")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.List{
			g.Map{"id": 1, "v4": net.ParseIP("192.168.1.1"), "v6": net.ParseIP("2001:db8::1")},
			g.Map{"id": 2, "v4": "10.0.0.1", "v6": "::1"},
		}).Insert()
		t.AssertNil(err)

		all, err := db.Model(table).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["v4"].String(), "192.168.1.1")
		t.Assert(all[0]["v6"].String(), "2001:db8::1")
		t.Assert(all[1]["v4"].String(), "10.0.0.1")
		t.Assert(all[1]["v6"].String(), "::1")

		count, err := db.Model(table).Where("v4", "10.0.0.1").Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}

// Test_ClickHouse_DataType_Date32 tests the Date and Date32 types, the latter at its bounds.
func Test_ClickHouse_DataType_Date32(t *testing.T) {
	table := chCCreateTypeTable("test_date32", "d Date, d32 Date32")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.List{
			g.Map{"id": 1, "d": "2024-01-15", "d32": "1900-01-01"},
			g.Map{"id": 2, "d": "2149-06-06", "d32": "2299-12-31"},
		}).Insert()
		t.AssertNil(err)

		all, err := db.Model(table).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["d"].String(), "2024-01-15")
		t.Assert(all[0]["d32"].GTime().Format("Y-m-d"), "1900-01-01")
		t.Assert(all[1]["d"].String(), "2149-06-06")
		t.Assert(all[1]["d32"].GTime().Format("Y-m-d"), "2299-12-31")

		count, err := db.Model(table).Where("d32 < ?", "2000-01-01").Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}

// Test_ClickHouse_DataType_DateTime64 tests the DateTime64 type with millisecond and nanosecond
// precisions and with an explicit time zone.
func Test_ClickHouse_DataType_DateTime64(t *testing.T) {
	table := chCCreateTypeTable("test_datetime64", "ms DateTime64(3), ns DateTime64(9), utc DateTime64(3, 'UTC')")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var (
			ms  = time.Date(2024, 1, 15, 12, 30, 45, 123000000, time.Local)
			ns  = time.Date(2024, 1, 15, 12, 30, 45, 123456789, time.Local)
			utc = time.Date(2024, 1, 15, 12, 30, 45, 123000000, time.UTC)
		)
		_, err := db.Model(table).Data(g.Map{
			"id":  1,
			"ms":  ms,
			"ns":  ns,
			"utc": utc,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["ms"].Time().Equal(ms), true)
		t.Assert(one["ns"].Time().Equal(ns), true)
		t.Assert(one["utc"].Time().Equal(utc), true)
		t.Assert(one["ms"].GTime().Format("Y-m-d H:i:s.u"), "2024-01-15 12:30:45.123")
		t.Assert(one["ns"].Time().Nanosecond(), 123456789)
	})
}

// Test_ClickHouse_DataType_Nullable tests nullable columns holding NULL and values.
func Test_ClickHouse_DataType_Nullable(t *testing.T) {
	table := chCCreateTypeTable("test_nullable", "n_int Nullable(Int32), n_str Nullable(String), n_float Nullable(Float64)")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.List{
			g.Map{"id": 1, "n_int": 10, "n_str": "a", "n_float": 1.5},
			g.Map{"id": 2, "n_int": nil, "n_str": nil, "n_float": nil},
		}).Insert()
		t.AssertNil(err)

		all, err := db.Model(table).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["n_int"].Int(), 10)
		t.Assert(all[0]["n_str"].String(), "a")
		t.Assert(all[0]["n_float"].Float64(), 1.5)
		t.Assert(all[1]["n_int"].IsNil(), true)
		t.Assert(all[1]["n_str"].IsNil(), true)
		t.Assert(all[1]["n_float"].IsNil(), true)

		count, err := db.Model(table).WhereNull("n_int").Count()
		t.AssertNil(err)
		t.Assert(count, 1)

		fields, err := db.TableFields(ctx, table)
		t.AssertNil(err)
		t.Assert(fields["n_int"].Null, true)
		t.Assert(fields["n_int"].Type, "Int32")
		t.Assert(fields["id"].Null, false)
	})
}

// Test_ClickHouse_DataType_LowCardinality tests the LowCardinality type over String and Nullable(String).
func Test_ClickHouse_DataType_LowCardinality(t *testing.T) {
	table := chCCreateTypeTable("test_low_cardinality", "lc LowCardinality(String), lcn LowCardinality(Nullable(String))")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.List{
			g.Map{"id": 1, "lc": "beijing", "lcn": "x"},
			g.Map{"id": 2, "lc": "shanghai", "lcn": nil},
			g.Map{"id": 3, "lc": "beijing", "lcn": "x"},
		}).Insert()
		t.AssertNil(err)

		all, err := db.Model(table).OrderAsc("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(all[0]["lc"].String(), "beijing")
		t.Assert(all[1]["lc"].String(), "shanghai")
		t.Assert(all[0]["lcn"].String(), "x")
		t.Assert(all[1]["lcn"].IsNil(), true)

		count, err := db.Model(table).Where("lc", "beijing").Count()
		t.AssertNil(err)
		t.Assert(count, 2)
	})
}
