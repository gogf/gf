package tpl

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/utils"
	_ "github.com/gogf/gf/contrib/drivers/clickhouse/v2"
	_ "github.com/gogf/gf/contrib/drivers/mssql/v2"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/drivers/oracle/v2"
	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/os/gview"
	"github.com/gogf/gf/v2/util/gtag"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
)

const (
	CGenTplConfig = `gfcli.gen.tpl`
	CGenTplUsage  = `gf gen tpl [OPTION]`
	CGenTplBrief  = `automatically generate template files`
	CGenTplEg     = `
gf gen tpl
gf gen tpl -t default -p ./template
`
	CGenTplAd = `
CONFIGURATION SUPPORT
    Options are also supported by configuration file.
    It's suggested using configuration file instead of command line arguments.
    The configuration node name is "gfcli.gen.tpl" which also supports multiple databases, for example(config.yaml):
	gfcli:
	  gen:
		tpl:
		- link:     "mysql:root:12345678@tcp(127.0.0.1:3306)/test"
		  tables:   "order,products"
		  jsonCase: "CamelLower"
		- link:     "mysql:root:12345678@tcp(127.0.0.1:3306)/primary"
		  path:     "./my-app"
		  prefix:   "primary_"
		  tables:   "user, userDetail"
`

	CGenTplBriefPath = `output directory path (default: "./template")`
)

func init() {
	gtag.Sets(g.MapStrStr{
		`CGenTplConfig`: CGenTplConfig,
		`CGenTplUsage`:  CGenTplUsage,
		`CGenTplBrief`:  CGenTplBrief,
		`CGenTplEg`:     CGenTplEg,
		`CGenTplAd`:     CGenTplAd,
	})
}

