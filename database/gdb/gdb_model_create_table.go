package gdb

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/text/gstr"
	"log"
	"reflect"
)

// 不同库对应 golang字段类型与数据库类型映射
var typeMapping = map[string]map[reflect.Kind]string{
	"mysql": {
		reflect.Int:     "INT",
		reflect.Int8:    "TINYINT",
		reflect.Int16:   "SMALLINT",
		reflect.Int32:   "INT",
		reflect.Int64:   "BIGINT",
		reflect.Uint:    "INT",
		reflect.Uint8:   "TINYINT",
		reflect.Uint16:  "SMALLINT",
		reflect.Uint32:  "INT",
		reflect.Uint64:  "BIGINT",
		reflect.Float32: "FLOAT",
		reflect.Float64: "DOUBLE",
		reflect.String:  "VARCHAR(255)",
		reflect.Bool:    "BOOLEAN",
		reflect.Ptr:     "DATETIME", // 假设是指向时间类型的指针
	},
	"mariadb": {
		reflect.Int:     "INT",
		reflect.Int8:    "TINYINT",
		reflect.Int16:   "SMALLINT",
		reflect.Int32:   "INT",
		reflect.Int64:   "BIGINT",
		reflect.Uint:    "INT",
		reflect.Uint8:   "TINYINT",
		reflect.Uint16:  "SMALLINT",
		reflect.Uint32:  "INT",
		reflect.Uint64:  "BIGINT",
		reflect.Float32: "FLOAT",
		reflect.Float64: "DOUBLE", // MariaDB 支持 DOUBLE
		reflect.String:  "VARCHAR(255)",
		reflect.Bool:    "BOOLEAN",
		reflect.Ptr:     "DATETIME", // 假设是指向时间类型的指针
	},
	"tidb": {
		reflect.Int:     "INT",
		reflect.Int8:    "TINYINT",
		reflect.Int16:   "SMALLINT",
		reflect.Int32:   "INT",
		reflect.Int64:   "BIGINT",
		reflect.Uint:    "INT",
		reflect.Uint8:   "TINYINT",
		reflect.Uint16:  "SMALLINT",
		reflect.Uint32:  "INT",
		reflect.Uint64:  "BIGINT",
		reflect.Float32: "FLOAT",
		reflect.Float64: "DOUBLE", // TiDB 支持 DOUBLE
		reflect.String:  "VARCHAR(255)",
		reflect.Bool:    "BOOLEAN",
		reflect.Ptr:     "DATETIME", // 假设是指向时间类型的指针
	},
	"mssql": {
		reflect.Int:     "INT",
		reflect.Int8:    "TINYINT",
		reflect.Int16:   "SMALLINT",
		reflect.Int32:   "INT",
		reflect.Int64:   "BIGINT",
		reflect.Uint:    "INT",
		reflect.Uint8:   "TINYINT",
		reflect.Uint16:  "SMALLINT",
		reflect.Uint32:  "INT",
		reflect.Uint64:  "BIGINT",
		reflect.Float32: "FLOAT",
		reflect.Float64: "FLOAT", // MSSQL 支持 FLOAT
		reflect.String:  "NVARCHAR(255)",
		reflect.Bool:    "BIT",
		reflect.Ptr:     "DATETIME2", // 假设是指向时间类型的指针
	},
	"sqlite": {
		reflect.Int:     "INTEGER",
		reflect.Int8:    "INTEGER",
		reflect.Int16:   "INTEGER",
		reflect.Int32:   "INTEGER",
		reflect.Int64:   "INTEGER",
		reflect.Uint:    "INTEGER",
		reflect.Uint8:   "INTEGER",
		reflect.Uint16:  "INTEGER",
		reflect.Uint32:  "INTEGER",
		reflect.Uint64:  "INTEGER",
		reflect.Float32: "REAL",
		reflect.Float64: "REAL", // SQLite 支持 REAL
		reflect.String:  "TEXT",
		reflect.Bool:    "INTEGER",
		reflect.Ptr:     "DATETIME", // 假设是指向时间类型的指针
	},
	"oracle": {
		reflect.Int:     "NUMBER",
		reflect.Int8:    "NUMBER",
		reflect.Int16:   "NUMBER",
		reflect.Int32:   "NUMBER",
		reflect.Int64:   "NUMBER",
		reflect.Uint:    "NUMBER",
		reflect.Uint8:   "NUMBER",
		reflect.Uint16:  "NUMBER",
		reflect.Uint32:  "NUMBER",
		reflect.Uint64:  "NUMBER",
		reflect.Float32: "FLOAT",
		reflect.Float64: "NUMBER", // Oracle 支持 NUMBER
		reflect.String:  "VARCHAR2(255)",
		reflect.Bool:    "NUMBER(1)",
		reflect.Ptr:     "TIMESTAMP", // 假设是指向时间类型的指针
	},
	"dm": {
		reflect.Int:     "INTEGER",
		reflect.Int8:    "TINYINT",
		reflect.Int16:   "SMALLINT",
		reflect.Int32:   "INTEGER",
		reflect.Int64:   "BIGINT",
		reflect.Uint:    "INTEGER",
		reflect.Uint8:   "TINYINT",
		reflect.Uint16:  "SMALLINT",
		reflect.Uint32:  "INTEGER",
		reflect.Uint64:  "BIGINT",
		reflect.Float32: "FLOAT",
		reflect.Float64: "DOUBLE", // DM 支持 DOUBLE
		reflect.String:  "VARCHAR(255)",
		reflect.Bool:    "BOOLEAN",
		reflect.Ptr:     "TIMESTAMP", // 假设是指向时间类型的指针
	},
	"duckdb": {
		reflect.Int:     "INTEGER",
		reflect.Int8:    "TINYINT",
		reflect.Int16:   "SMALLINT",
		reflect.Int32:   "INTEGER",
		reflect.Int64:   "BIGINT",
		reflect.Uint:    "INTEGER",
		reflect.Uint8:   "TINYINT",
		reflect.Uint16:  "SMALLINT",
		reflect.Uint32:  "INTEGER",
		reflect.Uint64:  "BIGINT",
		reflect.Float32: "REAL",
		reflect.Float64: "DOUBLE",
		reflect.String:  "VARCHAR(255)",
		reflect.Bool:    "BOOLEAN",
		reflect.Ptr:     "TIMESTAMP", // 假设是指向时间类型的指针
	},
	"pgsql": {
		reflect.Int:     "INTEGER",
		reflect.Int8:    "SMALLINT",
		reflect.Int16:   "SMALLINT",
		reflect.Int32:   "INTEGER",
		reflect.Int64:   "BIGINT",
		reflect.Uint:    "INTEGER",
		reflect.Uint8:   "SMALLINT",
		reflect.Uint16:  "SMALLINT",
		reflect.Uint32:  "INTEGER",
		reflect.Uint64:  "BIGINT",
		reflect.Float32: "REAL",
		reflect.Float64: "DOUBLE PRECISION",
		reflect.String:  "VARCHAR(255)",
		reflect.Bool:    "BOOLEAN",
		reflect.Ptr:     "TIMESTAMPTZ", // 假设是指向时间类型的指针
	},
	"clickhouse": {
		reflect.Int:     "Int32",
		reflect.Int8:    "Int8",
		reflect.Int16:   "Int16",
		reflect.Int32:   "Int32",
		reflect.Int64:   "Int64",
		reflect.Uint:    "UInt32",
		reflect.Uint8:   "UInt8",
		reflect.Uint16:  "UInt16",
		reflect.Uint32:  "UInt32",
		reflect.Uint64:  "UInt64",
		reflect.Float32: "Float32",
		reflect.Float64: "Float64",
		reflect.String:  "String",
		reflect.Bool:    "UInt8",
		reflect.Ptr:     "DateTime", // 假设是指向时间类型的指针
	},
}

