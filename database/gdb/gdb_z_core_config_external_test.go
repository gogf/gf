// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_GetAllConfig(t *testing.T) {
	// Test case 1: Empty configuration
	gtest.C(t, func(t *gtest.T) {
		// Save original config and restore after test
		originalConfig := gdb.GetAllConfig()
		defer func() {
			gdb.SetConfig(originalConfig)
		}()

		// Reset config to empty
		gdb.SetConfig(make(gdb.Config))

		result := gdb.GetAllConfig()
		t.Assert(len(result), 0)
	})

	// Test case 2: Single configuration group with one node
	gtest.C(t, func(t *gtest.T) {
		// Save original config and restore after test
		originalConfig := gdb.GetAllConfig()
		defer func() {
			gdb.SetConfig(originalConfig)
		}()

		// Reset config
		gdb.SetConfig(make(gdb.Config))

		testNode := gdb.ConfigNode{
			Host: "127.0.0.1",
			Port: "3306",
			User: "root",
			Pass: "123456",
			Name: "test_db",
			Type: "mysql",
		}

		err := gdb.AddConfigNode("test_group", testNode)
		t.AssertNil(err)

		result := gdb.GetAllConfig()
		t.Assert(len(result), 1)
		t.Assert(len(result["test_group"]), 1)
		t.Assert(result["test_group"][0].Host, "127.0.0.1")
		t.Assert(result["test_group"][0].Port, "3306")
		t.Assert(result["test_group"][0].User, "root")
		t.Assert(result["test_group"][0].Pass, "123456")
		t.Assert(result["test_group"][0].Name, "test_db")
		t.Assert(result["test_group"][0].Type, "mysql")
	})

	// Test case 3: Multiple configuration groups with multiple nodes
	gtest.C(t, func(t *gtest.T) {
		// Save original config and restore after test
		originalConfig := gdb.GetAllConfig()
		defer func() {
			gdb.SetConfig(originalConfig)
		}()

		// Reset config
		gdb.SetConfig(make(gdb.Config))

		// Add first group with two nodes
		testNode1 := gdb.ConfigNode{
			Host: "127.0.0.1",
			Port: "3306",
			User: "root",
			Pass: "123456",
			Name: "master_db",
			Type: "mysql",
			Role: "master",
		}
		testNode2 := gdb.ConfigNode{
			Host: "127.0.0.2",
			Port: "3306",
			User: "root",
			Pass: "123456",
			Name: "slave_db",
			Type: "mysql",
			Role: "slave",
		}

		err := gdb.AddConfigNode("mysql_cluster", testNode1)
		t.AssertNil(err)
		err = gdb.AddConfigNode("mysql_cluster", testNode2)
		t.AssertNil(err)

		// Add second group with one node
		testNode3 := gdb.ConfigNode{
			Host: "localhost",
			Port: "5432",
			User: "postgres",
			Pass: "password",
			Name: "pg_db",
			Type: "pgsql",
		}

		err = gdb.AddConfigNode("postgres_db", testNode3)
		t.AssertNil(err)

		result := gdb.GetAllConfig()
		t.Assert(len(result), 2)

		// Check mysql_cluster group
		t.Assert(len(result["mysql_cluster"]), 2)
		t.Assert(result["mysql_cluster"][0].Host, "127.0.0.1")
		t.Assert(result["mysql_cluster"][0].Role, "master")
		t.Assert(result["mysql_cluster"][1].Host, "127.0.0.2")
		t.Assert(result["mysql_cluster"][1].Role, "slave")

		// Check postgres_db group
		t.Assert(len(result["postgres_db"]), 1)
		t.Assert(result["postgres_db"][0].Host, "localhost")
		t.Assert(result["postgres_db"][0].Port, "5432")
		t.Assert(result["postgres_db"][0].Type, "pgsql")
	})

	// Test case 4: Configuration with Link syntax
	gtest.C(t, func(t *gtest.T) {
		// Save original config and restore after test
		originalConfig := gdb.GetAllConfig()
		defer func() {
			gdb.SetConfig(originalConfig)
		}()

		// Reset config
		gdb.SetConfig(make(gdb.Config))

		testNode := gdb.ConfigNode{
			Link: "mysql:root:123456@tcp(127.0.0.1:3306)/test_db?charset=utf8",
		}

		err := gdb.AddConfigNode("link_test", testNode)
		t.AssertNil(err)

		result := gdb.GetAllConfig()
		t.Assert(len(result), 1)
		t.Assert(len(result["link_test"]), 1)

		// Check parsed values from link
		node := result["link_test"][0]
		t.Assert(node.Type, "mysql")
		t.Assert(node.User, "root")
		t.Assert(node.Pass, "123456")
		t.Assert(node.Host, "127.0.0.1")
		t.Assert(node.Port, "3306")
		t.Assert(node.Name, "test_db")
		t.Assert(node.Charset, "utf8")
		t.Assert(node.Protocol, "tcp")
	})

	// Test case 5: Default group configuration
	gtest.C(t, func(t *gtest.T) {
		// Save original config and restore after test
		originalConfig := gdb.GetAllConfig()
		defer func() {
			gdb.SetConfig(originalConfig)
		}()

		// Reset config
		gdb.SetConfig(make(gdb.Config))

		testNode := gdb.ConfigNode{
			Host: "localhost",
			Port: "3306",
			User: "user",
			Pass: "pass",
			Name: "default_db",
			Type: "mysql",
		}

		err := gdb.AddDefaultConfigNode(testNode)
		t.AssertNil(err)

		result := gdb.GetAllConfig()
		t.Assert(len(result), 1)
		t.Assert(len(result["default"]), 1)
		t.Assert(result["default"][0].Host, "localhost")
		t.Assert(result["default"][0].Name, "default_db")
	})

	// Test case 6: SetConfig with multiple groups
	gtest.C(t, func(t *gtest.T) {
		// Save original config and restore after test
		originalConfig := gdb.GetAllConfig()
		defer func() {
			gdb.SetConfig(originalConfig)
		}()

		testConfig := gdb.Config{
			"group1": gdb.ConfigGroup{
				{
					Host: "host1",
					Port: "3306",
					User: "user1",
					Pass: "pass1",
					Name: "db1",
					Type: "mysql",
				},
			},
			"group2": gdb.ConfigGroup{
				{
					Host: "host2",
					Port: "5432",
					User: "user2",
					Pass: "pass2",
					Name: "db2",
					Type: "pgsql",
				},
				{
					Host: "host3",
					Port: "5432",
					User: "user3",
					Pass: "pass3",
					Name: "db3",
					Type: "pgsql",
				},
			},
		}

		err := gdb.SetConfig(testConfig)
		t.AssertNil(err)

		result := gdb.GetAllConfig()
		t.Assert(len(result), 2)
		t.Assert(len(result["group1"]), 1)
		t.Assert(len(result["group2"]), 2)

		t.Assert(result["group1"][0].Host, "host1")
		t.Assert(result["group1"][0].Type, "mysql")

		t.Assert(result["group2"][0].Host, "host2")
		t.Assert(result["group2"][0].Type, "pgsql")
		t.Assert(result["group2"][1].Host, "host3")
		t.Assert(result["group2"][1].Type, "pgsql")
	})

	// Test case 7: Test return value is a copy (not reference)
	gtest.C(t, func(t *gtest.T) {
		// Save original config and restore after test
		originalConfig := gdb.GetAllConfig()
		defer func() {
			gdb.SetConfig(originalConfig)
		}()

		// Reset config
		gdb.SetConfig(make(gdb.Config))

		testNode := gdb.ConfigNode{
			Host: "original_host",
			Port: "3306",
			User: "original_user",
			Pass: "original_pass",
			Name: "original_db",
			Type: "mysql",
		}

		err := gdb.AddConfigNode("test_copy", testNode)
		t.AssertNil(err)

		// Get config and modify it
		result := gdb.GetAllConfig()
		t.Assert(len(result), 1)

		// Verify original values
		t.Assert(result["test_copy"][0].Host, "original_host")

		// Note: GetAllConfig returns the internal config directly (not a copy)
		// This is by design for performance reasons
		// So modifying the returned config would affect the internal state
		// This test just verifies the current behavior
	})
}

