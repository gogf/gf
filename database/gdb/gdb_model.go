// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/container/gmap"
	"reflect"
	"time"

	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/text/gstr"

	"github.com/gogf/gf/util/gconv"
)

// Model is the DAO for ORM.
type Model struct {
	db            DB             // Underlying DB interface.
	tx            *TX            // Underlying TX interface.
	schema        string         // Custom database schema.
	linkType      int            // Mark for operation on master or slave.
	tablesInit    string         // Table names when model initialization.
	tables        string         // Operation table names, which can be more than one table names and aliases, like: "user", "user u", "user u, user_detail ud".
	fields        string         // Operation fields, multiple fields joined using char ','.
	fieldsEx      string         // Excluded operation fields, multiple fields joined using char ','.
	whereArgs     []interface{}  // Arguments for where operation.
	whereHolder   []*whereHolder // Condition strings for where operation.
	groupBy       string         // Used for "group by" statement.
	orderBy       string         // Used for "order by" statement.
	start         int            // Used for "select ... start, limit ..." statement.
	limit         int            // Used for "select ... start, limit ..." statement.
	option        int            // Option for extra operation features.
	offset        int            // Offset statement for some databases grammar.
	data          interface{}    // Data for operation, which can be type of map/[]map/struct/*struct/string, etc.
	batch         int            // Batch number for batch Insert/Replace/Save operations.
	filter        bool           // Filter data and where key-value pairs according to the fields of the table.
	cacheEnabled  bool           // Enable sql result cache feature.
	cacheDuration time.Duration  // Cache TTL duration.
	cacheName     string         // Cache name for custom operation.
	safe          bool           // If true, it clones and returns a new model object whenever operation done; or else it changes the attribute of current model.
}

// whereHolder is the holder for where condition preparing.
type whereHolder struct {
	operator int           // Operator for this holder.
	where    interface{}   // Where parameter.
	args     []interface{} // Arguments for where parameter.
}

const (
	gLINK_TYPE_MASTER   = 1
	gLINK_TYPE_SLAVE    = 2
	gWHERE_HOLDER_WHERE = 1
	gWHERE_HOLDER_AND   = 2
	gWHERE_HOLDER_OR    = 3
	OPTION_OMITEMPTY    = 1 << iota
	OPTION_ALLOWEMPTY
)

// Table creates and returns a new ORM model from given schema.
// The parameter <tables> can be more than one table names, like :
// "user", "user u", "user, user_detail", "user u, user_detail ud"
func (c *Core) Table(table string) *Model {
	table = c.DB.handleTableName(table)
	return &Model{
		db:         c.DB,
		tablesInit: table,
		tables:     table,
		fields:     "*",
		start:      -1,
		offset:     -1,
		safe:       false,
		option:     OPTION_ALLOWEMPTY,
	}
}

// Model is alias of Core.Table.
// See Core.Table.
func (c *Core) Model(table string) *Model {
	return c.DB.Table(table)
}

// From is alias of Core.Table.
// See Core.Table.
// Deprecated.
func (c *Core) From(table string) *Model {
	return c.DB.Table(table)
}

// Table acts like Core.Table except it operates on transaction.
// See Core.Table.
func (tx *TX) Table(table string) *Model {
	table = tx.db.handleTableName(table)
	return &Model{
		db:         tx.db,
		tx:         tx,
		tablesInit: table,
		tables:     table,
		fields:     "*",
		start:      -1,
		offset:     -1,
		safe:       false,
		option:     OPTION_ALLOWEMPTY,
	}
}

// Model is alias of tx.Table.
// See tx.Table.
func (tx *TX) Model(table string) *Model {
	return tx.Table(table)
}

// From is alias of tx.Table.
// See tx.Table.
// Deprecated.
func (tx *TX) From(table string) *Model {
	return tx.Table(table)
}

// As sets an alias name for current table.
func (m *Model) As(as string) *Model {
	if m.tables != "" {
		model := m.getModel()
		model.tables = gstr.TrimRight(model.tables) + " AS " + as
		return model
	}
	return m
}

