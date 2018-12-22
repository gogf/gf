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
	"gitee.com/johng/gf/g/util/gconv"
	_ "gitee.com/johng/gf/third/github.com/go-sql-driver/mysql"
    "reflect"
    "strings"
)

// 数据库链式操作模型对象
type Model struct {
	db           DB            // 数据库操作对象
	tx           *TX           // 数据库事务对象
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

// 链式操作，数据表字段，可支持多个表，以半角逗号连接
func (bs *dbBase) Table(tables string) (*Model) {
	return &Model {
		db         : bs.db,
        tablesInit : tables,
		tables     : tables,
		fields     : "*",
	}
}

// 链式操作，数据表字段，可支持多个表，以半角逗号连接
func (bs *dbBase) From(tables string) (*Model) {
	return bs.db.Table(tables)
}

// (事务)链式操作，数据表字段，可支持多个表，以半角逗号连接
func (tx *TX) Table(tables string) (*Model) {
	return &Model{
		db         : tx.db,
		tx         : tx,
        tablesInit : tables,
		tables     : tables,
	}
}

// (事务)链式操作，数据表字段，可支持多个表，以半角逗号连接
func (tx *TX) From(tables string) (*Model) {
	return tx.Table(tables)
}

// 克隆一个当前对象
func (md *Model) Clone() *Model {
    newModel := (*Model)(nil)
	if md.tx != nil {
        newModel = md.tx.Table(md.tablesInit)
	} else {
        newModel = md.db.Table(md.tablesInit)
	}
    *newModel = *md
    return newModel
}

// 链式操作，左联表
func (md *Model) LeftJoin(joinTable string, on string) (*Model) {
    model        := md.Clone()
    model.tables += fmt.Sprintf(" LEFT JOIN %s ON (%s)", joinTable, on)
	return model
}

// 链式操作，右联表
func (md *Model) RightJoin(joinTable string, on string) (*Model) {
    model        := md.Clone()
    model.tables += fmt.Sprintf(" RIGHT JOIN %s ON (%s)", joinTable, on)
	return model
}

// 链式操作，内联表
func (md *Model) InnerJoin(joinTable string, on string) (*Model) {
    model        := md.Clone()
    model.tables += fmt.Sprintf(" INNER JOIN %s ON (%s)", joinTable, on)
	return model
}

// 链式操作，查询字段
func (md *Model) Fields(fields string) (*Model) {
    model       := md.Clone()
    model.fields = fields
	return model
}

// 链式操作，过滤字段
func (md *Model) Filter() (*Model) {
    model       := md.Clone()
    model.filter = true
    return model
}

// 链式操作，condition，支持string & gdb.Map
func (md *Model) Where(where interface{}, args ...interface{}) (*Model) {
    model             := md.Clone()
    newWhere, newArgs := formatCondition(where, args)
    model.where        = newWhere
    model.whereArgs    = append(model.whereArgs, newArgs...)
	// 支持 Where("uid", 1)这种格式
	if len(args) == 1 && strings.Index(model.where , "?") < 0 {
        model.where += "=?"
    }
	return model
}

// 链式操作，添加AND条件到Where中
func (md *Model) And(where interface{}, args ...interface{}) (*Model) {
    model             := md.Clone()
    newWhere, newArgs := formatCondition(where, args)
    model.where       += " AND " + newWhere
    model.whereArgs    = append(model.whereArgs, newArgs...)
	return model
}

// 链式操作，添加OR条件到Where中
func (md *Model) Or(where interface{}, args ...interface{}) (*Model) {
    model             := md.Clone()
    newWhere, newArgs := formatCondition(where, args)
    model.where       += " OR " + newWhere
    model.whereArgs    = append(model.whereArgs, newArgs...)
	return model
}

// 链式操作，group by
func (md *Model) GroupBy(groupBy string) (*Model) {
    model        := md.Clone()
    model.groupBy = groupBy
	return model
}

// 链式操作，order by
func (md *Model) OrderBy(orderBy string) (*Model) {
    model        := md.Clone()
    model.orderBy = orderBy
	return model
}

// 链式操作，limit
func (md *Model) Limit(start int, limit int) (*Model) {
    model      := md.Clone()
    model.start = start
    model.limit = limit
	return model
}

// 链式操作，翻页
// @author ymrjqyy
func (md *Model) ForPage(page, limit int) (*Model) {
    model      := md.Clone()
    model.start = (page - 1) * limit
    model.limit = limit
	return model
}

// 设置批处理的大小
func (md *Model) Batch(batch int) *Model {
    model      := md.Clone()
    model.batch = batch
    return model
}

// 查询缓存/清除缓存操作，需要注意的是，事务查询不支持缓存。
// 当time < 0时表示清除缓存， time=0时表示不过期, time > 0时表示过期时间，time过期时间单位：秒；
// name表示自定义的缓存名称，便于业务层精准定位缓存项(如果业务层需要手动清理时，必须指定缓存名称)，
// 例如：查询缓存时设置名称，清理缓存时可以给定清理的缓存名称进行精准清理。
func (md *Model) Cache(time int, name ... string) *Model {
    model          := md.Clone()
    model.cacheTime = time
    if len(name) > 0 {
        model.cacheName = name[0]
    }
    // 查询缓存特性不支持事务操作
    if model.tx == nil {
        model.cacheEnabled = true
    }
    return model
}

// 链式操作，操作数据记录项，可以是string/Map, 也可以是：key,value,key,value,...
func (md *Model) Data(data ...interface{}) (*Model) {
    model := md.Clone()
	if len(data) > 1 {
		m := make(map[string]interface{})
		for i := 0; i < len(data); i += 2 {
			m[gconv.String(data[i])] = data[i+1]
		}
        model.data = m
	} else {
		switch data[0].(type) {
			case List:
                model.data = data[0]
			case Map:
                model.data = data[0]
			default:
                rv   := reflect.ValueOf(data[0])
                kind := rv.Kind()
                if kind == reflect.Ptr {
                    rv   = rv.Elem()
                    kind = rv.Kind()
                }
                switch kind {
                    case reflect.Slice: fallthrough
                    case reflect.Array:
                        list := make(List, rv.Len())
                        for i := 0; i < rv.Len(); i++ {
                            list[i] = gconv.Map(rv.Index(i).Interface())
                        }
                        model.data = list
                    case reflect.Map:
                        model.data = gconv.Map(data[0])
                    default:
                        model.data = data[0]
                }
		}
	}
	return model
}

// 链式操作， CURD - Insert/BatchInsert
func (md *Model) Insert() (result sql.Result, err error) {
	defer func() {
		if err == nil {
			md.checkAndRemoveCache()
		}
	}()
	if md.data == nil {
		return nil, errors.New("inserting into table with empty data")
	}
	// 批量操作
	if list, ok := md.data.(List); ok {
		batch := 10
		if md.batch > 0 {
			batch = md.batch
		}
		if md.filter {
		    for k, m := range list {
                list[k] = md.db.filterFields(md.tables, m)
            }
        }
		if md.tx == nil {
			return md.db.BatchInsert(md.tables, list, batch)
		} else {
			return md.tx.BatchInsert(md.tables, list, batch)
		}
	} else if data, ok := md.data.(Map); ok {
        if md.filter {
            data = md.db.filterFields(md.tables, data)
        }
		if md.tx == nil {
			return md.db.Insert(md.tables, data)
		} else {
			return md.tx.Insert(md.tables, data)
		}
	}
	return nil, errors.New("inserting into table with invalid data type")
}

// 链式操作， CURD - Replace/BatchReplace
func (md *Model) Replace() (result sql.Result, err error) {
	defer func() {
		if err == nil {
			md.checkAndRemoveCache()
		}
	}()
	if md.data == nil {
		return nil, errors.New("replacing into table with empty data")
	}
	// 批量操作
	if list, ok := md.data.(List); ok {
		batch := 10
		if md.batch > 0 {
			batch = md.batch
		}
        if md.filter {
            for k, m := range list {
                list[k] = md.db.filterFields(md.tables, m)
            }
        }
		if md.tx == nil {
			return md.db.BatchReplace(md.tables, list, batch)
		} else {
			return md.tx.BatchReplace(md.tables, list, batch)
		}
	} else if data, ok := md.data.(Map); ok {
        if md.filter {
            data = md.db.filterFields(md.tables, data)
        }
		if md.tx == nil {
			return md.db.Replace(md.tables, data)
		} else {
			return md.tx.Replace(md.tables, data)
		}
	}
	return nil, errors.New("replacing into table with invalid data type")
}

// 链式操作， CURD - Save/BatchSave
func (md *Model) Save() (result sql.Result, err error) {
	defer func() {
		if err == nil {
			md.checkAndRemoveCache()
		}
	}()
	if md.data == nil {
		return nil, errors.New("replacing into table with empty data")
	}
	// 批量操作
	if list, ok := md.data.(List); ok {
		batch := 10
		if md.batch > 0 {
			batch = md.batch
		}
        if md.filter {
            for k, m := range list {
                list[k] = md.db.filterFields(md.tables, m)
            }
        }
		if md.tx == nil {
			return md.db.BatchSave(md.tables, list, batch)
		} else {
			return md.tx.BatchSave(md.tables, list, batch)
		}
	} else if data, ok := md.data.(Map); ok {
        if md.filter {
            data = md.db.filterFields(md.tables, data)
        }
		if md.tx == nil {
			return md.db.Save(md.tables, data)
		} else {
			return md.tx.Save(md.tables, data)
		}
	}
	return nil, errors.New("saving into table with invalid data type")
}

// 链式操作， CURD - Update
func (md *Model) Update() (result sql.Result, err error) {
	defer func() {
		if err == nil {
			md.checkAndRemoveCache()
		}
	}()
	if md.data == nil {
		return nil, errors.New("updating table with empty data")
	}
    if md.filter {
        if data, ok := md.data.(Map); ok {
            if md.filter {
                md.data = md.db.filterFields(md.tables, data)
            }
        }
    }
	if md.tx == nil {
		return md.db.Update(md.tables, md.data, md.where, md.whereArgs ...)
	} else {
		return md.tx.Update(md.tables, md.data, md.where, md.whereArgs ...)
	}
}

// 链式操作， CURD - Delete
func (md *Model) Delete() (result sql.Result, err error) {
	defer func() {
		if err == nil {
			md.checkAndRemoveCache()
		}
	}()
	if md.tx == nil {
		return md.db.Delete(md.tables, md.where, md.whereArgs...)
	} else {
		return md.tx.Delete(md.tables, md.where, md.whereArgs...)
	}
}

// 链式操作，select
func (md *Model) Select() (Result, error) {
	return md.All()
}

// 链式操作，查询所有记录
func (md *Model) All() (Result, error) {
	return md.getAll(md.getFormattedSql(), md.whereArgs...)
}

// 链式操作，查询单条记录
func (md *Model) One() (Record, error) {
	list, err := md.All()
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		return list[0], nil
	}
	return nil, nil
}

