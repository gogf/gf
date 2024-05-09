// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gendao

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/mod/modfile"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/utils"
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gtag"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
)

const (
	CGenDaoConfig = `gfcli.gen.dao`
	CGenDaoUsage  = `gf gen dao [OPTION]`
	CGenDaoBrief  = `automatically generate go files for dao/do/entity`
	CGenDaoEg     = `
gf gen dao
gf gen dao -l "mysql:root:12345678@tcp(127.0.0.1:3306)/test"
gf gen dao -p ./model -g user-center -t user,user_detail,user_login
gf gen dao -r user_
`

	CGenDaoAd = `
CONFIGURATION SUPPORT
    Options are also supported by configuration file.
    It's suggested using configuration file instead of command line arguments making producing.
    The configuration node name is "gfcli.gen.dao", which also supports multiple databases, for example(config.yaml):
	gfcli:
	  gen:
		dao:
		- link:     "mysql:root:12345678@tcp(127.0.0.1:3306)/test"
		  tables:   "order,products"
		  jsonCase: "CamelLower"
		- link:   "mysql:root:12345678@tcp(127.0.0.1:3306)/primary"
		  path:   "./my-app"
		  prefix: "primary_"
		  tables: "user, userDetail"
		  typeMapping:
			decimal:
			  type:   decimal.Decimal
			  import: github.com/shopspring/decimal
			numeric:
			  type: string
		  fieldMapping:
			table_name.field_name:
			  type:   decimal.Decimal
			  import: github.com/shopspring/decimal
`
	CGenDaoBriefPath              = `directory path for generated files`
	CGenDaoBriefLink              = `database configuration, the same as the ORM configuration of GoFrame`
	CGenDaoBriefTables            = `generate models only for given tables, multiple table names separated with ','`
	CGenDaoBriefTablesEx          = `generate models excluding given tables, multiple table names separated with ','`
	CGenDaoBriefPrefix            = `add prefix for all table of specified link/database tables`
	CGenDaoBriefRemovePrefix      = `remove specified prefix of the table, multiple prefix separated with ','`
	CGenDaoBriefRemoveFieldPrefix = `remove specified prefix of the field, multiple prefix separated with ','`
	CGenDaoBriefStdTime           = `use time.Time from stdlib instead of gtime.Time for generated time/date fields of tables`
	CGenDaoBriefWithTime          = `add created time for auto produced go files`
	CGenDaoBriefGJsonSupport      = `use gJsonSupport to use *gjson.Json instead of string for generated json fields of tables`
	CGenDaoBriefImportPrefix      = `custom import prefix for generated go files`
	CGenDaoBriefDaoPath           = `directory path for storing generated dao files under path`
	CGenDaoBriefDoPath            = `directory path for storing generated do files under path`
	CGenDaoBriefEntityPath        = `directory path for storing generated entity files under path`
	CGenDaoBriefOverwriteDao      = `overwrite all dao files both inside/outside internal folder`
	CGenDaoBriefModelFile         = `custom file name for storing generated model content`
	CGenDaoBriefModelFileForDao   = `custom file name generating model for DAO operations like Where/Data. It's empty in default`
	CGenDaoBriefDescriptionTag    = `add comment to description tag for each field`
	CGenDaoBriefNoJsonTag         = `no json tag will be added for each field`
	CGenDaoBriefNoModelComment    = `no model comment will be added for each field`
	CGenDaoBriefClear             = `delete all generated go files that do not exist in database`
	CGenDaoBriefTypeMapping       = `custom local type mapping for generated struct attributes relevant to fields of table`
	CGenDaoBriefFieldMapping      = `custom local type mapping for generated struct attributes relevant to specific fields of table`
	CGenDaoBriefGroup             = `
specifying the configuration group name of database for generated ORM instance,
it's not necessary and the default value is "default"
`
	CGenDaoBriefJsonCase = `
generated json tag case for model struct, cases are as follows:
| Case            | Example            |
|---------------- |--------------------|
| Camel           | AnyKindOfString    |
| CamelLower      | anyKindOfString    | default
| Snake           | any_kind_of_string |
| SnakeScreaming  | ANY_KIND_OF_STRING |
| SnakeFirstUpper | rgb_code_md5       |
| Kebab           | any-kind-of-string |
| KebabScreaming  | ANY-KIND-OF-STRING |
`
	CGenDaoBriefTplDaoIndexPath    = `template file path for dao index file`
	CGenDaoBriefTplDaoInternalPath = `template file path for dao internal file`
	CGenDaoBriefTplDaoDoPathPath   = `template file path for dao do file`
	CGenDaoBriefTplDaoEntityPath   = `template file path for dao entity file`

	tplVarTableName               = `{TplTableName}`
	tplVarTableNameCamelCase      = `{TplTableNameCamelCase}`
	tplVarTableNameCamelLowerCase = `{TplTableNameCamelLowerCase}`
	tplVarPackageImports          = `{TplPackageImports}`
	tplVarImportPrefix            = `{TplImportPrefix}`
	tplVarStructDefine            = `{TplStructDefine}`
	tplVarColumnDefine            = `{TplColumnDefine}`
	tplVarColumnNames             = `{TplColumnNames}`
	tplVarGroupName               = `{TplGroupName}`
	tplVarDatetimeStr             = `{TplDatetimeStr}`
	tplVarCreatedAtDatetimeStr    = `{TplCreatedAtDatetimeStr}`
)

