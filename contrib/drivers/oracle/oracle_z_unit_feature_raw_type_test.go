// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package oracle_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gregex"
)

func rawTypeCreateTable(prefix, columns string) string {
	table := fmt.Sprintf("%s_%d", prefix, gtime.TimestampMicro()%1e9)
	dropTable(table)
	if _, err := db.Exec(ctx, fmt.Sprintf("CREATE TABLE %s (ID NUMBER(10) NOT NULL PRIMARY KEY, %s)", table, columns)); err != nil {
		gtest.Fatal(err)
	}
	createAutoIncrement(table, "ID", 1)
	return table
}

func Test_Raw_Insert(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		user := db.Model(table)
		result, err := user.Data(g.Map{
			"id":          gdb.Raw(table + "_ID_SEQ.NEXTVAL+1"),
			"passport":    "port_1",
			"password":    "pass_1",
			"nickname":    "name_1",
			"create_time": gdb.Raw("TO_CHAR(SYSDATE, 'YYYY-MM-DD HH24:MI:SS')"),
		}).Insert()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		one, err := db.Model(table).Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "port_1")
		t.Assert(gregex.IsMatchString(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`, one["CREATE_TIME"].String()), true)
	})

	table2 := createTable()
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		result, err := db.Model(table2).Data(g.Map{
			"passport":    "port_1",
			"password":    "pass_1",
			"nickname":    "name_1",
			"create_time": gdb.Raw("TO_CHAR(SYSDATE, 'YYYY-MM-DD HH24:MI:SS')"),
		}).Insert()
		t.AssertNil(err)
		n, _ := result.LastInsertId()
		t.Assert(n, 1)

		one, err := db.Model(table2).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(gregex.IsMatchString(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`, one["CREATE_TIME"].String()), true)
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
					"id":          gdb.Raw(table + "_ID_SEQ.NEXTVAL+1"),
					"passport":    "port_2",
					"password":    "pass_2",
					"nickname":    "name_2",
					"create_time": gdb.Raw("TO_CHAR(SYSDATE, 'YYYY-MM-DD HH24:MI:SS')"),
				},
				g.Map{
					"id":          gdb.Raw(table + "_ID_SEQ.NEXTVAL+2"),
					"passport":    "port_4",
					"password":    "pass_4",
					"nickname":    "name_4",
					"create_time": gdb.Raw("TO_CHAR(SYSDATE, 'YYYY-MM-DD HH24:MI:SS')"),
				},
			},
		).Insert()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 2)

		ids, err := db.Model(table).Order("id").Array("id")
		t.AssertNil(err)
		t.Assert(ids, g.Slice{2, 4})
	})
}

