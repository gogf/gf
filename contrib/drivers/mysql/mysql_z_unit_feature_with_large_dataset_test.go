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

// Test_WithAll_LargeDataset 大数据量级性能测试
// 数据规模：
// - 2000 个用户（是 Test_WithAll_PerformanceComparison 的 4 倍）
// - 每个用户 10 个 UserScore（共 20,000 个 score）
// - 每个 score 10 个 ScoreDetail（共 200,000 个 detail）
//
// 测试目标：
// 1. 验证大数据量下的查询正确性
// 2. 验证 BatchSize 分段查询的 SQL 正确性
// 3. 对比优化前后的性能差异
// 4. 验证不同 BatchSize 配置的影响
func Test_WithAll_LargeDataset(t *testing.T) {
	var (
		// 使用独立的表名避免与其他测试冲突
		tableUser             = "user_large"
		tableUserDetail       = "user_detail_large"
		tableUserScores       = "user_scores_large"
		tableUserScoreDetails = "user_score_details_large"
	)

	// 数据结构定义
	type UserScoreDetails struct {
		gmeta.Meta `orm:"table:user_score_details_large"`
		Id         int    `json:"id"`
		ScoreId    int    `json:"score_id"`
		DetailInfo string `json:"detail_info"`
	}

	type UserScores struct {
		gmeta.Meta   `orm:"table:user_scores_large"`
		Id           int                 `json:"id"`
		Uid          int                 `json:"uid"`
		Score        int                 `json:"score"`
		ScoreDetails []*UserScoreDetails `orm:"with:score_id=id"`
	}

	type UserDetail struct {
		gmeta.Meta `orm:"table:user_detail_large"`
		Uid        int    `json:"uid"`
		Address    string `json:"address"`
	}

	type User struct {
		gmeta.Meta `orm:"table:user_large"`
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

	_, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id int(10) unsigned NOT NULL AUTO_INCREMENT,
			name varchar(45) NOT NULL,
			PRIMARY KEY (id),
			KEY idx_name (name)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`, tableUser))
	gtest.AssertNil(err)

	_, err = db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			uid int(10) unsigned NOT NULL,
			address varchar(100) NOT NULL,
			PRIMARY KEY (uid),
			KEY idx_uid (uid)
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

	_, err = db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id int(10) unsigned NOT NULL AUTO_INCREMENT,
			score_id int(10) unsigned NOT NULL,
			detail_info varchar(200) NOT NULL,
			PRIMARY KEY (id),
			KEY idx_score_id (score_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`, tableUserScoreDetails))
	gtest.AssertNil(err)

	defer dropTable(tableUser)
	defer dropTable(tableUserDetail)
	defer dropTable(tableUserScores)
	defer dropTable(tableUserScoreDetails)

	// ========================================
	// 数据初始化（大数据量）
	// ========================================
	const (
		userCount      = 2000                         // 2000个用户
		scorePerUser   = 10                           // 每个用户10个score
		detailPerScore = 10                           // 每个score 10个detail（符合1:10规范）
		totalScores    = userCount * scorePerUser     // 20,000
		totalDetails   = totalScores * detailPerScore // 200,000
	)

	fmt.Println("\n========== 开始初始化大数据集 ==========")
	fmt.Printf("数据规模：%d 用户, %d scores, %d details\n", userCount, totalScores, totalDetails)

	// 1. 插入用户数据
	fmt.Println("→ 插入用户数据...")
	startTime := time.Now()
	usersData := make(g.List, 0, userCount)
	for i := 1; i <= userCount; i++ {
		usersData = append(usersData, g.Map{
			"id":   i,
			"name": fmt.Sprintf("user_%d", i),
		})
	}
	_, err = db.Model(tableUser).Data(usersData).Batch(1000).Insert()
	gtest.AssertNil(err)
	fmt.Printf("  用户数据插入完成，耗时: %v\n", time.Since(startTime))

	// 2. 插入用户详情
	fmt.Println("→ 插入用户详情...")
	startTime = time.Now()
	detailsData := make(g.List, 0, userCount)
	for i := 1; i <= userCount; i++ {
		detailsData = append(detailsData, g.Map{
			"uid":     i,
			"address": fmt.Sprintf("address_%d", i),
		})
	}
	_, err = db.Model(tableUserDetail).Data(detailsData).Batch(1000).Insert()
	gtest.AssertNil(err)
	fmt.Printf("  用户详情插入完成，耗时: %v\n", time.Since(startTime))

	// 3. 插入 UserScores
	fmt.Println("→ 插入 UserScores...")
	startTime = time.Now()
	scoresData := make(g.List, 0, totalScores)
	scoreId := 1
	for i := 1; i <= userCount; i++ {
		for j := 1; j <= scorePerUser; j++ {
			scoresData = append(scoresData, g.Map{
				"id":    scoreId,
				"uid":   i,
				"score": j * 10,
			})
			scoreId++
		}
	}
	_, err = db.Model(tableUserScores).Data(scoresData).Batch(1000).Insert()
	gtest.AssertNil(err)
	fmt.Printf("  UserScores 插入完成，耗时: %v\n", time.Since(startTime))

	// 4. 插入 ScoreDetails
	fmt.Println("→ 插入 ScoreDetails...")
	startTime = time.Now()
	scoreDetailsData := make(g.List, 0, totalDetails)
	for i := 1; i <= totalScores; i++ {
		for j := 1; j <= detailPerScore; j++ {
			scoreDetailsData = append(scoreDetailsData, g.Map{
				"score_id":    i,
				"detail_info": fmt.Sprintf("detail_%d_%d", i, j),
			})
		}
	}
	_, err = db.Model(tableUserScoreDetails).Data(scoreDetailsData).Batch(1000).Insert()
	gtest.AssertNil(err)
	fmt.Printf("  ScoreDetails 插入完成，耗时: %v\n", time.Since(startTime))

	fmt.Println("========== 数据初始化完成 ==========")

	// 数据验证辅助函数
	verifyUserData := func(users []*User, expectedCount int, checkDetails bool) {
		gtest.Assert(len(users), expectedCount)
		for _, user := range users {
			// 验证 UserDetail
			gtest.AssertNE(user.UserDetail, nil)
			gtest.Assert(user.UserDetail.Address, fmt.Sprintf("address_%d", user.Id))

			// 验证 UserScores
			gtest.Assert(len(user.UserScores), scorePerUser)
			for idx, score := range user.UserScores {
				gtest.Assert(score.Uid, user.Id)
				gtest.Assert(score.Score, (idx+1)*10)

				// 验证 ScoreDetails
				if checkDetails {
					gtest.Assert(len(score.ScoreDetails), detailPerScore)
					for detailIdx, detail := range score.ScoreDetails {
						gtest.Assert(detail.ScoreId, score.Id)
						expectedInfo := fmt.Sprintf("detail_%d_%d", score.Id, detailIdx+1)
						gtest.Assert(detail.DetailInfo, expectedInfo)
					}
				}
			}
		}
	}

	// 打开 SQL 日志（用于验证 SQL 正确性）
	db.SetDebug(true)
	defer db.SetDebug(false)

	// ========================================
	// Scenario 1: 查询前 100 个用户（不开启优化）
	// 验证：应该产生 N+1 查询
	// ========================================
	gtest.C(t, func(t *gtest.T) {
		fmt.Println("\n========== Scenario 1: 不开启优化（查询100个用户）==========")

		fmt.Println("执行查询...")
		startTime := time.Now()
		var users []*User
		err := db.Model(tableUser).WithAll().Where("id<=?", 100).Scan(&users)
		duration := time.Since(startTime)

		t.AssertNil(err)
		fmt.Printf("查询完成，耗时: %v\n", duration)

		fmt.Println("验证数据完整性...")
		verifyUserData(users, 100, true)

		fmt.Println("✓ 数据验证通过")
		fmt.Println("\n预期行为：应该看到大量的单次查询（N+1问题）")
	})

	// 临时关闭日志，避免输出过多
	db.SetDebug(false)

	// ========================================
	// Scenario 2: 查询前 100 个用户（开启优化，默认 BatchSize）
	// 验证：应该只有 4 次查询
	// ========================================
	gtest.C(t, func(t *gtest.T) {
		fmt.Println("\n========== Scenario 2: 开启优化，默认BatchSize（查询100个用户）==========")

		db.SetDebug(true)
		fmt.Println("执行查询...")
		startTime := time.Now()
		var users []*User
		err := db.Model(tableUser).WithAll().WithBatch().Where("id<=?", 100).Scan(&users)
		duration := time.Since(startTime)
		db.SetDebug(false)

		t.AssertNil(err)
		fmt.Printf("查询完成，耗时: %v\n", duration)

		fmt.Println("验证数据完整性...")
		verifyUserData(users, 100, true)

		fmt.Println("✓ 数据验证通过")
		fmt.Println("\n预期行为：")
		fmt.Println("  1. 主查询: SELECT ... FROM user_large WHERE id<=100")
		fmt.Println("  2. UserDetail: SELECT ... WHERE uid IN(1,2,...,100)")
		fmt.Println("  3. UserScores: SELECT ... WHERE uid IN(1,2,...,100)")
		fmt.Println("  4. ScoreDetails: SELECT ... WHERE score_id IN(...) [1000条]")
		fmt.Println("  总查询次数: 4 次")
	})

	// ========================================
	// Scenario 3: 查询全部 2000 个用户（BatchSize=50）
	// 验证：BatchSize 分段查询是否正确
	// ========================================
	gtest.C(t, func(t *gtest.T) {
		fmt.Println("\n========== Scenario 3: BatchSize=50（查询全部2000用户）==========")

		db.SetDebug(true)
		fmt.Println("执行查询...")
		startTime := time.Now()
		var users []*User
		err := db.Model(tableUser).WithAll().
			WithBatchOption(gdb.WithBatchOption{
				Layer:          0,
				Enabled:        true,
				BatchThreshold: 0,
				BatchSize:      50,
			}).
			WithBatch().
			Where("id<=?", 2000).
			Scan(&users)
		duration := time.Since(startTime)
		db.SetDebug(false)

		t.AssertNil(err)
		fmt.Printf("查询完成，耗时: %v\n", duration)

		fmt.Println("验证数据完整性（采样验证前50个用户）...")
		verifyUserData(users[:50], 50, true)

		fmt.Println("✓ 数据验证通过")
		fmt.Println("\n预期行为：")
		fmt.Println("  1. 主查询: 1 次")
		fmt.Println("  2. UserDetail: 40 次（2000/50=40）")
		fmt.Println("  3. UserScores: 40 次（2000/50=40）")
		fmt.Println("  4. ScoreDetails: 400 次（20000/50=400）")
		fmt.Println("  总查询次数: 481 次")
	})

	// ========================================
	// Scenario 4: 多层级不同 BatchSize 配置
	// 验证：不同层级使用不同的 BatchSize
	// ========================================
	gtest.C(t, func(t *gtest.T) {
		fmt.Println("\n========== Scenario 4: 多层级BatchSize配置（查询500用户）==========")
		fmt.Println("配置：Layer 1 BatchSize=100, Layer 2 BatchSize=200")

		db.SetDebug(true)
		fmt.Println("执行查询...")
		startTime := time.Now()
		var users []*User
		err := db.Model(tableUser).WithAll().
			WithBatchOption(
				gdb.WithBatchOption{Layer: 1, Enabled: true, BatchThreshold: 0, BatchSize: 100},
				gdb.WithBatchOption{Layer: 2, Enabled: true, BatchThreshold: 0, BatchSize: 200},
			).
			WithBatch().
			Where("id<=?", 500).
			Scan(&users)
		duration := time.Since(startTime)
		db.SetDebug(false)

		t.AssertNil(err)
		fmt.Printf("查询完成，耗时: %v\n", duration)

		fmt.Println("验证数据完整性（采样验证前50个用户）...")
		verifyUserData(users[:50], 50, true)

		fmt.Println("✓ 数据验证通过")
		fmt.Println("\n预期行为：")
		fmt.Println("  1. 主查询: 1 次")
		fmt.Println("  2. UserDetail: 5 次（500/100=5，Layer 1配置）")
		fmt.Println("  3. UserScores: 5 次（500/100=5，Layer 1配置）")
		fmt.Println("  4. ScoreDetails: 25 次（5000/200=25，Layer 2配置）")
		fmt.Println("  总查询次数: 36 次")
	})

	// ========================================
	// Scenario 5: 性能对比测试
	// 对比不同优化方案的性能
	// ========================================
	gtest.C(t, func(t *gtest.T) {
		fmt.Println("\n========== Scenario 5: 性能对比测试 ==========")

		testCount := 200 // 测试200个用户

		// 5.1 不开启优化
		fmt.Println("\n→ 测试 A: 不开启优化")
		startTime := time.Now()
		var usersA []*User
		err := db.Model(tableUser).WithAll().Where("id<=?", testCount).Scan(&usersA)
		durationA := time.Since(startTime)
		t.AssertNil(err)
		fmt.Printf("  耗时: %v\n", durationA)

		// 5.2 开启优化（默认配置）
		fmt.Println("\n→ 测试 B: 开启优化（默认BatchSize=1000）")
		startTime = time.Now()
		var usersB []*User
		err = db.Model(tableUser).WithAll().WithBatch().Where("id<=?", testCount).Scan(&usersB)
		durationB := time.Since(startTime)
		t.AssertNil(err)
		fmt.Printf("  耗时: %v\n", durationB)

		// 5.3 开启优化（BatchSize=50）
		fmt.Println("\n→ 测试 C: 开启优化（BatchSize=50）")
		startTime = time.Now()
		var usersC []*User
		err = db.Model(tableUser).WithAll().
			WithBatchOption(gdb.WithBatchOption{Layer: 0, Enabled: true, BatchThreshold: 0, BatchSize: 50}).
			WithBatch().
			Where("id<=?", testCount).
			Scan(&usersC)
		durationC := time.Since(startTime)
		t.AssertNil(err)
		fmt.Printf("  耗时: %v\n", durationC)

		// 验证数据一致性
		fmt.Println("\n验证三种方式的数据一致性...")
		verifyUserData(usersA, testCount, true)
		verifyUserData(usersB, testCount, true)
		verifyUserData(usersC, testCount, true)

		// 性能对比
		fmt.Println("\n========== 性能对比结果 ==========")
		fmt.Printf("不开启优化:              %v\n", durationA)
		fmt.Printf("优化(BatchSize=1000):    %v (提升 %.1fx)\n", durationB, float64(durationA)/float64(durationB))
		fmt.Printf("优化(BatchSize=50):      %v (提升 %.1fx)\n", durationC, float64(durationA)/float64(durationC))
		fmt.Println("================================")

		fmt.Println("✓ 所有数据验证通过")
	})

	fmt.Println("\n========== 大数据集测试全部完成 ==========")
}
