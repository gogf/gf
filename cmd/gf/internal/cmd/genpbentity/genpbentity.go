package genpbentity

import (
	"bytes"
	"context"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"strings"

	"github.com/gogf/gf/cmd/gf/v2/internal/consts"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gtag"
)

type (
	CGenPbEntity      struct{}
	CGenPbEntityInput struct {
		g.Meta       `name:"pbentity" config:"{CGenPbEntityConfig}" brief:"{CGenPbEntityBrief}" eg:"{CGenPbEntityEg}" ad:"{CGenPbEntityAd}"`
		Path         string `name:"path"         short:"p" brief:"{CGenPbEntityBriefPath}" d:"manifest/protobuf/pbentity"`
		Package      string `name:"package"      short:"k" brief:"{CGenPbEntityBriefPackage}"`
		Link         string `name:"link"         short:"l" brief:"{CGenPbEntityBriefLink}"`
		Tables       string `name:"tables"       short:"t" brief:"{CGenPbEntityBriefTables}"`
		Prefix       string `name:"prefix"       short:"f" brief:"{CGenPbEntityBriefPrefix}"`
		RemovePrefix string `name:"removePrefix" short:"r" brief:"{CGenPbEntityBriefRemovePrefix}"`
		NameCase     string `name:"nameCase"     short:"n" brief:"{CGenPbEntityBriefNameCase}" d:"Camel"`
		JsonCase     string `name:"jsonCase"     short:"j" brief:"{CGenPbEntityBriefJsonCase}" d:"CamelLower"`
		Option       string `name:"option"       short:"o" brief:"{CGenPbEntityBriefOption}"`
	}
	CGenPbEntityOutput struct{}

	CGenPbEntityInternalInput struct {
		CGenPbEntityInput
		DB           gdb.DB
		TableName    string // TableName specifies the table name of the table.
		NewTableName string // NewTableName specifies the prefix-stripped name of the table.
	}
)

