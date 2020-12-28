// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"fmt"
	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/internal/json"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"reflect"
)

// Select is alias of Model.All.
// See Model.All.
// Deprecated.
func (m *Model) Select(where ...interface{}) (Result, error) {
	return m.All(where...)
}

// All does "SELECT FROM ..." statement for the model.
// It retrieves the records from table and returns the result as slice type.
// It returns nil if there's no record retrieved with the given conditions from table.
//
// The optional parameter <where> is the same as the parameter of Model.Where function,
// see Model.Where.
func (m *Model) All(where ...interface{}) (Result, error) {
	return m.doGetAll(false, where...)
}

// doGetAll does "SELECT FROM ..." statement for the model.
// It retrieves the records from table and returns the result as slice type.
// It returns nil if there's no record retrieved with the given conditions from table.
//
// The parameter <limit1> specifies whether limits querying only one record if m.limit is not set.
// The optional parameter <where> is the same as the parameter of Model.Where function,
// see Model.Where.
func (m *Model) doGetAll(limit1 bool, where ...interface{}) (Result, error) {
	if len(where) > 0 {
		return m.Where(where[0], where[1:]...).All()
	}
	var (
		softDeletingCondition                         = m.getConditionForSoftDeleting()
		conditionWhere, conditionExtra, conditionArgs = m.formatCondition(limit1, false)
	)
	if !m.unscoped && softDeletingCondition != "" {
		if conditionWhere == "" {
			conditionWhere = " WHERE "
		} else {
			conditionWhere += " AND "
		}
		conditionWhere += softDeletingCondition
	}

	// DO NOT quote the m.fields where, in case of fields like:
	// DISTINCT t.user_id uid
	return m.doGetAllBySql(
		fmt.Sprintf(
			"SELECT %s FROM %s%s",
			m.getFieldsFiltered(),
			m.tables,
			conditionWhere+conditionExtra,
		),
		conditionArgs...,
	)
}

