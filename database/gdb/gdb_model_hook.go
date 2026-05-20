// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gogf/gf/v3/container/gvar"
	"github.com/gogf/gf/v3/text/gregex"
	"github.com/gogf/gf/v3/text/gstr"
)

type (
	HookBeforeSelect func(ctx context.Context, in *HookSelectInput) error
	HookAfterSelect  func(ctx context.Context, in *HookSelectInput, result Result, err error) (Result, error)

	HookBeforeInsert func(ctx context.Context, in *HookInsertInput) error
	HookAfterInsert  func(ctx context.Context, in *HookInsertInput, result sql.Result, err error) (sql.Result, error)

	HookBeforeUpdate func(ctx context.Context, in *HookUpdateInput) error
	HookAfterUpdate  func(ctx context.Context, in *HookUpdateInput, result sql.Result, err error) (sql.Result, error)

	HookBeforeDelete func(ctx context.Context, in *HookDeleteInput) error
	HookAfterDelete  func(ctx context.Context, in *HookDeleteInput, result sql.Result, err error) (sql.Result, error)
)

// HookHandler manages all supported hook functions for Model.
type HookHandler struct {
	selectBefore []HookBeforeSelect
	selectAfter  []HookAfterSelect

	insertBefore []HookBeforeInsert
	insertAfter  []HookAfterInsert

	updateBefore []HookBeforeUpdate
	updateAfter  []HookAfterUpdate

	deleteBefore []HookBeforeDelete
	deleteAfter  []HookAfterDelete
}

type HookType int

const (
	HookTypeSelect HookType = 1
	HookTypeInsert HookType = 2
	HookTypeUpdate HookType = 3
	HookTypeDelete HookType = 4
)

type HookStage int

const (
	HookStageBefore HookStage = 1
	HookStageAfter  HookStage = 2
)

type HookDescriptor struct {
	Type    HookType
	Stage   HookStage
	Handler any
}

func BeforeSelect(handler HookBeforeSelect) HookDescriptor {
	return HookDescriptor{Type: HookTypeSelect, Stage: HookStageBefore, Handler: handler}
}

func AfterSelect(handler HookAfterSelect) HookDescriptor {
	return HookDescriptor{Type: HookTypeSelect, Stage: HookStageAfter, Handler: handler}
}

func BeforeInsert(handler HookBeforeInsert) HookDescriptor {
	return HookDescriptor{Type: HookTypeInsert, Stage: HookStageBefore, Handler: handler}
}

func AfterInsert(handler HookAfterInsert) HookDescriptor {
	return HookDescriptor{Type: HookTypeInsert, Stage: HookStageAfter, Handler: handler}
}

func BeforeUpdate(handler HookBeforeUpdate) HookDescriptor {
	return HookDescriptor{Type: HookTypeUpdate, Stage: HookStageBefore, Handler: handler}
}

func AfterUpdate(handler HookAfterUpdate) HookDescriptor {
	return HookDescriptor{Type: HookTypeUpdate, Stage: HookStageAfter, Handler: handler}
}

func BeforeDelete(handler HookBeforeDelete) HookDescriptor {
	return HookDescriptor{Type: HookTypeDelete, Stage: HookStageBefore, Handler: handler}
}

func AfterDelete(handler HookAfterDelete) HookDescriptor {
	return HookDescriptor{Type: HookTypeDelete, Stage: HookStageAfter, Handler: handler}
}

// internalParamHook manages all internal parameters for hook operations.
// The `internal` obviously means you cannot access these parameters outside this package.
type internalParamHook struct {
	link               Link      // Connection object from third party sql driver.
	removedWhere       bool      // Removed mark for condition string that was removed `WHERE` prefix.
	originalTableName  *gvar.Var // The original table name.
	originalSchemaName *gvar.Var // The original schema name.
}

// HookSelectInput holds the parameters for select hook operation.
// Note that, COUNT statement will also be hooked by this feature,
// which is usually not be interesting for upper business hook handler.
type HookSelectInput struct {
	internalParamHook
	Model      *Model     // Current operation Model.
	Table      string     // The table name that to be used. Update this attribute to change target table name.
	Schema     string     // The schema name that to be used. Update this attribute to change target schema name.
	Sql        string     // The sql string that to be committed.
	Args       []any      // The arguments of sql.
	SelectType SelectType // The type of this SELECT operation.
}

