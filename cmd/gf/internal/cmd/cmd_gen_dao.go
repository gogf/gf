package cmd

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/cmd/gf/v2/internal/consts"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/utils"
	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gtag"
	"github.com/olekukonko/tablewriter"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/lib/pq"
	//_ "github.com/mattn/go-oci8"
	//_ "github.com/mattn/go-sqlite3"
)

const (
	defaultDaoPath    = `service/internal/dao`
	defaultDoPath     = `service/internal/do`
	defaultEntityPath = `model/entity`
	cGenDaoConfig     = `gfcli.gen.dao`
	cGenDaoUsage      = `gf gen dao [OPTION]`
	cGenDaoBrief      = `automatically generate go files for dao/do/entity`
	cGenDaoEg         = `
gf gen dao
gf gen dao -l "mysql:root:12345678@tcp(127.0.0.1:3306)/test"
gf gen dao -p ./model -c config.yaml -g user-center -t user,user_detail,user_login
gf gen dao -r user_
`

	cGenDaoAd = `
CONFIGURATION SUPPORT
    Options are also supported by configuration file.
    It's suggested using configuration file instead of command line arguments making producing.
    The configuration node name is "gf.gen.dao", which also supports multiple databases, for example(config.yaml):
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
	cGenDaoBriefPath            = `directory path for generated files`
	cGenDaoBriefLink            = `database configuration, the same as the ORM configuration of GoFrame`
	cGenDaoBriefTables          = `generate models only for given tables, multiple table names separated with ','`
	cGenDaoBriefTablesEx        = `generate models excluding given tables, multiple table names separated with ','`
	cGenDaoBriefPrefix          = `add prefix for all table of specified link/database tables`
	cGenDaoBriefRemovePrefix    = `remove specified prefix of the table, multiple prefix separated with ','`
	cGenDaoBriefStdTime         = `use time.Time from stdlib instead of gtime.Time for generated time/date fields of tables`
	cGenDaoBriefGJsonSupport    = `use gJsonSupport to use *gjson.Json instead of string for generated json fields of tables`
	cGenDaoBriefImportPrefix    = `custom import prefix for generated go files`
	cGenDaoBriefOverwriteDao    = `overwrite all dao files both inside/outside internal folder`
	cGenDaoBriefModelFile       = `custom file name for storing generated model content`
	cGenDaoBriefModelFileForDao = `custom file name generating model for DAO operations like Where/Data. It's empty in default`
	cGenDaoBriefDescriptionTag  = `add comment to description tag for each field`
	cGenDaoBriefNoJsonTag       = `no json tag will be added for each field`
	cGenDaoBriefNoModelComment  = `no model comment will be added for each field`
	cGenDaoBriefGroup           = `
specifying the configuration group name of database for generated ORM instance,
it's not necessary and the default value is "default"
`
	cGenDaoBriefJsonCase = `
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

	tplVarTableName               = `{TplTableName}`
	tplVarTableNameCamelCase      = `{TplTableNameCamelCase}`
	tplVarTableNameCamelLowerCase = `{TplTableNameCamelLowerCase}`
	tplVarPackageImports          = `{TplPackageImports}`
	tplVarImportPrefix            = `{TplImportPrefix}`
	tplVarStructDefine            = `{TplStructDefine}`
	tplVarColumnDefine            = `{TplColumnDefine}`
	tplVarColumnNames             = `{TplColumnNames}`
	tplVarGroupName               = `{TplGroupName}`
	tplVarDatetime                = `{TplDatetime}`
)

var (
	createdAt *gtime.Time
)

func init() {
	gtag.Sets(g.MapStrStr{
		`cGenDaoConfig`:               cGenDaoConfig,
		`cGenDaoUsage`:                cGenDaoUsage,
		`cGenDaoBrief`:                cGenDaoBrief,
		`cGenDaoEg`:                   cGenDaoEg,
		`cGenDaoAd`:                   cGenDaoAd,
		`cGenDaoBriefPath`:            cGenDaoBriefPath,
		`cGenDaoBriefLink`:            cGenDaoBriefLink,
		`cGenDaoBriefTables`:          cGenDaoBriefTables,
		`cGenDaoBriefTablesEx`:        cGenDaoBriefTablesEx,
		`cGenDaoBriefPrefix`:          cGenDaoBriefPrefix,
		`cGenDaoBriefRemovePrefix`:    cGenDaoBriefRemovePrefix,
		`cGenDaoBriefStdTime`:         cGenDaoBriefStdTime,
		`cGenDaoBriefGJsonSupport`:    cGenDaoBriefGJsonSupport,
		`cGenDaoBriefImportPrefix`:    cGenDaoBriefImportPrefix,
		`cGenDaoBriefOverwriteDao`:    cGenDaoBriefOverwriteDao,
		`cGenDaoBriefModelFile`:       cGenDaoBriefModelFile,
		`cGenDaoBriefModelFileForDao`: cGenDaoBriefModelFileForDao,
		`cGenDaoBriefDescriptionTag`:  cGenDaoBriefDescriptionTag,
		`cGenDaoBriefNoJsonTag`:       cGenDaoBriefNoJsonTag,
		`cGenDaoBriefNoModelComment`:  cGenDaoBriefNoModelComment,
		`cGenDaoBriefGroup`:           cGenDaoBriefGroup,
		`cGenDaoBriefJsonCase`:        cGenDaoBriefJsonCase,
	})

	createdAt = gtime.Now()
}

type (
	cGenDaoInput struct {
		g.Meta         `name:"dao" config:"{cGenDaoConfig}" usage:"{cGenDaoUsage}" brief:"{cGenDaoBrief}" eg:"{cGenDaoEg}" ad:"{cGenDaoAd}"`
		Path           string `name:"path"            short:"p" brief:"{cGenDaoBriefPath}" d:"internal"`
		Link           string `name:"link"            short:"l" brief:"{cGenDaoBriefLink}"`
		Tables         string `name:"tables"          short:"t" brief:"{cGenDaoBriefTables}"`
		TablesEx       string `name:"tablesEx"        short:"e" brief:"{cGenDaoBriefTablesEx}"`
		Group          string `name:"group"           short:"g" brief:"{cGenDaoBriefGroup}" d:"default"`
		Prefix         string `name:"prefix"          short:"f" brief:"{cGenDaoBriefPrefix}"`
		RemovePrefix   string `name:"removePrefix"    short:"r" brief:"{cGenDaoBriefRemovePrefix}"`
		JsonCase       string `name:"jsonCase"        short:"j" brief:"{cGenDaoBriefJsonCase}" d:"CamelLower"`
		ImportPrefix   string `name:"importPrefix"    short:"i" brief:"{cGenDaoBriefImportPrefix}"`
		StdTime        bool   `name:"stdTime"         short:"s" brief:"{cGenDaoBriefStdTime}"         orphan:"true"`
		GJsonSupport   bool   `name:"gJsonSupport"    short:"n" brief:"{cGenDaoBriefGJsonSupport}"    orphan:"true"`
		OverwriteDao   bool   `name:"overwriteDao"    short:"o" brief:"{cGenDaoBriefOverwriteDao}"    orphan:"true"`
		DescriptionTag bool   `name:"descriptionTag"  short:"d" brief:"{cGenDaoBriefDescriptionTag}"  orphan:"true"`
		NoJsonTag      bool   `name:"noJsonTag"       short:"k" brief:"{cGenDaoBriefNoJsonTag"        orphan:"true"`
		NoModelComment bool   `name:"noModelComment"  short:"m" brief:"{cGenDaoBriefNoModelComment}"  orphan:"true"`
	}
	cGenDaoOutput struct{}

	cGenDaoInternalInput struct {
		cGenDaoInput
		TableName    string // TableName specifies the table name of the table.
		NewTableName string // NewTableName specifies the prefix-stripped name of the table.
		ModName      string // ModName specifies the module name of current golang project, which is used for import purpose.
	}
)

func (c cGen) Dao(ctx context.Context, in cGenDaoInput) (out *cGenDaoOutput, err error) {
	if g.Cfg().Available(ctx) {
		v := g.Cfg().MustGet(ctx, cGenDaoConfig)
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
func doGenDaoForArray(ctx context.Context, index int, in cGenDaoInput) {
	var (
		err     error
		db      gdb.DB
		modName string // Go module name, eg: github.com/gogf/gf.
	)
	if index >= 0 {
		err = g.Cfg().MustGet(
			ctx,
			fmt.Sprintf(`%s.%d`, cGenDaoConfig, index),
		).Scan(&in)
		if err != nil {
			mlog.Fatalf(`invalid configuration of "%s": %+v`, cGenDaoConfig, err)
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
		tempGroup := gtime.TimestampNanoStr()
		match, _ := gregex.MatchString(`([a-z]+):(.+)`, in.Link)
		if len(match) == 3 {
			gdb.AddConfigNode(tempGroup, gdb.ConfigNode{
				Type: gstr.Trim(match[1]),
				Link: gstr.Trim(match[2]),
			})
			db, _ = gdb.Instance(tempGroup)
		}
	} else {
		db = g.DB(in.Group)
	}
	if db == nil {
		mlog.Fatal("database initialization failed")
	}

	var tableNames []string
	if in.Tables != "" {
		tableNames = gstr.SplitAndTrim(in.Tables, ",")
	} else {
		tableNames, err = db.Tables(context.TODO())
		if err != nil {
			mlog.Fatalf("fetching tables failed: \n %v", err)
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
		// Dao.
		generateDao(ctx, db, cGenDaoInternalInput{
			cGenDaoInput: in,
			TableName:    tableName,
			NewTableName: newTableName,
			ModName:      modName,
		})
	}
	// Do.
	generateDo(ctx, db, tableNames, newTableNames, cGenDaoInternalInput{
		cGenDaoInput: in,
		ModName:      modName,
	})
	// Entity.
	generateEntity(ctx, db, tableNames, newTableNames, cGenDaoInternalInput{
		cGenDaoInput: in,
		ModName:      modName,
	})
}

// generateDaoContentFile generates the dao and model content of given table.
func generateDao(ctx context.Context, db gdb.DB, in cGenDaoInternalInput) {
	// Generating table data preparing.
	fieldMap, err := db.TableFields(ctx, in.TableName)
	if err != nil {
		mlog.Fatalf("fetching tables fields failed for table '%s':\n%v", in.TableName, err)
	}
	var (
		dirRealPath             = gfile.RealPath(in.Path)
		dirPathDao              = gfile.Join(in.Path, defaultDaoPath)
		tableNameCamelCase      = gstr.CaseCamel(in.NewTableName)
		tableNameCamelLowerCase = gstr.CaseCamelLower(in.NewTableName)
		tableNameSnakeCase      = gstr.CaseSnake(in.NewTableName)
		importPrefix            = in.ImportPrefix
	)
	if importPrefix == "" {
		if dirRealPath == "" {
			dirRealPath = in.Path
			importPrefix = dirRealPath
			importPrefix = gstr.Trim(dirRealPath, "./")
		} else {
			importPrefix = gstr.Replace(dirRealPath, gfile.Pwd(), "")
		}
		importPrefix = gstr.Replace(importPrefix, gfile.Separator, "/")
		importPrefix = gstr.Join(g.SliceStr{in.ModName, importPrefix, defaultDaoPath}, "/")
		importPrefix, _ = gregex.ReplaceString(`\/{2,}`, `/`, gstr.Trim(importPrefix, "/"))
	}

	fileName := gstr.Trim(tableNameSnakeCase, "-_.")
	if len(fileName) > 5 && fileName[len(fileName)-5:] == "_test" {
		// Add suffix to avoid the table name which contains "_test",
		// which would make the go file a testing file.
		fileName += "_table"
	}

	// dao - index
	generateDaoIndex(tableNameCamelCase, tableNameCamelLowerCase, importPrefix, dirPathDao, fileName, in)

	// dao - internal
	generateDaoInternal(tableNameCamelCase, tableNameCamelLowerCase, importPrefix, dirPathDao, fileName, fieldMap, in)
}

func generateDo(ctx context.Context, db gdb.DB, tableNames, newTableNames []string, in cGenDaoInternalInput) {
	var (
		doDirPath = gfile.Join(in.Path, defaultDoPath)
	)
	in.NoJsonTag = true
	in.DescriptionTag = false
	in.NoModelComment = false
	// Model content.
	for i, tableName := range tableNames {
		in.TableName = tableName
		fieldMap, err := db.TableFields(ctx, tableName)
		if err != nil {
			mlog.Fatalf("fetching tables fields failed for table '%s':\n%v", in.TableName, err)
		}
		var (
			newTableName     = newTableNames[i]
			doFilePath       = gfile.Join(doDirPath, gstr.CaseSnake(newTableName)+".go")
			structDefinition = generateStructDefinition(generateStructDefinitionInput{
				cGenDaoInternalInput: in,
				StructName:           gstr.CaseCamel(newTableName),
				FieldMap:             fieldMap,
				IsDo:                 true,
			})
		)
		// replace all types to interface{}.
		structDefinition, _ = gregex.ReplaceStringFuncMatch(
			"([A-Z]\\w*?)\\s+([\\w\\*\\.]+?)\\s+(//)",
			structDefinition,
			func(match []string) string {
				// If the type is already a pointer/slice/map, it does nothing.
				if !gstr.HasPrefix(match[2], "*") && !gstr.HasPrefix(match[2], "[]") && !gstr.HasPrefix(match[2], "map") {
					return fmt.Sprintf(`%s interface{} %s`, match[1], match[3])
				}
				return match[0]
			},
		)
		modelContent := generateDoContent(
			tableName,
			gstr.CaseCamel(newTableName),
			structDefinition,
		)
		err = gfile.PutContents(doFilePath, strings.TrimSpace(modelContent))
		if err != nil {
			mlog.Fatalf("writing content to '%s' failed: %v", doFilePath, err)
		} else {
			utils.GoFmt(doFilePath)
			mlog.Print("generated:", doFilePath)
		}
	}
}

func generateEntity(ctx context.Context, db gdb.DB, tableNames, newTableNames []string, in cGenDaoInternalInput) {
	var (
		entityDirPath = gfile.Join(in.Path, defaultEntityPath)
	)

	// Model content.
	for i, tableName := range tableNames {
		fieldMap, err := db.TableFields(ctx, tableName)
		if err != nil {
			mlog.Fatalf("fetching tables fields failed for table '%s':\n%v", in.TableName, err)
		}
		var (
			newTableName   = newTableNames[i]
			entityFilePath = gfile.Join(entityDirPath, gstr.CaseSnake(newTableName)+".go")
			entityContent  = generateEntityContent(
				newTableName,
				gstr.CaseCamel(newTableName),
				generateStructDefinition(generateStructDefinitionInput{
					cGenDaoInternalInput: in,
					StructName:           gstr.CaseCamel(newTableName),
					FieldMap:             fieldMap,
					IsDo:                 false,
				}),
			)
		)
		err = gfile.PutContents(entityFilePath, strings.TrimSpace(entityContent))
		if err != nil {
			mlog.Fatalf("writing content to '%s' failed: %v", entityFilePath, err)
		} else {
			utils.GoFmt(entityFilePath)
			mlog.Print("generated:", entityFilePath)
		}
	}
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

	// Generate and write content to golang file.
	packageImportsStr := ""
	if packageImportsArray.Len() > 0 {
		packageImportsStr = fmt.Sprintf("import(\n%s\n)", packageImportsArray.Join("\n"))
	}
	return packageImportsStr
}

func generateEntityContent(tableName, tableNameCamelCase, structDefine string) string {
	entityContent := gstr.ReplaceByMap(consts.TemplateGenDaoEntityContent, g.MapStrStr{
		tplVarTableName:          tableName,
		tplVarPackageImports:     getImportPartContent(structDefine, false),
		tplVarTableNameCamelCase: tableNameCamelCase,
		tplVarStructDefine:       structDefine,
	})
	entityContent = replaceDefaultVar(entityContent)
	return entityContent
}

func generateDoContent(tableName, tableNameCamelCase, structDefine string) string {
	doContent := gstr.ReplaceByMap(consts.TemplateGenDaoDoContent, g.MapStrStr{
		tplVarTableName:          tableName,
		tplVarPackageImports:     getImportPartContent(structDefine, true),
		tplVarTableNameCamelCase: tableNameCamelCase,
		tplVarStructDefine:       structDefine,
	})
	doContent = replaceDefaultVar(doContent)
	return doContent
}

func generateDaoIndex(tableNameCamelCase, tableNameCamelLowerCase, importPrefix, dirPathDao, fileName string, in cGenDaoInternalInput) {
	path := gfile.Join(dirPathDao, fileName+".go")
	if in.OverwriteDao || !gfile.Exists(path) {
		indexContent := gstr.ReplaceByMap(getTplDaoIndexContent(""), g.MapStrStr{
			tplVarImportPrefix:            importPrefix,
			tplVarTableName:               in.TableName,
			tplVarTableNameCamelCase:      tableNameCamelCase,
			tplVarTableNameCamelLowerCase: tableNameCamelLowerCase,
		})
		indexContent = replaceDefaultVar(indexContent)
		if err := gfile.PutContents(path, strings.TrimSpace(indexContent)); err != nil {
			mlog.Fatalf("writing content to '%s' failed: %v", path, err)
		} else {
			utils.GoFmt(path)
			mlog.Print("generated:", path)
		}
	}
}

func generateDaoInternal(
	tableNameCamelCase, tableNameCamelLowerCase, importPrefix string,
	dirPathDao, fileName string,
	fieldMap map[string]*gdb.TableField,
	in cGenDaoInternalInput,
) {
	path := gfile.Join(dirPathDao, "internal", fileName+".go")
	modelContent := gstr.ReplaceByMap(getTplDaoInternalContent(""), g.MapStrStr{
		tplVarImportPrefix:            importPrefix,
		tplVarTableName:               in.TableName,
		tplVarGroupName:               in.Group,
		tplVarTableNameCamelCase:      tableNameCamelCase,
		tplVarTableNameCamelLowerCase: tableNameCamelLowerCase,
		tplVarColumnDefine:            gstr.Trim(generateColumnDefinitionForDao(fieldMap)),
		tplVarColumnNames:             gstr.Trim(generateColumnNamesForDao(fieldMap)),
	})
	modelContent = replaceDefaultVar(modelContent)
	if err := gfile.PutContents(path, strings.TrimSpace(modelContent)); err != nil {
		mlog.Fatalf("writing content to '%s' failed: %v", path, err)
	} else {
		utils.GoFmt(path)
		mlog.Print("generated:", path)
	}
}

func replaceDefaultVar(origin string) string {
	return gstr.ReplaceByMap(origin, g.MapStrStr{
		tplVarDatetime: createdAt.String(),
	})
}

type generateStructDefinitionInput struct {
	cGenDaoInternalInput
	StructName string                     // Struct name.
	FieldMap   map[string]*gdb.TableField // Table field map.
	IsDo       bool                       // Is generating DTO struct.
}

func generateStructDefinition(in generateStructDefinitionInput) string {
	buffer := bytes.NewBuffer(nil)
	array := make([][]string, len(in.FieldMap))
	names := sortFieldKeyForDao(in.FieldMap)
	for index, name := range names {
		field := in.FieldMap[name]
		array[index] = generateStructFieldDefinition(field, in)
	}
	tw := tablewriter.NewWriter(buffer)
	tw.SetBorder(false)
	tw.SetRowLine(false)
	tw.SetAutoWrapText(false)
	tw.SetColumnSeparator("")
	tw.AppendBulk(array)
	tw.Render()
	stContent := buffer.String()
	// Let's do this hack of table writer for indent!
	stContent = gstr.Replace(stContent, "  #", "")
	stContent = gstr.Replace(stContent, "` ", "`")
	stContent = gstr.Replace(stContent, "``", "")
	buffer.Reset()
	buffer.WriteString(fmt.Sprintf("type %s struct {\n", in.StructName))
	if in.IsDo {
		buffer.WriteString(fmt.Sprintf("g.Meta `orm:\"table:%s, do:true\"`\n", in.TableName))
	}
	buffer.WriteString(stContent)
	buffer.WriteString("}")
	return buffer.String()
}