// DB sets/changes the db object for current operation.
func (m *Model) DB(db DB) *Model {
	model := m.getModel()
	model.db = db
	return model
}

// TX sets/changes the transaction for current operation.
func (m *Model) TX(tx *TX) *Model {
	model := m.getModel()
	model.tx = tx
	return model
}

// Schema sets the schema for current operation.
func (m *Model) Schema(schema string) *Model {
	model := m.getModel()
	model.schema = schema
	return model
}

// Clone creates and returns a new model which is a clone of current model.
// Note that it uses deep-copy for the clone.
func (m *Model) Clone() *Model {
	newModel := (*Model)(nil)
	if m.tx != nil {
		newModel = m.tx.Table(m.tablesInit)
	} else {
		newModel = m.db.Table(m.tablesInit)
	}
	*newModel = *m
	// Deep copy slice attributes.
	if n := len(m.whereArgs); n > 0 {
		newModel.whereArgs = make([]interface{}, n)
		copy(newModel.whereArgs, m.whereArgs)
	}
	if n := len(m.whereHolder); n > 0 {
		newModel.whereHolder = make([]*whereHolder, n)
		copy(newModel.whereHolder, m.whereHolder)
	}
	return newModel
}

// Master marks the following operation on master node.
func (m *Model) Master() *Model {
	model := m.getModel()
	model.linkType = gLINK_TYPE_MASTER
	return model
}

// Slave marks the following operation on slave node.
// Note that it makes sense only if there's any slave node configured.
func (m *Model) Slave() *Model {
	model := m.getModel()
	model.linkType = gLINK_TYPE_SLAVE
	return model
}

// Safe marks this model safe or unsafe. If safe is true, it clones and returns a new model object
// whenever the operation done, or else it changes the attribute of current model.
func (m *Model) Safe(safe ...bool) *Model {
	if len(safe) > 0 {
		m.safe = safe[0]
	} else {
		m.safe = true
	}
	return m
}

// getModel creates and returns a cloned model of current model if <safe> is true, or else it returns
// the current model.
func (m *Model) getModel() *Model {
	if !m.safe {
		return m
	} else {
		return m.Clone()
	}
}

// LeftJoin does "LEFT JOIN ... ON ..." statement on the model.
func (m *Model) LeftJoin(table string, on string) *Model {
	model := m.getModel()
	model.tables += fmt.Sprintf(" LEFT JOIN %s ON (%s)", m.db.handleTableName(table), on)
	return model
}

// RightJoin does "RIGHT JOIN ... ON ..." statement on the model.
func (m *Model) RightJoin(table string, on string) *Model {
	model := m.getModel()
	model.tables += fmt.Sprintf(" RIGHT JOIN %s ON (%s)", m.db.handleTableName(table), on)
	return model
}

// InnerJoin does "INNER JOIN ... ON ..." statement on the model.
func (m *Model) InnerJoin(table string, on string) *Model {
	model := m.getModel()
	model.tables += fmt.Sprintf(" INNER JOIN %s ON (%s)", m.db.handleTableName(table), on)
	return model
}

// Fields sets the operation fields of the model, multiple fields joined using char ','.
func (m *Model) Fields(fields string) *Model {
	model := m.getModel()
	model.fields = fields
	return model
}

// FieldsEx sets the excluded operation fields of the model, multiple fields joined using char ','.
func (m *Model) FieldsEx(fields string) *Model {
	if gstr.Contains(m.tables, " ") {
		panic("function FieldsEx supports only single table operations")
	}
	model := m.getModel()
	model.fieldsEx = fields
	fieldsExSet := gset.NewStrSetFrom(gstr.SplitAndTrim(fields, ","))
	if m, err := m.db.TableFields(m.tables); err == nil {
		model.fields = ""
		for k, _ := range m {
			if fieldsExSet.Contains(k) {
				continue
			}
			if len(model.fields) > 0 {
				model.fields += ","
			}
			model.fields += k
		}
	}
	return model
}

