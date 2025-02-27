// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package genpbentity

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/olekukonko/tablewriter"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gtag"

	"github.com/gogf/gf/cmd/gf/v2/internal/consts"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/utils"
)

type (
	CGenPbEntity      struct{}
	CGenPbEntityInput struct {
		g.Meta            `name:"pbentity" config:"{CGenPbEntityConfig}" brief:"{CGenPbEntityBrief}" eg:"{CGenPbEntityEg}" ad:"{CGenPbEntityAd}"`
		Path              string `name:"path"              short:"p"  brief:"{CGenPbEntityBriefPath}" d:"manifest/protobuf/pbentity"`
		Package           string `name:"package"           short:"k"  brief:"{CGenPbEntityBriefPackage}"`
		GoPackage         string `name:"goPackage"           short:"g"  brief:"{CGenPbEntityBriefGoPackage}"`
		Link              string `name:"link"              short:"l"  brief:"{CGenPbEntityBriefLink}"`
		Tables            string `name:"tables"            short:"t"  brief:"{CGenPbEntityBriefTables}"`
		Prefix            string `name:"prefix"            short:"f"  brief:"{CGenPbEntityBriefPrefix}"`
		RemovePrefix      string `name:"removePrefix"      short:"r"  brief:"{CGenPbEntityBriefRemovePrefix}"`
		RemoveFieldPrefix string `name:"removeFieldPrefix" short:"rf" brief:"{CGenPbEntityBriefRemoveFieldPrefix}"`
		TablesEx          string `name:"tablesEx"          short:"x"  brief:"{CGenDaoBriefTablesEx}"`
		NameCase          string `name:"nameCase"          short:"n"  brief:"{CGenPbEntityBriefNameCase}" d:"Camel"`
		JsonCase          string `name:"jsonCase"          short:"j"  brief:"{CGenPbEntityBriefJsonCase}" d:"none"`
		Option            string `name:"option"            short:"o"  brief:"{CGenPbEntityBriefOption}"`

		TypeMapping  map[DBFieldTypeName]CustomAttributeType  `name:"typeMapping"  short:"y"  brief:"{CGenPbEntityBriefTypeMapping}"  orphan:"true"`
		FieldMapping map[DBTableFieldName]CustomAttributeType `name:"fieldMapping" short:"fm" brief:"{CGenPbEntityBriefFieldMapping}" orphan:"true"`
	}
	CGenPbEntityOutput struct{}

	CGenPbEntityInternalInput struct {
		CGenPbEntityInput
		DB           gdb.DB
		TableName    string // TableName specifies the table name of the table.
		NewTableName string // NewTableName specifies the prefix-stripped name of the table.
	}

	DBTableFieldName    = string
	DBFieldTypeName     = string
	CustomAttributeType struct {
		Type   string `brief:"custom attribute type name"`
		Import string `brief:"custom import for this type"`
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
            typeMapping:
              json:
                type: google.protobuf.Value
                import: google/protobuf/struct.proto
              jsonb:
                type: google.protobuf.Value
                import: google/protobuf/struct.proto
`
	CGenPbEntityBriefPath              = `directory path for generated files storing`
	CGenPbEntityBriefPackage           = `package path for all entity proto files`
	CGenPbEntityBriefGoPackage         = `go package path for all entity proto files`
	CGenPbEntityBriefLink              = `database configuration, the same as the ORM configuration of GoFrame`
	CGenPbEntityBriefTables            = `generate models only for given tables, multiple table names separated with ','`
	CGenPbEntityBriefPrefix            = `add specified prefix for all entity names and entity proto files`
	CGenPbEntityBriefRemovePrefix      = `remove specified prefix of the table, multiple prefix separated with ','`
	CGenPbEntityBriefTablesEx          = `generate all models exclude the specified tables, multiple prefix separated with ','`
	CGenPbEntityBriefRemoveFieldPrefix = `remove specified prefix of the field, multiple prefix separated with ','`
	CGenPbEntityBriefOption            = `extra protobuf options`
	CGenPbEntityBriefGroup             = `
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

	CGenPbEntityBriefTypeMapping  = `custom local type mapping for generated struct attributes relevant to fields of table`
	CGenPbEntityBriefFieldMapping = `custom local type mapping for generated struct attributes relevant to specific fields of table`
)

var defaultTypeMapping = map[DBFieldTypeName]CustomAttributeType{
	// gdb.LocalTypeString
	"string": {
		Type: "string",
	},
	// gdb.LocalTypeTime
	// "time": {
	// 	Type:   "google.protobuf.Duration",
	// 	Import: "google/protobuf/duration.proto",
	// },
	// gdb.LocalTypeDate
	"date": {
		Type:   "google.protobuf.Timestamp",
		Import: "google/protobuf/timestamp.proto",
	},
	// gdb.LocalTypeDatetime
	"datetime": {
		Type:   "google.protobuf.Timestamp",
		Import: "google/protobuf/timestamp.proto",
	},
	// gdb.LocalTypeInt
	"int": {
		Type: "int32",
	},
	// gdb.LocalTypeUint
	"uint": {
		Type: "uint32",
	},
	// gdb.LocalTypeInt64
	"int64": {
		Type: "int64",
	},
	// gdb.LocalTypeUint64
	"uint64": {
		Type: "uint64",
	},
	// gdb.LocalTypeIntSlice
	"[]int": {
		Type: "repeated int32",
	},
	// gdb.LocalTypeInt64Slice
	"[]int64": {
		Type: "repeated int64",
	},
	// gdb.LocalTypeUint64Slice
	"[]uint64": {
		Type: "repeated uint64",
	},
	// gdb.LocalTypeInt64Bytes
	"int64-bytes": {
		Type: "repeated int64",
	},
	// gdb.LocalTypeUint64Bytes
	"uint64-bytes": {
		Type: "repeated uint64",
	},
	// gdb.LocalTypeFloat32
	"float32": {
		Type: "float",
	},
	// gdb.LocalTypeFloat64
	"float64": {
		Type: "double",
	},
	// gdb.LocalTypeBytes
	"[]byte": {
		Type: "bytes",
	},
	// gdb.LocalTypeBool
	"bool": {
		Type: "bool",
	},
	// gdb.LocalTypeJson
	// "json": {
	// 	Type:   "google.protobuf.Value",
	// 	Import: "google/protobuf/struct.proto",
	// },
	// gdb.LocalTypeJsonb
	// "jsonb": {
	// 	Type:   "google.protobuf.Value",
	// 	Import: "google/protobuf/struct.proto",
	// },
}

func init() {
	gtag.Sets(g.MapStrStr{
		`CGenPbEntityConfig`:                 CGenPbEntityConfig,
		`CGenPbEntityBrief`:                  CGenPbEntityBrief,
		`CGenPbEntityEg`:                     CGenPbEntityEg,
		`CGenPbEntityAd`:                     CGenPbEntityAd,
		`CGenPbEntityBriefPath`:              CGenPbEntityBriefPath,
		`CGenPbEntityBriefPackage`:           CGenPbEntityBriefPackage,
		`CGenPbEntityBriefGoPackage`:         CGenPbEntityBriefGoPackage,
		`CGenPbEntityBriefLink`:              CGenPbEntityBriefLink,
		`CGenPbEntityBriefTables`:            CGenPbEntityBriefTables,
		`CGenPbEntityBriefPrefix`:            CGenPbEntityBriefPrefix,
		`CGenPbEntityBriefRemovePrefix`:      CGenPbEntityBriefRemovePrefix,
		`CGenPbEntityBriefTablesEx`:          CGenPbEntityBriefTablesEx,
		`CGenPbEntityBriefRemoveFieldPrefix`: CGenPbEntityBriefRemoveFieldPrefix,
		`CGenPbEntityBriefGroup`:             CGenPbEntityBriefGroup,
		`CGenPbEntityBriefNameCase`:          CGenPbEntityBriefNameCase,
		`CGenPbEntityBriefJsonCase`:          CGenPbEntityBriefJsonCase,
		`CGenPbEntityBriefOption`:            CGenPbEntityBriefOption,
		`CGenPbEntityBriefTypeMapping`:       CGenPbEntityBriefTypeMapping,
		`CGenPbEntityBriefFieldMapping`:      CGenPbEntityBriefFieldMapping,
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
		modName := utils.GetImportPath(gfile.Pwd())
		in.Package = modName + "/" + defaultPackageSuffix
	}
	removePrefixArray := gstr.SplitAndTrim(in.RemovePrefix, ",")

	excludeTables := gset.NewStrSetFrom(gstr.SplitAndTrim(in.TablesEx, ","))

	// It uses user passed database configuration.
	if in.Link != "" {
		var (
			tempGroup = gtime.TimestampNanoStr()
			match, _  = gregex.MatchString(`([a-z]+):(.+)`, in.Link)
		)
		if len(match) == 3 {
			gdb.AddConfigNode(tempGroup, gdb.ConfigNode{
				Type: gstr.Trim(match[1]),
				Link: in.Link,
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

	for _, tableName := range tableNames {
		if excludeTables.Contains(tableName) {
			continue
		}
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
		tableNameCamelCase                 = gstr.CaseCamel(newTableName)
		tableNameSnakeCase                 = gstr.CaseSnake(newTableName)
		entityMessageDefine, appendImports = generateEntityMessageDefinition(tableNameCamelCase, fieldMap, in)
		fileName                           = gstr.Trim(tableNameSnakeCase, "-_.")
		path                               = filepath.FromSlash(gfile.Join(in.Path, fileName+".proto"))
	)
	packageImportStr := ""
	var packageImportsArray = garray.NewStrArray()
	if len(appendImports) > 0 {
		for _, appendImport := range appendImports {
			packageImportStr = fmt.Sprintf(`import "%s";`, appendImport)
			if packageImportsArray.Search(packageImportStr) == -1 {
				packageImportsArray.Append(packageImportStr)
			}
		}
	}
	if in.GoPackage == "" {
		in.GoPackage = in.Package
	}
	entityContent := gstr.ReplaceByMap(getTplPbEntityContent(""), g.MapStrStr{
		"{Imports}":       packageImportsArray.Join("\n"),
		"{PackageName}":   gfile.Basename(in.Package),
		"{GoPackage}":     in.GoPackage,
		"{OptionContent}": in.Option,
		"{EntityMessage}": entityMessageDefine,
	})
	if err := gfile.PutContents(path, strings.TrimSpace(entityContent)); err != nil {
		mlog.Fatalf("writing content to '%s' failed: %v", path, err)
	} else {
		mlog.Print("generated:", gfile.RealPath(path))
	}
}

// generateEntityMessageDefinition generates and returns the message definition for specified table.
func generateEntityMessageDefinition(entityName string, fieldMap map[string]*gdb.TableField, in CGenPbEntityInternalInput) (string, []string) {
	var (
		appendImports []string
		buffer        = bytes.NewBuffer(nil)
		array         = make([][]string, len(fieldMap))
		names         = sortFieldKeyForPbEntity(fieldMap)
	)
	for index, name := range names {
		var imports string
		array[index], imports = generateMessageFieldForPbEntity(index+1, fieldMap[name], in)
		if imports != "" {
			appendImports = append(appendImports, imports)
		}
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
	stContent = regexp.MustCompile(`\s+\n`).ReplaceAllString(gstr.Replace(stContent, "  #", ""), "\n")
	buffer.Reset()
	buffer.WriteString(fmt.Sprintf("message %s {\n", entityName))
	buffer.WriteString(stContent)
	buffer.WriteString("}")
	return buffer.String(), appendImports
}

// generateMessageFieldForPbEntity generates and returns the message definition for specified field.
func generateMessageFieldForPbEntity(index int, field *gdb.TableField, in CGenPbEntityInternalInput) (attrLines []string, appendImport string) {
	var (
		localTypeNameStr string
		localTypeName    gdb.LocalType
		comment          string
		jsonTagStr       string
		err              error
		ctx              = gctx.GetInitCtx()
	)

	if in.TypeMapping != nil && len(in.TypeMapping) > 0 {
		localTypeName, err = in.DB.CheckLocalTypeForField(ctx, field.Type, nil)
		if err != nil {
			panic(err)
		}
		if localTypeName != "" {
			if typeMapping, ok := in.TypeMapping[strings.ToLower(string(localTypeName))]; ok {
				localTypeNameStr = typeMapping.Type
				appendImport = typeMapping.Import
			}
		}
	}

	if localTypeNameStr == "" {
		localTypeNameStr = "string"
	}

	comment = gstr.ReplaceByArray(field.Comment, g.SliceStr{
		"\n", " ",
		"\r", " ",
	})
	comment = gstr.Trim(comment)
	comment = gstr.Replace(comment, `\n`, " ")
	comment, _ = gregex.ReplaceString(`\s{2,}`, ` `, comment)
	if jsonTagName := formatCase(field.Name, in.JsonCase); jsonTagName != "" {
		jsonTagStr = fmt.Sprintf(`[json_name = "%s"]`, jsonTagName)
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

	removeFieldPrefixArray := gstr.SplitAndTrim(in.RemoveFieldPrefix, ",")
	newFiledName := field.Name
	for _, v := range removeFieldPrefixArray {
		newFiledName = gstr.TrimLeftStr(newFiledName, v, 1)
	}

	if in.FieldMapping != nil && len(in.FieldMapping) > 0 {
		if typeMapping, ok := in.FieldMapping[fmt.Sprintf("%s.%s", in.TableName, newFiledName)]; ok {
			localTypeNameStr = typeMapping.Type
			appendImport = typeMapping.Import
		}
	}

	return []string{
		"    #" + localTypeNameStr,
		" #" + formatCase(newFiledName, in.NameCase),
		" #= " + gconv.String(index) + jsonTagStr + ";",
		" #" + fmt.Sprintf(`// %s`, comment),
	}, appendImport
}

func getTplPbEntityContent(tplEntityPath string) string {
	if tplEntityPath != "" {
		return gfile.GetContents(tplEntityPath)
	}
	return consts.TemplatePbEntityMessageContent
}

// formatCase call gstr.Case* function to convert the s to specified case.
func formatCase(str, caseStr string) string {
	if caseStr == "none" {
		return ""
	}
	return gstr.CaseConvert(str, gstr.CaseTypeMatch(caseStr))
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
