// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/ClickHouse/clickhouse-go"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"reflect"
	"strings"
)

// DriverClickhouse is the driver for SQL server database.
type DriverClickhouse struct {
	*Core
}

func (d *DriverClickhouse) New(core *Core, node *ConfigNode) (DB, error) {
	return &DriverMssql{
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
	glog.Infof(context.Background(), "Open: %s %s", source, clickhouse.DefaultDatabase)
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
	result, err = d.DoGetAll(ctx, link, fmt.Sprintf("select name from `system`.tables where database '%s'", d.GetSchema()))
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

func (d *DriverClickhouse) DoUpdate(ctx context.Context, link Link, table string, data interface{}, condition string, args ...interface{}) (result sql.Result, err error) {
	table = d.QuotePrefixTableName(table)
	var (
		rv   = reflect.ValueOf(data)
		kind = rv.Kind()
	)
	if kind == reflect.Ptr {
		rv = rv.Elem()
		kind = rv.Kind()
	}
	var (
		params  []interface{}
		updates = ""
	)
	switch kind {
	case reflect.Map, reflect.Struct:
		var (
			fields         []string
			dataMap        = d.db.ConvertDataForRecord(ctx, data)
			counterHandler = func(column string, counter Counter) {
				if counter.Value != 0 {
					column = d.QuoteWord(column)
					var (
						columnRef = d.QuoteWord(counter.Field)
						columnVal = counter.Value
						operator  = "+"
					)
					if columnVal < 0 {
						operator = "-"
						columnVal = -columnVal
					}
					fields = append(fields, fmt.Sprintf("%s=%s%s?", column, columnRef, operator))
					params = append(params, columnVal)
				}
			}
		)

		for k, v := range dataMap {
			switch value := v.(type) {
			case *Counter:
				counterHandler(k, *value)

			case Counter:
				counterHandler(k, value)

			default:
				if s, ok := v.(Raw); ok {
					fields = append(fields, d.QuoteWord(k)+"="+gconv.String(s))
				} else {
					fields = append(fields, d.QuoteWord(k)+"=?")
					params = append(params, v)
				}
			}
		}
		updates = strings.Join(fields, ",")

	default:
		updates = gconv.String(data)
	}
	if len(updates) == 0 {
		return nil, gerror.NewCode(gcode.CodeMissingParameter, "data cannot be empty")
	}
	if len(params) > 0 {
		args = append(params, args...)
	}
	// If no link passed, it then uses the master link.
	if link == nil {
		if link, err = d.MasterLink(); err != nil {
			return nil, err
		}
	}
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
	return d.db.DoExec(ctx, link, fmt.Sprintf("ALTER TABLE %s DELETE FROM %s", table, condition), args...)
}

func (d *DriverClickhouse) Transaction(ctx context.Context, f func(ctx context.Context, tx *TX) error) error {
	return errors.New("transaction operations are not supported")
}

func (d *DriverClickhouse) DoInsert(ctx context.Context, link Link, table string, data List, option DoInsertOption) (result sql.Result, err error) {
	return nil, nil
}

func (d *DriverClickhouse) PingMaster() error {
	return d.db.PingMaster()
}