// FieldsStr retrieves and returns all fields from the table, joined with char ','.
// The optional parameter <prefix> specifies the prefix for each field, eg: FieldsStr("u.").
func (m *Model) FieldsStr(prefix ...string) string {
	prefixStr := ""
	if len(prefix) > 0 {
		prefixStr = prefix[0]
	}
	if m, err := m.db.TableFields(m.tables); err == nil {
		fieldsArray := garray.NewStrArraySize(len(m), len(m))
		for _, field := range m {
			fieldsArray.Set(field.Index, prefixStr+field.Name)
		}
		return fieldsArray.Join(",")
	}
	return ""
}

// FieldsExStr retrieves and returns fields which are not in parameter <fields> from the table,
// joined with char ','.
// The parameter <fields> specifies the fields that are excluded.
// The optional parameter <prefix> specifies the prefix for each field, eg: FieldsExStr("id", "u.").
func (m *Model) FieldsExStr(fields string, prefix ...string) string {
	prefixStr := ""
	if len(prefix) > 0 {
		prefixStr = prefix[0]
	}
	if m, err := m.db.TableFields(m.tables); err == nil {
		fieldsArray := garray.NewStrArraySize(len(m), len(m))
		fieldsExSet := gset.NewStrSetFrom(gstr.SplitAndTrim(fields, ","))
		for _, field := range m {
			if fieldsExSet.Contains(field.Name) {
				continue
			}
			fieldsArray.Set(field.Index, prefixStr+field.Name)
		}
		fieldsArray.FilterEmpty()
		return fieldsArray.Join(",")
	}
	return ""
}

// Option adds extra operation option for the model.
func (m *Model) Option(option int) *Model {
	model := m.getModel()
	model.option = model.option | option
	return model
}

// OptionOmitEmpty sets OPTION_OMITEMPTY option for the model, which automatically filers
// the data and where attributes for empty values.
// Deprecated, use OmitEmpty instead.
func (m *Model) OptionOmitEmpty() *Model {
	return m.Option(OPTION_OMITEMPTY)
}

// OmitEmpty sets OPTION_OMITEMPTY option for the model, which automatically filers
// the data and where attributes for empty values.
func (m *Model) OmitEmpty() *Model {
	return m.Option(OPTION_OMITEMPTY)
}

// Filter marks filtering the fields which does not exist in the fields of the operated table.
func (m *Model) Filter() *Model {
	if gstr.Contains(m.tables, " ") {
		panic("function Filter supports only single table operations")
	}
	model := m.getModel()
	model.filter = true
	return model
}

// Where sets the condition statement for the model. The parameter <where> can be type of
// string/map/gmap/slice/struct/*struct, etc. Note that, if it's called more than one times,
// multiple conditions will be joined into where statement using "AND".
// Eg:
// Where("uid=10000")
// Where("uid", 10000)
// Where("money>? AND name like ?", 99999, "vip_%")
// Where("uid", 1).Where("name", "john")
// Where("status IN (?)", g.Slice{1,2,3})
// Where("age IN(?,?)", 18, 50)
// Where(User{ Id : 1, UserName : "john"})
func (m *Model) Where(where interface{}, args ...interface{}) *Model {
	model := m.getModel()
	if model.whereHolder == nil {
		model.whereHolder = make([]*whereHolder, 0)
	}
	model.whereHolder = append(model.whereHolder, &whereHolder{
		operator: gWHERE_HOLDER_WHERE,
		where:    where,
		args:     args,
	})
	return model
}

// WherePri does the same logic as Model.Where except that if the parameter <where>
// is a single condition like int/string/float/slice, it treats the condition as the primary
// key value. That is, if primary key is "id" and given <where> parameter as "123", the
// WherePri function treats it as "id=123", but Model.Where treats it as string "123".
func (m *Model) WherePri(where interface{}, args ...interface{}) *Model {
	if len(args) > 0 {
		return m.Where(where, args...)
	}
	newWhere := GetPrimaryKeyCondition(m.getPrimaryKey(), where)
	return m.Where(newWhere[0], newWhere[1:]...)
}

