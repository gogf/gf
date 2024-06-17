// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/encoding/ghash"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/internal/empty"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/internal/reflection"
	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/os/gstructs"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gmeta"
	"github.com/gogf/gf/v2/util/gtag"
	"github.com/gogf/gf/v2/util/gutil"
)

// iString is the type assert api for String.
type iString interface {
	String() string
}

// iIterator is the type assert api for Iterator.
type iIterator interface {
	Iterator(f func(key, value interface{}) bool)
}

// iInterfaces is the type assert api for Interfaces.
type iInterfaces interface {
	Interfaces() []interface{}
}

// iNil if the type assert api for IsNil.
type iNil interface {
	IsNil() bool
}

// iTableName is the interface for retrieving table name for struct.
type iTableName interface {
	TableName() string
}

const (
	OrmTagForStruct       = "orm"
	OrmTagForTable        = "table"
	OrmTagForWith         = "with"
	OrmTagForWithWhere    = "where"
	OrmTagForWithOrder    = "order"
	OrmTagForWithUnscoped = "unscoped"
	OrmTagForDo           = "do"
)

var (
	// quoteWordReg is the regular expression object for a word check.
	quoteWordReg = regexp.MustCompile(`^[a-zA-Z0-9\-_]+$`)

	// structTagPriority tags for struct converting for orm field mapping.
	structTagPriority = append([]string{OrmTagForStruct}, gtag.StructTagPriority...)
)

// WithDB injects given db object into context and returns a new context.
func WithDB(ctx context.Context, db DB) context.Context {
	if db == nil {
		return ctx
	}
	dbCtx := db.GetCtx()
	if ctxDb := DBFromCtx(dbCtx); ctxDb != nil {
		return dbCtx
	}
	ctx = context.WithValue(ctx, ctxKeyForDB, db)
	return ctx
}

// DBFromCtx retrieves and returns DB object from context.
func DBFromCtx(ctx context.Context) DB {
	if ctx == nil {
		return nil
	}
	v := ctx.Value(ctxKeyForDB)
	if v != nil {
		return v.(DB)
	}
	return nil
}

// ToSQL formats and returns the last one of sql statements in given closure function
// WITHOUT TRULY EXECUTING IT.
// Be caution that, all the following sql statements should use the context object passing by function `f`.
func ToSQL(ctx context.Context, f func(ctx context.Context) error) (sql string, err error) {
	var manager = &CatchSQLManager{
		SQLArray: garray.NewStrArray(),
		DoCommit: false,
	}
	ctx = context.WithValue(ctx, ctxKeyCatchSQL, manager)
	err = f(ctx)
	sql, _ = manager.SQLArray.PopRight()
	return
}

// CatchSQL catches and returns all sql statements that are EXECUTED in given closure function.
// Be caution that, all the following sql statements should use the context object passing by function `f`.
func CatchSQL(ctx context.Context, f func(ctx context.Context) error) (sqlArray []string, err error) {
	var manager = &CatchSQLManager{
		SQLArray: garray.NewStrArray(),
		DoCommit: true,
	}
	ctx = context.WithValue(ctx, ctxKeyCatchSQL, manager)
	err = f(ctx)
	return manager.SQLArray.Slice(), err
}

// isDoStruct checks and returns whether given type is a DO struct.
func isDoStruct(object interface{}) bool {
	// It checks by struct name like "XxxForDao", to be compatible with old version.
	// TODO remove this compatible codes in future.
	reflectType := reflect.TypeOf(object)
	if gstr.HasSuffix(reflectType.String(), modelForDaoSuffix) {
		return true
	}
	// It checks by struct meta for DO struct in version.
	if ormTag := gmeta.Get(object, OrmTagForStruct); !ormTag.IsEmpty() {
		match, _ := gregex.MatchString(
			fmt.Sprintf(`%s\s*:\s*([^,]+)`, OrmTagForDo),
			ormTag.String(),
		)
		if len(match) > 1 {
			return gconv.Bool(match[1])
		}
	}
	return false
}

