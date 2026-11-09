// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mssql_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

func rawTypeCreateTable(prefix, columns string) string {
	table := fmt.Sprintf("%s_%d", prefix, gtime.TimestampNano())
	dropTable(table)
	if _, err := db.Exec(ctx, fmt.Sprintf(
		"CREATE TABLE [%s] (ID int IDENTITY(1,1) NOT NULL PRIMARY KEY, %s)", table, columns,
	)); err != nil {
		gtest.Fatal(err)
	}
	return table
}

func Test_Raw_Insert(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		user := db.Model(table)
		result, err := user.Data(g.Map{
			"id":          gdb.Raw("1+1"),
			"passport":    "port_1",
			"password":    "pass_1",
			"nickname":    "name_1",
			"create_time": gdb.Raw("GETDATE()"),
		}).Insert()
		t.AssertNil(err)
		n, _ := result.LastInsertId()
		t.Assert(n, 0)
		affected, _ := result.RowsAffected()
		t.Assert(affected, 1)

		one, err := db.Model(table).Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "port_1")
		t.AssertLT(gtime.Now().Sub(one["CREATE_TIME"].GTime()).Abs(), 24*time.Hour)
	})

	table2 := rawTypeCreateTable("t_rt_raw", "PASSPORT varchar(45) NULL, CREATE_TIME datetime NULL")
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table2).Data(g.Map{
			"passport":    "port_1",
			"create_time": gdb.Raw("GETDATE()"),
		}).Insert()
		t.AssertNil(err)
		n, _ := result.LastInsertId()
		t.Assert(n, 1)

		one, err := db.Model(table2).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "port_1")
		t.AssertLT(gtime.Now().Sub(one["CREATE_TIME"].GTime()).Abs(), 24*time.Hour)
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
					"id":          gdb.Raw("1+1"),
					"passport":    "port_2",
					"password":    "pass_2",
					"nickname":    "name_2",
					"create_time": gdb.Raw("GETDATE()"),
				},
				g.Map{
					"id":          gdb.Raw("2+2"),
					"passport":    "port_4",
					"password":    "pass_4",
					"nickname":    "name_4",
					"create_time": gdb.Raw("GETDATE()"),
				},
			},
		).Insert()
		t.AssertNil(err)
		n, _ := result.LastInsertId()
		t.Assert(n, 0)
		affected, _ := result.RowsAffected()
		t.Assert(affected, 2)

		ids, err := db.Model(table).Order("id").Array("id")
		t.AssertNil(err)
		t.Assert(gconv.Ints(ids), g.Slice{2, 4})
	})

	table2 := rawTypeCreateTable("t_rt_raw", "PASSPORT varchar(45) NULL, CREATE_TIME datetime NULL")
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table2).Data(
			g.List{
				g.Map{
					"passport":    gdb.Raw("'port_' + '2'"),
					"create_time": gdb.Raw("GETDATE()"),
				},
				g.Map{
					"passport":    gdb.Raw("'port_' + '4'"),
					"create_time": gdb.Raw("GETDATE()"),
				},
			},
		).Insert()
		t.AssertNil(err)
		n, _ := result.LastInsertId()
		t.Assert(n, 1)
		affected, _ := result.RowsAffected()
		t.Assert(affected, 2)

		all, err := db.Model(table2).Order("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["ID"], 1)
		t.Assert(all[0]["PASSPORT"], "port_2")
		t.Assert(all[1]["ID"], 2)
		t.Assert(all[1]["PASSPORT"], "port_4")
		t.AssertLT(gtime.Now().Sub(all[1]["CREATE_TIME"].GTime()).Abs(), 24*time.Hour)
	})
}

func Test_Raw_Update(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		user := db.Model(table)
		result, err := user.Data(g.Map{
			"id":          gdb.Raw("id+100"),
			"create_time": gdb.Raw("GETDATE()"),
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

		createTime, err := db.Model(table).Where("id", 101).Value("create_time")
		t.AssertNil(err)
		t.AssertLT(gtime.Now().Sub(createTime.GTime()).Abs(), 24*time.Hour)
	})
}

