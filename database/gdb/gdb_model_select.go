// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"fmt"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/utils"
	"reflect"

	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

// All does "SELECT FROM ..." statement for the model.
// It retrieves the records from table and returns the result as slice type.
// It returns nil if there's no record retrieved with the given conditions from table.
//
// The optional parameter `where` is the same as the parameter of Model.Where function,
// see Model.Where.
func (m *Model) All(where ...interface{}) (Result, error) {
	return m.doGetAll(false, where...)
}

// doGetAll does "SELECT FROM ..." statement for the model.
// It retrieves the records from table and returns the result as slice type.
// It returns nil if there's no record retrieved with the given conditions from table.
//
// The parameter `limit1` specifies whether limits querying only one record if m.limit is not set.
// The optional parameter `where` is the same as the parameter of Model.Where function,
// see Model.Where.
func (m *Model) doGetAll(limit1 bool, where ...interface{}) (Result, error) {
	if len(where) > 0 {
		return m.Where(where[0], where[1:]...).All()
	}
	sqlWithHolder, holderArgs := m.getFormattedSqlAndArgs(queryTypeNormal, limit1)
	return m.doGetAllBySql(sqlWithHolder, holderArgs...)
}

// getFieldsFiltered checks the fields and fieldsEx attributes, filters and returns the fields that will
// really be committed to underlying database driver.
func (m *Model) getFieldsFiltered() string {
	if m.fieldsEx == "" {
		// No filtering.
		if !gstr.Contains(m.fields, ".") && !gstr.Contains(m.fields, " ") {
			return m.db.GetCore().QuoteString(m.fields)
		}
		return m.fields
	}
	var (
		fieldsArray []string
		fieldsExSet = gset.NewStrSetFrom(gstr.SplitAndTrim(m.fieldsEx, ","))
	)
	if m.fields != "*" {
		// Filter custom fields with fieldEx.
		fieldsArray = make([]string, 0, 8)
		for _, v := range gstr.SplitAndTrim(m.fields, ",") {
			fieldsArray = append(fieldsArray, v[gstr.PosR(v, "-")+1:])
		}
	} else {
		if gstr.Contains(m.tables, " ") {
			panic("function FieldsEx supports only single table operations")
		}
		// Filter table fields with fieldEx.
		tableFields, err := m.TableFields(m.tablesInit)
		if err != nil {
			panic(err)
		}
		if len(tableFields) == 0 {
			panic(fmt.Sprintf(`empty table fields for table "%s"`, m.tables))
		}
		fieldsArray = make([]string, len(tableFields))
		for k, v := range tableFields {
			fieldsArray[v.Index] = k
		}
	}
	newFields := ""
	for _, k := range fieldsArray {
		if fieldsExSet.Contains(k) {
			continue
		}
		if len(newFields) > 0 {
			newFields += ","
		}
		newFields += m.db.GetCore().QuoteWord(k)
	}
	return newFields
}

// Chunk iterates the query result with given `size` and `handler` function.
func (m *Model) Chunk(size int, handler ChunkHandler) {
	page := m.start
	if page <= 0 {
		page = 1
	}
	model := m
	for {
		model = model.Page(page, size)
		data, err := model.All()
		if err != nil {
			handler(nil, err)
			break
		}
		if len(data) == 0 {
			break
		}
		if handler(data, err) == false {
			break
		}
		if len(data) < size {
			break
		}
		page++
	}
}

// One retrieves one record from table and returns the result as map type.
// It returns nil if there's no record retrieved with the given conditions from table.
//
// The optional parameter `where` is the same as the parameter of Model.Where function,
// see Model.Where.
func (m *Model) One(where ...interface{}) (Record, error) {
	if len(where) > 0 {
		return m.Where(where[0], where[1:]...).One()
	}
	all, err := m.doGetAll(true)
	if err != nil {
		return nil, err
	}
	if len(all) > 0 {
		return all[0], nil
	}
	return nil, nil
}