// getTableNameFromOrmTag retrieves and returns the table name from struct object.
func getTableNameFromOrmTag(object interface{}) string {
	var tableName string
	// Use the interface value.
	if r, ok := object.(iTableName); ok {
		tableName = r.TableName()
	}
	// User meta data tag "orm".
	if tableName == "" {
		if ormTag := gmeta.Get(object, OrmTagForStruct); !ormTag.IsEmpty() {
			match, _ := gregex.MatchString(
				fmt.Sprintf(`%s\s*:\s*([^,]+)`, OrmTagForTable),
				ormTag.String(),
			)
			if len(match) > 1 {
				tableName = match[1]
			}
		}
	}
	// Use the struct name of snake case.
	if tableName == "" {
		if t, err := gstructs.StructType(object); err != nil {
			panic(err)
		} else {
			tableName = gstr.CaseSnakeFirstUpper(
				gstr.StrEx(t.String(), "."),
			)
		}
	}
	return tableName
}

// ListItemValues retrieves and returns the elements of all item struct/map with key `key`.
// Note that the parameter `list` should be type of slice which contains elements of map or struct,
// or else it returns an empty slice.
//
// The parameter `list` supports types like:
// []map[string]interface{}
// []map[string]sub-map
// []struct
// []struct:sub-struct
// Note that the sub-map/sub-struct makes sense only if the optional parameter `subKey` is given.
// See gutil.ListItemValues.
func ListItemValues(list interface{}, key interface{}, subKey ...interface{}) (values []interface{}) {
	return gutil.ListItemValues(list, key, subKey...)
}

// ListItemValuesUnique retrieves and returns the unique elements of all struct/map with key `key`.
// Note that the parameter `list` should be type of slice which contains elements of map or struct,
// or else it returns an empty slice.
// See gutil.ListItemValuesUnique.
func ListItemValuesUnique(list interface{}, key string, subKey ...interface{}) []interface{} {
	return gutil.ListItemValuesUnique(list, key, subKey...)
}

// GetInsertOperationByOption returns proper insert option with given parameter `option`.
func GetInsertOperationByOption(option InsertOption) string {
	var operator string
	switch option {
	case InsertOptionReplace:
		operator = InsertOperationReplace
	case InsertOptionIgnore:
		operator = InsertOperationIgnore
	default:
		operator = InsertOperationInsert
	}
	return operator
}

func anyValueToMapBeforeToRecord(value interface{}) map[string]interface{} {
	convertedMap := gconv.Map(value, gconv.MapOption{
		Tags:      structTagPriority,
		OmitEmpty: true, // To be compatible with old version from v2.6.0.
	})
	if gutil.OriginValueAndKind(value).OriginKind != reflect.Struct {
		return convertedMap
	}
	// It here converts all struct/map slice attributes to json string.
	for k, v := range convertedMap {
		originValueAndKind := gutil.OriginValueAndKind(v)
		switch originValueAndKind.OriginKind {
		// Check map item slice item.
		case reflect.Array, reflect.Slice:
			mapItemValue := originValueAndKind.OriginValue
			if mapItemValue.Len() == 0 {
				break
			}
			// Check slice item type struct/map type.
			switch mapItemValue.Index(0).Kind() {
			case reflect.Struct, reflect.Map:
				mapItemJsonBytes, err := json.Marshal(v)
				if err != nil {
					// Do not eat any error.
					intlog.Error(context.TODO(), err)
				}
				convertedMap[k] = mapItemJsonBytes
			}
		}
	}
	return convertedMap
}

// MapOrStructToMapDeep converts `value` to map type recursively(if attribute struct is embedded).
// The parameter `value` should be type of *map/map/*struct/struct.
// It supports embedded struct definition for struct.
func MapOrStructToMapDeep(value interface{}, omitempty bool) map[string]interface{} {
	m := gconv.Map(value, gconv.MapOption{
		Tags:      structTagPriority,
		OmitEmpty: omitempty,
	})
	for k, v := range m {
		switch v.(type) {
		case time.Time, *time.Time, gtime.Time, *gtime.Time, gjson.Json, *gjson.Json:
			m[k] = v
		}
	}
	return m
}

