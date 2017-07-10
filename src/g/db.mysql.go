package g

import (
    "log"
    "fmt"
    "errors"
    "strconv"
    "strings"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

// 数据库全局对象，用于封装方法
var Db = gDb {
    // 数据库操作选项
    OPTION_INSERT  : 0,
    OPTION_REPLACE : 1,
    OPTION_SAVE    : 2,
    OPTION_IGNORE  : 3,
}

type gDb struct {
    OPTION_INSERT  uint8
    OPTION_REPLACE uint8
    OPTION_SAVE    uint8
    OPTION_IGNORE  uint8
}

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
    //fmt.Println(q)
    //fmt.Println(args)
    return d.db.Exec(q, args ...)
}

// 数据库查询，获取查询结果集，以列表结构返回
func (d *GDb) GetAll(q string, args ...interface{}) ([]map[string]string, error) {
    // 执行sql
    rows, err := d.Query(q, args ...)
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
func (d *GDb) GetOne(q string, args ...interface{}) (map[string]string, error) {
    list, err := d.GetAll(q, args ...)
    if err != nil {
        return nil, err
    }
    return list[0], nil
}

// 数据库查询，获取查询字段值
func (d *GDb) GetValue(q string, args ...interface{}) (string, error) {
    one, err := d.GetOne(q, args ...)
    if err != nil {
        return "", err
    }
    for _, v := range one {
        return v, nil
    }
    return "", nil
}

// sql预处理，执行完成后调用返回值sql.Stmt.Exec完成sql操作
// 记得调用sql.Stmt.Close关闭操作对象
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

// 根据insert选项获得操作名称
func (d *GDb) getInsertOperationByOption(option uint8) string {
    oper := "INSERT"
    switch option {
    case Db.OPTION_INSERT:
    case Db.OPTION_REPLACE:
        oper   = "REPLACE"
    case Db.OPTION_SAVE:
    case Db.OPTION_IGNORE:
        oper   = "INSERT IGNORE"
    }
    return oper
}

// insert、replace, save， ignore操作
// 0: insert:  仅仅执行写入操作，如果存在冲突的主键或者唯一索引，那么报错返回
// 1: replace: 如果数据存在(主键或者唯一索引)，那么删除后重新写入一条
// 2: save:    如果数据存在(主键或者唯一索引)，那么更新，否则写入一条新数据
// 3: ignore:  如果数据存在(主键或者唯一索引)，那么什么也不做
func (d *GDb) insert(table string, data map[string]string, option uint8) (sql.Result, error) {
    var keys   []string
    var values []string
    var params []interface{}
    for k, v := range data {
        keys   = append(keys,   fmt.Sprintf("`%s`", k))
        values = append(values, "?")
        params = append(params, v)
    }
    operation := d.getInsertOperationByOption(option)
    updatestr := ""
    if option == Db.OPTION_SAVE {
        var updates []string
        for k, _ := range data {
            updates = append(updates, fmt.Sprintf("`%s`=VALUES(%s)", k, k))
        }
        updatestr = fmt.Sprintf(" ON DUPLICATE KEY UPDATE %s", strings.Join(updates, ","))
    }
    return d.Exec(
        fmt.Sprintf("%s INTO `%s`(%s) VALUES(%s) %s",
            operation, table, strings.Join(keys,   ","), strings.Join(values, ","), updatestr), params...
    )
}

// CURD操作:单条数据写入, 仅仅执行写入操作，如果存在冲突的主键或者唯一索引，那么报错返回
func (d *GDb) Insert(table string, data map[string]string) (sql.Result, error) {
    return d.insert(table, data, Db.OPTION_INSERT)
}

// CURD操作:单条数据写入, 如果数据存在(主键或者唯一索引)，那么删除后重新写入一条
func (d *GDb) Replace(table string, data map[string]string) (sql.Result, error) {
    return d.insert(table, data, Db.OPTION_REPLACE)
}

// CURD操作:单条数据写入, 如果数据存在(主键或者唯一索引)，那么更新，否则写入一条新数据
func (d *GDb) Save(table string, data map[string]string) (sql.Result, error) {
    return d.insert(table, data, Db.OPTION_SAVE)
}

// 批量写入数据
func (d *GDb) batchInsert(table string, list []map[string]string, batch int, option uint8) error {
    var keys    []string
    var values  []string
    var bvalues []string
    var params  []interface{}
    var size int = len(list)
    // 判断长度
    if size < 1 {
        return errors.New("empty data list")
    }
    // 首先获取字段名称及记录长度
    for k, _ := range list[0] {
        keys   = append(keys,   k)
        values = append(values, "?")
    }
    var kstr = "`" + strings.Join(keys, "`,`") + "`"
    // 操作判断
    operation := d.getInsertOperationByOption(option)
    updatestr := ""
    if option == Db.OPTION_SAVE {
        var updates []string
        for _, k := range keys {
            updates = append(updates, fmt.Sprintf("`%s`=VALUES(%s)", k, k))
        }
        updatestr = fmt.Sprintf(" ON DUPLICATE KEY UPDATE %s", strings.Join(updates, ","))
    }
    // 构造批量写入数据格式(注意map的遍历是无序的)
    for i := 0; i < size; i++ {
        for _, k := range keys {
            params = append(params, list[i][k])
        }
        bvalues = append(bvalues, "(" + strings.Join(values, ",") + ")")
        if len(bvalues) == batch {
            _, err := d.Exec(fmt.Sprintf("%s INTO `%s`(%s) VALUES%s %s", operation, table, kstr, strings.Join(bvalues, ","), updatestr), params...)
            if err != nil {
                return err
            }
            bvalues = bvalues[:0]
        }
    }
    // 处理最后不构成指定批量的数据
    if (len(bvalues) > 0) {
        _, err := d.Exec(fmt.Sprintf("%s INTO `%s`(%s) VALUES%s %s", operation, table, kstr, strings.Join(bvalues, ","), updatestr), params...)
        if err != nil {
            return err
        }
    }
    return nil
}

// CURD操作:批量数据指定批次量写入
func (d *GDb) BatchInsert(table string, list []map[string]string, batch int) error {
    return d.batchInsert(table, list, batch, Db.OPTION_INSERT)
}

// CURD操作:批量数据指定批次量写入, 如果数据存在(主键或者唯一索引)，那么删除后重新写入一条
func (d *GDb) BatchReplace(table string, list []map[string]string, batch int) error {
    return d.batchInsert(table, list, batch, Db.OPTION_REPLACE)
}

// CURD操作:批量数据指定批次量写入, 如果数据存在(主键或者唯一索引)，那么更新，否则写入一条新数据
func (d *GDb) BatchSave(table string, list []map[string]string, batch int) error {
    return d.batchInsert(table, list, batch, Db.OPTION_SAVE)
}

// CURD操作:数据更新，统一采用sql预处理
func (d *GDb) Update(table string, data interface{}, condition string, args ...interface{}) (sql.Result, error) {
    var params  []interface{}
    var updates string
    switch data.(type) {
    case string:
        updates = data.(string)
    case map[string]string:
        var keys []string
        for k, v := range data.(map[string]string) {
            keys   = append(keys,   fmt.Sprintf("`%s`=?", k))
            params = append(params, v)
        }
        updates = strings.Join(keys,   ",")
    }
    for _, v := range args {
        params = append(params, v.(string))
    }
    return d.Exec(
        fmt.Sprintf("UPDATE `%s` SET %s WHERE %s", table, updates, condition), params...
    )
}

// CURD操作:删除数据
func (d *GDb) Delete(table string, condition string, args ...interface{}) (sql.Result, error) {
    return d.Exec(
        fmt.Sprintf("DELETE FROM `%s` WHERE %s", table, condition), args...
    )
}

