// Copyright 2017-2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gogf/gf/internal/empty"
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

// apiInterfacesis the type assert api for Interfaces.
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

// handleTableName adds prefix string and quote chars for the table. It handles table string like:
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

// 获得orm标签与属性的映射关系
func GetOrmMappingOfStruct(pointer interface{}) map[string]string {
	mapping := make(map[string]string)
	for tag, attr := range structs.TagMapName(pointer, []string{ORM_TAG_FOR_STRUCT}, true) {
		mapping[strings.Split(tag, ",")[0]] = attr
	}
	return mapping
}

// 格式化SQL语句.
func formatQuery(query string, args []interface{}) (newQuery string, newArgs []interface{}) {
	return handlerSliceArguments(query, args)
}

// 格式化Where查询条件。
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
		for key, value := range varToMapDeep(where) {
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
		for key, value := range varToMapDeep(where) {
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
		// It supports formats like: Where/And/Or("uid", 1) , Where/And/Or("uid>=", 1)
		if gstr.Pos(newWhere, "?") == -1 {
			if lastOperatorReg.MatchString(newWhere) {
				newWhere += "?"
			} else if gregex.IsMatchString(`^[\w\.\-]+$`, newWhere) {
				newWhere += "=?"
			}
		}
	}
	return handlerSliceArguments(newWhere, newArgs)
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
	key = db.quoteWord(key)
	if buffer.Len() > 0 {
		buffer.WriteString(" AND ")
	}
	// 支持slice键值/属性，如果只有一个?占位符号，那么作为IN查询，否则打散作为多个查询参数
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		count := gstr.Count(key, "?")
		if count == 0 {
			buffer.WriteString(key + " IN(?)")
			newArgs = append(newArgs, value)
		} else if count != rv.Len() {
			buffer.WriteString(key)
			newArgs = append(newArgs, value)
		} else {
			buffer.WriteString(key)
			// 如果键名/属性名称中带有多个?占位符号，那么将参数打散
			newArgs = append(newArgs, gconv.Interfaces(value)...)
		}
	default:
		if value == nil {
			buffer.WriteString(key)
		} else {
			// 支持key带操作符号，注意like也算是操作符号
			key = gstr.Trim(key)
			if gstr.Pos(key, "?") == -1 {
				like := " like"
				if len(key) > len(like) && gstr.Equal(key[len(key)-len(like):], like) {
					buffer.WriteString(key + " ?")
				} else if lastOperatorReg.MatchString(key) {
					buffer.WriteString(key + " ?")
				} else {
					buffer.WriteString(key + "=?")
				}
			} else {
				buffer.WriteString(key)
			}
			newArgs = append(newArgs, value)
		}
	}
	return newArgs
}

// 将对象转换为map，如果对象带有继承对象，那么执行递归转换。
// 该方法用于将变量传递给数据库执行之前。
func varToMapDeep(obj interface{}) map[string]interface{} {
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
			// 底层数据库引擎支持 time.Time/*time.Time 类型
			if _, ok := value.(time.Time); ok {
				continue
			}
			if _, ok := value.(*time.Time); ok {
				continue
			}
			// 如果执行String方法，那么执行字符串转换
			if s, ok := value.(apiString); ok {
				data[key] = s.String()
				continue
			}
			delete(data, key)
			for k, v := range varToMapDeep(value) {
				data[k] = v
			}
		}
	}
	return data
}

// 处理预处理占位符与slice类型的参数。
// 需要注意的是，
// 如果是链式操作，在条件参数中也会调用该方法处理查询参数，
// 如果是方法参数，在sql提交执行之前也会再次调用该方法处理查询语句和参数。
func handlerSliceArguments(query string, args []interface{}) (newQuery string, newArgs []interface{}) {
	newQuery = query
	// 查询条件参数处理，主要处理slice参数类型
	if len(args) > 0 {
		for index, arg := range args {
			rv := reflect.ValueOf(arg)
			kind := rv.Kind()
			if kind == reflect.Ptr {
				rv = rv.Elem()
				kind = rv.Kind()
			}
			switch kind {
			// '?'占位符支持slice类型, 这里会将slice参数拆散，并更新原有占位符'?'为多个'?'，使用','符号连接。
			case reflect.Slice, reflect.Array:
				// 不拆分[]byte类型(当做字符串处理)
				// Eg: table.Where("name = ?", []byte("john"))
				if _, ok := arg.([]byte); ok {
					newArgs = append(newArgs, arg)
					continue
				}
				for i := 0; i < rv.Len(); i++ {
					newArgs = append(newArgs, rv.Index(i).Interface())
				}
				// 如果参数直接传递slice，并且占位符数量与slice长度相等，
				// 那么不用替换扩展占位符数量，直接使用该slice作为查询参数
				// Eg: db.Query("SELECT ?+?", g.Slice{1, 2})
				if len(args) == 1 && gstr.Count(newQuery, "?") == rv.Len() {
					break
				}
				// counter用于匹配该参数的位置(与index对应)
				counter := 0
				newQuery, _ = gregex.ReplaceStringFunc(`\?`, newQuery, func(s string) string {
					counter++
					if counter == index+1 {
						return "?" + strings.Repeat(",?", rv.Len()-1)
					}
					return s
				})

			// Special struct handling.
			case reflect.Struct:
				if v, ok := arg.(apiString); ok {
					newArgs = append(newArgs, v.String())
				} else {
					newArgs = append(newArgs, arg)
				}

			default:
				newArgs = append(newArgs, arg)
			}
		}
	}
	return
}

// 将预处理参数转换为底层数据库引擎支持的格式。
// 主要是判断参数是否为复杂数据类型，如果是，那么转换为基础类型。
func convertParam(value interface{}) interface{} {
	rv := reflect.ValueOf(value)
	kind := rv.Kind()
	if kind == reflect.Ptr {
		rv = rv.Elem()
		kind = rv.Kind()
	}
	switch kind {
	case reflect.Struct:
		// 底层数据库引擎支持 time.Time/*time.Time 类型
		if v, ok := value.(time.Time); ok {
			if v.IsZero() {
				return "null"
			}
			return value
		}
		if v, ok := value.(*time.Time); ok {
			if v.IsZero() {
				return ""
			}
			return value
		}
		return gconv.String(value)
	}
	return value
}

// 格式化错误信息
func formatError(err error, query string, args ...interface{}) error {
	if err != nil && err != sql.ErrNoRows {
		return errors.New(fmt.Sprintf("%s, %s\n", err.Error(), bindArgsToQuery(query, args)))
	}
	return err
}

// 根据insert选项获得操作名称
func getInsertOperationByOption(option int) string {
	operator := "INSERT"
	switch option {
	case gINSERT_OPTION_REPLACE:
		operator = "REPLACE"
	case gINSERT_OPTION_SAVE:
	case gINSERT_OPTION_IGNORE:
		operator = "INSERT IGNORE"
	}
	return operator
}

// 将参数绑定到SQL语句中，仅用于调试打印。
func bindArgsToQuery(query string, args []interface{}) string {
	index := -1
	newQuery, _ := gregex.ReplaceStringFunc(`\?`, query, func(s string) string {
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

// 使用递归的方式将map键值对映射到struct对象上，注意参数<pointer>是一个指向struct的指针。
func mapToStruct(data map[string]interface{}, pointer interface{}) error {
	return gconv.StructDeep(data, pointer, GetOrmMappingOfStruct(pointer))
}