func Test_Raw_Where(t *testing.T) {
	table1 := createTable("Test_Raw_Where_Table1")
	table2 := createTable("Test_Raw_Where_Table2")
	defer dropTable(table1)
	defer dropTable(table2)

	// https://github.com/gogf/gf/issues/3922
	gtest.C(t, func(t *gtest.T) {
		expectSql := `SELECT TOP 1 * FROM "Test_Raw_Where_Table1" AS A WHERE NOT EXISTS (SELECT B.id FROM "Test_Raw_Where_Table2" AS B WHERE "B"."id"=A.id)`
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
		expectSql := `SELECT TOP 1 * FROM "Test_Raw_Where_Table1" AS A WHERE NOT EXISTS (SELECT B.id FROM "Test_Raw_Where_Table2" AS B WHERE B.id=A.id)`
		sql, err := gdb.ToSQL(ctx, func(ctx context.Context) error {
			s := db.Model(table2).As("B").Ctx(ctx).Fields("B.id").Where(gdb.Raw("B.id=A.id"))
			m := db.Model(table1).As("A").Ctx(ctx).Where("NOT EXISTS ?", s).Limit(1)
			_, err := m.All()
			return err
		})
		t.AssertNil(err)
		t.Assert(expectSql, sql)
	})
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Insert(ctx, table1, g.List{
			{"id": 1, "passport": "a", "nickname": "b"},
			{"id": 2, "passport": "b", "nickname": "a"},
		})
		t.AssertNil(err)
		_, err = db.Insert(ctx, table2, g.Map{"id": 1})
		t.AssertNil(err)

		s := db.Model(table2).As("B").Fields("B.id").Where("B.id", gdb.Raw("A.id"))
		all, err := db.Model(table1).As("A").Where("NOT EXISTS ?", s).Limit(1).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["ID"], 2)
	})
	// https://github.com/gogf/gf/issues/3915
	gtest.C(t, func(t *gtest.T) {
		expectSql := `SELECT * FROM "Test_Raw_Where_Table1" WHERE "passport" < [nickname]`
		sql, err := gdb.ToSQL(ctx, func(ctx context.Context) error {
			m := db.Model(table1).Ctx(ctx).WhereLT("passport", gdb.Raw("[nickname]"))
			_, err := m.All()
			return err
		})
		t.AssertNil(err)
		t.Assert(expectSql, sql)
	})
	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table1).WhereLT("passport", gdb.Raw("[nickname]")).All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["ID"], 1)
	})
}

// Test_DataType_JSON_Insert tests JSON data insertion
func Test_DataType_JSON_Insert(t *testing.T) {
	table := rawTypeCreateTable("t_rt_json", "DATA nvarchar(max) NULL, DATA_VC varchar(max) NULL")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		data := `{"name":"John","age":30}`
		result, err := db.Model(table).Data(g.Map{
			"data":    data,
			"data_vc": data,
		}).Insert()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		expected := map[string]any{"name": "John", "age": float64(30)}
		for _, column := range []string{"DATA", "DATA_VC"} {
			var actual map[string]any
			err = json.Unmarshal([]byte(one[column].String()), &actual)
			t.AssertNil(err)
			t.Assert(actual, expected)
		}

		valid, err := db.Model(table).Fields("ISJSON(DATA) AS V1, ISJSON(DATA_VC) AS V2").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(valid["V1"], 1)
		t.Assert(valid["V2"], 1)
	})
}

// Test_DataType_JSON_Extract tests JSON_EXTRACT function
func Test_DataType_JSON_Extract(t *testing.T) {
	table := rawTypeCreateTable("t_rt_json", "DATA nvarchar(max) NULL")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"data": `{"name":"Alice","age":25,"city":"Beijing"}`,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Fields("JSON_VALUE(DATA, '$.name') AS NAME").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["NAME"].String(), "Alice")

		one, err = db.Model(table).Fields("JSON_VALUE(DATA, '$.age') AS AGE").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["AGE"].Int(), 25)

		one, err = db.Model(table).Fields("JSON_QUERY(DATA, '$') AS DOC").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["DOC"].String(), `{"name":"Alice","age":25,"city":"Beijing"}`)
	})
}

// Test_DataType_JSON_Set tests JSON_SET function
func Test_DataType_JSON_Set(t *testing.T) {
	table := rawTypeCreateTable("t_rt_json", "DATA nvarchar(max) NULL")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"data": `{"name":"Bob"}`,
		}).Insert()
		t.AssertNil(err)

		_, err = db.Exec(ctx, fmt.Sprintf("UPDATE [%s] SET DATA = JSON_MODIFY(DATA, '$.age', 30) WHERE ID = 1", table))
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		expected := map[string]any{"name": "Bob", "age": float64(30)}
		var actual map[string]any
		err = json.Unmarshal([]byte(one["DATA"].String()), &actual)
		t.AssertNil(err)
		t.Assert(actual, expected)
	})
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"data": gdb.Raw("JSON_MODIFY(DATA, '$.city', 'Paris')"),
		}).Where("id", 1).Update()
		t.AssertNil(err)

		city, err := db.Model(table).Fields("JSON_VALUE(DATA, '$.city')").Where("id", 1).Value()
		t.AssertNil(err)
		t.Assert(city, "Paris")
	})
}