// Value retrieves a specified record value from table and returns the result as interface type.
// It returns nil if there's no record found with the given conditions from table.
//
// If the optional parameter `fieldsAndWhere` is given, the fieldsAndWhere[0] is the selected fields
// and fieldsAndWhere[1:] is treated as where condition fields.
// Also see Model.Fields and Model.Where functions.
func (m *Model) Value(fieldsAndWhere ...interface{}) (Value, error) {
	if len(fieldsAndWhere) > 0 {
		if len(fieldsAndWhere) > 2 {
			return m.Fields(gconv.String(fieldsAndWhere[0])).Where(fieldsAndWhere[1], fieldsAndWhere[2:]...).Value()
		} else if len(fieldsAndWhere) == 2 {
			return m.Fields(gconv.String(fieldsAndWhere[0])).Where(fieldsAndWhere[1]).Value()
		} else {
			return m.Fields(gconv.String(fieldsAndWhere[0])).Value()
		}
	}
	one, err := m.One()
	if err != nil {
		return gvar.New(nil), err
	}
	for _, v := range one {
		return v, nil
	}
	return gvar.New(nil), nil
}

// Array queries and returns data values as slice from database.
// Note that if there are multiple columns in the result, it returns just one column values randomly.
//
// If the optional parameter `fieldsAndWhere` is given, the fieldsAndWhere[0] is the selected fields
// and fieldsAndWhere[1:] is treated as where condition fields.
// Also see Model.Fields and Model.Where functions.
func (m *Model) Array(fieldsAndWhere ...interface{}) ([]Value, error) {
	if len(fieldsAndWhere) > 0 {
		if len(fieldsAndWhere) > 2 {
			return m.Fields(gconv.String(fieldsAndWhere[0])).Where(fieldsAndWhere[1], fieldsAndWhere[2:]...).Array()
		} else if len(fieldsAndWhere) == 2 {
			return m.Fields(gconv.String(fieldsAndWhere[0])).Where(fieldsAndWhere[1]).Array()
		} else {
			return m.Fields(gconv.String(fieldsAndWhere[0])).Array()
		}
	}
	all, err := m.All()
	if err != nil {
		return nil, err
	}
	return all.Array(), nil
}

// Struct retrieves one record from table and converts it into given struct.
// The parameter `pointer` should be type of *struct/**struct. If type **struct is given,
// it can create the struct internally during converting.
//
// The optional parameter `where` is the same as the parameter of Model.Where function,
// see Model.Where.
//
// Note that it returns sql.ErrNoRows if the given parameter `pointer` pointed to a variable that has
// default value and there's no record retrieved with the given conditions from table.
//
// Example:
// user := new(User)
// err  := db.Model("user").Where("id", 1).Scan(user)
//
// user := (*User)(nil)
// err  := db.Model("user").Where("id", 1).Scan(&user)
func (m *Model) doStruct(pointer interface{}, where ...interface{}) error {
	model := m
	// Auto selecting fields by struct attributes.
	if model.fieldsEx == "" && (model.fields == "" || model.fields == "*") {
		if v, ok := pointer.(reflect.Value); ok {
			model = m.Fields(v.Interface())
		} else {
			model = m.Fields(pointer)
		}
	}
	one, err := model.One(where...)
	if err != nil {
		return err
	}
	if err = one.Struct(pointer); err != nil {
		return err
	}
	return model.doWithScanStruct(pointer)
}

// Structs retrieves records from table and converts them into given struct slice.
// The parameter `pointer` should be type of *[]struct/*[]*struct. It can create and fill the struct
// slice internally during converting.
//
// The optional parameter `where` is the same as the parameter of Model.Where function,
// see Model.Where.
//
// Note that it returns sql.ErrNoRows if the given parameter `pointer` pointed to a variable that has
// default value and there's no record retrieved with the given conditions from table.
//
// Example:
// users := ([]User)(nil)
// err   := db.Model("user").Scan(&users)
//
// users := ([]*User)(nil)
// err   := db.Model("user").Scan(&users)
func (m *Model) doStructs(pointer interface{}, where ...interface{}) error {
	model := m
	// Auto selecting fields by struct attributes.
	if model.fieldsEx == "" && (model.fields == "" || model.fields == "*") {
		if v, ok := pointer.(reflect.Value); ok {
			model = m.Fields(
				reflect.New(
					v.Type().Elem(),
				).Interface(),
			)
		} else {
			model = m.Fields(
				reflect.New(
					reflect.ValueOf(pointer).Elem().Type().Elem(),
				).Interface(),
			)
		}
	}
	all, err := model.All(where...)
	if err != nil {
		return err
	}
	if err = all.Structs(pointer); err != nil {
		return err
	}
	return model.doWithScanStructs(pointer)
}

