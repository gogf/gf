// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gendao

import (
	"sort"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/os/gfile"

	"github.com/gogf/gf/cmd/gf/v2/internal/utility/mlog"
)

// SQLDialect defines supported SQL dialect types.
type SQLDialect string

const (
	SQLDialectMySQL      SQLDialect = "mysql"
	SQLDialectPgSQL      SQLDialect = "pgsql"
	SQLDialectMSSQL      SQLDialect = "mssql"
	SQLDialectOracle     SQLDialect = "oracle"
	SQLDialectSQLite     SQLDialect = "sqlite"
	SQLDialectClickHouse SQLDialect = "clickhouse"
)

// SQLStatementType identifies the type of a DDL statement.
type SQLStatementType int

const (
	SQLStatementUnknown     SQLStatementType = iota
	SQLStatementCreateTable                  // CREATE TABLE
	SQLStatementAlterTable                   // ALTER TABLE
	SQLStatementDropTable                    // DROP TABLE
	SQLStatementRenameTable                  // RENAME TABLE / ALTER TABLE ... RENAME TO
	SQLStatementComment                      // COMMENT ON COLUMN / sp_addextendedproperty
)

// SQLParser is the interface for parsing SQL DDL files into table field definitions.
// Each parser must implement CREATE TABLE parsing. ALTER TABLE, DROP TABLE, and
// comment handling are optional and can be delegated to the common layer.
type SQLParser interface {
	// ParseCreateTable parses a single CREATE TABLE statement and returns table name and fields.
	ParseCreateTable(stmt string) (tableName string, fields map[string]*gdb.TableField, err error)

	// ParseAlterTable parses a single ALTER TABLE statement and applies changes to existing tables.
	// Returns the affected table name.
	ParseAlterTable(stmt string, tables map[string]map[string]*gdb.TableField) error

	// ParseComment parses a comment statement (COMMENT ON COLUMN / sp_addextendedproperty)
	// and applies the comment to the corresponding field.
	ParseComment(stmt string, tables map[string]map[string]*gdb.TableField)
}

// GetSQLParser returns the appropriate SQL parser for the given dialect.
func GetSQLParser(dialect SQLDialect) SQLParser {
	switch dialect {
	case SQLDialectMySQL:
		return &MySQLParser{}
	case SQLDialectPgSQL:
		return &PgSQLParser{}
	case SQLDialectMSSQL:
		return &MSSQLParser{}
	case SQLDialectOracle:
		return &OracleParser{}
	case SQLDialectSQLite:
		return &SQLiteParser{}
	default:
		return nil
	}
}

// ParseSQLFilesFromDir parses all .sql files from the given directory using the specified
// dialect parser. Files are processed in sorted order (by filename) to ensure correct
// incremental migration order: CREATE TABLE first, then ALTER TABLE modifications.
func ParseSQLFilesFromDir(sqlDir string, dialect SQLDialect) (
	tableNames []string,
	tableFieldsMap map[string]map[string]*gdb.TableField,
) {
	parser := GetSQLParser(dialect)
	if parser == nil {
		mlog.Fatalf(`unsupported SQL dialect "%s"`, dialect)
	}

	sqlFiles, err := gfile.ScanDirFile(sqlDir, "*.sql", true)
	if err != nil {
		mlog.Fatalf(`scanning SQL directory "%s" failed: %+v`, sqlDir, err)
	}
	if len(sqlFiles) == 0 {
		mlog.Fatalf(`no .sql files found in directory "%s"`, sqlDir)
	}

	// Sort files by name to ensure deterministic migration order.
	// This supports naming conventions like:
	//   V001_create_tables.sql, V002_add_columns.sql, V003_modify_columns.sql
	//   001_init.sql, 002_alter.sql
	//   2024-01-01_create.sql, 2024-01-15_alter.sql
	sort.Strings(sqlFiles)

	tableFieldsMap = make(map[string]map[string]*gdb.TableField)

	for _, sqlFile := range sqlFiles {
		content := gfile.GetContents(sqlFile)
		if content == "" {
			continue
		}
		err := processSQL(parser, content, tableFieldsMap)
		if err != nil {
			mlog.Fatalf(`parsing SQL file "%s" failed: %+v`, sqlFile, err)
		}
	}

	for tableName := range tableFieldsMap {
		tableNames = append(tableNames, tableName)
	}
	sort.Strings(tableNames)
	return
}

