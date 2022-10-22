package gendao

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gtag"
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
`
	CGenDaoBriefPath            = `directory path for generated files`
	CGenDaoBriefLink            = `database configuration, the same as the ORM configuration of GoFrame`
	CGenDaoBriefTables          = `generate models only for given tables, multiple table names separated with ','`
	CGenDaoBriefTablesEx        = `generate models excluding given tables, multiple table names separated with ','`
	CGenDaoBriefPrefix          = `add prefix for all table of specified link/database tables`
	CGenDaoBriefRemovePrefix    = `remove specified prefix of the table, multiple prefix separated with ','`
	CGenDaoBriefStdTime         = `use time.Time from stdlib instead of gtime.Time for generated time/date fields of tables`
	CGenDaoBriefWithTime        = `add created time for auto produced go files`
	CGenDaoBriefGJsonSupport    = `use gJsonSupport to use *gjson.Json instead of string for generated json fields of tables`
	CGenDaoBriefImportPrefix    = `custom import prefix for generated go files`
	CGenDaoBriefDaoPath         = `directory path for storing generated dao files under path`
	CGenDaoBriefDoPath          = `directory path for storing generated do files under path`
	CGenDaoBriefEntityPath      = `directory path for storing generated entity files under path`
	CGenDaoBriefOverwriteDao    = `overwrite all dao files both inside/outside internal folder`
	CGenDaoBriefModelFile       = `custom file name for storing generated model content`
	CGenDaoBriefModelFileForDao = `custom file name generating model for DAO operations like Where/Data. It's empty in default`
	CGenDaoBriefDescriptionTag  = `add comment to description tag for each field`
	CGenDaoBriefNoJsonTag       = `no json tag will be added for each field`
	CGenDaoBriefNoModelComment  = `no model comment will be added for each field`
	CGenDaoBriefClear           = `delete all generated go files that do not exist in database`
	CGenDaoBriefGroup           = `
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
	createdAt = gtime.Now()
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
	}
	CGenDaoOutput struct{}

	CGenDaoInternalInput struct {
		CGenDaoInput
		DB            gdb.DB
		TableNames    []string
		NewTableNames []string
		ModName       string // Module name of current golang project, which is used for import purpose.
	}
)

func (c CGenDao) Dao(ctx context.Context, in CGenDaoInput) (out *CGenDaoOutput, err error) {
	if g.Cfg().Available(ctx) {
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
	mlog.Print("done!")
	return
}

// doGenDaoForArray implements the "gen dao" command for configuration array.
func doGenDaoForArray(ctx context.Context, index int, in CGenDaoInput) {
	var (
		err     error
		db      gdb.DB
		modName string // Go module name, eg: github.com/gogf/gf.
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
	if in.ImportPrefix == "" {
		if !gfile.Exists("go.mod") {
			mlog.Fatal("go.mod does not exist in current working directory")
		}
		var (
			goModContent = gfile.GetContents("go.mod")
			match, _     = gregex.MatchString(`^module\s+(.+)\s*`, goModContent)
		)
		if len(match) > 1 {
			modName = gstr.Trim(match[1])
		} else {
			mlog.Fatal("module name does not found in go.mod")
		}
	}

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
	// Dao: index and internal.
	generateDao(ctx, CGenDaoInternalInput{
		CGenDaoInput:  in,
		DB:            db,
		TableNames:    tableNames,
		NewTableNames: newTableNames,
		ModName:       modName,
	})
	// Do.
	generateDo(ctx, CGenDaoInternalInput{
		CGenDaoInput:  in,
		DB:            db,
		TableNames:    tableNames,
		NewTableNames: newTableNames,
		ModName:       modName,
	})
	// Entity.
	generateEntity(ctx, CGenDaoInternalInput{
		CGenDaoInput:  in,
		DB:            db,
		TableNames:    tableNames,
		NewTableNames: newTableNames,
		ModName:       modName,
	})
}

func getImportPartContent(source string, isDo bool) string {
	var (
		packageImportsArray = garray.NewStrArray()
	)

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

	// garray
	if strings.Contains(source, "garray.Array") {
		packageImportsArray.Append(`"github.com/gogf/gf/v2/container/garray"`)
	}

	// gmap
	if strings.Contains(source, "g.Map") {
		packageImportsArray.Append(`"github.com/gogf/gf/v2/frame/g"`)
	}

	// uuid
	if strings.Contains(source, "uuid.UUID") {
		packageImportsArray.Append(`"github.com/google/uuid"`)
	}

	// big int
	if strings.Contains(source, "big.Int") {
		packageImportsArray.Append(`"math/big"`)
	}

	// decimal
	if strings.Contains(source, "decimal.Decimal") {
		packageImportsArray.Append(`"github.com/shopspring/decimal"`)
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
