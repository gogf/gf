// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/test/gtest"
	"testing"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

func ExampleTransaction() {
	g.DB().Transaction(context.TODO(), func(ctx context.Context, tx gdb.TX) error {
		// user
		result, err := tx.Insert("user", g.Map{
			"passport": "john",
			"password": "12345678",
			"nickname": "JohnGuo",
		})
		if err != nil {
			return err
		}
		// user_detail
		id, err := result.LastInsertId()
		if err != nil {
			return err
		}
		_, err = tx.Insert("user_detail", g.Map{
			"uid":       id,
			"site":      "https://johng.cn",
			"true_name": "GuoQiang",
		})
		if err != nil {
			return err
		}
		return nil
	})
}

func Test_PartitionTable(t *testing.T) {
	setConfig()
	dropShopDBTable()
	createShopDBTable()
	insertShopDBData()
	//defer dropShopDBTable()
	gtest.C(t, func(t *gtest.T) {
		data, err := g.DB().Ctx(context.TODO()).Model("dbx_order").Partition("p1").All()
		t.AssertNil(err)
		t.Assert(len(data), 1)
	})
}
func setConfig() {
	gdb.SetConfig(gdb.Config{
		"default": gdb.ConfigGroup{
			gdb.ConfigNode{
				Host:                 "127.0.0.1",
				Port:                 "3306",
				User:                 "root",
				Pass:                 "111111",
				Name:                 "shop_db",
				Type:                 "mysql",
				Role:                 "master",
				Weight:               100,
				TimeMaintainDisabled: true,
				Debug:                true,
			},
		},
	})
}
func createShopDBTable() {
	sql := `CREATE TABLE dbx_order (
  id int(11) NOT NULL,
  sales_date date DEFAULT NULL,
  amount decimal(10,2) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
PARTITION BY RANGE (YEAR(sales_date))
(PARTITION p1 VALUES LESS THAN (2020) ENGINE = InnoDB,
 PARTITION p2 VALUES LESS THAN (2021) ENGINE = InnoDB,
 PARTITION p3 VALUES LESS THAN (2022) ENGINE = InnoDB,
 PARTITION p4 VALUES LESS THAN MAXVALUE ENGINE = InnoDB);`
	ctx := context.TODO()
	_, err := g.DB().Exec(ctx, sql)
	if err != nil {
		gtest.Fatal(err.Error())
	}
}
func insertShopDBData() {
	ctx := context.TODO()
	data := g.Slice{}
	year := 2020
	for i := 1; i <= 5; i++ {
		year++
		data = append(data, g.Map{
			"id":         i,
			"sales_date": fmt.Sprintf("%d-09-21", year),
			"amount":     fmt.Sprintf("1%d.21", i),
		})
	}
	_, err := g.DB().Model("dbx_order").Ctx(ctx).Data(data).Insert()
	if err != nil {
		gtest.Error(err)
	}
}
func dropShopDBTable() {
	ctx := context.TODO()
	if _, err := g.DB().Exec(ctx, "DROP TABLE IF EXISTS `dbx_order`"); err != nil {
		gtest.Error(err)
	}
}