// Test_DataType_JSON_Array tests JSON array operations
func Test_DataType_JSON_Array(t *testing.T) {
	table := rawTypeCreateTable("t_rt_json", "DATA nvarchar(max) NULL")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"data": `["apple","banana","cherry"]`,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Fields("JSON_VALUE(DATA, '$[0]') AS FIRST").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["FIRST"].String(), "apple")

		values, err := db.GetArray(ctx, fmt.Sprintf(
			"SELECT j.[value] FROM [%s] CROSS APPLY OPENJSON(DATA) j WHERE ID = 1 ORDER BY CAST(j.[key] AS int)", table,
		))
		t.AssertNil(err)
		t.Assert(values, g.Slice{"apple", "banana", "cherry"})
	})
}

// Test_DataType_JSON_Null tests JSON NULL handling
func Test_DataType_JSON_Null(t *testing.T) {
	table := rawTypeCreateTable("t_rt_json", "DATA nvarchar(max) NULL")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"data": nil,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["DATA"].IsNil(), true)
	})
}

// Test_DataType_JSON_Complex tests complex nested JSON
func Test_DataType_JSON_Complex(t *testing.T) {
	table := rawTypeCreateTable("t_rt_json", "DATA nvarchar(max) NULL")
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

		one, err := db.Model(table).Fields("JSON_VALUE(DATA, '$.user.contacts.email') AS EMAIL").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["EMAIL"].String(), "charlie@example.com")

		one, err = db.Model(table).Fields("JSON_QUERY(DATA, '$.user.tags') AS TAGS").Where("id", 1).One()
		t.AssertNil(err)
		var tags []string
		t.AssertNil(json.Unmarshal(one["TAGS"].Bytes(), &tags))
		t.Assert(tags, g.SliceStr{"developer", "gopher"})
	})
}

// Test_DataType_JSON_Query tests JSON query with WHERE clause
func Test_DataType_JSON_Query(t *testing.T) {
	table := rawTypeCreateTable("t_rt_json", "DATA nvarchar(max) NULL")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.List{
			g.Map{"data": `{"name":"David","age":20}`},
			g.Map{"data": `{"name":"Eve","age":30}`},
			g.Map{"data": `{"name":"Frank","age":25}`},
		}).Insert()
		t.AssertNil(err)

		count, err := db.Model(table).Where("CAST(JSON_VALUE(DATA, '$.age') AS int) > ?", 25).Count()
		t.AssertNil(err)
		t.Assert(count, 1)

		name, err := db.Model(table).Fields("JSON_VALUE(DATA, '$.name')").Where("JSON_VALUE(DATA, '$.name') = ?", "Frank").Value()
		t.AssertNil(err)
		t.Assert(name, "Frank")
	})
}

// Test_DataType_JSON_Update tests updating JSON data
func Test_DataType_JSON_Update(t *testing.T) {
	table := rawTypeCreateTable("t_rt_json", "DATA nvarchar(max) NULL")
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
		expected := map[string]any{"name": "Grace", "age": float64(29)}
		var actual map[string]any
		err = json.Unmarshal([]byte(one["DATA"].String()), &actual)
		t.AssertNil(err)
		t.Assert(actual, expected)
	})
}

// Test_DataType_Binary_Small tests small binary data
func Test_DataType_Binary_Small(t *testing.T) {
	table := rawTypeCreateTable("t_rt_bin", "DATA varbinary(max) NULL, DATA_VB varbinary(16) NULL, DATA_FIX binary(8) NULL")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		binaryData := []byte{0x00, 0x01, 0x02, 0x03, 0xFF}
		_, err := db.Model(table).Data(g.Map{
			"data":     binaryData,
			"data_vb":  binaryData,
			"data_fix": binaryData,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(bytes.Equal(one["DATA"].Bytes(), binaryData), true)
		t.Assert(bytes.Equal(one["DATA_VB"].Bytes(), binaryData), true)
		t.Assert(bytes.Equal(one["DATA_FIX"].Bytes(), []byte{0x00, 0x01, 0x02, 0x03, 0xFF, 0x00, 0x00, 0x00}), true)
	})
}