// 链式操作，查询字段值
func (md *Model) Value() (Value, error) {
	one, err := md.One()
	if err != nil {
		return nil, err
	}
	for _, v := range one {
		return v, nil
	}
	return nil, nil
}

// 链式操作，查询单条记录，并自动转换为struct对象
func (md *Model) Struct(obj interface{}) error {
	one, err := md.One()
	if err != nil {
		return err
	}
	return one.ToStruct(obj)
}

// 链式操作，查询数量，fields可以为空，也可以自定义查询字段，
// 当给定自定义查询字段时，该字段必须为数量结果，否则会引起歧义，使用如：md.Fields("COUNT(id)")
func (md *Model) Count() (int, error) {
    defer func(fields string) {
        md.fields = fields
    }(md.fields)
	if md.fields == "" || md.fields == "*" {
		md.fields = "COUNT(1)"
	} else {
        md.fields = fmt.Sprintf(`COUNT(%s)`, md.fields)
	}
	s := md.getFormattedSql()
	if len(md.groupBy) > 0 {
		s = fmt.Sprintf("SELECT COUNT(1) FROM (%s) count_alias", s)
	}
	list, err := md.getAll(s, md.whereArgs...)
	if err != nil {
		return 0, err
	}
	if len(list) > 0 {
		for _, v := range list[0] {
			return v.Int(), nil
		}
	}
	return 0, nil
}

