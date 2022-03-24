// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/longbridgeapp/sqlparser"
)

// ShardingInput is input parameters for custom sharding handler.
type ShardingInput struct {
	Table         string           // Current operation table name.
	Schema        string           // Current operation schema, usually empty string which means uses default schema from configuration.
	OperationData map[string]Value // Accurate readonly key-value data pairs from INSERT/UPDATE statement.
	ConditionData map[string]Value // Accurate readonly key-value condition pairs from SELECT/UPDATE/DELETE statement.
}

// ShardingOutput is output parameters for custom sharding handler.
type ShardingOutput struct {
	Table  string // New table name for current operation. Use empty string for no changes of table name.
	Schema string // New schema name for current operation. Use empty string for using default schema from configuration.
}

// ShardingHandler is a custom function for custom sharding table and schema for DB operation.
type ShardingHandler func(ctx context.Context, in ShardingInput) (out *ShardingOutput, err error)

const (
	ctxKeyForShardingHandler gctx.StrKey = "ShardingHandler"
)

// Sharding creates and returns a new model with sharding handler.
func (m *Model) Sharding(handler ShardingHandler) *Model {
	var (
		ctx   = m.GetCtx()
		model = m.getModel()
	)
	model.shardingHandler = handler
	// Inject sharding handler into context.
	model = model.Ctx(model.injectShardingInputCaller(ctx))
	return model
}

// injectShardingInputCaller injects custom sharding handler into context.
func (m *Model) injectShardingInputCaller(ctx context.Context) context.Context {
	if m.shardingHandler == nil {
		return ctx
	}
	if ctx.Value(ctxKeyForShardingHandler) != nil {
		return ctx
	}
	return context.WithValue(ctx, ctxKeyForShardingHandler, m.shardingHandler)
}

type callShardingHandlerFromCtxInput struct {
	Sql          string
	FormattedSql string
}

type callShardingHandlerFromCtxOutput struct {
	Sql             string
	Table           string
	Schema          string
	ParsedSqlOutput *parseFormattedSqlOutput
}

func (c *Core) callShardingHandlerFromCtx(
	ctx context.Context, in callShardingHandlerFromCtxInput,
) (out *callShardingHandlerFromCtxOutput, err error) {
	var (
		newSql          = in.Sql
		ctxValue        interface{}
		shardingHandler ShardingHandler
		ok              bool
	)
	// If no sharding handler, it does nothing.
	if ctxValue = ctx.Value(ctxKeyForShardingHandler); ctxValue == nil {
		return nil, nil
	}
	if shardingHandler, ok = ctxValue.(ShardingHandler); !ok {
		return nil, nil
	}
	parsedOut, err := c.parseFormattedSql(in.FormattedSql)
	if err != nil {
		return nil, err
	}
	var shardingIn = ShardingInput{
		Table:         parsedOut.Table,
		Schema:        c.db.GetSchema(),
		OperationData: parsedOut.OperationData,
		ConditionData: parsedOut.ConditionData,
	}
	shardingOut, err := shardingHandler(ctx, shardingIn)
	if err != nil {
		return nil, gerror.Wrap(err, `calling sharding handler failed`)
	}
	if shardingOut.Table != shardingIn.Table || shardingOut.Schema != shardingIn.Schema {
		if shardingOut.Table != shardingIn.Table {
			newSql, err = c.formatSqlWithNewTable(in.Sql, shardingOut.Table)
			if err != nil {
				return nil, err
			}
		}
		out = &callShardingHandlerFromCtxOutput{
			Sql:             newSql,
			Table:           shardingOut.Table,
			Schema:          shardingOut.Schema,
			ParsedSqlOutput: parsedOut,
		}
		return out, nil
	}
	return nil, nil
}

