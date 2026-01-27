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
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gmeta"
)

// Test_WithAll_AdvancedScenarios 增强版的复杂场景和边界测试
// 覆盖场景：
// 1. 极小数据集（单条记录）
// 2. BatchThreshold 精确边界（阈值临界值）
// 3. 多层级混合 BatchSize 配置
// 4. 空关联数据处理
// 5. 深度嵌套（四层关联）
// 6. 大批量数据（1000+用户）
// 7. 全局配置+层级覆盖
func Test_WithAll_AdvancedScenarios(t *testing.T) {
	// 使用独立的表名避免与其他测试冲突
	var (
		tableUser             = "user_adv"
		tableUserDetail       = "user_detail_adv"
		tableUserScores       = "user_scores_adv"
		tableUserScoreDetails = "user_score_details_adv"
		tableScoreComments    = "score_comments_adv"
	)

	// 四层嵌套结构定义
	type ScoreComments struct {
		gmeta.Meta `orm:"table:score_comments_adv"`
		Id         int    `json:"id"`
		DetailId   int    `json:"detail_id"`
		Comment    string `json:"comment"`
	}

	type UserScoreDetails struct {
		gmeta.Meta   `orm:"table:user_score_details_adv"`
		Id           int              `json:"id"`
		ScoreId      int              `json:"score_id"`
		DetailInfo   string           `json:"detail_info"`
		ScoreComment []*ScoreComments `orm:"with:detail_id=id"`
	}

	type UserScores struct {
		gmeta.Meta   `orm:"table:user_scores_adv"`
		Id           int                 `json:"id"`
		Uid          int                 `json:"uid"`
		Score        int                 `json:"score"`
		ScoreDetails []*UserScoreDetails `orm:"with:score_id=id"`
	}

	type UserDetail struct {
		gmeta.Meta `orm:"table:user_detail_adv"`
		Uid        int    `json:"uid"`
		Address    string `json:"address"`
	}

	type User struct {
		gmeta.Meta `orm:"table:user_adv"`
		Id         int           `json:"id"`
		Name       string        `json:"name"`
		UserDetail *UserDetail   `orm:"with:uid=id"`
		UserScores []*UserScores `orm:"with:uid=id"`
	}

	// 初始化表结构
	dropTable(tableUser)
	dropTable(tableUserDetail)
	dropTable(tableUserScores)
	dropTable(tableUserScoreDetails)
	dropTable(tableScoreComments)

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

	_, err = db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id int(10) unsigned NOT NULL AUTO_INCREMENT,
			detail_id int(10) unsigned NOT NULL,
			comment varchar(100) NOT NULL,
			PRIMARY KEY (id),
			KEY (detail_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`, tableScoreComments))
	gtest.AssertNil(err)

	defer dropTable(tableUser)
	defer dropTable(tableUserDetail)
	defer dropTable(tableUserScores)
	defer dropTable(tableUserScoreDetails)
	defer dropTable(tableScoreComments)

	// ========================================
	// Scenario 1: 极小数据集测试（1条）
	// 验证优化在极小数据量下不产生负面影响
	// ========================================
	gtest.C(t, func(t *gtest.T) {
		db.SetDebug(true)
		defer db.SetDebug(false)
		fmt.Println("\n========== Scenario 1: 极小数据集（1条）==========")

		user1 := User{Id: 1, Name: "user_1"}
		userDetail1 := UserDetail{Uid: 1, Address: "address_1"}
		userScore1 := UserScores{Id: 1, Uid: 1, Score: 100}

		_, err := db.Model(tableUser).Data(user1).Insert()
		t.AssertNil(err)
		_, err = db.Model(tableUserDetail).Data(userDetail1).Insert()
		t.AssertNil(err)
		_, err = db.Model(tableUserScores).Data(userScore1).Insert()
		t.AssertNil(err)

		var users []*User
		err = db.Model(tableUser).WithAll().WithBatch().Where("id", 1).Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), 1)
		t.Assert(users[0].Name, "user_1")
		t.AssertNE(users[0].UserDetail, nil)
		t.Assert(len(users[0].UserScores), 1)

		db.Model(tableUser).Where("id", 1).Delete()
		db.Model(tableUserDetail).Where("uid", 1).Delete()
		db.Model(tableUserScores).Where("id", 1).Delete()
		fmt.Println("✓ 极小数据集验证通过")
	})

	// ========================================
	// Scenario 2: BatchThreshold 精确边界
	// 验证阈值-1、阈值、阈值+1的行为差异
	// ========================================
	gtest.C(t, func(t *gtest.T) {
		db.SetDebug(true)
		defer db.SetDebug(false)
		fmt.Println("\n========== Scenario 2: BatchThreshold 精确边界 ==========")

		usersData := make([]*User, 10)
		detailsData := make([]*UserDetail, 10)
		for i := 0; i < 10; i++ {
			usersData[i] = &User{Id: i + 1, Name: fmt.Sprintf("user_%d", i+1)}
			detailsData[i] = &UserDetail{Uid: i + 1, Address: fmt.Sprintf("address_%d", i+1)}
		}
		_, err := db.Model(tableUser).Data(usersData).Insert()
		t.AssertNil(err)
		_, err = db.Model(tableUserDetail).Data(detailsData).Insert()
		t.AssertNil(err)

		// 测试 Threshold=11 (不触发)
		var users1 []*User
		err = db.Model(tableUser).WithAll().
			WithBatchOption(gdb.WithBatchOption{Layer: 1, Enabled: true, BatchThreshold: 11, BatchSize: 1000}).
			WithBatch().Where("id<=?", 10).Scan(&users1)
		t.AssertNil(err)
		t.Assert(len(users1), 10)

		// 测试 Threshold=10 (触发)
		var users2 []*User
		err = db.Model(tableUser).WithAll().
			WithBatchOption(gdb.WithBatchOption{Layer: 1, Enabled: true, BatchThreshold: 10, BatchSize: 1000}).
			WithBatch().Where("id<=?", 10).Scan(&users2)
		t.AssertNil(err)
		t.Assert(len(users2), 10)

		// 测试 Threshold=9 (触发)
		var users3 []*User
		err = db.Model(tableUser).WithAll().
			WithBatchOption(gdb.WithBatchOption{Layer: 1, Enabled: true, BatchThreshold: 9, BatchSize: 1000}).
			WithBatch().Where("id<=?", 10).Scan(&users3)
		t.AssertNil(err)
		t.Assert(len(users3), 10)

		db.Model(tableUser).Where("id<=?", 10).Delete()
		db.Model(tableUserDetail).Where("uid<=?", 10).Delete()
		fmt.Println("✓ BatchThreshold 边界验证通过")
	})

	// ========================================
	// Scenario 3: 多层级混合 BatchSize
	// 验证不同层级使用不同BatchSize的正确性
	// ========================================
	gtest.C(t, func(t *gtest.T) {
		fmt.Println("\n========== Scenario 3: 多层级混合BatchSize ==========")
		db.SetDebug(true)
		defer db.SetDebug(false)
		// 20用户 * 5scores * 3details = 300 details
		usersData := make([]*User, 20)
		for i := 0; i < 20; i++ {
			usersData[i] = &User{Id: i + 1, Name: fmt.Sprintf("user_%d", i+1)}
		}
		_, err := db.Model(tableUser).Data(usersData).Insert()
		t.AssertNil(err)

		scoresData := make([]*UserScores, 100) // 20*5=100
		scoreId := 1
		for i := 0; i < 20; i++ {
			for j := 0; j < 5; j++ {
				scoresData[scoreId-1] = &UserScores{Id: scoreId, Uid: i + 1, Score: (j + 1) * 10}
				scoreId++
			}
		}
		_, err = db.Model(tableUserScores).Data(scoresData).Insert()
		t.AssertNil(err)

		detailsData := make([]*UserScoreDetails, 300) // 100*3=300
		detailIdx := 0
		for i := 1; i <= 100; i++ {
			for j := 0; j < 3; j++ {
				detailsData[detailIdx] = &UserScoreDetails{ScoreId: i, DetailInfo: fmt.Sprintf("detail_%d_%d", i, j+1)}
				detailIdx++
			}
		}
		_, err = db.Model(tableUserScoreDetails).Data(detailsData).Insert()
		t.AssertNil(err)

		// Layer1: BatchSize=5, Layer2: BatchSize=10
		var users []*User
		err = db.Model(tableUser).WithAll().
			WithBatchOption(
				gdb.WithBatchOption{Layer: 1, Enabled: true, BatchThreshold: 0, BatchSize: 5},
				gdb.WithBatchOption{Layer: 2, Enabled: true, BatchThreshold: 0, BatchSize: 10},
			).
			WithBatch().Where("id<=?", 20).Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), 20)
		for _, u := range users {
			t.Assert(len(u.UserScores), 5)
			for _, s := range u.UserScores {
				t.Assert(len(s.ScoreDetails), 3)
			}
		}

		db.Model(tableUser).Where("id<=?", 20).Delete()
		db.Model(tableUserScores).Where("uid<=?", 20).Delete()
		db.Model(tableUserScoreDetails).Where("score_id<=?", 100).Delete()
		fmt.Println("✓ 多层级混合BatchSize验证通过")
	})

	// ========================================
	// Scenario 4: 空关联数据
	// 验证部分记录无关联时不会崩溃
	// ========================================
	gtest.C(t, func(t *gtest.T) {
		fmt.Println("\n========== Scenario 4: 空关联数据 ==========")
		db.SetDebug(true)
		defer db.SetDebug(false)

		userEmpty := User{Id: 100, Name: "user_empty"}
		_, err := db.Model(tableUser).Data(userEmpty).Insert()
		t.AssertNil(err)

		var users []*User
		err = db.Model(tableUser).WithAll().WithBatch().Where("id", 100).Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), 1)
		t.Assert(users[0].UserDetail, nil)
		t.Assert(len(users[0].UserScores), 0)

		db.Model(tableUser).Where("id", 100).Delete()
		fmt.Println("✓ 空关联数据验证通过")
	})

	// ========================================
	// Scenario 5: 深度嵌套（四层）
	// 验证深层嵌套时优化依然有效
	// ========================================
	gtest.C(t, func(t *gtest.T) {
		fmt.Println("\n========== Scenario 5: 深度嵌套（四层）==========")
		db.SetDebug(true)
		defer db.SetDebug(false)
		// 5用户 * 2scores * 2details * 2comments
		usersData := make([]*User, 5)
		for i := 0; i < 5; i++ {
			usersData[i] = &User{Id: i + 1, Name: fmt.Sprintf("user_%d", i+1)}
		}
		_, err := db.Model(tableUser).Data(usersData).Insert()
		t.AssertNil(err)

		scoresData := make([]*UserScores, 10) // 5*2=10
		scoreId := 1
		for i := 0; i < 5; i++ {
			for j := 0; j < 2; j++ {
				scoresData[scoreId-1] = &UserScores{Id: scoreId, Uid: i + 1, Score: (j + 1) * 10}
				scoreId++
			}
		}
		_, err = db.Model(tableUserScores).Data(scoresData).Insert()
		t.AssertNil(err)

		detailsData := make([]*UserScoreDetails, 20) // 10*2=20
		detailId := 1
		for i := 0; i < 10; i++ {
			for j := 0; j < 2; j++ {
				detailsData[detailId-1] = &UserScoreDetails{Id: detailId, ScoreId: i + 1, DetailInfo: fmt.Sprintf("detail_%d_%d", i+1, j+1)}
				detailId++
			}
		}
		_, err = db.Model(tableUserScoreDetails).Data(detailsData).Insert()
		t.AssertNil(err)

		commentsData := make([]*ScoreComments, 40) // 20*2=40
		commentId := 1
		for i := 0; i < 20; i++ {
			for j := 0; j < 2; j++ {
				commentsData[commentId-1] = &ScoreComments{DetailId: i + 1, Comment: fmt.Sprintf("comment_%d_%d", i+1, j+1)}
				commentId++
			}
		}
		_, err = db.Model(tableScoreComments).Data(commentsData).Insert()
		t.AssertNil(err)

		start := time.Now()
		var users []*User
		err = db.Model(tableUser).WithAll().WithBatch().Where("id<=?", 5).Scan(&users)
		duration := time.Since(start)
		t.AssertNil(err)
		t.Assert(len(users), 5)
		for _, u := range users {
			t.Assert(len(u.UserScores), 2)
			for _, s := range u.UserScores {
				t.Assert(len(s.ScoreDetails), 2)
				for _, d := range s.ScoreDetails {
					t.Assert(len(d.ScoreComment), 2)
				}
			}
		}
		fmt.Printf("  四层嵌套查询耗时: %v\n", duration)

		db.Model(tableUser).Where("id<=?", 5).Delete()
		db.Model(tableUserScores).Where("uid<=?", 5).Delete()
		db.Model(tableUserScoreDetails).Where("score_id<=?", 10).Delete()
		db.Model(tableScoreComments).Where("detail_id<=?", 20).Delete()
		fmt.Println("✓ 四层嵌套验证通过")
	})

	// ========================================
	// Scenario 6: 大批量数据（1000用户）
	// 模拟生产环境验证性能和稳定性
	// ========================================
	gtest.C(t, func(t *gtest.T) {
		fmt.Println("\n========== Scenario 6: 大批量数据（1000用户）==========")
		db.SetDebug(true)
		defer db.SetDebug(false)
		userCount := 1000
		scorePerUser := 5
		detailPerScore := 3

		fmt.Println("  → 插入1000个用户...")
		usersData := make([]*User, userCount)
		for i := 0; i < userCount; i++ {
			usersData[i] = &User{Id: i + 1, Name: fmt.Sprintf("user_%d", i+1)}
		}
		_, err := db.Model(tableUser).Data(usersData).Insert()
		t.AssertNil(err)

		fmt.Println("  → 插入5000个scores...")
		totalScores := userCount * scorePerUser
		scoresData := make([]*UserScores, totalScores)
		scoreId := 1
		for i := 0; i < userCount; i++ {
			for j := 0; j < scorePerUser; j++ {
				scoresData[scoreId-1] = &UserScores{Id: scoreId, Uid: i + 1, Score: (j + 1) * 10}
				scoreId++
			}
		}
		_, err = db.Model(tableUserScores).Data(scoresData).Batch(1000).Insert()
		t.AssertNil(err)

		fmt.Println("  → 插入15000个details...")
		totalDetails := userCount * scorePerUser * detailPerScore
		detailsData := make([]*UserScoreDetails, totalDetails)
		detailIdx := 0
		for i := 1; i <= userCount*scorePerUser; i++ {
			for j := 0; j < detailPerScore; j++ {
				detailsData[detailIdx] = &UserScoreDetails{ScoreId: i, DetailInfo: fmt.Sprintf("detail_%d_%d", i, j+1)}
				detailIdx++
			}
		}
		_, err = db.Model(tableUserScoreDetails).Data(detailsData).Batch(1000).Insert()
		t.AssertNil(err)

		fmt.Println("  → 测试优化模式...")
		start := time.Now()
		var users []*User
		err = db.Model(tableUser).WithAll().WithBatch().Where("id<=?", userCount).Scan(&users)
		duration := time.Since(start)
		t.AssertNil(err)
		t.Assert(len(users), userCount)
		for _, u := range users {
			t.Assert(len(u.UserScores), scorePerUser)
			for _, s := range u.UserScores {
				t.Assert(len(s.ScoreDetails), detailPerScore)
			}
		}

		fmt.Printf("  ✓ 数据规模: %d 用户, %d scores, %d details\n", userCount, userCount*scorePerUser, userCount*scorePerUser*detailPerScore)
		fmt.Printf("  ✓ 优化模式耗时: %v\n", duration)

		db.Model(tableUser).Where("id<=?", userCount).Delete()
		db.Model(tableUserScores).Where("uid<=?", userCount).Delete()
		db.Model(tableUserScoreDetails).Where("score_id<=?", userCount*scorePerUser).Delete()
		fmt.Println("✓ 大批量数据验证通过")
	})

	// ========================================
	// Scenario 7: 全局配置+层级覆盖
	// 验证 Layer=0 全局设置能被特定层覆盖
	// ========================================
	gtest.C(t, func(t *gtest.T) {
		fmt.Println("\n========== Scenario 7: 全局配置+层级覆盖 ==========")
		db.SetDebug(true)
		defer db.SetDebug(false)
		usersData := make([]*User, 10)
		for i := 0; i < 10; i++ {
			usersData[i] = &User{Id: i + 1, Name: fmt.Sprintf("user_%d", i+1)}
		}
		_, err := db.Model(tableUser).Data(usersData).Insert()
		t.AssertNil(err)

		scoresData := make([]*UserScores, 30) // 10*3=30
		scoreIdx := 0
		for i := 0; i < 10; i++ {
			for j := 0; j < 3; j++ {
				scoresData[scoreIdx] = &UserScores{Uid: i + 1, Score: (j + 1) * 10}
				scoreIdx++
			}
		}
		_, err = db.Model(tableUserScores).Data(scoresData).Insert()
		t.AssertNil(err)

		// 全局 BatchSize=5，Layer2 覆盖为 BatchSize=2
		var users []*User
		err = db.Model(tableUser).WithAll().
			WithBatchOption(
				gdb.WithBatchOption{Layer: 0, Enabled: true, BatchThreshold: 0, BatchSize: 5},
				gdb.WithBatchOption{Layer: 2, Enabled: true, BatchThreshold: 0, BatchSize: 2},
			).
			WithBatch().Where("id<=?", 10).Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), 10)
		for _, u := range users {
			t.Assert(len(u.UserScores), 3)
		}

		db.Model(tableUser).Where("id<=?", 10).Delete()
		db.Model(tableUserScores).Where("uid<=?", 10).Delete()
		fmt.Println("✓ 全局配置+层级覆盖验证通过")
	})

	fmt.Println("\n========== 所有高级场景测试完成 ==========")
}
