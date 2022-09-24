package dm

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	_ "gitee.com/chunanyong/dm"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
)

type DriverDM struct {
	*gdb.Core
}

var (
	tableFieldsMap = gmap.New(true)
)

func init() {
	var (
		err         error
		driverObj   = New()
		driverNames = g.SliceStr{"dm"}
	)
	for _, driverName := range driverNames {
		if err = gdb.Register(driverName, driverObj); err != nil {
			panic(err)
		}
	}
}

func New() gdb.Driver {
	return &DriverDM{}
}

func (d *DriverDM) New(core *gdb.Core, node gdb.ConfigNode) (gdb.DB, error) {
	return &DriverDM{
		Core: core,
	}, nil
}

func (d *DriverDM) Open(config gdb.ConfigNode) (db *sql.DB, err error) {
	var (
		source               string
		underlyingDriverName = "dm"
	)
	// Data Source Name of DM8:
	// dm://userName:password@ip:port/dbname
	if config.Link != "" {
		source = config.Link
		// Custom changing the schema in runtime.
		if config.Name != "" {
			source, _ = gregex.ReplaceString(`/([\w\.\-]+)+`, "/"+config.Name, source)
		}
	} else {
		source = fmt.Sprintf(
			"dm://%s:%s@%s:%s/%s?charset=%s",
			config.User, config.Pass, config.Host, config.Port, config.Name, config.Charset,
		)
		// Demo of timezong setting:
		// &loc=Asia/Shanghai
		if config.Timezone != "" {
			source = fmt.Sprintf("%s&loc%s", source, url.QueryEscape(config.Timezone))
		}
	}
	if db, err = sql.Open(underlyingDriverName, source); err != nil {
		err = gerror.WrapCodef(
			gcode.CodeDbOperationError, err,
			`dm.Open failed for driver "%s" by source "%s"`, underlyingDriverName, source,
		)
		return nil, err
	}
	return
}

func (d *DriverDM) GetChars() (charLeft string, charRight string) {
	return `"`, `"`
}

func (d *DriverDM) Tables(ctx context.Context, schema ...string) (tables []string, err error) {
	var result gdb.Result
	link, err := d.SlaveLink(schema...)
	if err != nil {
		return nil, err
	}
	// TODO support multiple schema
	if len(schema) == 0 {
		return nil, gerror.NewCode(gcode.CodeNotSupported, `Schema is empty`)
	}
	result, err = d.DoSelect(ctx, link, fmt.Sprintf(`SELECT * FROM ALL_TABLES WHERE OWNER IN ('%s')`, schema[0]))
	if err != nil {
		return
	}
	for _, m := range result {
		if v, ok := m["IOT_NAME"]; ok {
			tables = append(tables, v.String())
		}
	}
	return
}

func (d *DriverDM) TableFields(ctx context.Context, table string, schema ...string) (fields map[string]*gdb.TableField, err error) {
	// Format dm table
	charL, charR := d.GetChars()
	table = gstr.Trim(table, charL+charR)
	if gstr.Contains(table, " ") {
		return nil, gerror.NewCode(
			gcode.CodeInvalidParameter,
			"function TableFields supports only single table operations",
		)
	}

	// SET schema
	useSchema := d.GetSchema()
	if len(schema) > 0 && schema[0] != "" {
		useSchema = schema[0]
	}

	v := tableFieldsMap.GetOrSetFuncLock(
		fmt.Sprintf(`dm_table_fields_%s_%s@group:%s`, table, useSchema, d.GetGroup()),
		func() interface{} {
			var (
				result gdb.Result
				link   gdb.Link
			)
			if link, err = d.SlaveLink(useSchema); err != nil {
				return nil
			}
			result, err = d.DoSelect(
				ctx, link,
				fmt.Sprintf(`SELECT * FROM ALL_TAB_COLUMNS WHERE OWNER='%s' AND Table_Name= '%s'`, useSchema, strings.ToUpper(table)),
			)
			if err != nil {
				return nil
			}
			fields = make(map[string]*gdb.TableField)
			for _, m := range result {
				// m[NULLABLE] returns "N" "Y"
				// "N" means not null
				// "Y" means could be null
				var nullable bool
				if m["NULLABLE"].String() != "N" {
					nullable = true
				}
				fields[m["COLUMN_NAME"].String()] = &gdb.TableField{
					Index:   m["COLUMN_ID"].Int(),
					Name:    m["COLUMN_NAME"].String(),
					Type:    m["DATA_TYPE"].String(),
					Null:    nullable,
					Default: m["DATA_DEFAULT"].Val(),
					// Key:     m["Key"].String(),
					// Extra:   m["Extra"].String(),
					// Comment: m["Comment"].String(),
				}
			}
			return fields
		},
	)
	if v != nil {
		fields = v.(map[string]*gdb.TableField)
	}
	return
}

