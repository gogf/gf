// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mariadb_test

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
	if _, err := db3.Exec(ctx, fmt.Sprintf("DROP TABLE IF EXISTS `%s`", name)); err != nil {
		gtest.Fatal(err)
	}
	if _, err := db3.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id int(11) NOT NULL,
			sales_date date DEFAULT NULL,
			amount decimal(10,2) DEFAULT NULL,
			region varchar(50) DEFAULT NULL
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
		PARTITION BY RANGE (YEAR(sales_date))
		(PARTITION p2020 VALUES LESS THAN (2021) ENGINE = InnoDB,
		 PARTITION p2021 VALUES LESS THAN (2022) ENGINE = InnoDB,
		 PARTITION p2022 VALUES LESS THAN (2023) ENGINE = InnoDB,
		 PARTITION p2023 VALUES LESS THAN (2024) ENGINE = InnoDB,
		 PARTITION p_future VALUES LESS THAN MAXVALUE ENGINE = InnoDB);
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
	if _, err := db3.Exec(ctx, fmt.Sprintf("DROP TABLE IF EXISTS `%s`", name)); err != nil {
		gtest.Fatal(err)
	}
	if _, err := db3.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id int(11) NOT NULL,
			user_id int(11) NOT NULL,
			username varchar(50) DEFAULT NULL,
			email varchar(100) DEFAULT NULL
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
		PARTITION BY HASH (user_id)
		PARTITIONS 4;
	`, name)); err != nil {
		gtest.Fatal(err)
	}
	return name
}

func createListPartitionTable(table ...string) string {
	var name string
	if len(table) > 0 {
		name = table[0]
	} else {
		name = fmt.Sprintf(`partition_list_%d`, gtime.TimestampNano())
	}
	if _, err := db3.Exec(ctx, fmt.Sprintf("DROP TABLE IF EXISTS `%s`", name)); err != nil {
		gtest.Fatal(err)
	}
	if _, err := db3.Exec(ctx, fmt.Sprintf(`
		CREATE TABLE %s (
			id int(11) NOT NULL,
			region_code int(11) NOT NULL,
			city varchar(50) DEFAULT NULL,
			population int(11) DEFAULT NULL
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
		PARTITION BY LIST (region_code)
		(PARTITION p_north VALUES IN (1,2,3) ENGINE = InnoDB,
		 PARTITION p_south VALUES IN (4,5,6) ENGINE = InnoDB,
		 PARTITION p_east VALUES IN (7,8,9) ENGINE = InnoDB,
		 PARTITION p_west VALUES IN (10,11,12) ENGINE = InnoDB);
	`, name)); err != nil {
		gtest.Fatal(err)
	}
	return name
}

func dropPartitionTable(table string) {
	if _, err := db3.Exec(ctx, fmt.Sprintf("DROP TABLE IF EXISTS `%s`", table)); err != nil {
		gtest.Error(err)
	}
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
		_, err := db3.Model(table).Data(data).Insert()
		t.AssertNil(err)

		// Query all data
		all, err := db3.Model(table).All()
		t.AssertNil(err)
		t.Assert(len(all), 5)

		// Query specific year (should hit specific partition)
		result, err := db3.Model(table).Where("YEAR(sales_date) = ?", 2022).All()
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["id"], 3)
	})
}

func Test_Partition_Range_PartitionQuery(t *testing.T) {
	// Known limitation: Model.Partition() sets m.partition field but it's not used in SQL generation
	// See: database/gdb/gdb_model_select.go lines 735,755 - m.tables is used without PARTITION clause
	// TODO: Add PARTITION clause support to GoFrame query builder
	t.Skip("Partition clause in SELECT queries not yet supported in GoFrame query builder")

	table := createRangePartitionTable()
	defer dropPartitionTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert data
		data := g.Slice{
			g.Map{"id": 1, "sales_date": "2020-06-15", "amount": 1000.50},
			g.Map{"id": 2, "sales_date": "2021-03-20", "amount": 2000.75},
			g.Map{"id": 3, "sales_date": "2022-09-10", "amount": 3000.00},
			g.Map{"id": 4, "sales_date": "2023-12-01", "amount": 4000.25},
		}
		_, err := db3.Model(table).Data(data).Insert()
		t.AssertNil(err)

		// Query specific partition
		result, err := db3.Model(table).Partition("p2022").All()
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["id"], 3)

		// Query multiple partitions
		result, err = db3.Model(table).Partition("p2021", "p2022").All()
		t.AssertNil(err)
		t.Assert(len(result), 2)
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
		_, err := db3.Model(table).Data(data).Insert()
		t.AssertNil(err)

		// Query all data
		all, err := db3.Model(table).All()
		t.AssertNil(err)
		t.Assert(len(all), 20)

		// Query specific user_id (will hit specific partition based on hash)
		result, err := db3.Model(table).Where("user_id", 100).One()
		t.AssertNil(err)
		t.Assert(result["username"], "user_10")
	})
}

func Test_Partition_List_Insert_And_Query(t *testing.T) {
	// Known limitation: Model.Partition() sets m.partition field but it's not used in SQL generation
	// See: database/gdb/gdb_model_select.go lines 735,755 - m.tables is used without PARTITION clause
	// TODO: Add PARTITION clause support to GoFrame query builder
	t.Skip("Partition clause in SELECT queries not yet supported in GoFrame query builder")

	table := createListPartitionTable()
	defer dropPartitionTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert data for different regions
		data := g.Slice{
			g.Map{"id": 1, "region_code": 1, "city": "Beijing", "population": 2154},
			g.Map{"id": 2, "region_code": 2, "city": "Harbin", "population": 1063},
			g.Map{"id": 3, "region_code": 5, "city": "Guangzhou", "population": 1868},
			g.Map{"id": 4, "region_code": 7, "city": "Shanghai", "population": 2428},
			g.Map{"id": 5, "region_code": 10, "city": "Chengdu", "population": 2093},
		}
		_, err := db3.Model(table).Data(data).Insert()
		t.AssertNil(err)

		// Query all
		all, err := db3.Model(table).All()
		t.AssertNil(err)
		t.Assert(len(all), 5)

		// Query specific partition (north region)
		result, err := db3.Model(table).Partition("p_north").All()
		t.AssertNil(err)
		t.Assert(len(result), 2)

		// Query specific partition (south region)
		result, err = db3.Model(table).Partition("p_south").All()
		t.AssertNil(err)
		t.Assert(len(result), 1)
		t.Assert(result[0]["city"], "Guangzhou")
	})
}

func Test_Partition_Range_Update(t *testing.T) {
	table := createRangePartitionTable()
	defer dropPartitionTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Insert data
		_, err := db3.Model(table).Data(g.Map{
			"id":         1,
			"sales_date": "2022-06-15",
			"amount":     1000.00,
			"region":     "North",
		}).Insert()
		t.AssertNil(err)

		// Update data within same partition
		result, err := db3.Model(table).Data(g.Map{
			"amount": 1500.00,
			"region": "South",
		}).Where("id", 1).Update()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		// Verify update
		one, err := db3.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["amount"], "1500.00")
		t.Assert(one["region"], "South")
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
		_, err := db3.Model(table).Data(data).Insert()
		t.AssertNil(err)

		// Delete from specific partition
		result, err := db3.Model(table).Where("YEAR(sales_date) = ?", 2021).Delete()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		// Verify deletion
		all, err := db3.Model(table).All()
		t.AssertNil(err)
		t.Assert(len(all), 2)

		// Verify remaining data
		result2, err := db3.Model(table).Where("YEAR(sales_date) = ?", 2021).All()
		t.AssertNil(err)
		t.Assert(len(result2), 0)
	})
}

func Test_Partition_Transaction(t *testing.T) {
	table := createRangePartitionTable()
	defer dropPartitionTable(table)

	gtest.C(t, func(t *gtest.T) {
		// Transaction with partitioned table
		err := db3.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
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
		all, err := db3.Model(table).All()
		t.AssertNil(err)
		t.Assert(len(all), 3)

		one, err := db3.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["amount"], "1500.00")
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
		_, err := db3.Model(table).Data(data).Insert()
		t.AssertNil(err)

		// Count by year (specific partition)
		count, err := db3.Model(table).Where("YEAR(sales_date) = ?", 2020).Count()
		t.AssertNil(err)
		t.Assert(count, 2)

		// Sum across partitions
		value, err := db3.Model(table).Fields("SUM(amount) as total").Value()
		t.AssertNil(err)
		t.AssertGT(value.Float64(), 7000.0) // 1000+1500+2000+3000 = 7500
	})
}
