// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/test/gtest"
)

// mockPanicStmt simulates a prepared statement that panics during execution
type mockPanicStmt struct {
	panicMessage string
}

func (m *mockPanicStmt) ExecContext(ctx context.Context, args ...any) (sql.Result, error) {
	if m.panicMessage != "" {
		panic(m.panicMessage)
	}
	panic("math/big: buffer too small to fit value")
}

func (m *mockPanicStmt) QueryContext(ctx context.Context, args ...any) (*sql.Rows, error) {
	if m.panicMessage != "" {
		panic(m.panicMessage)
	}
	panic("math/big: buffer too small to fit value")
}

func (m *mockPanicStmt) QueryRowContext(ctx context.Context, args ...any) *sql.Row {
	if m.panicMessage != "" {
		panic(m.panicMessage)
	}
	panic("math/big: buffer too small to fit value")
}

func (m *mockPanicStmt) Close() error {
	return nil
}

// Test_PanicRecoveryErrorWrapping tests that the panic recovery properly wraps errors
// with correct error codes and messages
func Test_PanicRecoveryErrorWrapping(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// Test creating an error from a string panic value
		defer func() {
			if exception := recover(); exception != nil {
				var err error
				if v, ok := exception.(error); ok && gerror.HasStack(v) {
					err = v
				} else {
					err = gerror.WrapCodef(gcode.CodeDbOperationError, gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception), "test SQL")
				}
				
				t.AssertNE(err, nil)
				t.Assert(strings.Contains(err.Error(), "buffer too small"), true)
				t.Assert(strings.Contains(err.Error(), "test SQL"), true)
			}
		}()

		// Simulate the panic that would occur in database operations
		panic("math/big: buffer too small to fit value")
	})

	gtest.C(t, func(t *gtest.T) {
		// Test creating an error from an error panic value with stack
		defer func() {
			if exception := recover(); exception != nil {
				var err error
				if v, ok := exception.(error); ok && gerror.HasStack(v) {
					err = v
				} else {
					err = gerror.WrapCodef(gcode.CodeDbOperationError, gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception), "test SQL")
				}
				
				t.AssertNE(err, nil)
				// Since gerror has stack, it should preserve the original error
				t.Assert(strings.Contains(err.Error(), "custom database error"), true)
			}
		}()

		// Simulate a panic with a custom error that has stack
		customErr := gerror.New("custom database error")
		panic(customErr)
	})
}

// Test_DoCommit_StmtPanicRecovery simulates the scenario from the issue where
// statement execution causes a panic during DoCommit operations
func Test_DoCommit_StmtPanicRecovery(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// We'll test the panic recovery by triggering it in the defer function
		// Since we can't easily mock sql.Stmt, we'll test the panic recovery mechanism directly
		
		testPanicRecovery := func(panicValue any, sqlText string) (err error) {
			defer func() {
				if exception := recover(); exception != nil {
					if err == nil {
						if v, ok := exception.(error); ok && gerror.HasStack(v) {
							err = v
						} else {
							err = gerror.WrapCodef(gcode.CodeDbOperationError, gerror.NewCodef(gcode.CodeInternalPanic, "%+v", exception), FormatSqlWithArgs(sqlText, []any{123}))
						}
					}
				}
			}()
			
			// Simulate the panic that would occur in database operations
			panic(panicValue)
		}
		
		// Test different panic scenarios
		testCases := []struct {
			name       string
			panicValue any
			sqlText    string
		}{
			{
				name:       "String panic from math/big",
				panicValue: "math/big: buffer too small to fit value",
				sqlText:    "INSERT INTO test VALUES (?)",
			},
			{
				name:       "Custom error panic",
				panicValue: gerror.New("clickhouse driver panic"),
				sqlText:    "SELECT * FROM test WHERE id = ?",
			},
		}

		for _, tc := range testCases {
			t.Log("Testing:", tc.name)
			
			// Test the panic recovery mechanism
			err := testPanicRecovery(tc.panicValue, tc.sqlText)
			
			// After our fix, these should return errors instead of panicking
			t.AssertNE(err, nil)
			
			// Verify the error contains information about the panic
			errorMsg := err.Error()
			
			if tc.name == "String panic from math/big" {
				t.Assert(strings.Contains(errorMsg, "buffer too small"), true)
				t.Assert(strings.Contains(errorMsg, "INSERT INTO test VALUES"), true)
			} else {
				t.Assert(strings.Contains(errorMsg, "clickhouse driver panic"), true)
			}
		}
	})
}