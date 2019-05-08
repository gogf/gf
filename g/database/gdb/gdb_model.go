// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//
// @author john, ymrjqyy

package gdb

import (
    "database/sql"
    "errors"
    "fmt"
    "github.com/gogf/gf/g/util/gconv"
    "reflect"
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
    safe         bool          // 当前模型是否运行安全模式（可修改当前模型，否则每一次链式操作都是返回新的模型对象）
}

// 链式操作，数据表字段，可支持多个表，以半角逗号连接
func (bs *dbBase) Table(tables string) (*Model) {
	return &Model {
		db         : bs.db,
        tablesInit : tables,
		tables     : tables,
		fields     : "*",
        safe       : false,
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
        safe       : false,
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

// 标识当前对象运行安全模式(可被修改)。
// 1. 默认情况下，模型对象的对象属性无法被修改，
// 每一次链式操作都是克隆一个新的模型对象，这样所有的操作都不会污染模型对象。
// 但是链式操作如果需要分开执行，那么需要将新的克隆对象赋值给旧的模型对象继续操作。
// 2. 当标识模型对象为可修改，那么在当前模型对象的所有链式操作均会影响下一次的链式操作，
// 即使是链式操作分开执行。
// 3. 大部分ORM框架默认模型对象是可修改的，但是GF框架的ORM提供给开发者更灵活，更安全的链式操作选项。
func (md *Model) Safe(safe...bool) *Model {
    if len(safe) > 0 {
        md.safe = safe[0]
    } else {
        md.safe = true
    }
    return md
}

// 返回操作的模型对象，可能是当前对象，也可能是新的克隆对象，根据alterable决定。
func (md *Model) getModel() *Model {
    if !md.safe {
        return md
    } else {
        return md.Clone()
    }
}

// 链式操作，左联表
func (md *Model) LeftJoin(joinTable string, on string) (*Model) {
    model        := md.getModel()
    model.tables += fmt.Sprintf(" LEFT JOIN %s ON (%s)", joinTable, on)
	return model
}

// 链式操作，右联表
func (md *Model) RightJoin(joinTable string, on string) (*Model) {
    model        := md.getModel()
    model.tables += fmt.Sprintf(" RIGHT JOIN %s ON (%s)", joinTable, on)
	return model
}

// 链式操作，内联表
func (md *Model) InnerJoin(joinTable string, on string) (*Model) {
    model        := md.getModel()
    model.tables += fmt.Sprintf(" INNER JOIN %s ON (%s)", joinTable, on)
	return model
}

// 链式操作，查询字段
func (md *Model) Fields(fields string) (*Model) {
    model       := md.getModel()
    model.fields = fields
	return model
}

// 链式操作，过滤字段
func (md *Model) Filter() (*Model) {
    model       := md.getModel()
    model.filter = true
    return model
}

// 链式操作，condition，支持string & gdb.Map.
// 注意，多个Where调用时，会自动转换为And条件调用。
func (md *Model) Where(where interface{}, args ...interface{}) (*Model) {
    model := md.getModel()
    if model.where != "" {
        return md.And(where, args...)
    }
    newWhere, newArgs := formatCondition(where, args)
    model.where        = newWhere
    model.whereArgs    = newArgs
	return model
}

// 链式操作，添加AND条件到Where中
func (md *Model) And(where interface{}, args ...interface{}) (*Model) {
    model             := md.getModel()
    newWhere, newArgs := formatCondition(where, args)
    if len(model.where) > 0 && model.where[0] == '(' {
        model.where = fmt.Sprintf(`%s AND (%s)`, model.where, newWhere)
    } else {
        model.where = fmt.Sprintf(`(%s) AND (%s)`, model.where, newWhere)
    }
    model.whereArgs = append(model.whereArgs, newArgs...)
	return model
}

// 链式操作，添加OR条件到Where中
func (md *Model) Or(where interface{}, args ...interface{}) (*Model) {
    model             := md.getModel()
    newWhere, newArgs := formatCondition(where, args)
    if len(model.where) > 0 && model.where[0] == '(' {
        model.where = fmt.Sprintf(`%s OR (%s)`, model.where, newWhere)
    } else {
        model.where = fmt.Sprintf(`(%s) OR (%s)`, model.where, newWhere)
    }
    model.whereArgs = append(model.whereArgs, newArgs...)
	return model
}

// 链式操作，group by
func (md *Model) GroupBy(groupBy string) (*Model) {
    model        := md.getModel()
    model.groupBy = groupBy
	return model
}

// 链式操作，order by
func (md *Model) OrderBy(orderBy string) (*Model) {
    model        := md.getModel()
    model.orderBy = orderBy
	return model
}

// 链式操作，limit
func (md *Model) Limit(start int, limit int) (*Model) {
    model      := md.getModel()
    model.start = start
    model.limit = limit
	return model
}

// 链式操作，翻页，注意分页页码从1开始，而Limit方法从0开始。
func (md *Model) ForPage(page, limit int) (*Model) {
    model      := md.getModel()
    model.start = (page - 1) * limit
    model.limit = limit
	return model
}

// 设置批处理的大小
func (md *Model) Batch(batch int) *Model {
    model      := md.getModel()
    model.batch = batch
    return model
}

// 查询缓存/清除缓存操作，需要注意的是，事务查询不支持缓存。
// 当time < 0时表示清除缓存， time=0时表示不过期, time > 0时表示过期时间，time过期时间单位：秒；
// name表示自定义的缓存名称，便于业务层精准定位缓存项(如果业务层需要手动清理时，必须指定缓存名称)，
// 例如：查询缓存时设置名称，清理缓存时可以给定清理的缓存名称进行精准清理。
func (md *Model) Cache(time int, name ... string) *Model {
    model          := md.getModel()
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

// 链式操作，操作数据项，参数data类型支持 string/map/slice/struct/*struct ,
// 也可以是：key,value,key,value,...。
func (md *Model) Data(data ...interface{}) *Model {
    model := md.getModel()
	if len(data) > 1 {
		m := make(map[string]interface{})
		for i := 0; i < len(data); i += 2 {
			m[gconv.String(data[i])] = data[i + 1]
		}
        model.data = m
	} else {
		switch params := data[0].(type) {
			case Result:
				model.data = params.ToList()
			case Record:
				model.data = params.ToMap()
			case List:
                model.data = params
			case Map:
                model.data = params
			default:
                rv   := reflect.ValueOf(params)
                kind := rv.Kind()
                if kind == reflect.Ptr {
                    rv   = rv.Elem()
                    kind = rv.Kind()
                }
                switch kind {
                	// 如果是slice，那么转换为List类型
                    case reflect.Slice: fallthrough
                    case reflect.Array:
                        list := make(List, rv.Len())
                        for i := 0; i < rv.Len(); i++ {
                            list[i] = gconv.Map(rv.Index(i).Interface())
                        }
                        model.data = list
                    case reflect.Map:   fallthrough
                    case reflect.Struct:
                        model.data = Map(gconv.Map(data[0]))
                    default:
                        model.data = data[0]
                }
		}
	}
	return model
}

// 链式操作， CURD - Insert/BatchInsert。
// 根据Data方法传递的参数类型决定该操作是单条操作还是批量操作，
// 如果Data方法传递的是slice类型，那么为批量操作。
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

// 链式操作， CURD - Replace/BatchReplace。
// 根据Data方法传递的参数类型决定该操作是单条操作还是批量操作，
// 如果Data方法传递的是slice类型，那么为批量操作。
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

// 链式操作， CURD - Save/BatchSave。
// 根据Data方法传递的参数类型决定该操作是单条操作还是批量操作，
// 如果Data方法传递的是slice类型，那么为批量操作。
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
		batch := gDEFAULT_BATCH_NUM
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
		return md.db.doUpdate(nil, md.tables, md.data, md.where, md.whereArgs ...)
	} else {
		return md.tx.doUpdate(md.tables, md.data, md.where, md.whereArgs ...)
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
		return md.db.doDelete(nil, md.tables, md.where, md.whereArgs...)
	} else {
		return md.tx.doDelete(md.tables, md.where, md.whereArgs...)
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

// 链式操作，查询单条记录，并自动转换为struct对象, 参数必须为对象的指针，不能为空指针。
func (md *Model) Struct(objPointer interface{}) error {
	one, err := md.One()
	if err != nil {
		return err
	}
	return one.ToStruct(objPointer)
}

// 链式操作，查询多条记录，并自动转换为指定的slice对象, 如: []struct/[]*struct。
func (md *Model) Structs(objPointerSlice interface{}) error {
	r, err := md.All()
	if err != nil {
		return err
	}
	return r.ToStructs(objPointerSlice)
}

// 链式操作，将结果转换为指定的struct/*struct/[]struct/[]*struct,
// 参数应该为指针类型，否则返回失败。
// 该方法自动识别参数类型，调用Struct/Structs方法。
func (md *Model) Scan(objPointer interface{}) error {
    t := reflect.TypeOf(objPointer)
    k := t.Kind()
    if k != reflect.Ptr {
        return fmt.Errorf("params should be type of pointer, but got: %v", k)
    }
    k = t.Elem().Kind()
    switch k {
        case reflect.Array:
        case reflect.Slice:
            return md.Structs(objPointer)
        case reflect.Struct:
            return md.Struct(objPointer)
        default:
            return fmt.Errorf("element type should be type of struct/slice, unsupported: %v", k)
    }
    return nil
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

// 组块结果集。
func (md *Model) Chunk(limit int, callback func(result Result, err error) bool) {
	page  := 1
	model := md
	for {
        model      = model.ForPage(page, limit)
		data, err := model.All()
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
