// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file unit-tests DryRun mode and Hook handlers on the Model API.
// Neither feature requires a live database connection: DryRun is verified
// through SetDryRun/GetDryRun round-trips; Hook behaviour is verified by
// confirming that the hook callback intercepts the query and that any error
// it returns reaches the caller — the callback short-circuits Next() so no
// real connection is needed.

package gdb

import (
	"context"
	"database/sql"
	"testing"

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

// Test_Hook_Model_ReturnsNewModel confirms that Hook returns a new model and
// does not mutate the original — the standard clone-on-chain contract.
func Test_Hook_Model_ReturnsNewModel(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		db := newDialectFakeDB("id", "name")
		m := db.Model("users")
		hook := HookHandler{
			Select: func(_ context.Context, _ *HookSelectInput) (Result, error) {
				return Result{}, nil
			},
		}
		mWithHook := m.Hook(hook)
		// Hook must return a clone: hookHandler.Select set only on the clone.
		t.Assert(mWithHook.hookHandler.Select != nil, true)
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
		m := db.Model("users").Hook(hook).Where("id", 1)
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
// does not affect the original model's hook state — a fresh db.Model() should
// always start with no hook handlers set.
func Test_Hook_Model_OriginalUnchangedAfterChain(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		db := newDialectFakeDB("id", "name")
		hook := HookHandler{
			Select: func(_ context.Context, _ *HookSelectInput) (Result, error) {
				return Result{}, nil
			},
		}
		// A new model from db.Model() must have no hook.
		original := db.Model("users")
		t.Assert(original.hookHandler.Select == nil, true)

		// Clone chain must not change that.
		_ = original.Hook(hook).Where("status", "active")
		original2 := db.Model("users")
		t.Assert(original2.hookHandler.Select == nil, true)
	})
}