// DoFilter deals with the sql string before commits it to underlying sql driver.
func (d *DriverDM) DoFilter(ctx context.Context, link gdb.Link, sql string, args []interface{}) (newSql string, newArgs []interface{}, err error) {
	defer func() {
		newSql, newArgs, err = d.Core.DoFilter(ctx, link, newSql, newArgs)
	}()
	str, _ := gregex.ReplaceString("\"", "", sql)
	str, _ = gregex.ReplaceString("\n", "", str)
	str, _ = gregex.ReplaceString("\t", "", str)
	// There should be no need to capitalize, because it has been done from field processing before
	newSql = str
	// g.Dump("DriverDM.DoFilter()::newSql", newSql)
	newArgs = args
	// g.Dump("DriverDM.DoFilter()::newArgs", newArgs)

	return
}

func (d *DriverDM) DoInsert(
	ctx context.Context, link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption,
) (result sql.Result, err error) {
	switch option.InsertOption {
	case gdb.InsertOptionReplace:
		// TODO
		// TO BE Supported
		return nil, gerror.NewCode(gcode.CodeNotSupported, `Replace operation is not supported by dm driver`)
	case gdb.InsertOptionSave:
		// This syntax currently only supports design tables whose primary key is ID
		listLength := len(list)
		if listLength == 0 {
			return nil, gerror.NewCode(gcode.CodeInvalidRequest, `Save operation list is empty by dm driver`)
		}
		var (
			keysSort     []string
			charL, charR = d.GetChars()
		)
		// Column names need to be aligned in the syntax
		for k := range list[0] {
			keysSort = append(keysSort, k)
		}
		var char = struct {
			charL        string
			charR        string
			valuecharL   string
			valuecharR   string
			duplicateKey string
			keys         []string
		}{
			charL:      charL,
			charR:      charR,
			valuecharL: "'",
			valuecharR: "'",
			// TODO
			// Need to dynamically set the primary key of the table
			duplicateKey: "ID",
			keys:         keysSort,
		}

		// insertKeys: Handle valid keys that need to be inserted and updated
		// insertValues: Handle values ​​that need to be inserted
		// updateValues: Handle values ​​that need to be updated
		// queryValues: Handle only one insert with column name
		insertKeys, insertValues, updateValues, queryValues := parseValue(list[0], char)
		// unionValues: Handling values ​​that need to be inserted and updated
		unionValues := parseUnion(list[1:], char)

		batchResult := new(gdb.SqlResult)
		// parseSql():
		// MERGE INTO {{table}} T1
		// USING ( SELECT {{queryValues}} FROM DUAL
		// {{unionValues}} ) T2
		// ON (T1.{{duplicateKey}} = T2.{{duplicateKey}})
		// WHEN NOT MATCHED THEN
		// INSERT {{insertKeys}} VALUES {{insertValues}}
		// WHEN MATCHED THEN
		// UPDATE SET {{updateValues}}
		sqlStr := parseSql(insertKeys, insertValues, updateValues, queryValues, unionValues, table, char.duplicateKey)
		r, err := d.DoExec(ctx, link, sqlStr)
		if err != nil {
			return r, err
		}
		if n, err := r.RowsAffected(); err != nil {
			return r, err
		} else {
			batchResult.Result = r
			batchResult.Affected += n
		}
		return batchResult, nil
	}
	return d.Core.DoInsert(ctx, link, table, list, option)
}