// processSQL processes all SQL statements in a single file content,
// dispatching each statement to the appropriate handler based on its type.
func processSQL(
	parser SQLParser,
	sql string,
	tables map[string]map[string]*gdb.TableField,
) error {
	statements := splitSQLStatements(sql)
	for _, stmt := range statements {
		stmtType := classifyStatement(stmt)
		switch stmtType {
		case SQLStatementCreateTable:
			tableName, fields, err := parser.ParseCreateTable(stmt)
			if err != nil {
				return err
			}
			if tableName != "" && len(fields) > 0 {
				tables[tableName] = fields
			}

		case SQLStatementAlterTable:
			err := parser.ParseAlterTable(stmt, tables)
			if err != nil {
				return err
			}

		case SQLStatementDropTable:
			applyDropTable(stmt, tables)

		case SQLStatementRenameTable:
			applyRenameTable(stmt, tables)

		case SQLStatementComment:
			parser.ParseComment(stmt, tables)
		}
	}
	return nil
}

// classifyStatement determines the type of a SQL DDL statement.
func classifyStatement(stmt string) SQLStatementType {
	upper := strings.ToUpper(strings.TrimSpace(stmt))
	// Remove leading block comments
	for strings.HasPrefix(upper, "/*") {
		end := strings.Index(upper, "*/")
		if end < 0 {
			break
		}
		upper = strings.TrimSpace(upper[end+2:])
	}

	words := strings.Fields(upper)
	if len(words) < 2 {
		return SQLStatementUnknown
	}

	switch words[0] {
	case "CREATE":
		// CREATE [TEMPORARY|TEMP] TABLE
		for _, w := range words[1:] {
			if w == "TABLE" {
				return SQLStatementCreateTable
			}
			if w != "TEMPORARY" && w != "TEMP" && w != "GLOBAL" && w != "LOCAL" &&
				w != "UNLOGGED" {
				break
			}
		}

	case "ALTER":
		if words[1] == "TABLE" {
			// Check if it's ALTER TABLE ... RENAME TO
			if strings.Contains(upper, "RENAME TO") || strings.Contains(upper, "RENAME AS") {
				return SQLStatementRenameTable
			}
			return SQLStatementAlterTable
		}

	case "DROP":
		if words[1] == "TABLE" {
			return SQLStatementDropTable
		}

	case "RENAME":
		// RENAME TABLE old TO new (MySQL syntax)
		if words[1] == "TABLE" {
			return SQLStatementRenameTable
		}

	case "COMMENT":
		// COMMENT ON COLUMN / COMMENT ON TABLE
		if len(words) >= 3 && words[1] == "ON" {
			return SQLStatementComment
		}

	case "EXEC", "EXECUTE":
		// EXEC sp_addextendedproperty (MSSQL comments)
		if strings.Contains(upper, "SP_ADDEXTENDEDPROPERTY") &&
			strings.Contains(upper, "MS_DESCRIPTION") {
			return SQLStatementComment
		}
	}

	return SQLStatementUnknown
}

// applyDropTable removes a table from the tables map.
// Handles: DROP TABLE [IF EXISTS] [schema.]table_name
func applyDropTable(stmt string, tables map[string]map[string]*gdb.TableField) {
	upper := strings.ToUpper(strings.TrimSpace(stmt))
	upper = strings.TrimPrefix(upper, "DROP")
	upper = strings.TrimSpace(upper)
	upper = strings.TrimPrefix(upper, "TABLE")
	upper = strings.TrimSpace(upper)
	if strings.HasPrefix(upper, "IF EXISTS") {
		upper = strings.TrimPrefix(upper, "IF EXISTS")
	}

	remaining := stmt[len(stmt)-len(upper):]
	remaining = strings.TrimSpace(remaining)

	// May be comma-separated list: DROP TABLE t1, t2, t3
	for _, name := range strings.Split(remaining, ",") {
		name = strings.TrimSpace(name)
		parts := strings.Split(name, ".")
		tableName := unquoteIdentifier(parts[len(parts)-1])
		delete(tables, tableName)
	}
}

