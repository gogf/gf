package clickhouse

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
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
		Host:  "127.0.0.1",
		Port:  "9000",
		User:  "default",
		Name:  "default",
		Type:  "clickhouse",
		Debug: true,
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
	_, err = connect.Model("visits").Data(g.Map{
		"created": time.Now().Format("2006-01-02 15:04:05"),
	}).Update()
	gtest.AssertNE(err, nil)
	_, err = connect.Model("visits").Update()
	gtest.AssertNE(err, nil)
}

func TestDriverClickhouse_Select(t *testing.T) {
	connect := InitClickhouse()
	gtest.AssertEQ(createClickhouseTable(connect), nil)
	defer dropClickhouseTable(connect)
	data, err := connect.Model("visits").All()
	gtest.AssertNil(err)
	gtest.AssertEQ(len(data), 0)
}

func TestDriver_InsertIgnore(t *testing.T) {
	connect := InitClickhouse()
	_, err := connect.InsertIgnore(context.Background(), "", nil)
	gtest.AssertEQ(err, ErrUnsupportedInsertIgnore)
}

func TestDriver_InsertAndGetId(t *testing.T) {
	connect := InitClickhouse()
	_, err := connect.InsertAndGetId(context.Background(), "", nil)
	gtest.AssertEQ(err, ErrUnsupportedInsertGetId)
}

func TestDriver_Replace(t *testing.T) {
	connect := InitClickhouse()
	_, err := connect.Replace(context.Background(), "", nil)
	gtest.AssertEQ(err, ErrUnsupportedReplace)
}

func TestDriverClickhouse_DoInsertOne(t *testing.T) {
	connect := InitClickhouse()
	gtest.AssertEQ(createClickhouseTable(connect), nil)
	defer dropClickhouseTable(connect)
	_, err := connect.Model("visits").Data(g.Map{
		"id":       grand.Intn(999),
		"duration": float64(grand.Intn(999)),
		"url":      gconv.String(grand.Intn(999)),
		"created":  time.Now().Format("2006-01-02 15:04:05"),
	}).Insert()
	gtest.AssertNil(err)
}

func TestDriver_DoInsertMany(t *testing.T) {
	connect := InitClickhouse()
	gtest.AssertEQ(createClickhouseTable(connect), nil)
	defer dropClickhouseTable(connect)
	tx, err := connect.Begin(context.Background())
	gtest.AssertEQ(err, ErrUnsupportedBegin)
	gtest.AssertNil(tx)
}

func TestDriverClickhouse_DoInsert(t *testing.T) {
	connect := InitClickhouse()
	gtest.AssertEQ(createClickhouseTable(connect), nil)
	type insertItem struct {
		Id       int     `orm:"id"`
		Duration float64 `orm:"duration"`
		Url      string  `orm:"url"`
		Created  string  `orm:"created"`
	}
	var (
		insertUrl = "https://goframe.org"
		total     = 0
		item      = insertItem{
			Id:       0,
			Duration: 1,
			Url:      insertUrl,
			Created:  time.Now().Format("2006-01-02 15:04:05"),
		}
	)
	_, err := connect.Model("visits").Data(item).Insert()
	gtest.AssertNil(err)
	_, err = connect.Model("visits").Data(item).Save()
	gtest.AssertNil(err)
	total, err = connect.Model("visits").Count()
	gtest.AssertNil(err)
	gtest.AssertEQ(total, 2)
	list := []*insertItem{}
	for i := 0; i < 50; i++ {
		list = append(list, &insertItem{
			Id:       grand.Intn(999),
			Duration: float64(grand.Intn(999)),
			Url:      insertUrl,
			Created:  time.Now().Format("2006-01-02 15:04:05"),
		})
	}
	_, err = connect.Model("visits").Data(list).Insert()
	gtest.AssertNil(err)
	_, err = connect.Model("visits").Data(list).Save()
	gtest.AssertNil(err)
	total, err = connect.Model("visits").Count()
	gtest.AssertNil(err)
	gtest.AssertEQ(total, 102)
	dropClickhouseTable(connect)
}

func TestDriverClickhouse_DoExec(t *testing.T) {
	connect := InitClickhouse()
	gtest.AssertNil(createClickhouseTable(connect))
	defer dropClickhouseTable(connect)
	sqlStr := "OPTIMIZE table visits"
	_, err := connect.Exec(context.Background(), sqlStr)
	gtest.AssertNil(err)
}

func TestDriver_DoFilter(t *testing.T) {
	rawSQL := "select * from visits where 1 = 1"
	this := Driver{}
	replaceSQL, _, err := this.DoFilter(nil, nil, rawSQL, nil)
	gtest.AssertNil(err)
	gtest.AssertEQ(rawSQL, replaceSQL)
	rawSQL = "update visit set url = '1'"
	replaceSQL, _, err = this.DoFilter(nil, nil, rawSQL, nil)
	gtest.AssertNil(err)
	// this SQL can't run ,clickhouse will report an error because there is no WHERE statement
	gtest.AssertEQ(replaceSQL, "ALTER TABLE visit update url = '1'")
	rawSQL = "delete from visit"
	replaceSQL, _, err = this.DoFilter(nil, nil, rawSQL, nil)
	gtest.AssertNil(err)
	// this SQL can't run ,clickhouse will report an error because there is no WHERE statement
	gtest.AssertEQ(replaceSQL, "ALTER TABLE visit delete")

	rawSQL = "update visit set url = '1' where url = '0'"
	replaceSQL, _, err = this.DoFilter(nil, nil, rawSQL, nil)
	gtest.AssertNil(err)
	// this SQL can't run ,clickhouse will report an error because there is no WHERE statement
	gtest.AssertEQ(replaceSQL, "ALTER TABLE visit update url = '1' where url = '0'")
	rawSQL = "delete from visit where url='0'"
	replaceSQL, _, err = this.DoFilter(nil, nil, rawSQL, nil)
	gtest.AssertNil(err)
	// this SQL can't run ,clickhouse will report an error because there is no WHERE statement
	gtest.AssertEQ(replaceSQL, "ALTER TABLE visit delete where url='0'")
}

func TestDriver_TableFields(t *testing.T) {
	connect := InitClickhouse()
	gtest.AssertNil(createClickhouseTable(connect))
	defer dropClickhouseTable(connect)
	field, err := connect.TableFields(context.Background(), "visits")
	gtest.AssertNil(err)
	gtest.AssertEQ(len(field), 4)
}