var (
	createdAt          = gtime.Now()
	defaultTypeMapping = map[DBFieldTypeName]CustomAttributeType{
		"decimal": {
			Type: "float64",
		},
		"money": {
			Type: "float64",
		},
		"numeric": {
			Type: "float64",
		},
		"smallmoney": {
			Type: "float64",
		},
	}
)

func init() {
	gtag.Sets(g.MapStrStr{
		`CGenDaoConfig`:                  CGenDaoConfig,
		`CGenDaoUsage`:                   CGenDaoUsage,
		`CGenDaoBrief`:                   CGenDaoBrief,
		`CGenDaoEg`:                      CGenDaoEg,
		`CGenDaoAd`:                      CGenDaoAd,
		`CGenDaoBriefPath`:               CGenDaoBriefPath,
		`CGenDaoBriefLink`:               CGenDaoBriefLink,
		`CGenDaoBriefTables`:             CGenDaoBriefTables,
		`CGenDaoBriefTablesEx`:           CGenDaoBriefTablesEx,
		`CGenDaoBriefPrefix`:             CGenDaoBriefPrefix,
		`CGenDaoBriefRemovePrefix`:       CGenDaoBriefRemovePrefix,
		`CGenDaoBriefRemoveFieldPrefix`:  CGenDaoBriefRemoveFieldPrefix,
		`CGenDaoBriefStdTime`:            CGenDaoBriefStdTime,
		`CGenDaoBriefWithTime`:           CGenDaoBriefWithTime,
		`CGenDaoBriefDaoPath`:            CGenDaoBriefDaoPath,
		`CGenDaoBriefDoPath`:             CGenDaoBriefDoPath,
		`CGenDaoBriefEntityPath`:         CGenDaoBriefEntityPath,
		`CGenDaoBriefGJsonSupport`:       CGenDaoBriefGJsonSupport,
		`CGenDaoBriefImportPrefix`:       CGenDaoBriefImportPrefix,
		`CGenDaoBriefOverwriteDao`:       CGenDaoBriefOverwriteDao,
		`CGenDaoBriefModelFile`:          CGenDaoBriefModelFile,
		`CGenDaoBriefModelFileForDao`:    CGenDaoBriefModelFileForDao,
		`CGenDaoBriefDescriptionTag`:     CGenDaoBriefDescriptionTag,
		`CGenDaoBriefNoJsonTag`:          CGenDaoBriefNoJsonTag,
		`CGenDaoBriefNoModelComment`:     CGenDaoBriefNoModelComment,
		`CGenDaoBriefClear`:              CGenDaoBriefClear,
		`CGenDaoBriefTypeMapping`:        CGenDaoBriefTypeMapping,
		`CGenDaoBriefFieldMapping`:       CGenDaoBriefFieldMapping,
		`CGenDaoBriefGroup`:              CGenDaoBriefGroup,
		`CGenDaoBriefJsonCase`:           CGenDaoBriefJsonCase,
		`CGenDaoBriefTplDaoIndexPath`:    CGenDaoBriefTplDaoIndexPath,
		`CGenDaoBriefTplDaoInternalPath`: CGenDaoBriefTplDaoInternalPath,
		`CGenDaoBriefTplDaoDoPathPath`:   CGenDaoBriefTplDaoDoPathPath,
		`CGenDaoBriefTplDaoEntityPath`:   CGenDaoBriefTplDaoEntityPath,
	})
}