// applyRenameTable renames a table in the tables map.
// Handles:
//   - RENAME TABLE old TO new (MySQL)
//   - ALTER TABLE old RENAME TO new (PostgreSQL, SQLite, Oracle)
func applyRenameTable(stmt string, tables map[string]map[string]*gdb.TableField) {
	upper := strings.ToUpper(strings.TrimSpace(stmt))
	words := strings.Fields(stmt)
	upperWords := strings.Fields(upper)

	if upperWords[0] == "RENAME" && len(upperWords) >= 5 && upperWords[1] == "TABLE" {
		// RENAME TABLE old_name TO new_name
		oldName := unquoteIdentifier(words[2])
		newName := unquoteIdentifier(words[4])
		if fields, ok := tables[oldName]; ok {
			tables[newName] = fields
			delete(tables, oldName)
		}
	} else if upperWords[0] == "ALTER" && len(upperWords) >= 6 && upperWords[1] == "TABLE" {
		// ALTER TABLE old_name RENAME TO new_name
		oldName := unquoteIdentifier(words[2])
		for i, w := range upperWords {
			if w == "RENAME" && i+2 < len(upperWords) &&
				(upperWords[i+1] == "TO" || upperWords[i+1] == "AS") {
				newName := unquoteIdentifier(words[i+2])
				if fields, ok := tables[oldName]; ok {
					tables[newName] = fields
					delete(tables, oldName)
				}
				return
			}
		}
	}
}

// parseAlterTableCommon provides common ALTER TABLE parsing logic that works for
// most SQL dialects. Individual parsers can call this or override with dialect-specific logic.
//
// Supported operations:
//   - ADD [COLUMN] column_name type [constraints]
//   - DROP [COLUMN] column_name
//   - MODIFY [COLUMN] column_name type [constraints]         (MySQL, Oracle)
//   - ALTER [COLUMN] column_name TYPE type / SET / DROP       (PostgreSQL)
//   - CHANGE [COLUMN] old_name new_name type [constraints]    (MySQL)
//   - ADD PRIMARY KEY (col1, col2, ...)
//   - DROP PRIMARY KEY
//   - RENAME COLUMN old_name TO new_name
func parseAlterTableCommon(
	stmt string,
	tables map[string]map[string]*gdb.TableField,
	columnParser func(def string, index int) (*gdb.TableField, error),
) error {
	upper := strings.ToUpper(strings.TrimSpace(stmt))
	words := strings.Fields(stmt)
	upperWords := strings.Fields(upper)

	if len(upperWords) < 4 || upperWords[0] != "ALTER" || upperWords[1] != "TABLE" {
		return nil
	}

	tableName := unquoteIdentifier(words[2])
	fields, exists := tables[tableName]
	if !exists {
		// Table not yet created, skip
		return nil
	}

	// The rest after ALTER TABLE tableName
	restIdx := 3
	if restIdx >= len(upperWords) {
		return nil
	}

	// Process the ALTER TABLE actions. Some dialects allow multiple actions separated by commas
	// but we handle one action at a time for simplicity and split multi-action later.
	return processAlterAction(upperWords, words, restIdx, fields, tableName, tables, columnParser)
}

