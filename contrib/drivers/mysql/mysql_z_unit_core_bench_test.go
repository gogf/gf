// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// go test *.go -bench=".*"

package mysql_test

import (
	"testing"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/grand"
)

func Benchmark_BatchInsert(b *testing.B) {
	table := createTable()
	defer dropTable(table)
	type User struct {
		Id         int         `c:"id"`
		Passport   string      `c:"passport"`
		Password   string      `c:"password"`
		NickName   string      `c:"nickname"`
		CreateTime *gtime.Time `c:"create_time"`
	}
	var users []*User
	for i := 0; i < 10000; i++ {
		users = append(users, &User{
			Passport:   grand.S(10),
			Password:   grand.S(10),
			NickName:   grand.S(10),
			CreateTime: gtime.Now(),
		})
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result, err := db.Insert(ctx, table, users)
		if err != nil {
			b.Fatalf("insert error: %v", err)
		}
		n, _ := result.RowsAffected()
		b.Logf("insert %d rows", n)
	}
}