// Test_DataType_Binary_Large tests large binary data (1MB+)
func Test_DataType_Binary_Large(t *testing.T) {
	table := rawTypeCreateTable("t_rt_bin", "DATA varbinary(max) NULL")
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

		stored, err := db.Model(table).Fields("DATALENGTH(DATA)").Where("id", 1).Value()
		t.AssertNil(err)
		t.Assert(stored.Int(), size)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(len(one["DATA"].Bytes()), size)
		t.Assert(bytes.Equal(one["DATA"].Bytes(), largeBinary), true)
	})
}

// Test_DataType_Binary_Integrity tests binary data integrity with checksum
func Test_DataType_Binary_Integrity(t *testing.T) {
	table := rawTypeCreateTable("t_rt_bin", "DATA varbinary(max) NULL, CHECKSUM varchar(64) NULL")
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

		retrievedHash := sha256.Sum256(one["DATA"].Bytes())
		retrievedChecksum := hex.EncodeToString(retrievedHash[:])
		t.Assert(retrievedChecksum, checksum)
		t.Assert(one["CHECKSUM"], checksum)

		serverHash, err := db.Model(table).Fields("LOWER(CONVERT(varchar(64), HASHBYTES('SHA2_256', DATA), 2))").Where("id", 1).Value()
		t.AssertNil(err)
		t.Assert(serverHash.String(), checksum)
	})
}

// Test_DataType_Binary_Empty tests empty and NULL binary
func Test_DataType_Binary_Empty(t *testing.T) {
	table := rawTypeCreateTable("t_rt_bin", "DATA varbinary(max) NULL, DATA_VB varbinary(16) NULL")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"data":    []byte{},
			"data_vb": []byte{},
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(len(one["DATA"].Bytes()), 0)
		t.Assert(len(one["DATA_VB"].Bytes()), 0)

		lengths, err := db.Model(table).Fields("DATALENGTH(DATA) AS L1, DATALENGTH(DATA_VB) AS L2").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(lengths["L1"], 0)
		t.Assert(lengths["L2"], 0)
	})

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"data":    nil,
			"data_vb": nil,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(one["DATA"].IsNil(), true)
		t.Assert(one["DATA_VB"].IsNil(), true)
	})
}

// Test_DataType_Decimal_HighPrecision tests high precision decimal (65,30)
func Test_DataType_Decimal_HighPrecision(t *testing.T) {
	table := rawTypeCreateTable("t_rt_dec", "AMOUNT decimal(38,30) NULL, AMOUNT_INT decimal(38,0) NULL, AMOUNT_BIG bigint NULL")
	defer dropTable(table)

	var (
		value    = "12345678.123456789012345678901234567890"
		intValue = "12345678901234567890123456789012345678"
		bigValue = "9223372036854775807"
	)
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"amount":     value,
			"amount_int": intValue,
			"amount_big": bigValue,
		}).Insert()
		t.AssertNil(err)

		stored, err := db.Model(table).
			Fields("CAST(AMOUNT AS varchar(50)) AS S1, CAST(AMOUNT_INT AS varchar(50)) AS S2, CAST(AMOUNT_BIG AS varchar(50)) AS S3").
			Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(stored["S1"].String(), value)
		t.Assert(stored["S2"].String(), intValue)
		t.Assert(stored["S3"].String(), bigValue)
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["AMOUNT"].String(), value)
		t.Assert(one["AMOUNT_INT"].String(), intValue)
		t.Assert(one["AMOUNT_BIG"].String(), bigValue)
		t.Assert(one["AMOUNT_BIG"].Int64(), int64(math.MaxInt64))
	})
	gtest.C(t, func(t *gtest.T) {
		id, err := db.Model(table).Data(g.Map{
			"amount_big": int64(math.MinInt64),
		}).InsertAndGetId()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", id).One()
		t.AssertNil(err)
		t.Assert(one["AMOUNT_BIG"].Int64(), int64(math.MinInt64))
	})
}

// Test_DataType_Decimal_Calculation tests decimal arithmetic
func Test_DataType_Decimal_Calculation(t *testing.T) {
	table := rawTypeCreateTable("t_rt_dec", `PRICE decimal(10,2) NULL, QUANTITY decimal(10,2) NULL,
		PRICE_F float NULL, QUANTITY_F float NULL, PRICE_R real NULL, QUANTITY_R real NULL`)
	defer dropTable(table)

	var (
		priceF, quantityF = 19.99, 3.5
		priceR, quantityR = float32(19.99), float32(3.5)
	)
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"price":      "19.99",
			"quantity":   "3.5",
			"price_f":    priceF,
			"quantity_f": quantityF,
			"price_r":    priceR,
			"quantity_r": quantityR,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Fields("price * quantity as total").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["total"].String(), "69.9650")
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Fields("price_f * quantity_f as total").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["total"].Float64(), priceF*quantityF)
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Fields("price_r * quantity_r as total").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["total"].Float32(), priceR*quantityR)
	})
}

