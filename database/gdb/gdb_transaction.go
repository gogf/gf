// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/gogf/gf/text/gregex"
)

// TX is the struct for transaction management.
type TX struct {
	db     DB
	tx     *sql.Tx
	master *sql.DB
}

// Commit commits the transaction.
func (tx *TX) Commit() error {
	return tx.tx.Commit()
}

// Rollback aborts the transaction.
func (tx *TX) Rollback() error {
	return tx.tx.Rollback()
}

// Query does query operation on transaction.
// See dbBase.Query.
func (tx *TX) Query(query string, args ...interface{}) (rows *sql.Rows, err error) {
	return tx.db.doQuery(tx.tx, query, args...)
}

// Exec does none query operation on transaction.
// See dbBase.Exec.
func (tx *TX) Exec(query string, args ...interface{}) (sql.Result, error) {
	return tx.db.doExec(tx.tx, query, args...)
}

// Prepare creates a prepared statement for later queries or executions.
// Multiple queries or executions may be run concurrently from the
// returned statement.
// The caller must call the statement's Close method
// when the statement is no longer needed.
func (tx *TX) Prepare(query string) (*sql.Stmt, error) {
	return tx.db.doPrepare(tx.tx, query)
}

// GetAll queries and returns data records from database.
func (tx *TX) GetAll(query string, args ...interface{}) (Result, error) {
	rows, err := tx.Query(query, args...)
	if err != nil || rows == nil {
		return nil, err
	}
	defer rows.Close()
	return tx.db.rowsToResult(rows)
}

// GetOne queries and returns one record from database.
func (tx *TX) GetOne(query string, args ...interface{}) (Record, error) {
	list, err := tx.GetAll(query, args...)
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		return list[0], nil
	}
	return nil, nil
}

// GetStruct queries one record from database and converts it to given struct.
// The parameter <pointer> should be a pointer to struct.
func (tx *TX) GetStruct(obj interface{}, query string, args ...interface{}) error {
	one, err := tx.GetOne(query, args...)
	if err != nil {
		return err
	}
	return one.Struct(obj)
}

// GetStructs queries records from database and converts them to given struct.
// The parameter <pointer> should be type of struct slice: []struct/[]*struct.
func (tx *TX) GetStructs(objPointerSlice interface{}, query string, args ...interface{}) error {
	all, err := tx.GetAll(query, args...)
	if err != nil {
		return err
	}
	return all.Structs(objPointerSlice)
}

// GetScan queries one or more records from database and converts them to given struct or
// struct array.
//
// If parameter <pointer> is type of struct pointer, it calls GetStruct internally for
// the conversion. If parameter <pointer> is type of slice, it calls GetStructs internally
// for conversion.
func (tx *TX) GetScan(objPointer interface{}, query string, args ...interface{}) error {
	t := reflect.TypeOf(objPointer)
	k := t.Kind()
	if k != reflect.Ptr {
		return fmt.Errorf("params should be type of pointer, but got: %v", k)
	}
	k = t.Elem().Kind()
	switch k {
	case reflect.Array, reflect.Slice:
		return tx.db.GetStructs(objPointer, query, args...)
	case reflect.Struct:
		return tx.db.GetStruct(objPointer, query, args...)
	default:
		return fmt.Errorf("element type should be type of struct/slice, unsupported: %v", k)
	}
	return nil
}

// GetValue queries and returns the field value from database.
// The sql should queries only one field from database, or else it returns only one
// field of the result.
func (tx *TX) GetValue(query string, args ...interface{}) (Value, error) {
	one, err := tx.GetOne(query, args...)
	if err != nil {
		return nil, err
	}
	for _, v := range one {
		return v, nil
	}
	return nil, nil
}

// GetCount queries and returns the count from database.
func (tx *TX) GetCount(query string, args ...interface{}) (int, error) {
	if !gregex.IsMatchString(`(?i)SELECT\s+COUNT\(.+\)\s+FROM`, query) {
		query, _ = gregex.ReplaceString(`(?i)(SELECT)\s+(.+)\s+(FROM)`, `$1 COUNT($2) $3`, query)
	}
	value, err := tx.GetValue(query, args...)
	if err != nil {
		return 0, err
	}
	return value.Int(), nil
}

// Insert does "INSERT INTO ..." statement for the table.
// If there's already one unique record of the data in the table, it returns error.
//
// The parameter <data> can be type of map/gmap/struct/*struct/[]map/[]struct, etc.
// Eg:
// Data(g.Map{"uid": 10000, "name":"john"})
// Data(g.Slice{g.Map{"uid": 10000, "name":"john"}, g.Map{"uid": 20000, "name":"smith"})
//
// The parameter <batch> specifies the batch operation count when given data is slice.
func (tx *TX) Insert(table string, data interface{}, batch ...int) (sql.Result, error) {
	return tx.db.doInsert(tx.tx, table, data, gINSERT_OPTION_DEFAULT, batch...)
}