type (
	CGenDao      struct{}
	CGenDaoInput struct {
		g.Meta             `name:"dao" config:"{CGenDaoConfig}" usage:"{CGenDaoUsage}" brief:"{CGenDaoBrief}" eg:"{CGenDaoEg}" ad:"{CGenDaoAd}"`
		Path               string `name:"path"                short:"p"  brief:"{CGenDaoBriefPath}" d:"internal"`
		Link               string `name:"link"                short:"l"  brief:"{CGenDaoBriefLink}"`
		Tables             string `name:"tables"              short:"t"  brief:"{CGenDaoBriefTables}"`
		TablesEx           string `name:"tablesEx"            short:"x"  brief:"{CGenDaoBriefTablesEx}"`
		Group              string `name:"group"               short:"g"  brief:"{CGenDaoBriefGroup}" d:"default"`
		Prefix             string `name:"prefix"              short:"f"  brief:"{CGenDaoBriefPrefix}"`
		RemovePrefix       string `name:"removePrefix"        short:"r"  brief:"{CGenDaoBriefRemovePrefix}"`
		RemoveFieldPrefix  string `name:"removeFieldPrefix"   short:"rf" brief:"{CGenDaoBriefRemoveFieldPrefix}"`
		JsonCase           string `name:"jsonCase"            short:"j"  brief:"{CGenDaoBriefJsonCase}" d:"CamelLower"`
		ImportPrefix       string `name:"importPrefix"        short:"i"  brief:"{CGenDaoBriefImportPrefix}"`
		DaoPath            string `name:"daoPath"             short:"d"  brief:"{CGenDaoBriefDaoPath}" d:"dao"`
		DoPath             string `name:"doPath"              short:"o"  brief:"{CGenDaoBriefDoPath}" d:"model/do"`
		EntityPath         string `name:"entityPath"          short:"e"  brief:"{CGenDaoBriefEntityPath}" d:"model/entity"`
		TplDaoIndexPath    string `name:"tplDaoIndexPath"     short:"t1" brief:"{CGenDaoBriefTplDaoIndexPath}"`
		TplDaoInternalPath string `name:"tplDaoInternalPath"  short:"t2" brief:"{CGenDaoBriefTplDaoInternalPath}"`
		TplDaoDoPath       string `name:"tplDaoDoPath"        short:"t3" brief:"{CGenDaoBriefTplDaoDoPathPath}"`
		TplDaoEntityPath   string `name:"tplDaoEntityPath"    short:"t4" brief:"{CGenDaoBriefTplDaoEntityPath}"`
		StdTime            bool   `name:"stdTime"             short:"s"  brief:"{CGenDaoBriefStdTime}" orphan:"true"`
		WithTime           bool   `name:"withTime"            short:"w"  brief:"{CGenDaoBriefWithTime}" orphan:"true"`
		GJsonSupport       bool   `name:"gJsonSupport"        short:"n"  brief:"{CGenDaoBriefGJsonSupport}" orphan:"true"`
		OverwriteDao       bool   `name:"overwriteDao"        short:"v"  brief:"{CGenDaoBriefOverwriteDao}" orphan:"true"`
		DescriptionTag     bool   `name:"descriptionTag"      short:"c"  brief:"{CGenDaoBriefDescriptionTag}" orphan:"true"`
		NoJsonTag          bool   `name:"noJsonTag"           short:"k"  brief:"{CGenDaoBriefNoJsonTag}" orphan:"true"`
		NoModelComment     bool   `name:"noModelComment"      short:"m"  brief:"{CGenDaoBriefNoModelComment}" orphan:"true"`
		Clear              bool   `name:"clear"               short:"a"  brief:"{CGenDaoBriefClear}" orphan:"true"`

		TypeMapping  map[DBFieldTypeName]CustomAttributeType  `name:"typeMapping" short:"y" brief:"{CGenDaoBriefTypeMapping}" orphan:"true"`
		FieldMapping map[DBTableFieldName]CustomAttributeType `name:"fieldMapping" short:"fm"  brief:"{CGenDaoBriefFieldMapping}" orphan:"true"`
		genItems     *CGenDaoInternalGenItems
	}
	CGenDaoOutput struct{}

	CGenDaoInternalInput struct {
		CGenDaoInput
		DB            gdb.DB
		TableNames    []string
		NewTableNames []string
	}
	DBTableFieldName    = string
	DBFieldTypeName     = string
	CustomAttributeType struct {
		Type   string `brief:"custom attribute type name"`
		Import string `brief:"custom import for this type"`
	}
)