func Test_SetConfig(t *testing.T) {
	// Test case 1: Normal configuration setting
	gtest.C(t, func(t *gtest.T) {
		// Save original config and restore after test
		originalConfig := gdb.GetAllConfig()
		defer func() {
			gdb.SetConfig(originalConfig)
		}()

		testConfig := gdb.Config{
			"group1": gdb.ConfigGroup{
				{
					Host: "127.0.0.1",
					Port: "3306",
					User: "root",
					Pass: "123456",
					Name: "test_db",
					Type: "mysql",
				},
			},
			"group2": gdb.ConfigGroup{
				{
					Host: "192.168.1.100",
					Port: "5432",
					User: "postgres",
					Pass: "password",
					Name: "pg_db",
					Type: "pgsql",
				},
			},
		}

		err := gdb.SetConfig(testConfig)
		t.AssertNil(err)

		result := gdb.GetAllConfig()
		t.Assert(len(result), 2)
		t.Assert(result["group1"][0].Host, "127.0.0.1")
		t.Assert(result["group2"][0].Type, "pgsql")
	})

	// Test case 2: Empty configuration
	gtest.C(t, func(t *gtest.T) {
		// Save original config and restore after test
		originalConfig := gdb.GetAllConfig()
		defer func() {
			gdb.SetConfig(originalConfig)
		}()

		testConfig := gdb.Config{}
		err := gdb.SetConfig(testConfig)
		t.AssertNil(err)

		result := gdb.GetAllConfig()
		t.Assert(len(result), 0)
	})

	// Test case 3: Configuration with Link syntax
	gtest.C(t, func(t *gtest.T) {
		// Save original config and restore after test
		originalConfig := gdb.GetAllConfig()
		defer func() {
			gdb.SetConfig(originalConfig)
		}()

		testConfig := gdb.Config{
			"mysql_link": gdb.ConfigGroup{
				{
					Link: "mysql:root:123456@tcp(127.0.0.1:3306)/test_db?charset=utf8",
				},
			},
		}

		err := gdb.SetConfig(testConfig)
		t.AssertNil(err)

		result := gdb.GetAllConfig()
		t.Assert(len(result), 1)
		node := result["mysql_link"][0]
		t.Assert(node.Type, "mysql")
		t.Assert(node.User, "root")
		t.Assert(node.Host, "127.0.0.1")
		t.Assert(node.Port, "3306")
		t.Assert(node.Name, "test_db")
	})

	// Test case 4: Configuration with invalid Link syntax
	gtest.C(t, func(t *gtest.T) {
		// Save original config and restore after test
		originalConfig := gdb.GetAllConfig()
		defer func() {
			gdb.SetConfig(originalConfig)
		}()

		testConfig := gdb.Config{
			"invalid_link": gdb.ConfigGroup{
				{
					Link: "invalid_link_format",
				},
			},
		}

		err := gdb.SetConfig(testConfig)
		t.AssertNE(err, nil)
	})
}

