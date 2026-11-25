// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file unit-tests DryRun mode and Hook handlers on the Model API.
// Neither feature requires a live database connection: DryRun is verified
// through SetDryRun/GetDryRun round-trips; for Hook, the select handler is
// installed on a model and the model is queried directly — the handler
// short-circuits Next() before any statement is committed, so the error it
// returns is what the caller receives.

package gdb

import (
	"context"
	"database/sql"
	"testing"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/test/gtest"
)

// ─────────────────────────── DryRun ───────────────────────────────

// Test_DryRun_SetGet verifies that SetDryRun and GetDryRun are consistent.
func Test_DryRun_SetGet(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		db := newDialectFakeDB()
		t.Assert(db.GetDryRun(), false)

		db.SetDryRun(true)
		t.Assert(db.GetDryRun(), true)

		db.SetDryRun(false)
		t.Assert(db.GetDryRun(), false)
	})
}

// Test_DryRun_Toggle verifies repeated enable/disable cycles.
func Test_DryRun_Toggle(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		db := newDialectFakeDB()
		for i := 0; i < 5; i++ {
			db.SetDryRun(true)
			t.Assert(db.GetDryRun(), true)
			db.SetDryRun(false)
			t.Assert(db.GetDryRun(), false)
		}
	})
}

// ─────────────────────────── Hook ─────────────────────────────────

// Test_Hook_Model_UnsafeSetsHandlerInPlace confirms that Hook on the default
// (unsafe) model sets the handler on the receiver and returns the receiver.
func Test_Hook_Model_UnsafeSetsHandlerInPlace(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		db := newDialectFakeDB("id", "name")
		m := db.Model("users")
		hook := HookHandler{
			Select: func(_ context.Context, _ *HookSelectInput) (Result, error) {
				return Result{}, nil
			},
		}
		mWithHook := m.Hook(hook)
		t.Assert(mWithHook == m, true)
		t.Assert(m.hookHandler.Select != nil, true)
	})
}

// Test_Hook_Model_SafeReturnsClone confirms that Hook on a safe model returns a
// clone carrying the handler and leaves the original model without one.
func Test_Hook_Model_SafeReturnsClone(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		db := newDialectFakeDB("id", "name")
		original := db.Model("users").Safe()
		hook := HookHandler{
			Select: func(_ context.Context, _ *HookSelectInput) (Result, error) {
				return Result{}, nil
			},
		}
		mWithHook := original.Hook(hook)
		t.Assert(mWithHook == original, false)
		t.Assert(mWithHook.hookHandler.Select != nil, true)
		t.Assert(original.hookHandler.Select == nil, true)
	})
}

// Test_Hook_Model_MultipleHandlers verifies that only specified hooks are set;
// unspecified ones remain nil on the clone.
func Test_Hook_Model_MultipleHandlers(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		db := newDialectFakeDB("id", "name")
		hook := HookHandler{
			Insert: func(_ context.Context, _ *HookInsertInput) (sql.Result, error) {
				return new(SqlResult), nil
			},
		}
		m := db.Model("users").Hook(hook)
		t.Assert(m.hookHandler.Select == nil, true)
		t.Assert(m.hookHandler.Insert != nil, true)
		t.Assert(m.hookHandler.Update == nil, true)
		t.Assert(m.hookHandler.Delete == nil, true)
	})
}

// Test_Hook_Model_ChainPreservesHandler confirms that chaining further model
// methods after Hook preserves the hook handler on the new clone.
func Test_Hook_Model_ChainPreservesHandler(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		db := newDialectFakeDB("id", "name")
		hook := HookHandler{
			Select: func(_ context.Context, _ *HookSelectInput) (Result, error) {
				return Result{}, nil
			},
		}
		// .Where() clones again after Hook — handler must survive.
		m := db.Model("users").Safe().Hook(hook).Where("id", 1)
		t.Assert(m.hookHandler.Select != nil, true)
	})
}

// Test_Hook_AllHandlerTypes verifies that all four hook types (Select, Insert,
// Update, Delete) can be assigned at once.
func Test_Hook_AllHandlerTypes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		db := newDialectFakeDB("id", "name")
		hook := HookHandler{
			Select: func(_ context.Context, _ *HookSelectInput) (Result, error) {
				return Result{}, nil
			},
			Insert: func(_ context.Context, _ *HookInsertInput) (sql.Result, error) {
				return new(SqlResult), nil
			},
			Update: func(_ context.Context, _ *HookUpdateInput) (sql.Result, error) {
				return new(SqlResult), nil
			},
			Delete: func(_ context.Context, _ *HookDeleteInput) (sql.Result, error) {
				return new(SqlResult), nil
			},
		}
		m := db.Model("users").Hook(hook)
		t.Assert(m.hookHandler.Select != nil, true)
		t.Assert(m.hookHandler.Insert != nil, true)
		t.Assert(m.hookHandler.Update != nil, true)
		t.Assert(m.hookHandler.Delete != nil, true)
	})
}

// Test_Hook_Model_OriginalUnchangedAfterChain verifies that chaining Hook+Where
// on a safe model keeps the handler on the derived model and leaves the original
// model's hook state untouched.
func Test_Hook_Model_OriginalUnchangedAfterChain(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		db := newDialectFakeDB("id", "name")
		hook := HookHandler{
			Select: func(_ context.Context, _ *HookSelectInput) (Result, error) {
				return Result{}, nil
			},
		}
		original := db.Model("users").Safe()
		derived := original.Hook(hook).Where("status", "active")
		t.Assert(derived.hookHandler.Select != nil, true)
		t.Assert(original.hookHandler.Select == nil, true)
	})
}

// Test_Hook_Model_SelectHandlerErrorReachesCaller verifies that querying a model
// with a select handler calls the handler with the statement and returns the
// handler's error to the caller.
func Test_Hook_Model_SelectHandlerErrorReachesCaller(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		db := newDialectFakeDB("id", "name")
		db.links = gmap.NewKVMap[ConfigNode, *sql.DB]()
		var (
			hookedSql   string
			hookedTable string
		)
		m := db.Model("users").Hook(HookHandler{
			Select: func(_ context.Context, in *HookSelectInput) (Result, error) {
				hookedSql = in.Sql
				hookedTable = in.Table
				return nil, gerror.New("hook intercepted the query")
			},
		})
		result, err := m.Where("id", 1).All()
		t.Assert(len(result), 0)
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "hook intercepted the query")
		t.Assert(hookedTable != "", true)
		t.Assert(hookedSql != "", true)
	})
}