// And adds "AND" condition to the where statement.
func (m *Model) And(where interface{}, args ...interface{}) *Model {
	model := m.getModel()
	if model.whereHolder == nil {
		model.whereHolder = make([]*whereHolder, 0)
	}
	model.whereHolder = append(model.whereHolder, &whereHolder{
		operator: gWHERE_HOLDER_AND,
		where:    where,
		args:     args,
	})
	return model
}

// Or adds "OR" condition to the where statement.
func (m *Model) Or(where interface{}, args ...interface{}) *Model {
	model := m.getModel()
	if model.whereHolder == nil {
		model.whereHolder = make([]*whereHolder, 0)
	}
	model.whereHolder = append(model.whereHolder, &whereHolder{
		operator: gWHERE_HOLDER_OR,
		where:    where,
		args:     args,
	})
	return model
}

// Group sets the "GROUP BY" statement for the model.
func (m *Model) Group(groupBy string) *Model {
	model := m.getModel()
	model.groupBy = m.db.QuoteString(groupBy)
	return model
}

// GroupBy is alias of Model.Group.
// See Model.Group.
// Deprecated.
func (m *Model) GroupBy(groupBy string) *Model {
	return m.Group(groupBy)
}

// Order sets the "ORDER BY" statement for the model.
func (m *Model) Order(orderBy string) *Model {
	model := m.getModel()
	model.orderBy = m.db.QuoteString(orderBy)
	return model
}

// OrderBy is alias of Model.Order.
// See Model.Order.
// Deprecated.
func (m *Model) OrderBy(orderBy string) *Model {
	return m.Order(orderBy)
}

// Limit sets the "LIMIT" statement for the model.
// The parameter <limit> can be either one or two number, if passed two number is passed,
// it then sets "LIMIT limit[0],limit[1]" statement for the model, or else it sets "LIMIT limit[0]"
// statement.
func (m *Model) Limit(limit ...int) *Model {
	model := m.getModel()
	switch len(limit) {
	case 1:
		model.limit = limit[0]
	case 2:
		model.start = limit[0]
		model.limit = limit[1]
	}
	return model
}

// Offset sets the "OFFSET" statement for the model.
// It only makes sense for some databases like SQLServer, PostgreSQL, etc.
func (m *Model) Offset(offset int) *Model {
	model := m.getModel()
	model.offset = offset
	return model
}

// Page sets the paging number for the model.
// The parameter <page> is started from 1 for paging.
// Note that, it differs that the Limit function start from 0 for "LIMIT" statement.
func (m *Model) Page(page, limit int) *Model {
	model := m.getModel()
	if page <= 0 {
		page = 1
	}
	model.start = (page - 1) * limit
	model.limit = limit
	return model
}

// ForPage is alias of Model.Page.
// See Model.Page.
// Deprecated.
func (m *Model) ForPage(page, limit int) *Model {
	return m.Page(page, limit)
}

// Batch sets the batch operation number for the model.
func (m *Model) Batch(batch int) *Model {
	model := m.getModel()
	model.batch = batch
	return model
}

// Cache sets the cache feature for the model. It caches the result of the sql, which means
// if there's another same sql request, it just reads and returns the result from cache, it
// but not committed and executed into the database.
//
// If the parameter <duration> < 0, which means it clear the cache with given <name>.
// If the parameter <duration> = 0, which means it never expires.
// If the parameter <duration> > 0, which means it expires after <duration>.
//
// The optional parameter <name> is used to bind a name to the cache, which means you can later
// control the cache like changing the <duration> or clearing the cache with specified <name>.
//
// Note that, the cache feature is disabled if the model is operating on a transaction.
func (m *Model) Cache(duration time.Duration, name ...string) *Model {
	model := m.getModel()
	model.cacheDuration = duration
	if len(name) > 0 {
		model.cacheName = name[0]
	}
	// It does not support cache on transaction.
	if model.tx == nil {
		model.cacheEnabled = true
	}
	return model
}

