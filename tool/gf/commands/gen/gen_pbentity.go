package gen

import (
	"bytes"
	"context"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/tool/gf/library/mlog"
	"github.com/gogf/gf/util/gconv"
	_ "github.com/lib/pq"
	//_ "github.com/mattn/go-oci8"
	//_ "github.com/mattn/go-sqlite3"
	"github.com/olekukonko/tablewriter"
	"strings"
)

// generatePbEntityReq is the input parameter for generating entity protobuf files.
type generatePbEntityReq struct {
	TableName     string // TableName specifies the table name of the table.
	NewTableName  string // NewTableName specifies the prefix-stripped name of the table.
	PrefixName    string // PrefixName specifies the custom prefix name for generated protobuf entity.
	GroupName     string // GroupName specifies the group name of database configuration node for generated protobuf entity.
	PkgName       string // PkgName specifies package name for generated protobuf.
	NameCase      string // NameCase specifies the case of generated attribute name for entity message, value from gstr.Case* function names.
	JsonCase      string // JsonCase specifies the case of json tag for attribute name of entity message, value from gstr.Case* function names.
	DirPath       string // DirPath specifies the directory path for generated files.
	OptionContent string // OptionContent specifies the extra option configuration content for protobuf.
	TplEntityPath string // TplEntityPath specifies the file path for generating protobuf entity files.
}

const (
	nodeNameGenPbEntityInConfigFile = "gfcli.gen.pbentity"
)

func HelpPbEntity() {
	mlog.Print(gstr.TrimLeft(`
USAGE 
    gf gen pbentity [OPTION]

OPTION
    -/--path             directory path for generated files.
    -/--package          package name for all entity proto files.
    -l, --link           database configuration, the same as the ORM configuration of GoFrame.
    -t, --tables         generate models only for given tables, multiple table names separated with ','
    -c, --config         used to specify the configuration file for database, it's commonly not necessary.
                         If "-l" is not passed, it will search "./config.toml" and "./config/config.toml" 
                         in current working directory in default.
    -p, --prefix         add specified prefix for all entity names and entity proto files.
    -r, --removePrefix   remove specified prefix of the table, multiple prefix separated with ',' 
    -n, --nameCase       case for message attribute names, default is "Camel": 
                         | Case            | Example            |
                         |---------------- |--------------------|
                         | Camel           | AnyKindOfString    | default
                         | CamelLower      | anyKindOfString    |
                         | Snake           | any_kind_of_string |
                         | SnakeScreaming  | ANY_KIND_OF_STRING |
                         | SnakeFirstUpper | rgb_code_md5       |
                         | Kebab           | any-kind-of-string |
                         | KebabScreaming  | ANY-KIND-OF-STRING |
    -j, --jsonCase       case for message json tag, cases are the same as "nameCase", default "CamelLower".
                         set it to "none" to ignore json tag generating.
    -o, --option         extra protobuf options.
    -/--tplEntity        template content for protobuf entity files generating.
                  
CONFIGURATION SUPPORT
    Options are also supported by configuration file.
    It's suggested using configuration file instead of command line arguments making producing. 
    The configuration node name is "gf.gen.pbentity", which also supports multiple databases, for example:
    [gfcli]
        [[gfcli.gen.pbentity]]
            link    = "mysql:root:12345678@tcp(127.0.0.1:3306)/test"
            path    = "protocol/demos/entity"
            tables  = "order,products"
            package = "demos"
        [[gfcli.gen.pbentity]]
            link    = "mysql:root:12345678@tcp(127.0.0.1:3306)/primary"
            path    = "protocol/demos/entity"
            prefix  = "primary_"
            tables  = "user, userDetail"
            package = "demos"
            option  = """
option go_package    = "protobuf/demos";
option java_package  = "protobuf/demos";
option php_namespace = "protobuf/demos";
"""

EXAMPLES
    gf gen pbentity
    gf gen pbentity -l "mysql:root:12345678@tcp(127.0.0.1:3306)/test"
    gf gen pbentity -path ./protocol/demos/entity -c config.yaml -g user-center -t user,user_detail,user_login
    gf gen pbentity -r user_
`))
}

// doGenPbEntity implements the "gen pbentity" command.
func doGenPbEntity() {
	parser, err := gcmd.Parse(g.MapStrBool{
		"path":           true,
		"package":        true,
		"l,link":         true,
		"t,tables":       true,
		"c,config":       true,
		"p,prefix":       true,
		"r,removePrefix": true,
		"o,option":       true,
		"n,nameCase":     true,
		"j,jsonCase":     true,
		"tplEntity":      true,
	})
	if err != nil {
		mlog.Fatal(err)
	}
	config := g.Cfg()
	if config.Available() {
		v := config.GetVar(nodeNameGenPbEntityInConfigFile)
		if v.IsEmpty() && g.IsEmpty(parser.GetOptAll()) {
			mlog.Fatal(`command arguments and configurations not found for generating protobuf entity files`)
		}
		if v.IsSlice() {
			for i := 0; i < len(v.Interfaces()); i++ {
				doGenPbEntityForArray(i, parser)
			}
		} else {
			doGenPbEntityForArray(-1, parser)
		}
	} else {
		doGenPbEntityForArray(-1, parser)
	}
	mlog.Print("done!")
}

