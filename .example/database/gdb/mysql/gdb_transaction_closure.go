package main

import (
	"context"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func main() {
	var (
		err   error
		db    = g.DB()
		ctx   = gctx.New()
		table = "user"
	)
	if err = db.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
		// Nested transaction 1.
		if err = tx.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
			_, err = tx.Model(table).Data(g.Map{"id": 1, "name": "john"}).Insert()
			return err
		}); err != nil {
			return err
		}
		// Nested transaction 2, panic.
		if err = tx.Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
			_, err = tx.Model(table).Data(g.Map{"id": 2, "name": "smith"}).Insert()
			// Create a panic that can make this transaction rollback automatically.
			panic("error")
		}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		panic(err)
	}
}
