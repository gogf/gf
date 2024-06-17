// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package gdb

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/internal/reflection"
	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
)

// GetCore returns the underlying *Core object.
func (c *Core) GetCore() *Core {
	return c
}

// Ctx is a chaining function, which creates and returns a new DB that is a shallow copy
// of current DB object and with given context in it.
// Note that this returned DB object can be used only once, so do not assign it to
// a global or package variable for long using.
func (c *Core) Ctx(ctx context.Context) DB {
	if ctx == nil {
		return c.db
	}
	// It makes a shallow copy of current db and changes its context for next chaining operation.
	var (
		err        error
		newCore    = &Core{}
		configNode = c.db.GetConfig()
	)
	*newCore = *c
	// It creates a new DB object(NOT NEW CONNECTION), which is commonly a wrapper for object `Core`.
	newCore.db, err = driverMap[configNode.Type].New(newCore, configNode)
	if err != nil {
		// It is really a serious error here.
		// Do not let it continue.
		panic(err)
	}
	newCore.ctx = WithDB(ctx, newCore.db)
	newCore.ctx = c.injectInternalCtxData(newCore.ctx)
	return newCore.db
}

// GetCtx returns the context for current DB.
// It returns `context.Background()` is there's no context previously set.
func (c *Core) GetCtx() context.Context {
	ctx := c.ctx
	if ctx == nil {
		ctx = context.TODO()
	}
	return c.injectInternalCtxData(ctx)
}

// GetCtxTimeout returns the context and cancel function for specified timeout type.
func (c *Core) GetCtxTimeout(ctx context.Context, timeoutType int) (context.Context, context.CancelFunc) {
	if ctx == nil {
		ctx = c.db.GetCtx()
	} else {
		ctx = context.WithValue(ctx, "WrappedByGetCtxTimeout", nil)
	}
	switch timeoutType {
	case ctxTimeoutTypeExec:
		if c.db.GetConfig().ExecTimeout > 0 {
			return context.WithTimeout(ctx, c.db.GetConfig().ExecTimeout)
		}
	case ctxTimeoutTypeQuery:
		if c.db.GetConfig().QueryTimeout > 0 {
			return context.WithTimeout(ctx, c.db.GetConfig().QueryTimeout)
		}
	case ctxTimeoutTypePrepare:
		if c.db.GetConfig().PrepareTimeout > 0 {
			return context.WithTimeout(ctx, c.db.GetConfig().PrepareTimeout)
		}
	default:
		panic(gerror.NewCodef(gcode.CodeInvalidParameter, "invalid context timeout type: %d", timeoutType))
	}
	return ctx, func() {}
}

// Close closes the database and prevents new queries from starting.
// Close then waits for all queries that have started processing on the server
// to finish.
//
// It is rare to Close a DB, as the DB handle is meant to be
// long-lived and shared between many goroutines.
func (c *Core) Close(ctx context.Context) (err error) {
	if err = c.cache.Close(ctx); err != nil {
		return err
	}
	c.links.LockFunc(func(m map[any]any) {
		for k, v := range m {
			if db, ok := v.(*sql.DB); ok {
				err = db.Close()
				if err != nil {
					err = gerror.WrapCode(gcode.CodeDbOperationError, err, `db.Close failed`)
				}
				intlog.Printf(ctx, `close link: %s, err: %v`, k, err)
				if err != nil {
					return
				}
				delete(m, k)
			}
		}
	})
	return
}

// Master creates and returns a connection from master node if master-slave configured.
// It returns the default connection if master-slave not configured.
func (c *Core) Master(schema ...string) (*sql.DB, error) {
	var (
		usedSchema   = gutil.GetOrDefaultStr(c.schema, schema...)
		charL, charR = c.db.GetChars()
	)
	return c.getSqlDb(true, gstr.Trim(usedSchema, charL+charR))
}

// Slave creates and returns a connection from slave node if master-slave configured.
// It returns the default connection if master-slave not configured.
func (c *Core) Slave(schema ...string) (*sql.DB, error) {
	var (
		usedSchema   = gutil.GetOrDefaultStr(c.schema, schema...)
		charL, charR = c.db.GetChars()
	)
	return c.getSqlDb(false, gstr.Trim(usedSchema, charL+charR))
}