func processAlterAction(
	upperWords, words []string,
	startIdx int,
	fields map[string]*gdb.TableField,
	tableName string,
	tables map[string]map[string]*gdb.TableField,
	columnParser func(def string, index int) (*gdb.TableField, error),
) error {
	if startIdx >= len(upperWords) {
		return nil
	}

	action := upperWords[startIdx]
	switch action {
	case "ADD":
		return processAlterAdd(upperWords, words, startIdx+1, fields, columnParser)

	case "DROP":
		processAlterDrop(upperWords, words, startIdx+1, fields)

	case "MODIFY":
		// MODIFY [COLUMN] column_name type [constraints] (MySQL, Oracle)
		return processAlterModify(upperWords, words, startIdx+1, fields, columnParser)

	case "CHANGE":
		// CHANGE [COLUMN] old_name new_name type [constraints] (MySQL)
		return processAlterChange(upperWords, words, startIdx+1, fields, columnParser)

	case "ALTER":
		// ALTER [COLUMN] column_name ... (PostgreSQL: SET DEFAULT, DROP DEFAULT, SET NOT NULL, etc.)
		processAlterColumn(upperWords, words, startIdx+1, fields)

	case "RENAME":
		// RENAME COLUMN old_name TO new_name
		processAlterRenameColumn(upperWords, words, startIdx+1, fields)
	}

	return nil
}

// processAlterAdd handles ALTER TABLE ... ADD [COLUMN] ... or ADD PRIMARY KEY ...
func processAlterAdd(
	upperWords, words []string,
	idx int,
	fields map[string]*gdb.TableField,
	columnParser func(def string, index int) (*gdb.TableField, error),
) error {
	if idx >= len(upperWords) {
		return nil
	}

	// Skip optional COLUMN keyword
	colIdx := idx
	if upperWords[colIdx] == "COLUMN" {
		colIdx++
	}

	// ADD PRIMARY KEY (col1, col2)
	if upperWords[idx] == "PRIMARY" || upperWords[idx] == "UNIQUE" ||
		upperWords[idx] == "INDEX" || upperWords[idx] == "KEY" ||
		upperWords[idx] == "CONSTRAINT" || upperWords[idx] == "FOREIGN" {
		if strings.Contains(strings.Join(upperWords[idx:], " "), "PRIMARY KEY") {
			fullStmt := strings.Join(words[idx:], " ")
			pkCols := findPrimaryKeysFromConstraints([]string{fullStmt})
			for _, pkCol := range pkCols {
				if f, ok := fields[pkCol]; ok {
					f.Key = "PRI"
				}
			}
		}
		return nil
	}

	if colIdx >= len(words) {
		return nil
	}

	// ADD [COLUMN] column_def ...
	// Build the column definition from remaining words
	def := strings.Join(words[colIdx:], " ")
	nextIndex := getNextFieldIndex(fields)
	field, err := columnParser(def, nextIndex)
	if err != nil {
		return nil // skip unparseable
	}
	if field != nil {
		fields[field.Name] = field
	}
	return nil
}

// processAlterDrop handles ALTER TABLE ... DROP [COLUMN] column_name or DROP PRIMARY KEY
func processAlterDrop(upperWords, words []string, idx int, fields map[string]*gdb.TableField) {
	if idx >= len(upperWords) {
		return
	}

	// DROP PRIMARY KEY
	if upperWords[idx] == "PRIMARY" {
		for _, f := range fields {
			if f.Key == "PRI" {
				f.Key = ""
			}
		}
		return
	}

	// DROP [COLUMN] column_name [CASCADE|RESTRICT]
	colIdx := idx
	if upperWords[colIdx] == "COLUMN" {
		colIdx++
	}
	if colIdx >= len(words) {
		return
	}

	colName := unquoteIdentifier(words[colIdx])
	delete(fields, colName)

	// Reindex remaining fields
	reindexFields(fields)
}

// processAlterModify handles ALTER TABLE ... MODIFY [COLUMN] column_name type [constraints]
func processAlterModify(
	upperWords, words []string,
	idx int,
	fields map[string]*gdb.TableField,
	columnParser func(def string, index int) (*gdb.TableField, error),
) error {
	if idx >= len(upperWords) {
		return nil
	}

	colIdx := idx
	if upperWords[colIdx] == "COLUMN" {
		colIdx++
	}
	if colIdx >= len(words) {
		return nil
	}

	colName := unquoteIdentifier(words[colIdx])
	def := strings.Join(words[colIdx:], " ")

	existingIndex := 0
	if existing, ok := fields[colName]; ok {
		existingIndex = existing.Index
	}

	field, err := columnParser(def, existingIndex)
	if err != nil {
		return nil
	}
	if field != nil {
		// Preserve the original index
		if existing, ok := fields[colName]; ok {
			field.Index = existing.Index
		}
		fields[field.Name] = field
	}
	return nil
}