// doGenPbEntityForArray implements the "gen pbentity" command for configuration array.
func doGenPbEntityForArray(index int, parser *gcmd.Parser) {
	var (
		err           error
		db            gdb.DB
		dirPath       = getOptionOrConfigForPbEntity(index, parser, "path")                   // Generated directory path.
		pkgName       = getOptionOrConfigForPbEntity(index, parser, "package")                // Package name for protobuf.
		tablesStr     = getOptionOrConfigForPbEntity(index, parser, "tables")                 // Tables that will be generated.
		prefixName    = getOptionOrConfigForPbEntity(index, parser, "prefix")                 // Add prefix to entity name.
		linkInfo      = getOptionOrConfigForPbEntity(index, parser, "link")                   // Custom database link.
		configPath    = getOptionOrConfigForPbEntity(index, parser, "config")                 // Config file path, eg: ./config/db.toml.
		configGroup   = getOptionOrConfigForPbEntity(index, parser, "group", "default")       // Group name of database configuration node for generated protobuf entity.
		removePrefix  = getOptionOrConfigForPbEntity(index, parser, "removePrefix")           // Remove prefix from table name.
		nameCase      = getOptionOrConfigForPbEntity(index, parser, "nameCase", "Camel")      // Case configuration for message name.
		jsonCase      = getOptionOrConfigForPbEntity(index, parser, "jsonCase", "CamelLower") // Case configuration for message json tag.
		optionContent = getOptionOrConfigForPbEntity(index, parser, "option")                 // Option content for protobuf.
		tplEntityPath = getOptionOrConfigForPbEntity(index, parser, "tplEntity")              // Specifies the file path for generating protobuf entity files.
	)
	if tplEntityPath != "" && (!gfile.Exists(tplEntityPath) || !gfile.IsReadable(tplEntityPath)) {
		mlog.Fatalf("template file for entity files generating does not exist or is not readable: %s", tplEntityPath)
	}
	// Make it compatible with old CLI version for option name: remove-prefix
	if removePrefix == "" {
		removePrefix = getOptionOrConfigForPbEntity(index, parser, "remove-prefix")
	}
	removePrefixArray := gstr.SplitAndTrim(removePrefix, ",")
	if pkgName == "" {
		mlog.Fatal("package name should not be empty")
	}
	// It reads database configuration from project configuration file.
	if configPath != "" {
		path, err := gfile.Search(configPath)
		if err != nil {
			mlog.Fatalf("search configuration file '%s' failed: %v", configPath, err)
		}
		if err := g.Cfg().SetPath(gfile.Dir(path)); err != nil {
			mlog.Fatalf("set configuration path '%s' failed: %v", path, err)
		}
		g.Cfg().SetFileName(gfile.Basename(path))
	}
	// It uses user passed database configuration.
	if linkInfo != "" {
		tempGroup := gtime.TimestampNanoStr()
		match, _ := gregex.MatchString(`([a-z]+):(.+)`, linkInfo)
		if len(match) == 3 {
			gdb.AddConfigNode(tempGroup, gdb.ConfigNode{
				Type: gstr.Trim(match[1]),
				Link: gstr.Trim(match[2]),
			})
			db, _ = gdb.Instance(tempGroup)
		}
	} else {
		db = g.DB(configGroup)
	}
	if db == nil {
		mlog.Fatal("database initialization failed")
	}

	tableNames := ([]string)(nil)
	if tablesStr != "" {
		tableNames = gstr.SplitAndTrim(tablesStr, ",")
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
		req := &generatePbEntityReq{
			TableName:     tableName,
			NewTableName:  newTableName,
			PrefixName:    prefixName,
			GroupName:     configGroup,
			PkgName:       pkgName,
			NameCase:      nameCase,
			JsonCase:      jsonCase,
			DirPath:       dirPath,
			OptionContent: gstr.Trim(optionContent),
			TplEntityPath: tplEntityPath,
		}
		generatePbEntityContentFile(db, req)
	}
}

// generatePbEntityContentFile generates the protobuf files for given table.
func generatePbEntityContentFile(db gdb.DB, req *generatePbEntityReq) {
	fieldMap, err := db.TableFields(db.GetCtx(), req.TableName)
	if err != nil {
		mlog.Fatalf("fetching tables fields failed for table '%s':\n%v", req.TableName, err)
	}
	// Change the `newTableName` if `prefixName` is given.
	newTableName := "Entity_" + req.PrefixName + req.NewTableName
	var (
		tableNameCamelCase  = gstr.CaseCamel(newTableName)
		tableNameSnakeCase  = gstr.CaseSnake(newTableName)
		entityMessageDefine = generateEntityMessageDefinition(tableNameCamelCase, fieldMap, req)
		fileName            = gstr.Trim(tableNameSnakeCase, "-_.")
		path                = gfile.Join(req.DirPath, fileName+".proto")
	)
	entityContent := gstr.ReplaceByMap(getTplPbEntityContent(req.TplEntityPath), g.MapStrStr{
		"{PackageName}":   req.PkgName,
		"{OptionContent}": req.OptionContent,
		"{EntityMessage}": entityMessageDefine,
	})
	if err := gfile.PutContents(path, strings.TrimSpace(entityContent)); err != nil {
		mlog.Fatalf("writing content to '%s' failed: %v", path, err)
	} else {
		mlog.Print("generated:", path)
	}
}

