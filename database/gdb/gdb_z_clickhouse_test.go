package gdb

import (
	"context"
	"fmt"
	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/gogf/gf/v2/util/grand"
	"testing"
	"time"
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
func InitClickhouse() (DB, error) {
	return New(ConfigNode{
		Host:   "127.0.0.1",
		Port:   "9000",
		Name:   "default",
		Type:   "clickhouse",
		User:   "default",
		Debug:  true,
		DryRun: false,
	})
}

func TestDriverClickhouse_Create(t *testing.T) {
	connect, err := InitClickhouse()
	if err != nil {
		t.Error(err.Error())
		return
	}
	sqlStr := "CREATE TABLE IF NOT EXISTS visits (id UInt64,duration Float64,url String,created DateTime) ENGINE = MergeTree()  PRIMARY KEY id ORDER BY id"
	_, err = connect.Exec(context.Background(), sqlStr)
	if err != nil {
		t.Error(err.Error())
	}
}

func createClickhouseTable(connect DB) {
	sqlStr := "CREATE TABLE IF NOT EXISTS visits (id UInt64,duration Float64,url String,created DateTime) ENGINE = MergeTree()  PRIMARY KEY id ORDER BY id"
	_, _ = connect.Exec(context.Background(), sqlStr)
}

func dropClickhouseTable(conn DB) {
	sqlStr := fmt.Sprintf("DROP TABLE IF EXISTS `visits`")
	_, _ = conn.Exec(context.Background(), sqlStr)
}

func TestDriverClickhouse_New(t *testing.T) {
	connect, err := InitClickhouse()
	if err != nil {
		t.Error(err.Error())
		return
	}
	err = connect.PingMaster()
	if err != nil {
		t.Error(err.Error())
		return
	}
	err = connect.PingSlave()
	if err != nil {
		t.Error(err.Error())
		return
	}
}

func TestDriverClickhouse_Tables(t *testing.T) {
	connect, err := InitClickhouse()
	if err != nil {
		t.Error(err.Error())
		return
	}
	createClickhouseTable(connect)
	defer dropClickhouseTable(connect)
	tables, err := connect.Tables(context.Background())
	if err != nil {
		t.Error(err.Error())
	}
	t.Logf("%+v", tables)
}

func TestDriverClickhouse_Transaction(t *testing.T) {
	connect, err := InitClickhouse()
	if err != nil {
		t.Error(err.Error())
		return
	}
	createClickhouseTable(connect)
	defer dropClickhouseTable(connect)
	err = connect.Transaction(context.Background(), func(ctx context.Context, tx *TX) error {
		_, err = tx.Update("", nil, nil)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Log(err.Error())
	}
}

func TestDriverClickhouse_DoDelete(t *testing.T) {
	connect, err := InitClickhouse()
	if err != nil {
		t.Error(err.Error())
		return
	}
	createClickhouseTable(connect)
	defer dropClickhouseTable(connect)
	_, err = connect.Model("visits").Where("created >", "2021-01-01 00:00:00").Delete()
	if err != nil {
		t.Error(err.Error())
		return
	}
}

func TestDriverClickhouse_DoCommit(t *testing.T) {
	connect, err := InitClickhouse()
	if err != nil {
		t.Error(err.Error())
		return
	}
	createClickhouseTable(connect)
	defer dropClickhouseTable(connect)
	data, err := connect.Model("visits").All()
	if err != nil {
		t.Error(err.Error())
		return
	}
	for _, item := range data {
		t.Logf("%+v\n", item)
	}
}

func TestDriverClickhouse_DoUpdate(t *testing.T) {
	connect, err := InitClickhouse()
	if err != nil {
		t.Error(err.Error())
		return
	}
	createClickhouseTable(connect)
	defer dropClickhouseTable(connect)
	result, err := connect.Model("visits").Where("created > ", "2021-01-01 15:15:15").Data(Map{
		"created": time.Now().Format("2006-01-02 15:04:05"),
	}).Update()
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Logf("%+v\n", result)
}

func TestDriverClickhouse_DoExec(t *testing.T) {
	connect, err := InitClickhouse()
	if err != nil {
		t.Error(err.Error())
		return
	}
	createClickhouseTable(connect)
	defer dropClickhouseTable(connect)
	sqlStr := "OPTIMIZE table visits"
	_, err = connect.Exec(context.Background(), sqlStr)
	if err != nil {
		t.Log(err.Error())
		return
	}
}

func TestDriverClickhouse_DoInsert(t *testing.T) {
	connect, err := InitClickhouse()
	if err != nil {
		t.Error(err.Error())
		return
	}
	createClickhouseTable(connect)
	defer dropClickhouseTable(connect)
	ctx := context.Background()
	type insertItem struct {
		Id       int     `orm:"id"`
		Duration float64 `orm:"duration"`
		Url      string  `orm:"url"`
		Created  string  `orm:"created"`
	}
	insertUrl := "https://goframe.org"
	// insert one data
	item := insertItem{
		Id:       0,
		Duration: 1,
		Url:      insertUrl,
		Created:  time.Now().Format("2006-01-02 15:04:05"),
	}
	_, err = connect.Model("visits").Ctx(ctx).Data(item).Insert()
	if err != nil {
		t.Error(err.Error())
	}
	_, err = connect.Model("visits").Ctx(ctx).Data(item).InsertIgnore()
	if err != nil {
		t.Error(err.Error())
	}
	_, err = connect.Model("visits").Ctx(ctx).Data(item).InsertAndGetId()
	if err != nil {
		t.Error(err.Error())
	}
	_, err = connect.Model("visits").Ctx(ctx).Data(item).Save()
	if err != nil {
		t.Error(err.Error())
	}
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
	if err != nil {
		t.Error(err.Error())
	}
	_, err = connect.Model("visits").Ctx(ctx).Data(list).InsertIgnore()
	if err != nil {
		t.Error(err.Error())
	}
	_, err = connect.Model("visits").Ctx(ctx).Data(list).InsertAndGetId()
	if err != nil {
		t.Error(err.Error())
	}
	_, err = connect.Model("visits").Ctx(ctx).Data(list).Save()
	if err != nil {
		t.Error(err.Error())
	}
}
