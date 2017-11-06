package gkvdb

import (
    "errors"
    "strconv"
    "fmt"
)

// 关闭数据库链接
func (db *DB) Close() {
    db.ixfp.Close()
    db.mtfp.Close()
    db.dbfp.Close()
}

// 查询KV数据
func (db *DB) Get(key []byte) []byte {
    value, _ := db.getValueByKey(key)
    return value
}


// 设置KV数据
func (db *DB) Set(key []byte, value []byte) error {
    if len(key) > gMAX_KEY_SIZE {
        return errors.New("too large key size, max allowed: " + strconv.Itoa(gMAX_KEY_SIZE) + " bytes")
    }
    if len(value) > gMAX_VALUE_SIZE {
        return errors.New("too large value size, max allowed: " + strconv.Itoa(gMAX_VALUE_SIZE) + " bytes")
    }

    // 查询索引信息
    record, err := db.getRecordByKey(key)
    if err != nil {
        return err
    }
    //oldcap := record.mt.cap
    // 写入数据文件，并更新record信息
    if err := db.insertDataByRecord(key, value, record); err != nil {
        return errors.New("inserting data error: " + err.Error())
    }
    //if oldcap > 0 && record.mt.cap > oldcap {
    //    fmt.Printf("new cap %d VS %d\n", record.mt.cap, oldcap)
    //}
    // 根据record信息更新索引文件
    if err := db.updateIndexByRecord(record); err != nil {
        return errors.New("creating index error: " + err.Error())
    }
    return nil
}

// 删除KV数据
func (db *DB) Remove(key []byte) error {
    // 查询索引信息
    record, err := db.getRecordByKey(key)
    if err != nil {
        return err
    }
    // 如果找到匹配才执行删除操作
    if record.mt.match {
        return db.removeDataByRecord(record)
    }
    return nil
}

// 打印数据库状态(调试使用)
func (db *DB) PrintState() {
    mtbysize  := db.mtsp.GetAllBlocksBySize()
    mtbyindex := db.mtsp.GetAllBlocksByIndex()
    dbbysize  := db.dbsp.GetAllBlocksBySize()
    dbbyindex := db.dbsp.GetAllBlocksByIndex()
    fmt.Println("meta pieces:")
    fmt.Println("   by index:", len(mtbyindex))
    //fmt.Println("      list:", mtbyindex)
    fmt.Println("    by size:", len(mtbysize))
    //fmt.Println("      list:", mtbysize)

    fmt.Println("data pieces:")
    fmt.Println("   by index:", len(dbbyindex))
    //fmt.Println("      list:", dbbyindex)
    fmt.Println("    by size:", len(dbbysize))
    //fmt.Println("      list:", dbbysize)
}
