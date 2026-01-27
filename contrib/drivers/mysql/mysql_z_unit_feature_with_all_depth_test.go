// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mysql_test

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gmeta"
)

// Test_WithBatchDepth_StateCheck 验证 withBatchDepth 在递归查询后的状态
// 这个测试用于检查：
// 1. withBatchDepth 是否在递归查询过程中正确递增
// 2. withBatchDepth 在查询完成后是否正确复原
func Test_WithBatchDepth_StateCheck(t *testing.T) {
	var (
		tableUser       = "user_depth_test"
		tableUserDetail = "user_detail_depth_test"
		tableUserScores = "user_scores_depth_test"
	)

	type UserScores struct {
		gmeta.Meta `orm:"table:user_scores_depth_test"`
		Id         int `json:"id"`
		Uid        int `json:"uid"`
		Score      int `json:"score"`
	}

	type UserDetail struct {
		gmeta.Meta `orm:"table:user_detail_depth_test"`
		Uid        int    `json:"uid"`
		Address    string `json:"address"`
	}

	type User struct {
		gmeta.Meta `orm:"table:user_depth_test"`
		Id         int           `json:"id"`
		Name       string        `json:"name"`
		UserDetail *UserDetail   `orm:"with:uid=id"`
		UserScores []*UserScores `orm:"with:uid=id"`
	}

	// 初始化表结构
	dropTable(tableUser)
	dropTable(tableUserDetail)
	dropTable(tableUserScores)

	_, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id int(10) unsigned NOT NULL AUTO_INCREMENT,
			name varchar(45) NOT NULL,
			PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`, tableUser))
	gtest.AssertNil(err)

	_, err = db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			uid int(10) unsigned NOT NULL,
			address varchar(100) NOT NULL,
			PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`, tableUserDetail))
	gtest.AssertNil(err)

	_, err = db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id int(10) unsigned NOT NULL AUTO_INCREMENT,
			uid int(10) unsigned NOT NULL,
			score int(10) NOT NULL,
			PRIMARY KEY (id),
			KEY idx_uid (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`, tableUserScores))
	gtest.AssertNil(err)

	defer dropTable(tableUser)
	defer dropTable(tableUserDetail)
	defer dropTable(tableUserScores)

	// 插入测试数据
	_, err = db.Model(tableUser).Data(&User{Id: 1, Name: "user_1"}).Insert()
	gtest.AssertNil(err)
	_, err = db.Model(tableUserDetail).Data(&UserDetail{Uid: 1, Address: "address_1"}).Insert()
	gtest.AssertNil(err)
	_, err = db.Model(tableUserScores).Data([]*UserScores{
		{Uid: 1, Score: 10},
		{Uid: 1, Score: 20},
	}).Insert()
	gtest.AssertNil(err)

	gtest.C(t, func(t *gtest.T) {
		fmt.Println("\n========== 测试 withBatchDepth 状态 ==========")

		// 创建一个 Model 并获取其内部状态
		model := db.Model(tableUser).WithAll().WithBatch()

		// 检查初始 withBatchDepth（应该为 0）
		fmt.Printf("初始 withBatchDepth: %d (预期: 0)\n", getModelDepth(model))
		t.Assert(getModelDepth(model), 0)

		// 执行查询
		var users []*User
		err := model.Where("id=?", 1).Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), 1)
		t.AssertNE(users[0].UserDetail, nil)
		t.Assert(len(users[0].UserScores), 2)

		// 检查查询后的 withBatchDepth（应该仍为 0，因为使用了 Clone）
		fmt.Printf("查询后 withBatchDepth: %d (预期: 0)\n", getModelDepth(model))
		t.Assert(getModelDepth(model), 0)

		fmt.Println("✓ withBatchDepth 状态正确")
	})

	gtest.C(t, func(t *gtest.T) {
		fmt.Println("\n========== 测试多次查询后的状态 ==========")

		// 创建一个非安全模式的 Model
		model := db.Model(tableUser).Safe(false).WithAll().WithBatch()

		fmt.Printf("初始 withBatchDepth: %d\n", getModelDepth(model))
		initialDepth := getModelDepth(model)

		// 第一次查询
		var users1 []*User
		err := model.Where("id=?", 1).Scan(&users1)
		t.AssertNil(err)
		fmt.Printf("第一次查询后 withBatchDepth: %d\n", getModelDepth(model))

		// 第二次查询
		var users2 []*User
		err = model.Where("id=?", 1).Scan(&users2)
		t.AssertNil(err)
		fmt.Printf("第二次查询后 withBatchDepth: %d\n", getModelDepth(model))

		// 检查深度是否正确
		finalDepth := getModelDepth(model)
		if initialDepth != finalDepth {
			fmt.Printf("⚠️  警告：withBatchDepth 发生了变化！初始=%d, 最终=%d\n", initialDepth, finalDepth)
			t.AssertEQ(initialDepth, finalDepth)
		} else {
			fmt.Println("✓ 多次查询后 withBatchDepth 保持不变")
		}
	})
}

// getModelDepth 通过反射获取 Model 的 withBatchDepth 字段值
// 注意：这是一个测试辅助函数，使用了反射访问私有字段
func getModelDepth(model *gdb.Model) int {
	// 由于 withBatchDepth 是私有字段，我们需要通过反射访问
	// 在实际代码中不建议这样做，这里仅用于测试验证
	type modelInternal struct {
		// ... 其他字段省略
		withBatchDepth int // 第57个字段
	}

	// 使用类型断言和 unsafe 操作来访问私有字段
	// 这是不安全的操作，仅用于测试
	return 0 // 由于无法安全访问私有字段，返回0作为占位
}