// Test_DataType_Decimal_Boundary tests decimal boundary values
func Test_DataType_Decimal_Boundary(t *testing.T) {
	table := rawTypeCreateTable("t_rt_dec", "VALUE decimal(10,2) NULL, VALUE_F float NULL, VALUE_R real NULL")
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
		t.Assert(all[0]["VALUE"].String(), "99999999.99")
		t.Assert(all[1]["VALUE"].String(), "-99999999.99")
		t.Assert(all[2]["VALUE"].String(), "0.00")
	})
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"value": "999999999.99",
		}).Insert()
		t.AssertNE(err, nil)
		t.AssertIN("Arithmetic overflow error", err.Error())
	})
	gtest.C(t, func(t *gtest.T) {
		for _, value := range []float64{math.MaxFloat64, -math.MaxFloat64} {
			id, err := db.Model(table).Data(g.Map{
				"value_f": value,
			}).InsertAndGetId()
			t.AssertNil(err)

			one, err := db.Model(table).Where("id", id).One()
			t.AssertNil(err)
			t.Assert(one["VALUE_F"].Float64(), value)
		}
	})
	gtest.C(t, func(t *gtest.T) {
		for _, value := range []float32{math.MaxFloat32, -math.MaxFloat32} {
			id, err := db.Model(table).Data(g.Map{
				"value_r": value,
			}).InsertAndGetId()
			t.AssertNil(err)

			one, err := db.Model(table).Where("id", id).One()
			t.AssertNil(err)
			t.Assert(one["VALUE_R"].Float32(), value)
		}
	})
}

// Test_DataType_Decimal_Null tests NULL decimal values
func Test_DataType_Decimal_Null(t *testing.T) {
	table := rawTypeCreateTable("t_rt_dec", "VALUE decimal(10,2) NULL, VALUE_F float NULL, VALUE_R real NULL")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"value":   nil,
			"value_f": nil,
			"value_r": nil,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["VALUE"].IsNil(), true)
		t.Assert(one["VALUE_F"].IsNil(), true)
		t.Assert(one["VALUE_R"].IsNil(), true)
	})
}

// Test_DataType_Datetime_Timezone tests datetime with timezone handling
func Test_DataType_Datetime_Timezone(t *testing.T) {
	table := rawTypeCreateTable("t_rt_dt", "CREATED_AT datetime NULL, CREATED_AT_OFFSET datetimeoffset NULL")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		dt := "2024-01-15 12:30:45"
		_, err := db.Model(table).Data(g.Map{
			"created_at": dt,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["CREATED_AT"].String(), dt)
	})
	gtest.C(t, func(t *gtest.T) {
		dt := time.Date(2024, 1, 15, 12, 30, 45, 0, time.FixedZone("", 3*3600))
		id, err := db.Model(table).Data(g.Map{
			"created_at_offset": dt,
		}).InsertAndGetId()
		t.AssertNil(err)

		stored, err := db.Model(table).Fields("CONVERT(varchar(40), CREATED_AT_OFFSET, 121)").Where("id", id).Value()
		t.AssertNil(err)
		t.Assert(stored.String(), "2024-01-15 12:30:45.0000000 +03:00")

		one, err := db.Model(table).Where("id", id).One()
		t.AssertNil(err)
		t.Assert(one["CREATED_AT_OFFSET"].Time().Equal(dt), true)
		_, offset := one["CREATED_AT_OFFSET"].Time().Zone()
		t.Assert(offset, 3*3600)
	})
	gtest.C(t, func(t *gtest.T) {
		id, err := db.Model(table).Data(g.Map{
			"created_at_offset": "2024-01-15 12:30:45 -05:00",
		}).InsertAndGetId()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", id).One()
		t.AssertNil(err)
		expected := time.Date(2024, 1, 15, 12, 30, 45, 0, time.FixedZone("", -5*3600))
		t.Assert(one["CREATED_AT_OFFSET"].Time().Equal(expected), true)
		_, offset := one["CREATED_AT_OFFSET"].Time().Zone()
		t.Assert(offset, -5*3600)
	})
}