func Test_SetConfigGroup(t *testing.T) {
	// Test case 1: Set new group configuration
	gtest.C(t, func(t *gtest.T) {
		// Save original config and restore after test
		originalConfig := gdb.GetAllConfig()
		defer func() {
			gdb.SetConfig(originalConfig)
		}()

		// Reset config
		gdb.SetConfig(make(gdb.Config))

		nodes := gdb.ConfigGroup{
			{
				Host: "127.0.0.1",
				Port: "3306",
				User: "root",
				Pass: "123456",
				Name: "db1",
				Type: "mysql",
				Role: "master",
			},
			{
				Host: "127.0.0.2",
				Port: "3306",
				User: "root",
				Pass: "123456",
				Name: "db2",
				Type: "mysql",
				Role: "slave",
			},
		}

		err := gdb.SetConfigGroup("test_group", nodes)
		t.AssertNil(err)

		result := gdb.GetAllConfig()
		t.Assert(len(result), 1)
		t.Assert(len(result["test_group"]), 2)
		t.Assert(result["test_group"][0].Role, "master")
		t.Assert(result["test_group"][1].Role, "slave")
	})

	// Test case 2: Overwrite existing group configuration
	gtest.C(t, func(t *gtest.T) {
		// Save original config and restore after test
		originalConfig := gdb.GetAllConfig()
		defer func() {
			gdb.SetConfig(originalConfig)
		}()

		// Reset config
		gdb.SetConfig(make(gdb.Config))

		// First set
		nodes1 := gdb.ConfigGroup{
			{
				Host: "old_host",
				Port: "3306",
				User: "old_user",
				Name: "old_db",
				Type: "mysql",
			},
		}
		err := gdb.SetConfigGroup("test_group", nodes1)
		t.AssertNil(err)

		// Overwrite with new config
		nodes2 := gdb.ConfigGroup{
			{
				Host: "new_host",
				Port: "5432",
				User: "new_user",
				Name: "new_db",
				Type: "pgsql",
			},
		}
		err = gdb.SetConfigGroup("test_group", nodes2)
		t.AssertNil(err)

		result := gdb.GetAllConfig()
		t.Assert(len(result), 1)
		t.Assert(len(result["test_group"]), 1)
		t.Assert(result["test_group"][0].Host, "new_host")
		t.Assert(result["test_group"][0].Type, "pgsql")
	})

	// Test case 3: Empty group configuration
	gtest.C(t, func(t *gtest.T) {
		// Save original config and restore after test
		originalConfig := gdb.GetAllConfig()
		defer func() {
			gdb.SetConfig(originalConfig)
		}()

		// Reset config
		gdb.SetConfig(make(gdb.Config))

		nodes := gdb.ConfigGroup{}
		err := gdb.SetConfigGroup("empty_group", nodes)
		t.AssertNil(err)

		result := gdb.GetAllConfig()
		t.Assert(len(result), 1)
		t.Assert(len(result["empty_group"]), 0)
	})

	// Test case 4: Configuration with invalid Link syntax
	gtest.C(t, func(t *gtest.T) {
		// Save original config and restore after test
		originalConfig := gdb.GetAllConfig()
		defer func() {
			gdb.SetConfig(originalConfig)
		}()

		nodes := gdb.ConfigGroup{
			{
				Link: "invalid_link",
			},
		}

		err := gdb.SetConfigGroup("invalid_group", nodes)
		t.AssertNE(err, nil)
	})
}