// doQuoteTableName adds prefix string and quote chars for table name. It handles table string like:
// "user", "user u", "user,user_detail", "user u, user_detail ut", "user as u, user_detail as ut",
// "user.user u", "`user`.`user` u".
//
// Note that, this will automatically check the table prefix whether already added, if true it does
// nothing to the table name, or else adds the prefix to the table name and returns new table name with prefix.
func doQuoteTableName(table, prefix, charLeft, charRight string) string {
	var (
		index  int
		chars  = charLeft + charRight
		array1 = gstr.SplitAndTrim(table, ",")
	)
	for k1, v1 := range array1 {
		array2 := gstr.SplitAndTrim(v1, " ")
		// Trim the security chars.
		array2[0] = gstr.Trim(array2[0], chars)
		// Check whether it has database name.
		array3 := gstr.Split(gstr.Trim(array2[0]), ".")
		for k, v := range array3 {
			array3[k] = gstr.Trim(v, chars)
		}
		index = len(array3) - 1
		// If the table name already has the prefix, skips the prefix adding.
		if len(array3[index]) <= len(prefix) || array3[index][:len(prefix)] != prefix {
			array3[index] = prefix + array3[index]
		}
		array2[0] = gstr.Join(array3, ".")
		// Add the security chars.
		array2[0] = doQuoteString(array2[0], charLeft, charRight)
		array1[k1] = gstr.Join(array2, " ")
	}
	return gstr.Join(array1, ",")
}

// doQuoteWord checks given string `s` a word, if true quotes it with `charLeft` and `charRight`
// and returns the quoted string; or else returns `s` without any change.
func doQuoteWord(s, charLeft, charRight string) string {
	if quoteWordReg.MatchString(s) && !gstr.ContainsAny(s, charLeft+charRight) {
		return charLeft + s + charRight
	}
	return s
}

// doQuoteString quotes string with quote chars.
// For example, if quote char is '`':
// "null"                             => "NULL"
// "user"                             => "`user`"
// "user u"                           => "`user` u"
// "user,user_detail"                 => "`user`,`user_detail`"
// "user u, user_detail ut"           => "`user` u,`user_detail` ut"
// "user.user u, user.user_detail ut" => "`user`.`user` u,`user`.`user_detail` ut"
// "u.id, u.name, u.age"              => "`u`.`id`,`u`.`name`,`u`.`age`"
// "u.id asc"                         => "`u`.`id` asc".
func doQuoteString(s, charLeft, charRight string) string {
	array1 := gstr.SplitAndTrim(s, ",")
	for k1, v1 := range array1 {
		array2 := gstr.SplitAndTrim(v1, " ")
		array3 := gstr.Split(gstr.Trim(array2[0]), ".")
		if len(array3) == 1 {
			if strings.EqualFold(array3[0], "NULL") {
				array3[0] = doQuoteWord(array3[0], "", "")
			} else {
				array3[0] = doQuoteWord(array3[0], charLeft, charRight)
			}
		} else if len(array3) >= 2 {
			array3[0] = doQuoteWord(array3[0], charLeft, charRight)
			// Note:
			// mysql: u.uid
			// mssql double dots: Database..Table
			array3[len(array3)-1] = doQuoteWord(array3[len(array3)-1], charLeft, charRight)
		}
		array2[0] = gstr.Join(array3, ".")
		array1[k1] = gstr.Join(array2, " ")
	}
	return gstr.Join(array1, ",")
}

func getFieldsFromStructOrMap(structOrMap interface{}) (fields []string) {
	fields = []string{}
	if utils.IsStruct(structOrMap) {
		structFields, _ := gstructs.Fields(gstructs.FieldsInput{
			Pointer:         structOrMap,
			RecursiveOption: gstructs.RecursiveOptionEmbeddedNoTag,
		})
		var ormTagValue string
		for _, structField := range structFields {
			ormTagValue = structField.Tag(OrmTagForStruct)
			ormTagValue = gstr.Split(gstr.Trim(ormTagValue), ",")[0]
			if ormTagValue != "" && gregex.IsMatchString(regularFieldNameRegPattern, ormTagValue) {
				fields = append(fields, ormTagValue)
			} else {
				fields = append(fields, structField.Name())
			}
		}
	} else {
		fields = gutil.Keys(structOrMap)
	}
	return
}