// GetAll queries and returns data records from database.
func (c *Core) GetAll(ctx context.Context, sql string, args ...interface{}) (Result, error) {
	return c.db.DoSelect(ctx, nil, sql, args...)
}

// DoSelect queries and returns data records from database.
func (c *Core) DoSelect(ctx context.Context, link Link, sql string, args ...interface{}) (result Result, err error) {
	return c.db.DoQuery(ctx, link, sql, args...)
}

// GetOne queries and returns one record from database.
func (c *Core) GetOne(ctx context.Context, sql string, args ...interface{}) (Record, error) {
	list, err := c.db.GetAll(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		return list[0], nil
	}
	return nil, nil
}

// GetArray queries and returns data values as slice from database.
// Note that if there are multiple columns in the result, it returns just one column values randomly.
func (c *Core) GetArray(ctx context.Context, sql string, args ...interface{}) ([]Value, error) {
	all, err := c.db.DoSelect(ctx, nil, sql, args...)
	if err != nil {
		return nil, err
	}
	return all.Array(), nil
}

// doGetStruct queries one record from database and converts it to given struct.
// The parameter `pointer` should be a pointer to struct.
func (c *Core) doGetStruct(ctx context.Context, pointer interface{}, sql string, args ...interface{}) error {
	one, err := c.db.GetOne(ctx, sql, args...)
	if err != nil {
		return err
	}
	return one.Struct(pointer)
}

// doGetStructs queries records from database and converts them to given struct.
// The parameter `pointer` should be type of struct slice: []struct/[]*struct.
func (c *Core) doGetStructs(ctx context.Context, pointer interface{}, sql string, args ...interface{}) error {
	all, err := c.db.GetAll(ctx, sql, args...)
	if err != nil {
		return err
	}
	return all.Structs(pointer)
}

// GetScan queries one or more records from database and converts them to given struct or
// struct array.
//
// If parameter `pointer` is type of struct pointer, it calls GetStruct internally for
// the conversion. If parameter `pointer` is type of slice, it calls GetStructs internally
// for conversion.
func (c *Core) GetScan(ctx context.Context, pointer interface{}, sql string, args ...interface{}) error {
	reflectInfo := reflection.OriginTypeAndKind(pointer)
	if reflectInfo.InputKind != reflect.Ptr {
		return gerror.NewCodef(
			gcode.CodeInvalidParameter,
			"params should be type of pointer, but got: %v",
			reflectInfo.InputKind,
		)
	}
	switch reflectInfo.OriginKind {
	case reflect.Array, reflect.Slice:
		return c.db.GetCore().doGetStructs(ctx, pointer, sql, args...)

	case reflect.Struct:
		return c.db.GetCore().doGetStruct(ctx, pointer, sql, args...)
	}
	return gerror.NewCodef(
		gcode.CodeInvalidParameter,
		`in valid parameter type "%v", of which element type should be type of struct/slice`,
		reflectInfo.InputType,
	)
}

// GetValue queries and returns the field value from database.
// The sql should query only one field from database, or else it returns only one
// field of the result.
func (c *Core) GetValue(ctx context.Context, sql string, args ...interface{}) (Value, error) {
	one, err := c.db.GetOne(ctx, sql, args...)
	if err != nil {
		return gvar.New(nil), err
	}
	for _, v := range one {
		return v, nil
	}
	return gvar.New(nil), nil
}