func Test_Raw_Update(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		user := db.Model(table)
		result, err := user.Data(g.Map{
			"id":          gdb.Raw("id+100"),
			"create_time": gdb.Raw("TO_CHAR(SYSDATE, 'YYYY-MM-DD HH24:MI:SS')"),
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
		t.Assert(gregex.IsMatchString(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`, createTime.String()), true)
	})
}

func Test_Raw_Where(t *testing.T) {
	table1 := createTable("Test_Raw_Where_Table1")
	table2 := createTable("Test_Raw_Where_Table2")
	defer dropTable(table1)
	defer dropTable(table2)

	// https://github.com/gogf/gf/issues/3922
	gtest.C(t, func(t *gtest.T) {
		expectSql := "SELECT * FROM (SELECT * FROM Test_Raw_Where_Table1 A WHERE NOT EXISTS (SELECT B.id FROM Test_Raw_Where_Table2 B WHERE B.id=A.id)) WHERE ROWNUM <= 1"
		sql, err := gdb.ToSQL(ctx, func(ctx context.Context) error {
			s := db.Model(table2).As("B").Ctx(ctx).Fields("B.id").Where("B.id", gdb.Raw("A.id"))
			m := db.Model(table1).As("A").Ctx(ctx).Where("NOT EXISTS ?", s).Limit(1)
			_, err := m.All()
			return err
		})
		t.AssertNil(err)
		t.Assert(sql, expectSql)
	})
	gtest.C(t, func(t *gtest.T) {
		expectSql := "SELECT * FROM (SELECT * FROM Test_Raw_Where_Table1 A WHERE NOT EXISTS (SELECT B.id FROM Test_Raw_Where_Table2 B WHERE B.id=A.id)) WHERE ROWNUM <= 1"
		sql, err := gdb.ToSQL(ctx, func(ctx context.Context) error {
			s := db.Model(table2).As("B").Ctx(ctx).Fields("B.id").Where(gdb.Raw("B.id=A.id"))
			m := db.Model(table1).As("A").Ctx(ctx).Where("NOT EXISTS ?", s).Limit(1)
			_, err := m.All()
			return err
		})
		t.AssertNil(err)
		t.Assert(sql, expectSql)
	})
	gtest.C(t, func(t *gtest.T) {
		s := db.Model(table2).As("B").Fields("B.id").Where("B.id", gdb.Raw("A.id"))
		_, err := db.Model(table1).As("A").Where("NOT EXISTS ?", s).Limit(1).All()
		t.AssertNil(err)
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
	table := rawTypeCreateTable("t_rt_json", "DATA CLOB, DATA_NCLOB NCLOB, DATA_VC VARCHAR2(4000), DATA_LONG LONG")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		data := `{"name":"John","age":30}`
		result, err := db.Model(table).Data(g.Map{
			"data":       data,
			"data_nclob": data,
			"data_vc":    data,
			"data_long":  data,
		}).Insert()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		expected := map[string]interface{}{"name": "John", "age": float64(30)}
		for _, column := range []string{"DATA", "DATA_NCLOB", "DATA_VC", "DATA_LONG"} {
			var actual map[string]interface{}
			err = json.Unmarshal([]byte(one[column].String()), &actual)
			t.AssertNil(err)
			t.Assert(actual, expected)
		}
	})
}

// Test_DataType_JSON_Null tests JSON NULL handling
func Test_DataType_JSON_Null(t *testing.T) {
	table := rawTypeCreateTable("t_rt_json", "DATA CLOB")
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

// Test_DataType_JSON_Update tests updating JSON data
func Test_DataType_JSON_Update(t *testing.T) {
	table := rawTypeCreateTable("t_rt_json", "DATA CLOB")
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
		err = json.Unmarshal([]byte(one["DATA"].String()), &actual)
		t.AssertNil(err)
		t.Assert(actual, expected)
	})
}

// Test_DataType_Binary_Small tests small binary data
func Test_DataType_Binary_Small(t *testing.T) {
	table := rawTypeCreateTable("t_rt_bin", "DATA BLOB, DATA_RAW RAW(16)")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		binaryData := []byte{0x00, 0x01, 0x02, 0x03, 0xFF}
		_, err := db.Model(table).Data(g.Map{
			"data":     binaryData,
			"data_raw": binaryData,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(bytes.Equal(one["DATA"].Bytes(), binaryData), true)
		t.Assert(bytes.Equal(one["DATA_RAW"].Bytes(), binaryData), true)
	})
}

// Test_DataType_Binary_Large tests large binary data (1MB+)
func Test_DataType_Binary_Large(t *testing.T) {
	table := rawTypeCreateTable("t_rt_bin", "DATA BLOB")
	defer dropTable(table)

	size := 1024 * 1024
	largeBinary := make([]byte, size)
	for i := 0; i < size; i++ {
		largeBinary[i] = byte(i % 256)
	}

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"data": largeBinary,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(len(one["DATA"].Bytes()), size)
		t.Assert(bytes.Equal(one["DATA"].Bytes(), largeBinary), true)
	})

	gtest.C(t, func(t *gtest.T) {
		pattern := make([]byte, 256)
		for i := range pattern {
			pattern[i] = byte(i)
		}
		_, err := db.Exec(ctx, fmt.Sprintf(`DECLARE
    b     BLOB;
    chunk RAW(16384) := UTL_RAW.COPIES(HEXTORAW('%s'), 64);
BEGIN
    DBMS_LOB.CREATETEMPORARY(b, TRUE);
    FOR i IN 1..64 LOOP
        DBMS_LOB.WRITEAPPEND(b, 16384, chunk);
    END LOOP;
    INSERT INTO %s (ID, DATA) VALUES (100, b);
    DBMS_LOB.FREETEMPORARY(b);
END;`, hex.EncodeToString(pattern), table))
		t.AssertNil(err)

		stored, err := db.Model(table).Fields("DBMS_LOB.GETLENGTH(DATA) AS DATA_LENGTH").Where("id", 100).Value()
		t.AssertNil(err)
		t.Assert(stored.Int(), size)

		one, err := db.Model(table).Where("id", 100).One()
		t.AssertNil(err)
		t.Assert(len(one["DATA"].Bytes()), size)
		t.Assert(bytes.Equal(one["DATA"].Bytes(), largeBinary), true)
	})
}

// Test_DataType_Binary_Integrity tests binary data integrity with checksum
func Test_DataType_Binary_Integrity(t *testing.T) {
	table := rawTypeCreateTable("t_rt_bin", "DATA BLOB, CHECKSUM VARCHAR2(64)")
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
	})
}

// Test_DataType_Binary_Empty tests empty and NULL binary
func Test_DataType_Binary_Empty(t *testing.T) {
	table := rawTypeCreateTable("t_rt_bin", "DATA BLOB, DATA_RAW RAW(16)")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"data":     []byte{},
			"data_raw": []byte{},
		}).Insert()
		t.AssertNil(err)

		_, err = db.Model(table).Data(g.Map{
			"data":     nil,
			"data_raw": nil,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(len(one["DATA"].Bytes()), 0)
		t.Assert(one["DATA"].IsNil(), true)
		t.Assert(one["DATA_RAW"].IsNil(), true)

		one, err = db.Model(table).Where("id", 2).One()
		t.AssertNil(err)
		t.Assert(one["DATA"].IsNil(), true)
		t.Assert(one["DATA_RAW"].IsNil(), true)
	})
}

// Test_DataType_Decimal_HighPrecision tests high precision decimal (65,30)
func Test_DataType_Decimal_HighPrecision(t *testing.T) {
	table := rawTypeCreateTable("t_rt_dec", "AMOUNT NUMBER(38,30), AMOUNT_INT NUMBER(19), AMOUNT_ANY NUMBER")
	defer dropTable(table)

	var (
		value    = "12345678.123456789012345678901234567891"
		intValue = "9223372036854775807"
		anyValue = "12345678901234567890.123456789"
	)
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"amount":     value,
			"amount_int": intValue,
			"amount_any": anyValue,
		}).Insert()
		t.AssertNil(err)

		stored, err := db.Model(table).
			Fields("TO_CHAR(AMOUNT) AS S1, TO_CHAR(AMOUNT_INT) AS S2, TO_CHAR(AMOUNT_ANY) AS S3").
			Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(stored["S1"].String(), value)
		t.Assert(stored["S2"].String(), intValue)
		t.Assert(stored["S3"].String(), anyValue)
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["AMOUNT"].String(), value)
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["AMOUNT_INT"].String(), intValue)
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["AMOUNT_ANY"].String(), anyValue)
	})
}

// Test_DataType_Decimal_Calculation tests decimal arithmetic
func Test_DataType_Decimal_Calculation(t *testing.T) {
	table := rawTypeCreateTable("t_rt_dec", `PRICE NUMBER(10,2), QUANTITY NUMBER(10,2),
		PRICE_BF BINARY_FLOAT, QUANTITY_BF BINARY_FLOAT, PRICE_BD BINARY_DOUBLE, QUANTITY_BD BINARY_DOUBLE`)
	defer dropTable(table)

	var (
		priceBF, quantityBF = float32(19.99), float32(3.5)
		priceBD, quantityBD = 19.99, 3.5
	)
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"price":       "19.99",
			"quantity":    "3.5",
			"price_bf":    priceBF,
			"quantity_bf": quantityBF,
			"price_bd":    priceBD,
			"quantity_bd": quantityBD,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Fields("price * quantity as total").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["TOTAL"].String(), "69.965")
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Fields("price_bf * quantity_bf as total").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["TOTAL"].Float32(), priceBF*quantityBF)
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Fields("price_bd * quantity_bd as total").Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["TOTAL"].Float64(), priceBD*quantityBD)
	})
}

// Test_DataType_Decimal_Boundary tests decimal boundary values
func Test_DataType_Decimal_Boundary(t *testing.T) {
	table := rawTypeCreateTable("t_rt_dec", "VALUE NUMBER(10,2), VALUE_BF BINARY_FLOAT, VALUE_BD BINARY_DOUBLE")
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
		t.Assert(all[2]["VALUE"].String(), "0")
	})
	gtest.C(t, func(t *gtest.T) {
		for _, value := range []float32{math.MaxFloat32, -math.MaxFloat32, math.SmallestNonzeroFloat32} {
			id, err := db.Model(table).Data(g.Map{
				"value_bf": value,
			}).InsertAndGetId()
			t.AssertNil(err)

			one, err := db.Model(table).Where("id", id).One()
			t.AssertNil(err)
			t.Assert(one["VALUE_BF"].Float32(), value)
		}
	})
	gtest.C(t, func(t *gtest.T) {
		for _, value := range []float64{math.MaxFloat64, -math.MaxFloat64, math.SmallestNonzeroFloat64} {
			id, err := db.Model(table).Data(g.Map{
				"value_bd": value,
			}).InsertAndGetId()
			t.AssertNil(err)

			one, err := db.Model(table).Where("id", id).One()
			t.AssertNil(err)
			t.Assert(one["VALUE_BD"].Float64(), value)
		}
	})
}

// Test_DataType_Decimal_Null tests NULL decimal values
func Test_DataType_Decimal_Null(t *testing.T) {
	table := rawTypeCreateTable("t_rt_dec", "VALUE NUMBER(10,2), VALUE_BF BINARY_FLOAT, VALUE_BD BINARY_DOUBLE")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"value":    nil,
			"value_bf": nil,
			"value_bd": nil,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["VALUE"].IsNil(), true)
		t.Assert(one["VALUE_BF"].IsNil(), true)
		t.Assert(one["VALUE_BD"].IsNil(), true)
	})
}

// Test_DataType_Datetime_Timezone tests datetime with timezone handling
func Test_DataType_Datetime_Timezone(t *testing.T) {
	table := rawTypeCreateTable("t_rt_dt", "CREATED_AT TIMESTAMP WITH TIME ZONE")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		dt := "2024-01-15 12:30:45"
		_, err := db.Model(table).Data(g.Map{
			"created_at": gtime.NewFromStr(dt),
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["CREATED_AT"].String(), dt)
	})
	gtest.C(t, func(t *gtest.T) {
		_, localOffset := time.Now().Zone()
		dt := time.Date(2024, 1, 15, 12, 30, 45, 0, time.FixedZone("", localOffset+3*3600))
		id, err := db.Model(table).Data(g.Map{
			"created_at": dt,
		}).InsertAndGetId()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", id).One()
		t.AssertNil(err)
		t.Assert(one["CREATED_AT"].Time().Equal(dt), true)
	})
}

// Test_DataType_Datetime_Precision tests datetime with microsecond precision
func Test_DataType_Datetime_Precision(t *testing.T) {
	table := rawTypeCreateTable("t_rt_dt", "CREATED_AT TIMESTAMP(6)")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		dt := "2024-01-15 12:30:45.123456"
		_, err := db.Model(table).Data(g.Map{
			"created_at": gtime.NewFromStr(dt),
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		expected := "2024-01-15 12:30:45"
		actual := one["CREATED_AT"].String()[:19]
		t.Assert(actual, expected)
		t.Assert(one["CREATED_AT"].Time().Nanosecond(), 123456000)
	})
}

// Test_DataType_Datetime_Boundary tests datetime boundary values
func Test_DataType_Datetime_Boundary(t *testing.T) {
	table := rawTypeCreateTable("t_rt_dt", "DT TIMESTAMP")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"dt": gtime.NewFromStr("1000-01-01 00:00:00"),
		}).Insert()
		t.AssertNil(err)

		_, err = db.Model(table).Data(g.Map{
			"dt": gtime.NewFromStr("9999-12-31 23:59:59"),
		}).Insert()
		t.AssertNil(err)

		all, err := db.Model(table).Order("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["DT"].String(), "1000-01-01 00:00:00")
		t.Assert(all[1]["DT"].String(), "9999-12-31 23:59:59")
	})
}

// Test_DataType_Datetime_Null tests NULL datetime
func Test_DataType_Datetime_Null(t *testing.T) {
	table := rawTypeCreateTable("t_rt_dt", "DT TIMESTAMP")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"dt": nil,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["DT"].IsNil(), true)
	})
}

// Test_DataType_Datetime_Update tests datetime updates
func Test_DataType_Datetime_Update(t *testing.T) {
	table := rawTypeCreateTable("t_rt_dt", "DT TIMESTAMP, SHIFT INTERVAL DAY(3) TO SECOND, SHIFT_YM INTERVAL YEAR TO MONTH")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		dt1 := "2024-01-01 10:00:00"
		_, err := db.Model(table).Data(g.Map{
			"dt": gtime.NewFromStr(dt1),
		}).Insert()
		t.AssertNil(err)

		dt2 := "2024-12-31 23:59:59"
		_, err = db.Model(table).Data(g.Map{
			"dt": gtime.NewFromStr(dt2),
		}).Where("id", 1).Update()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["DT"].String(), dt2)
	})
	gtest.C(t, func(t *gtest.T) {
		id, err := db.Model(table).Data(g.Map{
			"dt":       gtime.NewFromStr("2024-01-01 10:00:00"),
			"shift":    "365 13:59:59",
			"shift_ym": "1-2",
		}).InsertAndGetId()
		t.AssertNil(err)

		_, err = db.Model(table).Data(g.Map{
			"dt": gdb.Raw("dt + shift"),
		}).Where("id", id).Update()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", id).One()
		t.AssertNil(err)
		t.Assert(one["DT"].String(), "2024-12-31 23:59:59")
		t.Assert(one["SHIFT"].String(), "+365 13:59:59.000000")
		t.Assert(one["SHIFT_YM"].String(), "+01-02")
	})
}

// Test_DataType_Enum_Valid tests valid ENUM values
func Test_DataType_Enum_Valid(t *testing.T) {
	table := rawTypeCreateTable("t_rt_enum", `STATUS VARCHAR2(20) CHECK (STATUS IN ('pending','approved','rejected')),
		STATUS_CHAR CHAR(8) CHECK (STATUS_CHAR IN ('pending','approved','rejected')),
		STATUS_NCHAR NCHAR(8) CHECK (STATUS_NCHAR IN ('pending','approved','rejected'))`)
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
	table := rawTypeCreateTable("t_rt_enum", "STATUS VARCHAR2(20) CHECK (STATUS IN ('pending','approved','rejected'))")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"status": "invalid_status",
		}).Insert()
		t.AssertNE(err, nil)
		t.AssertIN("ORA-02290", err.Error())

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}

// Test_DataType_Set_Valid tests valid SET values
func Test_DataType_Set_Valid(t *testing.T) {
	table := rawTypeCreateTable("t_rt_set", `PERMISSIONS VARCHAR2(30)
		CHECK (REGEXP_LIKE(PERMISSIONS, '^(read|write|execute)(,(read|write|execute))*$')),
		PERMISSIONS_N NVARCHAR2(30)
		CHECK (REGEXP_LIKE(PERMISSIONS_N, '^(read|write|execute)(,(read|write|execute))*$'))`)
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
}

// Test_DataType_Set_Empty tests empty SET values
func Test_DataType_Set_Empty(t *testing.T) {
	table := rawTypeCreateTable("t_rt_set", "PERMISSIONS VARCHAR2(30) CHECK (REGEXP_LIKE(PERMISSIONS, '^(read|write|execute)(,(read|write|execute))*$'))")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"permissions": "",
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["PERMISSIONS"].String(), "")
		t.Assert(one["PERMISSIONS"].IsNil(), true)
	})
}

// Test_DataType_Geometry_Point tests POINT geometry type
func Test_DataType_Geometry_Point(t *testing.T) {
	table := rawTypeCreateTable("t_rt_geo", "LOCATION SDO_GEOMETRY")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Exec(ctx, fmt.Sprintf("INSERT INTO %s (LOCATION) VALUES (SDO_GEOMETRY(2001, NULL, SDO_POINT_TYPE(116.4074, 39.9042, NULL), NULL, NULL))", table))
		t.AssertNil(err)

		one, err := db.Model(table+" T").
			Fields("T.LOCATION.SDO_GTYPE AS GTYPE, T.LOCATION.SDO_POINT.X AS X, T.LOCATION.SDO_POINT.Y AS Y").
			Where("ID", 1).One()
		t.AssertNil(err)
		t.Assert(one["GTYPE"].Int(), 2001)
		t.Assert(one["X"].String(), "116.4074")
		t.Assert(one["Y"].String(), "39.9042")
	})
}

// Test_DataType_Geometry_Polygon tests POLYGON geometry type
func Test_DataType_Geometry_Polygon(t *testing.T) {
	table := rawTypeCreateTable("t_rt_geo", "AREA SDO_GEOMETRY")
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Exec(ctx, fmt.Sprintf(
			"INSERT INTO %s (AREA) VALUES (SDO_GEOMETRY(2003, NULL, NULL, SDO_ELEM_INFO_ARRAY(1, 1003, 1), SDO_ORDINATE_ARRAY(0, 0, 10, 0, 10, 10, 0, 10, 0, 0)))",
			table,
		))
		t.AssertNil(err)

		one, err := db.Model(table+" T").
			Fields("T.AREA.SDO_GTYPE AS GTYPE, SDO_GEOM.SDO_AREA(T.AREA, 0.005) AS AREA_SIZE").
			Where("ID", 1).One()
		t.AssertNil(err)
		t.Assert(one["GTYPE"].Int(), 2003)
		t.Assert(one["AREA_SIZE"].Float64(), 100)

		vertices, err := db.GetAll(ctx, fmt.Sprintf(
			"SELECT V.X, V.Y FROM %s T, TABLE(SDO_UTIL.GETVERTICES(T.AREA)) V WHERE T.ID = 1 ORDER BY V.ID",
			table,
		))
		t.AssertNil(err)
		points := make([]string, 0, len(vertices))
		for _, vertex := range vertices {
			points = append(points, vertex["X"].String()+" "+vertex["Y"].String())
		}
		expected := "POLYGON((0 0,10 0,10 10,0 10,0 0))"
		actual := "POLYGON((" + strings.Join(points, ",") + "))"
		t.Assert(actual, expected)
	})
}

// Test_DataType_Geometry_Null tests NULL geometry values
func Test_DataType_Geometry_Null(t *testing.T) {
	table := rawTypeCreateTable("t_rt_geo", "LOCATION SDO_GEOMETRY")
	defer dropTable(table)

	// Note: inserting nil through the model is not asserted. go-ora v2.9.0 binds a nil value as a
	// NCHAR NULL, which Oracle rejects for an object column (ORA-00932), see
	// https://github.com/sijms/go-ora/blob/v2.9.0/v2/parameter_encode.go#L20-L21.

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Exec(ctx, fmt.Sprintf("INSERT INTO %s (ID, LOCATION) VALUES (100, NULL)", table))
		t.AssertNil(err)

		count, err := db.Model(table).Where("id", 100).WhereNull("location").Count()
		t.AssertNil(err)
		t.Assert(count, 1)
		// Note: reading the NULL back with One() is not asserted. go-ora v2.9.0 fails to decode an
		// SDO_GEOMETRY column, even a NULL one ("invalid size for GetInt64"), see
		// https://github.com/sijms/go-ora/blob/v2.9.0/v2/network/session.go#L1608.
	})
}

func Test_Raw_Save(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"id":          1,
			"passport":    gdb.Raw("'port_' || 'raw'"),
			"password":    "pass_1",
			"nickname":    gdb.Raw("UPPER('name_raw')"),
			"create_time": gdb.Raw("TO_CHAR(SYSDATE, 'YYYY-MM-DD HH24:MI:SS')"),
		}).Save()
		t.AssertNil(err)

		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		t.Assert(one["PASSPORT"], "port_raw")
		t.Assert(strings.TrimRight(one["PASSWORD"].String(), " "), "pass_1")
		t.Assert(one["NICKNAME"], "NAME_RAW")
		t.Assert(gregex.IsMatchString(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`, one["CREATE_TIME"].String()), true)
	})
}

// Test_DataType_Text_Large tests character LOB values longer than a VARCHAR2 bind
func Test_DataType_Text_Large(t *testing.T) {
	table := rawTypeCreateTable("t_rt_text", "DATA CLOB, DATA_NCLOB NCLOB")
	defer dropTable(table)

	var (
		text      = strings.Repeat("0123456789", 5000)
		nationals = strings.Repeat("中文", 3000)
	)
	gtest.C(t, func(t *gtest.T) {
		id, err := db.Model(table).Data(g.Map{
			"data":       text,
			"data_nclob": nationals,
		}).InsertAndGetId()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", id).One()
		t.AssertNil(err)
		t.Assert(one["DATA"].String(), text)
		t.Assert(one["DATA_NCLOB"].String(), nationals)
	})
	gtest.C(t, func(t *gtest.T) {
		_, err := db.Model(table).Data(g.Map{
			"data":       nationals,
			"data_nclob": text,
		}).Where("id", 1).Update()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["DATA"].String(), nationals)
		t.Assert(one["DATA_NCLOB"].String(), text)
	})
}
