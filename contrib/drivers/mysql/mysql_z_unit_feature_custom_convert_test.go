// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mysql_test

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

type sqlScanMyDecimal struct {
	v float64
}

func (m *sqlScanMyDecimal) Scan(src any) (err error) {
	var v float64
	switch sv := src.(type) {
	case []byte:
		v, err = strconv.ParseFloat(string(sv), 64)
	case string:
		v, err = strconv.ParseFloat(sv, 64)
	case float64:
		v = sv
	case float32:
		v = float64(sv)
	default:
		err = fmt.Errorf("unknown type: %v(%T)", src, src)
	}
	if err != nil {
		return err
	}
	m.v = v
	return nil
}

func testCustomConvert(t *gtest.T, f func(t *gtest.T)) {
	sql := `CREATE TABLE IF NOT EXISTS %s (
	id int(10) unsigned NOT NULL AUTO_INCREMENT primary key,
	my_decimal1 decimal(5,2) NOT NULL,
	my_decimal2 decimal(5,2) NOT NULL                                 
	) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`
	tableName := "test_decimal"
	_, err := db.Exec(ctx, fmt.Sprintf(sql, tableName))
	t.AssertNil(err)
	defer dropTable(tableName)
	table := db.Model(tableName)
	data := g.Map{
		"my_decimal1": 777.333,
		"my_decimal2": 888.444,
	}
	_, err = table.Data(data).Insert()
	t.AssertNil(err)
	f(t)
}

func Test_Custom_Convert_SqlScanner(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		tableName := "test_decimal"
		table := db.Model(tableName)
		type TestDecimal struct {
			MyDecimal1 sqlScanMyDecimal  `orm:"my_decimal1"`
			MyDecimal2 *sqlScanMyDecimal `orm:"my_decimal2"`
		}
		testCustomConvert(t, func(t *gtest.T) {
			var res *TestDecimal
			err := table.Scan(&res)
			t.AssertNil(err)
			t.Assert(res.MyDecimal1.v, 777.33)
			t.Assert(res.MyDecimal2.v, 888.44)
		})
	})
}

func testCustomFieldConvertFunc(dest reflect.Value, src any) (err error) {
	var v float64
	switch sv := src.(type) {
	case []byte:
		v, err = strconv.ParseFloat(string(sv), 64)
	case string:
		v, err = strconv.ParseFloat(sv, 64)
	case float64:
		v = sv
	case float32:
		v = float64(sv)
	default:
		return fmt.Errorf("unknown type: %v(%T)", src, src)
	}
	if err != nil {
		return err
	}

	if dest.Kind() == reflect.Ptr {
		if dest.IsNil() {
			dest.Set(reflect.New(dest.Type().Elem()))
		}
		dest = dest.Elem()
	}
	switch dest.Kind() {
	case reflect.Float32, reflect.Float64:
		dest.SetFloat(v)
	case reflect.Int, reflect.Int64 /* reflect.Int16,reflect.Int32 */ :
		dest.SetInt(int64(v))
	case reflect.Uint, reflect.Uint64 /* reflect.Uint16,reflect.Uint32 */ :
		dest.SetUint(uint64(v))
	// case other types...:
	default:
		if dest.Kind() != reflect.Ptr {
			dest = dest.Addr()
			switch impl := dest.Interface().(type) {
			case sql.Scanner:
				return impl.Scan(v)
			default:
			}
		}
		return fmt.Errorf("Error: mysql_z_unit_feature_custom_convert_test.go:125  unsupported types: %v", dest.Type())
	}
	return nil
}

func Test_Custom_Convert_RegisterDatabaseConvertFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		gdb.RegisterDatabaseConvertFunc("mysql", "decimal", testCustomFieldConvertFunc)
		tableName := "test_decimal"
		table := db.Model(tableName)
		type MyDecimal float64
		type TestDecimal struct {
			MyDecimal1 MyDecimal  `orm:"my_decimal1"`
			MyDecimal2 *MyDecimal `orm:"my_decimal2"`
		}
		testCustomConvert(t, func(t *gtest.T) {
			var res *TestDecimal
			err := table.Scan(&res)
			t.AssertNil(err)
			t.Assert(res.MyDecimal1, 777.33)
			t.Assert(res.MyDecimal2, 888.44)
		})
	})
}

func Test_Custom_Convert_RegisterStructFieldConvertFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type MyDecimal float64
		type TestDecimal struct {
			MyDecimal1 MyDecimal         `orm:"my_decimal1"`
			MyDecimal2 *sqlScanMyDecimal `orm:"my_decimal2"`
		}
		gdb.RegisterStructFieldConvertFunc(reflect.TypeOf(TestDecimal{}), "MyDecimal1", testCustomFieldConvertFunc)
		tableName := "test_decimal"
		table := db.Model(tableName)
		testCustomConvert(t, func(t *gtest.T) {
			var res *TestDecimal
			err := table.Scan(&res)
			t.AssertNil(err)
			t.Assert(res.MyDecimal1, 777.33)
			t.Assert(res.MyDecimal2.v, 888.44)
		})
	})
}
