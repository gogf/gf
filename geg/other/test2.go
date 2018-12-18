package main

import "fmt"

type Model struct {
    tablesInit   string        // 初始化Model时的表名称(可以是多个)
    tables       string        // 数据库操作表
    fields       string        // 操作字段
    where        string        // 操作条件
    whereArgs    []interface{} // 操作条件参数
    groupBy      string        // 分组语句
    orderBy      string        // 排序语句
    start        int           // 分页开始
    limit        int           // 分页条数
    data         interface{}   // 操作记录(支持Map/List/string类型)
    batch        int           // 批量操作条数
    filter       bool          // 是否按照表字段过滤data参数
    cacheEnabled bool          // 当前SQL操作是否开启查询缓存功能
    cacheTime    int           // 查询缓存时间
    cacheName    string        // 查询缓存名称
}

func main() {
    m1 := &Model{
        tables : "1",
    }
    m2 := &Model{
        tables : "2",
    }
    *m2 = *m1
    fmt.Println(m1)
    fmt.Println(m2)
}