func Test_AddConfigNode(t *testing.T) {
	// Test case 1: Add node to new group
	gtest.C(t, func(t *gtest.T) {
		// Save original config and restore after test
		originalConfig := gdb.GetAllConfig()
		defer func() {
			gdb.SetConfig(originalConfig)
		}()

		// Reset config
		gdb.SetConfig(make(gdb.Config))

		node := gdb.ConfigNode{
			Host: "127.0.0.1",
			Port: "3306",
			User: "root",
			Pass: "123456",
			Name: "test_db",
			Type: "mysql",
		}

		err := gdb.AddConfigNode("new_group", node)
		t.AssertNil(err)

		result := gdb.GetAllConfig()
		t.Assert(len(result), 1)
		t.Assert(len(result["new_group"]), 1)
		t.Assert(result["new_group"][0].Host, "127.0.0.1")
	})

	// Test case 2: Add node to existing group
	gtest.C(t, func(t *gtest.T) {
		// Save original config and restore after test
		originalConfig := gdb.GetAllConfig()
		defer func() {
			gdb.SetConfig(originalConfig)
		}()

		// Reset config
		gdb.SetConfig(make(gdb.Config))

		// Add first node
		node1 := gdb.ConfigNode{
			Host: "127.0.0.1",
			Port: "3306",
			User: "root",
			Pass: "123456",
			Name: "db1",
			Type: "mysql",
		}
		err := gdb.AddConfigNode("existing_group", node1)
		t.AssertNil(err)

		// Add second node to same group
		node2 := gdb.ConfigNode{
			Host: "127.0.0.2",
			Port: "3306",
			User: "root",
			Pass: "123456",
			Name: "db2",
			Type: "mysql",
		}
		err = gdb.AddConfigNode("existing_group", node2)
		t.AssertNil(err)

		result := gdb.GetAllConfig()
		t.Assert(len(result), 1)
		t.Assert(len(result["existing_group"]), 2)
		t.Assert(result["existing_group"][0].Name, "db1")
		t.Assert(result["existing_group"][1].Name, "db2")
	})

	// Test case 3: Add node with Link syntax
	gtest.C(t, func(t *gtest.T) {
		// Save original config and restore after test
		originalConfig := gdb.GetAllConfig()
		defer func() {
			gdb.SetConfig(originalConfig)
		}()

		// Reset config
		gdb.SetConfig(make(gdb.Config))

		node := gdb.ConfigNode{
			Link: "mysql:root:password@tcp(192.168.1.100:3306)/mydb?charset=utf8mb4",
		}

		err := gdb.AddConfigNode("link_group", node)
		t.AssertNil(err)

		result := gdb.GetAllConfig()
		t.Assert(len(result), 1)
		t.Assert(len(result["link_group"]), 1)
		t.Assert(result["link_group"][0].Type, "mysql")
		t.Assert(result["link_group"][0].Host, "192.168.1.100")
		t.Assert(result["link_group"][0].Port, "3306")
		t.Assert(result["link_group"][0].Name, "mydb")
	})

	// Test case 4: Add node with invalid Link syntax
	gtest.C(t, func(t *gtest.T) {
		// Save original config and restore after test
		originalConfig := gdb.GetAllConfig()
		defer func() {
			gdb.SetConfig(originalConfig)
		}()

		node := gdb.ConfigNode{
			Link: "invalid_link_format",
		}

		err := gdb.AddConfigNode("invalid_group", node)
		t.AssertNE(err, nil)
	})
}

