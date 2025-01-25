package tpl

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	_ "github.com/gogf/gf/contrib/drivers/clickhouse/v2"
	_ "github.com/gogf/gf/contrib/drivers/mssql/v2"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/drivers/oracle/v2"
	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/os/gview"
	"github.com/gogf/gf/v2/util/gtag"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/utils"
)

const (
	CGenTplConfig = `gfcli.gen.tpl`
	CGenTplUsage  = `gf gen tpl [OPTION]`
	CGenTplBrief  = `automatically generate template files`
	CGenTplEg     = `
gf gen tpl
gf gen tpl -l "mysql:root:12345678@tcp(127.0.0.1:3306)/test"
gf gen tpl -p ./model -g user-center -t user,user_detail,user_login
gf gen tpl -r user_
`

	CGenTplAd = `
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
	CGenTplBriefPath              = `directory path for generated files`
	CGenTplBriefLink              = `database configuration, the same as the ORM configuration of GoFrame`
	CGenTplBriefTables            = `generate models only for given tables, multiple table names separated with ','`
	CGenTplBriefTablesEx          = `generate models excluding given tables, multiple table names separated with ','`
	CGenTplBriefPrefix            = `add prefix for all table of specified link/database tables`
	CGenTplBriefRemovePrefix      = `remove specified prefix of the table, multiple prefix separated with ','`
	CGenTplBriefRemoveFieldPrefix = `remove specified prefix of the field, multiple prefix separated with ','`
	CGenTplBriefStdTime           = `use time.Time from stdlib instead of gtime.Time for generated time/date fields of tables`
	CGenTplBriefWithTime          = `add created time for auto produced go files`
	CGenTplBriefGJsonSupport      = `use gJsonSupport to use *gjson.Json instead of string for generated json fields of tables`
	CGenTplBriefImportPrefix      = `custom import prefix for generated go files`
	CGenTplBriefDaoPath           = `directory path for storing generated dao files under path`
	CGenTplBriefDoPath            = `directory path for storing generated do files under path`
	CGenTplBriefEntityPath        = `directory path for storing generated entity files under path`
	CGenTplBriefOverwriteDao      = `overwrite all dao files both inside/outside internal folder`
	CGenTplBriefModelFile         = `custom file name for storing generated model content`
	CGenTplBriefModelFileForDao   = `custom file name generating model for DAO operations like Where/Data. It's empty in default`
	CGenTplBriefDescriptionTag    = `add comment to description tag for each field`
	CGenTplBriefNoJsonTag         = `no json tag will be added for each field`
	CGenTplBriefNoModelComment    = `no model comment will be added for each field`
	CGenTplBriefClear             = `delete all generated go files that do not exist in database`
	CGenTplBriefTypeMapping       = `custom local type mapping for generated struct attributes relevant to fields of table`
	CGenTplBriefFieldMapping      = `custom local type mapping for generated struct attributes relevant to specific fields of table`
	CGenTplBriefGroup             = `
specifying the configuration group name of database for generated ORM instance,
it's not necessary and the default value is "default"
`
	CGenTplBriefJsonCase = `
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
	CGenTplBriefTplDaoIndexPath    = `template file path for dao index file`
	CGenTplBriefTplDaoInternalPath = `template file path for dao internal file`
	CGenTplBriefTplDaoDoPathPath   = `template file path for dao do file`
	CGenTplBriefTplDaoEntityPath   = `template file path for dao entity file`
)

func init() {
	gtag.Sets(g.MapStrStr{
		`CGenTplConfig`:                  CGenTplConfig,
		`CGenTplUsage`:                   CGenTplUsage,
		`CGenTplBrief`:                   CGenTplBrief,
		`CGenTplEg`:                      CGenTplEg,
		`CGenTplAd`:                      CGenTplAd,
		`CGenTplBriefPath`:               CGenTplBriefPath,
		`CGenTplBriefLink`:               CGenTplBriefLink,
		`CGenTplBriefTables`:             CGenTplBriefTables,
		`CGenTplBriefTablesEx`:           CGenTplBriefTablesEx,
		`CGenTplBriefPrefix`:             CGenTplBriefPrefix,
		`CGenTplBriefRemovePrefix`:       CGenTplBriefRemovePrefix,
		`CGenTplBriefRemoveFieldPrefix`:  CGenTplBriefRemoveFieldPrefix,
		`CGenTplBriefStdTime`:            CGenTplBriefStdTime,
		`CGenTplBriefWithTime`:           CGenTplBriefWithTime,
		`CGenTplBriefDaoPath`:            CGenTplBriefDaoPath,
		`CGenTplBriefDoPath`:             CGenTplBriefDoPath,
		`CGenTplBriefEntityPath`:         CGenTplBriefEntityPath,
		`CGenTplBriefGJsonSupport`:       CGenTplBriefGJsonSupport,
		`CGenTplBriefImportPrefix`:       CGenTplBriefImportPrefix,
		`CGenTplBriefOverwriteDao`:       CGenTplBriefOverwriteDao,
		`CGenTplBriefModelFile`:          CGenTplBriefModelFile,
		`CGenTplBriefModelFileForDao`:    CGenTplBriefModelFileForDao,
		`CGenTplBriefDescriptionTag`:     CGenTplBriefDescriptionTag,
		`CGenTplBriefNoJsonTag`:          CGenTplBriefNoJsonTag,
		`CGenTplBriefNoModelComment`:     CGenTplBriefNoModelComment,
		`CGenTplBriefClear`:              CGenTplBriefClear,
		`CGenTplBriefTypeMapping`:        CGenTplBriefTypeMapping,
		`CGenTplBriefFieldMapping`:       CGenTplBriefFieldMapping,
		`CGenTplBriefGroup`:              CGenTplBriefGroup,
		`CGenTplBriefJsonCase`:           CGenTplBriefJsonCase,
		`CGenTplBriefTplDaoIndexPath`:    CGenTplBriefTplDaoIndexPath,
		`CGenTplBriefTplDaoInternalPath`: CGenTplBriefTplDaoInternalPath,
		`CGenTplBriefTplDaoDoPathPath`:   CGenTplBriefTplDaoDoPathPath,
		`CGenTplBriefTplDaoEntityPath`:   CGenTplBriefTplDaoEntityPath,
	})
}

type (
	CGenTpl      struct{}
	CGenTplInput struct {
		g.Meta            `name:"tpl" config:"{CGenTplConfig}" usage:"{CGenTplUsage}" brief:"{CGenTplBrief}" eg:"{CGenTplEg}" ad:"{CGenTplAd}"`
		Path              string `name:"path"                short:"p"  brief:"{CGenTplBriefPath}" d:"./output"`
		TplPath           string `name:"tplPath"             short:"tp" brief:"模板目录路径"`
		Link              string `name:"link"                short:"l"  brief:"{CGenTplBriefLink}"`
		Tables            string `name:"tables"              short:"t"  brief:"{CGenTplBriefTables}"`
		TablesEx          string `name:"tablesEx"            short:"x"  brief:"{CGenTplBriefTablesEx}"`
		Group             string `name:"group"               short:"g"  brief:"{CGenTplBriefGroup}" d:"default"`
		Prefix            string `name:"prefix"              short:"f"  brief:"{CGenTplBriefPrefix}"`
		RemovePrefix      string `name:"removePrefix"        short:"r"  brief:"{CGenTplBriefRemovePrefix}"`
		RemoveFieldPrefix string `name:"removeFieldPrefix"   short:"rf" brief:"{CGenTplBriefRemoveFieldPrefix}"`
		JsonCase          string `name:"jsonCase"            short:"j"  brief:"{CGenTplBriefJsonCase}" d:"CamelLower"`
		ImportPrefix      string `name:"importPrefix"        short:"i"  brief:"{CGenTplBriefImportPrefix}"`
		// 新增过滤参数
		TableNamePattern string `name:"tableNamePattern"    short:"tn" brief:"表名匹配模式，支持通配符"`
		// DaoPath            string                         `name:"daoPath"             short:"d"  brief:"{CGenTplBriefDaoPath}" d:"dao"`
		// DoPath             string                         `name:"doPath"              short:"o"  brief:"{CGenTplBriefDoPath}" d:"model/do"`
		// EntityPath         string                         `name:"entityPath"          short:"e"  brief:"{CGenTplBriefEntityPath}" d:"model/entity"`
		// TplDaoIndexPath    string                         `name:"tplDaoIndexPath"     short:"t1" brief:"{CGenTplBriefTplDaoIndexPath}"`
		// TplDaoInternalPath string                         `name:"tplDaoInternalPath"  short:"t2" brief:"{CGenTplBriefTplDaoInternalPath}"`
		// TplDaoDoPath       string                         `name:"tplDaoDoPath"        short:"t3" brief:"{CGenTplBriefTplDaoDoPathPath}"`
		// TplDaoEntityPath   string                         `name:"tplDaoEntityPath"    short:"t4" brief:"{CGenTplBriefTplDaoEntityPath}"`
		StdTime        bool                           `name:"stdTime"             short:"s"  brief:"{CGenTplBriefStdTime}" orphan:"true"`
		WithTime       bool                           `name:"withTime"            short:"w"  brief:"{CGenTplBriefWithTime}" orphan:"true"`
		GJsonSupport   bool                           `name:"gJsonSupport"        short:"n"  brief:"{CGenTplBriefGJsonSupport}" orphan:"true"`
		OverwriteDao   bool                           `name:"overwriteDao"        short:"v"  brief:"{CGenTplBriefOverwriteDao}" orphan:"true"`
		DescriptionTag bool                           `name:"descriptionTag"      short:"c"  brief:"{CGenTplBriefDescriptionTag}" orphan:"true"`
		NoJsonTag      bool                           `name:"noJsonTag"           short:"k"  brief:"{CGenTplBriefNoJsonTag}" orphan:"true"`
		NoModelComment bool                           `name:"noModelComment"      short:"m"  brief:"{CGenTplBriefNoModelComment}" orphan:"true"`
		Clear          bool                           `name:"clear"               short:"a"  brief:"{CGenTplBriefClear}" orphan:"true"`
		TypeMapping    map[string]CustomAttributeType `name:"typeMapping"         short:"y"  brief:"{CGenTplBriefTypeMapping}"  orphan:"true"`
		FieldMapping   map[string]CustomAttributeType `name:"fieldMapping"        short:"fm" brief:"{CGenTplBriefFieldMapping}" orphan:"true"`
	}
	CGenTplOutput struct{}

	CustomAttributeType struct {
		Type   string `brief:"custom attribute type name"`
		Import string `brief:"custom import for this type"`
	}
)

var (
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

type (
	DBFieldTypeName = string
)

// TplObj description
type TplObj struct {
	ctx        context.Context
	in         CGenTplInput
	db         gdb.DB
	TplPathAbs string
}

// NewTpl description
//
// createTime: 2025-01-25 16:36:43
func NewTpl(ctx context.Context, in CGenTplInput) (*TplObj, error) {
	db, err := in.GetDB()
	if err != nil {
		return nil, err
	}
	return &TplObj{
		ctx:        ctx,
		in:         in,
		db:         db,
		TplPathAbs: gfile.Abs(in.TplPath),
	}, nil
}

func (t *TplObj) ShowParams() {
	mlog.Debug("tplPath:", t.in.TplPath)
	mlog.Debug("output:", t.in.Path)
}
func (t *TplObj) Format() {
	utils.GoFmt(t.in.Path)
}

// GetTplFileList description
//
// createTime: 2025-01-25 16:43:06
func (t *TplObj) GetTplFileList() ([]string, error) {
	tplList, err := gfile.ScanDirFile(t.TplPathAbs, "*.tpl", true)
	if err != nil {
		return nil, err
	}
	return tplList, nil
}

func (c CGenTpl) Tpl(ctx context.Context, in CGenTplInput) (out *CGenTplOutput, err error) {
	if in.TplPath == "" {
		return nil, gerror.New("tplPath is required")
	}

	// Clear old files
	if in.Clear {
		if err := gfile.Remove(in.Path); err != nil {
			return nil, gerror.Wrapf(err, "clear output path failed")
		}
	}

	// Create output directory
	if !gfile.Exists(in.Path) {
		if err := gfile.Mkdir(in.Path); err != nil {
			return nil, gerror.Wrapf(err, "create output directory failed")
		}
	}

	tplObj, err := NewTpl(ctx, in)
	if err != nil {
		return nil, err
	}

	tplList, err := tplObj.GetTplFileList()
	if err != nil {
		panic(err)
	}
	fmt.Println(tplList)

	fmt.Printf("%#v\n", Table{})
	fmt.Printf("%#v\n", TableField{})

	tables, err := tplObj.GetTables()
	if err != nil {
		return nil, err
	}
	view := gview.New()

	for _, table := range tables {

		tplData := g.Map{
			"table":  table,
			"tables": tables,
		}
		fmt.Println(table.FieldsJsonStr(in.JsonCase))
		for _, tpl := range tplList {
			mlog.Print("generating template file:", tpl)
			// 相对路径
			relativePath := strings.TrimPrefix(gfile.Dir(tpl), tplObj.TplPathAbs)
			mlog.Print("relativePath:", relativePath)
			table.PackageName = filepath.ToSlash(filepath.Join(in.ImportPrefix, relativePath))
			filePath := filepath.Join(relativePath, table.FileName())
			mlog.Print("generating table filePath:", filePath)

			res, err := view.Parse(ctx, tpl, tplData)
			if err != nil {
				mlog.Fatal(err)
			}
			fmt.Println(len(res), err)

			err = tplObj.SaveFile(ctx, filePath, res)
			if err != nil {
				panic(err)
			}
		}

	}

	// Format generated files
	tplObj.Format()

	mlog.Print("template files generated successfully!")
	return &CGenTplOutput{}, nil
}

// SaveFile description
//
// createTime: 2025-01-25 17:05:25
func (t *TplObj) SaveFile(ctx context.Context, path, content string) error {
	mlog.Print("saving file:", path)
	path = filepath.Join(t.in.Path, path)
	mlog.Print("saving file:", path)
	path = filepath.FromSlash(path)
	mlog.Print("saving file:", path)
	if err := gfile.PutContents(path, content); err != nil {
		return err
	}
	return nil
}

// GetDB description
//
// createTime: 2025-01-24 16:58:46
func (in CGenTplInput) GetDB() (db gdb.DB, err error) {
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
	return
}