func parseValue(listOne gdb.Map, char struct {
	charL        string
	charR        string
	valuecharL   string
	valuecharR   string
	duplicateKey string
	keys         []string
}) (insertKeys []string, insertValues []string, updateValues []string, queryValues []string) {
	for _, column := range char.keys {
		if listOne[column] == nil {
			// remove unassigned struct object
			continue
		}
		insertKeys = append(insertKeys, char.charL+column+char.charR)
		insertValues = append(insertValues, "T2."+char.charL+column+char.charR)
		if column != char.duplicateKey {
			updateValues = append(updateValues, fmt.Sprintf(`T1.%s = T2.%s`, char.charL+column+char.charR, char.charL+column+char.charR))
		}

		va := reflect.ValueOf(listOne[column])
		ty := reflect.TypeOf(listOne[column])
		saveValue := ""
		switch ty.Kind() {
		case reflect.String:
			saveValue = va.String()
		case reflect.Int:
			saveValue = strconv.FormatInt(va.Int(), 10)
		case reflect.Int64:
			saveValue = strconv.FormatInt(va.Int(), 10)
		default:
			g.Dump(fmt.Sprintf("column: %v, kind:%v", column, ty.Kind()))
		}
		queryValues = append(queryValues, fmt.Sprintf(char.valuecharL+"%s"+char.valuecharR+" AS "+char.charL+"%s"+char.charR, saveValue, column))
	}
	return
}

func parseUnion(list gdb.List, char struct {
	charL        string
	charR        string
	valuecharL   string
	valuecharR   string
	duplicateKey string
	keys         []string
}) (unionValues []string) {
	for _, mapper := range list {
		var saveValue []string
		for _, column := range char.keys {
			if mapper[column] == nil {
				continue
			}
			va := reflect.ValueOf(mapper[column])
			ty := reflect.TypeOf(mapper[column])
			switch ty.Kind() {
			case reflect.String:
				saveValue = append(saveValue, char.valuecharL+va.String()+char.valuecharR)
			case reflect.Int:
				saveValue = append(saveValue, strconv.FormatInt(va.Int(), 10))
			case reflect.Int64:
				saveValue = append(saveValue, strconv.FormatInt(va.Int(), 10))
			default:
				g.Dump(fmt.Sprintf("column: %v, kind:%v", column, ty.Kind()))
			}
		}
		unionValues = append(unionValues, fmt.Sprintf(`UNION ALL SELECT %s FROM DUAL`, strings.Join(saveValue, ",")))
	}
	return
}

func parseSql(insertKeys, insertValues, updateValues, queryValues, unionValues []string, table, duplicateKey string) (sqlStr string) {
	queryValueStr := strings.Join(queryValues, ",")
	unionValueStr := strings.Join(unionValues, " ")
	insertKeyStr := strings.Join(insertKeys, ",")
	insertValueStr := strings.Join(insertValues, ",")
	updateValueStr := strings.Join(updateValues, ",")
	return fmt.Sprintf(`MERGE INTO %s T1 USING (SELECT %s FROM DUAL %s) T2 ON %s WHEN NOT MATCHED THEN INSERT(%s) VALUES (%s) WHEN MATCHED THEN UPDATE SET %s; 
	COMMIT;`, table, queryValueStr, unionValueStr, fmt.Sprintf("(T1.%s = T2.%s)", duplicateKey, duplicateKey), insertKeyStr, insertValueStr, updateValueStr)
}
