package gdb

import (
    "fmt"
    "database/sql"
    "errors"
    _ "github.com/go-sql-driver/mysql"
)

// gf的数据库操作支持普通方法操作及链式操作两种方式，本文件是链式操作的封装，提供非常简便的CURD方法


// 数据库链式操作对象
type gLinkOp struct {
    link          Link
    tables        string
    fields        string
    condition     string
    conditionArgs []interface{}
    groupby       string
    orderby       string
    start         int
    limit         int
    data          interface{}
    dataList      *DataList
    batch         int
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

// 链式操作，操作数据记录项列表
func (op *gLinkOp) DataList(list *DataList) (*gLinkOp) {
    op.dataList = list
    return op
}

// 链式操作， CURD - Insert
func (op *gLinkOp) Insert() (sql.Result, error) {
    if op.data == nil {
        return nil, errors.New("inserting into table with empty data")
    }
    if d, ok :=  op.data.(*DataMap); ok {
        return op.link.Insert(op.tables, d)
    }
    return nil, errors.New("inserting into table with invalid data type")
}

// 链式操作， CURD - Replace
func (op *gLinkOp) Replace() (sql.Result, error) {
    if op.data == nil {
        return nil, errors.New("replacing into table with empty data")
    }
    if d, ok :=  op.data.(*DataMap); ok {
        return op.link.Insert(op.tables, d)
    }
    return nil, errors.New("replacing into table with invalid data type")
}

// 链式操作， CURD - Save
func (op *gLinkOp) Save() (sql.Result, error) {
    if op.data == nil {
        return nil, errors.New("saving into table with empty data")
    }
    if d, ok :=  op.data.(*DataMap); ok {
        return op.link.Insert(op.tables, d)
    }
    return nil, errors.New("saving into table with invalid data type")
}

// 设置批处理的大小
func (op *gLinkOp) Batch(batch int) *gLinkOp {
    op.batch = batch
    return op
}

// 链式操作， CURD - BatchInsert
func (op *gLinkOp) BatchInsert() error {
    if op.dataList == nil || len(*op.dataList) < 1 {
        return errors.New("batch inserting into table with empty data list")
    }
    batch := 10
    if op.batch > 0 {
        batch = op.batch
    }
    return op.link.BatchInsert(op.tables, op.dataList, batch)
}

// 链式操作， CURD - BatchReplace
func (op *gLinkOp) BatchReplace() error {
    if op.dataList == nil || len(*op.dataList) < 1 {
        return errors.New("batch replacing into table with empty data list")
    }
    batch := 10
    if op.batch > 0 {
        batch = op.batch
    }
    return op.link.BatchReplace(op.tables, op.dataList, batch)
}

// 链式操作， CURD - BatchSave
func (op *gLinkOp) BatchSave() error {
    if op.dataList == nil || len(*op.dataList) < 1 {
        return errors.New("batch saving into table with empty data list")
    }
    batch := 10
    if op.batch > 0 {
        batch = op.batch
    }
    return op.link.BatchSave(op.tables, op.dataList, batch)
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

// 链式操作，select
func (op *gLinkOp) Select() (*DataList, error) {
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
func (op *gLinkOp) All() (*DataList, error) {
    return op.Select()
}

// 链式操作，查询单条记录
func (op *gLinkOp) One() (*DataMap, error) {
    list, err := op.All()
    if err != nil {
        return nil, err
    }
    return &(*list)[0], nil
}

// 链式操作，查询字段值
func (op *gLinkOp) Value() (string, error) {
    one, err := op.One()
    if err != nil {
        return "", err
    }
    for _, v := range *one {
        return v, nil
    }
    return "", nil
}

