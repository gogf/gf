// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package pgsql_test

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

// https://github.com/gogf/gf/issues/3330
func Test_Issue3330(t *testing.T) {
	var (
		table      = fmt.Sprintf(`%s_%d`, TablePrefix+"test", gtime.TimestampNano())
		uniqueName = fmt.Sprintf(`%s_%d`, TablePrefix+"test_unique", gtime.TimestampNano())
	)
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
		   	id bigserial  NOT NULL,
		   	passport varchar(45) NOT NULL,
		   	password varchar(32) NOT NULL,
		   	nickname varchar(45) NOT NULL,
		   	create_time timestamp NOT NULL,
		   	PRIMARY KEY (id),
			CONSTRAINT %s unique ("password")
		) ;`, table, uniqueName,
	)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var (
			list []map[string]any
			one  gdb.Record
			err  error
		)

		fields, err := db.TableFields(ctx, table)
		t.AssertNil(err)

		t.Assert(fields["id"].Key, "pri")
		t.Assert(fields["password"].Key, "uni")

		for i := 1; i <= 10; i++ {
			list = append(list, g.Map{
				"id":          i,
				"passport":    fmt.Sprintf("p%d", i),
				"password":    fmt.Sprintf("pw%d", i),
				"nickname":    fmt.Sprintf("n%d", i),
				"create_time": "2016-06-01 00:00:00",
			})
		}

		_, err = db.Model(table).Data(list).Insert()
		t.AssertNil(err)

		for i := 1; i <= 10; i++ {
			one, err = db.Model(table).WherePri(i).One()
			t.AssertNil(err)
			t.Assert(one["id"], list[i-1]["id"])
			t.Assert(one["passport"], list[i-1]["passport"])
			t.Assert(one["password"], list[i-1]["password"])
			t.Assert(one["nickname"], list[i-1]["nickname"])
		}
	})
}

// https://github.com/gogf/gf/issues/3632
func Test_Issue3632(t *testing.T) {
	type Member struct {
		One []int64    `json:"one" orm:"one"`
		Two [][]string `json:"two" orm:"two"`
	}
	var (
		sqlText = gtest.DataContent("issues", "issue3632.sql")
		table   = fmt.Sprintf(`%s_%d`, TablePrefix+"issue3632", gtime.TimestampNano())
	)
	if _, err := db.Exec(ctx, fmt.Sprintf(sqlText, table)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var (
			dao    = db.Model(table)
			member = Member{
				One: []int64{1, 2, 3},
				Two: [][]string{{"a", "b"}, {"c", "d"}},
			}
		)

		_, err := dao.Ctx(ctx).Data(&member).Insert()
		t.AssertNil(err)
	})
}

// https://github.com/gogf/gf/issues/3671
func Test_Issue3671(t *testing.T) {
	type SubMember struct {
		Seven string
		Eight int64
	}
	type Member struct {
		One   []int64     `json:"one" orm:"one"`
		Two   [][]string  `json:"two" orm:"two"`
		Three []string    `json:"three" orm:"three"`
		Four  []int64     `json:"four" orm:"four"`
		Five  []SubMember `json:"five" orm:"five"`
	}
	var (
		sqlText = gtest.DataContent("issues", "issue3671.sql")
		table   = fmt.Sprintf(`%s_%d`, TablePrefix+"issue3632", gtime.TimestampNano())
	)
	if _, err := db.Exec(ctx, fmt.Sprintf(sqlText, table)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var (
			dao    = db.Model(table)
			member = Member{
				One:   []int64{1, 2, 3},
				Two:   [][]string{{"a", "b"}, {"c", "d"}},
				Three: []string{"x", "y", "z"},
				Four:  []int64{1, 2, 3},
				Five:  []SubMember{{Seven: "1", Eight: 2}, {Seven: "3", Eight: 4}},
			}
		)

		_, err := dao.Ctx(ctx).Data(&member).Insert()
		t.AssertNil(err)
	})
}

// https://github.com/gogf/gf/issues/3668
func Test_Issue3668(t *testing.T) {
	type Issue3668 struct {
		Text   any
		Number any
	}
	var (
		sqlText = gtest.DataContent("issues", "issue3668.sql")
		table   = fmt.Sprintf(`%s_%d`, TablePrefix+"issue3668", gtime.TimestampNano())
	)
	if _, err := db.Exec(ctx, fmt.Sprintf(sqlText, table)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var (
			dao  = db.Model(table)
			data = Issue3668{
				Text:   "我们都是自然的婴儿，卧在宇宙的摇篮里",
				Number: nil,
			}
		)
		_, err := dao.Ctx(ctx).
			Data(data).
			Insert()
		t.AssertNil(err)
	})
}

type Issue4033Status int

const (
	Issue4033StatusA Issue4033Status = 1
)

func (s Issue4033Status) String() string {
	return "somevalue"
}

func (s Issue4033Status) Int64() int64 {
	return int64(s)
}

// https://github.com/gogf/gf/issues/4033
func Test_Issue4033(t *testing.T) {
	var (
		sqlText = gtest.DataContent("issues", "issue4033.sql")
		table   = "test_enum"
	)
	if _, err := db.Exec(ctx, sqlText); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		query := g.Map{
			"status": g.Slice{Issue4033StatusA},
		}
		_, err := db.Model(table).Ctx(ctx).Where(query).All()
		t.AssertNil(err)
	})
}

// https://github.com/gogf/gf/issues/4500
// Raw() Count ignores Where condition
func Test_Issue4500(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	// Test 1: Raw SQL with WHERE + external Where condition + Count
	// This tests that formatCondition correctly uses AND when Raw SQL already has WHERE
	gtest.C(t, func(t *gtest.T) {
		count, err := db.
			Raw(fmt.Sprintf("SELECT * FROM %s WHERE id IN (?)", table), g.Slice{1, 5, 7, 8, 9, 10}).
			WhereLT("id", 8).
			Count()
		t.AssertNil(err)
		// Raw SQL: id IN (1,5,7,8,9,10) = 6 records
		// Where: id < 8 filters to {1,5,7} = 3 records
		t.Assert(count, 3)
	})

	// Test 2: Raw SQL without WHERE + external Where condition + Count
	// This tests that formatCondition correctly adds WHERE
	gtest.C(t, func(t *gtest.T) {
		count, err := db.
			Raw(fmt.Sprintf("SELECT * FROM %s", table)).
			WhereLT("id", 5).
			Count()
		t.AssertNil(err)
		// Raw SQL: all 10 records
		// Where: id < 5 = {1,2,3,4} = 4 records
		t.Assert(count, 4)
	})

	// Test 3: Raw + Where + ScanAndCount
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id       int
			Passport string
		}
		var users []User
		var total int
		err := db.
			Raw(fmt.Sprintf("SELECT * FROM %s WHERE id IN (?)", table), g.Slice{1, 5, 7, 8, 9, 10}).
			WhereLT("id", 8).
			ScanAndCount(&users, &total, false)
		t.AssertNil(err)
		// Both scan result and count should respect Where condition
		t.Assert(len(users), 3)
		t.Assert(total, 3)
	})

	// Test 4: Raw + multiple Where conditions + Count
	gtest.C(t, func(t *gtest.T) {
		count, err := db.
			Raw(fmt.Sprintf("SELECT * FROM %s WHERE id > ?", table), 0).
			WhereLT("id", 5).
			WhereGTE("id", 2).
			Count()
		t.AssertNil(err)
		// Raw: id > 0 (all 10 records)
		// Where: id < 5 AND id >= 2 = {2, 3, 4} = 3 records
		t.Assert(count, 3)
	})

	// Test 5: Raw SQL with no external Where + Count (baseline test)
	gtest.C(t, func(t *gtest.T) {
		count, err := db.
			Raw(fmt.Sprintf("SELECT * FROM %s WHERE id IN (?)", table), g.Slice{1, 2, 3}).
			Count()
		t.AssertNil(err)
		// Should count 3 records
		t.Assert(count, 3)
	})

	// Test 6: Verify All() still works correctly with Raw + Where
	gtest.C(t, func(t *gtest.T) {
		all, err := db.
			Raw(fmt.Sprintf("SELECT * FROM %s WHERE id IN (?)", table), g.Slice{1, 5, 7, 8, 9, 10}).
			WhereLT("id", 8).
			All()
		t.AssertNil(err)
		t.Assert(len(all), 3)
	})
}

// https://github.com/gogf/gf/issues/4677
// record.Get().Bytes() corrupts bytea data on retrieval from PostgreSQL.
func Test_Issue4677(t *testing.T) {
	table := fmt.Sprintf(`%s_%d`, TablePrefix+"issue4677", gtime.TimestampNano())
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id bigserial PRIMARY KEY,
			bin_data bytea
		);`, table,
	)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Test 1: Binary data with various byte values including 0x00, 0x5D(']'), 0x5B('[')
		originalBytes := []byte{
			0xDE, 0xAD, 0xBE, 0xEF, 0x00, 0x01, 0x5B, 0x5D,
			0xFF, 0x7B, 0x7D, 0x80, 0xCA, 0xFE, 0xBA, 0xBE,
		}

		_, err := db.Model(table).Data(g.Map{
			"bin_data": originalBytes,
		}).Insert()
		t.AssertNil(err)

		record, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)

		retrievedBytes := record["bin_data"].Bytes()
		t.Assert(len(retrievedBytes), len(originalBytes))
		t.Assert(retrievedBytes, originalBytes)
	})

	gtest.C(t, func(t *gtest.T) {
		// Test 2: Larger binary data (simulating gob/protobuf encoded payload)
		largeBytes := make([]byte, 1024)
		for i := range largeBytes {
			largeBytes[i] = byte(i % 256)
		}

		_, err := db.Model(table).Data(g.Map{
			"bin_data": largeBytes,
		}).Insert()
		t.AssertNil(err)

		record, err := db.Model(table).OrderDesc("id").One()
		t.AssertNil(err)

		retrievedBytes := record["bin_data"].Bytes()
		t.Assert(len(retrievedBytes), len(largeBytes))
		t.Assert(retrievedBytes, largeBytes)
	})
}