// Data sets the operation data for the model.
// The parameter <data> can be type of string/map/gmap/slice/struct/*struct, etc.
// Eg:
// Data("uid=10000")
// Data("uid", 10000)
// Data(g.Map{"uid": 10000, "name":"john"})
// Data(g.Slice{g.Map{"uid": 10000, "name":"john"}, g.Map{"uid": 20000, "name":"smith"})
func (m *Model) Data(data ...interface{}) *Model {
	model := m.getModel()
	if len(data) > 1 {
		m := make(map[string]interface{})
		for i := 0; i < len(data); i += 2 {
			m[gconv.String(data[i])] = data[i+1]
		}
		model.data = m
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
			rv := reflect.ValueOf(params)
			kind := rv.Kind()
			if kind == reflect.Ptr {
				rv = rv.Elem()
				kind = rv.Kind()
			}
			switch kind {
			case reflect.Slice, reflect.Array:
				list := make(List, rv.Len())
				for i := 0; i < rv.Len(); i++ {
					list[i] = varToMapDeep(rv.Index(i).Interface())
				}
				model.data = list
			case reflect.Map, reflect.Struct:
				model.data = varToMapDeep(data[0])
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
	return m.doInsertWithOption(gINSERT_OPTION_DEFAULT, data...)
}

// InsertIgnore does "INSERT IGNORE INTO ..." statement for the model.
// The optional parameter <data> is the same as the parameter of Model.Data function,
// see Model.Data.
func (m *Model) InsertIgnore(data ...interface{}) (result sql.Result, err error) {
	return m.doInsertWithOption(gINSERT_OPTION_IGNORE, data...)
}

// doInsertWithOption inserts data with option parameter.
func (m *Model) doInsertWithOption(option int, data ...interface{}) (result sql.Result, err error) {
	if len(data) > 0 {
		return m.Data(data...).Insert()
	}
	defer func() {
		if err == nil {
			m.checkAndRemoveCache()
		}
	}()
	if m.data == nil {
		return nil, errors.New("inserting into table with empty data")
	}
	if list, ok := m.data.(List); ok {
		// Batch insert.
		batch := 10
		if m.batch > 0 {
			batch = m.batch
		}
		return m.db.DoBatchInsert(
			m.getLink(true),
			m.tables,
			m.filterDataForInsertOrUpdate(list),
			option,
			batch,
		)
	} else if data, ok := m.data.(Map); ok {
		// Single insert.
		return m.db.DoInsert(
			m.getLink(true),
			m.tables,
			m.filterDataForInsertOrUpdate(data),
			option,
		)
	}
	return nil, errors.New("inserting into table with invalid data type")
}

// Replace does "REPLACE INTO ..." statement for the model.
// The optional parameter <data> is the same as the parameter of Model.Data function,
// see Model.Data.
func (m *Model) Replace(data ...interface{}) (result sql.Result, err error) {
	if len(data) > 0 {
		return m.Data(data...).Replace()
	}
	defer func() {
		if err == nil {
			m.checkAndRemoveCache()
		}
	}()
	if m.data == nil {
		return nil, errors.New("replacing into table with empty data")
	}
	if list, ok := m.data.(List); ok {
		// Batch replace.
		batch := 10
		if m.batch > 0 {
			batch = m.batch
		}
		return m.db.DoBatchInsert(
			m.getLink(true),
			m.tables,
			m.filterDataForInsertOrUpdate(list),
			gINSERT_OPTION_REPLACE,
			batch,
		)
	} else if data, ok := m.data.(Map); ok {
		// Single insert.
		return m.db.DoInsert(
			m.getLink(true),
			m.tables,
			m.filterDataForInsertOrUpdate(data),
			gINSERT_OPTION_REPLACE,
		)
	}
	return nil, errors.New("replacing into table with invalid data type")
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
	defer func() {
		if err == nil {
			m.checkAndRemoveCache()
		}
	}()
	if m.data == nil {
		return nil, errors.New("saving into table with empty data")
	}
	if list, ok := m.data.(List); ok {
		// Batch save.
		batch := gDEFAULT_BATCH_NUM
		if m.batch > 0 {
			batch = m.batch
		}
		return m.db.DoBatchInsert(
			m.getLink(true),
			m.tables,
			m.filterDataForInsertOrUpdate(list),
			gINSERT_OPTION_SAVE,
			batch,
		)
	} else if data, ok := m.data.(Map); ok {
		// Single save.
		return m.db.DoInsert(
			m.getLink(true),
			m.tables,
			m.filterDataForInsertOrUpdate(data),
			gINSERT_OPTION_SAVE,
		)
	}
	return nil, errors.New("saving into table with invalid data type")
}

// Update does "UPDATE ... " statement for the model.
//
// If the optional parameter <dataAndWhere> is given, the dataAndWhere[0] is the updated data field,
// and dataAndWhere[1:] is treated as where condition fields.
// Also see Model.Data and Model.Where functions.
func (m *Model) Update(dataAndWhere ...interface{}) (result sql.Result, err error) {
	if len(dataAndWhere) > 0 {
		if len(dataAndWhere) > 2 {
			return m.Data(dataAndWhere[0]).Where(dataAndWhere[1], dataAndWhere[2:]...).Update()
		} else if len(dataAndWhere) == 2 {
			return m.Data(dataAndWhere[0]).Where(dataAndWhere[1]).Update()
		} else {
			return m.Data(dataAndWhere[0]).Update()
		}
	}
	defer func() {
		if err == nil {
			m.checkAndRemoveCache()
		}
	}()
	if m.data == nil {
		return nil, errors.New("updating table with empty data")
	}
	condition, conditionArgs := m.formatCondition(false)
	return m.db.DoUpdate(
		m.getLink(true),
		m.tables,
		m.filterDataForInsertOrUpdate(m.data),
		condition,
		conditionArgs...,
	)
}

// Delete does "DELETE FROM ... " statement for the model.
// The optional parameter <where> is the same as the parameter of Model.Where function,
// see Model.Where.
func (m *Model) Delete(where ...interface{}) (result sql.Result, err error) {
	if len(where) > 0 {
		return m.Where(where[0], where[1:]...).Delete()
	}
	defer func() {
		if err == nil {
			m.checkAndRemoveCache()
		}
	}()
	condition, conditionArgs := m.formatCondition(false)
	return m.db.DoDelete(m.getLink(true), m.tables, condition, conditionArgs...)
}

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
	condition, conditionArgs := m.formatCondition(false)
	return m.getAll(fmt.Sprintf("SELECT %s FROM %s%s", m.fields, m.tables, condition), conditionArgs...)
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
	condition, conditionArgs := m.formatCondition(true)
	all, err := m.getAll(fmt.Sprintf("SELECT %s FROM %s%s", m.fields, m.tables, condition), conditionArgs...)
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
		return nil, err
	}
	for _, v := range one {
		return v, nil
	}
	return nil, nil
}