// processAlterChange handles ALTER TABLE ... CHANGE [COLUMN] old_name new_name type [constraints]
func processAlterChange(
	upperWords, words []string,
	idx int,
	fields map[string]*gdb.TableField,
	columnParser func(def string, index int) (*gdb.TableField, error),
) error {
	if idx >= len(upperWords) {
		return nil
	}

	colIdx := idx
	if upperWords[colIdx] == "COLUMN" {
		colIdx++
	}
	if colIdx+1 >= len(words) {
		return nil
	}

	oldName := unquoteIdentifier(words[colIdx])
	// New definition starts from the new column name
	def := strings.Join(words[colIdx+1:], " ")

	existingIndex := 0
	if existing, ok := fields[oldName]; ok {
		existingIndex = existing.Index
	}

	field, err := columnParser(def, existingIndex)
	if err != nil {
		return nil
	}
	if field != nil {
		// Remove old field
		delete(fields, oldName)
		if existing, ok := fields[oldName]; ok {
			field.Index = existing.Index
		} else {
			field.Index = existingIndex
		}
		fields[field.Name] = field
	}
	return nil
}

// processAlterColumn handles ALTER TABLE ... ALTER [COLUMN] column_name ...
// PostgreSQL style: SET DEFAULT, DROP DEFAULT, SET NOT NULL, DROP NOT NULL, TYPE
func processAlterColumn(upperWords, words []string, idx int, fields map[string]*gdb.TableField) {
	if idx >= len(upperWords) {
		return
	}

	colIdx := idx
	if upperWords[colIdx] == "COLUMN" {
		colIdx++
	}
	if colIdx >= len(words) {
		return
	}

	colName := unquoteIdentifier(words[colIdx])
	field, ok := fields[colName]
	if !ok {
		return
	}

	actionIdx := colIdx + 1
	if actionIdx >= len(upperWords) {
		return
	}

	switch upperWords[actionIdx] {
	case "SET":
		if actionIdx+1 < len(upperWords) {
			switch upperWords[actionIdx+1] {
			case "NOT":
				// SET NOT NULL
				if actionIdx+2 < len(upperWords) && upperWords[actionIdx+2] == "NULL" {
					field.Null = false
				}
			case "DEFAULT":
				// SET DEFAULT value
				if actionIdx+2 < len(words) {
					defaultVal, _ := extractDefaultValue("DEFAULT " + strings.Join(words[actionIdx+2:], " "))
					field.Default = defaultVal
				}
			case "DATA":
				// SET DATA TYPE type_name (PostgreSQL)
				if actionIdx+2 < len(upperWords) && upperWords[actionIdx+2] == "TYPE" &&
					actionIdx+3 < len(words) {
					field.Type = strings.Join(words[actionIdx+3:], " ")
				}
			}
		}
	case "DROP":
		if actionIdx+1 < len(upperWords) {
			switch upperWords[actionIdx+1] {
			case "NOT":
				// DROP NOT NULL
				if actionIdx+2 < len(upperWords) && upperWords[actionIdx+2] == "NULL" {
					field.Null = true
				}
			case "DEFAULT":
				// DROP DEFAULT
				field.Default = nil
			}
		}
	case "TYPE":
		// TYPE new_type (PostgreSQL: ALTER COLUMN col TYPE varchar(200))
		if actionIdx+1 < len(words) {
			// Collect the type, which may include USING clause
			typeParts := make([]string, 0)
			for j := actionIdx + 1; j < len(words); j++ {
				if strings.ToUpper(words[j]) == "USING" {
					break
				}
				typeParts = append(typeParts, words[j])
			}
			if len(typeParts) > 0 {
				field.Type = strings.Join(typeParts, " ")
			}
		}
	}
}

