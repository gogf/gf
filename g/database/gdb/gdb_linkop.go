// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gdb

import (
    "fmt"
    "errors"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

// 数据库链式操作对象
type gLinkOp struct {
    link          Link          // 数据库链接对象
    tables        string        // 数据库操作表
    fields        string        // 操作字段
    condition     string        // 操作条件
    conditionArgs []interface{} // 操作条件参数
    groupby       string        // 分组语句
    orderby       string        // 排序语句
    start         int           // 分页开始
    limit         int           // 分页条数
    data          interface{}   // 操作记录(支持Map/List/string类型)
    batch         int           // 批量操作条数
}

// 链式操作，数据表字段，可支持多个表，以半角逗号连接
func (l *dbLink) Table(tables string) (*gLinkOp) {
    return &gLinkOp{
        link  : l.link,
        tables: tables,
    }
}

// 链式操作，左联表
func (op *gLinkOp) LeftJoin(joinTable string, on string) (*gLinkOp) {
    op.tables += fmt.Sprintf(" LEFT JOIN %s ON (%s)", joinTable, on)
    return op
}

// 链式操作，右联表
func (op *gLinkOp) RightJoin(joinTable string, on string) (*gLinkOp) {
    op.tables += fmt.Sprintf(" RIGHT JOIN %s ON (%s)", joinTable, on)
    return op
}

// 链式操作，内联表
func (op *gLinkOp) InnerJoin(joinTable string, on string) (*gLinkOp) {
    op.tables += fmt.Sprintf(" INNER JOIN %s ON (%s)", joinTable, on)
    return op
}

// 链式操作，查询字段
func (op *gLinkOp) Fields(fields string) (*gLinkOp) {
    op.fields = fields
    return op
}

// 链式操作，consition
func (op *gLinkOp) Condition(condition string, args...interface{}) (*gLinkOp) {
    op.condition     = condition
    op.conditionArgs = args
    return op
}

// 链式操作，group by
func (op *gLinkOp) GroupBy(groupby string) (*gLinkOp) {
    op.groupby = groupby
    return op
}

// 链式操作，order by
func (op *gLinkOp) OrderBy(orderby string) (*gLinkOp) {
    op.orderby = orderby
    return op
}

// 链式操作，limit
func (op *gLinkOp) Limit(start int, limit int) (*gLinkOp) {
    op.start = start
    op.limit = limit
    return op
}

// 链式操作，操作数据记录项
func (op *gLinkOp) Data(data interface{}) (*gLinkOp) {
    op.data = data
    return op
}

// 链式操作， CURD - Insert/BatchInsert
func (op *gLinkOp) Insert() (sql.Result, error) {
    // 批量操作
    if list, ok :=  op.data.(List); ok {
        batch := 10
        if op.batch > 0 {
            batch = op.batch
        }
        return op.link.BatchInsert(op.tables, list, batch)
    }
    // 记录操作
    if op.data == nil {
        return nil, errors.New("inserting into table with empty data")
    }
    if d, ok :=  op.data.(Map); ok {
        return op.link.Insert(op.tables, d)
    }
    return nil, errors.New("inserting into table with invalid data type")
}

// 链式操作， CURD - Replace/BatchReplace
func (op *gLinkOp) Replace() (sql.Result, error) {
    // 批量操作
    if list, ok :=  op.data.(List); ok {
        batch := 10
        if op.batch > 0 {
            batch = op.batch
        }
        return op.link.BatchReplace(op.tables, list, batch)
    }
    // 记录操作
    if op.data == nil {
        return nil, errors.New("replacing into table with empty data")
    }
    if d, ok :=  op.data.(Map); ok {
        return op.link.Insert(op.tables, d)
    }
    return nil, errors.New("replacing into table with invalid data type")
}

// 链式操作， CURD - Save/BatchSave
func (op *gLinkOp) Save() (sql.Result, error) {
    // 批量操作
    if list, ok :=  op.data.(List); ok {
        batch := 10
        if op.batch > 0 {
            batch = op.batch
        }
        return op.link.BatchSave(op.tables, list, batch)
    }
    // 记录操作
    if op.data == nil {
        return nil, errors.New("saving into table with empty data")
    }
    if d, ok :=  op.data.(Map); ok {
        return op.link.Insert(op.tables, d)
    }
    return nil, errors.New("saving into table with invalid data type")
}

// 链式操作， CURD - Update
func (op *gLinkOp) Update() (sql.Result, error) {
    if op.data == nil {
        return nil, errors.New("updating table with empty data")
    }
    return op.link.Update(op.tables, op.data, op.condition, op.conditionArgs ...)
}

// 链式操作， CURD - Delete
func (op *gLinkOp) Delete() (sql.Result, error) {
    if op.condition == "" {
        return nil, errors.New("condition is required while deleting")
    }
    return op.link.Delete(op.tables, op.condition, op.conditionArgs...)
}

// 设置批处理的大小
func (op *gLinkOp) Batch(batch int) *gLinkOp {
    op.batch = batch
    return op
}

// 链式操作，select
func (op *gLinkOp) Select() (List, error) {
    if op.fields == "" {
        op.fields = "*"
    }
    s := fmt.Sprintf("SELECT %s FROM %s", op.fields, op.tables)
    if op.condition != "" {
        s += " WHERE " + op.condition
    }
    if op.groupby != "" {
        s += " GROUP BY " + op.groupby
    }
    if op.orderby != "" {
        s += " ORDER BY " + op.orderby
    }
    if op.limit != 0 {
        s += fmt.Sprintf(" LIMIT %d, %d", op.start, op.limit)
    }
    return op.link.GetAll(s, op.conditionArgs...)
}

// 链式操作，查询所有记录
func (op *gLinkOp) All() (List, error) {
    return op.Select()
}

// 链式操作，查询单条记录
func (op *gLinkOp) One() (Map, error) {
    list, err := op.All()
    if err != nil {
        return nil, err
    }
    return list[0], nil
}

// 链式操作，查询字段值
func (op *gLinkOp) Value() (interface{}, error) {
    one, err := op.One()
    if err != nil {
        return "", err
    }
    for _, v := range one {
        return v, nil
    }
    return "", nil
}