// generateEntityMessageDefinition generates and returns the message definition for specified table.
func generateEntityMessageDefinition(name string, fieldMap map[string]*gdb.TableField, req *generatePbEntityReq) string {
	var (
		buffer = bytes.NewBuffer(nil)
		array  = make([][]string, len(fieldMap))
		names  = sortFieldKeyForPbEntity(fieldMap)
	)
	for index, name := range names {
		array[index] = generateMessageFieldForPbEntity(index+1, fieldMap[name], req)
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
	buffer.WriteString(fmt.Sprintf("message %s {\n", name))
	buffer.WriteString(stContent)
	buffer.WriteString("}")
	return buffer.String()
}

// generateMessageFieldForPbEntity generates and returns the message definition for specified field.
func generateMessageFieldForPbEntity(index int, field *gdb.TableField, req *generatePbEntityReq) []string {
	var (
		typeName   string
		comment    string
		jsonTagStr string
	)
	t, _ := gregex.ReplaceString(`\(.+\)`, "", field.Type)
	t = gstr.Split(gstr.Trim(t), " ")[0]
	t = gstr.ToLower(t)
	switch t {
	case "binary", "varbinary", "blob", "tinyblob", "mediumblob", "longblob":
		typeName = "bytes"

	case "bit", "int", "tinyint", "small_int", "smallint", "medium_int", "mediumint", "serial":
		if gstr.ContainsI(field.Type, "unsigned") {
			typeName = "uint32"
		} else {
			typeName = "int32"
		}

	case "int8", "big_int", "bigint", "bigserial":
		if gstr.ContainsI(field.Type, "unsigned") {
			typeName = "uint64"
		} else {
			typeName = "int64"
		}

	case "real":
		typeName = "float"

	case "float", "double", "decimal", "smallmoney":
		typeName = "double"

	case "bool":
		typeName = "bool"

	case "datetime", "timestamp", "date", "time":
		typeName = "int64"

	default:
		// Auto detecting type.
		switch {
		case strings.Contains(t, "int"):
			typeName = "int"
		case strings.Contains(t, "text") || strings.Contains(t, "char"):
			typeName = "string"
		case strings.Contains(t, "float") || strings.Contains(t, "double"):
			typeName = "double"
		case strings.Contains(t, "bool"):
			typeName = "bool"
		case strings.Contains(t, "binary") || strings.Contains(t, "blob"):
			typeName = "bytes"
		case strings.Contains(t, "date") || strings.Contains(t, "time"):
			typeName = "int64"
		default:
			typeName = "string"
		}
	}
	comment = gstr.ReplaceByArray(field.Comment, g.SliceStr{
		"\n", " ",
		"\r", " ",
	})
	comment = gstr.Trim(comment)
	comment = gstr.Replace(comment, `\n`, " ")
	comment, _ = gregex.ReplaceString(`\s{2,}`, ` `, comment)
	if jsonTagName := formatCase(field.Name, req.JsonCase); jsonTagName != "" {
		jsonTagStr = fmt.Sprintf(`[(gogoproto.jsontag) = "%s"]`, jsonTagName)
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
		" #" + formatCase(field.Name, req.NameCase),
		" #= " + gconv.String(index) + jsonTagStr + ";",
		" #" + fmt.Sprintf(`// %s`, comment),
	}
}

func getTplPbEntityContent(tplEntityPath string) string {
	if tplEntityPath != "" {
		return gfile.GetContents(tplEntityPath)
	}
	return templatePbEntityMessageContent
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

// getOptionOrConfigForPbEntity retrieves option value from parser and configuration file.
// It returns the default value specified by parameter <value> is no value found.
func getOptionOrConfigForPbEntity(index int, parser *gcmd.Parser, name string, defaultValue ...string) (result string) {
	result = parser.GetOpt(name)
	if result == "" && g.Config().Available() {
		g.Cfg().SetViolenceCheck(true)
		if index >= 0 {
			result = g.Cfg().GetString(fmt.Sprintf(`%s.%d.%s`, nodeNameGenPbEntityInConfigFile, index, name))
		} else {
			result = g.Cfg().GetString(fmt.Sprintf(`%s.%s`, nodeNameGenPbEntityInConfigFile, name))
		}
	}
	if result == "" && len(defaultValue) > 0 {
		result = defaultValue[0]
	}
	return
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