func (c CGenDao) Dao(ctx context.Context, in CGenDaoInput) (out *CGenDaoOutput, err error) {
	in.genItems = newCGenDaoInternalGenItems()
	if in.Link != "" {
		doGenDaoForArray(ctx, -1, in)
	} else if g.Cfg().Available(ctx) {
		v := g.Cfg().MustGet(ctx, CGenDaoConfig)
		if v.IsSlice() {
			for i := 0; i < len(v.Interfaces()); i++ {
				doGenDaoForArray(ctx, i, in)
			}
		} else {
			doGenDaoForArray(ctx, -1, in)
		}
	} else {
		doGenDaoForArray(ctx, -1, in)
	}
	doClear(in.genItems)
	mlog.Print("done!")
	return
}

// doGenDaoForArray implements the "gen dao" command for configuration array.
func doGenDaoForArray(ctx context.Context, index int, in CGenDaoInput) {
	var (
		err error
		db  gdb.DB
	)
	if index >= 0 {
		err = g.Cfg().MustGet(
			ctx,
			fmt.Sprintf(`%s.%d`, CGenDaoConfig, index),
		).Scan(&in)
		if err != nil {
			mlog.Fatalf(`invalid configuration of "%s": %+v`, CGenDaoConfig, err)
		}
	}
	if dirRealPath := gfile.RealPath(in.Path); dirRealPath == "" {
		mlog.Fatalf(`path "%s" does not exist`, in.Path)
	}
	removePrefixArray := gstr.SplitAndTrim(in.RemovePrefix, ",")

	// It uses user passed database configuration.
	if in.Link != "" {
		var tempGroup = gtime.TimestampNanoStr()
		gdb.AddConfigNode(tempGroup, gdb.ConfigNode{
			Link: in.Link,
		})
		if db, err = gdb.Instance(tempGroup); err != nil {
			mlog.Fatalf(`database initialization failed: %+v`, err)
		}
	} else {
		db = g.DB(in.Group)
	}
	if db == nil {
		mlog.Fatal(`database initialization failed, may be invalid database configuration`)
	}

	var tableNames []string
	if in.Tables != "" {
		tableNames = gstr.SplitAndTrim(in.Tables, ",")
	} else {
		tableNames, err = db.Tables(context.TODO())
		if err != nil {
			mlog.Fatalf("fetching tables failed: %+v", err)
		}
	}
	// Table excluding.
	if in.TablesEx != "" {
		array := garray.NewStrArrayFrom(tableNames)
		for _, v := range gstr.SplitAndTrim(in.TablesEx, ",") {
			array.RemoveValue(v)
		}
		tableNames = array.Slice()
	}

	// merge default typeMapping to input typeMapping.
	if in.TypeMapping == nil {
		in.TypeMapping = defaultTypeMapping
	} else {
		for key, typeMapping := range defaultTypeMapping {
			if _, ok := in.TypeMapping[key]; !ok {
				in.TypeMapping[key] = typeMapping
			}
		}
	}

	// Generating dao & model go files one by one according to given table name.
	newTableNames := make([]string, len(tableNames))
	for i, tableName := range tableNames {
		newTableName := tableName
		for _, v := range removePrefixArray {
			newTableName = gstr.TrimLeftStr(newTableName, v, 1)
		}
		newTableName = in.Prefix + newTableName
		newTableNames[i] = newTableName
	}

	in.genItems.Scale()

	// Dao: index and internal.
	generateDao(ctx, CGenDaoInternalInput{
		CGenDaoInput:  in,
		DB:            db,
		TableNames:    tableNames,
		NewTableNames: newTableNames,
	})
	// Do.
	generateDo(ctx, CGenDaoInternalInput{
		CGenDaoInput:  in,
		DB:            db,
		TableNames:    tableNames,
		NewTableNames: newTableNames,
	})
	// Entity.
	generateEntity(ctx, CGenDaoInternalInput{
		CGenDaoInput:  in,
		DB:            db,
		TableNames:    tableNames,
		NewTableNames: newTableNames,
	})

	in.genItems.SetClear(in.Clear)
}

