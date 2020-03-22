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

const (
	ORM_TAG_FOR_STRUCT  = "orm"
	ORM_TAG_FOR_UNIQUE  = "unique"
	ORM_TAG_FOR_PRIMARY = "primary"
)

var (
	// quoteWordReg is the regular expression object for a word check.
	quoteWordReg = regexp.MustCompile(`^[a-zA-Z0-9\-_]+$`)
)

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
func DataToMapDeep(obj interface{}) map[string]interface{} {
	data := gconv.Map(obj, ORM_TAG_FOR_STRUCT)
	for key, value := range data {
		rv := reflect.ValueOf(value)
		kind := rv.Kind()
		if kind == reflect.Ptr {
			rv = rv.Elem()
			kind = rv.Kind()
		}
		switch kind {
		case reflect.Struct:
			// The underlying driver supports time.Time/*time.Time types.
			if _, ok := value.(time.Time); ok {
				continue
			}
			if _, ok := value.(*time.Time); ok {
				continue
			}
			// Use string conversion in default.
			if s, ok := value.(apiString); ok {
				data[key] = s.String()
				continue
			}
			delete(data, key)
			for k, v := range DataToMapDeep(value) {
				data[k] = v
			}
		}
	}
	return data
}

// QuotePrefixTableName adds prefix string and quote chars for the table. It handles table string like:
// "user", "user u", "user,user_detail", "user u, user_detail ut", "user as u, user_detail as ut", "user.user u".
//
// Note that, this will automatically checks the table prefix whether already added, if true it does
// nothing to the table name, or else adds the prefix to the table name.
func doHandleTableName(table, prefix, charLeft, charRight string) string {
	index := 0
	array1 := gstr.SplitAndTrim(table, ",")
	for k1, v1 := range array1 {
		array2 := gstr.SplitAndTrim(v1, " ")
		// Trim the security chars.
		array2[0] = gstr.TrimLeftStr(array2[0], charLeft)
		array2[0] = gstr.TrimRightStr(array2[0], charRight)
		// Check whether it has database name.
		array3 := gstr.Split(gstr.Trim(array2[0]), ".")
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
// "user", "user u", "user,user_detail", "user u, user_detail ut",
// "user.user u, user.user_detail ut", "u.id asc".
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
// 123, []int{1, 2, 3}, "john", []string{"john", "smith"}
// g.Map{"id": g.Slice{1,2,3}}, g.Map{"id": 1, "name": "john"}, etc.
//
// Note that it returns the given <where> parameter directly if there's the <primary> is empty.
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
func formatSql(sql string, args []interface{}) (newQuery string, newArgs []interface{}) {
	return handleArguments(sql, args)
}

// formatWhere formats where statement and its arguments.
// TODO []interface{} type support for parameter <where> does not completed yet.
func formatWhere(db DB, where interface{}, args []interface{}, omitEmpty bool) (newWhere string, newArgs []interface{}) {
	buffer := bytes.NewBuffer(nil)
	rv := reflect.ValueOf(where)
	kind := rv.Kind()
	if kind == reflect.Ptr {
		rv = rv.Elem()
		kind = rv.Kind()
	}
	switch kind {
	case reflect.Array, reflect.Slice:
		newArgs = formatWhereInterfaces(db, gconv.Interfaces(where), buffer, newArgs)

	case reflect.Map:
		for key, value := range DataToMapDeep(where) {
			if omitEmpty && empty.IsEmpty(value) {
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
				if omitEmpty && empty.IsEmpty(value) {
					return true
				}
				newArgs = formatWhereKeyValue(db, buffer, newArgs, gconv.String(key), value)
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
			if lastOperatorReg.MatchString(newWhere) {
				// Eg: Where/And/Or("uid>=", 1)
				newWhere += "?"
			} else if gregex.IsMatchString(`^[\w\.\-]+$`, newWhere) {
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
// TODO []interface{} type support for parameter <where> does not completed yet.
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
	rv := reflect.ValueOf(value)
	kind := rv.Kind()
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
			if gregex.IsMatchString(`^[\w\.\-]+$`, key) {
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
					buffer.WriteString(quotedKey + " ?")
				} else if lastOperatorReg.MatchString(quotedKey) {
					buffer.WriteString(quotedKey + " ?")
				} else {
					buffer.WriteString(quotedKey + "=?")
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
	// Handles the slice arguments.
	if len(args) > 0 {
		for index, arg := range args {
			rv := reflect.ValueOf(arg)
			kind := rv.Kind()
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
				for i := 0; i < rv.Len(); i++ {
					newArgs = append(newArgs, rv.Index(i).Interface())
				}
				// It the '?' holder count equals the length of the slice,
				// it does not implement the arguments splitting logic.
				// Eg: db.Query("SELECT ?+?", g.Slice{1, 2})
				if len(args) == 1 && gstr.Count(newSql, "?") == rv.Len() {
					break
				}
				// counter is used to finding the inserting position for the '?' holder.
				counter := 0
				newSql, _ = gregex.ReplaceStringFunc(`\?`, newSql, func(s string) string {
					counter++
					if counter == index+1 {
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
				// It converts the struct to string in default
				// if it implements the String interface.
				if v, ok := arg.(apiString); ok {
					newArgs = append(newArgs, v.String())
					continue
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
		`(\?|:\d+|\$\d+|@p\d+)`, sql, func(s string) string {
			index++
			if len(args) > index {
				if args[index] == nil {
					return "null"
				}
				rv := reflect.ValueOf(args[index])
				kind := rv.Kind()
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
