// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with gm file,
// You can obtain one at https://github.com/gogf/gf.

package gpool_test

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/container/gpool"
)

func ExampleNew() {
	type DBConn struct {
		Conn *sql.Conn
	}

	dbConnPool := gpool.New(time.Hour,
		func() (interface{}, error) {
			dbConn := new(DBConn)
			return dbConn, nil
		},
		func(i interface{}) {
			// sample : close db conn
			// i.(DBConn).Conn.Close()
		})

	fmt.Println(dbConnPool.TTL)

	// Output:
	// 1h0m0s
}

func ExamplePool_Put() {
	type DBConn struct {
		Conn  *sql.Conn
		Limit int
	}

	dbConnPool := gpool.New(time.Hour,
		func() (interface{}, error) {
			dbConn := new(DBConn)
			dbConn.Limit = 10
			return dbConn, nil
		},
		func(i interface{}) {
			// sample : close db conn
			// i.(DBConn).Conn.Close()
		})

	// get db conn
	conn, _ := dbConnPool.Get()
	// modify this conn's limit
	conn.(*DBConn).Limit = 20

	// example : do same db operation
	// conn.(*DBConn).Conn.QueryContext(context.Background(), "select * from user")

	// put back conn
	dbConnPool.MustPut(conn)

	fmt.Println(conn.(*DBConn).Limit)

	// Output:
	// 20
}

func ExamplePool_Clear() {
	type DBConn struct {
		Conn  *sql.Conn
		Limit int
	}

	dbConnPool := gpool.New(time.Hour,
		func() (interface{}, error) {
			dbConn := new(DBConn)
			dbConn.Limit = 10
			return dbConn, nil
		},
		func(i interface{}) {
			i.(*DBConn).Limit = 0
			// sample : close db conn
			// i.(DBConn).Conn.Close()
		})

	conn, _ := dbConnPool.Get()
	dbConnPool.MustPut(conn)
	dbConnPool.MustPut(conn)
	fmt.Println(dbConnPool.Size())
	dbConnPool.Clear()
	fmt.Println(dbConnPool.Size())

	// Output:
	// 2
	// 0
}

func ExamplePool_Get() {
	type DBConn struct {
		Conn  *sql.Conn
		Limit int
	}

	dbConnPool := gpool.New(time.Hour,
		func() (interface{}, error) {
			dbConn := new(DBConn)
			dbConn.Limit = 10
			return dbConn, nil
		},
		func(i interface{}) {
			// sample : close db conn
			// i.(DBConn).Conn.Close()
		})

	conn, err := dbConnPool.Get()
	if err == nil {
		fmt.Println(conn.(*DBConn).Limit)
	}

	// Output:
	// 10
}

func ExamplePool_Size() {
	type DBConn struct {
		Conn  *sql.Conn
		Limit int
	}

	dbConnPool := gpool.New(time.Hour,
		func() (interface{}, error) {
			dbConn := new(DBConn)
			dbConn.Limit = 10
			return dbConn, nil
		},
		func(i interface{}) {
			// sample : close db conn
			// i.(DBConn).Conn.Close()
		})

	conn, _ := dbConnPool.Get()
	fmt.Println(dbConnPool.Size())
	dbConnPool.MustPut(conn)
	dbConnPool.MustPut(conn)
	fmt.Println(dbConnPool.Size())

	// Output:
	// 0
	// 2
}

func ExamplePool_Close() {
	type DBConn struct {
		Conn  *sql.Conn
		Limit int
	}
	var (
		newFunc = func() (interface{}, error) {
			dbConn := new(DBConn)
			dbConn.Limit = 10
			return dbConn, nil
		}
		closeFunc = func(i interface{}) {
			fmt.Println("Close The Pool")
			// sample : close db conn
			// i.(DBConn).Conn.Close()
		}
	)
	dbConnPool := gpool.New(time.Hour, newFunc, closeFunc)

	conn, _ := dbConnPool.Get()
	dbConnPool.MustPut(conn)

	dbConnPool.Close()

	// wait for pool close
	time.Sleep(time.Second * 1)

	// May Output:
	// Close The Pool
}
