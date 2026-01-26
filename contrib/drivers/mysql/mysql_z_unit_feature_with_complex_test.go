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
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gmeta"
)

// Test_WithAll_Complex 大数据量级复杂条件及软删除测试
func Test_WithAll_Complex(t *testing.T) {
	var (
		tableUser       = "user_complex"
		tableUserDetail = "user_detail_complex"
		tableUserScores = "user_scores_complex"
	)

	// 1. 软删除相关数据结构
	type UserDetailSoft struct {
		gmeta.Meta `orm:"table:user_detail_complex"`
		Uid        int         `json:"uid"`
		Address    string      `json:"address"`
		DeleteAt   *gtime.Time `json:"delete_at" orm:"delete_at"`
	}

	type UserScoresSoft struct {
		gmeta.Meta `orm:"table:user_scores_complex"`
		Id         int         `json:"id"`
		Uid        int         `json:"uid"`
		Score      int         `json:"score"`
		DeleteAt   *gtime.Time `json:"delete_at" orm:"delete_at"`
	}

	type UserSoft struct {
		gmeta.Meta `orm:"table:user_complex"`
		Id         int               `json:"id"`
		Name       string            `json:"name"`
		UserDetail *UserDetailSoft   `orm:"with:uid=id, unscoped:true"`
		UserScores []*UserScoresSoft `orm:"with:uid=id, unscoped:true"`
	}

	// 2. 复杂条件筛选数据结构
	type UserDetailCond struct {
		gmeta.Meta `orm:"table:user_detail_complex"`
		Uid        int    `json:"uid"`
		Address    string `json:"address"`
	}

	type UserScoresCond struct {
		gmeta.Meta `orm:"table:user_scores_complex"`
		Id         int `json:"id"`
		Uid        int `json:"uid"`
		Score      int `json:"score"`
	}

	type UserCond struct {
		gmeta.Meta `orm:"table:user_complex"`
		Id         int               `json:"id"`
		Name       string            `json:"name"`
		UserDetail *UserDetailCond   `orm:"with:uid=id, where:uid > 3"`
		UserScores []*UserScoresCond `orm:"with:uid=id, where:score>1 and score<5, order:score desc"`
	}

	// 初始化表结构
	dropTable(tableUser)
	dropTable(tableUserDetail)
	dropTable(tableUserScores)
	db.SetDebug(true)
	defer db.SetDebug(false)

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
			delete_at datetime DEFAULT NULL,
			PRIMARY KEY (uid)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`, tableUserDetail))
	gtest.AssertNil(err)

	_, err = db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id int(10) unsigned NOT NULL AUTO_INCREMENT,
			uid int(10) unsigned NOT NULL,
			score int(10) NOT NULL,
			delete_at datetime DEFAULT NULL,
			PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`, tableUserScores))
	gtest.AssertNil(err)

	defer dropTable(tableUser)
	defer dropTable(tableUserDetail)
	defer dropTable(tableUserScores)

	db.SetDebug(true)
	defer db.SetDebug(false)

	// ========================================
	// 数据初始化
	// ========================================
	const (
		userCount    = 100
		scorePerUser = 10
	)

	fmt.Println("\n========== 开始初始化数据 ==========")

	// 1. 插入用户
	usersData := make(g.List, 0, userCount)
	for i := 1; i <= userCount; i++ {
		usersData = append(usersData, g.Map{"id": i, "name": fmt.Sprintf("user_%d", i)})
	}
	_, err = db.Model(tableUser).Data(usersData).Insert()
	gtest.AssertNil(err)

	// 2. 插入详情（部分标记为删除）
	detailsData := make(g.List, 0, userCount)
	for i := 1; i <= userCount; i++ {
		deleteAt := interface{}(nil)
		if i%2 == 0 { // 偶数用户详情被软删除
			deleteAt = gtime.Now()
		}
		detailsData = append(detailsData, g.Map{
			"uid":       i,
			"address":   fmt.Sprintf("address_%d", i),
			"delete_at": deleteAt,
		})
	}
	_, err = db.Model(tableUserDetail).Data(detailsData).Insert()
	gtest.AssertNil(err)

	// 3. 插入成绩（部分标记为删除，且分数各异）
	scoresData := make(g.List, 0, userCount*scorePerUser)
	for i := 1; i <= userCount; i++ {
		for j := 1; j <= scorePerUser; j++ {
			deleteAt := interface{}(nil)
			if j%2 == 0 { // 每个用户的偶数项成绩被软删除
				deleteAt = gtime.Now()
			}
			scoresData = append(scoresData, g.Map{
				"uid":       i,
				"score":     j, // score 从 1 到 10
				"delete_at": deleteAt,
			})
		}
	}
	_, err = db.Model(tableUserScores).Data(scoresData).Batch(500).Insert()
	gtest.AssertNil(err)

	fmt.Println("========== 数据初始化完成 ==========")

	// ========================================
	// Scenario 1: 验证 unscoped:true
	// ========================================
	gtest.C(t, func(t *gtest.T) {
		fmt.Println("\nScenario 1: 验证 unscoped:true (包含已软删除数据)")
		var users []*UserSoft
		err := db.Model(tableUser).WithAll().Where("id", 1).Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), 1)

		// user := users[0]
		// UserDetail 虽然 i=1 是奇数没删，但我们验证 i=2 的
		var usersAll []*UserSoft
		err = db.Model(tableUser).WithAll().Where("id IN (1,2)").Order("id").Scan(&usersAll)
		t.AssertNil(err)
		t.Assert(len(usersAll), 2)

		// user 1: detail not deleted
		t.AssertNE(usersAll[0].UserDetail, nil)
		t.AssertNil(usersAll[0].UserDetail.DeleteAt)

		// user 2: detail deleted, but unscoped:true should find it
		t.AssertNE(usersAll[1].UserDetail, nil)
		t.AssertNE(usersAll[1].UserDetail.DeleteAt, nil)

		// scores: half deleted, but unscoped:true should find all 10
		t.Assert(len(usersAll[0].UserScores), 10)
		deletedCount := 0
		for _, s := range usersAll[0].UserScores {
			if s.DeleteAt != nil {
				deletedCount++
			}
		}
		t.Assert(deletedCount, 5)
	})

	// ========================================
	// Scenario 2: 验证 where 和 order 条件
	// ========================================
	gtest.C(t, func(t *gtest.T) {
		fmt.Println("\nScenario 2: 验证 where 和 order 条件")
		// UserDetailCond: where:uid > 3
		// UserScoresCond: where:score>1 and score<5, order:score desc
		// Note: Normal scan will respect soft delete by default if not unscoped.
		// So for UserScores, score 2 and 4 are deleted (even items).
		// score 1, 3, 5, 7, 9 are NOT deleted.
		// where: score > 1 and score < 5 => scores 2, 3, 4
		// But 2 and 4 are deleted. So only score 3 should remain.

		var users []*UserCond
		err := db.Model(tableUser).WithAll().Where("id", 4).Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), 1)

		user := users[0]
		// Detail: uid=4 > 3, should be there.
		// BUT wait, user 4's detail is deleted (4%2==0).
		// Since UserDetailCond DOES NOT have unscoped:true, it should be nil!
		t.Assert(user.UserDetail, nil)

		// User 5
		users = nil
		err = db.Model(tableUser).WithAll().Where("id", 5).Scan(&users)
		t.AssertNil(err)
		user = users[0]
		// uid=5 > 3 and not deleted.
		t.AssertNE(user.UserDetail, nil)
		t.Assert(user.UserDetail.Uid, 5)

		// User 3
		users = nil
		err = db.Model(tableUser).WithAll().Where("id", 3).Scan(&users)
		t.AssertNil(err)
		user = users[0]
		// uid=3 NOT > 3.
		t.Assert(user.UserDetail, nil)

		// UserScores: score > 1 and score < 5 => 2, 3, 4.
		// 2 and 4 are deleted. Only 3 remains.
		t.Assert(len(user.UserScores), 1)
		t.Assert(user.UserScores[0].Score, 3)
	})

	// ========================================
	// Scenario 3: 大数据量分批查询与复杂条件结合
	// ========================================
	gtest.C(t, func(t *gtest.T) {
		fmt.Println("\nScenario 3: 大数据量分批查询与复杂条件结合")

		// 增加数据量以测试分批
		const largeUserCount = 500
		usersData := make(g.List, 0, largeUserCount)
		for i := 101; i <= 101+largeUserCount; i++ {
			usersData = append(usersData, g.Map{"id": i, "name": fmt.Sprintf("user_%d", i)})
		}
		_, err = db.Model(tableUser).Data(usersData).Insert()
		t.AssertNil(err)

		detailsData := make(g.List, 0, largeUserCount)
		for i := 101; i <= 101+largeUserCount; i++ {
			detailsData = append(detailsData, g.Map{
				"uid":     i,
				"address": fmt.Sprintf("address_%d", i),
			})
		}
		_, err = db.Model(tableUserDetail).Data(detailsData).Insert()
		t.AssertNil(err)

		scoresData := make(g.List, 0, largeUserCount*scorePerUser)
		for i := 101; i <= 101+largeUserCount; i++ {
			for j := 1; j <= scorePerUser; j++ {
				scoresData = append(scoresData, g.Map{
					"uid":   i,
					"score": j,
				})
			}
		}
		_, err = db.Model(tableUserScores).Data(scoresData).Batch(1000).Insert()
		t.AssertNil(err)

		// 执行分批查询
		var users []*UserCond
		err = db.Model(tableUser).WithAll().
			WithBatchOption(gdb.WithBatchOption{
				Layer:          0,
				Enabled:        true,
				BatchThreshold: 0,
				BatchSize:      100,
			}).
			WithBatch().
			Where("id > ?", 100).
			Scan(&users)

		t.AssertNil(err)
		t.Assert(len(users), largeUserCount+1)

		// 验证复杂条件在分批下是否依然有效
		// UserScoresCond: where:score>1 and score<5 => 2, 3, 4
		// 对于新插入的数据，没有标记删除，所以应该有 3 条 (2, 3, 4)
		// 且 order:score desc => 4, 3, 2
		for _, u := range users {
			t.Assert(len(u.UserScores), 3)
			t.Assert(u.UserScores[0].Score, 4)
			t.Assert(u.UserScores[1].Score, 3)
			t.Assert(u.UserScores[2].Score, 2)

			if u.Id > 3 {
				t.AssertNE(u.UserDetail, nil)
			} else {
				t.Assert(u.UserDetail, nil)
			}
		}
	})

	// ========================================
	// Scenario 4: 边界与异常情况
	// ========================================
	gtest.C(t, func(t *gtest.T) {
		fmt.Println("\nScenario 4: 边界与异常情况")

		// 1. 查询不存在的用户
		var users []*UserCond
		err := db.Model(tableUser).WithAll().Where("id", 99999).Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), 0)

		// 2. 关联表完全没数据 (清空 scores)
		_, err = db.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s", tableUserScores))
		t.AssertNil(err)

		err = db.Model(tableUser).WithAll().Where("id", 1).Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), 1)
		t.Assert(len(users[0].UserScores), 0)
	})
}
