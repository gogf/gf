// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
)

// Begin starts and returns the transaction object.
func (d *Driver) Begin(ctx context.Context) (tx gdb.TX, err error) {
	return nil, errUnsupportedBegin
}

// BeginWithOptions starts and returns the transaction object with given options.
func (d *Driver) BeginWithOptions(ctx context.Context, opts gdb.TxOptions) (tx gdb.TX, err error) {
	return nil, errUnsupportedBegin
}

// Transaction wraps the transaction logic using function `f`.
func (d *Driver) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) error {
	return errUnsupportedTransaction
}

// TransactionWithOptions wraps the transaction logic with propagation options using function `f`.
func (d *Driver) TransactionWithOptions(
	ctx context.Context, opts gdb.TxOptions, f func(ctx context.Context, tx gdb.TX) error,
) error {
	return errUnsupportedTransaction
}
