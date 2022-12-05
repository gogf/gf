// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package clickhouse implements gdb.Driver, which supports operations for database ClickHouse.
package clickhouse

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gtag"
	"github.com/gogf/gf/v2/util/gutil"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Driver is the driver for postgresql database.
type Driver struct {
	*gdb.Core
}

var (
	errUnsupportedInsertIgnore = errors.New("unsupported method:InsertIgnore")
	errUnsupportedInsertGetId  = errors.New("unsupported method:InsertGetId")
	errUnsupportedReplace      = errors.New("unsupported method:Replace")
	errUnsupportedBegin        = errors.New("unsupported method:Begin")
	errUnsupportedTransaction  = errors.New("unsupported method:Transaction")
)

const (
	updateFilterPattern              = `(?i)UPDATE[\s]+?(\w+[\.]?\w+)[\s]+?SET`
	deleteFilterPattern              = `(?i)DELETE[\s]+?FROM[\s]+?(\w+[\.]?\w+)`
	filterTypePattern                = `(?i)^UPDATE|DELETE`
	replaceSchemaPattern             = `@(.+?)/([\w\.\-]+)+`
	needParsedSqlInCtx   gctx.StrKey = "NeedParsedSql"
	OrmTagForStruct                  = gtag.ORM
	driverName                       = "clickhouse"
)

func init() {
	if err := gdb.Register(`clickhouse`, New()); err != nil {
		panic(err)
	}
}

// New create and returns a driver that implements gdb.Driver, which supports operations for clickhouse.
func New() gdb.Driver {
	return &Driver{}
}

// New creates and returns a database object for clickhouse.
// It implements the interface of gdb.Driver for extra database driver installation.
func (d *Driver) New(core *gdb.Core, node *gdb.ConfigNode) (gdb.DB, error) {
	return &Driver{
		Core: core,
	}, nil
}

// Open creates and returns an underlying sql.DB object for clickhouse.
func (d *Driver) Open(config *gdb.ConfigNode) (db *sql.DB, err error) {
	source := config.Link
	// clickhouse://username:password@host1:9000,host2:9000/database?dial_timeout=200ms&max_execution_time=60
	if config.Link != "" {
		// ============================================================================
		// Deprecated from v2.2.0.
		// ============================================================================
		// Custom changing the schema in runtime.
		if config.Name != "" {
			source, _ = gregex.ReplaceString(replaceSchemaPattern, "@$1/"+config.Name, config.Link)
		} else {
			// If no schema, the link is matched for replacement
			dbName, _ := gregex.MatchString(replaceSchemaPattern, config.Link)
			if len(dbName) > 0 {
				config.Name = dbName[len(dbName)-1]
			}
		}
	} else {
		if config.Pass != "" {
			source = fmt.Sprintf(
				"clickhouse://%s:%s@%s:%s/%s?charset=%s&debug=%t",
				config.User, url.PathEscape(config.Pass),
				config.Host, config.Port, config.Name, config.Charset, config.Debug,
			)
		} else {
			source = fmt.Sprintf(
				"clickhouse://%s@%s:%s/%s?charset=%s&debug=%t",
				config.User, config.Host, config.Port, config.Name, config.Charset, config.Debug,
			)
		}
		if config.Extra != "" {
			source = fmt.Sprintf("%s&%s", source, config.Extra)
		}
	}
	if db, err = sql.Open(driverName, source); err != nil {
		err = gerror.WrapCodef(
			gcode.CodeDbOperationError, err,
			`sql.Open failed for driver "%s" by source "%s"`, driverName, source,
		)
		return nil, err
	}
	return
}

// Tables retrieves and returns the tables of current schema.
// It's mainly used in cli tool chain for automatically generating the models.
func (d *Driver) Tables(ctx context.Context, schema ...string) (tables []string, err error) {
	var result gdb.Result
	link, err := d.SlaveLink(schema...)
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf("select name from `system`.tables where database = '%s'", d.GetConfig().Name)
	result, err = d.DoSelect(ctx, link, query)
	if err != nil {
		return
	}
	for _, m := range result {
		tables = append(tables, m["name"].String())
	}
	return
}