// GetPrimaryKeyCondition returns a new where condition by primary field name.
// The optional parameter `where` is like follows:
// 123                             => primary=123
// []int{1, 2, 3}                  => primary IN(1,2,3)
// "john"                          => primary='john'
// []string{"john", "smith"}       => primary IN('john','smith')
// g.Map{"id": g.Slice{1,2,3}}     => id IN(1,2,3)
// g.Map{"id": 1, "name": "john"}  => id=1 AND name='john'
// etc.
//
// Note that it returns the given `where` parameter directly if the `primary` is empty
// or length of `where` > 1.
func GetPrimaryKeyCondition(primary string, where ...interface{}) (newWhereCondition []interface{}) {
	if len(where) == 0 {
		return nil
	}
	if primary == "" {
		return where
	}
	if len(where) == 1 {
		var (
			rv   = reflect.ValueOf(where[0])
			kind = rv.Kind()
		)
		if kind == reflect.Ptr {
			rv = rv.Elem()
			kind = rv.Kind()
		}
		switch kind {
		case reflect.Map, reflect.Struct:
			// Ignore the parameter `primary`.
			break

		default:
			return []interface{}{map[string]interface{}{
				primary: where[0],
			}}
		}
	}
	return where
}

type formatWhereHolderInput struct {
	WhereHolder
	OmitNil   bool
	OmitEmpty bool
	Schema    string
	Table     string // Table is used for fields mapping and filtering internally.
}

func isKeyValueCanBeOmitEmpty(omitEmpty bool, whereType string, key, value interface{}) bool {
	if !omitEmpty {
		return false
	}
	// Eg:
	// Where("id", []int{}).All()             -> SELECT xxx FROM xxx WHERE 0=1
	// Where("name", "").All()                -> SELECT xxx FROM xxx WHERE `name`=''
	// OmitEmpty().Where("id", []int{}).All() -> SELECT xxx FROM xxx
	// OmitEmpty().Where("name", "").All()    -> SELECT xxx FROM xxx
	// OmitEmpty().Where("1").All()           -> SELECT xxx FROM xxx WHERE 1
	switch whereType {
	case whereHolderTypeNoArgs:
		return false

	case whereHolderTypeIn:
		return gutil.IsEmpty(value)

	default:
		if gstr.Count(gconv.String(key), "?") == 0 && gutil.IsEmpty(value) {
			return true
		}
	}
	return false
}

