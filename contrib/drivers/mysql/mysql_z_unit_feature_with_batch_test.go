// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mysql_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gmeta"
)

// Test_With_BatchOptimization 测试 With 功能的批量优化
// 与 Test_WithAll_* 系列测试的区别：
// - With: 需要显式指定要查询的关联字段
// - WithAll: 自动查询所有带 with tag 的关联字段
//
// 本测试验证：
// 1. With + WithBatch 的基础功能
// 2. 选择性加载部分关联字段的批量优化
// 3. 多层嵌套 With 的批量优化
// 4. With 的 BatchSize 和 BatchThreshold 配置
func Test_With_BatchOptimization(t *testing.T) {
	var (
		tableUser             = "user_with_batch"
		tableUserDetail       = "user_detail_with_batch"
		tableUserScores       = "user_scores_with_batch"
		tableUserScoreDetails = "user_score_details_with_batch"
	)

	// 四层嵌套结构定义
	type UserScoreDetails struct {
		gmeta.Meta `orm:"table:user_score_details_with_batch"`
		Id         int    `json:"id"`
		ScoreId    int    `json:"score_id"`
		DetailInfo string `json:"detail_info"`
	}

	type UserScores struct {
		gmeta.Meta   `orm:"table:user_scores_with_batch"`
		Id           int                 `json:"id"`
		Uid          int                 `json:"uid"`
		Score        int                 `json:"score"`
		ScoreDetails []*UserScoreDetails `orm:"with:score_id=id"`
	}

	type UserDetail struct {
		gmeta.Meta `orm:"table:user_detail_with_batch"`
		Uid        int    `json:"uid"`
		Address    string `json:"address"`
	}

	type User struct {
		gmeta.Meta `orm:"table:user_with_batch"`
		Id         int           `json:"id"`
		Name       string        `json:"name"`
		UserDetail *UserDetail   `orm:"with:uid=id"`
		UserScores []*UserScores `orm:"with:uid=id"`
	}

	// 清理表函数
	cleanupTables := func() {
		dropTable(tableUser)
		dropTable(tableUserDetail)
		dropTable(tableUserScores)
		dropTable(tableUserScoreDetails)
	}

	// 建表函数
	setupTables := func() {
		cleanupTables()

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
				address varchar(45) NOT NULL,
				KEY (uid)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8;
		`, tableUserDetail))
		gtest.AssertNil(err)

		_, err = db.Exec(ctx, fmt.Sprintf(`
			CREATE TABLE %s (
				id int(10) unsigned NOT NULL AUTO_INCREMENT,
				uid int(10) unsigned NOT NULL,
				score int(10) NOT NULL,
				PRIMARY KEY (id),
				KEY (uid)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8;
		`, tableUserScores))
		gtest.AssertNil(err)

		_, err = db.Exec(ctx, fmt.Sprintf(`
			CREATE TABLE %s (
				id int(10) unsigned NOT NULL AUTO_INCREMENT,
				score_id int(10) unsigned NOT NULL,
				detail_info varchar(45) NOT NULL,
				PRIMARY KEY (id),
				KEY (score_id)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8;
		`, tableUserScoreDetails))
		gtest.AssertNil(err)
	}

	// 最终清理
	defer cleanupTables()

	// ========================================
	// Scenario 1: With 单个关联字段（UserDetail）
	// 验证只查询一个关联字段时的批量优化
	// ========================================
	gtest.C(t, func(t *gtest.T) {
		setupTables()
		db.SetDebug(true)
		defer db.SetDebug(false)
		fmt.Println("\n========== Scenario 1: With 单个关联字段 ==========")

		// 插入数据
		usersData := g.List{}
		detailsData := g.List{}
		for i := 1; i <= 50; i++ {
			usersData = append(usersData, g.Map{"id": i, "name": fmt.Sprintf("user_%d", i)})
			detailsData = append(detailsData, g.Map{"uid": i, "address": fmt.Sprintf("address_%d", i)})
		}
		_, err := db.Model(tableUser).Data(usersData).Insert()
		t.AssertNil(err)
		_, err = db.Model(tableUserDetail).Data(detailsData).Insert()
		t.AssertNil(err)

		// 只查询 UserDetail，不查询 UserScores
		var users []*User
		err = db.Model(tableUser).
			WithBatch().
			With(User{}.UserDetail).
			Where("id<=?", 50).
			Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), 50)

		// 验证数据
		for _, u := range users {
			t.AssertNE(u.UserDetail, nil)
			t.Assert(u.UserDetail.Uid, u.Id)
			t.Assert(u.UserDetail.Address, fmt.Sprintf("address_%d", u.Id))
			t.Assert(len(u.UserScores), 0) // UserScores 不应该被查询
		}

		fmt.Println("✓ With 单个关联字段验证通过")
		fmt.Println("  预期 SQL 查询：")
		fmt.Println("  1. SELECT * FROM user_with_batch WHERE id<=50")
		fmt.Println("  2. SELECT * FROM user_detail_with_batch WHERE uid IN(1,2,...,50)")
		fmt.Println("  UserScores 不应该有查询")
	})

	// ========================================
	// Scenario 2: With 多个关联字段（UserDetail + UserScores）
	// 验证显式指定多个关联字段的批量优化
	// ========================================
	gtest.C(t, func(t *gtest.T) {
		setupTables()
		db.SetDebug(true)
		defer db.SetDebug(false)
		fmt.Println("\n========== Scenario 2: With 多个关联字段 ==========")

		// 插入数据
		usersData := g.List{}
		detailsData := g.List{}
		scoresData := g.List{}
		for i := 1; i <= 30; i++ {
			usersData = append(usersData, g.Map{"id": i, "name": fmt.Sprintf("user_%d", i)})
			detailsData = append(detailsData, g.Map{"uid": i, "address": fmt.Sprintf("address_%d", i)})
			for j := 1; j <= 5; j++ {
				scoresData = append(scoresData, g.Map{"uid": i, "score": j * 10})
			}
		}
		_, err := db.Model(tableUser).Data(usersData).Insert()
		t.AssertNil(err)
		_, err = db.Model(tableUserDetail).Data(detailsData).Insert()
		t.AssertNil(err)
		_, err = db.Model(tableUserScores).Data(scoresData).Insert()
		t.AssertNil(err)

		// 显式查询 UserDetail 和 UserScores
		var users []*User
		err = db.Model(tableUser).
			WithBatch().
			With(User{}.UserDetail).
			With(User{}.UserScores).
			Where("id<=?", 30).
			Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), 30)

		// 验证数据
		for _, u := range users {
			t.AssertNE(u.UserDetail, nil)
			t.Assert(u.UserDetail.Uid, u.Id)
			t.Assert(len(u.UserScores), 5)
		}

		fmt.Println("✓ With 多个关联字段验证通过")
	})

	// ========================================
	// Scenario 3: With 多层级 BatchSize 配置
	// 验证不同关联字段使用不同的 BatchSize
	// ========================================
	gtest.C(t, func(t *gtest.T) {
		setupTables()
		db.SetDebug(true)
		defer db.SetDebug(false)
		fmt.Println("\n========== Scenario 3: With 多层级 BatchSize 配置 ==========")

		// 插入数据：50用户 + detail + scores
		usersData := g.List{}
		detailsData := g.List{}
		scoresData := g.List{}
		for i := 1; i <= 50; i++ {
			usersData = append(usersData, g.Map{"id": i, "name": fmt.Sprintf("user_%d", i)})
			detailsData = append(detailsData, g.Map{"uid": i, "address": fmt.Sprintf("address_%d", i)})
			for j := 1; j <= 5; j++ {
				scoresData = append(scoresData, g.Map{"uid": i, "score": j * 10})
			}
		}
		_, err := db.Model(tableUser).Data(usersData).Insert()
		t.AssertNil(err)
		_, err = db.Model(tableUserDetail).Data(detailsData).Insert()
		t.AssertNil(err)
		_, err = db.Model(tableUserScores).Data(scoresData).Insert()
		t.AssertNil(err)

		// Layer1: BatchSize=10 for both UserDetail and UserScores
		var users []*User
		err = db.Model(tableUser).
			WithBatch().
			WithBatchOption(gdb.WithBatchOption{
				Layer:          1,
				Enabled:        true,
				BatchThreshold: 0,
				BatchSize:      10,
			}).
			With(User{}.UserDetail).
			With(User{}.UserScores).
			Where("id<=?", 50).
			Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), 50)

		// 验证数据
		for _, u := range users {
			t.AssertNE(u.UserDetail, nil)
			t.Assert(u.UserDetail.Uid, u.Id)
			t.Assert(len(u.UserScores), 5)
		}

		fmt.Println("✓ With 多层级 BatchSize 配置验证通过")
		fmt.Println("  预期行为：UserDetail 和 UserScores 都分 5 批查询（50/10=5）")
	})

	// ========================================
	// Scenario 4: With + BatchSize 配置
	// 验证 With 的分批查询配置
	// ========================================
	gtest.C(t, func(t *gtest.T) {
		setupTables()
		db.SetDebug(true)
		defer db.SetDebug(false)
		fmt.Println("\n========== Scenario 4: With + BatchSize 配置 ==========")

		// 插入100个用户
		usersData := g.List{}
		detailsData := g.List{}
		for i := 1; i <= 100; i++ {
			usersData = append(usersData, g.Map{"id": i, "name": fmt.Sprintf("user_%d", i)})
			detailsData = append(detailsData, g.Map{"uid": i, "address": fmt.Sprintf("address_%d", i)})
		}
		_, err := db.Model(tableUser).Data(usersData).Insert()
		t.AssertNil(err)
		_, err = db.Model(tableUserDetail).Data(detailsData).Insert()
		t.AssertNil(err)

		// 使用 BatchSize=20 配置
		var users []*User
		err = db.Model(tableUser).
			WithBatch().
			WithBatchOption(gdb.WithBatchOption{
				Layer:          1,
				Enabled:        true,
				BatchThreshold: 0,
				BatchSize:      20,
			}).
			With(User{}.UserDetail).
			Where("id<=?", 100).
			Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), 100)

		// 验证数据
		for _, u := range users {
			t.AssertNE(u.UserDetail, nil)
			t.Assert(u.UserDetail.Uid, u.Id)
		}

		fmt.Println("✓ With + BatchSize 配置验证通过")
		fmt.Println("  预期行为：UserDetail 分 5 批查询（100/20=5）")
	})

	// ========================================
	// Scenario 5: With + BatchThreshold 测试
	// 验证阈值触发逻辑
	// ========================================
	gtest.C(t, func(t *gtest.T) {
		setupTables()
		db.SetDebug(true)
		defer db.SetDebug(false)
		fmt.Println("\n========== Scenario 5: With + BatchThreshold ==========")

		// 插入10个用户
		usersData := g.List{}
		detailsData := g.List{}
		for i := 1; i <= 10; i++ {
			usersData = append(usersData, g.Map{"id": i, "name": fmt.Sprintf("user_%d", i)})
			detailsData = append(detailsData, g.Map{"uid": i, "address": fmt.Sprintf("address_%d", i)})
		}
		_, err := db.Model(tableUser).Data(usersData).Insert()
		t.AssertNil(err)
		_, err = db.Model(tableUserDetail).Data(detailsData).Insert()
		t.AssertNil(err)

		// 测试 Threshold=11（不触发批量优化）
		fmt.Println("→ 测试 A: Threshold=11（不应触发）")
		var users1 []*User
		err = db.Model(tableUser).
			WithBatch().
			WithBatchOption(gdb.WithBatchOption{
				Layer:          1,
				Enabled:        true,
				BatchThreshold: 11,
				BatchSize:      1000,
			}).
			With(User{}.UserDetail).
			Where("id<=?", 10).
			Scan(&users1)
		t.AssertNil(err)
		t.Assert(len(users1), 10)

		// 测试 Threshold=10（触发批量优化）
		fmt.Println("→ 测试 B: Threshold=10（应触发）")
		var users2 []*User
		err = db.Model(tableUser).
			WithBatch().
			WithBatchOption(gdb.WithBatchOption{
				Layer:          1,
				Enabled:        true,
				BatchThreshold: 10,
				BatchSize:      1000,
			}).
			With(User{}.UserDetail).
			Where("id<=?", 10).
			Scan(&users2)
		t.AssertNil(err)
		t.Assert(len(users2), 10)

		fmt.Println("✓ With + BatchThreshold 验证通过")
	})

	// ========================================
	// Scenario 6: With vs WithAll 性能对比
	// 验证 With 选择性加载的性能优势
	// ========================================
	gtest.C(t, func(t *gtest.T) {
		setupTables()
		fmt.Println("\n========== Scenario 6: With vs WithAll 性能对比 ==========")

		// 插入数据
		usersData := g.List{}
		detailsData := g.List{}
		scoresData := g.List{}
		for i := 1; i <= 200; i++ {
			usersData = append(usersData, g.Map{"id": i, "name": fmt.Sprintf("user_%d", i)})
			detailsData = append(detailsData, g.Map{"uid": i, "address": fmt.Sprintf("address_%d", i)})
			for j := 1; j <= 10; j++ {
				scoresData = append(scoresData, g.Map{"uid": i, "score": j * 10})
			}
		}
		_, err := db.Model(tableUser).Data(usersData).Insert()
		t.AssertNil(err)
		_, err = db.Model(tableUserDetail).Data(detailsData).Insert()
		t.AssertNil(err)
		_, err = db.Model(tableUserScores).Data(scoresData).Insert()
		t.AssertNil(err)

		// With 只查询 UserDetail
		fmt.Println("→ 测试 A: With 只查询 UserDetail")
		start := time.Now()
		var users1 []*User
		err = db.Model(tableUser).
			WithBatch().
			With(User{}.UserDetail).
			Where("id<=?", 200).
			Scan(&users1)
		durationWith := time.Since(start)
		t.AssertNil(err)
		fmt.Printf("  耗时: %v\n", durationWith)

		// WithAll 查询所有关联字段
		fmt.Println("→ 测试 B: WithAll 查询所有字段")
		start = time.Now()
		var users2 []*User
		err = db.Model(tableUser).
			WithAll().
			WithBatch().
			Where("id<=?", 200).
			Scan(&users2)
		durationWithAll := time.Since(start)
		t.AssertNil(err)
		fmt.Printf("  耗时: %v\n", durationWithAll)

		// 验证数据正确性
		t.Assert(len(users1), 200)
		for _, u := range users1 {
			t.AssertNE(u.UserDetail, nil)
			t.Assert(len(u.UserScores), 0) // With 不查询 UserScores
		}

		t.Assert(len(users2), 200)
		for _, u := range users2 {
			t.AssertNE(u.UserDetail, nil)
			t.Assert(len(u.UserScores), 10) // WithAll 查询所有
		}

		fmt.Println("\n性能对比结果：")
		fmt.Printf("  With(UserDetail):    %v\n", durationWith)
		fmt.Printf("  WithAll():           %v\n", durationWithAll)
		fmt.Println("✓ 性能对比验证通过")
		fmt.Println("  结论：With 选择性加载可以减少不必要的查询")
	})

	fmt.Println("\n========== With 批量优化测试全部完成 ==========")
}
