// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"database/sql"
	"github.com/ClickHouse/clickhouse-go"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/guid"
)

// DriverClickhouse is the driver for SQL server database.
type DriverClickhouse struct {
	*Core
}

func (d *DriverClickhouse) New(core *Core, node *ConfigNode) (DB, error) {
	return &DriverClickhouse{
		Core: core,
	}, nil
}

func (d *DriverClickhouse) Open(config *ConfigNode) (db *sql.DB, err error) {
	var (
		source string
		driver = "clickhouse"
	)
	if config.Pass != "" {
		source = fmt.Sprintf(
			"tcp://%s:%s?database=%s&password=%s&charset=%s&debug=%s",
			config.Host, config.Port, config.Name, config.Pass, config.Charset, gconv.String(config.Debug),
		)
	} else {
		source = fmt.Sprintf(
			"tcp://%s:%s?database=%s&charset=%s&debug=%s",
			config.Host, config.Port, config.Name, config.Charset, gconv.String(config.Debug),
		)
	}
	intlog.Printf(d.GetCtx(), "Open: %s", source)
	if db, err := sql.Open(driver, source); err == nil {
		d.SetSchema(config.Name)
		return db, nil
	} else {
		return nil, err
	}
}

// Tables Get all tables from system tables record.
func (d *DriverClickhouse) Tables(ctx context.Context, schema ...string) (tables []string, err error) {
	var result Result
	link, err := d.SlaveLink(schema...)
	if err != nil {
		return nil, err
	}
	result, err = d.DoGetAll(ctx, link, fmt.Sprintf("select name from `system`.tables where database = '%s'", d.GetSchema()))
	if err != nil {
		return
	}
	for _, m := range result {
		tables = append(tables, m["name"].String())
	}
	return
}

// TableFields Get
func (d *DriverClickhouse) TableFields(ctx context.Context, table string, schema ...string) (fields map[string]*TableField, err error) {
	charL, charR := d.GetChars()
	table = gstr.Trim(table, charL+charR)
	if gstr.Contains(table, " ") {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "function TableFields supports only single table operations")
	}
	useSchema := d.schema.Val()
	if len(schema) > 0 && schema[0] != "" {
		useSchema = schema[0]
	}
	v := tableFieldsMap.GetOrSetFuncLock(
		fmt.Sprintf(`clickhouse_table_fields_%s_%s@group:%s`, table, useSchema, d.GetGroup()),
		func() interface{} {
			var (
				result Result
				link   Link
			)
			if link, err = d.SlaveLink(useSchema); err != nil {
				return nil
			}
			getColumnsSql := fmt.Sprintf("select name,position,default_expression,comment from `system`.columns c where database = '%s' and `table` = '%s'", d.GetSchema(), table)
			result, err := d.DoGetAll(ctx, link, getColumnsSql)
			if err != nil {
				return nil
			}
			fields = make(map[string]*TableField)
			for _, m := range result {
				var (
					isNull    = false
					fieldType = m["type"].String()
				)
				// in clickhouse , filed type like is Nullable(int)
				fieldsResult, _ := gregex.MatchString(`^Nullable\((.*?)\)`, fieldType)
				if len(fieldsResult) == 2 {
					isNull = true
					fieldType = fieldsResult[1]
				}
				fields[m["name"].String()] = &TableField{
					Index:   m["position"].Int(),
					Name:    m["name"].String(),
					Default: m["default_expression"].Val(),
					Comment: m["comment"].String(),
					//Key:     m["Key"].String(),
					Type: fieldType,
					Null: isNull,
				}
			}
			return fields
		},
	)
	if v != nil {
		fields = v.(map[string]*TableField)
	}
	return
}

func (d *DriverClickhouse) FilteredLink() string {
	linkInfo := d.GetConfig().Link
	if linkInfo == "" {
		return ""
	}
	s, _ := gregex.ReplaceString(
		`(.+?):(.+)@tcp(.+)`,
		`$1:xxx@tcp$3`,
		linkInfo,
	)
	return s
}

func (d *DriverClickhouse) PingMaster() error {
	conn, err := d.Master()
	if err != nil {
		return err
	}
	return d.ping(conn)

}

func (d *DriverClickhouse) PingSlave() error {
	conn, err := d.Slave()
	if err != nil {
		return err
	}
	return d.ping(conn)
}

func (d *DriverClickhouse) ping(conn *sql.DB) error {
	err := conn.Ping()
	if exception, ok := err.(*clickhouse.Exception); ok {
		return errors.New(fmt.Sprintf("[%d]%s", exception.Code, exception.Message))
	}
	return err
}

func (d *DriverClickhouse) DoUpdateSQL(ctx context.Context, link Link, table string, updates interface{}, condition string, args ...interface{}) (result sql.Result, err error) {
	// in clickhouse ,use update must use alter
	// ALTER TABLE [db.]table UPDATE column1 = expr1 [, ...] WHERE filter_expr
	return d.db.DoExec(ctx, link, fmt.Sprintf("ALTER TABLE %s UPDATE %s%s", table, updates, condition), args...)
}

func (d *DriverClickhouse) DoDelete(ctx context.Context, link Link, table string, condition string, args ...interface{}) (result sql.Result, err error) {
	if link == nil {
		if link, err = d.MasterLink(); err != nil {
			return nil, err
		}
	}
	table = d.QuotePrefixTableName(table)
	// in clickhouse , delete must use alter
	// ALTER TABLE [db.]table DELETE WHERE filter_expr
	return d.db.DoExec(ctx, link, fmt.Sprintf("ALTER TABLE %s DELETE %s", table, condition), args...)
}

func (d *DriverClickhouse) Transaction(ctx context.Context, f func(ctx context.Context, tx *TX) error) error {
	return errors.New("transaction operations are not supported")
}

