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
	"github.com/ClickHouse/clickhouse-go"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/util/gconv"
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
	source := ""
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
	if db, err := sql.Open("clickhouse", source); err == nil {
		d.SetSchema(config.Name)
		return db, nil
	} else {
		return nil, err
	}
}

func (d *DriverClickhouse) Insert(ctx context.Context, table string, data interface{}, batch ...int) (sql.Result, error) {
	return nil, nil
}

func (d *DriverClickhouse) Tables(ctx context.Context, schema ...string) (tables []string, err error) {
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

func (d *DriverClickhouse) DoInsert(ctx context.Context, link Link, table string, list List, option DoInsertOption) (result sql.Result, err error) {
	return nil, nil
}

func (d *DriverClickhouse) GetChars() (charLeft string, charRight string) {
	return "", ""
}

func (d *DriverClickhouse) TableFields(ctx context.Context, table string, schema ...string) (fields map[string]*TableField, err error) {
	link, err := d.SlaveLink(schema...)
	if err != nil {
		return nil, err
	}
	getColumnsSql := fmt.Sprintf("select * from `system`.columns c where database = '%s' and `table` = '%s'", d.GetSchema(), table)
	result, err := d.DoGetAll(ctx, link, getColumnsSql)
	if err != nil {
		return nil, err
	}
	fields = make(map[string]*TableField)
	for _, m := range result {
		fields[m["name"].String()] = &TableField{
			Index:   m["position"].Int(),
			Name:    m["name"].String(),
			Type:    m["type"].String(),
			Comment: m["comment"].String(),
		}
	}
	return fields, nil
}

func (d *DriverClickhouse) InsertIgnore(ctx context.Context, table string, data interface{}, batch ...int) (sql.Result, error) {
	return nil, nil
}

func (d *DriverClickhouse) InsertAndGetId(ctx context.Context, table string, data interface{}, batch ...int) (int64, error) {
	return 0, nil
}

func (d *DriverClickhouse) Replace(ctx context.Context, table string, data interface{}, batch ...int) (sql.Result, error) {
	return nil, nil
}

func (d *DriverClickhouse) Save(ctx context.Context, table string, data interface{}, batch ...int) (sql.Result, error) {
	return nil, nil
}

func (d *DriverClickhouse) Update(ctx context.Context, table string, data interface{}, condition interface{}, args ...interface{}) (sql.Result, error) {
	return nil, nil
}

func (d *DriverClickhouse) Delete(ctx context.Context, table string, condition interface{}, args ...interface{}) (sql.Result, error) {
	return nil, nil
}
