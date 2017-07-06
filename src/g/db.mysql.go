package g

import (
    "log"
    "fmt"
    "errors"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "strconv"
)

// 数据库全局空对象，用于封装方法
var Db gDb

type gDb struct {}

type gDbTransaction struct {
    db *sql.DB
    tx *sql.Tx
}

type GDb struct {
    Transaction gDbTransaction
    db *sql.DB
}

// 数据库配置
type GDbConfig struct {
    Host string
    Port string
    User string
    Pass string
    Name string
}

// 获得一个数据库操作对象
func (d gDb) New(c GDbConfig) (*GDb) {
    db, err := sql.Open(
        "mysql",
        fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", c.User, c.Pass, c.Host, c.Port, c.Name),
    )
    if err != nil {
        log.Fatal(err)
    }
    return &GDb {
        db : db,
        Transaction: gDbTransaction {
            db : db,
        },
    }
}

// 关闭链接
func (d *GDb) Close() {
    d.db.Close()
}

// 数据库sql查询
func (d *GDb) Query(q string, args ...interface{}) (*sql.Rows, error) {
    rows, err := d.db.Query(q, args ...)
    if (err == nil) {
        return rows, nil
    }
    return nil, err
}

// 执行一条sql，并返回执行情况
func (d *GDb) Exec(q string, args ...interface{}) (sql.Result, error) {
    return d.db.Exec(q, args ...)
}

// 数据库查询，获取查询结果集，以列表结构返回
func (d *GDb) GetAll(q string) ([]map[string]string, error) {
    // 执行sql
    rows, err := d.Query(q)
    if err != nil || rows == nil {
        return nil, err
    }
    // 列名称列表
    columns, err := rows.Columns()
    if err != nil {
        return nil, err
    }
    // 返回结构组装
    values   := make([]sql.RawBytes, len(columns))
    scanArgs := make([]interface{}, len(values))
    var list []map[string]string
    for i := range values {
        scanArgs[i] = &values[i]
    }
    for rows.Next() {
        err = rows.Scan(scanArgs...)
        if err != nil {
            return list, err
        }
        row := make(map[string]string)
        for i, col := range values {
            row[columns[i]] = string(col)
        }
        list = append(list, row)
    }
    return list, nil
}

// 数据库查询，获取查询结果集，以关联数组结构返回
func (d *GDb) GetOne(q string) (map[string]string, error) {
    list, err := d.GetAll(q)
    if err != nil {
        return nil, err
    }
    return list[0], nil
}

// 数据库查询，获取查询字段值
func (d *GDb) GetValue(q string) (string, error) {
    one, err := d.GetOne(q)
    if err != nil {
        return "", err
    }
    for _, v := range one {
        return v, nil
    }
    return "", nil
}

// sql预处理
func (d *GDb) Prepare(q string) (*sql.Stmt, error) {
    return d.db.Prepare(q)
}

// 获取上一次数据库写入产生的自增主键id，另外也可以使用Exec来实现
func (d *GDb) LastInsertId() (int, error) {
    one, err := d.GetOne("SELECT last_insert_id()")
    if err != nil {
        return 0, err
    }
    for _, v := range one {
        return strconv.Atoi(v)
    }
    return 0, nil
}

// ping一下，判断或保持数据库链接
func (d *GDb) Ping() error {
    err := d.db.Ping();
    return err
}

// 设置数据库连接池中空闲链接的大小
func (d *GDb) SetMaxIdleConns(n int) {
    d.db.SetMaxIdleConns(n);
}

// 设置数据库连接池最大打开的链接数量
func (d *GDb) SetMaxOpenConns(n int) {
    d.db.SetMaxOpenConns(n);
}

// 事务操作，开启
func (t *gDbTransaction) Begin() (*sql.Tx, error) {
    tx, err := t.db.Begin()
    t.tx = tx
    return tx, err
}

// 事务操作，提交
func (t *gDbTransaction) Commit() error {
    if t.tx == nil {
        errors.New("transaction not start")
    }
    err := t.tx.Commit()
    return err
}

// 事务操作，回滚
func (t *gDbTransaction) Rollback() error {
    if t.tx == nil {
        errors.New("transaction not start")
    }
    err := t.tx.Rollback()
    return err
}




