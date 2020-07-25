// Copyright 2017-2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gogf/gf/internal/empty"
	"github.com/gogf/gf/internal/utils"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gutil"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/gogf/gf/internal/structs"

	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
)

// apiString is the type assert api for String.
type apiString interface {
	String() string
}

// apiIterator is the type assert api for Iterator.
type apiIterator interface {
	Iterator(f func(key, value interface{}) bool)
}

// apiInterfaces is the type assert api for Interfaces.
type apiInterfaces interface {
	Interfaces() []interface{}
}

// apiMapStrAny is the interface support for converting struct parameter to map.
type apiMapStrAny interface {
	MapStrAny() map[string]interface{}
}

const (
	ORM_TAG_FOR_STRUCT  = "orm"
	ORM_TAG_FOR_UNIQUE  = "unique"
	ORM_TAG_FOR_PRIMARY = "primary"
)

var (
	// quoteWordReg is the regular expression object for a word check.
	quoteWordReg = regexp.MustCompile(`^[a-zA-Z0-9\-_]+$`)

	// Priority tags for struct converting for orm field mapping.
	structTagPriority = append([]string{ORM_TAG_FOR_STRUCT}, gconv.StructTagPriority...)
)

// ListItemValues retrieves and returns the elements of all item struct/map with key <key>.
// Note that the parameter <list> should be type of slice which contains elements of map or struct,
// or else it returns an empty slice.
//
// The parameter <list> supports types like:
// []map[string]interface{}
// []map[string]sub-map
// []struct
// []struct:sub-struct
// Note that the sub-map/sub-struct makes sense only if the optional parameter <subKey> is given.
// See gutil.ListItemValues.
func ListItemValues(list interface{}, key interface{}, subKey ...interface{}) (values []interface{}) {
	return gutil.ListItemValues(list, key, subKey...)
}

// ListItemValuesUnique retrieves and returns the unique elements of all struct/map with key <key>.
// Note that the parameter <list> should be type of slice which contains elements of map or struct,
// or else it returns an empty slice.
// See gutil.ListItemValuesUnique.
func ListItemValuesUnique(list interface{}, key string, subKey ...interface{}) []interface{} {
	return gutil.ListItemValuesUnique(list, key, subKey...)
}

// GetInsertOperationByOption returns proper insert option with given parameter <option>.
func GetInsertOperationByOption(option int) string {
	var operator string
	switch option {
	case gINSERT_OPTION_REPLACE:
		operator = "REPLACE"
	case gINSERT_OPTION_IGNORE:
		operator = "INSERT IGNORE"
	default:
		operator = "INSERT"
	}
	return operator
}

