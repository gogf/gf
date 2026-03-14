// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gendao

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	"golang.org/x/mod/modfile"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/os/gview"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
	"github.com/gogf/gf/cmd/gf/v2/internal/utility/utils"
)

type (
	// CGenDao is the command handler struct for "gen dao" command.
	CGenDao struct{}

	// CGenDaoInput defines all input parameters for the "gen dao" command.
	// It supports both command-line arguments and configuration file options.
	CGenDaoInput struct {
		g.Meta             `name:"dao" config:"{CGenDaoConfig}" usage:"{CGenDaoUsage}" brief:"{CGenDaoBrief}" eg:"{CGenDaoEg}" ad:"{CGenDaoAd}"`
		Path               string   `name:"path"                short:"p"  brief:"{CGenDaoBriefPath}" d:"internal"`          // Base directory path for generated files.
		Link               string   `name:"link"                short:"l"  brief:"{CGenDaoBriefLink}"`                        // Database connection string (e.g., "mysql:root:pass@tcp(127.0.0.1:3306)/db").
		Tables             string   `name:"tables"              short:"t"  brief:"{CGenDaoBriefTables}"`                      // Comma-separated table names or wildcard patterns to include.
		TablesEx           string   `name:"tablesEx"            short:"x"  brief:"{CGenDaoBriefTablesEx}"`                    // Comma-separated table names or wildcard patterns to exclude.
		ShardingPattern    []string `name:"shardingPattern"     short:"sp" brief:"{CGenDaoBriefShardingPattern}"`             // Patterns for sharding tables (e.g., "users_?" merges users_001, users_002 into one dao).
		Group              string   `name:"group"               short:"g"  brief:"{CGenDaoBriefGroup}" d:"default"`           // Database configuration group name for ORM instance.
		Prefix             string   `name:"prefix"              short:"f"  brief:"{CGenDaoBriefPrefix}"`                      // Prefix to add to all generated table names.
		RemovePrefix       string   `name:"removePrefix"        short:"r"  brief:"{CGenDaoBriefRemovePrefix}"`                // Comma-separated prefixes to remove from table names.
		RemoveFieldPrefix  string   `name:"removeFieldPrefix"   short:"rf" brief:"{CGenDaoBriefRemoveFieldPrefix}"`           // Comma-separated prefixes to remove from field names.
		JsonCase           string   `name:"jsonCase"            short:"j"  brief:"{CGenDaoBriefJsonCase}" d:"CamelLower"`     // Naming convention for JSON tags (e.g., CamelLower, Snake).
		ImportPrefix       string   `name:"importPrefix"        short:"i"  brief:"{CGenDaoBriefImportPrefix}"`                // Custom Go import path prefix for generated files.
		DaoPath            string   `name:"daoPath"             short:"d"  brief:"{CGenDaoBriefDaoPath}" d:"dao"`             // Sub-directory under Path for dao files.
		TablePath          string   `name:"tablePath"           short:"tp" brief:"{CGenDaoBriefTablePath}" d:"table"`         // Sub-directory under Path for table field definition files.
		DoPath             string   `name:"doPath"              short:"o"  brief:"{CGenDaoBriefDoPath}" d:"model/do"`         // Sub-directory under Path for DO (Data Object) files.
		EntityPath         string   `name:"entityPath"          short:"e"  brief:"{CGenDaoBriefEntityPath}" d:"model/entity"` // Sub-directory under Path for entity struct files.
		TplDaoTablePath    string   `name:"tplDaoTablePath"     short:"t0" brief:"{CGenDaoBriefTplDaoTablePath}"`             // Custom template file for dao table generation.
		TplDaoIndexPath    string   `name:"tplDaoIndexPath"     short:"t1" brief:"{CGenDaoBriefTplDaoIndexPath}"`             // Custom template file for dao index generation.
		TplDaoInternalPath string   `name:"tplDaoInternalPath"  short:"t2" brief:"{CGenDaoBriefTplDaoInternalPath}"`          // Custom template file for dao internal generation.
		TplDaoDoPath       string   `name:"tplDaoDoPath"        short:"t3" brief:"{CGenDaoBriefTplDaoDoPathPath}"`            // Custom template file for DO generation.
		TplDaoEntityPath   string   `name:"tplDaoEntityPath"    short:"t4" brief:"{CGenDaoBriefTplDaoEntityPath}"`            // Custom template file for entity generation.
		StdTime            bool     `name:"stdTime"             short:"s"  brief:"{CGenDaoBriefStdTime}" orphan:"true"`       // Use stdlib time.Time instead of gtime.Time for time fields.
		WithTime           bool     `name:"withTime"            short:"w"  brief:"{CGenDaoBriefWithTime}" orphan:"true"`      // Add creation timestamp to generated file headers.
		GJsonSupport       bool     `name:"gJsonSupport"        short:"n"  brief:"{CGenDaoBriefGJsonSupport}" orphan:"true"`  // Use *gjson.Json instead of string for JSON fields.
		OverwriteDao       bool     `name:"overwriteDao"        short:"v"  brief:"{CGenDaoBriefOverwriteDao}" orphan:"true"`  // Overwrite existing dao files (both index and internal).
		DescriptionTag     bool     `name:"descriptionTag"      short:"c"  brief:"{CGenDaoBriefDescriptionTag}" orphan:"true"`// Add description struct tag with field comment.
		NoJsonTag          bool     `name:"noJsonTag"           short:"k"  brief:"{CGenDaoBriefNoJsonTag}" orphan:"true"`     // Omit json struct tags from generated structs.
		NoModelComment     bool     `name:"noModelComment"      short:"m"  brief:"{CGenDaoBriefNoModelComment}" orphan:"true"`// Omit inline comments from generated struct fields.
		Clear              bool     `name:"clear"               short:"a"  brief:"{CGenDaoBriefClear}" orphan:"true"`         // Delete generated files that no longer correspond to database tables.
		GenTable           bool     `name:"genTable"            short:"gt" brief:"{CGenDaoBriefGenTable}" orphan:"true"`      // Enable generation of table field definition files.
		SqlDir             string   `name:"sqlDir"              short:"sd" brief:"{CGenDaoBriefSqlDir}"`                      // Directory of SQL DDL files for offline generation (no DB connection needed).
		SqlType            string   `name:"sqlType"             short:"st" brief:"{CGenDaoBriefSqlType}" d:"mysql"`           // SQL dialect when using SqlDir (mysql, pgsql, mssql, oracle, sqlite).

		// TypeMapping maps database field type names to custom Go types.
		// For example, mapping "decimal" to "float64" or "uuid" to "uuid.UUID".
		TypeMapping map[DBFieldTypeName]CustomAttributeType `name:"typeMapping"  short:"y"  brief:"{CGenDaoBriefTypeMapping}"  orphan:"true"`
		// FieldMapping maps specific table.field combinations to custom Go types.
		// For example, mapping "user.balance" to "decimal.Decimal".
		FieldMapping map[DBTableFieldName]CustomAttributeType `name:"fieldMapping" short:"fm" brief:"{CGenDaoBriefFieldMapping}" orphan:"true"`

		// genItems tracks all generated file paths and directories for cleanup purposes.
		genItems *CGenDaoInternalGenItems
	}

	// CGenDaoOutput is the output of the "gen dao" command (currently empty).
	CGenDaoOutput struct{}

	// CGenDaoInternalInput extends CGenDaoInput with runtime-resolved fields
	// used during the actual generation process.
	CGenDaoInternalInput struct {
		CGenDaoInput
		DB               gdb.DB       // Database connection instance (nil in SQL file mode).
		TableNames       []string     // Original table names from database or SQL files.
		NewTableNames    []string     // Processed table names after prefix removal and sharding.
		ShardingTableSet *gset.StrSet // Set of table names identified as sharding tables.
		// TableFieldsMap stores pre-parsed table fields from SQL files.
		// When this is set (SQL file mode), DB may be nil.
		TableFieldsMap map[string]map[string]*gdb.TableField
	}

	// DBTableFieldName is the fully-qualified field name in "table.field" format.
	DBTableFieldName = string
	// DBFieldTypeName is the database column type name (e.g., "varchar", "decimal").
	DBFieldTypeName = string
	// CustomAttributeType defines a custom Go type mapping with its import path.
	CustomAttributeType struct {
		Type   string `brief:"custom attribute type name"` // Go type name (e.g., "decimal.Decimal").
		Import string `brief:"custom import for this type"` // Go import path (e.g., "github.com/shopspring/decimal").
	}
)