// TableFields retrieves and returns the fields' information of specified table of current schema.
// Also see DriverMysql.TableFields.
func (d *Driver) TableFields(
	ctx context.Context, table string, schema ...string,
) (fields map[string]*gdb.TableField, err error) {
	var (
		result    gdb.Result
		link      gdb.Link
		useSchema = gutil.GetOrDefaultStr(d.GetSchema(), schema...)
	)
	if link, err = d.SlaveLink(useSchema); err != nil {
		return nil, err
	}
	var (
		columns       = "name,position,default_expression,comment,type,is_in_partition_key,is_in_sorting_key,is_in_primary_key,is_in_sampling_key"
		getColumnsSql = fmt.Sprintf(
			"select %s from `system`.columns c where `table` = '%s'",
			columns, table,
		)
	)
	result, err = d.DoSelect(ctx, link, getColumnsSql)
	if err != nil {
		return nil, err
	}
	fields = make(map[string]*gdb.TableField)
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
		fields[m["name"].String()] = &gdb.TableField{
			Index:   m["position"].Int() - 1,
			Name:    m["name"].String(),
			Default: m["default_expression"].Val(),
			Comment: m["comment"].String(),
			// Key:     m["Key"].String(),
			Type: fieldType,
			Null: isNull,
		}
	}
	return fields, nil
}

// PingMaster pings the master node to check authentication or keeps the connection alive.
func (d *Driver) PingMaster() error {
	conn, err := d.Master()
	if err != nil {
		return err
	}
	return d.ping(conn)
}

// PingSlave pings the slave node to check authentication or keeps the connection alive.
func (d *Driver) PingSlave() error {
	conn, err := d.Slave()
	if err != nil {
		return err
	}
	return d.ping(conn)
}

// ping Returns the Clickhouse specific error.
func (d *Driver) ping(conn *sql.DB) error {
	err := conn.Ping()
	if exception, ok := err.(*clickhouse.Exception); ok {
		return fmt.Errorf("[%d]%s", exception.Code, exception.Message)
	}
	return err
}

// DoFilter handles the sql before posts it to database.
func (d *Driver) DoFilter(
	ctx context.Context, link gdb.Link, originSql string, args []interface{},
) (newSql string, newArgs []interface{}, err error) {
	if len(args) == 0 {
		return originSql, args, nil
	}

	var index int
	// Convert placeholder char '?' to string "$x".
	originSql, _ = gregex.ReplaceStringFunc(`\?`, originSql, func(s string) string {
		index++
		return fmt.Sprintf(`$%d`, index)
	})

	// Only SQL generated through the framework is processed.
	if !d.getNeedParsedSqlFromCtx(ctx) {
		return originSql, args, nil
	}

	// replace STD SQL to Clickhouse SQL grammar
	modeRes, err := gregex.MatchString(filterTypePattern, strings.TrimSpace(originSql))
	if err != nil {
		return "", nil, err
	}
	if len(modeRes) == 0 {
		return originSql, args, nil
	}

	// Only delete/ UPDATE statements require filter
	switch strings.ToUpper(modeRes[0]) {
	case "UPDATE":
		// MySQL eg: UPDATE table_name SET field1=new-value1, field2=new-value2 [WHERE Clause]
		// Clickhouse eg: ALTER TABLE [db.]table UPDATE column1 = expr1 [, ...] WHERE filter_expr
		newSql, err = gregex.ReplaceStringFuncMatch(updateFilterPattern, originSql, func(s []string) string {
			return fmt.Sprintf("ALTER TABLE %s UPDATE", s[1])
		})
		if err != nil {
			return "", nil, err
		}
		return newSql, args, nil

	case "DELETE":
		// MySQL eg: DELETE FROM table_name [WHERE Clause]
		// Clickhouse eg: ALTER TABLE [db.]table [ON CLUSTER cluster] DELETE WHERE filter_expr
		newSql, err = gregex.ReplaceStringFuncMatch(deleteFilterPattern, originSql, func(s []string) string {
			return fmt.Sprintf("ALTER TABLE %s DELETE", s[1])
		})
		if err != nil {
			return "", nil, err
		}
		return newSql, args, nil

	}
	return originSql, args, nil
}

// DoCommit commits current sql and arguments to underlying sql driver.
func (d *Driver) DoCommit(ctx context.Context, in gdb.DoCommitInput) (out gdb.DoCommitOutput, err error) {
	ctx = d.InjectIgnoreResult(ctx)
	return d.Core.DoCommit(ctx, in)
}

func (d *Driver) DoInsert(
	ctx context.Context, link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption,
) (result sql.Result, err error) {
	var (
		keys        []string // Field names.
		valueHolder = make([]string, 0)
	)
	// Handle the field names and placeholders.
	for k := range list[0] {
		keys = append(keys, k)
		valueHolder = append(valueHolder, "?")
	}
	// Prepare the batch result pointer.
	var (
		charL, charR = d.Core.GetChars()
		keysStr      = charL + strings.Join(keys, charR+","+charL) + charR
		holderStr    = strings.Join(valueHolder, ",")
		tx           = &gdb.TX{}
		stdSqlResult sql.Result
		stmt         *gdb.Stmt
	)
	tx, err = d.Core.Begin(ctx)
	if err != nil {
		return
	}
	stmt, err = tx.Prepare(fmt.Sprintf(
		"INSERT INTO %s(%s) VALUES (%s)",
		d.QuotePrefixTableName(table), keysStr,
		holderStr,
	))
	if err != nil {
		return
	}
	for i := 0; i < len(list); i++ {
		params := make([]interface{}, 0) // Values that will be committed to underlying database driver.
		for _, k := range keys {
			params = append(params, list[i][k])
		}
		// Prepare is allowed to execute only once in a transaction opened by clickhouse
		stdSqlResult, err = stmt.ExecContext(ctx, params...)
		if err != nil {
			return stdSqlResult, err
		}
	}
	return stdSqlResult, tx.Commit()
}

