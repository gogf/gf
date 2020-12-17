// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"database/sql"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gutil"
	"reflect"
)

// Batch sets the batch operation number for the model.
func (m *Model) Batch(batch int) *Model {
	model := m.getModel()
	model.batch = batch
	return model
}

// Data sets the operation data for the model.
// The parameter <data> can be type of string/map/gmap/slice/struct/*struct, etc.
// Eg:
// Data("uid=10000")
// Data("uid", 10000)
// Data("uid=? AND name=?", 10000, "john")
// Data(g.Map{"uid": 10000, "name":"john"})
// Data(g.Slice{g.Map{"uid": 10000, "name":"john"}, g.Map{"uid": 20000, "name":"smith"})
func (m *Model) Data(data ...interface{}) *Model {
	model := m.getModel()
	if len(data) > 1 {
		s := gconv.String(data[0])
		if gstr.Contains(s, "?") {
			model.data = s
			model.extraArgs = data[1:]
		} else {
			m := make(map[string]interface{})
			for i := 0; i < len(data); i += 2 {
				m[gconv.String(data[i])] = data[i+1]
			}
			model.data = m
		}
	} else {
		switch params := data[0].(type) {
		case Result:
			model.data = params.List()
		case Record:
			model.data = params.Map()
		case List:
			model.data = params
		case Map:
			model.data = params
		default:
			var (
				rv   = reflect.ValueOf(params)
				kind = rv.Kind()
			)
			if kind == reflect.Ptr {
				rv = rv.Elem()
				kind = rv.Kind()
			}
			switch kind {
			case reflect.Slice, reflect.Array:
				list := make(List, rv.Len())
				for i := 0; i < rv.Len(); i++ {
					list[i] = ConvertDataForTableRecord(rv.Index(i).Interface())
				}
				model.data = list
			case reflect.Map:
				model.data = ConvertDataForTableRecord(data[0])
			case reflect.Struct:
				if v, ok := data[0].(apiInterfaces); ok {
					var (
						array = v.Interfaces()
						list  = make(List, len(array))
					)
					for i := 0; i < len(array); i++ {
						list[i] = ConvertDataForTableRecord(array[i])
					}
					model.data = list
				} else {
					model.data = ConvertDataForTableRecord(data[0])
				}
			default:
				model.data = data[0]
			}
		}
	}
	return model
}

// Insert does "INSERT INTO ..." statement for the model.
// The optional parameter <data> is the same as the parameter of Model.Data function,
// see Model.Data.
func (m *Model) Insert(data ...interface{}) (result sql.Result, err error) {
	if len(data) > 0 {
		return m.Data(data...).Insert()
	}
	return m.doInsertWithOption(insertOptionDefault, data...)
}

// InsertIgnore does "INSERT IGNORE INTO ..." statement for the model.
// The optional parameter <data> is the same as the parameter of Model.Data function,
// see Model.Data.
func (m *Model) InsertIgnore(data ...interface{}) (result sql.Result, err error) {
	if len(data) > 0 {
		return m.Data(data...).InsertIgnore()
	}
	return m.doInsertWithOption(insertOptionIgnore, data...)
}

// Replace does "REPLACE INTO ..." statement for the model.
// The optional parameter <data> is the same as the parameter of Model.Data function,
// see Model.Data.
func (m *Model) Replace(data ...interface{}) (result sql.Result, err error) {
	if len(data) > 0 {
		return m.Data(data...).Replace()
	}
	return m.doInsertWithOption(insertOptionReplace, data...)
}

// Save does "INSERT INTO ... ON DUPLICATE KEY UPDATE..." statement for the model.
// The optional parameter <data> is the same as the parameter of Model.Data function,
// see Model.Data.
//
// It updates the record if there's primary or unique index in the saving data,
// or else it inserts a new record into the table.
func (m *Model) Save(data ...interface{}) (result sql.Result, err error) {
	if len(data) > 0 {
		return m.Data(data...).Save()
	}
	return m.doInsertWithOption(insertOptionSave, data...)
}

// doInsertWithOption inserts data with option parameter.
func (m *Model) doInsertWithOption(option int, data ...interface{}) (result sql.Result, err error) {
	defer func() {
		if err == nil {
			m.checkAndRemoveCache()
		}
	}()
	if m.data == nil {
		return nil, gerror.New("inserting into table with empty data")
	}
	var (
		nowString       = gtime.Now().String()
		fieldNameCreate = m.getSoftFieldNameCreated()
		fieldNameUpdate = m.getSoftFieldNameUpdated()
		fieldNameDelete = m.getSoftFieldNameDeleted()
	)
	// Batch operation.
	if list, ok := m.data.(List); ok {
		batch := defaultBatchNumber
		if m.batch > 0 {
			batch = m.batch
		}
		// Automatic handling for creating/updating time.
		if !m.unscoped && (fieldNameCreate != "" || fieldNameUpdate != "") {
			for k, v := range list {
				gutil.MapDelete(v, fieldNameCreate, fieldNameUpdate, fieldNameDelete)
				if fieldNameCreate != "" {
					v[fieldNameCreate] = nowString
				}
				if fieldNameUpdate != "" {
					v[fieldNameUpdate] = nowString
				}
				list[k] = v
			}
		}
		newData, err := m.filterDataForInsertOrUpdate(list)
		if err != nil {
			return nil, err
		}
		return m.db.DoBatchInsert(
			m.getLink(true),
			m.tables,
			newData,
			option,
			batch,
		)
	}
	// Single operation.
	if data, ok := m.data.(Map); ok {
		// Automatic handling for creating/updating time.
		if !m.unscoped && (fieldNameCreate != "" || fieldNameUpdate != "") {
			gutil.MapDelete(data, fieldNameCreate, fieldNameUpdate, fieldNameDelete)
			if fieldNameCreate != "" {
				data[fieldNameCreate] = nowString
			}
			if fieldNameUpdate != "" {
				data[fieldNameUpdate] = nowString
			}
		}
		newData, err := m.filterDataForInsertOrUpdate(data)
		if err != nil {
			return nil, err
		}
		return m.db.DoInsert(
			m.getLink(true),
			m.tables,
			newData,
			option,
		)
	}
	return nil, gerror.New("inserting into table with invalid data type")
}
