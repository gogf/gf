// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"database/sql"
)

type DriverClickhouse struct {
	*Core
}

func (d *DriverClickhouse) New(core *Core, node *ConfigNode) (DB, error) {
	return &DriverMssql{
		Core: core,
	}, nil
}

func (d *DriverClickhouse) Open(config *ConfigNode) (db *sql.DB, err error) {
	return nil, nil
}

func (d *DriverClickhouse) Insert(ctx context.Context, table string, data interface{}, batch ...int) (sql.Result, error) {
	return nil, nil
}

func (d *DriverClickhouse) Ctx(ctx context.Context) DB {
	return nil
}

func (d *DriverClickhouse) Tables(ctx context.Context, schema ...string) (tables []string, err error) {
	return nil, nil
}

func (d *DriverClickhouse) FilteredLink() string {
	return ""
}

func (d *DriverClickhouse) GetCore() *Core {
	return nil
}

func (d *DriverClickhouse) DoInsert(ctx context.Context, link Link, table string, list List, option DoInsertOption) (result sql.Result, err error) {
	return nil, nil
}

func (d *DriverClickhouse) Model(tableNameOrStruct ...interface{}) *Model {
	return nil
}

func (d *DriverClickhouse) GetChars() (charLeft string, charRight string) {
	return "", ""
}

func (d *DriverClickhouse) DoFilter(ctx context.Context, link Link, sql string, args []interface{}) (newSql string, newArgs []interface{}, err error) {
	return "", nil, nil
}

func (d *DriverClickhouse) TableFields(ctx context.Context, table string, schema ...string) (fields map[string]*TableField, err error) {
	return nil, nil
}