// HookInsertInput holds the parameters for insert hook operation.
type HookInsertInput struct {
	internalParamHook
	Model  *Model         // Current operation Model.
	Table  string         // The table name that to be used. Update this attribute to change target table name.
	Schema string         // The schema name that to be used. Update this attribute to change target schema name.
	Data   List           // The data records list to be inserted/saved into table.
	Option DoInsertOption // The extra option for data inserting.
}

// HookUpdateInput holds the parameters for update hook operation.
type HookUpdateInput struct {
	internalParamHook
	Model     *Model // Current operation Model.
	Table     string // The table name that to be used. Update this attribute to change target table name.
	Schema    string // The schema name that to be used. Update this attribute to change target schema name.
	Data      any    // Data can be type of: map[string]any/string. You can use type assertion on `Data`.
	Condition string // The where condition string for updating.
	Args      []any  // The arguments for sql place-holders.
}

// HookDeleteInput holds the parameters for delete hook operation.
type HookDeleteInput struct {
	internalParamHook
	Model     *Model // Current operation Model.
	Table     string // The table name that to be used. Update this attribute to change target table name.
	Schema    string // The schema name that to be used. Update this attribute to change target schema name.
	Condition string // The where condition string for deleting.
	Args      []any  // The arguments for sql place-holders.
}

const (
	whereKeyInCondition = " WHERE "
)

// IsTransaction checks and returns whether current operation is during transaction.
func (h *internalParamHook) IsTransaction() bool {
	return h.link.IsTransaction()
}

// Hook sets the hook functions for current model.
// Can be used multiple times without overwriting.
func (m *Model) Hook(descriptors ...HookDescriptor) *Model {
	model := m

	for _, descriptor := range descriptors {
		model.hookHandler = model.hookHandler.append(descriptor)
	}

	return model
}

func (h HookHandler) append(descriptor HookDescriptor) HookHandler {
	switch descriptor.Type {
	case HookTypeSelect:
		switch descriptor.Stage {
		case HookStageBefore:
			handler, ok := descriptor.Handler.(HookBeforeSelect)
			if !ok {
				panic(fmt.Sprintf("invalid hook handler type for select/before: %T", descriptor.Handler))
			}
			h.selectBefore = append(h.selectBefore, handler)
		case HookStageAfter:
			handler, ok := descriptor.Handler.(HookAfterSelect)
			if !ok {
				panic(fmt.Sprintf("invalid hook handler type for select/after: %T", descriptor.Handler))
			}
			h.selectAfter = append(h.selectAfter, handler)
		default:
			panic(fmt.Sprintf("invalid hook stage: %d", descriptor.Stage))
		}

	case HookTypeInsert:
		switch descriptor.Stage {
		case HookStageBefore:
			handler, ok := descriptor.Handler.(HookBeforeInsert)
			if !ok {
				panic(fmt.Sprintf("invalid hook handler type for insert/before: %T", descriptor.Handler))
			}
			h.insertBefore = append(h.insertBefore, handler)
		case HookStageAfter:
			handler, ok := descriptor.Handler.(HookAfterInsert)
			if !ok {
				panic(fmt.Sprintf("invalid hook handler type for insert/after: %T", descriptor.Handler))
			}
			h.insertAfter = append(h.insertAfter, handler)
		default:
			panic(fmt.Sprintf("invalid hook stage: %d", descriptor.Stage))
		}

	case HookTypeUpdate:
		switch descriptor.Stage {
		case HookStageBefore:
			handler, ok := descriptor.Handler.(HookBeforeUpdate)
			if !ok {
				panic(fmt.Sprintf("invalid hook handler type for update/before: %T", descriptor.Handler))
			}
			h.updateBefore = append(h.updateBefore, handler)
		case HookStageAfter:
			handler, ok := descriptor.Handler.(HookAfterUpdate)
			if !ok {
				panic(fmt.Sprintf("invalid hook handler type for update/after: %T", descriptor.Handler))
			}
			h.updateAfter = append(h.updateAfter, handler)
		default:
			panic(fmt.Sprintf("invalid hook stage: %d", descriptor.Stage))
		}

	case HookTypeDelete:
		switch descriptor.Stage {
		case HookStageBefore:
			handler, ok := descriptor.Handler.(HookBeforeDelete)
			if !ok {
				panic(fmt.Sprintf("invalid hook handler type for delete/before: %T", descriptor.Handler))
			}
			h.deleteBefore = append(h.deleteBefore, handler)
		case HookStageAfter:
			handler, ok := descriptor.Handler.(HookAfterDelete)
			if !ok {
				panic(fmt.Sprintf("invalid hook handler type for delete/after: %T", descriptor.Handler))
			}
			h.deleteAfter = append(h.deleteAfter, handler)
		default:
			panic(fmt.Sprintf("invalid hook stage: %d", descriptor.Stage))
		}

	default:
		panic(fmt.Sprintf("invalid hook type: %d", descriptor.Type))
	}
	return h
}

