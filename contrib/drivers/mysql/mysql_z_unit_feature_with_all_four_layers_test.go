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

	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gmeta"
)

// Test_WithAll_MediumDataset_FourLayers 中小型四层数据分批查询测试
// 数据规模：
// - 20 个用户
// - 每个用户 5 个 UserScore（共 100 个 score）
// - 每个 score 4 个 ScoreDetail（共 400 个 detail）
// - 每个 detail 3 个 DetailMeta（共 1200 个 meta）
//
// 测试目标：
// 1. 验证四层关联查询的正确性
// 2. 验证每层都能正确触发分批查询
// 3. 输出较短的 SQL 便于验证查询逻辑
// 4. 对比默认查询和优化查询的差异
func Test_WithAll_MediumDataset_FourLayers(t *testing.T) {
	var (
		// 使用独立的表名避免与其他测试冲突
		tableUser             = "user_four_layers"
		tableUserDetail       = "user_detail_four_layers"
		tableUserScores       = "user_scores_four_layers"
		tableUserScoreDetails = "user_score_details_four_layers"
		tableDetailMeta       = "detail_meta_four_layers"
	)

	// 数据结构定义（四层）
	type DetailMeta struct {
		gmeta.Meta `orm:"table:detail_meta_four_layers"`
		Id         int    `json:"id"`
		DetailId   int    `json:"detail_id"`
		MetaKey    string `json:"meta_key"`
		MetaValue  string `json:"meta_value"`
	}

	type UserScoreDetails struct {
		gmeta.Meta `orm:"table:user_score_details_four_layers"`
		Id         int           `json:"id"`
		ScoreId    int           `json:"score_id"`
		DetailInfo string        `json:"detail_info"`
		DetailMeta []*DetailMeta `orm:"with:detail_id=id,batch:threshold=100,batchSize=200"`
	}

	type UserScores struct {
		gmeta.Meta   `orm:"table:user_scores_four_layers"`
		Id           int                 `json:"id"`
		Uid          int                 `json:"uid"`
		Score        int                 `json:"score"`
		ScoreDetails []*UserScoreDetails `orm:"with:score_id=id"`
	}

	type UserDetail struct {
		gmeta.Meta `orm:"table:user_detail_four_layers"`
		Uid        int    `json:"uid"`
		Address    string `json:"address"`
	}

	type User struct {
		gmeta.Meta `orm:"table:user_four_layers"`
		Id         int           `json:"id"`
		Name       string        `json:"name"`
		UserDetail *UserDetail   `orm:"with:uid=id, batch:threshold=5,batchSize=5"`
		UserScores []*UserScores `orm:"with:uid=id,batch:threshold=10,batchSize=10"`
	}

	// 初始化表结构
	dropTable(tableUser)
	dropTable(tableUserDetail)
	dropTable(tableUserScores)
	dropTable(tableUserScoreDetails)
	dropTable(tableDetailMeta)

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

	_, err = db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id int(10) unsigned NOT NULL AUTO_INCREMENT,
			detail_id int(10) unsigned NOT NULL,
			meta_key varchar(50) NOT NULL,
			meta_value varchar(100) NOT NULL,
			PRIMARY KEY (id),
			KEY idx_detail_id (detail_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`, tableDetailMeta))
	gtest.AssertNil(err)

	defer dropTable(tableUser)
	defer dropTable(tableUserDetail)
	defer dropTable(tableUserScores)
	defer dropTable(tableUserScoreDetails)
	defer dropTable(tableDetailMeta)

	// ========================================
	// 数据初始化（中小数据量）
	// ========================================
	const (
		userCount      = 20                           // 20个用户
		scorePerUser   = 5                            // 每个用户5个score
		detailPerScore = 4                            // 每个score 4个detail
		metaPerDetail  = 3                            // 每个detail 3个meta
		totalScores    = userCount * scorePerUser     // 100
		totalDetails   = totalScores * detailPerScore // 400
		totalMeta      = totalDetails * metaPerDetail // 1200
	)

	fmt.Println("\n========== 开始初始化中小型四层数据集 ==========")
	fmt.Printf("数据规模：%d 用户, %d scores, %d details, %d meta\n", userCount, totalScores, totalDetails, totalMeta)

	// 1. 插入用户数据
	fmt.Println("→ 插入用户数据...")
	startTime := time.Now()
	usersData := make([]*User, 0, userCount)
	for i := 1; i <= userCount; i++ {
		usersData = append(usersData, &User{
			Id:   i,
			Name: fmt.Sprintf("user_%d", i),
		})
	}
	_, err = db.Model(tableUser).Data(usersData).Insert()
	gtest.AssertNil(err)
	fmt.Printf("  用户数据插入完成，耗时: %v\n", time.Since(startTime))

	// 2. 插入用户详情
	fmt.Println("→ 插入用户详情...")
	startTime = time.Now()
	detailsData := make([]*UserDetail, 0, userCount)
	for i := 1; i <= userCount; i++ {
		detailsData = append(detailsData, &UserDetail{
			Uid:     i,
			Address: fmt.Sprintf("address_%d", i),
		})
	}
	_, err = db.Model(tableUserDetail).Data(detailsData).Insert()
	gtest.AssertNil(err)
	fmt.Printf("  用户详情插入完成，耗时: %v\n", time.Since(startTime))

	// 3. 插入 UserScores
	fmt.Println("→ 插入 UserScores...")
	startTime = time.Now()
	scoresData := make([]*UserScores, 0, totalScores)
	scoreId := 1
	for i := 1; i <= userCount; i++ {
		for j := 1; j <= scorePerUser; j++ {
			scoresData = append(scoresData, &UserScores{
				Id:    scoreId,
				Uid:   i,
				Score: j * 10,
			})
			scoreId++
		}
	}
	_, err = db.Model(tableUserScores).Data(scoresData).Insert()
	gtest.AssertNil(err)
	fmt.Printf("  UserScores 插入完成，耗时: %v\n", time.Since(startTime))

	// 4. 插入 ScoreDetails
	fmt.Println("→ 插入 ScoreDetails...")
	startTime = time.Now()
	scoreDetailsData := make([]*UserScoreDetails, 0, totalDetails)
	detailId := 1
	for i := 1; i <= totalScores; i++ {
		for j := 1; j <= detailPerScore; j++ {
			scoreDetailsData = append(scoreDetailsData, &UserScoreDetails{
				Id:         detailId,
				ScoreId:    i,
				DetailInfo: fmt.Sprintf("detail_%d_%d", i, j),
			})
			detailId++
		}
	}
	_, err = db.Model(tableUserScoreDetails).Data(scoreDetailsData).Insert()
	gtest.AssertNil(err)
	fmt.Printf("  ScoreDetails 插入完成，耗时: %v\n", time.Since(startTime))

	// 5. 插入 DetailMeta（第四层）
	fmt.Println("→ 插入 DetailMeta...")
	startTime = time.Now()
	metaData := make([]*DetailMeta, 0, totalMeta)
	for i := 1; i <= totalDetails; i++ {
		for j := 1; j <= metaPerDetail; j++ {
			metaData = append(metaData, &DetailMeta{
				DetailId:  i,
				MetaKey:   fmt.Sprintf("key_%d", j),
				MetaValue: fmt.Sprintf("value_%d_%d", i, j),
			})
		}
	}
	_, err = db.Model(tableDetailMeta).Data(metaData).Insert()
	gtest.AssertNil(err)
	fmt.Printf("  DetailMeta 插入完成，耗时: %v\n", time.Since(startTime))

	fmt.Println("========== 数据初始化完成 ==========")

	// 数据验证辅助函数
	verifyUserData := func(users []*User, expectedCount int, scenario string) {
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
				gtest.Assert(len(score.ScoreDetails), detailPerScore)
				for detailIdx, detail := range score.ScoreDetails {
					gtest.Assert(detail.ScoreId, score.Id)
					expectedInfo := fmt.Sprintf("detail_%d_%d", score.Id, detailIdx+1)
					gtest.Assert(detail.DetailInfo, expectedInfo)

					// 验证 DetailMeta（第四层）
					gtest.Assert(len(detail.DetailMeta), metaPerDetail)
					for metaIdx, meta := range detail.DetailMeta {
						gtest.Assert(meta.DetailId, detail.Id)
						gtest.Assert(meta.MetaKey, fmt.Sprintf("key_%d", metaIdx+1))
						expectedValue := fmt.Sprintf("value_%d_%d", detail.Id, metaIdx+1)
						gtest.Assert(meta.MetaValue, expectedValue)
					}
				}
			}
		}
		fmt.Printf("✓ %s - 数据验证通过（验证了 %d 个用户的四层完整数据）\n", scenario, expectedCount)
	}

	gtest.C(t, func(t *gtest.T) {

		db.SetDebug(true)
		fmt.Println("\n开始执行查询...")
		startTime := time.Now()
		var users []*User
		err := db.Model(tableUser).WithAll().Scan(&users)
		duration := time.Since(startTime)
		db.SetDebug(false)

		t.AssertNil(err)
		fmt.Printf("\n查询完成,耗时: %v\n", duration)

		verifyUserData(users, 20, "Scenario 1")
	})

	gtest.C(t, func(t *gtest.T) {

		db.SetDebug(true)
		fmt.Println("\n开始执行查询...")
		startTime := time.Now()
		var users []*User
		err := db.Model(tableUser).WithBatch().WithAll().Scan(&users)
		duration := time.Since(startTime)
		db.SetDebug(false)

		t.AssertNil(err)
		fmt.Printf("\n查询完成,耗时: %v\n", duration)

		verifyUserData(users, 20, "Scenario 2")
	})
	fmt.Println("\n========== 中小型数据集测试全部完成 ==========")
}
