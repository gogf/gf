package gkvdb

// 设置是否开启缓存
func (db *DB) SetCache(enabled bool) {
    if enabled {
        db.setCache(1)
    } else {
        db.setCache(0)
    }
}

// 关闭数据库链接
func (db *DB) Close() {
    db.ixfp.Close()
    db.mtfp.Close()
    db.dbfp.Close()
}

// 查询KV数据
func (db *DB) Get(key []byte) []byte {
    if v, ok := db.memt.get(key); ok {
        return v
    }
    return db.get(key)
}

// 设置KV数据
func (db *DB) Set(key []byte, value []byte) error {
    if db.getCache() {
        if err := db.memt.set(key, value); err != nil {
            return err
        }
        return nil
    }
    return db.set(key, value)
}

// 设置KV数据(强制不使用缓存)
func (db *DB) SetWithoutCache(key []byte, value []byte) error {
    return db.set(key, value)
}

// 删除KV数据
func (db *DB) Remove(key []byte) error {
    if db.getCache() {
        if err := db.memt.remove(key); err != nil {
            return err
        }
        return nil
    }
    return db.remove(key)
}

// 删除KV数据(强制不使用缓存)
func (db *DB) RemoveWithoutCache(key []byte) error {
    return db.remove(key)
}

// 打印数据库状态(调试使用)
//func (db *DB) PrintState() {
//    mtblocks := db.mtsp.GetAllBlocks()
//    dbblocks := db.dbsp.GetAllBlocks()
//    fmt.Println("meta pieces:")
//    fmt.Println("       size:", len(mtblocks))
//    fmt.Println("       list:", mtblocks)
//
//    fmt.Println("data pieces:")
//    fmt.Println("       size:", len(dbblocks))
//    fmt.Println("       list:", dbblocks)
//
//    fmt.Println("=======================================")
//}

//// 获取所有的碎片(调试使用)
//func (db *DB) GetBlocks() []gfilespace.Block {
//    return db.mtsp.GetAllBlocks()
//}