// GetCount queries and returns the count from database.
func (c *Core) GetCount(ctx context.Context, sql string, args ...interface{}) (int, error) {
	// If the query fields do not contain function "COUNT",
	// it replaces the sql string and adds the "COUNT" function to the fields.
	if !gregex.IsMatchString(`(?i)SELECT\s+COUNT\(.+\)\s+FROM`, sql) {
		sql, _ = gregex.ReplaceString(`(?i)(SELECT)\s+(.+)\s+(FROM)`, `$1 COUNT($2) $3`, sql)
	}
	value, err := c.db.GetValue(ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	return value.Int(), nil
}

// Union does "(SELECT xxx FROM xxx) UNION (SELECT xxx FROM xxx) ..." statement.
func (c *Core) Union(unions ...*Model) *Model {
	var ctx = c.db.GetCtx()
	return c.doUnion(ctx, unionTypeNormal, unions...)
}

// UnionAll does "(SELECT xxx FROM xxx) UNION ALL (SELECT xxx FROM xxx) ..." statement.
func (c *Core) UnionAll(unions ...*Model) *Model {
	var ctx = c.db.GetCtx()
	return c.doUnion(ctx, unionTypeAll, unions...)
}

func (c *Core) doUnion(ctx context.Context, unionType int, unions ...*Model) *Model {
	var (
		unionTypeStr   string
		composedSqlStr string
		composedArgs   = make([]interface{}, 0)
	)
	if unionType == unionTypeAll {
		unionTypeStr = "UNION ALL"
	} else {
		unionTypeStr = "UNION"
	}
	for _, v := range unions {
		sqlWithHolder, holderArgs := v.getFormattedSqlAndArgs(ctx, queryTypeNormal, false)
		if composedSqlStr == "" {
			composedSqlStr += fmt.Sprintf(`(%s)`, sqlWithHolder)
		} else {
			composedSqlStr += fmt.Sprintf(` %s (%s)`, unionTypeStr, sqlWithHolder)
		}
		composedArgs = append(composedArgs, holderArgs...)
	}
	return c.db.Raw(composedSqlStr, composedArgs...)
}

// PingMaster pings the master node to check authentication or keeps the connection alive.
func (c *Core) PingMaster() error {
	var ctx = c.db.GetCtx()
	if master, err := c.db.Master(); err != nil {
		return err
	} else {
		if err = master.PingContext(ctx); err != nil {
			err = gerror.WrapCode(gcode.CodeDbOperationError, err, `master.Ping failed`)
		}
		return err
	}
}

// PingSlave pings the slave node to check authentication or keeps the connection alive.
func (c *Core) PingSlave() error {
	var ctx = c.db.GetCtx()
	if slave, err := c.db.Slave(); err != nil {
		return err
	} else {
		if err = slave.PingContext(ctx); err != nil {
			err = gerror.WrapCode(gcode.CodeDbOperationError, err, `slave.Ping failed`)
		}
		return err
	}
}

// Insert does "INSERT INTO ..." statement for the table.
// If there's already one unique record of the data in the table, it returns error.
//
// The parameter `data` can be type of map/gmap/struct/*struct/[]map/[]struct, etc.
// Eg:
// Data(g.Map{"uid": 10000, "name":"john"})
// Data(g.Slice{g.Map{"uid": 10000, "name":"john"}, g.Map{"uid": 20000, "name":"smith"})
//
// The parameter `batch` specifies the batch operation count when given data is slice.
func (c *Core) Insert(ctx context.Context, table string, data interface{}, batch ...int) (sql.Result, error) {
	if len(batch) > 0 {
		return c.Model(table).Ctx(ctx).Data(data).Batch(batch[0]).Insert()
	}
	return c.Model(table).Ctx(ctx).Data(data).Insert()
}

// InsertIgnore does "INSERT IGNORE INTO ..." statement for the table.
// If there's already one unique record of the data in the table, it ignores the inserting.
//
// The parameter `data` can be type of map/gmap/struct/*struct/[]map/[]struct, etc.
// Eg:
// Data(g.Map{"uid": 10000, "name":"john"})
// Data(g.Slice{g.Map{"uid": 10000, "name":"john"}, g.Map{"uid": 20000, "name":"smith"})
//
// The parameter `batch` specifies the batch operation count when given data is slice.
func (c *Core) InsertIgnore(ctx context.Context, table string, data interface{}, batch ...int) (sql.Result, error) {
	if len(batch) > 0 {
		return c.Model(table).Ctx(ctx).Data(data).Batch(batch[0]).InsertIgnore()
	}
	return c.Model(table).Ctx(ctx).Data(data).InsertIgnore()
}

// InsertAndGetId performs action Insert and returns the last insert id that automatically generated.
func (c *Core) InsertAndGetId(ctx context.Context, table string, data interface{}, batch ...int) (int64, error) {
	if len(batch) > 0 {
		return c.Model(table).Ctx(ctx).Data(data).Batch(batch[0]).InsertAndGetId()
	}
	return c.Model(table).Ctx(ctx).Data(data).InsertAndGetId()
}

// Replace does "REPLACE INTO ..." statement for the table.
// If there's already one unique record of the data in the table, it deletes the record
// and inserts a new one.
//
// The parameter `data` can be type of map/gmap/struct/*struct/[]map/[]struct, etc.
// Eg:
// Data(g.Map{"uid": 10000, "name":"john"})
// Data(g.Slice{g.Map{"uid": 10000, "name":"john"}, g.Map{"uid": 20000, "name":"smith"})
//
// The parameter `data` can be type of map/gmap/struct/*struct/[]map/[]struct, etc.
// If given data is type of slice, it then does batch replacing, and the optional parameter
// `batch` specifies the batch operation count.
func (c *Core) Replace(ctx context.Context, table string, data interface{}, batch ...int) (sql.Result, error) {
	if len(batch) > 0 {
		return c.Model(table).Ctx(ctx).Data(data).Batch(batch[0]).Replace()
	}
	return c.Model(table).Ctx(ctx).Data(data).Replace()
}

// Save does "INSERT INTO ... ON DUPLICATE KEY UPDATE..." statement for the table.
// It updates the record if there's primary or unique index in the saving data,
// or else it inserts a new record into the table.
//
// The parameter `data` can be type of map/gmap/struct/*struct/[]map/[]struct, etc.
// Eg:
// Data(g.Map{"uid": 10000, "name":"john"})
// Data(g.Slice{g.Map{"uid": 10000, "name":"john"}, g.Map{"uid": 20000, "name":"smith"})
//
// If given data is type of slice, it then does batch saving, and the optional parameter
// `batch` specifies the batch operation count.
func (c *Core) Save(ctx context.Context, table string, data interface{}, batch ...int) (sql.Result, error) {
	if len(batch) > 0 {
		return c.Model(table).Ctx(ctx).Data(data).Batch(batch[0]).Save()
	}
	return c.Model(table).Ctx(ctx).Data(data).Save()
}

func (c *Core) fieldsToSequence(ctx context.Context, table string, fields []string) ([]string, error) {
	var (
		fieldSet               = gset.NewStrSetFrom(fields)
		fieldsResultInSequence = make([]string, 0)
		tableFields, err       = c.db.TableFields(ctx, table)
	)
	if err != nil {
		return nil, err
	}
	// Sort the fields in order.
	var fieldsOfTableInSequence = make([]string, len(tableFields))
	for _, field := range tableFields {
		fieldsOfTableInSequence[field.Index] = field.Name
	}
	// Sort the input fields.
	for _, fieldName := range fieldsOfTableInSequence {
		if fieldSet.Contains(fieldName) {
			fieldsResultInSequence = append(fieldsResultInSequence, fieldName)
		}
	}
	return fieldsResultInSequence, nil
}

// DoInsert inserts or updates data for given table.
// This function is usually used for custom interface definition, you do not need call it manually.
// The parameter `data` can be type of map/gmap/struct/*struct/[]map/[]struct, etc.
// Eg:
// Data(g.Map{"uid": 10000, "name":"john"})
// Data(g.Slice{g.Map{"uid": 10000, "name":"john"}, g.Map{"uid": 20000, "name":"smith"})
//
// The parameter `option` values are as follows:
// InsertOptionDefault: just insert, if there's unique/primary key in the data, it returns error;
// InsertOptionReplace: if there's unique/primary key in the data, it deletes it from table and inserts a new one;
// InsertOptionSave:    if there's unique/primary key in the data, it updates it or else inserts a new one;
// InsertOptionIgnore:  if there's unique/primary key in the data, it ignores the inserting;
func (c *Core) DoInsert(ctx context.Context, link Link, table string, list List, option DoInsertOption) (result sql.Result, err error) {
	var (
		keys           []string      // Field names.
		values         []string      // Value holder string array, like: (?,?,?)
		params         []interface{} // Values that will be committed to underlying database driver.
		onDuplicateStr string        // onDuplicateStr is used in "ON DUPLICATE KEY UPDATE" statement.
	)
	// ============================================================================================
	// Group the list by fields. Different fields to different list.
	// It here uses ListMap to keep sequence for data inserting.
	// ============================================================================================
	var keyListMap = gmap.NewListMap()
	for _, item := range list {
		var (
			tmpKeys              = make([]string, 0)
			tmpKeysInSequenceStr string
		)
		for k := range item {
			tmpKeys = append(tmpKeys, k)
		}
		keys, err = c.fieldsToSequence(ctx, table, tmpKeys)
		if err != nil {
			return nil, err
		}
		tmpKeysInSequenceStr = gstr.Join(keys, ",")
		if !keyListMap.Contains(tmpKeysInSequenceStr) {
			keyListMap.Set(tmpKeysInSequenceStr, make(List, 0))
		}
		tmpKeysInSequenceList := keyListMap.Get(tmpKeysInSequenceStr).(List)
		tmpKeysInSequenceList = append(tmpKeysInSequenceList, item)
		keyListMap.Set(tmpKeysInSequenceStr, tmpKeysInSequenceList)
	}
	if keyListMap.Size() > 1 {
		var (
			tmpResult    sql.Result
			sqlResult    SqlResult
			rowsAffected int64
		)
		keyListMap.Iterator(func(key, value interface{}) bool {
			tmpResult, err = c.DoInsert(ctx, link, table, value.(List), option)
			if err != nil {
				return false
			}
			rowsAffected, err = tmpResult.RowsAffected()
			if err != nil {
				return false
			}
			sqlResult.Result = tmpResult
			sqlResult.Affected += rowsAffected
			return true
		})
		return &sqlResult, err
	}

	// Prepare the batch result pointer.
	var (
		charL, charR = c.db.GetChars()
		batchResult  = new(SqlResult)
		keysStr      = charL + strings.Join(keys, charR+","+charL) + charR
		operation    = GetInsertOperationByOption(option.InsertOption)
	)
	// Upsert clause only takes effect on Save operation.
	if option.InsertOption == InsertOptionSave {
		onDuplicateStr, err = c.db.FormatUpsert(keys, list, option)
		if err != nil {
			return nil, err
		}
	}
	var (
		listLength   = len(list)
		valueHolders = make([]string, 0)
	)
	for i := 0; i < listLength; i++ {
		values = values[:0]
		// Note that the map type is unordered,
		// so it should use slice+key to retrieve the value.
		for _, k := range keys {
			if s, ok := list[i][k].(Raw); ok {
				values = append(values, gconv.String(s))
			} else {
				values = append(values, "?")
				params = append(params, list[i][k])
			}
		}
		valueHolders = append(valueHolders, "("+gstr.Join(values, ",")+")")
		// Batch package checks: It meets the batch number, or it is the last element.
		if len(valueHolders) == option.BatchCount || (i == listLength-1 && len(valueHolders) > 0) {
			var (
				stdSqlResult sql.Result
				affectedRows int64
			)
			stdSqlResult, err = c.db.DoExec(ctx, link, fmt.Sprintf(
				"%s INTO %s(%s) VALUES%s %s",
				operation, c.QuotePrefixTableName(table), keysStr,
				gstr.Join(valueHolders, ","),
				onDuplicateStr,
			), params...)
			if err != nil {
				return stdSqlResult, err
			}
			if affectedRows, err = stdSqlResult.RowsAffected(); err != nil {
				err = gerror.WrapCode(gcode.CodeDbOperationError, err, `sql.Result.RowsAffected failed`)
				return stdSqlResult, err
			} else {
				batchResult.Result = stdSqlResult
				batchResult.Affected += affectedRows
			}
			params = params[:0]
			valueHolders = valueHolders[:0]
		}
	}
	return batchResult, nil
}

// Update does "UPDATE ... " statement for the table.
//
// The parameter `data` can be type of string/map/gmap/struct/*struct, etc.
// Eg: "uid=10000", "uid", 10000, g.Map{"uid": 10000, "name":"john"}
//
// The parameter `condition` can be type of string/map/gmap/slice/struct/*struct, etc.
// It is commonly used with parameter `args`.
// Eg:
// "uid=10000",
// "uid", 10000
// "money>? AND name like ?", 99999, "vip_%"
// "status IN (?)", g.Slice{1,2,3}
// "age IN(?,?)", 18, 50
// User{ Id : 1, UserName : "john"}.
func (c *Core) Update(ctx context.Context, table string, data interface{}, condition interface{}, args ...interface{}) (sql.Result, error) {
	return c.Model(table).Ctx(ctx).Data(data).Where(condition, args...).Update()
}

// DoUpdate does "UPDATE ... " statement for the table.
// This function is usually used for custom interface definition, you do not need to call it manually.
func (c *Core) DoUpdate(ctx context.Context, link Link, table string, data interface{}, condition string, args ...interface{}) (result sql.Result, err error) {
	table = c.QuotePrefixTableName(table)
	var (
		rv   = reflect.ValueOf(data)
		kind = rv.Kind()
	)
	if kind == reflect.Ptr {
		rv = rv.Elem()
		kind = rv.Kind()
	}
	var (
		params  []interface{}
		updates string
	)
	switch kind {
	case reflect.Map, reflect.Struct:
		var (
			fields         []string
			dataMap        map[string]interface{}
			counterHandler = func(column string, counter Counter) {
				if counter.Value != 0 {
					column = c.QuoteWord(column)
					var (
						columnRef = c.QuoteWord(counter.Field)
						columnVal = counter.Value
						operator  = "+"
					)
					if columnVal < 0 {
						operator = "-"
						columnVal = -columnVal
					}
					fields = append(fields, fmt.Sprintf("%s=%s%s?", column, columnRef, operator))
					params = append(params, columnVal)
				}
			}
		)
		dataMap, err = c.ConvertDataForRecord(ctx, data, table)
		if err != nil {
			return nil, err
		}
		// Sort the data keys in sequence of table fields.
		var (
			dataKeys       = make([]string, 0)
			keysInSequence = make([]string, 0)
		)
		for k := range dataMap {
			dataKeys = append(dataKeys, k)
		}
		keysInSequence, err = c.fieldsToSequence(ctx, table, dataKeys)
		if err != nil {
			return nil, err
		}
		for _, k := range keysInSequence {
			v := dataMap[k]
			switch value := v.(type) {
			case *Counter:
				counterHandler(k, *value)

			case Counter:
				counterHandler(k, value)

			default:
				if s, ok := v.(Raw); ok {
					fields = append(fields, c.QuoteWord(k)+"="+gconv.String(s))
				} else {
					fields = append(fields, c.QuoteWord(k)+"=?")
					params = append(params, v)
				}
			}
		}
		updates = strings.Join(fields, ",")

	default:
		updates = gconv.String(data)
	}
	if len(updates) == 0 {
		return nil, gerror.NewCode(gcode.CodeMissingParameter, "data cannot be empty")
	}
	if len(params) > 0 {
		args = append(params, args...)
	}
	// If no link passed, it then uses the master link.
	if link == nil {
		if link, err = c.MasterLink(); err != nil {
			return nil, err
		}
	}
	return c.db.DoExec(ctx, link, fmt.Sprintf(
		"UPDATE %s SET %s%s",
		table, updates, condition,
	),
		args...,
	)
}

// Delete does "DELETE FROM ... " statement for the table.
//
// The parameter `condition` can be type of string/map/gmap/slice/struct/*struct, etc.
// It is commonly used with parameter `args`.
// Eg:
// "uid=10000",
// "uid", 10000
// "money>? AND name like ?", 99999, "vip_%"
// "status IN (?)", g.Slice{1,2,3}
// "age IN(?,?)", 18, 50
// User{ Id : 1, UserName : "john"}.
func (c *Core) Delete(ctx context.Context, table string, condition interface{}, args ...interface{}) (result sql.Result, err error) {
	return c.Model(table).Ctx(ctx).Where(condition, args...).Delete()
}

// DoDelete does "DELETE FROM ... " statement for the table.
// This function is usually used for custom interface definition, you do not need call it manually.
func (c *Core) DoDelete(ctx context.Context, link Link, table string, condition string, args ...interface{}) (result sql.Result, err error) {
	if link == nil {
		if link, err = c.MasterLink(); err != nil {
			return nil, err
		}
	}
	table = c.QuotePrefixTableName(table)
	return c.db.DoExec(ctx, link, fmt.Sprintf("DELETE FROM %s%s", table, condition), args...)
}

// FilteredLink retrieves and returns filtered `linkInfo` that can be using for
// logging or tracing purpose.
func (c *Core) FilteredLink() string {
	return fmt.Sprintf(
		`%s@%s(%s:%s)/%s`,
		c.config.User, c.config.Protocol, c.config.Host, c.config.Port, c.config.Name,
	)
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
// It just returns the pointer address.
//
// Note that this interface implements mainly for workaround for a json infinite loop bug
// of Golang version < v1.14.
func (c Core) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`%+v`, c)), nil
}

