// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"database/sql"
	"reflect"

	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/reflection"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
)

// Batch sets the batch operation number for the model.
func (m *Model) Batch(batch int) *Model {
	model := m.getModel()
	model.batch = batch
	return model
}

// Data sets the operation data for the model.
// The parameter `data` can be type of string/map/gmap/slice/struct/*struct, etc.
// Note that, it uses shallow value copying for `data` if `data` is type of map/slice
// to avoid changing it inside function.
// Eg:
// Data("uid=10000")
// Data("uid", 10000)
// Data("uid=? AND name=?", 10000, "john")
// Data(g.Map{"uid": 10000, "name":"john"})
// Data(g.Slice{g.Map{"uid": 10000, "name":"john"}, g.Map{"uid": 20000, "name":"smith"}).
func (m *Model) Data(data ...interface{}) *Model {
	model := m.getModel()
	if len(data) > 1 {
		if s := gconv.String(data[0]); gstr.Contains(s, "?") {
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
		switch value := data[0].(type) {
		case Result:
			model.data = value.List()

		case Record:
			model.data = value.Map()

		case List:
			list := make(List, len(value))
			for k, v := range value {
				list[k] = gutil.MapCopy(v)
			}
			model.data = list

		case Map:
			model.data = gutil.MapCopy(value)

		default:
			reflectInfo := reflection.OriginValueAndKind(value)
			switch reflectInfo.OriginKind {
			case reflect.Slice, reflect.Array:
				if reflectInfo.OriginValue.Len() > 0 {
					// If the `data` parameter is a DO struct,
					// it then adds `OmitNilData` option for this condition,
					// which will filter all nil parameters in `data`.
					if isDoStruct(reflectInfo.OriginValue.Index(0).Interface()) {
						model = model.OmitNilData()
						model.option |= optionOmitNilDataInternal
					}
				}
				list := make(List, reflectInfo.OriginValue.Len())
				for i := 0; i < reflectInfo.OriginValue.Len(); i++ {
					list[i] = m.db.ConvertDataForRecord(m.GetCtx(), reflectInfo.OriginValue.Index(i).Interface())
				}
				model.data = list

			case reflect.Struct:
				// If the `data` parameter is a DO struct,
				// it then adds `OmitNilData` option for this condition,
				// which will filter all nil parameters in `data`.
				if isDoStruct(value) {
					model = model.OmitNilData()
				}
				if v, ok := data[0].(iInterfaces); ok {
					var (
						array = v.Interfaces()
						list  = make(List, len(array))
					)
					for i := 0; i < len(array); i++ {
						list[i] = m.db.ConvertDataForRecord(m.GetCtx(), array[i])
					}
					model.data = list
				} else {
					model.data = m.db.ConvertDataForRecord(m.GetCtx(), data[0])
				}

			case reflect.Map:
				model.data = m.db.ConvertDataForRecord(m.GetCtx(), data[0])

			default:
				model.data = data[0]
			}
		}
	}
	return model
}

// OnDuplicate sets the operations when columns conflicts occurs.
// In MySQL, this is used for "ON DUPLICATE KEY UPDATE" statement.
// The parameter `onDuplicate` can be type of string/Raw/*Raw/map/slice.
// Example:
// OnDuplicate("nickname, age")
// OnDuplicate("nickname", "age")
// OnDuplicate(g.Map{
//     "nickname": gdb.Raw("CONCAT('name_', VALUES(`nickname`))"),
// })
// OnDuplicate(g.Map{
//     "nickname": "passport",
// }).
func (m *Model) OnDuplicate(onDuplicate ...interface{}) *Model {
	model := m.getModel()
	if len(onDuplicate) > 1 {
		model.onDuplicate = onDuplicate
	} else {
		model.onDuplicate = onDuplicate[0]
	}
	return model
}

// OnDuplicateEx sets the excluding columns for operations when columns conflict occurs.
// In MySQL, this is used for "ON DUPLICATE KEY UPDATE" statement.
// The parameter `onDuplicateEx` can be type of string/map/slice.
// Example:
// OnDuplicateEx("passport, password")
// OnDuplicateEx("passport", "password")
// OnDuplicateEx(g.Map{
//     "passport": "",
//     "password": "",
// }).
func (m *Model) OnDuplicateEx(onDuplicateEx ...interface{}) *Model {
	model := m.getModel()
	if len(onDuplicateEx) > 1 {
		model.onDuplicateEx = onDuplicateEx
	} else {
		model.onDuplicateEx = onDuplicateEx[0]
	}
	return model
}

// Insert does "INSERT INTO ..." statement for the model.
// The optional parameter `data` is the same as the parameter of Model.Data function,
// see Model.Data.
func (m *Model) Insert(data ...interface{}) (result sql.Result, err error) {
	if len(data) > 0 {
		return m.Data(data...).Insert()
	}
	return m.doInsertWithOption(InsertOptionDefault)
}

// InsertAndGetId performs action Insert and returns the last insert id that automatically generated.
func (m *Model) InsertAndGetId(data ...interface{}) (lastInsertId int64, err error) {
	if len(data) > 0 {
		return m.Data(data...).InsertAndGetId()
	}
	result, err := m.doInsertWithOption(InsertOptionDefault)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// InsertIgnore does "INSERT IGNORE INTO ..." statement for the model.
// The optional parameter `data` is the same as the parameter of Model.Data function,
// see Model.Data.
func (m *Model) InsertIgnore(data ...interface{}) (result sql.Result, err error) {
	if len(data) > 0 {
		return m.Data(data...).InsertIgnore()
	}
	return m.doInsertWithOption(InsertOptionIgnore)
}

// Replace does "REPLACE INTO ..." statement for the model.
// The optional parameter `data` is the same as the parameter of Model.Data function,
// see Model.Data.
func (m *Model) Replace(data ...interface{}) (result sql.Result, err error) {
	if len(data) > 0 {
		return m.Data(data...).Replace()
	}
	return m.doInsertWithOption(InsertOptionReplace)
}

// Save does "INSERT INTO ... ON DUPLICATE KEY UPDATE..." statement for the model.
// The optional parameter `data` is the same as the parameter of Model.Data function,
// see Model.Data.
//
// It updates the record if there's primary or unique index in the saving data,
// or else it inserts a new record into the table.
func (m *Model) Save(data ...interface{}) (result sql.Result, err error) {
	if len(data) > 0 {
		return m.Data(data...).Save()
	}
	return m.doInsertWithOption(InsertOptionSave)
}

// doInsertWithOption inserts data with option parameter.
func (m *Model) doInsertWithOption(insertOption int) (result sql.Result, err error) {
	defer func() {
		if err == nil {
			m.checkAndRemoveCache()
		}
	}()
	if m.data == nil {
		return nil, gerror.NewCode(gcode.CodeMissingParameter, "inserting into table with empty data")
	}
	var (
		list            List
		nowString       = gtime.Now().String()
		fieldNameCreate = m.getSoftFieldNameCreated()
		fieldNameUpdate = m.getSoftFieldNameUpdated()
	)
	newData, err := m.filterDataForInsertOrUpdate(m.data)
	if err != nil {
		return nil, err
	}
	// It converts any data to List type for inserting.
	switch value := newData.(type) {
	case Result:
		list = value.List()

	case Record:
		list = List{value.Map()}

	case List:
		list = value
		for i, v := range list {
			list[i] = m.db.ConvertDataForRecord(m.GetCtx(), v)
		}

	case Map:
		list = List{m.db.ConvertDataForRecord(m.GetCtx(), value)}

	default:
		reflectInfo := reflection.OriginValueAndKind(newData)
		switch reflectInfo.OriginKind {
		// If it's slice type, it then converts it to List type.
		case reflect.Slice, reflect.Array:
			list = make(List, reflectInfo.OriginValue.Len())
			for i := 0; i < reflectInfo.OriginValue.Len(); i++ {
				list[i] = m.db.ConvertDataForRecord(m.GetCtx(), reflectInfo.OriginValue.Index(i).Interface())
			}

		case reflect.Map:
			list = List{m.db.ConvertDataForRecord(m.GetCtx(), value)}

		case reflect.Struct:
			if v, ok := value.(iInterfaces); ok {
				array := v.Interfaces()
				list = make(List, len(array))
				for i := 0; i < len(array); i++ {
					list[i] = m.db.ConvertDataForRecord(m.GetCtx(), array[i])
				}
			} else {
				list = List{m.db.ConvertDataForRecord(m.GetCtx(), value)}
			}

		default:
			return result, gerror.NewCodef(
				gcode.CodeInvalidParameter,
				"unsupported data list type: %v",
				reflectInfo.InputValue.Type(),
			)
		}
	}

	if len(list) < 1 {
		return result, gerror.NewCode(gcode.CodeMissingParameter, "data list cannot be empty")
	}

	// Automatic handling for creating/updating time.
	if !m.unscoped && (fieldNameCreate != "" || fieldNameUpdate != "") {
		for k, v := range list {
			if fieldNameCreate != "" {
				v[fieldNameCreate] = nowString
			}
			if fieldNameUpdate != "" {
				v[fieldNameUpdate] = nowString
			}
			list[k] = v
		}
	}
	// Format DoInsertOption, especially for "ON DUPLICATE KEY UPDATE" statement.
	columnNames := make([]string, 0, len(list[0]))
	for k := range list[0] {
		columnNames = append(columnNames, k)
	}
	doInsertOption, err := m.formatDoInsertOption(insertOption, columnNames)
	if err != nil {
		return result, err
	}
	return m.db.DoInsert(m.GetCtx(), m.getLink(true), m.tables, list, doInsertOption)
}

func (m *Model) formatDoInsertOption(insertOption int, columnNames []string) (option DoInsertOption, err error) {
	option = DoInsertOption{
		InsertOption: insertOption,
		BatchCount:   m.getBatch(),
	}
	if insertOption == InsertOptionSave {
		onDuplicateExKeys, err := m.formatOnDuplicateExKeys(m.onDuplicateEx)
		if err != nil {
			return option, err
		}
		onDuplicateExKeySet := gset.NewStrSetFrom(onDuplicateExKeys)
		if m.onDuplicate != nil {
			switch m.onDuplicate.(type) {
			case Raw, *Raw:
				option.OnDuplicateStr = gconv.String(m.onDuplicate)

			default:
				reflectInfo := reflection.OriginValueAndKind(m.onDuplicate)
				switch reflectInfo.OriginKind {
				case reflect.String:
					option.OnDuplicateMap = make(map[string]interface{})
					for _, v := range gstr.SplitAndTrim(reflectInfo.OriginValue.String(), ",") {
						if onDuplicateExKeySet.Contains(v) {
							continue
						}
						option.OnDuplicateMap[v] = v
					}

				case reflect.Map:
					option.OnDuplicateMap = make(map[string]interface{})
					for k, v := range gconv.Map(m.onDuplicate) {
						if onDuplicateExKeySet.Contains(k) {
							continue
						}
						option.OnDuplicateMap[k] = v
					}

				case reflect.Slice, reflect.Array:
					option.OnDuplicateMap = make(map[string]interface{})
					for _, v := range gconv.Strings(m.onDuplicate) {
						if onDuplicateExKeySet.Contains(v) {
							continue
						}
						option.OnDuplicateMap[v] = v
					}

				default:
					return option, gerror.NewCodef(
						gcode.CodeInvalidParameter,
						`unsupported OnDuplicate parameter type "%s"`,
						reflect.TypeOf(m.onDuplicate),
					)
				}
			}
		} else if onDuplicateExKeySet.Size() > 0 {
			option.OnDuplicateMap = make(map[string]interface{})
			for _, v := range columnNames {
				if onDuplicateExKeySet.Contains(v) {
					continue
				}
				option.OnDuplicateMap[v] = v
			}
		}
	}
	return
}

func (m *Model) formatOnDuplicateExKeys(onDuplicateEx interface{}) ([]string, error) {
	if onDuplicateEx == nil {
		return nil, nil
	}

	reflectInfo := reflection.OriginValueAndKind(onDuplicateEx)
	switch reflectInfo.OriginKind {
	case reflect.String:
		return gstr.SplitAndTrim(reflectInfo.OriginValue.String(), ","), nil

	case reflect.Map:
		return gutil.Keys(onDuplicateEx), nil

	case reflect.Slice, reflect.Array:
		return gconv.Strings(onDuplicateEx), nil

	default:
		return nil, gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`unsupported OnDuplicateEx parameter type "%s"`,
			reflect.TypeOf(onDuplicateEx),
		)
	}
}

func (m *Model) getBatch() int {
	batch := defaultBatchNumber
	if m.batch > 0 {
		batch = m.batch
	}
	return batch
}