func getImportPartContent(ctx context.Context, source string, isDo bool, appendImports []string) string {
	var packageImportsArray = garray.NewStrArray()
	if isDo {
		packageImportsArray.Append(`"github.com/gogf/gf/v2/frame/g"`)
	}

	// Time package recognition.
	if strings.Contains(source, "gtime.Time") {
		packageImportsArray.Append(`"github.com/gogf/gf/v2/os/gtime"`)
	} else if strings.Contains(source, "time.Time") {
		packageImportsArray.Append(`"time"`)
	}

	// Json type.
	if strings.Contains(source, "gjson.Json") {
		packageImportsArray.Append(`"github.com/gogf/gf/v2/encoding/gjson"`)
	}

	// Check and update imports in go.mod
	if len(appendImports) > 0 {
		goModPath := utils.GetModPath()
		if goModPath == "" {
			mlog.Fatal("go.mod not found in current project")
		}
		mod, err := modfile.Parse(goModPath, gfile.GetBytes(goModPath), nil)
		if err != nil {
			mlog.Fatalf("parse go.mod failed: %+v", err)
		}
		for _, appendImport := range appendImports {
			found := false
			for _, require := range mod.Require {
				if gstr.Contains(appendImport, require.Mod.Path) {
					found = true
					break
				}
			}
			if !found {
				if err = gproc.ShellRun(ctx, `go get `+appendImport); err != nil {
					mlog.Fatalf(`%+v`, err)
				}
			}
			packageImportsArray.Append(fmt.Sprintf(`"%s"`, appendImport))
		}
	}

	// Generate and write content to golang file.
	packageImportsStr := ""
	if packageImportsArray.Len() > 0 {
		packageImportsStr = fmt.Sprintf("import(\n%s\n)", packageImportsArray.Join("\n"))
	}
	return packageImportsStr
}

func replaceDefaultVar(in CGenDaoInternalInput, origin string) string {
	var tplCreatedAtDatetimeStr string
	var tplDatetimeStr string = createdAt.String()
	if in.WithTime {
		tplCreatedAtDatetimeStr = fmt.Sprintf(`Created at %s`, tplDatetimeStr)
	}
	return gstr.ReplaceByMap(origin, g.MapStrStr{
		tplVarDatetimeStr:          tplDatetimeStr,
		tplVarCreatedAtDatetimeStr: tplCreatedAtDatetimeStr,
	})
}

func sortFieldKeyForDao(fieldMap map[string]*gdb.TableField) []string {
	names := make(map[int]string)
	for _, field := range fieldMap {
		names[field.Index] = field.Name
	}
	var (
		i      = 0
		j      = 0
		result = make([]string, len(names))
	)
	for {
		if len(names) == 0 {
			break
		}
		if val, ok := names[i]; ok {
			result[j] = val
			j++
			delete(names, i)
		}
		i++
	}
	return result
}

func getTemplateFromPathOrDefault(filePath string, def string) string {
	if filePath != "" {
		if contents := gfile.GetContents(filePath); contents != "" {
			return contents
		}
	}
	return def
}