// DataToMapDeep converts struct object to map type recursively.
// The parameter <obj> should be type of *map/map/*struct/struct.
// It supports inherit struct definition for struct.
func DataToMapDeep(value interface{}) map[string]interface{} {
	if v, ok := value.(apiMapStrAny); ok {
		return v.MapStrAny()
	}

	var (
		rvValue reflect.Value
		rvField reflect.Value
		rvKind  reflect.Kind
		rtField reflect.StructField
	)
	if v, ok := value.(reflect.Value); ok {
		rvValue = v
	} else {
		rvValue = reflect.ValueOf(value)
	}
	rvKind = rvValue.Kind()
	if rvKind == reflect.Ptr {
		rvValue = rvValue.Elem()
		rvKind = rvValue.Kind()
	}
	// If given <value> is not a struct, it uses gconv.Map for converting.
	if rvKind != reflect.Struct {
		return gconv.Map(value, structTagPriority...)
	}
	// Struct handling.
	var (
		fieldTag reflect.StructTag
		rvType   = rvValue.Type()
		name     = ""
		data     = make(map[string]interface{})
	)
	for i := 0; i < rvValue.NumField(); i++ {
		rtField = rvType.Field(i)
		rvField = rvValue.Field(i)
		fieldName := rtField.Name
		if !utils.IsLetterUpper(fieldName[0]) {
			continue
		}
		// Struct attribute inherit
		if rtField.Anonymous {
			for k, v := range DataToMapDeep(rvField) {
				data[k] = v
			}
			continue
		}
		// Other attributes.
		name = ""
		fieldTag = rtField.Tag
		for _, tag := range structTagPriority {
			if s := fieldTag.Get(tag); s != "" {
				name = s
				break
			}
		}
		if name == "" {
			name = fieldName
		} else {
			// The "orm" tag supports json tag feature: -, omitempty
			// The "orm" tag would be like: "id,priority", so it should use splitting handling.
			name = gstr.Trim(name)
			if name == "-" {
				continue
			}
			array := gstr.SplitAndTrim(name, ",")
			if len(array) > 1 {
				switch array[1] {
				case "omitempty":
					if empty.IsEmpty(rvField.Interface()) {
						continue
					} else {
						name = array[0]
					}
				default:
					name = array[0]
				}
			}
		}

		// The underlying driver supports time.Time/*time.Time types.
		fieldValue := rvField.Interface()
		switch fieldValue.(type) {
		case time.Time, *time.Time:
			data[name] = fieldValue
		default:
			// Use string conversion in default.
			if s, ok := fieldValue.(apiString); ok {
				data[name] = s.String()
			} else {
				data[name] = fieldValue
			}
		}
	}
	return data
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

// doQuoteWord checks given string <s> a word, if true quotes it with <charLeft> and <charRight>
// and returns the quoted string; or else returns <s> without any change.
func doQuoteWord(s, charLeft, charRight string) string {
	if quoteWordReg.MatchString(s) && !gstr.ContainsAny(s, charLeft+charRight) {
		return charLeft + s + charRight
	}
	return s
}

// doQuoteString quotes string with quote chars. It handles strings like:
// "user",
// "user u",
// "user,user_detail",
// "user u, user_detail ut",
// "user.user u, user.user_detail ut",
// "u.id, u.name, u.age",
// "u.id asc".
func doQuoteString(s, charLeft, charRight string) string {
	array1 := gstr.SplitAndTrim(s, ",")
	for k1, v1 := range array1 {
		array2 := gstr.SplitAndTrim(v1, " ")
		array3 := gstr.Split(gstr.Trim(array2[0]), ".")
		if len(array3) == 1 {
			array3[0] = doQuoteWord(array3[0], charLeft, charRight)
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

// GetWhereConditionOfStruct returns the where condition sql and arguments by given struct pointer.
// This function automatically retrieves primary or unique field and its attribute value as condition.
func GetWhereConditionOfStruct(pointer interface{}) (where string, args []interface{}) {
	array := ([]string)(nil)
	for _, field := range structs.TagFields(pointer, []string{ORM_TAG_FOR_STRUCT}, true) {
		array = strings.Split(field.Tag, ",")
		if len(array) > 1 && gstr.InArray([]string{ORM_TAG_FOR_UNIQUE, ORM_TAG_FOR_PRIMARY}, array[1]) {
			return array[0], []interface{}{field.Value()}
		}
		if len(where) > 0 {
			where += " "
		}
		where += field.Tag + "=?"
		args = append(args, field.Value())
	}
	return
}

// GetPrimaryKey retrieves and returns primary key field name from given struct.
func GetPrimaryKey(pointer interface{}) string {
	array := ([]string)(nil)
	for _, field := range structs.TagFields(pointer, []string{ORM_TAG_FOR_STRUCT}, true) {
		array = strings.Split(field.Tag, ",")
		if len(array) > 1 && array[1] == ORM_TAG_FOR_PRIMARY {
			return array[0]
		}
	}
	return ""
}

// GetPrimaryKeyCondition returns a new where condition by primary field name.
// The optional parameter <where> is like follows:
// 123                             => primary=123
// []int{1, 2, 3}                  => primary IN(1,2,3)
// "john"                          => primary='john'
// []string{"john", "smith"}       => primary IN('john','smith')
// g.Map{"id": g.Slice{1,2,3}}     => id IN(1,2,3)
// g.Map{"id": 1, "name": "john"}  => id=1 AND name='john'
// etc.
//
// Note that it returns the given <where> parameter directly if the <primary> is empty
// or length of <where> > 1.
func GetPrimaryKeyCondition(primary string, where ...interface{}) (newWhereCondition []interface{}) {
	if len(where) == 0 {
		return nil
	}
	if primary == "" {
		return where
	}
	if len(where) == 1 {
		rv := reflect.ValueOf(where[0])
		kind := rv.Kind()
		if kind == reflect.Ptr {
			rv = rv.Elem()
			kind = rv.Kind()
		}
		switch kind {
		case reflect.Map, reflect.Struct:
			// Ignore the parameter <primary>.
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
	sql = gstr.Trim(sql)
	sql = gstr.Replace(sql, "\n", " ")
	sql, _ = gregex.ReplaceString(`\s{2,}`, ` `, sql)
	return handleArguments(sql, args)
}

// formatWhere formats where statement and its arguments.
// TODO []interface{} type support for parameter <where> does not completed yet.
func formatWhere(db DB, where interface{}, args []interface{}, omitEmpty bool) (newWhere string, newArgs []interface{}) {
	var (
		buffer = bytes.NewBuffer(nil)
		rv     = reflect.ValueOf(where)
		kind   = rv.Kind()
	)
	if kind == reflect.Ptr {
		rv = rv.Elem()
		kind = rv.Kind()
	}
	switch kind {
	case reflect.Array, reflect.Slice:
		newArgs = formatWhereInterfaces(db, gconv.Interfaces(where), buffer, newArgs)

	case reflect.Map:
		for key, value := range DataToMapDeep(where) {
			if gregex.IsMatchString(regularFieldNameRegPattern, key) && omitEmpty && empty.IsEmpty(value) {
				continue
			}
			newArgs = formatWhereKeyValue(db, buffer, newArgs, key, value)
		}

	case reflect.Struct:
		// If <where> struct implements apiIterator interface,
		// it then uses its Iterate function to iterates its key-value pairs.
		// For example, ListMap and TreeMap are ordered map,
		// which implement apiIterator interface and are index-friendly for where conditions.
		if iterator, ok := where.(apiIterator); ok {
			iterator.Iterator(func(key, value interface{}) bool {
				ketStr := gconv.String(key)
				if gregex.IsMatchString(regularFieldNameRegPattern, ketStr) && omitEmpty && empty.IsEmpty(value) {
					return true
				}
				newArgs = formatWhereKeyValue(db, buffer, newArgs, ketStr, value)
				return true
			})
			break
		}
		for key, value := range DataToMapDeep(where) {
			if omitEmpty && empty.IsEmpty(value) {
				continue
			}
			newArgs = formatWhereKeyValue(db, buffer, newArgs, key, value)
		}

	default:
		buffer.WriteString(gconv.String(where))
	}

	if buffer.Len() == 0 {
		return "", args
	}
	newArgs = append(newArgs, args...)
	newWhere = buffer.String()
	if len(newArgs) > 0 {
		if gstr.Pos(newWhere, "?") == -1 {
			if gregex.IsMatchString(lastOperatorRegPattern, newWhere) {
				// Eg: Where/And/Or("uid>=", 1)
				newWhere += "?"
			} else if gregex.IsMatchString(regularFieldNameRegPattern, newWhere) {
				newWhere = db.QuoteString(newWhere)
				if len(newArgs) > 0 {
					if utils.IsArray(newArgs[0]) {
						// Eg: Where("id", []int{1,2,3})
						newWhere += " IN (?)"
					} else if empty.IsNil(newArgs[0]) {
						// Eg: Where("id", nil)
						newWhere += " IS NULL"
						newArgs = nil
					} else {
						// Eg: Where/And/Or("uid", 1)
						newWhere += "=?"
					}
				}
			}
		}
	}
	return handleArguments(newWhere, newArgs)
}

// formatWhereInterfaces formats <where> as []interface{}.
// TODO supporting for parameter <where> with []interface{} type is not completed yet.
func formatWhereInterfaces(db DB, where []interface{}, buffer *bytes.Buffer, newArgs []interface{}) []interface{} {
	var str string
	var array []interface{}
	var holderCount int
	for i := 0; i < len(where); {
		if holderCount > 0 {
			array = gconv.Interfaces(where[i])
			newArgs = append(newArgs, array...)
			holderCount -= len(array)
		} else {
			str = gconv.String(where[i])
			holderCount = gstr.Count(str, "?")
			buffer.WriteString(str)
		}
	}
	return newArgs
}

// formatWhereKeyValue handles each key-value pair of the parameter map.
func formatWhereKeyValue(db DB, buffer *bytes.Buffer, newArgs []interface{}, key string, value interface{}) []interface{} {
	quotedKey := db.QuoteWord(key)
	if buffer.Len() > 0 {
		buffer.WriteString(" AND ")
	}
	// If the value is type of slice, and there's only one '?' holder in
	// the key string, it automatically adds '?' holder chars according to its arguments count
	// and converts it to "IN" statement.
	var (
		rv   = reflect.ValueOf(value)
		kind = rv.Kind()
	)
	switch kind {
	case reflect.Slice, reflect.Array:
		count := gstr.Count(quotedKey, "?")
		if count == 0 {
			buffer.WriteString(quotedKey + " IN(?)")
			newArgs = append(newArgs, value)
		} else if count != rv.Len() {
			buffer.WriteString(quotedKey)
			newArgs = append(newArgs, value)
		} else {
			buffer.WriteString(quotedKey)
			newArgs = append(newArgs, gconv.Interfaces(value)...)
		}
	default:
		if value == nil || empty.IsNil(rv) {
			if gregex.IsMatchString(regularFieldNameRegPattern, key) {
				// The key is a single field name.
				buffer.WriteString(quotedKey + " IS NULL")
			} else {
				// The key may have operation chars.
				buffer.WriteString(quotedKey)
			}
		} else {
			// It also supports "LIKE" statement, which we considers it an operator.
			quotedKey = gstr.Trim(quotedKey)
			if gstr.Pos(quotedKey, "?") == -1 {
				like := " like"
				if len(quotedKey) > len(like) && gstr.Equal(quotedKey[len(quotedKey)-len(like):], like) {
					// Eg: Where(g.Map{"name like": "john%"})
					buffer.WriteString(quotedKey + " ?")
				} else if gregex.IsMatchString(lastOperatorRegPattern, quotedKey) {
					// Eg: Where(g.Map{"age > ": 16})
					buffer.WriteString(quotedKey + " ?")
				} else if gregex.IsMatchString(regularFieldNameRegPattern, key) {
					// The key is a regular field name.
					buffer.WriteString(quotedKey + "=?")
				} else {
					// The key is not a regular field name.
					// Eg: Where(g.Map{"age > 16": nil})
					// Issue: https://github.com/gogf/gf/issues/765
					if empty.IsEmpty(value) {
						buffer.WriteString(quotedKey)
						break
					} else {
						buffer.WriteString(quotedKey + "=?")
					}
				}
			} else {
				buffer.WriteString(quotedKey)
			}
			newArgs = append(newArgs, value)
		}
	}
	return newArgs
}

// handleArguments is a nice function which handles the query and its arguments before committing to
// underlying driver.
func handleArguments(sql string, args []interface{}) (newSql string, newArgs []interface{}) {
	newSql = sql
	// insertHolderCount is used to calculate the inserting position for the '?' holder.
	insertHolderCount := 0
	// Handles the slice arguments.
	if len(args) > 0 {
		for index, arg := range args {
			var (
				rv   = reflect.ValueOf(arg)
				kind = rv.Kind()
			)
			if kind == reflect.Ptr {
				rv = rv.Elem()
				kind = rv.Kind()
			}
			switch kind {
			case reflect.Slice, reflect.Array:
				// It does not split the type of []byte.
				// Eg: table.Where("name = ?", []byte("john"))
				if _, ok := arg.([]byte); ok {
					newArgs = append(newArgs, arg)
					continue
				}

				if rv.Len() == 0 {
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
					for i := 0; i < rv.Len(); i++ {
						newArgs = append(newArgs, rv.Index(i).Interface())
					}
				}

				// If the '?' holder count equals the length of the slice,
				// it does not implement the arguments splitting logic.
				// Eg: db.Query("SELECT ?+?", g.Slice{1, 2})
				if len(args) == 1 && gstr.Count(newSql, "?") == rv.Len() {
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
						insertHolderCount += rv.Len() - 1
						return "?" + strings.Repeat(",?", rv.Len()-1)
					}
					return s
				})

			// Special struct handling.
			case reflect.Struct:
				// The underlying driver supports time.Time/*time.Time types.
				if _, ok := arg.(time.Time); ok {
					newArgs = append(newArgs, arg)
					continue
				}
				if _, ok := arg.(*time.Time); ok {
					newArgs = append(newArgs, arg)
					continue
				}
				switch v := arg.(type) {
				case time.Time, *time.Time:
					newArgs = append(newArgs, arg)
					continue

				// Special handling for gtime.Time.
				case gtime.Time:
					newArgs = append(newArgs, v.String())
					continue

				default:
					// It converts the struct to string in default
					// if it implements the String interface.
					if v, ok := arg.(apiString); ok {
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
func formatError(err error, sql string, args ...interface{}) error {
	if err != nil && err != ErrNoRows {
		return errors.New(fmt.Sprintf("%s, %s\n", err.Error(), FormatSqlWithArgs(sql, args)))
	}
	return err
}

// FormatSqlWithArgs binds the arguments to the sql string and returns a complete
// sql string, just for debugging.
func FormatSqlWithArgs(sql string, args []interface{}) string {
	index := -1
	newQuery, _ := gregex.ReplaceStringFunc(
		`(\?|:v\d+|\$\d+|@p\d+)`, sql, func(s string) string {
			index++
			if len(args) > index {
				if args[index] == nil {
					return "null"
				}
				var (
					rv   = reflect.ValueOf(args[index])
					kind = rv.Kind()
				)
				if kind == reflect.Ptr {
					if rv.IsNil() || !rv.IsValid() {
						return "null"
					}
					rv = rv.Elem()
					kind = rv.Kind()
				}
				switch kind {
				case reflect.String, reflect.Map, reflect.Slice, reflect.Array:
					return `'` + gstr.QuoteMeta(gconv.String(args[index]), `'`) + `'`
				case reflect.Struct:
					if t, ok := args[index].(time.Time); ok {
						return `'` + gtime.NewFromTime(t).String() + `'`
					}
					return `'` + gstr.QuoteMeta(gconv.String(args[index]), `'`) + `'`
				}
				return gconv.String(args[index])
			}
			return s
		})
	return newQuery
}

// mapToStruct maps the <data> to given struct.
// Note that the given parameter <pointer> should be a pointer to s struct.
func mapToStruct(data map[string]interface{}, pointer interface{}) error {
	// It retrieves and returns the mapping between orm tag and the struct attribute name.
	mapping := make(map[string]string)
	for tag, attr := range structs.TagMapName(pointer, []string{ORM_TAG_FOR_STRUCT}, true) {
		mapping[strings.Split(tag, ",")[0]] = attr
	}
	return gconv.StructDeep(data, pointer, mapping)
}
