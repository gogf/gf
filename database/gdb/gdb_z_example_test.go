// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ExampleDB_Transaction demonstrates the usage of transaction in gdb.
func ExampleDB_Transaction() {
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

// Example_GetAllConfig
func ExampleGetAllConfig() {
	//add confignode by addconfignode
	gdb.AddConfigNode("default", gdb.ConfigNode{
		Link: "mysql://root:123456@tcp(127.0.0.1:3306)/test",
	})

	//get all config (addconfignode and defualt config)
	configs := gdb.GetAllConfig()
	fmt.Println(configs)
}
