// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

func createRangePartitionTable(table ...string) string {
	var name string
	if len(table) > 0 {
		name = table[0]
	} else {
		name = fmt.Sprintf(`partition_range_%d`, gtime.TimestampNano())
	}
	dropTable(name)
	// Note: DM does not partition by an expression like YEAR(sales_date), so the ranges are
	// bounds of the DATE column itself.
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id INT NOT NULL,
			sales_date DATE DEFAULT NULL,
			amount DECIMAL(10,2) DEFAULT NULL,
			region VARCHAR(50) DEFAULT NULL
		)
		PARTITION BY RANGE (sales_date)
		(PARTITION p2020 VALUES LESS THAN ('2021-01-01'),
		 PARTITION p2021 VALUES LESS THAN ('2022-01-01'),
		 PARTITION p2022 VALUES LESS THAN ('2023-01-01'),
		 PARTITION p2023 VALUES LESS THAN ('2024-01-01'),
		 PARTITION p_future VALUES LESS THAN (MAXVALUE))
	`, name)); err != nil {
		gtest.Fatal(err)
	}
	return name
}

func createHashPartitionTable(table ...string) string {
	var name string
	if len(table) > 0 {
		name = table[0]
	} else {
		name = fmt.Sprintf(`partition_hash_%d`, gtime.TimestampNano())
	}
	dropTable(name)
	if _, err := db.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id INT NOT NULL,
			user_id INT NOT NULL,
			username VARCHAR(50) DEFAULT NULL,
			email VARCHAR(100) DEFAULT NULL
		)
		PARTITION BY HASH (user_id)
		PARTITIONS 4
	`, name)); err != nil {
		gtest.Fatal(err)
	}
	return name
}

func dropPartitionTable(table string) {
	dropTable(table)
}

func Test_Partition_Range_Insert_And_Query(t *testing.T) {
	table := createRangePartitionTable()
	defer dropPartitionTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert data across different partitions
		data := g.Slice{
			g.Map{"id": 1, "sales_date": "2020-06-15", "amount": 1000.50, "region": "North"},
			g.Map{"id": 2, "sales_date": "2021-03-20", "amount": 2000.75, "region": "South"},
			g.Map{"id": 3, "sales_date": "2022-09-10", "amount": 3000.00, "region": "East"},
			g.Map{"id": 4, "sales_date": "2023-12-01", "amount": 4000.25, "region": "West"},
			g.Map{"id": 5, "sales_date": "2024-01-15", "amount": 5000.00, "region": "North"},
		}
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)

		// Query all data
		all, err := db.Model(table).All()
		t.AssertNil(err)
		t.Assert(len(all), 5)

		// Query specific year (should hit specific partition)
		result, err := db.Model(table).Where("YEAR(sales_date) = ?", 2022).All()
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["ID"], 3)
	})
}

func Test_Partition_Hash_Insert_And_Distribution(t *testing.T) {
	table := createHashPartitionTable()
	defer dropPartitionTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert data that will be distributed across hash partitions
		data := g.Slice{}
		for i := 1; i <= 20; i++ {
			data = append(data, g.Map{
				"id":       i,
				"user_id":  i * 10,
				"username": fmt.Sprintf("user_%d", i),
				"email":    fmt.Sprintf("user%d@example.com", i),
			})
		}
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)

		// Query all data
		all, err := db.Model(table).All()
		t.AssertNil(err)
		t.Assert(len(all), 20)

		// Query specific user_id (will hit specific partition based on hash)
		result, err := db.Model(table).Where("user_id", 100).One()
		t.AssertNil(err)
		t.Assert(result["USERNAME"], "user_10")
	})
}

func Test_Partition_Range_Update(t *testing.T) {
	table := createRangePartitionTable()
	defer dropPartitionTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert data
		_, err := db.Model(table).Data(g.Map{
			"id":         1,
			"sales_date": "2022-06-15",
			"amount":     1000.00,
			"region":     "North",
		}).Insert()
		t.AssertNil(err)

		// Update data within same partition
		result, err := db.Model(table).Data(g.Map{
			"amount": 1500.00,
			"region": "South",
		}).Where("id", 1).Update()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		// Verify update
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["AMOUNT"], "1500.00")
		t.Assert(one["REGION"], "South")
	})
}

func Test_Partition_Range_Delete(t *testing.T) {
	table := createRangePartitionTable()
	defer dropPartitionTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert data
		data := g.Slice{
			g.Map{"id": 1, "sales_date": "2020-06-15", "amount": 1000.50},
			g.Map{"id": 2, "sales_date": "2021-03-20", "amount": 2000.75},
			g.Map{"id": 3, "sales_date": "2022-09-10", "amount": 3000.00},
		}
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)

		// Delete from specific partition
		result, err := db.Model(table).Where("YEAR(sales_date) = ?", 2021).Delete()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		// Verify deletion
		all, err := db.Model(table).All()
		t.AssertNil(err)
		t.Assert(len(all), 2)

		// Verify remaining data
		result2, err := db.Model(table).Where("YEAR(sales_date) = ?", 2021).All()
		t.AssertNil(err)
		t.Assert(len(result2), 0)
	})
}

func Test_Partition_Transaction(t *testing.T) {
	table := createRangePartitionTable()
	defer dropPartitionTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Transaction with partitioned table
		err := db.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			// Insert across multiple partitions
			data := g.Slice{
				g.Map{"id": 1, "sales_date": "2020-06-15", "amount": 1000.50},
				g.Map{"id": 2, "sales_date": "2021-03-20", "amount": 2000.75},
				g.Map{"id": 3, "sales_date": "2022-09-10", "amount": 3000.00},
			}
			_, err := tx.Model(table).Ctx(ctx).Data(data).Insert()
			if err != nil {
				return err
			}

			// Update in transaction
			_, err = tx.Model(table).Ctx(ctx).Data(g.Map{
				"amount": 1500.00,
			}).Where("id", 1).Update()
			return err
		})
		t.AssertNil(err)

		// Verify transaction committed
		all, err := db.Model(table).All()
		t.AssertNil(err)
		t.Assert(len(all), 3)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["AMOUNT"], "1500.00")
	})
}

func Test_Partition_Range_Count_And_Sum(t *testing.T) {
	table := createRangePartitionTable()
	defer dropPartitionTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert data
		data := g.Slice{
			g.Map{"id": 1, "sales_date": "2020-06-15", "amount": 1000.00},
			g.Map{"id": 2, "sales_date": "2020-09-20", "amount": 1500.00},
			g.Map{"id": 3, "sales_date": "2021-03-20", "amount": 2000.00},
			g.Map{"id": 4, "sales_date": "2022-09-10", "amount": 3000.00},
		}
		_, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)

		// Count by year (specific partition)
		count, err := db.Model(table).Where("YEAR(sales_date) = ?", 2020).Count()
		t.AssertNil(err)
		t.Assert(count, 2)

		// Sum across partitions
		value, err := db.Model(table).Fields("SUM(amount) as total").Value()
		t.AssertNil(err)
		t.AssertGT(value.Float64(), 7000.0) // 1000+1500+2000+3000 = 7500
	})
}
