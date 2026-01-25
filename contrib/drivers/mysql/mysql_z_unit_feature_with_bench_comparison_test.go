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

func Test_WithAll_PerformanceComparison(t *testing.T) {
	var (
		tableUser             = "user_bench_disabled"
		tableUserDetail       = "user_detail_bench_disabled"
		tableUserScores       = "user_scores_bench_disabled"
		tableUserScoreDetails = "user_score_details_bench_disabled"
	)

	type UserScoreDetails struct {
		gmeta.Meta `orm:"table:user_score_details_bench_disabled"`
		ScoreId    int    `json:"score_id"`
		DetailInfo string `json:"detail_info"`
	}

	type UserScores struct {
		gmeta.Meta   `orm:"table:user_scores_bench_disabled"`
		Id           int                 `json:"id"`
		Uid          int                 `json:"uid"`
		Score        int                 `json:"score"`
		ScoreDetails []*UserScoreDetails `orm:"with:score_id=id"`
	}

	type UserDetail struct {
		gmeta.Meta `orm:"table:user_detail_bench_disabled"`
		Uid        int    `json:"uid"`
		Address    string `json:"address"`
	}

	type User struct {
		gmeta.Meta `orm:"table:user_bench_disabled"`
		Id         int           `json:"id"`
		Name       string        `json:"name"`
		UserDetail *UserDetail   `orm:"with:uid=id"`
		UserScores []*UserScores `orm:"with:uid=id"`
	}

	// 1. Initialize Tables and Data
	dropTable(tableUser)
	dropTable(tableUserDetail)
	dropTable(tableUserScores)
	dropTable(tableUserScoreDetails)

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

	defer dropTable(tableUser)
	defer dropTable(tableUserDetail)
	defer dropTable(tableUserScores)
	defer dropTable(tableUserScoreDetails)

	// Insert Data
	fmt.Println("Initializing data...")
	userCount := 500
	detailPerUser := 1
	scorePerUser := 10
	scoreDetailPerScore := 10

	// Batch insert users
	usersData := g.List{}
	for i := 1; i <= userCount; i++ {
		usersData = append(usersData, g.Map{"id": i, "name": fmt.Sprintf("name_%d", i)})
	}
	_, err = db.Model(tableUser).Data(usersData).Insert()
	gtest.AssertNil(err)

	// Batch insert details
	detailsData := g.List{}
	for i := 1; i <= userCount; i++ {
		for j := 1; j <= detailPerUser; j++ {
			detailsData = append(detailsData, g.Map{"uid": i, "address": fmt.Sprintf("address_%d_%d", i, j)})
		}
	}
	_, err = db.Model(tableUserDetail).Data(detailsData).Insert()
	gtest.AssertNil(err)

	// Batch insert scores
	scoresData := g.List{}
	for i := 1; i <= userCount; i++ {
		for j := 1; j <= scorePerUser; j++ {
			scoresData = append(scoresData, g.Map{"uid": i, "score": j})
		}
	}
	_, err = db.Model(tableUserScores).Data(scoresData).Insert()
	gtest.AssertNil(err)

	// Batch insert score details
	scoreDetailsData := g.List{}
	for i := 1; i <= userCount; i++ {
		for j := 1; j <= scorePerUser; j++ {
			actualScoreId := (i-1)*scorePerUser + j
			for k := 1; k <= scoreDetailPerScore; k++ {
				scoreDetailsData = append(scoreDetailsData, g.Map{
					"score_id":    actualScoreId,
					"detail_info": fmt.Sprintf("detail_info_%d_%d_%d", i, j, k),
				})
			}
		}
	}

	_, err = db.Model(tableUserScoreDetails).Data(scoreDetailsData).Batch(1000).Insert()
	gtest.AssertNil(err)

	fmt.Println("Data initialization completed.")

	// Prepare Query Condition
	var userIds []int
	for i := 1; i <= userCount; i++ {
		userIds = append(userIds, i)
	}

	// Helper for verification
	verifyData := func(t *gtest.T, users []*User, expectedUserCount ...int) {
		count := userCount
		if len(expectedUserCount) > 0 {
			count = expectedUserCount[0]
		}
		t.Assert(len(users), count)
		for _, u := range users {
			t.AssertNE(u.UserDetail, nil)
			t.Assert(u.UserDetail.Uid, u.Id)
			t.Assert(u.UserDetail.Address == fmt.Sprintf("address_%d_1", u.Id) || u.UserDetail.Address == fmt.Sprintf("address_%d_2", u.Id), true)
			t.Assert(len(u.UserScores), scorePerUser)
			for _, s := range u.UserScores {
				t.Assert(s.Uid, u.Id)
				t.Assert(len(s.ScoreDetails), scoreDetailPerScore)
				if len(s.ScoreDetails) > 0 {
					t.Assert(s.ScoreDetails[0].ScoreId, s.Id)
				}
			}
		}
		fmt.Printf("Scenario: Data verification passed.\n")
	}

	db.SetDebug(true)
	defer db.SetDebug(false)

	gtest.C(t, func(t *gtest.T) {
		var (
			users []*User
			start = time.Now()
		)
		// We use small BatchThreshold and BatchSize to force many batches for 500 users
		err := db.Model(tableUser).
			WithAll().
			WithBatch().
			Where("id", userIds).
			Scan(&users)
		duration := time.Since(start)
		t.AssertNil(err)
		verifyData(t, users)
		fmt.Printf("Total Time: %v\n", duration)
	})

	gtest.C(t, func(t *gtest.T) {
		var users []*User
		// Scenario 2: BatchThreshold.
		// Set BatchThreshold to a large value, so batch optimization should NOT be triggered.
		// We can verify this by checking the logs (should see many SELECT statements instead of one big IN).
		fmt.Println("\nTesting BatchThreshold (Optimization should be DISABLED)...")
		err := db.Model(tableUser).
			WithAll().
			WithBatchOption(gdb.WithBatchOption{
				Layer:          2, // UserScores -> UserScoreDetails is Layer 2
				Enabled:        true,
				BatchThreshold: 100000, // Very high threshold
			}).
			WithBatch().
			Where("id", userIds[:5]). // Only 5 users
			Scan(&users)
		t.AssertNil(err)
		verifyData(t, users, 5)

		// Scenario 3: Disable batching for a specific layer.
		fmt.Println("\nTesting Layer-specific Disable (Layer 1 should be DISABLED)...")
		users = nil
		err = db.Model(tableUser).
			WithAll().
			WithBatchOption(gdb.WithBatchOption{
				Layer:   1, // User -> UserDetail/UserScores is Layer 1
				Enabled: false,
			}).
			WithBatch().
			Where("id", userIds[:5]).
			Scan(&users)
		t.AssertNil(err)
		verifyData(t, users, 5)
	})
}
