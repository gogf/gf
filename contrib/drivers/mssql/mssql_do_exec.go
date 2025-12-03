package mssql

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

const (
	insertPrefixDefault = "INSERT INTO"
	insertPrefixIgnore  = "INSERT IGNORE INTO"

	fieldExtraIdentity     = "IDENTITY"
	fieldKeyPrimary        = "PRI"
	outputKeyword          = "OUTPUT"
	insertedObjectName     = "INSERTED"
	affectCountExpression  = " 1 as AffectCount"
	affectCountFieldName   = "AffectCount"
	lastInsertIdFieldAlias = "ID"
	insertValuesMarker     = ") VALUES" // find the position of the string "VALUES" in the INSERT SQL statement to embed output code for retrieving the last inserted ID
)

// DoExec commits the sql string and its arguments to underlying driver
// through given link object and returns the execution result.
func (d *Driver) DoExec(ctx context.Context, link gdb.Link, sqlStr string, args ...interface{}) (result sql.Result, err error) {
	// Transaction checks.
	if link == nil {
		if tx := gdb.TXFromCtx(ctx, d.GetGroup()); tx != nil {
			// Firstly, check and retrieve transaction link from context.
			link = &txLinkMssql{tx.GetSqlTX()}
		} else if link, err = d.Core.MasterLink(); err != nil {
			// Or else it creates one from master node.
			return nil, err
		}
	} else if !link.IsTransaction() {
		// If current link is not transaction link, it checks and retrieves transaction from context.
		if tx := gdb.TXFromCtx(ctx, d.GetGroup()); tx != nil {
			link = &txLinkMssql{tx.GetSqlTX()}
		}
	}

	// SQL filtering.
	sqlStr, args = d.Core.FormatSqlBeforeExecuting(sqlStr, args)
	sqlStr, args, err = d.DoFilter(ctx, link, sqlStr, args)
	if err != nil {
		return nil, err
	}

	if !(strings.HasPrefix(sqlStr, insertPrefixDefault) || strings.HasPrefix(sqlStr, insertPrefixIgnore)) {
		return d.Core.DoExec(ctx, link, sqlStr, args)
	}
	// find the first pos
	pos := strings.Index(sqlStr, insertValuesMarker)

	table := d.GetTableNameFromSql(sqlStr)
	outPutSql := d.GetInsertOutputSql(ctx, table)
	// rebuild sql add output
	var (
		sqlValueBefore = sqlStr[:pos+1]
		sqlValueAfter  = sqlStr[pos+1:]
	)

	sqlStr = fmt.Sprintf("%s%s%s", sqlValueBefore, outPutSql, sqlValueAfter)

	// fmt.Println("sql str:", sqlStr)
	// Link execution.
	var out gdb.DoCommitOutput
	out, err = d.DoCommit(ctx, gdb.DoCommitInput{
		Link:          link,
		Sql:           sqlStr,
		Args:          args,
		Stmt:          nil,
		Type:          gdb.SqlTypeQueryContext,
		IsTransaction: link.IsTransaction(),
	})
	if err != nil {
		return &InsertResult{lastInsertId: 0, rowsAffected: 0, err: err}, err
	}
	var (
		lId int64 // last insert id
	)
	stdSqlResult := out.Records
	if len(stdSqlResult) == 0 {
		err = gerror.WrapCode(gcode.CodeDbOperationError, gerror.New("affectcount is zero"), `sql.Result.RowsAffected failed`)
		return &InsertResult{lastInsertId: 0, rowsAffected: 0, err: err}, err
	}
	// get affect count from the number of returned rows
	aCount := int64(len(stdSqlResult))
	// get last_insert_id from the first returned row
	lId = stdSqlResult[0].GMap().GetVar(lastInsertIdFieldAlias).Int64()

	return &InsertResult{lastInsertId: lId, rowsAffected: aCount}, err
}

// GetTableNameFromSql get table name from sql statement
// It handles table string like:
// "user"
// "user u"
// "DbLog.dbo.user",
// "user as u".
func (d *Driver) GetTableNameFromSql(sqlStr string) (table string) {
	// INSERT INTO "ip_to_id"("ip") OUTPUT  1 as AffectCount,INSERTED.id as ID VALUES(?)
	leftChars, rightChars := d.GetChars()
	trimStr := leftChars + rightChars + "[] "
	pattern := "INTO(.+?)\\("
	regCompile := regexp.MustCompile(pattern)
	tableInfo := regCompile.FindStringSubmatch(sqlStr)
	//get the first one. after the first it may be content of the value, it's not table name.
	table = tableInfo[1]
	table = strings.Trim(table, " ")
	if strings.Contains(table, ".") {
		tmpAry := strings.Split(table, ".")
		// the last one is tablename
		table = tmpAry[len(tmpAry)-1]
	} else if strings.Contains(table, "as") || strings.Contains(table, " ") {
		tmpAry := strings.Split(table, "as")
		if len(tmpAry) < 2 {
			tmpAry = strings.Split(table, " ")
		}
		// get the first one
		table = tmpAry[0]
	}
	table = strings.Trim(table, trimStr)
	return table
}

// txLink is used to implement interface Link for TX.
type txLinkMssql struct {
	*sql.Tx
}

// IsTransaction returns if current Link is a transaction.
func (l *txLinkMssql) IsTransaction() bool {
	return true
}

// IsOnMaster checks and returns whether current link is operated on master node.
// Note that, transaction operation is always operated on master node.
func (l *txLinkMssql) IsOnMaster() bool {
	return true
}

// InsertResult instance of sql.Result
type InsertResult struct {
	lastInsertId int64
	rowsAffected int64
	err          error
}

func (r *InsertResult) LastInsertId() (int64, error) {
	return r.lastInsertId, r.err
}

func (r *InsertResult) RowsAffected() (int64, error) {
	return r.rowsAffected, r.err
}

// GetInsertOutputSql  gen get last_insert_id code
func (m *Driver) GetInsertOutputSql(ctx context.Context, table string) string {
	fds, errFd := m.GetDB().TableFields(ctx, table)
	if errFd != nil {
		return ""
	}
	extraSqlAry := make([]string, 0)
	extraSqlAry = append(extraSqlAry, fmt.Sprintf(" %s %s", outputKeyword, affectCountExpression))
	incrNo := 0
	if len(fds) > 0 {
		for _, fd := range fds {
			// has primary key and is auto-increment
			if fd.Extra == fieldExtraIdentity && fd.Key == fieldKeyPrimary && !fd.Null {
				incrNoStr := ""
				if incrNo == 0 { // fixed first field named id, convenient to get
					incrNoStr = fmt.Sprintf(" as %s", lastInsertIdFieldAlias)
				}

				extraSqlAry = append(extraSqlAry, fmt.Sprintf("%s.%s%s", insertedObjectName, fd.Name, incrNoStr))
				incrNo++
			}
			// fmt.Printf("null:%t name:%s key:%s k:%s \n", fd.Null, fd.Name, fd.Key, k)
		}
	}
	return strings.Join(extraSqlAry, ",")
	// sql example:INSERT INTO "ip_to_id"("ip") OUTPUT  1 as AffectCount,INSERTED.id as ID VALUES(?)
}
