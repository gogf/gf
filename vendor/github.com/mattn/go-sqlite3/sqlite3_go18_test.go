// Copyright (C) 2014 Yasuhiro Matsumoto <mattn.jp@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// +build go1.8

package sqlite3

import (
	"database/sql"
	"os"
	"testing"
)

func TestNamedParams(t *testing.T) {
	tempFilename := TempFilename(t)
	defer os.Remove(tempFilename)
	db, err := sql.Open("sqlite3", tempFilename)
	if err != nil {
		t.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	_, err = db.Exec(`
	create table foo (id integer, name text, extra text);
	`)
	if err != nil {
		t.Error("Failed to call db.Query:", err)
	}

	_, err = db.Exec(`insert into foo(id, name, extra) values(:id, :name, :name)`, sql.Named("name", "foo"), sql.Named("id", 1))
	if err != nil {
		t.Error("Failed to call db.Exec:", err)
	}

	row := db.QueryRow(`select id, extra from foo where id = :id and extra = :extra`, sql.Named("id", 1), sql.Named("extra", "foo"))
	if row == nil {
		t.Error("Failed to call db.QueryRow")
	}
	var id int
	var extra string
	err = row.Scan(&id, &extra)
	if err != nil {
		t.Error("Failed to db.Scan:", err)
	}
	if id != 1 || extra != "foo" {
		t.Error("Failed to db.QueryRow: not matched results")
	}
}