// formatWhereHolder formats where statement and its arguments for `Where` and `Having` statements.
func formatWhereHolder(ctx context.Context, db DB, in formatWhereHolderInput) (newWhere string, newArgs []interface{}) {
	var (
		buffer      = bytes.NewBuffer(nil)
		reflectInfo = reflection.OriginValueAndKind(in.Where)
	)
	switch reflectInfo.OriginKind {
	case reflect.Array, reflect.Slice:
		newArgs = formatWhereInterfaces(db, gconv.Interfaces(in.Where), buffer, newArgs)

	case reflect.Map:
		for key, value := range MapOrStructToMapDeep(in.Where, true) {
			if in.OmitNil && empty.IsNil(value) {
				continue
			}
			if in.OmitEmpty && empty.IsEmpty(value) {
				continue
			}
			newArgs = formatWhereKeyValue(formatWhereKeyValueInput{
				Db:     db,
				Buffer: buffer,
				Args:   newArgs,
				Key:    key,
				Value:  value,
				Prefix: in.Prefix,
				Type:   in.Type,
			})
		}

	case reflect.Struct:
		// If the `where` parameter is `DO` struct, it then adds `OmitNil` option for this condition,
		// which will filter all nil parameters in `where`.
		if isDoStruct(in.Where) {
			in.OmitNil = true
		}
		// If `where` struct implements `iIterator` interface,
		// it then uses its Iterate function to iterate its key-value pairs.
		// For example, ListMap and TreeMap are ordered map,
		// which implement `iIterator` interface and are index-friendly for where conditions.
		if iterator, ok := in.Where.(iIterator); ok {
			iterator.Iterator(func(key, value interface{}) bool {
				ketStr := gconv.String(key)
				if in.OmitNil && empty.IsNil(value) {
					return true
				}
				if in.OmitEmpty && empty.IsEmpty(value) {
					return true
				}
				newArgs = formatWhereKeyValue(formatWhereKeyValueInput{
					Db:        db,
					Buffer:    buffer,
					Args:      newArgs,
					Key:       ketStr,
					Value:     value,
					OmitEmpty: in.OmitEmpty,
					Prefix:    in.Prefix,
					Type:      in.Type,
				})
				return true
			})
			break
		}
		// Automatically mapping and filtering the struct attribute.
		var (
			reflectType = reflectInfo.OriginValue.Type()
			structField reflect.StructField
			data        = MapOrStructToMapDeep(in.Where, true)
		)
		// If `Prefix` is given, it checks and retrieves the table name.
		if in.Prefix != "" {
			hasTable, _ := db.GetCore().HasTable(in.Prefix)
			if hasTable {
				in.Table = in.Prefix
			} else {
				ormTagTableName := getTableNameFromOrmTag(in.Where)
				if ormTagTableName != "" {
					in.Table = ormTagTableName
				}
			}
		}
		// Mapping and filtering fields if `Table` is given.
		if in.Table != "" {
			data, _ = db.GetCore().mappingAndFilterData(ctx, in.Schema, in.Table, data, true)
		}
		// Put the struct attributes in sequence in Where statement.
		var ormTagValue string
		for i := 0; i < reflectType.NumField(); i++ {
			structField = reflectType.Field(i)
			// Use tag value from `orm` as field name if specified.
			ormTagValue = structField.Tag.Get(OrmTagForStruct)
			ormTagValue = gstr.Split(gstr.Trim(ormTagValue), ",")[0]
			if ormTagValue == "" {
				ormTagValue = structField.Name
			}
			foundKey, foundValue := gutil.MapPossibleItemByKey(data, ormTagValue)
			if foundKey != "" {
				if in.OmitNil && empty.IsNil(foundValue) {
					continue
				}
				if in.OmitEmpty && empty.IsEmpty(foundValue) {
					continue
				}
				newArgs = formatWhereKeyValue(formatWhereKeyValueInput{
					Db:        db,
					Buffer:    buffer,
					Args:      newArgs,
					Key:       foundKey,
					Value:     foundValue,
					OmitEmpty: in.OmitEmpty,
					Prefix:    in.Prefix,
					Type:      in.Type,
				})
			}
		}

	default:
		// Where filter.
		var omitEmptyCheckValue interface{}
		if len(in.Args) == 1 {
			omitEmptyCheckValue = in.Args[0]
		} else {
			omitEmptyCheckValue = in.Args
		}
		if isKeyValueCanBeOmitEmpty(in.OmitEmpty, in.Type, in.Where, omitEmptyCheckValue) {
			return
		}
		// Usually a string.
		whereStr := gstr.Trim(gconv.String(in.Where))
		// Is `whereStr` a field name which composed as a key-value condition?
		// Eg:
		// Where("id", 1)
		// Where("id", g.Slice{1,2,3})
		if gregex.IsMatchString(regularFieldNameWithoutDotRegPattern, whereStr) && len(in.Args) == 1 {
			newArgs = formatWhereKeyValue(formatWhereKeyValueInput{
				Db:        db,
				Buffer:    buffer,
				Args:      newArgs,
				Key:       whereStr,
				Value:     in.Args[0],
				OmitEmpty: in.OmitEmpty,
				Prefix:    in.Prefix,
				Type:      in.Type,
			})
			in.Args = in.Args[:0]
			break
		}
		// If the first part is column name, it automatically adds prefix to the column.
		if in.Prefix != "" {
			array := gstr.Split(whereStr, " ")
			if ok, _ := db.GetCore().HasField(ctx, in.Table, array[0]); ok {
				whereStr = in.Prefix + "." + whereStr
			}
		}
		// Regular string and parameter place holder handling.
		// Eg:
		// Where("id in(?) and name=?", g.Slice{1,2,3}, "john")
		i := 0
		for {
			if i >= len(in.Args) {
				break
			}
			// ===============================================================
			// Sub query, which is always used along with a string condition.
			// ===============================================================
			if subModel, ok := in.Args[i].(*Model); ok {
				index := -1
				whereStr, _ = gregex.ReplaceStringFunc(`(\?)`, whereStr, func(s string) string {
					index++
					if i+len(newArgs) == index {
						sqlWithHolder, holderArgs := subModel.getHolderAndArgsAsSubModel(ctx)
						in.Args = gutil.SliceInsertAfter(in.Args, i, holderArgs...)
						// Automatically adding the brackets.
						return "(" + sqlWithHolder + ")"
					}
					return s
				})
				in.Args = gutil.SliceDelete(in.Args, i)
				continue
			}
			i++
		}
		buffer.WriteString(whereStr)
	}

	if buffer.Len() == 0 {
		return "", in.Args
	}
	if len(in.Args) > 0 {
		newArgs = append(newArgs, in.Args...)
	}
	newWhere = buffer.String()
	if len(newArgs) > 0 {
		if gstr.Pos(newWhere, "?") == -1 {
			if gregex.IsMatchString(lastOperatorRegPattern, newWhere) {
				// Eg: Where/And/Or("uid>=", 1)
				newWhere += "?"
			} else if gregex.IsMatchString(regularFieldNameRegPattern, newWhere) {
				newWhere = db.GetCore().QuoteString(newWhere)
				if len(newArgs) > 0 {
					if utils.IsArray(newArgs[0]) {
						// Eg:
						// Where("id", []int{1,2,3})
						// Where("user.id", []int{1,2,3})
						newWhere += " IN (?)"
					} else if empty.IsNil(newArgs[0]) {
						// Eg:
						// Where("id", nil)
						// Where("user.id", nil)
						newWhere += " IS NULL"
						newArgs = nil
					} else {
						// Eg:
						// Where/And/Or("uid", 1)
						// Where/And/Or("user.uid", 1)
						newWhere += "=?"
					}
				}
			}
		}
	}
	return handleSliceAndStructArgsForSql(newWhere, newArgs)
}

