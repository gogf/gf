package gdb

import (
	"context"
	"database/sql"
	_ "github.com/ClickHouse/clickhouse-go"
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
	config := ConfigNode{
		Host:   "127.0.0.1",
		Port:   "9000",
		Name:   "default",
		Type:   "clickhouse",
		Debug:  true,
		DryRun: false,
	}
	AddDefaultConfigNode(config)
	return New()
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

func TestDriverClickhouse_TableFields(t *testing.T) {
	connect, err := InitClickhouse()
	if err != nil {
		t.Error(err.Error())
		return
	}
	tables, err := connect.TableFields(context.Background(), "visits")
	if err != nil {
		t.Error(err.Error())
	}
	for k, item := range tables {
		t.Logf("tables:%v,data=%+v\n", k, item)
	}
}

func TestDriverClickhouse_Tables(t *testing.T) {
	connect, err := InitClickhouse()
	if err != nil {
		t.Error(err.Error())
		return
	}
	tables, err := connect.Tables(context.Background(), "")
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
	result, err := connect.Model("visits").Where("created > ", "2021-01-01 15:15:15").Data(Map{
		"created": time.Now().Format("2006-01-02 15:04:05"),
	}).Update()
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Logf("%+v\n", result)
}

func TestDriverClickhouse_DoInsert(t *testing.T) {
	connect, err := InitClickhouse()
	if err != nil {
		t.Error(err.Error())
		return
	}
	ctx := context.Background()
	// insert one data
	type insertItem struct {
	}
	var (
		insertOneDataList = []interface{}{
			insertItem{},
			map[string]string{},
		}
		result   sql.Result
		resultId int64
	)
	for _, item := range insertOneDataList {
		result, err = connect.Model().Ctx(ctx).Data(item).Insert()
		if err != nil {
			t.Error(err.Error())
		}
		t.Logf("%+v\n", result)
		result, err = connect.Model().Ctx(ctx).Data(item).InsertIgnore()
		if err != nil {
			t.Error(err.Error())
		}
		t.Logf("%+v\n", result)
		resultId, err = connect.Model().Ctx(ctx).Data(item).InsertAndGetId()
		if err != nil {
			t.Error(err.Error())
		}
		t.Logf("%+v\n", resultId)
		result, err = connect.Model().Ctx(ctx).Data(item).Save()
		if err != nil {
			t.Error(err.Error())
		}
		t.Logf("%+v\n", result)
	}

}