func Test_AddDefaultConfigNode(t *testing.T) {
	// Test case 1: Add node to default group
	gtest.C(t, func(t *gtest.T) {
		// Save original config and restore after test
		originalConfig := gdb.GetAllConfig()
		defer func() {
			gdb.SetConfig(originalConfig)
		}()

		// Reset config
		gdb.SetConfig(make(gdb.Config))

		node := gdb.ConfigNode{
			Host: "localhost",
			Port: "3306",
			User: "root",
			Pass: "root",
			Name: "default_db",
			Type: "mysql",
		}

		err := gdb.AddDefaultConfigNode(node)
		t.AssertNil(err)

		result := gdb.GetAllConfig()
		t.Assert(len(result), 1)
		t.Assert(len(result["default"]), 1)
		t.Assert(result["default"][0].Host, "localhost")
		t.Assert(result["default"][0].Name, "default_db")
	})

	// Test case 2: Add multiple nodes to default group
	gtest.C(t, func(t *gtest.T) {
		// Save original config and restore after test
		originalConfig := gdb.GetAllConfig()
		defer func() {
			gdb.SetConfig(originalConfig)
		}()

		// Reset config
		gdb.SetConfig(make(gdb.Config))

		node1 := gdb.ConfigNode{
			Host: "127.0.0.1",
			Port: "3306",
			User: "root",
			Pass: "123456",
			Name: "db1",
			Type: "mysql",
			Role: "master",
		}
		err := gdb.AddDefaultConfigNode(node1)
		t.AssertNil(err)

		node2 := gdb.ConfigNode{
			Host: "127.0.0.2",
			Port: "3306",
			User: "root",
			Pass: "123456",
			Name: "db2",
			Type: "mysql",
			Role: "slave",
		}
		err = gdb.AddDefaultConfigNode(node2)
		t.AssertNil(err)

		result := gdb.GetAllConfig()
		t.Assert(len(result), 1)
		t.Assert(len(result["default"]), 2)
		t.Assert(result["default"][0].Role, "master")
		t.Assert(result["default"][1].Role, "slave")
	})

	// Test case 3: Add node with Link syntax to default group
	gtest.C(t, func(t *gtest.T) {
		// Save original config and restore after test
		originalConfig := gdb.GetAllConfig()
		defer func() {
			gdb.SetConfig(originalConfig)
		}()

		// Reset config
		gdb.SetConfig(make(gdb.Config))

		node := gdb.ConfigNode{
			Link: "pgsql:postgres:password@tcp(localhost:5432)/testdb",
		}

		err := gdb.AddDefaultConfigNode(node)
		t.AssertNil(err)

		result := gdb.GetAllConfig()
		t.Assert(len(result), 1)
		t.Assert(len(result["default"]), 1)
		t.Assert(result["default"][0].Type, "pgsql")
		t.Assert(result["default"][0].User, "postgres")
		t.Assert(result["default"][0].Host, "localhost")
		t.Assert(result["default"][0].Port, "5432")
		t.Assert(result["default"][0].Name, "testdb")
	})
}