// formatWhereInterfaces formats `where` as []interface{}.
func formatWhereInterfaces(db DB, where []interface{}, buffer *bytes.Buffer, newArgs []interface{}) []interface{} {
	if len(where) == 0 {
		return newArgs
	}
	if len(where)%2 != 0 {
		buffer.WriteString(gstr.Join(gconv.Strings(where), ""))
		return newArgs
	}
	var str string
	for i := 0; i < len(where); i += 2 {
		str = gconv.String(where[i])
		if buffer.Len() > 0 {
			buffer.WriteString(" AND " + db.GetCore().QuoteWord(str) + "=?")
		} else {
			buffer.WriteString(db.GetCore().QuoteWord(str) + "=?")
		}
		if s, ok := where[i+1].(Raw); ok {
			buffer.WriteString(gconv.String(s))
		} else {
			newArgs = append(newArgs, where[i+1])
		}
	}
	return newArgs
}

type formatWhereKeyValueInput struct {
	Db        DB            // Db is the underlying DB object for current operation.
	Buffer    *bytes.Buffer // Buffer is the sql statement string without Args for current operation.
	Args      []interface{} // Args is the full arguments of current operation.
	Key       string        // The field name, eg: "id", "name", etc.
	Value     interface{}   // The field value, can be any types.
	Type      string        // The value in Where type.
	OmitEmpty bool          // Ignores current condition key if `value` is empty.
	Prefix    string        // Field prefix, eg: "user", "order", etc.
}

