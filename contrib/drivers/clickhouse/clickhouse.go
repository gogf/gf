// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package clickhouse implements gdb.Driver, which supports operations for ClickHouse.
package clickhouse

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/ClickHouse/clickhouse-go"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

// Driver is the driver for postgresql database.
type Driver struct {
	*gdb.Core
}

var (
	// tableFieldsMap caches the table information retrieved from database.
	tableFieldsMap             = gmap.New(true)
	ErrUnsupportedInsertIgnore = errors.New("unsupported method:InsertIgnore")
	ErrUnsupportedInsertGetId  = errors.New("unsupported method:InsertGetId")
	ErrUnsupportedReplace      = errors.New("unsupported method:Replace")
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
func (d *Driver) Open(config *gdb.ConfigNode) (*sql.DB, error) {
	var (
		source string
		driver = "clickhouse"
	)
	if config.Pass != "" {
		source = fmt.Sprintf(
			"clickhouse://%s:%s@%s:%s/%s",
			config.User, config.Pass, config.Host, config.Port, config.Name)
	} else {
		source = fmt.Sprintf(
			"clickhouse://%s@%s:%s/%s",
			config.User, config.Host, config.Port, config.Name)
	}
	source += fmt.Sprintf(
		"?charset=%s&debug=%s&compress=%s",
		config.Charset, gconv.String(config.Debug), gconv.String(config.Compress),
	)
	db, err := sql.Open(driver, source)
	if err != nil {
		return nil, err
	}

	return db, nil
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
	result, err = d.DoGetAll(ctx, link, query)
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
func (d *Driver) TableFields(ctx context.Context, table string, schema ...string) (fields map[string]*gdb.TableField, err error) {
	charL, charR := d.GetChars()
	table = gstr.Trim(table, charL+charR)
	if gstr.Contains(table, " ") {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "function TableFields supports only single table operations")
	}
	useSchema := d.GetSchema()
	if len(schema) > 0 && schema[0] != "" {
		useSchema = schema[0]
	}
	v := tableFieldsMap.GetOrSetFuncLock(
		fmt.Sprintf(`clickhouse_table_fields_%s_%s@group:%s`, table, useSchema, d.GetGroup()),
		func() interface{} {
			var (
				result gdb.Result
				link   gdb.Link
			)
			if link, err = d.SlaveLink(useSchema); err != nil {
				return nil
			}
			getColumnsSql := fmt.Sprintf("select name,position,default_expression,comment from `system`.columns c where database = '%s' and `table` = '%s'", d.GetConfig().Name, table)
			result, err = d.DoGetAll(ctx, link, getColumnsSql)
			if err != nil {
				return nil
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
		fields = v.(map[string]*gdb.TableField)
	}
	return
}

// FilteredLink retrieves and returns filtered `linkInfo` that can be using for
// logging or tracing purpose.
func (d *Driver) FilteredLink() string {
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
		return errors.New(fmt.Sprintf("[%d]%s", exception.Code, exception.Message))
	}
	return err
}

// Transaction Clickhouse does not support transactions
// So when you call this method you get an error.
func (d *Driver) Transaction(ctx context.Context, f func(ctx context.Context, tx *gdb.TX) error) error {
	return errors.New("transaction operations are not supported")
}

// DoUpdateSQL in clickhouse ,use update must use alter
// eg.
// ALTER TABLE [db.]table UPDATE column1 = expr1 [, ...] WHERE filter_expr
func (d *Driver) DoUpdateSQL(ctx context.Context, link gdb.Link, table string, updates interface{}, condition string, args ...interface{}) (result sql.Result, err error) {
	return d.Core.DoExec(ctx, link, fmt.Sprintf("ALTER TABLE %s UPDATE %s%s", table, updates, condition), args...)
}

// DoDeleteSQL in clickhouse , delete must use alter
// eg.
// ALTER TABLE [db.]table DELETE WHERE filter_expr
func (d *Driver) DoDeleteSQL(ctx context.Context, link gdb.Link, table string, condition interface{}, args ...interface{}) (result sql.Result, err error) {
	return d.Core.DoExec(ctx, link, fmt.Sprintf("ALTER TABLE %s DELETE %s", table, condition), args...)
}

func (d *Driver) DoInsert(ctx context.Context, link gdb.Link, table string, data gdb.List, option gdb.DoInsertOption) (result sql.Result, err error) {
	return
}

func (d *Driver) DoCommit(ctx context.Context, in gdb.DoCommitInput) (out gdb.DoCommitOutput, err error) {
	in.IsIgnoreResult = true
	return d.Core.DoCommit(ctx, in)
}

// InsertIgnore Other queries for modifying data parts are not supported: REPLACE, MERGE, UPSERT, INSERT UPDATE.
func (d *Driver) InsertIgnore(ctx context.Context, table string, data interface{}, batch ...int) (sql.Result, error) {
	return nil, ErrUnsupportedInsertIgnore
}

// InsertAndGetId Other queries for modifying data parts are not supported: REPLACE, MERGE, UPSERT, INSERT UPDATE.
func (d *Driver) InsertAndGetId(ctx context.Context, table string, data interface{}, batch ...int) (int64, error) {
	return 0, ErrUnsupportedInsertGetId
}

// Replace Other queries for modifying data parts are not supported: REPLACE, MERGE, UPSERT, INSERT UPDATE.
func (d *Driver) Replace(ctx context.Context, table string, data interface{}, batch ...int) (sql.Result, error) {
	return nil, ErrUnsupportedReplace
}