type (
	CGenTpl      struct{}
	CGenTplInput struct {
		g.Meta            `name:"tpl" config:"{CGenTplConfig}" usage:"{CGenTplUsage}" brief:"{CGenTplBrief}" eg:"{CGenTplEg}" ad:"{CGenTplAd}"`
		Path              string `name:"path"                short:"p"  brief:"{CGenTplBriefPath}" d:"./output"`
		TplPath           string `name:"tplPath"             short:"tp" brief:"模板目录路径" d:"./template"`
		Link              string `name:"link"                short:"l"  brief:"{CGenDaoBriefLink}"`
		Tables            string `name:"tables"              short:"t"  brief:"{CGenDaoBriefTables}"`
		TablesEx          string `name:"tablesEx"            short:"x"  brief:"{CGenDaoBriefTablesEx}"`
		Group             string `name:"group"               short:"g"  brief:"{CGenDaoBriefGroup}" d:"default"`
		Prefix            string `name:"prefix"              short:"f"  brief:"{CGenDaoBriefPrefix}"`
		RemovePrefix      string `name:"removePrefix"        short:"r"  brief:"{CGenDaoBriefRemovePrefix}"`
		RemoveFieldPrefix string `name:"removeFieldPrefix"   short:"rf" brief:"{CGenDaoBriefRemoveFieldPrefix}"`
		JsonCase          string `name:"jsonCase"            short:"j"  brief:"{CGenDaoBriefJsonCase}" d:"CamelLower"`
		ImportPrefix      string `name:"importPrefix"        short:"i"  brief:"{CGenDaoBriefImportPrefix}"`
		// DaoPath            string                         `name:"daoPath"             short:"d"  brief:"{CGenDaoBriefDaoPath}" d:"dao"`
		// DoPath             string                         `name:"doPath"              short:"o"  brief:"{CGenDaoBriefDoPath}" d:"model/do"`
		// EntityPath         string                         `name:"entityPath"          short:"e"  brief:"{CGenDaoBriefEntityPath}" d:"model/entity"`
		// TplDaoIndexPath    string                         `name:"tplDaoIndexPath"     short:"t1" brief:"{CGenDaoBriefTplDaoIndexPath}"`
		// TplDaoInternalPath string                         `name:"tplDaoInternalPath"  short:"t2" brief:"{CGenDaoBriefTplDaoInternalPath}"`
		// TplDaoDoPath       string                         `name:"tplDaoDoPath"        short:"t3" brief:"{CGenDaoBriefTplDaoDoPathPath}"`
		// TplDaoEntityPath   string                         `name:"tplDaoEntityPath"    short:"t4" brief:"{CGenDaoBriefTplDaoEntityPath}"`
		StdTime        bool                           `name:"stdTime"             short:"s"  brief:"{CGenDaoBriefStdTime}" orphan:"true"`
		WithTime       bool                           `name:"withTime"            short:"w"  brief:"{CGenDaoBriefWithTime}" orphan:"true"`
		GJsonSupport   bool                           `name:"gJsonSupport"        short:"n"  brief:"{CGenDaoBriefGJsonSupport}" orphan:"true"`
		OverwriteDao   bool                           `name:"overwriteDao"        short:"v"  brief:"{CGenDaoBriefOverwriteDao}" orphan:"true"`
		DescriptionTag bool                           `name:"descriptionTag"      short:"c"  brief:"{CGenDaoBriefDescriptionTag}" orphan:"true"`
		NoJsonTag      bool                           `name:"noJsonTag"           short:"k"  brief:"{CGenDaoBriefNoJsonTag}" orphan:"true"`
		NoModelComment bool                           `name:"noModelComment"      short:"m"  brief:"{CGenDaoBriefNoModelComment}" orphan:"true"`
		Clear          bool                           `name:"clear"               short:"a"  brief:"{CGenDaoBriefClear}" orphan:"true"`
		TypeMapping    map[string]CustomAttributeType `name:"typeMapping"         short:"y"  brief:"{CGenDaoBriefTypeMapping}"  orphan:"true"`
		FieldMapping   map[string]CustomAttributeType `name:"fieldMapping"        short:"fm" brief:"{CGenDaoBriefFieldMapping}" orphan:"true"`
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

// GetTables 获取数据库表结构信息
func GetTables(ctx context.Context, db gdb.DB) Tables {
	tablesName, err := db.Tables(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(tablesName)
	tables := make(Tables, 0)
	for _, v := range tablesName {
		t, err := NewTable(ctx, db, v)
		if err != nil {
			panic(err)
		}
		t.SortFields(true)
		tables = append(tables, t)
	}
	return tables
}

func (c CGenTpl) Tpl(ctx context.Context, in CGenTplInput) (out *CGenTplOutput, err error) {
	// Clear old files
	// if in.Clear {
	// 	if err := gfile.Remove(in.Path); err != nil {
	// 		return nil, gerror.Wrapf(err, "clear output path failed")
	// 	}
	// }

	// Create output directory
	// if !gfile.Exists(in.Path) {
	// 	if err := gfile.Mkdir(in.Path); err != nil {
	// 		return nil, gerror.Wrapf(err, "create output directory failed")
	// 	}
	// }
	db, err := in.GetDB()
	if err != nil {
		return nil, err
	}
	outputDir := in.Path
	tplRootDir := in.TplPath
	tplRootDir = gfile.Abs(tplRootDir)
	fmt.Println(tplRootDir)
	tplList, err := gfile.ScanDirFile(tplRootDir, "*.tpl", true)
	if err != nil {
		panic(err)
	}
	fmt.Println(tplList)

	fmt.Printf("%#v\n", Table{})
	fmt.Printf("%#v\n", TableField{})

	tables := GetTables(ctx, db)
	view := gview.New()

	for _, table := range tables {
		table.PackageName = "github.com/gogf/gf/cmd/gf/v2"
		tplData := g.Map{
			"table":  table,
			"tables": tables,
		}
		fmt.Println(table.FieldsJsonStr("Snake"))
		for _, tpl := range tplList {
			tplDir := gfile.Dir(tpl)

			res, err := view.Parse(ctx, tpl, tplData)
			if err != nil {
				panic(err)
			}
			// fmt.Println(res, err)
			filePath := filepath.FromSlash(fmt.Sprintf(outputDir+"%s/%s.go", strings.TrimPrefix(tplDir, tplRootDir), table.Name))
			err = gfile.PutContents(filePath, res)
			if err != nil {
				panic(err)
			}
			utils.GoFmt(filePath)
		}

	}

	mlog.Print("template files generated successfully!")
	return &CGenTplOutput{}, nil
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