func (d *DriverClickhouse) DoCommit(ctx context.Context, in DoCommitInput) (out DoCommitOutput, err error) {
	var (
		sqlTx                *sql.Tx
		sqlStmt              *sql.Stmt
		sqlRows              *sql.Rows
		sqlResult            sql.Result
		stmtSqlRows          *sql.Rows
		stmtSqlRow           *sql.Row
		rowsAffected         int64
		cancelFuncForTimeout context.CancelFunc
		timestampMilli1      = gtime.TimestampMilli()
	)
	// Execution cased by type.
	switch in.Type {
	case SqlTypeBegin:
		if sqlTx, err = in.Db.Begin(); err == nil {
			out.Tx = &TX{
				db:            d.db,
				tx:            sqlTx,
				ctx:           context.WithValue(ctx, transactionIdForLoggerCtx, transactionIdGenerator.Add(1)),
				master:        in.Db,
				transactionId: guid.S(),
			}
			ctx = out.Tx.ctx
		}
		out.RawResult = sqlTx

	case SqlTypeTXCommit:
		err = in.Tx.Commit()

	case SqlTypeTXRollback:
		// Clickhouse does not support the transaction
		// But it is necessary to submit the transaction after entering the transaction
		// So shields the rollback event.
		err = nil

	case SqlTypeExecContext:
		if d.db.GetDryRun() {
			sqlResult = new(SqlResult)
		} else {
			sqlResult, err = in.Link.ExecContext(ctx, in.Sql, in.Args...)
		}
		out.RawResult = sqlResult

	case SqlTypeQueryContext:
		sqlRows, err = in.Link.QueryContext(ctx, in.Sql, in.Args...)
		out.RawResult = sqlRows

	case SqlTypePrepareContext:
		sqlStmt, err = in.Link.PrepareContext(ctx, in.Sql)
		out.RawResult = sqlStmt

	case SqlTypeStmtExecContext:
		ctx, cancelFuncForTimeout = d.GetCtxTimeout(ctxTimeoutTypeExec, ctx)
		defer cancelFuncForTimeout()
		if d.db.GetDryRun() {
			sqlResult = new(SqlResult)
		} else {
			sqlResult, err = in.Stmt.ExecContext(ctx, in.Args...)
		}
		out.RawResult = sqlResult

	case SqlTypeStmtQueryContext:
		ctx, cancelFuncForTimeout = d.GetCtxTimeout(ctxTimeoutTypeQuery, ctx)
		defer cancelFuncForTimeout()
		stmtSqlRows, err = in.Stmt.QueryContext(ctx, in.Args...)
		out.RawResult = stmtSqlRows

	case SqlTypeStmtQueryRowContext:
		ctx, cancelFuncForTimeout = d.GetCtxTimeout(ctxTimeoutTypeQuery, ctx)
		defer cancelFuncForTimeout()
		stmtSqlRow = in.Stmt.QueryRowContext(ctx, in.Args...)
		out.RawResult = stmtSqlRow

	default:
		panic(gerror.NewCodef(gcode.CodeInvalidParameter, `invalid SqlType "%s"`, in.Type))
	}
	// Result handling.
	switch {
	case sqlResult != nil:
		// RowsAffected is not supported , so return default result
		rowsAffected, err = 0, nil
		out.Result = sqlResult

	case sqlRows != nil:
		out.Records, err = d.RowsToResult(ctx, sqlRows)
		rowsAffected = int64(len(out.Records))

	case sqlStmt != nil:
		out.Stmt = &Stmt{
			Stmt: sqlStmt,
			core: d.Core,
			link: in.Link,
			sql:  in.Sql,
		}
	}
	var (
		timestampMilli2 = gtime.TimestampMilli()
		sqlObj          = &Sql{
			Sql:           in.Sql,
			Type:          in.Type,
			Args:          in.Args,
			Format:        FormatSqlWithArgs(in.Sql, in.Args),
			Error:         err,
			Start:         timestampMilli1,
			End:           timestampMilli2,
			Group:         d.db.GetGroup(),
			RowsAffected:  rowsAffected,
			IsTransaction: in.IsTransaction,
		}
	)
	// Tracing and logging.
	d.addSqlToTracing(ctx, sqlObj)
	if d.db.GetDebug() {
		d.writeSqlToLogger(ctx, sqlObj)
	}
	if err != nil && err != sql.ErrNoRows {
		err = gerror.NewCodef(
			gcode.CodeDbOperationError,
			"%s, %s\n",
			err.Error(),
			FormatSqlWithArgs(in.Sql, in.Args),
		)
	}
	return out, err
}

func (d *DriverClickhouse) DoInsert(ctx context.Context, link Link, table string, list List, option DoInsertOption) (result sql.Result, err error) {
	var (
		fields     []string
		question   = []string{}
		listLength = len(list)
	)
	for item := range list[0] {
		fields = append(fields, item)
		question = append(question, "?")
	}
	if link.IsTransaction() {
		return nil, errors.New("transaction operations are not supported")
	}
	tx, err := d.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	sqlStr := fmt.Sprintf("INSERT INTO %v(%v) VALUES(%v)", table, strings.Join(fields, ","), strings.Join(question, ","))
	stmt, err := tx.Prepare(sqlStr)
	if err != nil {
		return nil, err
	}
	for i := 0; i < listLength; i++ {
		valueInterfaceSlice := []interface{}{}
		for _, filed := range fields {
			valueInterfaceSlice = append(valueInterfaceSlice, list[i][filed])
		}
		// TODO Clickhouse does not support the number of inserts, but can rely on this to get
		result, err = stmt.ExecContext(ctx, valueInterfaceSlice...)
		if err != nil {
			return nil, err
		}
	}
	return result, tx.Commit()
}
