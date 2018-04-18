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
type DbOp struct {
    tx            *Tx           // 数据库事务对象
    db            *Db           // 数据库操作对象
    tables        string        // 数据库操作表
    fields        string        // 操作字段
    where         string        // 操作条件
    whereArgs     []interface{} // 操作条件参数
    groupby       string        // 分组语句
    orderby       string        // 排序语句
    start         int           // 分页开始
    limit         int           // 分页条数
    data          interface{}   // 操作记录(支持Map/List/string类型)
    batch         int           // 批量操作条数
}

// 链式操作，数据表字段，可支持多个表，以半角逗号连接
func (db *Db) Table(tables string) (*DbOp) {
    return &DbOp {
        db     : db,
        tables : tables,
    }
}

// 链式操作，数据表字段，可支持多个表，以半角逗号连接
func (db *Db) From(tables string) (*DbOp) {
    return db.Table(tables)
}

// (事务)链式操作，数据表字段，可支持多个表，以半角逗号连接
func (tx *Tx) Table(tables string) (*DbOp) {
    return &DbOp {
        tx     : tx,
        tables : tables,
    }
}

// (事务)链式操作，数据表字段，可支持多个表，以半角逗号连接
func (tx *Tx) From(tables string) (*DbOp) {
    return tx.Table(tables)
}

// 链式操作，左联表
func (op *DbOp) LeftJoin(joinTable string, on string) (*DbOp) {
    op.tables += fmt.Sprintf(" LEFT JOIN %s ON (%s)", joinTable, on)
    return op
}

// 链式操作，右联表
func (op *DbOp) RightJoin(joinTable string, on string) (*DbOp) {
    op.tables += fmt.Sprintf(" RIGHT JOIN %s ON (%s)", joinTable, on)
    return op
}

// 链式操作，内联表
func (op *DbOp) InnerJoin(joinTable string, on string) (*DbOp) {
    op.tables += fmt.Sprintf(" INNER JOIN %s ON (%s)", joinTable, on)
    return op
}

// 链式操作，查询字段
func (op *DbOp) Fields(fields string) (*DbOp) {
    op.fields = fields
    return op
}

// 链式操作，consition
func (op *DbOp) Where(where string, args...interface{}) (*DbOp) {
    op.where     = where
    op.whereArgs = args
    return op
}

// 链式操作，group by
func (op *DbOp) GroupBy(groupby string) (*DbOp) {
    op.groupby = groupby
    return op
}

// 链式操作，order by
func (op *DbOp) OrderBy(orderby string) (*DbOp) {
    op.orderby = orderby
    return op
}

// 链式操作，limit
func (op *DbOp) Limit(start int, limit int) (*DbOp) {
    op.start = start
    op.limit = limit
    return op
}

// 链式操作，操作数据记录项
func (op *DbOp) Data(data interface{}) (*DbOp) {
    op.data = data
    return op
}

// 链式操作， CURD - Insert/BatchInsert
func (op *DbOp) Insert() (sql.Result, error) {
    // 批量操作
    if list, ok :=  op.data.(List); ok {
        batch := 10
        if op.batch > 0 {
            batch = op.batch
        }
        if op.tx == nil {
            return op.db.BatchInsert(op.tables, list, batch)
        } else {
            return op.tx.BatchInsert(op.tables, list, batch)
        }
    }
    // 记录操作
    if op.data == nil {
        return nil, errors.New("inserting into table with empty data")
    }
    if dataMap, ok :=  op.data.(Map); ok {
        if op.tx == nil {
            return op.db.Insert(op.tables, dataMap)
        } else {
            return op.tx.Insert(op.tables, dataMap)
        }
    }
    return nil, errors.New("inserting into table with invalid data type")
}

// 链式操作， CURD - Replace/BatchReplace
func (op *DbOp) Replace() (sql.Result, error) {
    // 批量操作
    if list, ok :=  op.data.(List); ok {
        batch := 10
        if op.batch > 0 {
            batch = op.batch
        }
        if op.tx == nil {
            return op.db.BatchReplace(op.tables, list, batch)
        } else {
            return op.tx.BatchReplace(op.tables, list, batch)
        }
    }
    // 记录操作
    if op.data == nil {
        return nil, errors.New("replacing into table with empty data")
    }
    if dataMap, ok :=  op.data.(Map); ok {
        if op.tx == nil {
            return op.db.Insert(op.tables, dataMap)
        } else {
            return op.tx.Insert(op.tables, dataMap)
        }
    }
    return nil, errors.New("replacing into table with invalid data type")
}

// 链式操作， CURD - Save/BatchSave
func (op *DbOp) Save() (sql.Result, error) {
    // 批量操作
    if list, ok :=  op.data.(List); ok {
        batch := 10
        if op.batch > 0 {
            batch = op.batch
        }
        if op.tx == nil {
            return op.db.BatchSave(op.tables, list, batch)
        } else {
            return op.tx.BatchSave(op.tables, list, batch)
        }
    }
    // 记录操作
    if op.data == nil {
        return nil, errors.New("saving into table with empty data")
    }
    if dataMap, ok :=  op.data.(Map); ok {
        if op.tx == nil {
            return op.db.Save(op.tables, dataMap)
        } else {
            return op.tx.Save(op.tables, dataMap)
        }
    }
    return nil, errors.New("saving into table with invalid data type")
}

// 链式操作， CURD - Update
func (op *DbOp) Update() (sql.Result, error) {
    if op.data == nil {
        return nil, errors.New("updating table with empty data")
    }
    if op.tx == nil {
        return op.db.Update(op.tables, op.data, op.where, op.whereArgs ...)
    } else {
        return op.tx.Update(op.tables, op.data, op.where, op.whereArgs ...)
    }
}

// 链式操作， CURD - Delete
func (op *DbOp) Delete() (sql.Result, error) {
    if op.where == "" {
        return nil, errors.New("where is required while deleting")
    }
    if op.tx == nil {
        return op.db.Delete(op.tables, op.where, op.whereArgs...)
    } else {
        return op.tx.Delete(op.tables, op.where, op.whereArgs...)
    }
}

// 设置批处理的大小
func (op *DbOp) Batch(batch int) *DbOp {
    op.batch = batch
    return op
}

// 链式操作，select
func (op *DbOp) Select() (List, error) {
    if op.fields == "" {
        op.fields = "*"
    }
    s := fmt.Sprintf("SELECT %s FROM %s", op.fields, op.tables)
    if op.where != "" {
        s += " WHERE " + op.where
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
    if op.tx == nil {
        return op.db.GetAll(s, op.whereArgs...)
    } else {
        return op.tx.GetAll(s, op.whereArgs...)
    }
}

// 链式操作，查询所有记录
func (op *DbOp) All() (List, error) {
    return op.Select()
}

// 链式操作，查询单条记录
func (op *DbOp) One() (Map, error) {
    list, err := op.All()
    if err != nil {
        return nil, err
    }
    if len(list) > 0 {
        return list[0], nil
    }
    return nil, nil
}

// 链式操作，查询字段值
func (op *DbOp) Value() (interface{}, error) {
    one, err := op.One()
    if err != nil {
        return "", err
    }
    for _, v := range one {
        return v, nil
    }
    return "", nil
}