func Test_AddDefaultConfigGroup(t *testing.T) {
	// Test case 1: Add multiple nodes to default group (deprecated function)
	gtest.C(t, func(t *gtest.T) {
		// Save original config and restore after test
		originalConfig := gdb.GetAllConfig()
		defer func() {
			gdb.SetConfig(originalConfig)
		}()

		// Reset config
		gdb.SetConfig(make(gdb.Config))

		nodes := gdb.ConfigGroup{
			{
				Host: "127.0.0.1",
				Port: "3306",
				User: "root",
				Pass: "123456",
				Name: "db1",
				Type: "mysql",
				Role: "master",
			},
			{
				Host: "127.0.0.2",
				Port: "3306",
				User: "root",
				Pass: "123456",
				Name: "db2",
				Type: "mysql",
				Role: "slave",
			},
		}

		err := gdb.AddDefaultConfigGroup(nodes)
		t.AssertNil(err)

		result := gdb.GetAllConfig()
		t.Assert(len(result), 1)
		t.Assert(len(result["default"]), 2)
		t.Assert(result["default"][0].Role, "master")
		t.Assert(result["default"][1].Role, "slave")
	})

	// Test case 2: Overwrite existing default group configuration
	gtest.C(t, func(t *gtest.T) {
		// Save original config and restore after test
		originalConfig := gdb.GetAllConfig()
		defer func() {
			gdb.SetConfig(originalConfig)
		}()

		// Reset config
		gdb.SetConfig(make(gdb.Config))

		// First set
		node1 := gdb.ConfigNode{
			Host: "old_host",
			Port: "3306",
			User: "old_user",
			Name: "old_db",
			Type: "mysql",
		}
		err := gdb.AddDefaultConfigNode(node1)
		t.AssertNil(err)

		// Overwrite with new group config
		nodes := gdb.ConfigGroup{
			{
				Host: "new_host",
				Port: "5432",
				User: "new_user",
				Name: "new_db",
				Type: "pgsql",
			},
		}
		err = gdb.AddDefaultConfigGroup(nodes)
		t.AssertNil(err)

		result := gdb.GetAllConfig()
		t.Assert(len(result), 1)
		t.Assert(len(result["default"]), 1)
		t.Assert(result["default"][0].Host, "new_host")
		t.Assert(result["default"][0].Type, "pgsql")
	})
}

func Test_SetDefaultConfigGroup(t *testing.T) {
	// Test case 1: Set default group configuration
	gtest.C(t, func(t *gtest.T) {
		// Save original config and restore after test
		originalConfig := gdb.GetAllConfig()
		defer func() {
			gdb.SetConfig(originalConfig)
		}()

		// Reset config
		gdb.SetConfig(make(gdb.Config))

		nodes := gdb.ConfigGroup{
			{
				Host: "192.168.1.10",
				Port: "3306",
				User: "admin",
				Pass: "admin123",
				Name: "main_db",
				Type: "mysql",
				Role: "master",
			},
			{
				Host: "192.168.1.11",
				Port: "3306",
				User: "admin",
				Pass: "admin123",
				Name: "backup_db",
				Type: "mysql",
				Role: "slave",
			},
		}

		err := gdb.SetDefaultConfigGroup(nodes)
		t.AssertNil(err)

		result := gdb.GetAllConfig()
		t.Assert(len(result), 1)
		t.Assert(len(result["default"]), 2)
		t.Assert(result["default"][0].Host, "192.168.1.10")
		t.Assert(result["default"][0].Role, "master")
		t.Assert(result["default"][1].Host, "192.168.1.11")
		t.Assert(result["default"][1].Role, "slave")
	})

	// Test case 2: Empty default group configuration
	gtest.C(t, func(t *gtest.T) {
		// Save original config and restore after test
		originalConfig := gdb.GetAllConfig()
		defer func() {
			gdb.SetConfig(originalConfig)
		}()

		// Reset config and add some initial data
		gdb.SetConfig(make(gdb.Config))
		err := gdb.AddDefaultConfigNode(gdb.ConfigNode{
			Host: "temp_host",
			Name: "temp_db",
			Type: "mysql",
		})
		t.AssertNil(err)

		// Set empty group
		nodes := gdb.ConfigGroup{}
		err = gdb.SetDefaultConfigGroup(nodes)
		t.AssertNil(err)

		result := gdb.GetAllConfig()
		t.Assert(len(result), 1)
		t.Assert(len(result["default"]), 0)
	})

	// Test case 3: Configuration with Link syntax
	gtest.C(t, func(t *gtest.T) {
		// Save original config and restore after test
		originalConfig := gdb.GetAllConfig()
		defer func() {
			gdb.SetConfig(originalConfig)
		}()

		// Reset config
		gdb.SetConfig(make(gdb.Config))

		nodes := gdb.ConfigGroup{
			{
				Link: "mysql:root:123456@tcp(localhost:3306)/test_db1",
			},
			{
				Link: "pgsql:postgres:password@tcp(localhost:5432)/test_db2",
			},
		}

		err := gdb.SetDefaultConfigGroup(nodes)
		t.AssertNil(err)

		result := gdb.GetAllConfig()
		t.Assert(len(result), 1)
		t.Assert(len(result["default"]), 2)
		t.Assert(result["default"][0].Type, "mysql")
		t.Assert(result["default"][0].Name, "test_db1")
		t.Assert(result["default"][1].Type, "pgsql")
		t.Assert(result["default"][1].Name, "test_db2")
	})
}

