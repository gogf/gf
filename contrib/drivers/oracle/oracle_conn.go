// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file keeps a go-ora transaction off the caller context.
// go-ora v2.9.0 stores the BeginTx context and, on Rollback, starts a watchdog from it.
// When that context is already canceled, the watchdog writes Session.breakConn while
// the rollback read calls IsBreak without the session mutex. database/sql still watches
// the caller context and rolls the transaction back.

package oracle

import (
	"context"
	"database/sql/driver"

	gora "github.com/sijms/go-ora/v2"
)

// detachCancelConnector opens Oracle connections whose driver transactions ignore
// cancellation of the caller context.
type detachCancelConnector struct {
	driver.Connector
}

// Connect opens a connection. go-ora connections are wrapped so BeginTx does not
// observe cancellation of the caller context.
func (c detachCancelConnector) Connect(ctx context.Context) (driver.Conn, error) {
	conn, err := c.Connector.Connect(ctx)
	if err != nil {
		return nil, err
	}
	oraConn, ok := conn.(*gora.Connection)
	if !ok {
		return conn, nil
	}
	return &detachCancelConn{Connection: oraConn}, nil
}

// detachCancelConn is a go-ora connection that detaches transaction cancellation.
// Embedding *gora.Connection keeps the other driver.Conn interfaces promoted.
type detachCancelConn struct {
	*gora.Connection
}

// BeginTx starts a transaction with a context that does not cancel with the caller.
// Query and exec methods still receive the caller context through the promoted methods.
func (c *detachCancelConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	return c.Connection.BeginTx(detachTransactionContext(ctx), opts)
}

// detachTransactionContext returns a context that keeps ctx values and drops its
// cancellation. A nil context becomes context.Background because WithoutCancel
// rejects nil.
func detachTransactionContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return context.WithoutCancel(ctx)
}