// getFieldsFiltered checks the fields and fieldsEx attributes, filters and returns the fields that will
// really be committed to underlying database driver.
func (m *Model) getFieldsFiltered() string {
	if m.fieldsEx == "" {
		// No filtering.
		if !gstr.Contains(m.fields, ".") && !gstr.Contains(m.fields, " ") {
			return m.db.QuoteString(m.fields)
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
		tableFields, err := m.db.TableFields(m.tables)
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
		newFields += m.db.QuoteWord(k)
	}
	return newFields
}

// Chunk iterates the query result with given size and callback function.
func (m *Model) Chunk(limit int, callback func(result Result, err error) bool) {
	page := m.start
	if page <= 0 {
		page = 1
	}
	model := m
	for {
		model = model.Page(page, limit)
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

// One retrieves one record from table and returns the result as map type.
// It returns nil if there's no record retrieved with the given conditions from table.
//
// The optional parameter <where> is the same as the parameter of Model.Where function,
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
// If the optional parameter <fieldsAndWhere> is given, the fieldsAndWhere[0] is the selected fields
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
// Note that if there're multiple columns in the result, it returns just one column values randomly.
//
// If the optional parameter <fieldsAndWhere> is given, the fieldsAndWhere[0] is the selected fields
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
// The parameter <pointer> should be type of *struct/**struct. If type **struct is given,
// it can create the struct internally during converting.
//
// The optional parameter <where> is the same as the parameter of Model.Where function,
// see Model.Where.
//
// Note that it returns sql.ErrNoRows if there's no record retrieved with the given conditions
// from table and <pointer> is not nil.
//
// Eg:
// user := new(User)
// err  := db.Table("user").Where("id", 1).Struct(user)
//
// user := (*User)(nil)
// err  := db.Table("user").Where("id", 1).Struct(&user)
func (m *Model) Struct(pointer interface{}, where ...interface{}) error {
	one, err := m.One(where...)
	if err != nil {
		return err
	}
	return one.Struct(pointer)
}

// Structs retrieves records from table and converts them into given struct slice.
// The parameter <pointer> should be type of *[]struct/*[]*struct. It can create and fill the struct
// slice internally during converting.
//
// The optional parameter <where> is the same as the parameter of Model.Where function,
// see Model.Where.
//
// Note that it returns sql.ErrNoRows if there's no record retrieved with the given conditions
// from table and <pointer> is not empty.
//
// Eg:
// users := ([]User)(nil)
// err   := db.Table("user").Structs(&users)
//
// users := ([]*User)(nil)
// err   := db.Table("user").Structs(&users)
func (m *Model) Structs(pointer interface{}, where ...interface{}) error {
	all, err := m.All(where...)
	if err != nil {
		return err
	}
	return all.Structs(pointer)
}

// Scan automatically calls Struct or Structs function according to the type of parameter <pointer>.
// It calls function Struct if <pointer> is type of *struct/**struct.
// It calls function Structs if <pointer> is type of *[]struct/*[]*struct.
//
// The optional parameter <where> is the same as the parameter of Model.Where function,
// see Model.Where.
//
// Note that it returns sql.ErrNoRows if there's no record retrieved with the given conditions
// from table.
//
// Eg:
// user := new(User)
// err  := db.Table("user").Where("id", 1).Struct(user)
//
// user := (*User)(nil)
// err  := db.Table("user").Where("id", 1).Struct(&user)
//
// users := ([]User)(nil)
// err   := db.Table("user").Structs(&users)
//
// users := ([]*User)(nil)
// err   := db.Table("user").Structs(&users)
func (m *Model) Scan(pointer interface{}, where ...interface{}) error {
	t := reflect.TypeOf(pointer)
	k := t.Kind()
	if k != reflect.Ptr {
		return fmt.Errorf("params should be type of pointer, but got: %v", k)
	}
	switch t.Elem().Kind() {
	case reflect.Array, reflect.Slice:
		return m.Structs(pointer, where...)
	default:
		return m.Struct(pointer, where...)
	}
}

// ScanList converts <r> to struct slice which contains other complex struct attributes.
// Note that the parameter <listPointer> should be type of *[]struct/*[]*struct.
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
// struct attribute name. It automatically calculates the HasOne/HasMany relationship with given <relation>
// parameter.
// See the example or unit testing cases for clear understanding for this function.
func (m *Model) ScanList(listPointer interface{}, attributeName string, relation ...string) (err error) {
	all, err := m.All()
	if err != nil {
		return err
	}
	return all.ScanList(listPointer, attributeName, relation...)
}

// Count does "SELECT COUNT(x) FROM ..." statement for the model.
// The optional parameter <where> is the same as the parameter of Model.Where function,
// see Model.Where.
func (m *Model) Count(where ...interface{}) (int, error) {
	if len(where) > 0 {
		return m.Where(where[0], where[1:]...).Count()
	}
	countFields := "COUNT(1)"
	if m.fields != "" && m.fields != "*" {
		// DO NOT quote the m.fields here, in case of fields like:
		// DISTINCT t.user_id uid
		countFields = fmt.Sprintf(`COUNT(%s)`, m.fields)
	}
	var (
		softDeletingCondition                         = m.getConditionForSoftDeleting()
		conditionWhere, conditionExtra, conditionArgs = m.formatCondition(false, true)
	)
	if !m.unscoped && softDeletingCondition != "" {
		if conditionWhere == "" {
			conditionWhere = " WHERE "
		} else {
			conditionWhere += " AND "
		}
		conditionWhere += softDeletingCondition
	}

	s := fmt.Sprintf("SELECT %s FROM %s%s", countFields, m.tables, conditionWhere+conditionExtra)
	if len(m.groupBy) > 0 {
		s = fmt.Sprintf("SELECT COUNT(1) FROM (%s) count_alias", s)
	}
	list, err := m.doGetAllBySql(s, conditionArgs...)
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

// FindOne retrieves and returns a single Record by Model.WherePri and Model.One.
// Also see Model.WherePri and Model.One.
func (m *Model) FindOne(where ...interface{}) (Record, error) {
	if len(where) > 0 {
		return m.WherePri(where[0], where[1:]...).One()
	}
	return m.One()
}

// FindAll retrieves and returns Result by by Model.WherePri and Model.All.
// Also see Model.WherePri and Model.All.
func (m *Model) FindAll(where ...interface{}) (Result, error) {
	if len(where) > 0 {
		return m.WherePri(where[0], where[1:]...).All()
	}
	return m.All()
}

// FindValue retrieves and returns single field value by Model.WherePri and Model.Value.
// Also see Model.WherePri and Model.Value.
func (m *Model) FindValue(fieldsAndWhere ...interface{}) (Value, error) {
	if len(fieldsAndWhere) >= 2 {
		return m.WherePri(fieldsAndWhere[1], fieldsAndWhere[2:]...).Fields(gconv.String(fieldsAndWhere[0])).Value()
	}
	if len(fieldsAndWhere) == 1 {
		return m.Fields(gconv.String(fieldsAndWhere[0])).Value()
	}
	return m.Value()
}

// FindArray queries and returns data values as slice from database.
// Note that if there're multiple columns in the result, it returns just one column values randomly.
// Also see Model.WherePri and Model.Value.
func (m *Model) FindArray(fieldsAndWhere ...interface{}) ([]Value, error) {
	if len(fieldsAndWhere) >= 2 {
		return m.WherePri(fieldsAndWhere[1], fieldsAndWhere[2:]...).Fields(gconv.String(fieldsAndWhere[0])).Array()
	}
	if len(fieldsAndWhere) == 1 {
		return m.Fields(gconv.String(fieldsAndWhere[0])).Array()
	}
	return m.Array()
}

// FindCount retrieves and returns the record number by Model.WherePri and Model.Count.
// Also see Model.WherePri and Model.Count.
func (m *Model) FindCount(where ...interface{}) (int, error) {
	if len(where) > 0 {
		return m.WherePri(where[0], where[1:]...).Count()
	}
	return m.Count()
}

// FindScan retrieves and returns the record/records by Model.WherePri and Model.Scan.
// Also see Model.WherePri and Model.Scan.
func (m *Model) FindScan(pointer interface{}, where ...interface{}) error {
	if len(where) > 0 {
		return m.WherePri(where[0], where[1:]...).Scan(pointer)
	}
	return m.Scan(pointer)
}

// doGetAllBySql does the select statement on the database.
func (m *Model) doGetAllBySql(sql string, args ...interface{}) (result Result, err error) {
	cacheKey := ""
	cacheObj := m.db.GetCache()
	// Retrieve from cache.
	if m.cacheEnabled && m.tx == nil {
		cacheKey = m.cacheName
		if len(cacheKey) == 0 {
			cacheKey = sql + ", @PARAMS:" + gconv.String(args)
		}
		if v, _ := cacheObj.GetVar(cacheKey); !v.IsNil() {
			if result, ok := v.Val().(Result); ok {
				// In-memory cache.
				return result, nil
			} else {
				// Other cache, it needs conversion.
				var result Result
				if err = json.Unmarshal(v.Bytes(), &result); err != nil {
					return nil, err
				} else {
					return result, nil
				}
			}
		}
	}
	result, err = m.db.DoGetAll(m.getLink(false), sql, m.mergeArguments(args)...)
	// Cache the result.
	if cacheKey != "" && err == nil {
		if m.cacheDuration < 0 {
			cacheObj.Remove(cacheKey)
		} else {
			cacheObj.Set(cacheKey, result, m.cacheDuration)
		}
	}
	return result, err
}
