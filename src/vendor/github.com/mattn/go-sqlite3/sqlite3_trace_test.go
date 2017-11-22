// Copyright (C) 2016 Yasuhiro Matsumoto <mattn.jp@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
// +build trace

package sqlite3

import (
	"database/sql"
	"testing"
)

type sumAggregator int64

func (s *sumAggregator) Step(x int64) {
	*s += sumAggregator(x)
}

func (s *sumAggregator) Done() int64 {
	return int64(*s)
}

func TestAggregatorRegistration(t *testing.T) {
	customSum := func() *sumAggregator {
		var ret sumAggregator
		return &ret
	}

	sql.Register("sqlite3_AggregatorRegistration", &SQLiteDriver{
		ConnectHook: func(conn *SQLiteConn) error {
			if err := conn.RegisterAggregator("customSum", customSum, true); err != nil {
				return err
			}
			return nil
		},
	})
	db, err := sql.Open("sqlite3_AggregatorRegistration", ":memory:")
	if err != nil {
		t.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	_, err = db.Exec("create table foo (department integer, profits integer)")
	if err != nil {
		// trace feature is not implemented
		t.Skip("Failed to create table:", err)
	}

	_, err = db.Exec("insert into foo values (1, 10), (1, 20), (2, 42)")
	if err != nil {
		t.Fatal("Failed to insert records:", err)
	}

	tests := []struct {
		dept, sum int64
	}{
		{1, 30},
		{2, 42},
	}

	for _, test := range tests {
		var ret int64
		err = db.QueryRow("select customSum(profits) from foo where department = $1 group by department", test.dept).Scan(&ret)
		if err != nil {
			t.Fatal("Query failed:", err)
		}
		if ret != test.sum {
			t.Fatalf("Custom sum returned wrong value, got %d, want %d", ret, test.sum)
		}
	}
}
