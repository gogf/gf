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
	"reflect"
	"strings"
	"time"

	"github.com/gogf/gf/internal/structs"

	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
)

// Type assert api for String.
type apiString interface {
	String() string
}

// Type assert api for Iterator.
type apiIterator interface {
	Iterator(f func(key, value interface{}) bool)
}

// Type assert api for Interfaces.
type apiInterfaces interface {
	Interfaces() []interface{}
}

const (
	ORM_TAG_FOR_STRUCT  = "orm"
	ORM_TAG_FOR_UNIQUE  = "unique"
	ORM_TAG_FOR_PRIMARY = "primary"
)

// 获得struct对象对应的where查询条件
func GetWhereConditionOfStruct(pointer interface{}) (where string, args []interface{}) {
	array := ([]string)(nil)
	for tag, field := range structs.TagMapField(pointer, []string{ORM_TAG_FOR_STRUCT}, true) {
		array = strings.Split(tag, ",")
		if len(array) > 1 && gstr.InArray([]string{ORM_TAG_FOR_UNIQUE, ORM_TAG_FOR_PRIMARY}, array[1]) {
			return array[0], []interface{}{field.Value()}
		}
		if len(where) > 0 {
			where += " "
		}
		where += tag + "=?"
		args = append(args, field.Value())
	}
	return
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
		// TODO garray support.
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
			} else if wordReg.MatchString(newWhere) {
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
				if rv.Len() == 0 {
					continue
				}
				// 不拆分[]byte类型
				if _, ok := arg.([]byte); ok {
					newArgs = append(newArgs, arg)
					continue
				}
				for i := 0; i < rv.Len(); i++ {
					newArgs = append(newArgs, rv.Index(i).Interface())
				}
				// 如果参数直接传递slice，并且占位符数量与slice长度相等，
				// 那么不用替换扩展占位符数量，直接使用该slice作为查询参数
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
