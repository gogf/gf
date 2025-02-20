// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gendao

import (
	"github.com/gogf/gf/v2/frame/g"
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
	CGenDaoBriefPath              = `directory path for generated files`
	CGenDaoBriefLink              = `database configuration, the same as the ORM configuration of GoFrame`
	CGenDaoBriefTables            = `generate models only for given tables, multiple table names separated with ','`
	CGenDaoBriefTablesEx          = `generate models excluding given tables, multiple table names separated with ','`
	CGenDaoBriefPrefix            = `add prefix for all table of specified link/database tables`
	CGenDaoBriefRemovePrefix      = `remove specified prefix of the table, multiple prefix separated with ','`
	CGenDaoBriefRemoveFieldPrefix = `remove specified prefix of the field, multiple prefix separated with ','`
	CGenDaoBriefStdTime           = `use time.Time from stdlib instead of gtime.Time for generated time/date fields of tables`
	CGenDaoBriefWithTime          = `add created time for auto produced go files`
	CGenDaoBriefGJsonSupport      = `use gJsonSupport to use *gjson.Json instead of string for generated json fields of tables`
	CGenDaoBriefImportPrefix      = `custom import prefix for generated go files`
	CGenDaoBriefDaoPath           = `directory path for storing generated dao files under path`
	CGenDaoBriefDoPath            = `directory path for storing generated do files under path`
	CGenDaoBriefEntityPath        = `directory path for storing generated entity files under path`
	CGenDaoBriefOverwriteDao      = `overwrite all dao files both inside/outside internal folder`
	CGenDaoBriefModelFile         = `custom file name for storing generated model content`
	CGenDaoBriefModelFileForDao   = `custom file name generating model for DAO operations like Where/Data. It's empty in default`
	CGenDaoBriefDescriptionTag    = `add comment to description tag for each field`
	CGenDaoBriefNoJsonTag         = `no json tag will be added for each field`
	CGenDaoBriefNoModelComment    = `no model comment will be added for each field`
	CGenDaoBriefClear             = `delete all generated go files that do not exist in database`
	CGenDaoBriefTypeMapping       = `custom local type mapping for generated struct attributes relevant to fields of table`
	CGenDaoBriefFieldMapping      = `custom local type mapping for generated struct attributes relevant to specific fields of table`
	CGenDaoBriefShardingPattern   = `sharding pattern for table name, e.g. "users_?" will be replace tables "users_001,users_002,..." to "users" dao`
	CGenDaoBriefGroup             = `
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

	tplVarTableName               = `TplTableName`
	tplVarTableNameCamelCase      = `TplTableNameCamelCase`
	tplVarTableNameCamelLowerCase = `TplTableNameCamelLowerCase`
	tplVarTableSharding           = `TplTableSharding`
	tplVarPackageImports          = `TplPackageImports`
	tplVarImportPrefix            = `TplImportPrefix`
	tplVarStructDefine            = `TplStructDefine`
	tplVarColumnDefine            = `TplColumnDefine`
	tplVarColumnNames             = `TplColumnNames`
	tplVarGroupName               = `TplGroupName`
	tplVarDatetimeStr             = `TplDatetimeStr`
	tplVarCreatedAtDatetimeStr    = `TplCreatedAtDatetimeStr`
	tplVarPackageName             = `TplPackageName`
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
		`CGenDaoBriefRemoveFieldPrefix`:  CGenDaoBriefRemoveFieldPrefix,
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
		`CGenDaoBriefTypeMapping`:        CGenDaoBriefTypeMapping,
		`CGenDaoBriefFieldMapping`:       CGenDaoBriefFieldMapping,
		`CGenDaoBriefShardingPattern`:    CGenDaoBriefShardingPattern,
		`CGenDaoBriefGroup`:              CGenDaoBriefGroup,
		`CGenDaoBriefJsonCase`:           CGenDaoBriefJsonCase,
		`CGenDaoBriefTplDaoIndexPath`:    CGenDaoBriefTplDaoIndexPath,
		`CGenDaoBriefTplDaoInternalPath`: CGenDaoBriefTplDaoInternalPath,
		`CGenDaoBriefTplDaoDoPathPath`:   CGenDaoBriefTplDaoDoPathPath,
		`CGenDaoBriefTplDaoEntityPath`:   CGenDaoBriefTplDaoEntityPath,
	})
}