var (
	createdAt = gtime.Now()  // Timestamp captured at program start, used in generated file headers.
	tplView   = gview.New()  // Shared template view instance for rendering all Go file templates.
	// defaultTypeMapping provides built-in type mappings from database types to Go types.
	// User-provided TypeMapping takes precedence over these defaults.
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
		"uuid": {
			Type:   "uuid.UUID",
			Import: "github.com/google/uuid",
		},
	}

	// twRenderer configures the tablewriter to render without borders or separators,
	// producing clean aligned text output for generated Go source code.
	twRenderer = tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{
		Borders: tw.Border{Top: tw.Off, Bottom: tw.Off, Left: tw.Off, Right: tw.Off},
		Settings: tw.Settings{
			Separators: tw.Separators{BetweenRows: tw.Off, BetweenColumns: tw.Off},
		},
		Symbols: tw.NewSymbols(tw.StyleASCII),
	}))
	twConfig = tablewriter.WithConfig(tablewriter.Config{
		Row: tw.CellConfig{
			Formatting: tw.CellFormatting{AutoWrap: tw.WrapNone},
		},
	})
)

// Dao is the main entry point for the "gen dao" command.
// It dispatches to the appropriate generation mode based on input:
//   - SQL file mode (SqlDir is set): generates from DDL files without database connection.
//   - Link mode (Link is set): uses a direct database connection string.
//   - Config mode: reads database configuration from the application config file.
func (c CGenDao) Dao(ctx context.Context, in CGenDaoInput) (out *CGenDaoOutput, err error) {
	in.genItems = newCGenDaoInternalGenItems()
	if in.SqlDir != "" {
		// SQL file mode: generate from SQL DDL files without database connection.
		doGenDaoFromSQLFiles(ctx, in)
	} else if in.Link != "" {
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

// doGenDaoForArray implements the "gen dao" command for a single configuration entry.
// When index >= 0, it reads configuration from the array at that index.
// When index < 0, it uses the input as-is (for Link mode or single config mode).
// It performs the full generation pipeline: connect to DB, resolve tables,
// apply sharding patterns, and generate dao/table/do/entity files.
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
		err = gdb.AddConfigNode(tempGroup, gdb.ConfigNode{
			Link: in.Link,
		})
		if err != nil {
			mlog.Fatalf(`database configuration failed: %+v`, err)
		}
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
		inputTables := gstr.SplitAndTrim(in.Tables, ",")
		// Check if any table pattern contains wildcard characters.
		// https://github.com/gogf/gf/issues/4629
		var hasPattern bool
		for _, t := range inputTables {
			if containsWildcard(t) {
				hasPattern = true
				break
			}
		}
		if hasPattern {
			// Fetch all tables first, then filter by patterns.
			allTables, err := db.Tables(context.TODO())
			if err != nil {
				mlog.Fatalf("fetching tables failed: %+v", err)
			}
			tableNames = filterTablesByPatterns(allTables, inputTables)
		} else {
			// Use exact table names as before.
			tableNames = inputTables
		}
	} else {
		tableNames, err = db.Tables(context.TODO())
		if err != nil {
			mlog.Fatalf("fetching tables failed: %+v", err)
		}
	}
	// Table excluding.
	if in.TablesEx != "" {
		array := garray.NewStrArrayFrom(tableNames)
		for _, p := range gstr.SplitAndTrim(in.TablesEx, ",") {
			if containsWildcard(p) {
				// Use exact match with ^ and $ anchors for consistency with tables pattern.
				regPattern := "^" + patternToRegex(p) + "$"
				for _, v := range array.Clone().Slice() {
					if gregex.IsMatchString(regPattern, v) {
						array.RemoveValue(v)
					}
				}
			} else {
				array.RemoveValue(p)
			}
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
	var (
		newTableNames       = make([]string, len(tableNames))
		shardingNewTableSet = gset.NewStrSet()
	)
	// Sort sharding patterns by length descending, so that longer (more specific) patterns
	// are matched first. This prevents shorter patterns like "a_?" from incorrectly matching
	// tables that should match longer patterns like "a_b_?" or "a_c_?".
	// https://github.com/gogf/gf/issues/4603
	sortedShardingPatterns := make([]string, len(in.ShardingPattern))
	copy(sortedShardingPatterns, in.ShardingPattern)
	sort.Slice(sortedShardingPatterns, func(i, j int) bool {
		return len(sortedShardingPatterns[i]) > len(sortedShardingPatterns[j])
	})
	for i, tableName := range tableNames {
		newTableName := tableName
		for _, v := range removePrefixArray {
			newTableName = gstr.TrimLeftStr(newTableName, v, 1)
		}
		if len(sortedShardingPatterns) > 0 {
			for _, pattern := range sortedShardingPatterns {
				var (
					match      []string
					regPattern = gstr.Replace(pattern, "?", `(.+)`)
				)
				match, err = gregex.MatchString(regPattern, newTableName)
				if err != nil {
					mlog.Fatalf(`invalid sharding pattern "%s": %+v`, pattern, err)
				}
				if len(match) < 2 {
					continue
				}
				newTableName = gstr.Replace(pattern, "?", "")
				newTableName = gstr.Trim(newTableName, `_.-`)
				if shardingNewTableSet.Contains(newTableName) {
					tableNames[i] = ""
					break
				}
				// Add prefix to sharding table name, if not, the isSharding check would not match.
				shardingNewTableSet.Add(in.Prefix + newTableName)
				break
			}
		}
		newTableName = in.Prefix + newTableName
		if tableNames[i] != "" {
			// If shardingNewTableSet contains newTableName (tableName is empty), it should not be added to tableNames, make it empty and filter later.
			newTableNames[i] = newTableName
		}
	}
	tableNames = garray.NewStrArrayFrom(tableNames).FilterEmpty().Slice()
	newTableNames = garray.NewStrArrayFrom(newTableNames).FilterEmpty().Slice() // Filter empty table names. make sure that newTableNames and tableNames have the same length.
	in.genItems.Scale()

	// Dao: index and internal.
	generateDao(ctx, CGenDaoInternalInput{
		CGenDaoInput:     in,
		DB:               db,
		TableNames:       tableNames,
		NewTableNames:    newTableNames,
		ShardingTableSet: shardingNewTableSet,
	})
	// Table: table fields.
	generateTable(ctx, CGenDaoInternalInput{
		CGenDaoInput:     in,
		DB:               db,
		TableNames:       tableNames,
		NewTableNames:    newTableNames,
		ShardingTableSet: shardingNewTableSet,
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

// getImportPartContent analyzes the generated Go source code and builds the import block.
// It automatically detects usage of gtime.Time, time.Time, and gjson.Json in the source,
// and includes the corresponding import paths. Additional custom imports (from TypeMapping
// or FieldMapping) are appended and their dependencies are resolved via "go get" if needed.
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

// assignDefaultVar sets the default template variables for datetime strings
// used in generated file headers. The creation timestamp is only included
// when WithTime is enabled in the input configuration.
func assignDefaultVar(view *gview.View, in CGenDaoInternalInput) {
	var (
		tplCreatedAtDatetimeStr string
		tplDatetimeStr          = createdAt.String()
	)
	if in.WithTime {
		tplCreatedAtDatetimeStr = fmt.Sprintf(`Created at %s`, tplDatetimeStr)
	}
	view.Assigns(g.Map{
		tplVarDatetimeStr:          tplDatetimeStr,
		tplVarCreatedAtDatetimeStr: tplCreatedAtDatetimeStr,
	})
}

// sortFieldKeyForDao returns field names sorted by their Index in the TableField map.
// This preserves the original column order as defined in the database table schema.
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

// getTableFields retrieves table fields either from the pre-parsed TableFieldsMap (SQL file mode)
// or from the database connection. This abstracts the data source for generation functions.
func getTableFields(ctx context.Context, in CGenDaoInternalInput, tableName string) (map[string]*gdb.TableField, error) {
	if in.TableFieldsMap != nil {
		if fields, ok := in.TableFieldsMap[tableName]; ok {
			return fields, nil
		}
		return nil, fmt.Errorf("table '%s' not found in SQL files", tableName)
	}
	return in.DB.TableFields(ctx, tableName)
}

// getTemplateFromPathOrDefault returns the template content from the given file path.
// If the file path is empty or the file has no content, it falls back to the default template.
func getTemplateFromPathOrDefault(filePath string, def string) string {
	if filePath != "" {
		if contents := gfile.GetContents(filePath); contents != "" {
			return contents
		}
	}
	return def
}

// containsWildcard checks if the pattern contains wildcard characters (* or ?).
func containsWildcard(pattern string) bool {
	return gstr.Contains(pattern, "*") || gstr.Contains(pattern, "?")
}

// patternToRegex converts a wildcard pattern to a regex pattern.
// Wildcard characters: * matches any characters, ? matches single character.
func patternToRegex(pattern string) string {
	pattern = gstr.ReplaceByMap(pattern, map[string]string{
		"\r": "",
		"\n": "",
	})
	pattern = gstr.ReplaceByMap(pattern, map[string]string{
		"*": "\r",
		"?": "\n",
	})
	pattern = gregex.Quote(pattern)
	pattern = gstr.ReplaceByMap(pattern, map[string]string{
		"\r": ".*",
		"\n": ".",
	})
	return pattern
}

// filterTablesByPatterns filters tables by given patterns.
// Patterns support wildcard characters: * matches any characters, ? matches single character.
// https://github.com/gogf/gf/issues/4629
func filterTablesByPatterns(allTables []string, patterns []string) []string {
	var result []string
	matched := make(map[string]bool)
	allTablesSet := make(map[string]bool)
	for _, t := range allTables {
		allTablesSet[t] = true
	}
	for _, p := range patterns {
		if containsWildcard(p) {
			regPattern := "^" + patternToRegex(p) + "$"
			for _, table := range allTables {
				if !matched[table] && gregex.IsMatchString(regPattern, table) {
					result = append(result, table)
					matched[table] = true
				}
			}
		} else {
			// Exact table name, use direct string comparison.
			if !allTablesSet[p] {
				mlog.Printf(`table "%s" does not exist, skipped`, p)
				continue
			}
			if !matched[p] {
				result = append(result, p)
				matched[p] = true
			}
		}
	}
	return result
}

// doGenDaoFromSQLFiles implements the "gen dao" command for SQL file mode.
// It parses DDL SQL files to obtain table structures without requiring a database connection.
func doGenDaoFromSQLFiles(ctx context.Context, in CGenDaoInput) {
	if dirRealPath := gfile.RealPath(in.Path); dirRealPath == "" {
		mlog.Fatalf(`path "%s" does not exist`, in.Path)
	}
	if dirRealPath := gfile.RealPath(in.SqlDir); dirRealPath == "" {
		mlog.Fatalf(`SQL directory "%s" does not exist`, in.SqlDir)
	}

	dialect := SQLDialect(strings.ToLower(in.SqlType))
	tableNames, tableFieldsMap := ParseSQLFilesFromDir(in.SqlDir, dialect)

	removePrefixArray := gstr.SplitAndTrim(in.RemovePrefix, ",")

	// Table filtering by name patterns.
	if in.Tables != "" {
		inputTables := gstr.SplitAndTrim(in.Tables, ",")
		var hasPattern bool
		for _, t := range inputTables {
			if containsWildcard(t) {
				hasPattern = true
				break
			}
		}
		if hasPattern {
			tableNames = filterTablesByPatterns(tableNames, inputTables)
		} else {
			tableNames = inputTables
		}
	}

	// Table excluding.
	if in.TablesEx != "" {
		array := garray.NewStrArrayFrom(tableNames)
		for _, p := range gstr.SplitAndTrim(in.TablesEx, ",") {
			if containsWildcard(p) {
				regPattern := "^" + patternToRegex(p) + "$"
				for _, v := range array.Clone().Slice() {
					if gregex.IsMatchString(regPattern, v) {
						array.RemoveValue(v)
					}
				}
			} else {
				array.RemoveValue(p)
			}
		}
		tableNames = array.Slice()
	}

	// merge default typeMapping.
	if in.TypeMapping == nil {
		in.TypeMapping = defaultTypeMapping
	} else {
		for key, typeMapping := range defaultTypeMapping {
			if _, ok := in.TypeMapping[key]; !ok {
				in.TypeMapping[key] = typeMapping
			}
		}
	}

	// Process table names (prefix removal, sharding, etc.)
	var (
		newTableNames       = make([]string, len(tableNames))
		shardingNewTableSet = gset.NewStrSet()
	)
	sortedShardingPatterns := make([]string, len(in.ShardingPattern))
	copy(sortedShardingPatterns, in.ShardingPattern)
	sort.Slice(sortedShardingPatterns, func(i, j int) bool {
		return len(sortedShardingPatterns[i]) > len(sortedShardingPatterns[j])
	})
	for i, tableName := range tableNames {
		newTableName := tableName
		for _, v := range removePrefixArray {
			newTableName = gstr.TrimLeftStr(newTableName, v, 1)
		}
		if len(sortedShardingPatterns) > 0 {
			for _, pattern := range sortedShardingPatterns {
				var (
					match      []string
					regPattern = gstr.Replace(pattern, "?", `(.+)`)
					err        error
				)
				match, err = gregex.MatchString(regPattern, newTableName)
				if err != nil {
					mlog.Fatalf(`invalid sharding pattern "%s": %+v`, pattern, err)
				}
				if len(match) < 2 {
					continue
				}
				newTableName = gstr.Replace(pattern, "?", "")
				newTableName = gstr.Trim(newTableName, `_.-`)
				if shardingNewTableSet.Contains(newTableName) {
					tableNames[i] = ""
					break
				}
				shardingNewTableSet.Add(in.Prefix + newTableName)
				break
			}
		}
		newTableName = in.Prefix + newTableName
		if tableNames[i] != "" {
			newTableNames[i] = newTableName
		}
	}
	tableNames = garray.NewStrArrayFrom(tableNames).FilterEmpty().Slice()
	newTableNames = garray.NewStrArrayFrom(newTableNames).FilterEmpty().Slice()
	in.genItems.Scale()

	internalInput := CGenDaoInternalInput{
		CGenDaoInput:     in,
		DB:               nil,
		TableNames:       tableNames,
		NewTableNames:    newTableNames,
		ShardingTableSet: shardingNewTableSet,
		TableFieldsMap:   tableFieldsMap,
	}

	// Generate all files using the same flow as database mode.
	generateDao(ctx, internalInput)
	generateTable(ctx, internalInput)
	generateDo(ctx, internalInput)
	generateEntity(ctx, internalInput)

	in.genItems.SetClear(in.Clear)
}