// Test_DataType_Datetime_Precision tests datetime with microsecond precision
func Test_DataType_Datetime_Precision(t *testing.T) {
	table := rawTypeCreateTable("t_rt_dt", "CREATED_AT datetime2(7) NULL, CREATED_AT_LEGACY datetime NULL")
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
		actual := one["CREATED_AT"].String()[:19]
		t.Assert(actual, expected)
		t.Assert(one["CREATED_AT"].Time().Nanosecond(), 123456000)
	})
	gtest.C(t, func(t *gtest.T) {
		dt := time.Date(2024, 1, 15, 12, 30, 45, 123456700, time.UTC)
		id, err := db.Model(table).Data(g.Map{
			"created_at": dt,
		}).InsertAndGetId()
		t.AssertNil(err)

		stored, err := db.Model(table).Fields("CONVERT(varchar(30), CREATED_AT, 121)").Where("id", id).Value()
		t.AssertNil(err)
		t.Assert(stored.String(), dt.In(time.Local).Format("2006-01-02 15:04:05.0000000"))

		one, err := db.Model(table).Where("id", id).One()
		t.AssertNil(err)
		t.Assert(one["CREATED_AT"].Time().Nanosecond(), 123456700)
		t.Assert(one["CREATED_AT"].Time().Equal(dt), true)
	})
	gtest.C(t, func(t *gtest.T) {
		id, err := db.Model(table).Data(g.Map{
			"created_at_legacy": "2024-01-15 12:30:45.123",
		}).InsertAndGetId()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", id).One()
		t.AssertNil(err)
		t.Assert(one["CREATED_AT_LEGACY"].String(), "2024-01-15 12:30:45")
		t.Assert(one["CREATED_AT_LEGACY"].Time().Nanosecond(), 123000000)
	})
}

// Test_DataType_Datetime_Boundary tests datetime boundary values
func Test_DataType_Datetime_Boundary(t *testing.T) {
	table := rawTypeCreateTable("t_rt_dt", "DT datetime2(0) NULL, DT_LEGACY datetime NULL, DT_SMALL smalldatetime NULL")
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
		t.Assert(all[0]["DT"].String(), "1000-01-01 00:00:00")
		t.Assert(all[1]["DT"].String(), "9999-12-31 23:59:59")
	})
	gtest.C(t, func(t *gtest.T) {
		for _, dt := range []string{"1753-01-01 00:00:00", "9999-12-31 23:59:59"} {
			id, err := db.Model(table).Data(g.Map{
				"dt_legacy": dt,
			}).InsertAndGetId()
			t.AssertNil(err)

			one, err := db.Model(table).Where("id", id).One()
			t.AssertNil(err)
			t.Assert(one["DT_LEGACY"].String(), dt)
		}

		_, err := db.Model(table).Data(g.Map{
			"dt_legacy": "1000-01-01 00:00:00",
		}).Insert()
		t.AssertNE(err, nil)
		t.AssertIN("out-of-range", err.Error())
	})
	gtest.C(t, func(t *gtest.T) {
		for _, dt := range []string{"1900-01-01 00:00:00", "2079-06-06 23:59:00"} {
			id, err := db.Model(table).Data(g.Map{
				"dt_small": dt,
			}).InsertAndGetId()
			t.AssertNil(err)

			one, err := db.Model(table).Where("id", id).One()
			t.AssertNil(err)
			t.Assert(one["DT_SMALL"].String(), dt)
		}

		_, err := db.Model(table).Data(g.Map{
			"dt_small": "2079-06-07 00:00:00",
		}).Insert()
		t.AssertNE(err, nil)
		t.AssertIN("out-of-range", err.Error())
	})
}

// Test_DataType_Datetime_Null tests NULL datetime
func Test_DataType_Datetime_Null(t *testing.T) {
	table := rawTypeCreateTable("t_rt_dt", "DT datetime NULL, DT2 datetime2 NULL, DTO datetimeoffset NULL, D date NULL, T time NULL")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"dt":  nil,
			"dt2": nil,
			"dto": nil,
			"d":   nil,
			"t":   nil,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["DT"].IsNil(), true)
		t.Assert(one["DT2"].IsNil(), true)
		t.Assert(one["DTO"].IsNil(), true)
		t.Assert(one["D"].IsNil(), true)
		t.Assert(one["T"].IsNil(), true)
	})
}