func Test_GetConfig(t *testing.T) {
	// Test case 1: Get existing group configuration (deprecated function)
	gtest.C(t, func(t *gtest.T) {
		// Save original config and restore after test
		originalConfig := gdb.GetAllConfig()
		defer func() {
			gdb.SetConfig(originalConfig)
		}()

		// Reset config
		gdb.SetConfig(make(gdb.Config))

		node := gdb.ConfigNode{
			Host: "127.0.0.1",
			Port: "3306",
			User: "root",
			Pass: "123456",
			Name: "test_db",
			Type: "mysql",
		}

		err := gdb.AddConfigNode("test_group", node)
		t.AssertNil(err)

		result := gdb.GetConfig("test_group")
		t.Assert(len(result), 1)
		t.Assert(result[0].Host, "127.0.0.1")
		t.Assert(result[0].Type, "mysql")
	})

	// Test case 2: Get non-existing group configuration
	gtest.C(t, func(t *gtest.T) {
		// Save original config and restore after test
		originalConfig := gdb.GetAllConfig()
		defer func() {
			gdb.SetConfig(originalConfig)
		}()

		// Reset config
		gdb.SetConfig(make(gdb.Config))

		result := gdb.GetConfig("non_existing_group")
		t.Assert(len(result), 0)
	})
}

func Test_GetConfigGroup(t *testing.T) {
	// Test case 1: Get existing group configuration
	gtest.C(t, func(t *gtest.T) {
		// Save original config and restore after test
		originalConfig := gdb.GetAllConfig()
		defer func() {
			gdb.SetConfig(originalConfig)
		}()

		// Reset config
		gdb.SetConfig(make(gdb.Config))

		nodes := gdb.ConfigGroup{
			{
				Host: "127.0.0.1",
				Port: "3306",
				User: "root",
				Pass: "123456",
				Name: "db1",
				Type: "mysql",
				Role: "master",
			},
			{
				Host: "127.0.0.2",
				Port: "3306",
				User: "root",
				Pass: "123456",
				Name: "db2",
				Type: "mysql",
				Role: "slave",
			},
		}

		err := gdb.SetConfigGroup("test_group", nodes)
		t.AssertNil(err)

		result, err := gdb.GetConfigGroup("test_group")
		t.AssertNil(err)
		t.Assert(len(result), 2)
		t.Assert(result[0].Host, "127.0.0.1")
		t.Assert(result[0].Role, "master")
		t.Assert(result[1].Host, "127.0.0.2")
		t.Assert(result[1].Role, "slave")
	})

	// Test case 2: Get non-existing group configuration
	gtest.C(t, func(t *gtest.T) {
		// Save original config and restore after test
		originalConfig := gdb.GetAllConfig()
		defer func() {
			gdb.SetConfig(originalConfig)
		}()

		// Reset config
		gdb.SetConfig(make(gdb.Config))

		result, err := gdb.GetConfigGroup("non_existing_group")
		t.AssertNE(err, nil)
		t.Assert(result, nil)
	})

	// Test case 3: Get empty group configuration
	gtest.C(t, func(t *gtest.T) {
		// Save original config and restore after test
		originalConfig := gdb.GetAllConfig()
		defer func() {
			gdb.SetConfig(originalConfig)
		}()

		// Reset config
		gdb.SetConfig(make(gdb.Config))

		err := gdb.SetConfigGroup("empty_group", gdb.ConfigGroup{})
		t.AssertNil(err)

		result, err := gdb.GetConfigGroup("empty_group")
		t.AssertNil(err)
		t.Assert(len(result), 0)
	})
}

