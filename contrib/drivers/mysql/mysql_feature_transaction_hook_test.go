package mysql_test

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"testing"
)

func Test_Transaction_Hook_For_Begin(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		// test begin hook fail
		tx, err := db.Begin(ctx, gdb.TxHookHandler{
			Begin: func(ctx context.Context, in *gdb.HookBeginInput) (err error) {
				fmt.Println("First Begin() exec begin hook return err:,txId:" + in.TransactionId)
				return gerror.New("begin hook fail")
			},
			Commit: func(ctx context.Context, in *gdb.HookCommitInput) (err error) {
				return
			},
			Rollback: func(ctx context.Context, in *gdb.HookRollbackInput) (err error) {
				return
			},
		})
		t.Assert(err.Error(), "begin hook fail")

		// test begin hook success but commit hook fail
		tx, err = db.Begin(ctx, gdb.TxHookHandler{
			Begin: func(ctx context.Context, in *gdb.HookBeginInput) (err error) {
				fmt.Println("Second Begin() exec begin hook return nil,txId:" + in.TransactionId)
				return
			},
			Commit: func(ctx context.Context, in *gdb.HookCommitInput) (err error) {
				fmt.Println("Second Begin() exec commit hook return err,txId:" + in.TransactionId)
				return gerror.New("commit hook fail")
			},
			Rollback: func(ctx context.Context, in *gdb.HookRollbackInput) (err error) {
				return
			},
		})
		t.Assert(err, nil)
		if err = g.TryFunc(ctx, func(ctx context.Context) error {
			_, err = tx.Model(table).Data(g.Map{
				"passport": "user_tx2",
				"password": "pass_tx2",
				"nickname": "name_tx2",
			}).InsertAndGetId()
			if err != nil {
				return err
			}
			return nil
		}); err != nil {
			err = tx.Rollback()
		} else {
			err = tx.Commit()
		}
		t.Assert(err.Error(), "commit hook fail")
		// test begin hook success but rollback hook fail
		tx, err = db.Begin(ctx, gdb.TxHookHandler{
			Begin: func(ctx context.Context, in *gdb.HookBeginInput) (err error) {
				fmt.Println("Three Begin() exec begin hook return nil,txId:" + in.TransactionId)
				return
			},
			Commit: func(ctx context.Context, in *gdb.HookCommitInput) (err error) {
				return
			},
			Rollback: func(ctx context.Context, in *gdb.HookRollbackInput) (err error) {
				fmt.Println("Three Begin() exec rollback hook return err,txId:" + in.TransactionId)
				return gerror.New("rollback hook fail")
			},
		})
		t.Assert(err, nil)
		err = tx.Rollback()
		t.Assert(err.Error(), "rollback hook fail")

		// test all hook success
		tx, err = db.Begin(ctx, gdb.TxHookHandler{
			Begin: func(ctx context.Context, in *gdb.HookBeginInput) (err error) {
				fmt.Println("Four Begin() exec begin hook return nil,txId:" + in.TransactionId)
				return
			},
			Commit: func(ctx context.Context, in *gdb.HookCommitInput) (err error) {
				fmt.Println("Four Begin() exec commit hook return nil,txId:" + in.TransactionId)
				return
			},
			Rollback: func(ctx context.Context, in *gdb.HookRollbackInput) (err error) {
				return gerror.New("rollback hook fail")
			},
		})
		t.Assert(err, nil)
		var userId int64
		var count int
		if err = g.TryFunc(ctx, func(ctx context.Context) error {
			userId, err = tx.Model(table).Data(g.Map{
				"passport": "user_tx2",
				"password": "pass_tx2",
				"nickname": "name_tx2",
			}).InsertAndGetId()
			if err != nil {
				return err
			}
			return nil
		}); err != nil {
			err = tx.Rollback()
		} else {
			err = tx.Commit()
		}
		count, err = db.Model(table).Where("id", userId).Count()
		t.Assert(err, nil)
		t.Assert(count, 1)
	})
}