// 自增字段语法映射
var autoIncrementMapping = map[string]string{
	"mysql":      "AUTO_INCREMENT",
	"mariadb":    "AUTO_INCREMENT",
	"tidb":       "AUTO_INCREMENT",
	"mssql":      "IDENTITY(1,1)",
	"sqlite":     "AUTOINCREMENT",
	"oracle":     "GENERATED BY DEFAULT AS IDENTITY",
	"dm":         "AUTO_INCREMENT",
	"duckdb":     "",       // 内存库,不支持传统意义上的自增
	"pgsql":      "SERIAL", //最新版本psql 使用BIGSERIAL
	"clickhouse": "AUTO_INCREMENT",
}

// AutoMigrate 自动迁移模型到数据库。
// 该方法根据模型的结构体字段自动生成对应的数据库表结构。
// 参数:
//
//	ctx - 上下文，用于传递请求范围的信息。
//	model - 需要进行迁移的模型结构体。
//
// 返回值:
//
//	如果迁移过程中发生错误，则返回错误。
func (m *Model) AutoMigrate(ctx context.Context, model interface{}) error {
	// 获取结构体的字段信息
	t := reflect.TypeOf(model)
	if t.Kind() != reflect.Struct {
		return fmt.Errorf("model must be a struct")
	}
	// 获取数据库类型
	dbType := m.db.GetCore().config.Type

	//获取字段类型映射
	typeMap, ok := typeMapping[dbType]
	if !ok {
		return fmt.Errorf("unsupported database type: %s", dbType)
	}

	//获取自增字段语法映射
	autoIncrement, ok := autoIncrementMapping[dbType]
	if !ok {
		return fmt.Errorf("unsupported database type for auto increment: %s", dbType)
	}

	var columns []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// 将字段名称转换为小写并用下划线分隔
		fieldName := gstr.CaseSnake(field.Name)
		//如果存在 orm 中的配置 则取orm中配置的字段名
		ormTag := field.Tag.Get("orm")
		if ormTag != "" {
			fieldName = ormTag
		}

		//使用默认的字段类型映射
		fieldType, ok := typeMap[field.Type.Kind()]
		if !ok {
			return fmt.Errorf("unsupported field type: %s", field.Type.Kind())
		}

		// 构建列定义
		columnDef := fmt.Sprintf(" %s %s", m.db.GetCore().QuoteString(fieldName), fieldType)

		// 添加默认选项
		if field.Type.Kind() == reflect.String {
			columnDef += " NOT NULL DEFAULT ''"
		} else if field.Type.Kind() == reflect.Bool {
			columnDef += " NOT NULL DEFAULT FALSE"
		}

		columns = append(columns, columnDef)
	}

	//默认给 id 字段设置为主键,orm标签目前是配置的字段,后续考虑通过配置orm标签扩展主键自增
	for i, column := range columns {
		if gstr.Contains(column, "id") {
			columns[i] = column + " PRIMARY KEY UNIQUE " + autoIncrement
			break
		}
	}

	// 生成创建表的 SQL 语句
	sql := fmt.Sprintf("CREATE TABLE %s (%s)", m.db.GetCore().QuotePrefixTableName(m.tables), gstr.Join(columns, ", "))

	_, err := m.db.Exec(ctx, sql)
	if err != nil {
		return err
	} else {
		log.Printf("生成表:%s---成功", m.tables)
	}

	return nil
}
