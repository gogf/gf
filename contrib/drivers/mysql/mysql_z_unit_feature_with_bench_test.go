// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mysql_test

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gmeta"
)

func createBenchTable(b *testing.B, template, table string) {
	dropTable(table)
	if _, err := db.Exec(ctx, fmt.Sprintf(gtest.DataContent(template), table)); err != nil {
		b.Fatal(err)
	}
	b.Cleanup(func() {
		dropTable(table)
	})
}

func batchInsertBenchData(b *testing.B, table string, data []g.Map) {
	if _, err := db.Ctx(ctx).Model(table).Data(data).Batch(1000).Insert(); err != nil {
		b.Fatal(err)
	}
}

// Benchmark_Table_Relation_WithAll_List measures nested association scanning
// with 100 parent records, where nested queries are the dominant cost.
func Benchmark_Table_Relation_WithAll_List(b *testing.B) {
	b.Run("OneToOneThenSlice", func(b *testing.B) {
		const (
			tableUser       = "bench_with_user"
			tableUserDetail = "bench_with_user_detail"
			tableUserScores = "bench_with_user_scores"
			nUsers          = 100
			nScores         = 5
		)
		createBenchTable(b, "with_tpl_user.sql", tableUser)
		createBenchTable(b, "with_tpl_user_detail.sql", tableUserDetail)
		createBenchTable(b, "with_tpl_user_scores.sql", tableUserScores)

		var (
			users   = make([]g.Map, 0, nUsers)
			details = make([]g.Map, 0, nUsers)
			scores  = make([]g.Map, 0, nUsers*nScores)
			scoreId int
		)
		for i := 1; i <= nUsers; i++ {
			users = append(users, g.Map{
				"id":   i,
				"name": fmt.Sprintf(`name_%d`, i),
			})
			details = append(details, g.Map{
				"uid":     i,
				"address": fmt.Sprintf(`address_%d`, i),
			})
			for j := 1; j <= nScores; j++ {
				scoreId++
				scores = append(scores, g.Map{
					"id":    scoreId,
					"uid":   i,
					"score": j,
				})
			}
		}
		batchInsertBenchData(b, tableUser, users)
		batchInsertBenchData(b, tableUserDetail, details)
		batchInsertBenchData(b, tableUserScores, scores)

		type UserScores struct {
			gmeta.Meta `orm:"table:bench_with_user_scores"`
			Id         int `json:"id"`
			Uid        int `json:"uid"`
			Score      int `json:"score"`
		}
		type UserDetail struct {
			gmeta.Meta `orm:"table:bench_with_user_detail"`
			Uid        int           `json:"uid"`
			Address    string        `json:"address"`
			UserScores []*UserScores `orm:"with:uid"`
		}
		type User struct {
			gmeta.Meta `orm:"table:bench_with_user"`
			Id         int         `json:"id"`
			Name       string      `json:"name"`
			UserDetail *UserDetail `orm:"with:uid=id"`
		}

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var users []*User
			if err := db.Model(tableUser).WithAll().Scan(&users); err != nil {
				b.Fatal(err)
			}
			if len(users) != nUsers {
				b.Fatalf("expect %d users, got %d", nUsers, len(users))
			}
		}
	})

	b.Run("FourLayers", func(b *testing.B) {
		const (
			tableUser         = "bench_with_user"
			tableUserScores   = "bench_with_user_scores"
			tableScoreDetails = "bench_with_score_details"
			tableDetailMeta   = "bench_with_detail_meta"
			nUsers            = 100
			nScores           = 5
			nDetails          = 4
			nMetas            = 3
		)
		createBenchTable(b, "with_tpl_user.sql", tableUser)
		createBenchTable(b, "with_tpl_user_scores.sql", tableUserScores)
		createBenchTable(b, "with_tpl_score_details.sql", tableScoreDetails)
		createBenchTable(b, "with_tpl_detail_meta.sql", tableDetailMeta)

		var (
			users    = make([]g.Map, 0, nUsers)
			scores   = make([]g.Map, 0, nUsers*nScores)
			details  = make([]g.Map, 0, nUsers*nScores*nDetails)
			metas    = make([]g.Map, 0, nUsers*nScores*nDetails*nMetas)
			scoreId  int
			detailId int
		)
		for i := 1; i <= nUsers; i++ {
			users = append(users, g.Map{
				"id":   i,
				"name": fmt.Sprintf(`name_%d`, i),
			})
			for j := 1; j <= nScores; j++ {
				scoreId++
				scores = append(scores, g.Map{
					"id":    scoreId,
					"uid":   i,
					"score": j,
				})
				for k := 1; k <= nDetails; k++ {
					detailId++
					details = append(details, g.Map{
						"id":          detailId,
						"score_id":    scoreId,
						"detail_info": fmt.Sprintf(`detail_%d_%d_%d`, i, j, k),
					})
					for m := 1; m <= nMetas; m++ {
						metas = append(metas, g.Map{
							"detail_id": detailId,
							"meta_key":  fmt.Sprintf(`meta_%d_%d_%d_%d`, i, j, k, m),
						})
					}
				}
			}
		}
		batchInsertBenchData(b, tableUser, users)
		batchInsertBenchData(b, tableUserScores, scores)
		batchInsertBenchData(b, tableScoreDetails, details)
		batchInsertBenchData(b, tableDetailMeta, metas)

		type DetailMeta struct {
			gmeta.Meta `orm:"table:bench_with_detail_meta"`
			Id         int    `json:"id"`
			DetailId   int    `json:"detail_id"`
			MetaKey    string `json:"meta_key"`
		}
		type ScoreDetails struct {
			gmeta.Meta `orm:"table:bench_with_score_details"`
			Id         int           `json:"id"`
			ScoreId    int           `json:"score_id"`
			DetailInfo string        `json:"detail_info"`
			DetailMeta []*DetailMeta `orm:"with:detail_id=id"`
		}
		type UserScores struct {
			gmeta.Meta   `orm:"table:bench_with_user_scores"`
			Id           int             `json:"id"`
			Uid          int             `json:"uid"`
			Score        int             `json:"score"`
			ScoreDetails []*ScoreDetails `orm:"with:score_id=id"`
		}
		type User struct {
			gmeta.Meta `orm:"table:bench_with_user"`
			Id         int           `json:"id"`
			Name       string        `json:"name"`
			UserScores []*UserScores `orm:"with:uid=id"`
		}

		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var users []*User
			if err := db.Model(tableUser).WithAll().Scan(&users); err != nil {
				b.Fatal(err)
			}
			if len(users) != nUsers {
				b.Fatalf("expect %d users, got %d", nUsers, len(users))
			}
		}
	})
}