// Struct retrieves one record from table and converts it into given struct.
// The parameter <pointer> should be type of *struct/**struct. If type **struct is given,
// it can create the struct internally during converting.
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
func (m *Model) Struct(pointer interface{}, where ...interface{}) error {
	one, err := m.One(where...)
	if err != nil {
		return err
	}
	if len(one) == 0 {
		return sql.ErrNoRows
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
// from table.
//
// Eg:
// users := ([]User)(nil)
// err := db.Table("user").Structs(&users)
//
// users := ([]*User)(nil)
// err := db.Table("user").Structs(&users)
func (m *Model) Structs(pointer interface{}, where ...interface{}) error {
	all, err := m.All(where...)
	if err != nil {
		return err
	}
	if len(all) == 0 {
		return sql.ErrNoRows
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
// err := db.Table("user").Structs(&users)
//
// users := ([]*User)(nil)
// err := db.Table("user").Structs(&users)
func (m *Model) Scan(pointer interface{}, where ...interface{}) error {
	t := reflect.TypeOf(pointer)
	k := t.Kind()
	if k != reflect.Ptr {
		return fmt.Errorf("params should be type of pointer, but got: %v", k)
	}
	switch t.Elem().Kind() {
	case reflect.Array:
	case reflect.Slice:
		return m.Structs(pointer, where...)
	default:
		return m.Struct(pointer, where...)
	}
	return nil
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
		countFields = fmt.Sprintf(`COUNT(%s)`, m.fields)
	}
	condition, conditionArgs := m.formatCondition(false)
	s := fmt.Sprintf("SELECT %s FROM %s %s", countFields, m.tables, condition)
	if len(m.groupBy) > 0 {
		s = fmt.Sprintf("SELECT COUNT(1) FROM (%s) count_alias", s)
	}
	list, err := m.getAll(s, conditionArgs...)
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

// Chunk iterates the table with given size and callback function.
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

// filterDataForInsertOrUpdate does filter feature with data for inserting/updating operations.
// Note that, it does not filter list item, which is also type of map, for "omit empty" feature.
func (m *Model) filterDataForInsertOrUpdate(data interface{}) interface{} {
	if list, ok := m.data.(List); ok {
		for k, item := range list {
			list[k] = m.doFilterDataMapForInsertOrUpdate(item, false)
		}
		return list
	} else if item, ok := m.data.(Map); ok {
		return m.doFilterDataMapForInsertOrUpdate(item, true)
	}
	return data
}

// doFilterDataMapForInsertOrUpdate does the filter features for map.
// Note that, it does not filter list item, which is also type of map, for "omit empty" feature.
func (m *Model) doFilterDataMapForInsertOrUpdate(data Map, allowOmitEmpty bool) Map {
	if m.filter {
		data = m.db.filterFields(m.schema, m.tables, data)
	}
	// Remove key-value pairs of which the value is empty.
	if allowOmitEmpty && m.option&OPTION_OMITEMPTY > 0 {
		m := gmap.NewStrAnyMapFrom(data)
		m.FilterEmpty()
		data = m.Map()
	}

	if len(m.fields) > 0 && m.fields != "*" {
		// Keep specified fields.
		set := gset.NewStrSetFrom(gstr.SplitAndTrim(m.fields, ","))
		for k := range data {
			if !set.Contains(k) {
				delete(data, k)
			}
		}
	} else if len(m.fieldsEx) > 0 {
		// Filter specified fields.
		for _, v := range gstr.SplitAndTrim(m.fieldsEx, ",") {
			delete(data, v)
		}
	}
	return data
}

// getLink returns the underlying database link object with configured <linkType> attribute.
// The parameter <master> specifies whether using the master node if master-slave configured.
func (m *Model) getLink(master bool) dbLink {
	if m.tx != nil {
		return m.tx.tx
	}
	linkType := m.linkType
	if linkType == 0 {
		if master {
			linkType = gLINK_TYPE_MASTER
		} else {
			linkType = gLINK_TYPE_SLAVE
		}
	}
	switch linkType {
	case gLINK_TYPE_MASTER:
		link, err := m.db.GetMaster(m.schema)
		if err != nil {
			panic(err)
		}
		return link
	case gLINK_TYPE_SLAVE:
		link, err := m.db.GetSlave(m.schema)
		if err != nil {
			panic(err)
		}
		return link
	}
	return nil
}

// getAll does the query from database.
func (m *Model) getAll(query string, args ...interface{}) (result Result, err error) {
	cacheKey := ""
	// Retrieve from cache.
	if m.cacheEnabled {
		cacheKey = m.cacheName
		if len(cacheKey) == 0 {
			cacheKey = query + "/" + gconv.String(args)
		}
		if v := m.db.GetCache().Get(cacheKey); v != nil {
			return v.(Result), nil
		}
	}
	result, err = m.db.DoGetAll(m.getLink(false), query, args...)
	// Cache the result.
	if len(cacheKey) > 0 && err == nil {
		if m.cacheDuration < 0 {
			m.db.GetCache().Remove(cacheKey)
		} else {
			m.db.GetCache().Set(cacheKey, result, m.cacheDuration)
		}
	}
	return result, err
}

// getPrimaryKey retrieves and returns the primary key name of the model table.
// It parses m.tables to retrieve the primary table name, supporting m.tables like:
// "user", "user u", "user as u, user_detail as ud".
func (m *Model) getPrimaryKey() string {
	table := gstr.SplitAndTrim(m.tables, " ")[0]
	tableFields, err := m.db.TableFields(table)
	if err != nil {
		return ""
	}
	for name, field := range tableFields {
		if gstr.ContainsI(field.Key, "pri") {
			return name
		}
	}
	return ""
}

// checkAndRemoveCache checks and remove the cache if necessary.
func (m *Model) checkAndRemoveCache() {
	if m.cacheEnabled && m.cacheDuration < 0 && len(m.cacheName) > 0 {
		m.db.GetCache().Remove(m.cacheName)
	}
}

// formatCondition formats where arguments of the model and returns a new condition sql and its arguments.
// Note that this function does not change any attribute value of the <m>.
//
// The parameter <limit> specifies whether limits querying only one record if m.limit is not set.
func (m *Model) formatCondition(limit bool) (condition string, conditionArgs []interface{}) {
	var where string
	if len(m.whereHolder) > 0 {
		for _, v := range m.whereHolder {
			switch v.operator {
			case gWHERE_HOLDER_WHERE:
				if where == "" {
					newWhere, newArgs := formatWhere(m.db, v.where, v.args, m.option&OPTION_OMITEMPTY > 0)
					if len(newWhere) > 0 {
						where = newWhere
						conditionArgs = newArgs
					}
					continue
				}
				fallthrough

			case gWHERE_HOLDER_AND:
				newWhere, newArgs := formatWhere(m.db, v.where, v.args, m.option&OPTION_OMITEMPTY > 0)
				if len(newWhere) > 0 {
					if where[0] == '(' {
						where = fmt.Sprintf(`%s AND (%s)`, where, newWhere)
					} else {
						where = fmt.Sprintf(`(%s) AND (%s)`, where, newWhere)
					}
					conditionArgs = append(conditionArgs, newArgs...)
				}

			case gWHERE_HOLDER_OR:
				newWhere, newArgs := formatWhere(m.db, v.where, v.args, m.option&OPTION_OMITEMPTY > 0)
				if len(newWhere) > 0 {
					if where[0] == '(' {
						where = fmt.Sprintf(`%s OR (%s)`, where, newWhere)
					} else {
						where = fmt.Sprintf(`(%s) OR (%s)`, where, newWhere)
					}
					conditionArgs = append(conditionArgs, newArgs...)
				}
			}
		}
	}
	if where != "" {
		condition += " WHERE " + where
	}
	if m.groupBy != "" {
		condition += " GROUP BY " + m.groupBy
	}
	if m.orderBy != "" {
		condition += " ORDER BY " + m.orderBy
	}
	if m.limit != 0 {
		if m.start >= 0 {
			condition += fmt.Sprintf(" LIMIT %d,%d", m.start, m.limit)
		} else {
			condition += fmt.Sprintf(" LIMIT %d", m.limit)
		}
	} else if limit {
		condition += " LIMIT 1"
	}
	if m.offset >= 0 {
		condition += fmt.Sprintf(" OFFSET %d", m.offset)
	}
	return
}