const (
	defaultPackageSuffix = `api/pbentity`
	CGenPbEntityConfig   = `gfcli.gen.pbentity`
	CGenPbEntityBrief    = `generate entity message files in protobuf3 format`
	CGenPbEntityEg       = `
gf gen pbentity
gf gen pbentity -l "mysql:root:12345678@tcp(127.0.0.1:3306)/test"
gf gen pbentity -p ./protocol/demos/entity -t user,user_detail,user_login
gf gen pbentity -r user_ -k github.com/gogf/gf/example/protobuf
gf gen pbentity -r user_
`

	CGenPbEntityAd = `
CONFIGURATION SUPPORT
    Options are also supported by configuration file.
    It's suggested using configuration file instead of command line arguments making producing. 
    The configuration node name is "gf.gen.pbentity", which also supports multiple databases, for example(config.yaml):
    gfcli:
      gen:
      - pbentity:
            link:    "mysql:root:12345678@tcp(127.0.0.1:3306)/test"
            path:    "protocol/demos/entity"
            tables:  "order,products"
            package: "demos"
      - pbentity:
            link:    "mysql:root:12345678@tcp(127.0.0.1:3306)/primary"
            path:    "protocol/demos/entity"
            prefix:  "primary_"
            tables:  "user, userDetail"
            package: "demos"
            option:  |
			  option go_package    = "protobuf/demos";
			  option java_package  = "protobuf/demos";
			  option php_namespace = "protobuf/demos";
`
	CGenPbEntityBriefPath         = `directory path for generated files storing`
	CGenPbEntityBriefPackage      = `package path for all entity proto files`
	CGenPbEntityBriefLink         = `database configuration, the same as the ORM configuration of GoFrame`
	CGenPbEntityBriefTables       = `generate models only for given tables, multiple table names separated with ','`
	CGenPbEntityBriefPrefix       = `add specified prefix for all entity names and entity proto files`
	CGenPbEntityBriefRemovePrefix = `remove specified prefix of the table, multiple prefix separated with ','`
	CGenPbEntityBriefOption       = `extra protobuf options`
	CGenPbEntityBriefGroup        = `
specifying the configuration group name of database for generated ORM instance,
it's not necessary and the default value is "default"
`

	CGenPbEntityBriefNameCase = `
case for message attribute names, default is "Camel":
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

	CGenPbEntityBriefJsonCase = `
case for message json tag, cases are the same as "nameCase", default "CamelLower".
set it to "none" to ignore json tag generating.
`
)

func init() {
	gtag.Sets(g.MapStrStr{
		`CGenPbEntityConfig`:            CGenPbEntityConfig,
		`CGenPbEntityBrief`:             CGenPbEntityBrief,
		`CGenPbEntityEg`:                CGenPbEntityEg,
		`CGenPbEntityAd`:                CGenPbEntityAd,
		`CGenPbEntityBriefPath`:         CGenPbEntityBriefPath,
		`CGenPbEntityBriefPackage`:      CGenPbEntityBriefPackage,
		`CGenPbEntityBriefLink`:         CGenPbEntityBriefLink,
		`CGenPbEntityBriefTables`:       CGenPbEntityBriefTables,
		`CGenPbEntityBriefPrefix`:       CGenPbEntityBriefPrefix,
		`CGenPbEntityBriefRemovePrefix`: CGenPbEntityBriefRemovePrefix,
		`CGenPbEntityBriefGroup`:        CGenPbEntityBriefGroup,
		`CGenPbEntityBriefNameCase`:     CGenPbEntityBriefNameCase,
		`CGenPbEntityBriefJsonCase`:     CGenPbEntityBriefJsonCase,
		`CGenPbEntityBriefOption`:       CGenPbEntityBriefOption,
	})
}

func (c CGenPbEntity) PbEntity(ctx context.Context, in CGenPbEntityInput) (out *CGenPbEntityOutput, err error) {
	var (
		config = g.Cfg()
	)
	if config.Available(ctx) {
		v := config.MustGet(ctx, CGenPbEntityConfig)
		if v.IsSlice() {
			for i := 0; i < len(v.Interfaces()); i++ {
				doGenPbEntityForArray(ctx, i, in)
			}
		} else {
			doGenPbEntityForArray(ctx, -1, in)
		}
	} else {
		doGenPbEntityForArray(ctx, -1, in)
	}
	mlog.Print("done!")
	return
}

func doGenPbEntityForArray(ctx context.Context, index int, in CGenPbEntityInput) {
	var (
		err error
		db  gdb.DB
	)
	if index >= 0 {
		err = g.Cfg().MustGet(
			ctx,
			fmt.Sprintf(`%s.%d`, CGenPbEntityConfig, index),
		).Scan(&in)
		if err != nil {
			mlog.Fatalf(`invalid configuration of "%s": %+v`, CGenPbEntityConfig, err)
		}
	}
	if in.Package == "" {
		mlog.Debug(`package parameter is empty, trying calculating the package path using go.mod`)
		if !gfile.Exists("go.mod") {
			mlog.Fatal("go.mod does not exist in current working directory")
		}
		var (
			modName      string
			goModContent = gfile.GetContents("go.mod")
			match, _     = gregex.MatchString(`^module\s+(.+)\s*`, goModContent)
		)
		if len(match) > 1 {
			modName = gstr.Trim(match[1])
			in.Package = modName + "/" + defaultPackageSuffix
		} else {
			mlog.Fatal("module name does not found in go.mod")
		}
	}
	removePrefixArray := gstr.SplitAndTrim(in.RemovePrefix, ",")
	// It uses user passed database configuration.
	if in.Link != "" {
		var (
			tempGroup = gtime.TimestampNanoStr()
			match, _  = gregex.MatchString(`([a-z]+):(.+)`, in.Link)
		)
		if len(match) == 3 {
			gdb.AddConfigNode(tempGroup, gdb.ConfigNode{
				Type: gstr.Trim(match[1]),
				Link: gstr.Trim(match[2]),
			})
			db, _ = gdb.Instance(tempGroup)
		}
	} else {
		db = g.DB()
	}
	if db == nil {
		mlog.Fatal("database initialization failed")
	}

	tableNames := ([]string)(nil)
	if in.Tables != "" {
		tableNames = gstr.SplitAndTrim(in.Tables, ",")
	} else {
		tableNames, err = db.Tables(context.TODO())
		if err != nil {
			mlog.Fatalf("fetching tables failed: \n %v", err)
		}
	}

	for _, tableName := range tableNames {
		newTableName := tableName
		for _, v := range removePrefixArray {
			newTableName = gstr.TrimLeftStr(newTableName, v, 1)
		}
		generatePbEntityContentFile(ctx, CGenPbEntityInternalInput{
			CGenPbEntityInput: in,
			DB:                db,
			TableName:         tableName,
			NewTableName:      newTableName,
		})
	}
}

// generatePbEntityContentFile generates the protobuf files for given table.
func generatePbEntityContentFile(ctx context.Context, in CGenPbEntityInternalInput) {
	fieldMap, err := in.DB.TableFields(ctx, in.TableName)
	if err != nil {
		mlog.Fatalf("fetching tables fields failed for table '%s':\n%v", in.TableName, err)
	}
	// Change the `newTableName` if `Prefix` is given.
	newTableName := in.Prefix + in.NewTableName
	var (
		imports             string
		tableNameCamelCase  = gstr.CaseCamel(newTableName)
		tableNameSnakeCase  = gstr.CaseSnake(newTableName)
		entityMessageDefine = generateEntityMessageDefinition(tableNameCamelCase, fieldMap, in)
		fileName            = gstr.Trim(tableNameSnakeCase, "-_.")
		path                = gfile.Join(in.Path, fileName+".proto")
	)
	if gstr.Contains(entityMessageDefine, "google.protobuf.Timestamp") {
		imports = `import "google/protobuf/timestamp.proto";`
	}
	entityContent := gstr.ReplaceByMap(getTplPbEntityContent(""), g.MapStrStr{
		"{Imports}":       imports,
		"{PackageName}":   gfile.Basename(in.Package),
		"{GoPackage}":     in.Package,
		"{OptionContent}": in.Option,
		"{EntityMessage}": entityMessageDefine,
	})
	if err := gfile.PutContents(path, strings.TrimSpace(entityContent)); err != nil {
		mlog.Fatalf("writing content to '%s' failed: %v", path, err)
	} else {
		mlog.Print("generated:", path)
	}
}

// generateEntityMessageDefinition generates and returns the message definition for specified table.
func generateEntityMessageDefinition(entityName string, fieldMap map[string]*gdb.TableField, in CGenPbEntityInternalInput) string {
	var (
		buffer = bytes.NewBuffer(nil)
		array  = make([][]string, len(fieldMap))
		names  = sortFieldKeyForPbEntity(fieldMap)
	)
	for index, name := range names {
		array[index] = generateMessageFieldForPbEntity(index+1, fieldMap[name], in)
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
	buffer.Reset()
	buffer.WriteString(fmt.Sprintf("message %s {\n", entityName))
	buffer.WriteString(stContent)
	buffer.WriteString("}")
	return buffer.String()
}

// generateMessageFieldForPbEntity generates and returns the message definition for specified field.
func generateMessageFieldForPbEntity(index int, field *gdb.TableField, in CGenPbEntityInternalInput) []string {
	var (
		typeName   string
		comment    string
		jsonTagStr string
		err        error
		ctx        = gctx.GetInitCtx()
	)
	typeName, err = in.DB.CheckLocalTypeForField(ctx, field.Type, nil)
	if err != nil {
		panic(err)
	}
	var typeMapping = map[string]string{
		gdb.LocalTypeString:      "string",
		gdb.LocalTypeDate:        "google.protobuf.Timestamp",
		gdb.LocalTypeDatetime:    "google.protobuf.Timestamp",
		gdb.LocalTypeInt:         "int32",
		gdb.LocalTypeUint:        "uint32",
		gdb.LocalTypeInt64:       "int64",
		gdb.LocalTypeUint64:      "uint64",
		gdb.LocalTypeIntSlice:    "repeated int32",
		gdb.LocalTypeInt64Slice:  "repeated int64",
		gdb.LocalTypeUint64Slice: "repeated uint64",
		gdb.LocalTypeInt64Bytes:  "repeated int64",
		gdb.LocalTypeUint64Bytes: "repeated uint64",
		gdb.LocalTypeFloat32:     "float32",
		gdb.LocalTypeFloat64:     "float64",
		gdb.LocalTypeBytes:       "[]byte",
		gdb.LocalTypeBool:        "bool",
		gdb.LocalTypeJson:        "string",
		gdb.LocalTypeJsonb:       "string",
	}
	typeName = typeMapping[typeName]
	if typeName == "" {
		typeName = "string"
	}

	comment = gstr.ReplaceByArray(field.Comment, g.SliceStr{
		"\n", " ",
		"\r", " ",
	})
	comment = gstr.Trim(comment)
	comment = gstr.Replace(comment, `\n`, " ")
	comment, _ = gregex.ReplaceString(`\s{2,}`, ` `, comment)
	if jsonTagName := formatCase(field.Name, in.JsonCase); jsonTagName != "" {
		// beautiful indent.
		if index < 10 {
			// 3 spaces
			jsonTagStr = "   " + jsonTagStr
		} else if index < 100 {
			// 2 spaces
			jsonTagStr = "  " + jsonTagStr
		} else {
			// 1 spaces
			jsonTagStr = " " + jsonTagStr
		}
	}
	return []string{
		"    #" + typeName,
		" #" + formatCase(field.Name, in.NameCase),
		" #= " + gconv.String(index) + jsonTagStr + ";",
		" #" + fmt.Sprintf(`// %s`, comment),
	}
}

func getTplPbEntityContent(tplEntityPath string) string {
	if tplEntityPath != "" {
		return gfile.GetContents(tplEntityPath)
	}
	return consts.TemplatePbEntityMessageContent
}

// formatCase call gstr.Case* function to convert the s to specified case.
func formatCase(str, caseStr string) string {
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

	case "none":
		return ""
	}
	return str
}

func sortFieldKeyForPbEntity(fieldMap map[string]*gdb.TableField) []string {
	names := make(map[int]string)
	for _, field := range fieldMap {
		names[field.Index] = field.Name
	}
	var (
		result = make([]string, len(names))
		i      = 0
		j      = 0
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