func Test_Transaction_Hook_For_Transaction(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		// test begin hook fail
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			return nil
		}, gdb.TxHookHandler{
			Begin: func(ctx context.Context, in *gdb.HookBeginInput) (err error) {
				fmt.Println("First Transaction() exec begin hook return err,txId:" + in.TransactionId)
				return gerror.New("begin hook fail")
			},
			Commit: func(ctx context.Context, in *gdb.HookCommitInput) (err error) {
				return
			},
			Rollback: func(ctx context.Context, in *gdb.HookRollbackInput) (err error) {
				return
			},
		})
		t.Assert(err.Error(), "begin hook fail")
		// test commit hook fail
		err = db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			return nil
		}, gdb.TxHookHandler{
			Begin: func(ctx context.Context, in *gdb.HookBeginInput) (err error) {
				fmt.Println("Second Transaction() exec begin hook return nil,txId:" + in.TransactionId)
				return
			},
			Commit: func(ctx context.Context, in *gdb.HookCommitInput) (err error) {
				fmt.Println("Second Transaction() exec commit hook return err,txId:" + in.TransactionId)
				return gerror.New("after commit hook fail")
			},
			Rollback: func(ctx context.Context, in *gdb.HookRollbackInput) (err error) {
				return gerror.New("after rollback hook")
			},
		})
		t.Assert(err.Error(), "after commit hook fail")
		// test rollback hook fail
		err = db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			return gerror.New("Rollback")
		}, gdb.TxHookHandler{
			Begin: func(ctx context.Context, in *gdb.HookBeginInput) (err error) {
				fmt.Println("Three Transaction() exec begin hook return nil,txId:" + in.TransactionId)
				return
			},
			Commit: func(ctx context.Context, in *gdb.HookCommitInput) (err error) {
				return gerror.New("after commit hook")
			},
			Rollback: func(ctx context.Context, in *gdb.HookRollbackInput) (err error) {
				fmt.Println("Three Transaction() exec rollback hook return err,txId:" + in.TransactionId)
				return gerror.New("after rollback hook fail")
			},
		})
		t.Assert(err.Error(), "after rollback hook fail")
		// test all success
		var userId int64
		var count int
		err = db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			userId, err = tx.Model(table).Data(g.Map{
				"passport": "user_tx",
				"password": "pass_tx",
				"nickname": "name_tx",
			}).InsertAndGetId()
			return err
		}, gdb.TxHookHandler{
			Begin: func(ctx context.Context, in *gdb.HookBeginInput) (err error) {
				fmt.Println("Four Transaction() exec begin hook return nil,txId:" + in.TransactionId)
				return
			},
			Commit: func(ctx context.Context, in *gdb.HookCommitInput) (err error) {
				count, err = db.Model(table).Where("id", userId).Count()
				t.Assert(err, nil)
				t.Assert(count, 1)
				fmt.Println("Four Transaction() exec commit hook return nil,txId:" + in.TransactionId)
				return nil
			},
			Rollback: func(ctx context.Context, in *gdb.HookRollbackInput) (err error) {
				return nil
			},
		})
		t.Assert(err, nil)
	})
}

func Test_Transaction_Hook_For_Transaction_No_Begin(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		// test commit hook fail
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			tx.Hook(gdb.TxHookHandler{
				Commit: func(ctx context.Context, in *gdb.HookCommitInput) (err error) {
					fmt.Println("First tx Transaction() No Begin exec commit hook return err,txId:" + in.TransactionId)
					return gerror.New("after commit hook fail")
				},
				Rollback: func(ctx context.Context, in *gdb.HookRollbackInput) (err error) {
					return gerror.New("after rollback hook")
				},
			})
			return nil
		})
		t.Assert(err.Error(), "after commit hook fail")
		// test rollback hook fail
		err = db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			tx.Hook(gdb.TxHookHandler{
				Commit: func(ctx context.Context, in *gdb.HookCommitInput) (err error) {
					return gerror.New("after commit hook")
				},
				Rollback: func(ctx context.Context, in *gdb.HookRollbackInput) (err error) {
					fmt.Println("Second tx Transaction() No Begin exec rollback hook return err,txId:" + in.TransactionId)
					return gerror.New("after rollback hook fail")
				},
			})
			return gerror.New("Rollback")
		})
		t.Assert(err.Error(), "after rollback hook fail")
		// test all success
		err = db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			var userId int64
			var count int
			tx.Hook(gdb.TxHookHandler{
				Commit: func(ctx context.Context, in *gdb.HookCommitInput) (err error) {
					count, err = db.Model(table).Where("id", userId).Count()
					t.Assert(err, nil)
					t.Assert(count, 1)
					fmt.Println("Three tx Transaction() No Begin exec commit hook return nil,txId:" + in.TransactionId)
					return nil
				},
				Rollback: func(ctx context.Context, in *gdb.HookRollbackInput) (err error) {
					return nil
				},
			})
			userId, err = tx.Model(table).Data(g.Map{
				"passport": "user_tx",
				"password": "pass_tx",
				"nickname": "name_tx",
			}).InsertAndGetId()
			return err
		})
		t.Assert(err, nil)
	})
}