// formatWhereKeyValue handles each key-value pair of the parameter map.
func formatWhereKeyValue(in formatWhereKeyValueInput) (newArgs []interface{}) {
	var (
		quotedKey   = in.Db.GetCore().QuoteWord(in.Key)
		holderCount = gstr.Count(quotedKey, "?")
	)
	if isKeyValueCanBeOmitEmpty(in.OmitEmpty, in.Type, quotedKey, in.Value) {
		return in.Args
	}
	if in.Prefix != "" && !gstr.Contains(quotedKey, ".") {
		quotedKey = in.Prefix + "." + quotedKey
	}
	if in.Buffer.Len() > 0 {
		in.Buffer.WriteString(" AND ")
	}
	// If the value is type of slice, and there's only one '?' holder in
	// the key string, it automatically adds '?' holder chars according to its arguments count
	// and converts it to "IN" statement.
	var (
		reflectValue = reflect.ValueOf(in.Value)
		reflectKind  = reflectValue.Kind()
	)
	switch reflectKind {
	// Slice argument.
	case reflect.Slice, reflect.Array:
		if holderCount == 0 {
			in.Buffer.WriteString(quotedKey + " IN(?)")
			in.Args = append(in.Args, in.Value)
		} else {
			if holderCount != reflectValue.Len() {
				in.Buffer.WriteString(quotedKey)
				in.Args = append(in.Args, in.Value)
			} else {
				in.Buffer.WriteString(quotedKey)
				in.Args = append(in.Args, gconv.Interfaces(in.Value)...)
			}
		}

	default:
		if in.Value == nil || empty.IsNil(reflectValue) {
			if gregex.IsMatchString(regularFieldNameRegPattern, in.Key) {
				// The key is a single field name.
				in.Buffer.WriteString(quotedKey + " IS NULL")
			} else {
				// The key may have operation chars.
				in.Buffer.WriteString(quotedKey)
			}
		} else {
			// It also supports "LIKE" statement, which we consider it an operator.
			quotedKey = gstr.Trim(quotedKey)
			if gstr.Pos(quotedKey, "?") == -1 {
				like := " LIKE"
				if len(quotedKey) > len(like) && gstr.Equal(quotedKey[len(quotedKey)-len(like):], like) {
					// Eg: Where(g.Map{"name like": "john%"})
					in.Buffer.WriteString(quotedKey + " ?")
				} else if gregex.IsMatchString(lastOperatorRegPattern, quotedKey) {
					// Eg: Where(g.Map{"age > ": 16})
					in.Buffer.WriteString(quotedKey + " ?")
				} else if gregex.IsMatchString(regularFieldNameRegPattern, in.Key) {
					// The key is a regular field name.
					in.Buffer.WriteString(quotedKey + "=?")
				} else {
					// The key is not a regular field name.
					// Eg: Where(g.Map{"age > 16": nil})
					// Issue: https://github.com/gogf/gf/issues/765
					if empty.IsEmpty(in.Value) {
						in.Buffer.WriteString(quotedKey)
						break
					} else {
						in.Buffer.WriteString(quotedKey + "=?")
					}
				}
			} else {
				in.Buffer.WriteString(quotedKey)
			}
			if s, ok := in.Value.(Raw); ok {
				in.Buffer.WriteString(gconv.String(s))
			} else {
				in.Args = append(in.Args, in.Value)
			}
		}
	}
	return in.Args
}

// handleSliceAndStructArgsForSql is an important function, which handles the sql and all its arguments
// before committing them to underlying driver.
func handleSliceAndStructArgsForSql(
	oldSql string, oldArgs []interface{},
) (newSql string, newArgs []interface{}) {
	newSql = oldSql
	if len(oldArgs) == 0 {
		return
	}
	// insertHolderCount is used to calculate the inserting position for the '?' holder.
	insertHolderCount := 0
	// Handles the slice and struct type argument item.
	for index, oldArg := range oldArgs {
		argReflectInfo := reflection.OriginValueAndKind(oldArg)
		switch argReflectInfo.OriginKind {
		case reflect.Slice, reflect.Array:
			// It does not split the type of []byte.
			// Eg: table.Where("name = ?", []byte("john"))
			if _, ok := oldArg.([]byte); ok {
				newArgs = append(newArgs, oldArg)
				continue
			}
			var (
				valueHolderCount = gstr.Count(newSql, "?")
				argSliceLength   = argReflectInfo.OriginValue.Len()
			)
			if argSliceLength == 0 {
				// Empty slice argument, it converts the sql to a false sql.
				// Example:
				// Query("select * from xxx where id in(?)", g.Slice{}) -> select * from xxx where 0=1
				// Where("id in(?)", g.Slice{}) -> WHERE 0=1
				if gstr.Contains(newSql, "?") {
					whereKeyWord := " WHERE "
					if p := gstr.PosI(newSql, whereKeyWord); p == -1 {
						return "0=1", []interface{}{}
					} else {
						return gstr.SubStr(newSql, 0, p+len(whereKeyWord)) + "0=1", []interface{}{}
					}
				}
			} else {
				// Example:
				// Query("SELECT ?+?", g.Slice{1,2})
				// WHERE("id=?", g.Slice{1,2})
				for i := 0; i < argSliceLength; i++ {
					newArgs = append(newArgs, argReflectInfo.OriginValue.Index(i).Interface())
				}
			}

			// If the '?' holder count equals the length of the slice,
			// it does not implement the arguments splitting logic.
			// Eg: db.Query("SELECT ?+?", g.Slice{1, 2})
			if len(oldArgs) == 1 && valueHolderCount == argSliceLength {
				break
			}

			// counter is used to finding the inserting position for the '?' holder.
			var (
				counter  = 0
				replaced = false
			)
			newSql, _ = gregex.ReplaceStringFunc(`\?`, newSql, func(s string) string {
				if replaced {
					return s
				}
				counter++
				if counter == index+insertHolderCount+1 {
					replaced = true
					insertHolderCount += argSliceLength - 1
					return "?" + strings.Repeat(",?", argSliceLength-1)
				}
				return s
			})

		// Special struct handling.
		case reflect.Struct:
			switch oldArg.(type) {
			// The underlying driver supports time.Time/*time.Time types.
			case time.Time, *time.Time:
				newArgs = append(newArgs, oldArg)
				continue

			case gtime.Time:
				newArgs = append(newArgs, oldArg.(gtime.Time).Time)
				continue

			case *gtime.Time:
				newArgs = append(newArgs, oldArg.(*gtime.Time).Time)
				continue

			default:
				// It converts the struct to string in default
				// if it has implemented the String interface.
				if v, ok := oldArg.(iString); ok {
					newArgs = append(newArgs, v.String())
					continue
				}
			}
			newArgs = append(newArgs, oldArg)

		default:
			newArgs = append(newArgs, oldArg)
		}
	}
	return
}