// writeSqlToLogger outputs the Sql object to logger.
// It is enabled only if configuration "debug" is true.
func (c *Core) writeSqlToLogger(ctx context.Context, sql *Sql) {
	var transactionIdStr string
	if sql.IsTransaction {
		if v := ctx.Value(transactionIdForLoggerCtx); v != nil {
			transactionIdStr = fmt.Sprintf(`[txid:%d] `, v.(uint64))
		}
	}
	s := fmt.Sprintf(
		"[%3d ms] [%s] [%s] [rows:%-3d] %s%s",
		sql.End-sql.Start, sql.Group, sql.Schema, sql.RowsAffected, transactionIdStr, sql.Format,
	)
	if sql.Error != nil {
		s += "\nError: " + sql.Error.Error()
		c.logger.Error(ctx, s)
	} else {
		c.logger.Debug(ctx, s)
	}
}

// HasTable determine whether the table name exists in the database.
func (c *Core) HasTable(name string) (bool, error) {
	tables, err := c.GetTablesWithCache()
	if err != nil {
		return false, err
	}
	charL, charR := c.db.GetChars()
	name = gstr.Trim(name, charL+charR)
	for _, table := range tables {
		if table == name {
			return true, nil
		}
	}
	return false, nil
}

func (c *Core) GetInnerMemCache() *gcache.Cache {
	return c.innerMemCache
}