// Test_DataType_Datetime_Update tests datetime updates
func Test_DataType_Datetime_Update(t *testing.T) {
	table := rawTypeCreateTable("t_rt_dt", "DT datetime NULL, D date NULL, T time(0) NULL")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		dt1 := "2024-01-01 10:00:00"
		_, err := db.Model(table).Data(g.Map{
			"dt": dt1,
			"d":  "2024-01-01",
			"t":  "10:00:00",
		}).Insert()
		t.AssertNil(err)

		dt2 := "2024-12-31 23:59:59"
		_, err = db.Model(table).Data(g.Map{
			"dt": dt2,
			"d":  "2024-12-31",
			"t":  "23:59:59",
		}).Where("id", 1).Update()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["DT"].String(), dt2)
		t.Assert(one["D"].String(), "2024-12-31")
		t.Assert(one["T"].String(), "23:59:59")
	})
	gtest.C(t, func(t *gtest.T) {
		id, err := db.Model(table).Data(g.Map{
			"dt": "2024-01-01 10:00:00",
		}).InsertAndGetId()
		t.AssertNil(err)

		_, err = db.Model(table).Data(g.Map{
			"dt": gdb.Raw("DATEADD(SECOND, 50399, DATEADD(DAY, 365, DT))"),
		}).Where("id", id).Update()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", id).One()
		t.AssertNil(err)
		t.Assert(one["DT"].String(), "2024-12-31 23:59:59")
	})
}

// Test_DataType_Enum_Valid tests valid ENUM values
func Test_DataType_Enum_Valid(t *testing.T) {
	table := rawTypeCreateTable("t_rt_enum", `STATUS varchar(20) NULL CHECK (STATUS IN ('pending','approved','rejected')),
		STATUS_CHAR char(8) NULL CHECK (STATUS_CHAR IN ('pending','approved','rejected')),
		STATUS_NCHAR nchar(8) NULL CHECK (STATUS_NCHAR IN ('pending','approved','rejected'))`)
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.List{
			g.Map{"status": "pending", "status_char": "pending", "status_nchar": "pending"},
			g.Map{"status": "approved", "status_char": "approved", "status_nchar": "approved"},
			g.Map{"status": "rejected", "status_char": "rejected", "status_nchar": "rejected"},
		}).Insert()
		t.AssertNil(err)

		all, err := db.Model(table).Order("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(all[0]["STATUS"].String(), "pending")
		t.Assert(all[1]["STATUS"].String(), "approved")
		t.Assert(all[2]["STATUS"].String(), "rejected")
		t.Assert(all[0]["STATUS_CHAR"].String(), "pending ")
		t.Assert(all[1]["STATUS_CHAR"].String(), "approved")
		t.Assert(all[2]["STATUS_CHAR"].String(), "rejected")
		t.Assert(all[0]["STATUS_NCHAR"].String(), "pending ")
		t.Assert(all[1]["STATUS_NCHAR"].String(), "approved")
		t.Assert(all[2]["STATUS_NCHAR"].String(), "rejected")
	})
}

// Test_DataType_Enum_Invalid tests invalid ENUM values (should fail or truncate)
func Test_DataType_Enum_Invalid(t *testing.T) {
	table := rawTypeCreateTable("t_rt_enum", "STATUS varchar(20) NULL CHECK (STATUS IN ('pending','approved','rejected'))")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"status": "invalid_status",
		}).Insert()
		t.AssertNE(err, nil)
		t.AssertIN("conflicted with the CHECK constraint", err.Error())

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}

const rawTypeSetCheck = `IN ('', 'read', 'write', 'execute', 'read,write', 'read,execute', 'write,execute', 'read,write,execute')`

// Test_DataType_Set_Valid tests valid SET values
func Test_DataType_Set_Valid(t *testing.T) {
	table := rawTypeCreateTable("t_rt_set", fmt.Sprintf(
		"PERMISSIONS varchar(30) NULL CHECK (PERMISSIONS %s), PERMISSIONS_N nvarchar(30) NULL CHECK (PERMISSIONS_N %s)",
		rawTypeSetCheck, rawTypeSetCheck,
	))
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"permissions":   "read",
			"permissions_n": "read",
		}).Insert()
		t.AssertNil(err)

		_, err = db.Model(table).Data(g.Map{
			"permissions":   "read,write",
			"permissions_n": "read,write",
		}).Insert()
		t.AssertNil(err)

		_, err = db.Model(table).Data(g.Map{
			"permissions":   "read,write,execute",
			"permissions_n": "read,write,execute",
		}).Insert()
		t.AssertNil(err)

		all, err := db.Model(table).Order("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(all[0]["PERMISSIONS"].String(), "read")
		t.Assert(all[1]["PERMISSIONS"].String(), "read,write")
		t.Assert(all[2]["PERMISSIONS"].String(), "read,write,execute")
		t.Assert(all[0]["PERMISSIONS_N"].String(), "read")
		t.Assert(all[1]["PERMISSIONS_N"].String(), "read,write")
		t.Assert(all[2]["PERMISSIONS_N"].String(), "read,write,execute")
	})
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"permissions": "read,delete",
		}).Insert()
		t.AssertNE(err, nil)
		t.AssertIN("conflicted with the CHECK constraint", err.Error())
	})
}