// generateStructFieldForModel generates and returns the attribute definition for specified field.
func generateStructFieldDefinition(field *gdb.TableField, in generateStructDefinitionInput) []string {
	var (
		typeName string
		jsonTag  = getJsonTagFromCase(field.Name, in.JsonCase)
	)
	t, _ := gregex.ReplaceString(`\(.+\)`, "", field.Type)
	t = gstr.Split(gstr.Trim(t), " ")[0]
	t = gstr.ToLower(t)
	switch t {
	case "binary", "varbinary", "blob", "tinyblob", "mediumblob", "longblob":
		typeName = "[]byte"

	case "bit", "int", "int2", "tinyint", "small_int", "smallint", "medium_int", "mediumint", "serial":
		if gstr.ContainsI(field.Type, "unsigned") {
			typeName = "uint"
		} else {
			typeName = "int"
		}

	case "int4", "int8", "big_int", "bigint", "bigserial":
		if gstr.ContainsI(field.Type, "unsigned") {
			typeName = "uint64"
		} else {
			typeName = "int64"
		}

	case "real":
		typeName = "float32"

	case "float", "double", "decimal", "smallmoney", "numeric":
		typeName = "float64"

	case "bool":
		typeName = "bool"

	case "datetime", "timestamp", "date", "time":
		if in.StdTime {
			typeName = "time.Time"
		} else {
			typeName = "*gtime.Time"
		}
	case "json", "jsonb":
		if in.GJsonSupport {
			typeName = "*gjson.Json"
		} else {
			typeName = "string"
		}
	default:
		// Automatically detect its data type.
		switch {
		case strings.Contains(t, "int"):
			typeName = "int"
		case strings.Contains(t, "text") || strings.Contains(t, "char"):
			typeName = "string"
		case strings.Contains(t, "float") || strings.Contains(t, "double"):
			typeName = "float64"
		case strings.Contains(t, "bool"):
			typeName = "bool"
		case strings.Contains(t, "binary") || strings.Contains(t, "blob"):
			typeName = "[]byte"
		case strings.Contains(t, "date") || strings.Contains(t, "time"):
			if in.StdTime {
				typeName = "time.Time"
			} else {
				typeName = "*gtime.Time"
			}
		default:
			typeName = "string"
		}
	}

	var (
		tagKey = "`"
		result = []string{
			"    #" + gstr.CaseCamel(field.Name),
			" #" + typeName,
		}
		descriptionTag = gstr.Replace(formatComment(field.Comment), `"`, `\"`)
	)

	result = append(result, " #"+fmt.Sprintf(tagKey+`json:"%s"`, jsonTag))
	result = append(result, " #"+fmt.Sprintf(`description:"%s"`+tagKey, descriptionTag))
	result = append(result, " #"+fmt.Sprintf(`// %s`, formatComment(field.Comment)))

	for k, v := range result {
		if in.NoJsonTag {
			v, _ = gregex.ReplaceString(`json:".+"`, ``, v)
		}
		if !in.DescriptionTag {
			v, _ = gregex.ReplaceString(`description:".*"`, ``, v)
		}
		if in.NoModelComment {
			v, _ = gregex.ReplaceString(`//.+`, ``, v)
		}
		result[k] = v
	}
	return result
}