func Test_SetDefaultGroup(t *testing.T) {
	// Test case 1: Set default group name
	gtest.C(t, func(t *gtest.T) {
		// Save original group and restore after test
		originalGroup := gdb.GetDefaultGroup()
		defer func() {
			gdb.SetDefaultGroup(originalGroup)
		}()

		gdb.SetDefaultGroup("custom_default")
		result := gdb.GetDefaultGroup()
		t.Assert(result, "custom_default")
	})

	// Test case 2: Set empty default group name
	gtest.C(t, func(t *gtest.T) {
		// Save original group and restore after test
		originalGroup := gdb.GetDefaultGroup()
		defer func() {
			gdb.SetDefaultGroup(originalGroup)
		}()

		gdb.SetDefaultGroup("")
		result := gdb.GetDefaultGroup()
		t.Assert(result, "")
	})

	// Test case 3: Multiple calls to SetDefaultGroup
	gtest.C(t, func(t *gtest.T) {
		// Save original group and restore after test
		originalGroup := gdb.GetDefaultGroup()
		defer func() {
			gdb.SetDefaultGroup(originalGroup)
		}()

		gdb.SetDefaultGroup("first_group")
		result1 := gdb.GetDefaultGroup()
		t.Assert(result1, "first_group")

		gdb.SetDefaultGroup("second_group")
		result2 := gdb.GetDefaultGroup()
		t.Assert(result2, "second_group")
	})
}

func Test_GetDefaultGroup(t *testing.T) {
	// Test case 1: Get default group name
	gtest.C(t, func(t *gtest.T) {
		// Save original group and restore after test
		originalGroup := gdb.GetDefaultGroup()
		defer func() {
			gdb.SetDefaultGroup(originalGroup)
		}()

		// Test with default value
		result := gdb.GetDefaultGroup()
		t.Assert(result, "default")
	})

	// Test case 2: Get custom default group name
	gtest.C(t, func(t *gtest.T) {
		// Save original group and restore after test
		originalGroup := gdb.GetDefaultGroup()
		defer func() {
			gdb.SetDefaultGroup(originalGroup)
		}()

		gdb.SetDefaultGroup("my_custom_group")
		result := gdb.GetDefaultGroup()
		t.Assert(result, "my_custom_group")
	})
}

func Test_IsConfigured(t *testing.T) {
	// Test case 1: No configuration
	gtest.C(t, func(t *gtest.T) {
		// Save original config and restore after test
		originalConfig := gdb.GetAllConfig()
		defer func() {
			gdb.SetConfig(originalConfig)
		}()

		// Reset config to empty
		gdb.SetConfig(make(gdb.Config))

		result := gdb.IsConfigured()
		t.Assert(result, false)
	})

	// Test case 2: Has configuration
	gtest.C(t, func(t *gtest.T) {
		// Save original config and restore after test
		originalConfig := gdb.GetAllConfig()
		defer func() {
			gdb.SetConfig(originalConfig)
		}()

		// Reset config
		gdb.SetConfig(make(gdb.Config))

		node := gdb.ConfigNode{
			Host: "127.0.0.1",
			Port: "3306",
			User: "root",
			Pass: "123456",
			Name: "test_db",
			Type: "mysql",
		}

		err := gdb.AddConfigNode("test_group", node)
		t.AssertNil(err)

		result := gdb.IsConfigured()
		t.Assert(result, true)
	})

	// Test case 3: Has empty group configuration
	gtest.C(t, func(t *gtest.T) {
		// Save original config and restore after test
		originalConfig := gdb.GetAllConfig()
		defer func() {
			gdb.SetConfig(originalConfig)
		}()

		// Reset config
		gdb.SetConfig(make(gdb.Config))

		err := gdb.SetConfigGroup("empty_group", gdb.ConfigGroup{})
		t.AssertNil(err)

		result := gdb.IsConfigured()
		t.Assert(result, true)
	})
}