// processAlterRenameColumn handles ALTER TABLE ... RENAME COLUMN old TO new
func processAlterRenameColumn(upperWords, words []string, idx int, fields map[string]*gdb.TableField) {
	if idx >= len(upperWords) {
		return
	}

	colIdx := idx
	if upperWords[colIdx] == "COLUMN" {
		colIdx++
	}
	if colIdx+2 >= len(words) {
		return
	}

	// RENAME [COLUMN] old_name TO new_name
	oldName := unquoteIdentifier(words[colIdx])
	// Find "TO"
	for i := colIdx + 1; i < len(upperWords)-1; i++ {
		if upperWords[i] == "TO" {
			newName := unquoteIdentifier(words[i+1])
			if field, ok := fields[oldName]; ok {
				field.Name = newName
				fields[newName] = field
				delete(fields, oldName)
			}
			return
		}
	}
}

// getNextFieldIndex returns the next available field index.
func getNextFieldIndex(fields map[string]*gdb.TableField) int {
	maxIndex := -1
	for _, f := range fields {
		if f.Index > maxIndex {
			maxIndex = f.Index
		}
	}
	return maxIndex + 1
}

// reindexFields re-assigns sequential indices to all fields after a deletion.
func reindexFields(fields map[string]*gdb.TableField) {
	type indexedField struct {
		name  string
		index int
	}
	sorted := make([]indexedField, 0, len(fields))
	for name, f := range fields {
		sorted = append(sorted, indexedField{name: name, index: f.Index})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].index < sorted[j].index
	})
	for i, sf := range sorted {
		fields[sf.name].Index = i
	}
}

// splitSQLStatements splits SQL content into individual statements by semicolons,
// handling quoted strings and parentheses properly.
func splitSQLStatements(sql string) []string {
	var (
		statements []string
		current    strings.Builder
		inSingle   bool
		inDouble   bool
		inBlock    bool // block comment
		depth      int
		prev       byte
	)
	for i := 0; i < len(sql); i++ {
		ch := sql[i]
		switch {
		case inBlock:
			current.WriteByte(ch)
			if ch == '/' && prev == '*' {
				inBlock = false
			}
		case ch == '/' && i+1 < len(sql) && sql[i+1] == '*' && !inSingle && !inDouble:
			inBlock = true
			current.WriteByte(ch)
		case ch == '-' && i+1 < len(sql) && sql[i+1] == '-' && !inSingle && !inDouble:
			// Line comment - skip to end of line
			for i < len(sql) && sql[i] != '\n' {
				i++
			}
		case ch == '\'' && !inDouble:
			inSingle = !inSingle
			current.WriteByte(ch)
		case ch == '"' && !inSingle:
			inDouble = !inDouble
			current.WriteByte(ch)
		case ch == '(' && !inSingle && !inDouble:
			depth++
			current.WriteByte(ch)
		case ch == ')' && !inSingle && !inDouble:
			depth--
			current.WriteByte(ch)
		case ch == ';' && !inSingle && !inDouble && depth == 0:
			stmt := strings.TrimSpace(current.String())
			if stmt != "" {
				statements = append(statements, stmt)
			}
			current.Reset()
		default:
			current.WriteByte(ch)
		}
		prev = ch
	}
	// Last statement without semicolon
	if stmt := strings.TrimSpace(current.String()); stmt != "" {
		statements = append(statements, stmt)
	}
	return statements
}

// unquoteIdentifier removes quotes from SQL identifiers.
// Handles: `name`, "name", [name], 'name'
func unquoteIdentifier(name string) string {
	name = strings.TrimSpace(name)
	if len(name) < 2 {
		return name
	}
	switch {
	case name[0] == '`' && name[len(name)-1] == '`':
		return name[1 : len(name)-1]
	case name[0] == '"' && name[len(name)-1] == '"':
		return name[1 : len(name)-1]
	case name[0] == '[' && name[len(name)-1] == ']':
		return name[1 : len(name)-1]
	}
	return name
}