// https://github.com/gogf/gf/issues/4231
// ConvertValueForField corrupts bytea data containing 0x5D on write.
func Test_Issue4231(t *testing.T) {
	table := fmt.Sprintf(`%s_%d`, TablePrefix+"issue4231", gtime.TimestampNano())
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id bigserial PRIMARY KEY,
			bin_data bytea
		);`, table,
	)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Bytes containing 0x5D (ASCII ']') which was being converted to 0x7D ('}')
		originalBytes := []byte{0x01, 0x5D, 0x02, 0x5B, 0x03}

		_, err := db.Model(table).Data(g.Map{
			"bin_data": originalBytes,
		}).Insert()
		t.AssertNil(err)

		record, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)

		retrievedBytes := record["bin_data"].Bytes()
		t.Assert(len(retrievedBytes), len(originalBytes))
		t.Assert(retrievedBytes, originalBytes)
	})
}

// https://github.com/gogf/gf/issues/4595
// FieldsPrefix silently drops fields when using table alias before LeftJoin.
func Test_Issue4595(t *testing.T) {
	var (
		tableUser       = fmt.Sprintf(`%s_%d`, TablePrefix+"issue4595_user", gtime.TimestampNano())
		tableUserDetail = fmt.Sprintf(`%s_%d`, TablePrefix+"issue4595_user_detail", gtime.TimestampNano())
	)

	// Create user table
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id bigserial PRIMARY KEY,
			name varchar(100),
			email varchar(100)
		);`, tableUser,
	)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(tableUser)

	// Create user_detail table
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id bigserial PRIMARY KEY,
			user_id bigint,
			phone varchar(20),
			address varchar(200)
		);`, tableUserDetail,
	)); err != nil {
		gtest.Fatal(err)
	}
	defer dropTable(tableUserDetail)

	// Insert test data
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		INSERT INTO %s (id, name, email) VALUES (1, 'john', 'john@example.com');
		INSERT INTO %s (id, user_id, phone, address) VALUES (1, 1, '1234567890', '123 Main St');
	`, tableUser, tableUserDetail)); err != nil {
		gtest.Fatal(err)
	}

	gtest.C(t, func(t *gtest.T) {
		// Test case 1: FieldsPrefix called before LeftJoin
		// Both t1 and t2 fields should be present
		r, err := db.Model(tableUser).As("t1").
			FieldsPrefix("t2", "phone", "address").
			FieldsPrefix("t1", "id", "name", "email").
			LeftJoin(tableUserDetail, "t2", "t1.id=t2.user_id").
			All()

		t.AssertNil(err)
		t.Assert(len(r), 1)
		t.Assert(r[0]["id"], 1)
		t.Assert(r[0]["name"], "john")
		t.Assert(r[0]["email"], "john@example.com")
		t.Assert(r[0]["phone"], "1234567890")
		t.Assert(r[0]["address"], "123 Main St")
	})

	gtest.C(t, func(t *gtest.T) {
		// Test case 2: Using Fields() with prefix
		r, err := db.Model(tableUser).As("t1").
			Fields("t2.phone", "t2.address", "t1.id", "t1.name", "t1.email").
			LeftJoin(tableUserDetail, "t2", "t1.id=t2.user_id").
			All()
		t.AssertNil(err)
		t.Assert(len(r), 1)
		t.Assert(r[0]["id"], 1)
		t.Assert(r[0]["name"], "john")
		t.Assert(r[0]["email"], "john@example.com")
		t.Assert(r[0]["phone"], "1234567890")
		t.Assert(r[0]["address"], "123 Main St")
	})

	gtest.C(t, func(t *gtest.T) {
		// Test case 3: FieldsPrefix called after LeftJoin
		r, err := db.Model(tableUser).As("t1").
			LeftJoin(tableUserDetail, "t2", "t1.id=t2.user_id").
			FieldsPrefix("t2", "phone", "address").
			FieldsPrefix("t1", "id", "name", "email").
			All()
		t.AssertNil(err)
		t.Assert(len(r), 1)
		t.Assert(r[0]["id"], 1)
		t.Assert(r[0]["name"], "john")
		t.Assert(r[0]["email"], "john@example.com")
		t.Assert(r[0]["phone"], "1234567890")
		t.Assert(r[0]["address"], "123 Main St")
	})
}