// InsertIgnore does "INSERT IGNORE INTO ..." statement for the table.
// If there's already one unique record of the data in the table, it ignores the inserting.
//
// The parameter <data> can be type of map/gmap/struct/*struct/[]map/[]struct, etc.
// Eg:
// Data(g.Map{"uid": 10000, "name":"john"})
// Data(g.Slice{g.Map{"uid": 10000, "name":"john"}, g.Map{"uid": 20000, "name":"smith"})
//
// The parameter <batch> specifies the batch operation count when given data is slice.
func (tx *TX) InsertIgnore(table string, data interface{}, batch ...int) (sql.Result, error) {
	return tx.db.doInsert(tx.tx, table, data, gINSERT_OPTION_IGNORE, batch...)
}

// Replace does "REPLACE INTO ..." statement for the table.
// If there's already one unique record of the data in the table, it deletes the record
// and inserts a new one.
//
// The parameter <data> can be type of map/gmap/struct/*struct/[]map/[]struct, etc.
// Eg:
// Data(g.Map{"uid": 10000, "name":"john"})
// Data(g.Slice{g.Map{"uid": 10000, "name":"john"}, g.Map{"uid": 20000, "name":"smith"})
//
// The parameter <data> can be type of map/gmap/struct/*struct/[]map/[]struct, etc.
// If given data is type of slice, it then does batch replacing, and the optional parameter
// <batch> specifies the batch operation count.
func (tx *TX) Replace(table string, data interface{}, batch ...int) (sql.Result, error) {
	return tx.db.doInsert(tx.tx, table, data, gINSERT_OPTION_REPLACE, batch...)
}

// Save does "INSERT INTO ... ON DUPLICATE KEY UPDATE..." statement for the table.
// It updates the record if there's primary or unique index in the saving data,
// or else it inserts a new record into the table.
//
// The parameter <data> can be type of map/gmap/struct/*struct/[]map/[]struct, etc.
// Eg:
// Data(g.Map{"uid": 10000, "name":"john"})
// Data(g.Slice{g.Map{"uid": 10000, "name":"john"}, g.Map{"uid": 20000, "name":"smith"})
//
// If given data is type of slice, it then does batch saving, and the optional parameter
// <batch> specifies the batch operation count.
func (tx *TX) Save(table string, data interface{}, batch ...int) (sql.Result, error) {
	return tx.db.doInsert(tx.tx, table, data, gINSERT_OPTION_SAVE, batch...)
}

// BatchInsert batch inserts data.
// The parameter <list> must be type of slice of map or struct.
func (tx *TX) BatchInsert(table string, list interface{}, batch ...int) (sql.Result, error) {
	return tx.db.doBatchInsert(tx.tx, table, list, gINSERT_OPTION_DEFAULT, batch...)
}

// BatchInsert batch inserts data with ignore option.
// The parameter <list> must be type of slice of map or struct.
func (tx *TX) BatchInsertIgnore(table string, list interface{}, batch ...int) (sql.Result, error) {
	return tx.db.doBatchInsert(tx.tx, table, list, gINSERT_OPTION_IGNORE, batch...)
}

// BatchReplace batch replaces data.
// The parameter <list> must be type of slice of map or struct.
func (tx *TX) BatchReplace(table string, list interface{}, batch ...int) (sql.Result, error) {
	return tx.db.doBatchInsert(tx.tx, table, list, gINSERT_OPTION_REPLACE, batch...)
}

// BatchSave batch replaces data.
// The parameter <list> must be type of slice of map or struct.
func (tx *TX) BatchSave(table string, list interface{}, batch ...int) (sql.Result, error) {
	return tx.db.doBatchInsert(tx.tx, table, list, gINSERT_OPTION_SAVE, batch...)
}

// Update does "UPDATE ... " statement for the table.
//
// The parameter <data> can be type of string/map/gmap/struct/*struct, etc.
// Eg: "uid=10000", "uid", 10000, g.Map{"uid": 10000, "name":"john"}
//
// The parameter <condition> can be type of string/map/gmap/slice/struct/*struct, etc.
// It is commonly used with parameter <args>.
// Eg:
// "uid=10000",
// "uid", 10000
// "money>? AND name like ?", 99999, "vip_%"
// "status IN (?)", g.Slice{1,2,3}
// "age IN(?,?)", 18, 50
// User{ Id : 1, UserName : "john"}
func (tx *TX) Update(table string, data interface{}, condition interface{}, args ...interface{}) (sql.Result, error) {
	newWhere, newArgs := formatWhere(tx.db, condition, args, false)
	if newWhere != "" {
		newWhere = " WHERE " + newWhere
	}
	return tx.db.doUpdate(tx.tx, table, data, newWhere, newArgs...)
}

// Delete does "DELETE FROM ... " statement for the table.
//
// The parameter <condition> can be type of string/map/gmap/slice/struct/*struct, etc.
// It is commonly used with parameter <args>.
// Eg:
// "uid=10000",
// "uid", 10000
// "money>? AND name like ?", 99999, "vip_%"
// "status IN (?)", g.Slice{1,2,3}
// "age IN(?,?)", 18, 50
// User{ Id : 1, UserName : "john"}
func (tx *TX) Delete(table string, condition interface{}, args ...interface{}) (sql.Result, error) {
	newWhere, newArgs := formatWhere(tx.db, condition, args, false)
	if newWhere != "" {
		newWhere = " WHERE " + newWhere
	}
	return tx.db.doDelete(tx.tx, table, newWhere, newArgs...)
}