// GetTablesWithCache retrieves and returns the table names of current database with cache.
func (c *Core) GetTablesWithCache() ([]string, error) {
	var (
		ctx           = c.db.GetCtx()
		cacheKey      = fmt.Sprintf(`Tables:%s`, c.db.GetGroup())
		cacheDuration = gcache.DurationNoExpire
		innerMemCache = c.GetInnerMemCache()
	)
	result, err := innerMemCache.GetOrSetFuncLock(
		ctx, cacheKey,
		func(ctx context.Context) (interface{}, error) {
			tableList, err := c.db.Tables(ctx)
			if err != nil {
				return nil, err
			}
			return tableList, nil
		}, cacheDuration,
	)
	if err != nil {
		return nil, err
	}
	return result.Strings(), nil
}

// IsSoftCreatedFieldName checks and returns whether given field name is an automatic-filled created time.
func (c *Core) IsSoftCreatedFieldName(fieldName string) bool {
	if fieldName == "" {
		return false
	}
	if config := c.db.GetConfig(); config.CreatedAt != "" {
		if utils.EqualFoldWithoutChars(fieldName, config.CreatedAt) {
			return true
		}
		return gstr.InArray(append([]string{config.CreatedAt}, createdFieldNames...), fieldName)
	}
	for _, v := range createdFieldNames {
		if utils.EqualFoldWithoutChars(fieldName, v) {
			return true
		}
	}
	return false
}

// FormatSqlBeforeExecuting formats the sql string and its arguments before executing.
// The internal handleArguments function might be called twice during the SQL procedure,
// but do not worry about it, it's safe and efficient.
func (c *Core) FormatSqlBeforeExecuting(sql string, args []interface{}) (newSql string, newArgs []interface{}) {
	// DO NOT do this as there may be multiple lines and comments in the sql.
	// sql = gstr.Trim(sql)
	// sql = gstr.Replace(sql, "\n", " ")
	// sql, _ = gregex.ReplaceString(`\s{2,}`, ` `, sql)
	return handleSliceAndStructArgsForSql(sql, args)
}