// formatComment formats the comment string to fit the golang code without any lines.
func formatComment(comment string) string {
	comment = gstr.ReplaceByArray(comment, g.SliceStr{
		"\n", " ",
		"\r", " ",
	})
	comment = gstr.Replace(comment, `\n`, " ")
	comment = gstr.Trim(comment)
	return comment
}

// generateColumnDefinitionForDao generates and returns the column names definition for specified table.
func generateColumnDefinitionForDao(fieldMap map[string]*gdb.TableField) string {
	var (
		buffer = bytes.NewBuffer(nil)
		array  = make([][]string, len(fieldMap))
		names  = sortFieldKeyForDao(fieldMap)
	)
	for index, name := range names {
		var (
			field   = fieldMap[name]
			comment = gstr.Trim(gstr.ReplaceByArray(field.Comment, g.SliceStr{
				"\n", " ",
				"\r", " ",
			}))
		)
		array[index] = []string{
			"    #" + gstr.CaseCamel(field.Name),
			" # " + "string",
			" #" + fmt.Sprintf(`// %s`, comment),
		}
	}
	tw := tablewriter.NewWriter(buffer)
	tw.SetBorder(false)
	tw.SetRowLine(false)
	tw.SetAutoWrapText(false)
	tw.SetColumnSeparator("")
	tw.AppendBulk(array)
	tw.Render()
	defineContent := buffer.String()
	// Let's do this hack of table writer for indent!
	defineContent = gstr.Replace(defineContent, "  #", "")
	buffer.Reset()
	buffer.WriteString(defineContent)
	return buffer.String()
}