// formatSqlWithNewTable modifies given `sql` and returns a sql with new table name `table`.
func (c *Core) formatSqlWithNewTable(sql, table string) (newSql string, err error) {
	parsedStmt, err := sqlparser.NewParser(strings.NewReader(sql)).ParseStatement()
	if err != nil {
		return "", gerror.Wrapf(err, `parse failed for SQL: %s`, sql)
	}
	newTable := &sqlparser.TableName{Name: &sqlparser.Ident{Name: table}}
	switch stmt := parsedStmt.(type) {
	case *sqlparser.SelectStatement:
		stmt.FromItems = newTable
		return stmt.String(), nil
	case *sqlparser.InsertStatement:
		stmt.TableName = newTable
		return stmt.String(), nil
	case *sqlparser.UpdateStatement:
		stmt.TableName = newTable
		return stmt.String(), nil
	case *sqlparser.DeleteStatement:
		stmt.TableName = newTable
		return stmt.String(), nil
	default:
		return "", gerror.Wrapf(err, `unsupported SQL: %s`, sql)
	}
}

type parseFormattedSqlOutput struct {
	Table          string
	OperationData  map[string]Value
	ConditionData  map[string]Value
	ParsedStmt     sqlparser.Statement
	SelectedFields []string
}

func (c *Core) parseFormattedSql(formattedSql string) (*parseFormattedSqlOutput, error) {
	var (
		condition sqlparser.Expr
		err       error
		out       = &parseFormattedSqlOutput{
			SelectedFields: make([]string, 0),
			OperationData:  make(map[string]Value),
			ConditionData:  make(map[string]Value),
		}
	)
	out.ParsedStmt, err = sqlparser.NewParser(strings.NewReader(formattedSql)).ParseStatement()
	if err != nil {
		return nil, gerror.Wrapf(err, `parse failed for SQL: %s`, formattedSql)
	}
	switch stmt := out.ParsedStmt.(type) {
	case *sqlparser.SelectStatement:
		if stmt.FromItems != nil {
			table, ok := stmt.FromItems.(*sqlparser.TableName)
			if !ok {
				return nil, gerror.Newf(
					`invalid table name "%s" in SQL: %s`,
					stmt.FromItems.String(), formattedSql,
				)
			}
			out.Table = table.TableName()
		}
		condition = stmt.Condition
		if stmt.Columns != nil {
			for _, column := range *stmt.Columns {
				if column.Alias != nil {
					out.SelectedFields = append(out.SelectedFields, column.Alias.Name)
				} else if column.Expr != nil {
					out.SelectedFields = append(out.SelectedFields, column.Expr.String())
				}
			}
		}

	case *sqlparser.InsertStatement:
		out.Table = stmt.TableName.TableName()
		if len(stmt.Expressions) > 0 && len(stmt.ColumnNames) > 0 {
			names := make([]string, len(stmt.ColumnNames))
			for i, ident := range stmt.ColumnNames {
				names[i] = ident.Name
			}
			// It just uses the first item.
			for i, expr := range stmt.Expressions[0].Exprs {
				c.injectDataByExpr(out.OperationData, names[i], expr)
			}
		}
	case *sqlparser.UpdateStatement:
		out.Table = stmt.TableName.TableName()
		condition = stmt.Condition
		if len(stmt.Assignments) > 0 {
			for _, assignment := range stmt.Assignments {
				if len(assignment.Columns) > 0 {
					c.injectDataByExpr(out.OperationData, assignment.Columns[0].Name, assignment.Expr)
				}
			}
		}
	case *sqlparser.DeleteStatement:
		out.Table = stmt.TableName.TableName()
		condition = stmt.Condition

	default:
		return nil, gerror.Wrapf(err, `unsupported SQL: %s`, formattedSql)
	}

	err = sqlparser.Walk(sqlparser.VisitFunc(func(node sqlparser.Node) error {
		if n, ok := node.(*sqlparser.BinaryExpr); ok {
			if x, ok := n.X.(*sqlparser.Ident); ok {
				if n.Op == sqlparser.EQ {
					c.injectDataByExpr(out.ConditionData, x.Name, n.Y)
				}
			}
		}
		return nil
	}), condition)
	return out, err
}

func (c *Core) injectDataByExpr(data map[string]Value, name string, expr sqlparser.Expr) {
	switch exprImp := expr.(type) {
	case *sqlparser.StringLit:
		data[name] = gvar.New(exprImp.Value)
	case *sqlparser.NumberLit:
		data[name] = gvar.New(exprImp.Value)
	default:
		data[name] = gvar.New(exprImp.String())
	}
}
