package gkvdb

import (
    "g/os/gcache"
    "errors"
)


// 查询
func (db *DB) get(key []byte) []byte {
    ckey := "value_cache_" + string(key)
    if v := gcache.Get(ckey); v != nil {
        return v.([]byte)
    }
    db.mu.RLock()
    defer db.mu.RUnlock()

    value, _ := db.getValueByKey(key)
    gcache.Set(ckey, value, gCACHE_DEFAULT_TIMEOUT)
    return value
}

// 保存
func (db *DB) set(key []byte, value []byte) error {
    defer gcache.Remove("value_cache_" + string(key))

    db.mu.Lock()
    defer db.mu.Unlock()

    // 查询索引信息
    record, err := db.getRecordByKey(key)
    if err != nil {
        return err
    }

    oldr := *record

    // 写入数据文件，并更新record信息
    if err := db.insertDataByRecord(key, value, record); err != nil {
        return errors.New("inserting data error: " + err.Error())
    }

    // 根据record信息更新索引文件
    if oldr.mt.start != record.mt.start || oldr.mt.size != record.mt.size {
        if err := db.updateIndexByRecord(record); err != nil {
            return errors.New("creating index error: " + err.Error())
        }
    }
    return nil
}


// 删除
func (db *DB) remove(key []byte) error {
    defer gcache.Remove("value_cache_" + string(key))

    db.mu.Lock()
    defer db.mu.Unlock()

    // 查询索引信息
    record, err := db.getRecordByKey(key)
    if err != nil {
        return err
    }
    // 如果找到匹配才执行删除操作
    if record.mt.match == 0 {
        return db.removeDataByRecord(record)
    }
    return nil
}