// extractTableName extracts the table name from a CREATE TABLE statement header.
// It handles: CREATE TABLE name, CREATE TABLE IF NOT EXISTS name,
// CREATE TABLE schema.name, etc.
func extractTableName(header string) string {
	header = strings.TrimSpace(header)
	// Remove CREATE [TEMPORARY] TABLE [IF NOT EXISTS]
	upper := strings.ToUpper(header)
	upper = strings.TrimPrefix(upper, "CREATE")
	upper = strings.TrimSpace(upper)
	if strings.HasPrefix(upper, "TEMPORARY") || strings.HasPrefix(upper, "TEMP") {
		idx := strings.Index(upper, "TABLE")
		if idx >= 0 {
			upper = upper[idx:]
		}
	}
	upper = strings.TrimPrefix(upper, "TABLE")
	upper = strings.TrimSpace(upper)
	if strings.HasPrefix(upper, "IF NOT EXISTS") {
		upper = strings.TrimPrefix(upper, "IF NOT EXISTS")
	}

	// Now get the actual name from original string at same offset
	remaining := header[len(header)-len(upper):]
	remaining = strings.TrimSpace(remaining)

	// Handle schema.table
	parts := strings.Split(remaining, ".")
	tableName := parts[len(parts)-1]
	tableName = strings.TrimSpace(tableName)

	return unquoteIdentifier(tableName)
}

// splitColumns splits the column definitions part of a CREATE TABLE statement
// into individual column/constraint definitions, properly handling nested parentheses.
func splitColumns(body string) []string {
	var (
		result  []string
		current strings.Builder
		depth   int
		inStr   bool
		quote   byte
	)
	for i := 0; i < len(body); i++ {
		ch := body[i]
		switch {
		case inStr:
			current.WriteByte(ch)
			if ch == quote && (i+1 >= len(body) || body[i+1] != quote) {
				inStr = false
			}
		case ch == '\'' || ch == '"':
			inStr = true
			quote = ch
			current.WriteByte(ch)
		case ch == '(':
			depth++
			current.WriteByte(ch)
		case ch == ')':
			depth--
			current.WriteByte(ch)
		case ch == ',' && depth == 0:
			if s := strings.TrimSpace(current.String()); s != "" {
				result = append(result, s)
			}
			current.Reset()
		default:
			current.WriteByte(ch)
		}
	}
	if s := strings.TrimSpace(current.String()); s != "" {
		result = append(result, s)
	}
	return result
}

// isConstraintKeyword checks if the given word starts a table-level constraint
// (not a column definition).
func isConstraintKeyword(word string) bool {
	upper := strings.ToUpper(word)
	switch upper {
	case "PRIMARY", "UNIQUE", "INDEX", "KEY", "CHECK", "FOREIGN", "CONSTRAINT",
		"CLUSTERED", "NONCLUSTERED", "SPATIAL", "FULLTEXT":
		return true
	}
	return false
}

// findPrimaryKeysFromConstraints scans constraint definitions for PRIMARY KEY
// and returns the column names that form the primary key.
func findPrimaryKeysFromConstraints(columnDefs []string) []string {
	var pkColumns []string
	for _, def := range columnDefs {
		upper := strings.ToUpper(strings.TrimSpace(def))
		if !strings.Contains(upper, "PRIMARY KEY") {
			continue
		}
		// Extract column names from PRIMARY KEY (col1, col2, ...)
		idx := strings.Index(upper, "PRIMARY KEY")
		rest := def[idx+len("PRIMARY KEY"):]
		rest = strings.TrimSpace(rest)
		// Skip optional CLUSTERED/NONCLUSTERED keyword (MSSQL).
		upperRest := strings.ToUpper(rest)
		if strings.HasPrefix(upperRest, "CLUSTERED") {
			rest = strings.TrimSpace(rest[len("CLUSTERED"):])
		} else if strings.HasPrefix(upperRest, "NONCLUSTERED") {
			rest = strings.TrimSpace(rest[len("NONCLUSTERED"):])
		}
		if len(rest) > 0 && rest[0] == '(' {
			end := strings.Index(rest, ")")
			if end > 0 {
				cols := rest[1:end]
				for _, col := range strings.Split(cols, ",") {
					col = strings.TrimSpace(col)
					// Remove ASC/DESC
					parts := strings.Fields(col)
					if len(parts) > 0 {
						pkColumns = append(pkColumns, unquoteIdentifier(parts[0]))
					}
				}
			}
		}
	}
	return pkColumns
}