// Scan automatically calls Struct or Structs function according to the type of parameter `pointer`.
// It calls function doStruct if `pointer` is type of *struct/**struct.
// It calls function doStructs if `pointer` is type of *[]struct/*[]*struct.
//
// The optional parameter `where` is the same as the parameter of Model.Where function,  see Model.Where.
//
// Note that it returns sql.ErrNoRows if the given parameter `pointer` pointed to a variable that has
// default value and there's no record retrieved with the given conditions from table.
//
// Example:
// user := new(User)
// err  := db.Model("user").Where("id", 1).Scan(user)
//
// user := (*User)(nil)
// err  := db.Model("user").Where("id", 1).Scan(&user)
//
// users := ([]User)(nil)
// err   := db.Model("user").Scan(&users)
//
// users := ([]*User)(nil)
// err   := db.Model("user").Scan(&users)
func (m *Model) Scan(pointer interface{}, where ...interface{}) error {
	var (
		reflectInfo = utils.OriginTypeAndKind(pointer)
	)
	if reflectInfo.InputKind != reflect.Ptr {
		return gerror.NewCode(
			gcode.CodeInvalidParameter,
			`the parameter "pointer" for function Scan should type of pointer`,
		)
	}
	switch reflectInfo.OriginKind {
	case reflect.Slice, reflect.Array:
		return m.doStructs(pointer, where...)

	case reflect.Struct, reflect.Invalid:
		return m.doStruct(pointer, where...)

	default:
		return gerror.NewCode(
			gcode.CodeInvalidParameter,
			`element of parameter "pointer" for function Scan should type of struct/*struct/[]struct/[]*struct`,
		)
	}
}

// ScanList converts `r` to struct slice which contains other complex struct attributes.
// Note that the parameter `listPointer` should be type of *[]struct/*[]*struct.
// Usage example:
//
// type Entity struct {
// 	   User       *EntityUser
// 	   UserDetail *EntityUserDetail
//	   UserScores []*EntityUserScores
// }
// var users []*Entity
// or
// var users []Entity
//
// ScanList(&users, "User")
// ScanList(&users, "UserDetail", "User", "uid:Uid")
// ScanList(&users, "UserScores", "User", "uid:Uid")
// The parameters "User"/"UserDetail"/"UserScores" in the example codes specify the target attribute struct
// that current result will be bound to.
// The "uid" in the example codes is the table field name of the result, and the "Uid" is the relational
// struct attribute name. It automatically calculates the HasOne/HasMany relationship with given `relation`
// parameter.
// See the example or unit testing cases for clear understanding for this function.
func (m *Model) ScanList(listPointer interface{}, attributeName string, relation ...string) (err error) {
	result, err := m.All()
	if err != nil {
		return err
	}
	return doScanList(m, result, listPointer, attributeName, relation...)
}