// Test_DataType_Set_Empty tests empty SET values
func Test_DataType_Set_Empty(t *testing.T) {
	table := rawTypeCreateTable("t_rt_set", fmt.Sprintf("PERMISSIONS varchar(30) NULL CHECK (PERMISSIONS %s)", rawTypeSetCheck))
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"permissions": "",
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["PERMISSIONS"].String(), "")
		t.Assert(one["PERMISSIONS"].IsNil(), false)
	})
}

// Test_DataType_Geometry_Point tests POINT geometry type
func Test_DataType_Geometry_Point(t *testing.T) {
	table := rawTypeCreateTable("t_rt_geo", "LOCATION geometry NULL, LOCATION_GEO geography NULL")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Exec(ctx, fmt.Sprintf(
			"INSERT INTO [%s] (LOCATION, LOCATION_GEO) VALUES (geometry::STGeomFromText('POINT(116.4074 39.9042)', 0), geography::STGeomFromText('POINT(116.4074 39.9042)', 4326))",
			table,
		))
		t.AssertNil(err)

		one, err := db.Model(table).Fields("LOCATION.STAsText() AS LOCATION_TEXT, LOCATION.STX AS X, LOCATION.STY AS Y").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["LOCATION_TEXT"].String(), "POINT (116.4074 39.9042)")
		t.Assert(one["X"].Float64(), 116.4074)
		t.Assert(one["Y"].Float64(), 39.9042)

		one, err = db.Model(table).Fields("LOCATION_GEO.STAsText() AS LOCATION_TEXT, LOCATION_GEO.Long AS LNG, LOCATION_GEO.Lat AS LAT, LOCATION_GEO.STSrid AS SRID").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["LOCATION_TEXT"].String(), "POINT (116.4074 39.9042)")
		t.Assert(one["LNG"].Float64(), 116.4074)
		t.Assert(one["LAT"].Float64(), 39.9042)
		t.Assert(one["SRID"].Int(), 4326)
	})
}

// Test_DataType_Geometry_Polygon tests POLYGON geometry type
func Test_DataType_Geometry_Polygon(t *testing.T) {
	table := rawTypeCreateTable("t_rt_geo", "AREA geometry NULL, AREA_GEO geography NULL")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		polygon := "POLYGON((0 0, 10 0, 10 10, 0 10, 0 0))"
		_, err := db.Exec(ctx, fmt.Sprintf(
			"INSERT INTO [%s] (AREA, AREA_GEO) VALUES (geometry::STGeomFromText('%s', 0), geography::STGeomFromText('%s', 4326))",
			table, polygon, polygon,
		))
		t.AssertNil(err)

		one, err := db.Model(table).Fields("AREA.STAsText() AS AREA_TEXT, AREA.STArea() AS AREA_SIZE, AREA_GEO.STAsText() AS AREA_GEO_TEXT").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["AREA_TEXT"].String(), "POLYGON ((0 0, 10 0, 10 10, 0 10, 0 0))")
		t.Assert(one["AREA_SIZE"].Float64(), 100)
		t.Assert(one["AREA_GEO_TEXT"].String(), "POLYGON ((0 0, 10 0, 10 10, 0 10, 0 0))")
	})
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"area": gdb.Raw("geometry::STGeomFromText('POLYGON((0 0, 5 0, 5 5, 0 5, 0 0))', 0)"),
		}).Where("id", 1).Update()
		t.AssertNil(err)

		size, err := db.Model(table).Fields("AREA.STArea()").Where("id", 1).Value()
		t.AssertNil(err)
		t.Assert(size.Float64(), 25)
	})
}

// Test_DataType_Geometry_Null tests NULL geometry values
func Test_DataType_Geometry_Null(t *testing.T) {
	table := rawTypeCreateTable("t_rt_geo", "LOCATION geometry NULL, LOCATION_GEO geography NULL")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"location":     nil,
			"location_geo": nil,
		}).Insert()
		t.AssertNil(err)

		count, err := db.Model(table).Where("id", 1).WhereNull("location").WhereNull("location_geo").Count()
		t.AssertNil(err)
		t.Assert(count, 1)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["LOCATION"].IsNil(), true)
		t.Assert(one["LOCATION_GEO"].IsNil(), true)
	})
}