// extractBodyAndTrailing splits CREATE TABLE ... (...) ... into body and trailing parts.
// It returns the content inside the outermost parentheses and anything after.
func extractBodyAndTrailing(sql string) (body, trailing string, ok bool) {
	// Find the first '(' that starts the column definitions
	depth := 0
	startIdx := -1
	endIdx := -1
	inStr := false
	var quote byte
	for i := 0; i < len(sql); i++ {
		ch := sql[i]
		if inStr {
			if ch == quote && (i+1 >= len(sql) || sql[i+1] != quote) {
				inStr = false
			}
			continue
		}
		if ch == '\'' || ch == '"' {
			inStr = true
			quote = ch
			continue
		}
		if ch == '(' {
			if depth == 0 {
				startIdx = i
			}
			depth++
		} else if ch == ')' {
			depth--
			if depth == 0 {
				endIdx = i
				break
			}
		}
	}
	if startIdx < 0 || endIdx < 0 {
		return "", "", false
	}
	body = sql[startIdx+1 : endIdx]
	trailing = strings.TrimSpace(sql[endIdx+1:])
	return body, trailing, true
}

// extractDefaultValue extracts the default value string from a column definition fragment.
// Returns the default value and the remaining string after the default clause.
func extractDefaultValue(s string) (defaultVal any, rest string) {
	upper := strings.ToUpper(strings.TrimSpace(s))
	if !strings.HasPrefix(upper, "DEFAULT") {
		return nil, s
	}
	s = strings.TrimSpace(s[7:]) // skip "DEFAULT"
	if len(s) == 0 {
		return nil, ""
	}

	// NULL
	if strings.HasPrefix(strings.ToUpper(s), "NULL") {
		return nil, strings.TrimSpace(s[4:])
	}

	// Quoted string
	if s[0] == '\'' {
		end := 1
		for end < len(s) {
			if s[end] == '\'' {
				if end+1 < len(s) && s[end+1] == '\'' {
					end += 2
					continue
				}
				val := s[1:end]
				val = strings.ReplaceAll(val, "''", "'")
				return val, strings.TrimSpace(s[end+1:])
			}
			end++
		}
		return s[1:], ""
	}

	// Parenthesized expression like (getdate()), ((0))
	if s[0] == '(' {
		depth := 0
		for i := 0; i < len(s); i++ {
			if s[i] == '(' {
				depth++
			} else if s[i] == ')' {
				depth--
				if depth == 0 {
					return s[:i+1], strings.TrimSpace(s[i+1:])
				}
			}
		}
		return s, ""
	}

	// Unquoted value (number, keyword like CURRENT_TIMESTAMP, etc.)
	parts := strings.Fields(s)
	if len(parts) > 0 {
		val := parts[0]
		// Remove trailing comma if any
		val = strings.TrimRight(val, ",")
		rest = strings.TrimSpace(s[len(parts[0]):])
		return val, rest
	}
	return nil, ""
}

// mysqlTokenize performs simple tokenization of a column definition string,
// respecting quoted identifiers and parenthesized type parameters.
func mysqlTokenize(def string) []string {
	var (
		tokens  []string
		current strings.Builder
		inStr   bool
		quote   byte
		depth   int
	)
	for i := 0; i < len(def); i++ {
		ch := def[i]
		switch {
		case inStr:
			current.WriteByte(ch)
			if ch == quote && (i+1 >= len(def) || def[i+1] != quote) {
				inStr = false
			}
		case ch == '\'' || ch == '"' || ch == '`':
			if current.Len() == 0 || depth > 0 {
				inStr = true
				quote = ch
			}
			current.WriteByte(ch)
		case ch == '(':
			depth++
			current.WriteByte(ch)
		case ch == ')':
			depth--
			current.WriteByte(ch)
		case (ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r') && depth == 0:
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
		default:
			current.WriteByte(ch)
		}
	}
	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}
	return tokens
}