func (h HookHandler) Clone() HookHandler {
	return HookHandler{
		selectBefore: append([]HookBeforeSelect(nil), h.selectBefore...),
		selectAfter:  append([]HookAfterSelect(nil), h.selectAfter...),
		insertBefore: append([]HookBeforeInsert(nil), h.insertBefore...),
		insertAfter:  append([]HookAfterInsert(nil), h.insertAfter...),
		updateBefore: append([]HookBeforeUpdate(nil), h.updateBefore...),
		updateAfter:  append([]HookAfterUpdate(nil), h.updateAfter...),
		deleteBefore: append([]HookBeforeDelete(nil), h.deleteBefore...),
		deleteAfter:  append([]HookAfterDelete(nil), h.deleteAfter...),
	}
}

func (h HookHandler) runSelect(ctx context.Context, in *HookSelectInput) (result Result, err error) {
	in.initOriginalNames()
	if err = in.applySharding(ctx); err != nil {
		return nil, err
	}
	for _, before := range h.selectBefore {
		if err = before(ctx, in); err != nil {
			return nil, err
		}
	}
	if result, err = in.doSelect(ctx); err != nil {
		return nil, err
	}
	for _, after := range h.selectAfter {
		result, err = after(ctx, in, result, err)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (h HookHandler) runInsert(ctx context.Context, in *HookInsertInput) (result sql.Result, err error) {
	in.initOriginalNames()
	if err = in.applySharding(ctx); err != nil {
		return nil, err
	}
	for _, before := range h.insertBefore {
		if err = before(ctx, in); err != nil {
			return nil, err
		}
	}
	if result, err = in.doInsert(ctx); err != nil {
		return nil, err
	}
	for _, after := range h.insertAfter {
		result, err = after(ctx, in, result, err)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (h HookHandler) runUpdate(ctx context.Context, in *HookUpdateInput) (result sql.Result, err error) {
	in.initOriginalNames()
	if err = in.applySharding(ctx); err != nil {
		return nil, err
	}
	in.normalizeWhereForHooks()
	for _, before := range h.updateBefore {
		if err = before(ctx, in); err != nil {
			return nil, err
		}
	}
	if result, err = in.doUpdate(ctx); err != nil {
		return nil, err
	}
	for _, after := range h.updateAfter {
		result, err = after(ctx, in, result, err)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (h HookHandler) runDelete(ctx context.Context, in *HookDeleteInput) (result sql.Result, err error) {
	in.initOriginalNames()
	if err = in.applySharding(ctx); err != nil {
		return nil, err
	}
	in.normalizeWhereForHooks()
	for _, before := range h.deleteBefore {
		if err = before(ctx, in); err != nil {
			return nil, err
		}
	}
	if result, err = in.doDelete(ctx); err != nil {
		return nil, err
	}
	for _, after := range h.deleteAfter {
		result, err = after(ctx, in, result, err)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (h *internalParamHook) initOriginalNames(schema, table string) {
	if h.originalTableName.IsNil() {
		h.originalTableName = gvar.New(table)
	}
	if h.originalSchemaName.IsNil() {
		h.originalSchemaName = gvar.New(schema)
	}
}

func (in *HookSelectInput) initOriginalNames() {
	in.internalParamHook.initOriginalNames(in.Schema, in.Table)
}

func (in *HookInsertInput) initOriginalNames() {
	in.internalParamHook.initOriginalNames(in.Schema, in.Table)
}

func (in *HookUpdateInput) initOriginalNames() {
	in.internalParamHook.initOriginalNames(in.Schema, in.Table)
}

func (in *HookDeleteInput) initOriginalNames() {
	in.internalParamHook.initOriginalNames(in.Schema, in.Table)
}

func (in *HookSelectInput) applySharding(ctx context.Context) (err error) {
	in.Schema, err = in.Model.getActualSchema(ctx, in.Schema)
	if err != nil {
		return err
	}
	in.Table, err = in.Model.getActualTable(ctx, in.Table)
	if err != nil {
		return err
	}
	return nil
}

func (in *HookInsertInput) applySharding(ctx context.Context) (err error) {
	in.Schema, err = in.Model.getActualSchema(ctx, in.Schema)
	if err != nil {
		return err
	}
	in.Table, err = in.Model.getActualTable(ctx, in.Table)
	if err != nil {
		return err
	}
	return nil
}

func (in *HookUpdateInput) applySharding(ctx context.Context) (err error) {
	in.Schema, err = in.Model.getActualSchema(ctx, in.Schema)
	if err != nil {
		return err
	}
	in.Table, err = in.Model.getActualTable(ctx, in.Table)
	if err != nil {
		return err
	}
	return nil
}

func (in *HookDeleteInput) applySharding(ctx context.Context) (err error) {
	in.Schema, err = in.Model.getActualSchema(ctx, in.Schema)
	if err != nil {
		return err
	}
	in.Table, err = in.Model.getActualTable(ctx, in.Table)
	if err != nil {
		return err
	}
	return nil
}

func (in *HookUpdateInput) normalizeWhereForHooks() {
	if gstr.HasPrefix(in.Condition, whereKeyInCondition) {
		in.removedWhere = true
		in.Condition = gstr.TrimLeftStr(in.Condition, whereKeyInCondition)
	}
}

func (in *HookDeleteInput) normalizeWhereForHooks() {
	if gstr.HasPrefix(in.Condition, whereKeyInCondition) {
		in.removedWhere = true
		in.Condition = gstr.TrimLeftStr(in.Condition, whereKeyInCondition)
	}
}

func (in *HookUpdateInput) restoreWhereForCommit() string {
	if in.removedWhere {
		return whereKeyInCondition + in.Condition
	}
	return in.Condition
}

func (in *HookDeleteInput) restoreWhereForCommit() string {
	if in.removedWhere {
		return whereKeyInCondition + in.Condition
	}
	return in.Condition
}

func (in *HookSelectInput) doSelect(ctx context.Context) (result Result, err error) {
	// Table change.
	if in.Table != in.originalTableName.String() {
		in.Sql, err = gregex.ReplaceStringFuncMatch(
			`(?i) FROM ([\S]+)`,
			in.Sql,
			func(match []string) string {
				charL, charR := in.Model.db.GetChars()
				return fmt.Sprintf(` FROM %s%s%s`, charL, in.Table, charR)
			},
		)
		if err != nil {
			return nil, err
		}
	}

	// Schema change.
	if in.Schema != "" && in.Schema != in.originalSchemaName.String() {
		in.link, err = in.Model.db.GetCore().SlaveLink(in.Schema)
		if err != nil {
			return nil, err
		}
		in.Model.db.GetCore().schema = in.Schema
		defer func() {
			in.Model.db.GetCore().schema = in.originalSchemaName.String()
		}()
	}
	return in.Model.db.DoSelect(ctx, in.link, in.Sql, in.Args...)
}

func (in *HookInsertInput) doInsert(ctx context.Context) (result sql.Result, err error) {
	// Schema change.
	if in.Schema != "" && in.Schema != in.originalSchemaName.String() {
		in.link, err = in.Model.db.GetCore().MasterLink(in.Schema)
		if err != nil {
			return nil, err
		}
		in.Model.db.GetCore().schema = in.Schema
		defer func() {
			in.Model.db.GetCore().schema = in.originalSchemaName.String()
		}()
	}
	return in.Model.db.DoInsert(ctx, in.link, in.Table, in.Data, in.Option)
}

func (in *HookUpdateInput) doUpdate(ctx context.Context) (result sql.Result, err error) {
	condition := in.restoreWhereForCommit()
	// Schema change.
	if in.Schema != "" && in.Schema != in.originalSchemaName.String() {
		in.link, err = in.Model.db.GetCore().MasterLink(in.Schema)
		if err != nil {
			return nil, err
		}
		in.Model.db.GetCore().schema = in.Schema
		defer func() {
			in.Model.db.GetCore().schema = in.originalSchemaName.String()
		}()
	}
	return in.Model.db.DoUpdate(ctx, in.link, in.Table, in.Data, condition, in.Args...)
}

func (in *HookDeleteInput) doDelete(ctx context.Context) (result sql.Result, err error) {
	condition := in.restoreWhereForCommit()
	// Schema change.
	if in.Schema != "" && in.Schema != in.originalSchemaName.String() {
		in.link, err = in.Model.db.GetCore().MasterLink(in.Schema)
		if err != nil {
			return nil, err
		}
		in.Model.db.GetCore().schema = in.Schema
		defer func() {
			in.Model.db.GetCore().schema = in.originalSchemaName.String()
		}()
	}
	return in.Model.db.DoDelete(ctx, in.link, in.Table, condition, in.Args...)
}