// 查询操作，对底层SQL操作的封装
func (md *Model) getAll(query string, args ...interface{}) (result Result, err error) {
	cacheKey := ""
	// 查询缓存查询处理
	if md.cacheEnabled {
		cacheKey = md.cacheName
		if len(cacheKey) == 0 {
			cacheKey = query + "/" + gconv.String(args)
		}
		if v := md.db.getCache().Get(cacheKey); v != nil {
			return v.(Result), nil
		}
	}

	if md.tx == nil {
		result, err = md.db.GetAll(query, args...)
	} else {
		result, err = md.tx.GetAll(query, args...)
	}
	// 查询缓存保存处理
	if len(cacheKey) > 0 && err == nil {
		if md.cacheTime < 0 {
			md.db.getCache().Remove(cacheKey)
		} else {
			md.db.getCache().Set(cacheKey, result, md.cacheTime*1000)
		}
	}
	return result, err
}

// 检查是否需要查询查询缓存
func (md *Model) checkAndRemoveCache() {
	if md.cacheEnabled && md.cacheTime < 0 && len(md.cacheName) > 0 {
		md.db.getCache().Remove(md.cacheName)
	}
}

// 格式化当前输入参数，返回可执行的SQL语句（不带参数）
func (md *Model) getFormattedSql() string {
	if md.fields == "" {
		md.fields = "*"
	}
	s := fmt.Sprintf("SELECT %s FROM %s", md.fields, md.tables)
	if md.where != "" {
		s += " WHERE " + md.where
	}
	if md.groupBy != "" {
		s += " GROUP BY " + md.groupBy
	}
	if md.orderBy != "" {
		s += " ORDER BY " + md.orderBy
	}
	if md.limit != 0 {
		s += fmt.Sprintf(" LIMIT %d, %d", md.start, md.limit)
	}
	return s
}

// 组块结果集
// @author ymrjqyy
// @author 2018-08-15
func (md *Model) Chunk(limit int, callback func(result Result, err error) bool) {
	page := 1
	for {
		md.ForPage(page, limit)
		data, err := md.getAll(md.getFormattedSql(), md.whereArgs...)
		if err != nil {
			callback(nil, err)
			break
		}
		if len(data) == 0 {
			break
		}
		if callback(data, err) == false {
			break
		}
		if len(data) < limit {
			break
		}
		page++
	}
}
