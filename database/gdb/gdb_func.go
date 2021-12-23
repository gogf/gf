// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"bytes"
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/empty"
	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/os/gstructs"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gmeta"
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

// iMapStrAny is the interface support for converting struct parameter to map.
type iMapStrAny interface {
	MapStrAny() map[string]interface{}
}

// iTableName is the interface for retrieving table name fro struct.
type iTableName interface {
	TableName() string
}

const (
	OrmTagForStruct    = "orm"
	OrmTagForTable     = "table"
	OrmTagForWith      = "with"
	OrmTagForWithWhere = "where"
	OrmTagForWithOrder = "order"
	OrmTagForDto       = "dto"
)

var (
	// quoteWordReg is the regular expression object for a word check.
	quoteWordReg = regexp.MustCompile(`^[a-zA-Z0-9\-_]+$`)

	// Priority tags for struct converting for orm field mapping.
	structTagPriority = append([]string{OrmTagForStruct}, gconv.StructTagPriority...)
)

// isDtoStruct checks and returns whether given type is a DTO struct.
func isDtoStruct(object interface{}) bool {
	// It checks by struct name like "XxxForDao", to be compatible with old version.
	// TODO remove this compatible codes in future.
	reflectType := reflect.TypeOf(object)
	if gstr.HasSuffix(reflectType.String(), modelForDaoSuffix) {
		return true
	}
	// It checks by struct meta for DTO struct in version.
	if ormTag := gmeta.Get(object, OrmTagForStruct); !ormTag.IsEmpty() {
		match, _ := gregex.MatchString(
			fmt.Sprintf(`%s\s*:\s*([^,]+)`, OrmTagForDto),
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
func GetInsertOperationByOption(option int) string {
	var operator string
	switch option {
	case insertOptionReplace:
		operator = "REPLACE"
	case insertOptionIgnore:
		operator = "INSERT IGNORE"
	default:
		operator = "INSERT"
	}
	return operator
}

// DataToMapDeep converts `value` to map type recursively(if attribute struct is embedded).
// The parameter `value` should be type of *map/map/*struct/struct.
// It supports embedded struct definition for struct.
func DataToMapDeep(value interface{}) map[string]interface{} {
	m := gconv.Map(value, structTagPriority...)
	for k, v := range m {
		switch v.(type) {
		case time.Time, *time.Time, gtime.Time, *gtime.Time:
			m[k] = v

		default:
			// Use string conversion in default.
			if s, ok := v.(iString); ok {
				m[k] = s.String()
			} else {
				m[k] = v
			}
		}
	}
	return m
}

// doHandleTableName adds prefix string and quote chars for the table. It handles table string like:
// "user", "user u", "user,user_detail", "user u, user_detail ut", "user as u, user_detail as ut",
// "user.user u", "`user`.`user` u".
//
// Note that, this will automatically checks the table prefix whether already added, if true it does
// nothing to the table name, or else adds the prefix to the table name.
func doHandleTableName(table, prefix, charLeft, charRight string) string {
	var (
		index  = 0
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
		for _, structField := range structFields {
			if tag := structField.Tag(OrmTagForStruct); tag != "" && gregex.IsMatchString(regularFieldNameRegPattern, tag) {
				fields = append(fields, tag)
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

// formatSql formats the sql string and its arguments before executing.
// The internal handleArguments function might be called twice during the SQL procedure,
// but do not worry about it, it's safe and efficient.
func formatSql(sql string, args []interface{}) (newSql string, newArgs []interface{}) {
	// DO NOT do this as there may be multiple lines and comments in the sql.
	// sql = gstr.Trim(sql)
	// sql = gstr.Replace(sql, "\n", " ")
	// sql, _ = gregex.ReplaceString(`\s{2,}`, ` `, sql)
	return handleArguments(sql, args)
}

type formatWhereHolderInput struct {
	Where     interface{}
	Args      []interface{}
	OmitNil   bool
	OmitEmpty bool
	Schema    string
	Table     string // Table is used for fields mapping and filtering internally.
	Prefix    string // Field prefix, eg: "user.", "order.".
}

// formatWhereHolder formats where statement and its arguments for `Where` and `Having` statements.
func formatWhereHolder(db DB, in formatWhereHolderInput) (newWhere string, newArgs []interface{}) {
	var (
		buffer      = bytes.NewBuffer(nil)
		reflectInfo = utils.OriginValueAndKind(in.Where)
	)
	switch reflectInfo.OriginKind {
	case reflect.Array, reflect.Slice:
		newArgs = formatWhereInterfaces(db, gconv.Interfaces(in.Where), buffer, newArgs)

	case reflect.Map:
		for key, value := range DataToMapDeep(in.Where) {
			if gregex.IsMatchString(regularFieldNameRegPattern, key) {
				if in.OmitNil && empty.IsNil(value) {
					continue
				}
				if in.OmitEmpty && empty.IsEmpty(value) {
					continue
				}
			}
			newArgs = formatWhereKeyValue(formatWhereKeyValueInput{
				Db:     db,
				Buffer: buffer,
				Args:   newArgs,
				Key:    key,
				Value:  value,
				Prefix: in.Prefix,
			})
		}

	case reflect.Struct:
		// If the `where` parameter is DTO struct, it then adds `OmitNil` option for this condition,
		// which will filter all nil parameters in `where`.
		if isDtoStruct(in.Where) {
			in.OmitNil = true
		}
		// If `where` struct implements `iIterator` interface,
		// it then uses its Iterate function to iterate its key-value pairs.
		// For example, ListMap and TreeMap are ordered map,
		// which implement `iIterator` interface and are index-friendly for where conditions.
		if iterator, ok := in.Where.(iIterator); ok {
			iterator.Iterator(func(key, value interface{}) bool {
				ketStr := gconv.String(key)
				if gregex.IsMatchString(regularFieldNameRegPattern, ketStr) {
					if in.OmitNil && empty.IsNil(value) {
						return true
					}
					if in.OmitEmpty && empty.IsEmpty(value) {
						return true
					}
				}
				newArgs = formatWhereKeyValue(formatWhereKeyValueInput{
					Db:        db,
					Buffer:    buffer,
					Args:      newArgs,
					Key:       ketStr,
					Value:     value,
					OmitEmpty: in.OmitEmpty,
					Prefix:    in.Prefix,
				})
				return true
			})
			break
		}
		// Automatically mapping and filtering the struct attribute.
		var (
			reflectType = reflectInfo.OriginValue.Type()
			structField reflect.StructField
			data        = DataToMapDeep(in.Where)
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
			data, _ = db.GetCore().mappingAndFilterData(in.Schema, in.Table, data, true)
		}
		// Put the struct attributes in sequence in Where statement.
		for i := 0; i < reflectType.NumField(); i++ {
			structField = reflectType.Field(i)
			foundKey, foundValue := gutil.MapPossibleItemByKey(data, structField.Name)
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
				})
			}
		}

	default:
		// Usually a string.
		whereStr := gconv.String(in.Where)
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
			})
			in.Args = in.Args[:0]
			break
		}
		// If the first part is column name, it automatically adds prefix to the column.
		if in.Prefix != "" {
			array := gstr.Split(whereStr, " ")
			if ok, _ := db.GetCore().HasField(in.Table, array[0]); ok {
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
			// Sub query, which is always used along with a string condition.
			if model, ok := in.Args[i].(*Model); ok {
				index := -1
				whereStr, _ = gregex.ReplaceStringFunc(`(\?)`, whereStr, func(s string) string {
					index++
					if i+len(newArgs) == index {
						sqlWithHolder, holderArgs := model.getFormattedSqlAndArgs(queryTypeNormal, false)
						newArgs = append(newArgs, holderArgs...)
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
	return handleArguments(newWhere, newArgs)
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
	OmitEmpty bool          // Ignores current condition key if `value` is empty.
	Prefix    string        // Field prefix, eg: "user", "order", etc.
}

// formatWhereKeyValue handles each key-value pair of the parameter map.
func formatWhereKeyValue(in formatWhereKeyValueInput) (newArgs []interface{}) {
	var (
		quotedKey   = in.Db.GetCore().QuoteWord(in.Key)
		holderCount = gstr.Count(quotedKey, "?")
	)
	// Eg:
	// Where("id", []int{}).All()             -> SELECT xxx FROM xxx WHERE 0=1
	// Where("name", "").All()                -> SELECT xxx FROM xxx WHERE `name`=''
	// OmitEmpty().Where("id", []int{}).All() -> SELECT xxx FROM xxx
	// OmitEmpty().("name", "").All()         -> SELECT xxx FROM xxx
	if in.OmitEmpty && holderCount == 0 && gutil.IsEmpty(in.Value) {
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

// handleArguments is an important function, which handles the sql and all its arguments
// before committing them to underlying driver.
func handleArguments(sql string, args []interface{}) (newSql string, newArgs []interface{}) {
	newSql = sql
	// insertHolderCount is used to calculate the inserting position for the '?' holder.
	insertHolderCount := 0
	// Handles the slice arguments.
	if len(args) > 0 {
		for index, arg := range args {
			reflectInfo := utils.OriginValueAndKind(arg)
			switch reflectInfo.OriginKind {
			case reflect.Slice, reflect.Array:
				// It does not split the type of []byte.
				// Eg: table.Where("name = ?", []byte("john"))
				if _, ok := arg.([]byte); ok {
					newArgs = append(newArgs, arg)
					continue
				}

				if reflectInfo.OriginValue.Len() == 0 {
					// Empty slice argument, it converts the sql to a false sql.
					// Eg:
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
					for i := 0; i < reflectInfo.OriginValue.Len(); i++ {
						newArgs = append(newArgs, reflectInfo.OriginValue.Index(i).Interface())
					}
				}

				// If the '?' holder count equals the length of the slice,
				// it does not implement the arguments splitting logic.
				// Eg: db.Query("SELECT ?+?", g.Slice{1, 2})
				if len(args) == 1 && gstr.Count(newSql, "?") == reflectInfo.OriginValue.Len() {
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
						insertHolderCount += reflectInfo.OriginValue.Len() - 1
						return "?" + strings.Repeat(",?", reflectInfo.OriginValue.Len()-1)
					}
					return s
				})

			// Special struct handling.
			case reflect.Struct:
				switch v := arg.(type) {
				// The underlying driver supports time.Time/*time.Time types.
				case time.Time, *time.Time:
					newArgs = append(newArgs, arg)
					continue

				// Special handling for gtime.Time/*gtime.Time.
				//
				// DO NOT use its underlying gtime.Time.Time as its argument,
				// because the std time.Time will be converted to certain timezone
				// according to underlying driver. And the underlying driver also
				// converts the time.Time to string automatically as the following does.
				case gtime.Time:
					newArgs = append(newArgs, v.String())
					continue

				case *gtime.Time:
					newArgs = append(newArgs, v.String())
					continue

				default:
					// It converts the struct to string in default
					// if it has implemented the String interface.
					if v, ok := arg.(iString); ok {
						newArgs = append(newArgs, v.String())
						continue
					}
				}
				newArgs = append(newArgs, arg)

			default:
				newArgs = append(newArgs, arg)
			}
		}
	}
	return
}

// formatError customizes and returns the SQL error.
func formatError(err error, s string, args ...interface{}) error {
	if err != nil && err != sql.ErrNoRows {
		return gerror.NewCodef(gcode.CodeDbOperationError, "%s, %s\n", err.Error(), FormatSqlWithArgs(s, args))
	}
	return err
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
				reflectInfo := utils.OriginValueAndKind(args[index])
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