// FormatSqlWithArgs binds the arguments to the sql string and returns a complete
// sql string, just for debugging.
func FormatSqlWithArgs(sql string, args []interface{}) string {
	index := -1
	newQuery, _ := gregex.ReplaceStringFunc(
		`(\?|:v\d+|\$\d+|@p\d+)`,
		sql,
		func(s string) string {
			index++
			if len(args) > index {
				if args[index] == nil {
					return "null"
				}
				// Parameters of type Raw do not require special treatment
				if v, ok := args[index].(Raw); ok {
					return gconv.String(v)
				}
				reflectInfo := reflection.OriginValueAndKind(args[index])
				if reflectInfo.OriginKind == reflect.Ptr &&
					(reflectInfo.OriginValue.IsNil() || !reflectInfo.OriginValue.IsValid()) {
					return "null"
				}
				switch reflectInfo.OriginKind {
				case reflect.String, reflect.Map, reflect.Slice, reflect.Array:
					return `'` + gstr.QuoteMeta(gconv.String(args[index]), `'`) + `'`

				case reflect.Struct:
					if t, ok := args[index].(time.Time); ok {
						return `'` + t.Format(`2006-01-02 15:04:05`) + `'`
					}
					return `'` + gstr.QuoteMeta(gconv.String(args[index]), `'`) + `'`
				}
				return gconv.String(args[index])
			}
			return s
		})
	return newQuery
}

// FormatMultiLineSqlToSingle formats sql template string into one line.
func FormatMultiLineSqlToSingle(sql string) (string, error) {
	var err error
	// format sql template string.
	sql, err = gregex.ReplaceString(`[\n\r\s]+`, " ", gstr.Trim(sql))
	if err != nil {
		return "", err
	}
	sql, err = gregex.ReplaceString(`\s{2,}`, " ", gstr.Trim(sql))
	if err != nil {
		return "", err
	}
	return sql, nil
}

func genTableFieldsCacheKey(group, schema, table string) string {
	return fmt.Sprintf(
		`%s%s@%s#%s`,
		cachePrefixTableFields,
		group,
		schema,
		table,
	)
}

func genSelectCacheKey(table, group, schema, name, sql string, args ...interface{}) string {
	if name == "" {
		name = fmt.Sprintf(
			`%s@%s#%s:%d`,
			table,
			group,
			schema,
			ghash.BKDR64([]byte(sql+", @PARAMS:"+gconv.String(args))),
		)
	}
	return fmt.Sprintf(`%s%s`, cachePrefixSelectCache, name)
}
