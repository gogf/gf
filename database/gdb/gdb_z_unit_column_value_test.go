// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file tests result value conversion with driver-provided column metadata.

package gdb

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

// mockColumnValueConnector supplies column metadata without a registered SQL driver.
type mockColumnValueConnector struct {
	scanType     reflect.Type
	databaseType string
}

// Connect creates a connection that exposes the configured column metadata.
func (c *mockColumnValueConnector) Connect(_ context.Context) (driver.Conn, error) {
	return &mockColumnValueConn{columnType: c}, nil
}

// Driver returns the driver associated with the connector.
func (c *mockColumnValueConnector) Driver() driver.Driver { return c }

// Open implements driver.Driver; sql.OpenDB uses Connect instead.
func (c *mockColumnValueConnector) Open(_ string) (driver.Conn, error) {
	return c.Connect(context.Background())
}

// mockColumnValueConn only supports the metadata query used by the tests.
type mockColumnValueConn struct {
	columnType *mockColumnValueConnector
}

// Prepare is unused because the connection implements QueryContext.
func (c *mockColumnValueConn) Prepare(_ string) (driver.Stmt, error) {
	return nil, driver.ErrSkip
}

// Close releases the connection, which holds no external resources.
func (c *mockColumnValueConn) Close() error { return nil }

// Begin rejects transactions, which are not needed for metadata queries.
func (c *mockColumnValueConn) Begin() (driver.Tx, error) { return nil, driver.ErrSkip }

// QueryContext returns an empty result set with the configured metadata.
func (c *mockColumnValueConn) QueryContext(_ context.Context, _ string, _ []driver.NamedValue) (driver.Rows, error) {
	return &mockColumnValueRows{columnType: c.columnType}, nil
}

// mockColumnValueRows exposes one column without requiring any result values.
type mockColumnValueRows struct {
	columnType *mockColumnValueConnector
}

// Columns returns the name of the test column.
func (r *mockColumnValueRows) Columns() []string { return []string{"v"} }

// Close releases the rows, which hold no external resources.
func (r *mockColumnValueRows) Close() error { return nil }

// Next ends the result set because only its metadata is needed.
func (r *mockColumnValueRows) Next(_ []driver.Value) error { return io.EOF }

// ColumnTypeScanType reports the scan type selected by the test case.
func (r *mockColumnValueRows) ColumnTypeScanType(_ int) reflect.Type {
	return r.columnType.scanType
}

// ColumnTypeDatabaseTypeName reports the database type selected by the test case.
func (r *mockColumnValueRows) ColumnTypeDatabaseTypeName(_ int) string {
	return r.columnType.databaseType
}

// mockColumnValueDB records delegation to the driver's local value conversion.
type mockColumnValueDB struct {
	DB
	called       bool
	databaseType string
	err          error
}

// ConvertValueForLocal preserves the value or returns the configured driver error.
func (d *mockColumnValueDB) ConvertValueForLocal(_ context.Context, databaseType string, value any) (any, error) {
	d.called = true
	d.databaseType = databaseType
	if d.err != nil {
		return nil, d.err
	}
	return value, nil
}

// Test_Core_ColumnValueToLocalValue checks inferred types and declared numeric types.
func Test_Core_ColumnValueToLocalValue(t *testing.T) {
	var (
		conversionError = errors.New("driver conversion failed")
		cases           = []struct {
			name         string
			scanValue    any
			databaseType string
			value        any
			want         any
			fallback     bool
			err          error
		}{
			{name: "declared smallint", scanValue: int16(0), databaseType: "SMALLINT", value: []byte("42"), want: int16(42)},
			{name: "declared int4", scanValue: int32(0), databaseType: "INT4", value: int64(42), want: int32(42)},
			{name: "declared year", scanValue: int16(0), databaseType: "YEAR", value: []byte("2026"), want: int16(2026)},
			{name: "declared float", scanValue: float32(0), databaseType: "FLOAT", value: []byte("1.5"), want: float32(1.5)},
			{name: "integer to real", scanValue: int64(0), value: float64(1.5), want: float64(1.5), fallback: true},
			{name: "real to integer", scanValue: float64(0), value: int64(42), want: int64(42), fallback: true},
			{name: "matching integer", scanValue: int64(0), value: int64(42), want: int64(42)},
			{name: "matching real", scanValue: float64(0), value: float64(1.5), want: float64(1.5)},
			{name: "integer to text", scanValue: int64(0), value: "hello", want: "hello", fallback: true},
			{name: "integer to blob", scanValue: int64(0), value: []byte{1, 2, 3}, want: []byte{1, 2, 3}, fallback: true},
			{name: "driver error", scanValue: int64(0), value: "hello", fallback: true, err: conversionError},
		}
	)
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gtest.C(t, func(t *gtest.T) {
				var (
					ctx = context.Background()
					db  = sql.OpenDB(&mockColumnValueConnector{
						scanType:     reflect.TypeOf(tc.scanValue),
						databaseType: tc.databaseType,
					})
					localDB = &mockColumnValueDB{err: tc.err}
					core    = &Core{db: localDB}
				)
				t.Cleanup(func() { t.AssertNil(db.Close()) })
				rows, err := db.QueryContext(ctx, "select")
				t.AssertNil(err)
				t.Cleanup(func() { t.AssertNil(rows.Close()) })
				columnTypes, err := rows.ColumnTypes()
				t.AssertNil(err)
				t.AssertEQ(len(columnTypes), 1)
				value, err := core.columnValueToLocalValue(ctx, tc.value, columnTypes[0])
				t.AssertEQ(errors.Is(err, tc.err), true)
				t.AssertEQ(value, tc.want)
				t.AssertEQ(localDB.called, tc.fallback)
				if tc.fallback {
					t.AssertEQ(localDB.databaseType, tc.databaseType)
				}
			})
		})
	}
}
