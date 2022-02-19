package clickhouse

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/grand"
)

// table DDL
// CREATE TABLE visits
// (
//    id UInt64,
//    duration Float64,
//    url String,
//    created DateTime
//)
// ENGINE = MergeTree()
// PRIMARY KEY id
// ORDER BY id
func InitClickhouse() gdb.DB {
	connect, err := gdb.New(gdb.ConfigNode{
		Host:     "127.0.0.1",
		Port:     "9000",
		User:     "default",
		Name:     "default",
		Type:     "clickhouse",
		Debug:    true,
		Compress: true,
	})
	gtest.AssertNil(err)
	gtest.AssertNE(connect, nil)
	return connect
}

func TestDriverClickhouse_Create(t *testing.T) {
	gtest.AssertNil(createClickhouseTable(InitClickhouse()))
}

func createClickhouseTable(connect gdb.DB) error {
	sqlStr := "CREATE TABLE IF NOT EXISTS visits (id UInt64,duration Float64,url String,created DateTime) ENGINE = MergeTree()  PRIMARY KEY id ORDER BY id"
	_, err := connect.Exec(context.Background(), sqlStr)
	return err
}

func dropClickhouseTable(conn gdb.DB) {
	sqlStr := fmt.Sprintf("DROP TABLE IF EXISTS `visits`")
	_, _ = conn.Exec(context.Background(), sqlStr)
}

func TestDriverClickhouse_New(t *testing.T) {
	connect := InitClickhouse()
	gtest.AssertNE(connect, nil)
	gtest.AssertNil(connect.PingMaster())
	gtest.AssertNil(connect.PingSlave())
}

func TestDriverClickhouse_Tables(t *testing.T) {
	connect := InitClickhouse()
	gtest.AssertEQ(createClickhouseTable(connect), nil)
	defer dropClickhouseTable(connect)
	tables, err := connect.Tables(context.Background())
	gtest.AssertNil(err)
	gtest.AssertNE(len(tables), 0)
}

func TestDriverClickhouse_Transaction(t *testing.T) {
	connect := InitClickhouse()
	defer dropClickhouseTable(connect)
	gtest.AssertNE(connect.Transaction(context.Background(), func(ctx context.Context, tx *gdb.TX) error {
		return nil
	}), nil)
}

func TestDriverClickhouse_DoDelete(t *testing.T) {
	connect := InitClickhouse()
	gtest.AssertEQ(createClickhouseTable(connect), nil)
	defer dropClickhouseTable(connect)
	_, err := connect.Model("visits").Where("created >", "2021-01-01 00:00:00").Delete()
	gtest.AssertNil(err)
}

func TestDriverClickhouse_DoUpdate(t *testing.T) {
	connect := InitClickhouse()
	gtest.AssertEQ(createClickhouseTable(connect), nil)
	defer dropClickhouseTable(connect)
	_, err := connect.Model("visits").Where("created > ", "2021-01-01 15:15:15").Data(g.Map{
		"created": time.Now().Format("2006-01-02 15:04:05"),
	}).Update()
	gtest.AssertNil(err)
}

func TestDriverClickhouse_Select(t *testing.T) {
	connect := InitClickhouse()
	gtest.AssertEQ(createClickhouseTable(connect), nil)
	defer dropClickhouseTable(connect)
	data, err := connect.Model("visits").All()
	gtest.AssertNil(err)
	gtest.AssertEQ(len(data), 0)
}

func TestDriverClickhouse_DoInsert(t *testing.T) {
	connect := InitClickhouse()
	gtest.AssertEQ(createClickhouseTable(connect), nil)
	defer dropClickhouseTable(connect)
	type insertItem struct {
		Id       int     `orm:"id"`
		Duration float64 `orm:"duration"`
		Url      string  `orm:"url"`
		Created  string  `orm:"created"`
	}
	var (
		ctx       = context.Background()
		insertUrl = "https://goframe.org"
		// insert one data
		item = insertItem{
			Id:       0,
			Duration: 1,
			Url:      insertUrl,
			Created:  time.Now().Format("2006-01-02 15:04:05"),
		}
	)
	_, err := connect.Model("visits").Ctx(ctx).Data(item).Insert()
	gtest.AssertNil(err)
	_, err = connect.Model("visits").Ctx(ctx).Data(item).InsertIgnore()
	gtest.AssertNil(err)

	_, err = connect.Model("visits").Ctx(ctx).Data(item).InsertAndGetId()
	_, err = connect.Model("visits").Ctx(ctx).Data(item).Save()
	gtest.AssertNil(err)
	// insert array data
	list := []*insertItem{}
	for i := 0; i < 999; i++ {
		list = append(list, &insertItem{
			Id:       grand.Intn(999),
			Duration: float64(grand.Intn(999)),
			Url:      insertUrl,
			Created:  time.Now().Format("2006-01-02 15:04:05"),
		})
	}
	_, err = connect.Model("visits").Ctx(ctx).Data(list).Insert()
	gtest.AssertNil(err)
	_, err = connect.Model("visits").Ctx(ctx).Data(list).InsertIgnore()
	gtest.AssertNil(err)
	_, err = connect.Model("visits").Ctx(ctx).Data(list).InsertAndGetId()
	_, err = connect.Model("visits").Ctx(ctx).Data(list).Save()
	gtest.AssertNil(err)
}

func TestDriverClickhouse_DoExec(t *testing.T) {
	connect := InitClickhouse()
	gtest.AssertNil(createClickhouseTable(connect))
	defer dropClickhouseTable(connect)
	sqlStr := "OPTIMIZE table visits"
	_, err := connect.Exec(context.Background(), sqlStr)
	gtest.AssertNil(err)
}
