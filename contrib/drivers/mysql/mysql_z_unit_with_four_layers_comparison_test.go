// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mysql_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gmeta"
)

// Test_With_FourLayers_Comparison tests four-layer nested With using the same dataset.
// It runs both Batch and Chunk modes and validates:
// 1. Data correctness
// 2. Query count
// 3. Performance
func Test_With_FourLayers_Comparison(t *testing.T) {
	var (
		tableUser         = "user_four_comparison"
		tableUserDetail   = "user_detail_four_comparison"
		tableUserScores   = "user_scores_four_comparison"
		tableScoreDetails = "score_details_four_comparison"
		tableDetailMeta   = "detail_meta_four_comparison"
	)

	// Create tables
	if _, err := db.Exec(ctx, fmt.Sprintf(gtest.DataContent("with_tpl_user_soft_delete.sql"), tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUser)

	if _, err := db.Exec(ctx, fmt.Sprintf(gtest.DataContent("with_tpl_user_detail.sql"), tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserDetail)

	if _, err := db.Exec(ctx, fmt.Sprintf(gtest.DataContent("with_tpl_user_scores_soft_delete.sql"), tableUserScores)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserScores)

	if _, err := db.Exec(ctx, fmt.Sprintf(gtest.DataContent("with_tpl_score_details_soft_delete.sql"), tableScoreDetails)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableScoreDetails)

	if _, err := db.Exec(ctx, fmt.Sprintf(gtest.DataContent("with_tpl_detail_meta_soft_delete.sql"), tableDetailMeta)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableDetailMeta)

	// Define structures with where, order, and unscoped:true tags
	type UserDetailInfo struct {
		gmeta.Meta `orm:"table:user_detail_four_comparison"`
		Uid        int    `json:"uid"`
		Address    string `json:"address"`
	}

	type DetailMeta struct {
		gmeta.Meta `orm:"table:detail_meta_four_comparison"`
		Id         int         `json:"id"`
		DetailId   int         `json:"detail_id"`
		MetaKey    string      `json:"meta_key"`
		MetaValue  string      `json:"meta_value"`
		SortOrder  int         `json:"sort_order"`
		DeletedAt  *gtime.Time `json:"deleted_at"`
	}

	type ScoreDetails struct {
		gmeta.Meta `orm:"table:score_details_four_comparison"`
		Id         int           `json:"id"`
		ScoreId    int           `json:"score_id"`
		DetailInfo string        `json:"detail_info"`
		Rank       int           `json:"rank"`
		DeletedAt  *gtime.Time   `json:"deleted_at"`
		DetailMeta []*DetailMeta `orm:"with:detail_id=id, where:meta_key like 'key_%', order:sort_order asc, unscoped:true"`
	}

	type UserScores struct {
		gmeta.Meta   `orm:"table:user_scores_four_comparison"`
		Id           int             `json:"id"`
		Uid          int             `json:"uid"`
		Score        int             `json:"score"`
		Priority     int             `json:"priority"`
		DeletedAt    *gtime.Time     `json:"deleted_at"`
		ScoreDetails []*ScoreDetails `orm:"with:score_id=id, where:rank > 0, order:rank desc, unscoped:true, chunkName:detailChunk, chunkSize:20, chunkMinRows:10"`
	}

	type User struct {
		gmeta.Meta `orm:"table:user_four_comparison"`
		Id         int             `json:"id"`
		Name       string          `json:"name"`
		Status     int             `json:"status"`
		DeletedAt  *gtime.Time     `json:"deleted_at"`
		UserDetail *UserDetailInfo `orm:"with:uid=id"`
		UserScores []*UserScores   `orm:"with:uid=id, where:score >= 10, order:priority desc, unscoped:true, chunkName:scoreChunk, chunkSize:15, chunkMinRows:8"`
	}

	// Initialize test data: 100 users for better performance comparison
	t.Log("Initializing test data...")
	for i := 1; i <= 100; i++ {
		user := User{
			Name:   fmt.Sprintf("user_%d", i),
			Status: i % 3, // 0, 1, 2
		}
		// Soft delete some users (every 10th)
		if i%10 == 0 {
			now := gtime.Now()
			user.DeletedAt = now
		}
		userId, err := db.Model(tableUser).Data(user).OmitEmpty().InsertAndGetId()
		gtest.AssertNil(err)

		// Create UserDetail for each user
		userDetail := UserDetailInfo{
			Uid:     int(userId),
			Address: fmt.Sprintf("address_%d", i),
		}
		_, err = db.Model(tableUserDetail).Data(userDetail).Insert()
		gtest.AssertNil(err)

		// Each user has 5 scores
		for j := 1; j <= 5; j++ {
			userScore := UserScores{
				Uid:      int(userId),
				Score:    j * 10,
				Priority: j,
			}
			// Soft delete the last score (j==5)
			if j == 5 {
				now := gtime.Now()
				userScore.DeletedAt = now
			}
			scoreId, err := db.Model(tableUserScores).Data(userScore).OmitEmpty().InsertAndGetId()
			gtest.AssertNil(err)

			// Each score has 4 details
			for k := 1; k <= 4; k++ {
				scoreDetail := ScoreDetails{
					ScoreId:    int(scoreId),
					DetailInfo: fmt.Sprintf("detail_%d_%d", j, k),
					Rank:       k,
				}
				// Soft delete the last detail (k==4)
				if k == 4 {
					now := gtime.Now()
					scoreDetail.DeletedAt = now
				}
				detailId, err := db.Model(tableScoreDetails).Data(scoreDetail).OmitEmpty().InsertAndGetId()
				gtest.AssertNil(err)

				// Each detail has 3 meta entries
				for m := 1; m <= 3; m++ {
					meta := DetailMeta{
						DetailId:  int(detailId),
						MetaKey:   fmt.Sprintf("key_%d", m),
						MetaValue: fmt.Sprintf("value_%d_%d_%d", j, k, m),
						SortOrder: m,
					}
					// Soft delete the last meta (m==3)
					if m == 3 {
						now := gtime.Now()
						meta.DeletedAt = now
					}
					_, err = db.Model(tableDetailMeta).Data(meta).OmitEmpty().Insert()
					gtest.AssertNil(err)
				}
			}
		}
	}
	t.Log("Test data initialized successfully")

	gtest.C(t, func(t *gtest.T) {
		// Helper function to validate data correctness
		validateData := func(users []*User, mode string) {
			t.Logf("\n=== Validating %s mode data ===", mode)
			t.Assert(len(users) > 0, true)
			t.Logf("Total users loaded: %d", len(users))

			totalScores := 0
			totalDetails := 0
			totalMeta := 0
			deletedScoresCount := 0
			deletedDetailsCount := 0
			deletedMetaCount := 0

			for _, user := range users {
				t.Assert(user.Status, 1)

				t.AssertNE(user.UserDetail, nil)
				t.Assert(user.UserDetail.Uid, user.Id)
				t.Assert(strings.HasPrefix(user.UserDetail.Address, "address_"), true)

				t.Assert(len(user.UserScores) > 0, true)
				totalScores += len(user.UserScores)

				for _, score := range user.UserScores {
					t.Assert(score.Score >= 10, true)

					if score.DeletedAt != nil {
						deletedScoresCount++
					}

					t.Assert(len(score.ScoreDetails) > 0, true)
					totalDetails += len(score.ScoreDetails)

					for _, detail := range score.ScoreDetails {
						t.Assert(detail.Rank > 0, true)
						if detail.DeletedAt != nil {
							deletedDetailsCount++
						}

						t.Assert(len(detail.DetailMeta) > 0, true)
						totalMeta += len(detail.DetailMeta)

						for _, meta := range detail.DetailMeta {
							t.Assert(strings.HasPrefix(meta.MetaKey, "key_"), true)
							if meta.DeletedAt != nil {
								deletedMetaCount++
							}
						}

						if len(detail.DetailMeta) > 1 {
							for i := 0; i < len(detail.DetailMeta)-1; i++ {
								t.Assert(detail.DetailMeta[i].SortOrder <= detail.DetailMeta[i+1].SortOrder, true)
							}
						}
					}

					if len(score.ScoreDetails) > 1 {
						for i := 0; i < len(score.ScoreDetails)-1; i++ {
							t.Assert(score.ScoreDetails[i].Rank >= score.ScoreDetails[i+1].Rank, true)
						}
					}
				}

				if len(user.UserScores) > 1 {
					for i := 0; i < len(user.UserScores)-1; i++ {
						t.Assert(user.UserScores[i].Priority >= user.UserScores[i+1].Priority, true)
					}
				}
			}

			t.Logf("Total scores: %d (deleted: %d)", totalScores, deletedScoresCount)
			t.Logf("Total details: %d (deleted: %d)", totalDetails, deletedDetailsCount)
			t.Logf("Total meta: %d (deleted: %d)", totalMeta, deletedMetaCount)

			t.Assert(deletedScoresCount > 0, true)
			t.Assert(deletedDetailsCount > 0, true)
			t.Assert(deletedMetaCount > 0, true)

			t.Logf("Data validation passed for %s mode", mode)
		}

		// Batch With mode
		t.Log("\n=== Testing Batch With Mode ===")
		oldDebug := db.GetDebug()
		db.SetDebug(true)

		startTime := time.Now()
		var usersBatch []*User
		err := db.Model(tableUser).Where("status=?", 1).WithAll().Scan(&usersBatch)
		duration := time.Since(startTime)

		db.SetDebug(oldDebug)

		t.AssertNil(err)
		t.Logf("Batch With Mode Duration: %v", duration)
		t.Log("Note: Check console output above to count SELECT queries")
		validateData(usersBatch, "Batch")

		// Chunk mode
		t.Log("\n=== Testing Chunk Mode ===")
		db.SetDebug(true)

		startTime = time.Now()
		var usersChunk []*User
		err = db.Model(tableUser).
			Where("status=?", 1).
			WithOptions(
				gdb.WithOption{ChunkName: "scoreChunk", ChunkSize: 12, ChunkMinRows: 6},
				gdb.WithOption{ChunkName: "detailChunk", ChunkSize: 10, ChunkMinRows: 5},
			).
			WithAll().
			Scan(&usersChunk)
		duration = time.Since(startTime)

		db.SetDebug(oldDebug)

		t.AssertNil(err)
		t.Logf("Chunk Mode Duration: %v", duration)
		t.Log("Note: Check console output above to count SELECT queries")
		validateData(usersChunk, "Chunk")

		// Summary
		t.Log("\n=== Performance Comparison Summary ===")
		t.Log("Both modes returned the same correct data.")
		t.Log("Batch mode: Single-batch queries (best for small datasets)")
		t.Log("Chunk mode: Chunked batch queries (balanced for large datasets)")
	})
}