// ConvertDataForRecord converting for any data that will be inserted into table/collection as a record.
func (d *Driver) ConvertDataForRecord(ctx context.Context, value interface{}) (map[string]interface{}, error) {
	m := gconv.Map(value, OrmTagForStruct)

	// transforms a value of a particular type
	for k, v := range m {
		switch itemValue := v.(type) {

		case time.Time:
			m[k] = itemValue
			// If the time is zero, it then updates it to nil,
			// which will insert/update the value to database as "null".
			if itemValue.IsZero() {
				m[k] = nil
			}

		case uuid.UUID:
			m[k] = itemValue

		case *time.Time:
			m[k] = itemValue
			// If the time is zero, it then updates it to nil,
			// which will insert/update the value to database as "null".
			if itemValue == nil || itemValue.IsZero() {
				m[k] = nil
			}

		case gtime.Time:
			// for gtime type, needs to get time.Time
			m[k] = itemValue.Time
			// If the time is zero, it then updates it to nil,
			// which will insert/update the value to database as "null".
			if itemValue.IsZero() {
				m[k] = nil
			}

		case *gtime.Time:
			// for gtime type, needs to get time.Time
			if itemValue != nil {
				m[k] = itemValue.Time
			}
			// If the time is zero, it then updates it to nil,
			// which will insert/update the value to database as "null".
			if itemValue == nil || itemValue.IsZero() {
				m[k] = nil
			}

		case decimal.Decimal:
			m[k] = itemValue

		case *decimal.Decimal:
			m[k] = nil
			if itemValue != nil {
				m[k] = *itemValue
			}

		default:
			// if the other type implements valuer for the driver package
			// the converted result is used
			// otherwise the interface data is committed
			valuer, ok := itemValue.(driver.Valuer)
			if !ok {
				m[k] = itemValue
				continue
			}
			convertedValue, err := valuer.Value()
			if err != nil {
				return nil, err
			}
			m[k] = convertedValue
		}
	}
	return m, nil
}

func (d *Driver) DoDelete(ctx context.Context, link gdb.Link, table string, condition string, args ...interface{}) (result sql.Result, err error) {
	ctx = d.injectNeedParsedSql(ctx)
	return d.Core.DoDelete(ctx, link, table, condition, args...)
}

func (d *Driver) DoUpdate(ctx context.Context, link gdb.Link, table string, data interface{}, condition string, args ...interface{}) (result sql.Result, err error) {
	ctx = d.injectNeedParsedSql(ctx)
	return d.Core.DoUpdate(ctx, link, table, data, condition, args...)
}

// InsertIgnore Other queries for modifying data parts are not supported: REPLACE, MERGE, UPSERT, INSERT UPDATE.
func (d *Driver) InsertIgnore(ctx context.Context, table string, data interface{}, batch ...int) (sql.Result, error) {
	return nil, errUnsupportedInsertIgnore
}

// InsertAndGetId Other queries for modifying data parts are not supported: REPLACE, MERGE, UPSERT, INSERT UPDATE.
func (d *Driver) InsertAndGetId(ctx context.Context, table string, data interface{}, batch ...int) (int64, error) {
	return 0, errUnsupportedInsertGetId
}

// Replace Other queries for modifying data parts are not supported: REPLACE, MERGE, UPSERT, INSERT UPDATE.
func (d *Driver) Replace(ctx context.Context, table string, data interface{}, batch ...int) (sql.Result, error) {
	return nil, errUnsupportedReplace
}

func (d *Driver) Begin(ctx context.Context) (tx *gdb.TX, err error) {
	return nil, errUnsupportedBegin
}

func (d *Driver) Transaction(ctx context.Context, f func(ctx context.Context, tx *gdb.TX) error) error {
	return errUnsupportedTransaction
}

func (d *Driver) injectNeedParsedSql(ctx context.Context) context.Context {
	if ctx.Value(needParsedSqlInCtx) != nil {
		return ctx
	}
	return context.WithValue(ctx, needParsedSqlInCtx, true)
}

func (d *Driver) getNeedParsedSqlFromCtx(ctx context.Context) bool {
	if ctx.Value(needParsedSqlInCtx) != nil {
		return true
	}
	return false
}