// Count does "SELECT COUNT(x) FROM ..." statement for the model.
// The optional parameter `where` is the same as the parameter of Model.Where function,
// see Model.Where.
func (m *Model) Count(where ...interface{}) (int, error) {
	if len(where) > 0 {
		return m.Where(where[0], where[1:]...).Count()
	}
	var (
		sqlWithHolder, holderArgs = m.getFormattedSqlAndArgs(queryTypeCount, false)
		list, err                 = m.doGetAllBySql(sqlWithHolder, holderArgs...)
	)
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

// CountColumn does "SELECT COUNT(x) FROM ..." statement for the model.
func (m *Model) CountColumn(column string) (int, error) {
	if len(column) == 0 {
		return 0, nil
	}
	return m.Fields(column).Count()
}

// Min does "SELECT MIN(x) FROM ..." statement for the model.
func (m *Model) Min(column string) (float64, error) {
	if len(column) == 0 {
		return 0, nil
	}
	value, err := m.Fields(fmt.Sprintf(`MIN(%s)`, m.db.GetCore().QuoteWord(column))).Value()
	if err != nil {
		return 0, err
	}
	return value.Float64(), err
}

// Max does "SELECT MAX(x) FROM ..." statement for the model.
func (m *Model) Max(column string) (float64, error) {
	if len(column) == 0 {
		return 0, nil
	}
	value, err := m.Fields(fmt.Sprintf(`MAX(%s)`, m.db.GetCore().QuoteWord(column))).Value()
	if err != nil {
		return 0, err
	}
	return value.Float64(), err
}

// Avg does "SELECT AVG(x) FROM ..." statement for the model.
func (m *Model) Avg(column string) (float64, error) {
	if len(column) == 0 {
		return 0, nil
	}
	value, err := m.Fields(fmt.Sprintf(`AVG(%s)`, m.db.GetCore().QuoteWord(column))).Value()
	if err != nil {
		return 0, err
	}
	return value.Float64(), err
}

// Sum does "SELECT SUM(x) FROM ..." statement for the model.
func (m *Model) Sum(column string) (float64, error) {
	if len(column) == 0 {
		return 0, nil
	}
	value, err := m.Fields(fmt.Sprintf(`SUM(%s)`, m.db.GetCore().QuoteWord(column))).Value()
	if err != nil {
		return 0, err
	}
	return value.Float64(), err
}

// Union does "(SELECT xxx FROM xxx) UNION (SELECT xxx FROM xxx) ..." statement for the model.
func (m *Model) Union(unions ...*Model) *Model {
	return m.db.Union(unions...)
}

// UnionAll does "(SELECT xxx FROM xxx) UNION ALL (SELECT xxx FROM xxx) ..." statement for the model.
func (m *Model) UnionAll(unions ...*Model) *Model {
	return m.db.UnionAll(unions...)
}

// doGetAllBySql does the select statement on the database.
func (m *Model) doGetAllBySql(sql string, args ...interface{}) (result Result, err error) {
	var (
		ctx      = m.GetCtx()
		cacheKey = ""
		cacheObj = m.db.GetCache()
	)
	// Retrieve from cache.
	if m.cacheEnabled && m.tx == nil {
		cacheKey = m.cacheName
		if len(cacheKey) == 0 {
			cacheKey = sql + ", @PARAMS:" + gconv.String(args)
		}
		if v, _ := cacheObj.Get(ctx, cacheKey); !v.IsNil() {
			if result, ok := v.Val().(Result); ok {
				// In-memory cache.
				return result, nil
			} else {
				// Other cache, it needs conversion.
				var result Result
				if err = json.UnmarshalUseNumber(v.Bytes(), &result); err != nil {
					return nil, err
				} else {
					return result, nil
				}
			}
		}
	}
	result, err = m.db.DoGetAll(
		m.GetCtx(), m.getLink(false), sql, m.mergeArguments(args)...,
	)
	// Cache the result.
	if cacheKey != "" && err == nil {
		if m.cacheDuration < 0 {
			if _, err := cacheObj.Remove(ctx, cacheKey); err != nil {
				intlog.Error(m.GetCtx(), err)
			}
		} else {
			// In case of Cache Penetration.
			if result == nil {
				result = Result{}
			}
			if err := cacheObj.Set(ctx, cacheKey, result, m.cacheDuration); err != nil {
				intlog.Error(m.GetCtx(), err)
			}
		}
	}
	return result, err
}

func (m *Model) getFormattedSqlAndArgs(queryType int, limit1 bool) (sqlWithHolder string, holderArgs []interface{}) {
	switch queryType {
	case queryTypeCount:
		countFields := "COUNT(1)"
		if m.fields != "" && m.fields != "*" {
			// DO NOT quote the m.fields here, in case of fields like:
			// DISTINCT t.user_id uid
			countFields = fmt.Sprintf(`COUNT(%s%s)`, m.distinct, m.fields)
		}
		// Raw SQL Model.
		if m.rawSql != "" {
			sqlWithHolder = fmt.Sprintf("SELECT %s FROM (%s) AS T", countFields, m.rawSql)
			return sqlWithHolder, nil
		}
		conditionWhere, conditionExtra, conditionArgs := m.formatCondition(false, true)
		sqlWithHolder = fmt.Sprintf("SELECT %s FROM %s%s", countFields, m.tables, conditionWhere+conditionExtra)
		if len(m.groupBy) > 0 {
			sqlWithHolder = fmt.Sprintf("SELECT COUNT(1) FROM (%s) count_alias", sqlWithHolder)
		}
		return sqlWithHolder, conditionArgs

	default:
		conditionWhere, conditionExtra, conditionArgs := m.formatCondition(limit1, false)
		// Raw SQL Model, especially for UNION/UNION ALL featured SQL.
		if m.rawSql != "" {
			sqlWithHolder = fmt.Sprintf(
				"%s%s",
				m.rawSql,
				conditionWhere+conditionExtra,
			)
			return sqlWithHolder, conditionArgs
		}
		// DO NOT quote the m.fields where, in case of fields like:
		// DISTINCT t.user_id uid
		sqlWithHolder = fmt.Sprintf(
			"SELECT %s%s FROM %s%s",
			m.distinct,
			m.getFieldsFiltered(),
			m.tables,
			conditionWhere+conditionExtra,
		)
		return sqlWithHolder, conditionArgs
	}
}
