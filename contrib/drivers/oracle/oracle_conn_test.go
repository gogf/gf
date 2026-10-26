// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file tests the Oracle connection wrapper that detaches transaction cancellation.

package oracle

import (
	"context"
	"database/sql/driver"
	"errors"
	"testing"

	gora "github.com/sijms/go-ora/v2"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/test/gtest"
)

// contextValueKey is the context key used to check that values survive detachment.
type contextValueKey struct{}

// Compile-time checks that the wrapper still exposes the driver interfaces go-ora implements.
var (
	_ driver.Conn               = (*detachCancelConn)(nil)
	_ driver.ConnBeginTx        = (*detachCancelConn)(nil)
	_ driver.ExecerContext      = (*detachCancelConn)(nil)
	_ driver.QueryerContext     = (*detachCancelConn)(nil)
	_ driver.ConnPrepareContext = (*detachCancelConn)(nil)
	_ driver.Pinger             = (*detachCancelConn)(nil)
	_ driver.SessionResetter    = (*detachCancelConn)(nil)
	_ driver.NamedValueChecker  = (*detachCancelConn)(nil)
	_ driver.Connector          = detachCancelConnector{}
)

// stubConnector is a driver.Connector used to exercise the wrapper without a database.
type stubConnector struct {
	conn driver.Conn
	err  error
}

// Connect returns the stub connection or error.
func (c stubConnector) Connect(context.Context) (driver.Conn, error) {
	if c.err != nil {
		return nil, c.err
	}
	return c.conn, nil
}

// Driver returns no driver. The wrapper tests do not use it.
func (stubConnector) Driver() driver.Driver {
	return nil
}

// stubConn is a driver.Conn that is not a go-ora connection.
type stubConn struct{}

// Prepare satisfies driver.Conn.
func (stubConn) Prepare(string) (driver.Stmt, error) { return nil, nil }

// Close satisfies driver.Conn.
func (stubConn) Close() error { return nil }

// Begin satisfies driver.Conn.
func (stubConn) Begin() (driver.Tx, error) { return nil, nil }

// TestDetachTransactionContextKeepsValuesAndDropsCancel verifies cancellation stops
// at the transaction context while values remain visible.
func TestDetachTransactionContextKeepsValuesAndDropsCancel(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		parent, cancel := context.WithCancel(context.WithValue(context.Background(), contextValueKey{}, "kept"))
		defer cancel()

		detached := detachTransactionContext(parent)
		cancel()

		t.AssertNE(parent.Err(), nil)
		t.AssertNil(detached.Err())
		t.Assert(detached.Value(contextValueKey{}), "kept")
	})
}

// TestDetachTransactionContextNilUsesBackground verifies a nil context does not panic.
func TestDetachTransactionContextNilUsesBackground(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		detached := detachTransactionContext(nil)
		t.AssertNE(detached, nil)
		t.AssertNil(detached.Err())
	})
}

// TestDetachCancelConnectorPassesThroughForeignConn verifies non-go-ora connections stay unchanged.
func TestDetachCancelConnectorPassesThroughForeignConn(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		raw := stubConn{}
		conn, err := (detachCancelConnector{Connector: stubConnector{conn: raw}}).Connect(context.Background())
		t.AssertNil(err)
		_, wrapped := conn.(*detachCancelConn)
		t.Assert(wrapped, false)
		t.Assert(conn, raw)
	})
}

// TestDetachCancelConnectorReturnsConnectError verifies a connector error is returned as-is.
func TestDetachCancelConnectorReturnsConnectError(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		want := errors.New("connect failed")
		conn, err := (detachCancelConnector{Connector: stubConnector{err: want}}).Connect(context.Background())
		t.Assert(conn, nil)
		t.Assert(err, want)
	})
}

// TestDetachCancelConnectorWrapsGoOraConnection verifies a go-ora connection is wrapped and
// BeginTx still delegates to the underlying driver.
func TestDetachCancelConnectorWrapsGoOraConnection(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		raw, err := gora.NewConnection("oracle://user:pass@127.0.0.1:1521/service", nil)
		t.AssertNil(err)

		conn, err := (detachCancelConnector{Connector: stubConnector{conn: raw}}).Connect(context.Background())
		t.AssertNil(err)
		wrapped, ok := conn.(*detachCancelConn)
		t.Assert(ok, true)
		t.Assert(wrapped.Connection, raw)

		parent, cancel := context.WithCancel(context.Background())
		defer cancel()
		tx, err := wrapped.BeginTx(parent, driver.TxOptions{ReadOnly: true})
		t.Assert(tx, nil)
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "readonly transaction is not supported")
	})
}

// TestOpenUsesDefaultOracleDriver verifies Open returns a database handle backed by go-ora
// without dialing the server.
func TestOpenUsesDefaultOracleDriver(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		db, err := new(Driver).Open(&gdb.ConfigNode{
			Host: "127.0.0.1",
			Port: "1521",
			Name: "service",
			User: "user",
			Pass: "pass",
		})
		t.AssertNil(err)
		t.AssertNE(db, nil)
		defer func() {
			t.AssertNil(db.Close())
		}()
		_, ok := db.Driver().(*gora.OracleDriver)
		t.Assert(ok, true)
	})
}