// generateColumnNamesForDao generates and returns the column names assignment content of column struct
// for specified table.
func generateColumnNamesForDao(fieldMap map[string]*gdb.TableField) string {
	var (
		buffer = bytes.NewBuffer(nil)
		array  = make([][]string, len(fieldMap))
		names  = sortFieldKeyForDao(fieldMap)
	)
	for index, name := range names {
		field := fieldMap[name]
		array[index] = []string{
			"            #" + gstr.CaseCamel(field.Name) + ":",
			fmt.Sprintf(` #"%s",`, field.Name),
		}
	}
	tw := tablewriter.NewWriter(buffer)
	tw.SetBorder(false)
	tw.SetRowLine(false)
	tw.SetAutoWrapText(false)
	tw.SetColumnSeparator("")
	tw.AppendBulk(array)
	tw.Render()
	namesContent := buffer.String()
	// Let's do this hack of table writer for indent!
	namesContent = gstr.Replace(namesContent, "  #", "")
	buffer.Reset()
	buffer.WriteString(namesContent)
	return buffer.String()
}

func getTplDaoIndexContent(tplDaoIndexPath string) string {
	if tplDaoIndexPath != "" {
		return gfile.GetContents(tplDaoIndexPath)
	}
	return consts.TemplateDaoDaoIndexContent
}

func getTplDaoInternalContent(tplDaoInternalPath string) string {
	if tplDaoInternalPath != "" {
		return gfile.GetContents(tplDaoInternalPath)
	}
	return consts.TemplateDaoDaoInternalContent
}

// getJsonTagFromCase call gstr.Case* function to convert the s to specified case.
func getJsonTagFromCase(str, caseStr string) string {
	switch gstr.ToLower(caseStr) {
	case gstr.ToLower("Camel"):
		return gstr.CaseCamel(str)

	case gstr.ToLower("CamelLower"):
		return gstr.CaseCamelLower(str)

	case gstr.ToLower("Kebab"):
		return gstr.CaseKebab(str)

	case gstr.ToLower("KebabScreaming"):
		return gstr.CaseKebabScreaming(str)

	case gstr.ToLower("Snake"):
		return gstr.CaseSnake(str)

	case gstr.ToLower("SnakeFirstUpper"):
		return gstr.CaseSnakeFirstUpper(str)

	case gstr.ToLower("SnakeScreaming"):
		return gstr.CaseSnakeScreaming(str)
	}
	return str
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
