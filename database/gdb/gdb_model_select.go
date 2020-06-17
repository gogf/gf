// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"fmt"
	"github.com/gogf/gf/container/gvar"
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
	if len(where) > 0 {
		return m.Where(where[0], where[1:]...).All()
	}
	var (
		softDeletingCondition                         = m.getConditionForSoftDeleting()
		conditionWhere, conditionExtra, conditionArgs = m.formatCondition(false)
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
	return m.doGetAll(
		fmt.Sprintf(
			"SELECT %s FROM %s%s",
			m.fields,
			m.tables,
			conditionWhere+conditionExtra,
		),
		conditionArgs...,
	)
}

// Chunk iterates the query result with given size and callback function.
func (m *Model) Chunk(limit int, callback func(result Result, err error) bool) {
	page := m.start
	if page == 0 {
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
	all, err := m.All()
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
		conditionWhere, conditionExtra, conditionArgs = m.formatCondition(false)
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
	list, err := m.doGetAll(s, conditionArgs...)
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

// doGetAll does the select statement on the database.
func (m *Model) doGetAll(sql string, args ...interface{}) (result Result, err error) {
	cacheKey := ""
	// Retrieve from cache.
	if m.cacheEnabled && m.tx == nil {
		cacheKey = m.cacheName
		if len(cacheKey) == 0 {
			cacheKey = sql + "/" + gconv.String(args)
		}
		if v := m.db.GetCache().Get(cacheKey); v != nil {
			return v.(Result), nil
		}
	}
	result, err = m.db.DoGetAll(m.getLink(false), sql, m.mergeArguments(args)...)
	// Cache the result.
	if cacheKey != "" && err == nil {
		if m.cacheDuration < 0 {
			m.db.GetCache().Remove(cacheKey)
		} else {
			m.db.GetCache().Set(cacheKey, result, m.cacheDuration)
		}
	}
	return result, err
}
