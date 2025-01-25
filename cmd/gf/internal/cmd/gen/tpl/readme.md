# 代码生成器设计文档

## 功能概述
基于数据库表结构，通过自定义模板生成Go代码的工具。

## 功能设计

生成流程：
1. 读取数据库表结构
2. 解析出表结构信息，包括表名、表注释、字段列表
3. 根据规则裁切表数据，生成模板数据
4. 根据模板生成代码

## 命令参数设计

```shell
$ gf gen tpl -h
USAGE
    gf gen tpl [OPTION]

OPTION
    -p, --path                  directory path for generated files
    -l, --link                  database configuration, the same as the ORM configuration of GoFrame
    -t, --tables                generate templates only for given tables, multiple table names separated with ','
    -x, --tablesEx              generate templates excluding given tables, multiple table names separated with ','
    -g, --group                 specifying the configuration group name of database for generated ORM instance,
                                it's not necessary and the default value is "default"
    -f, --prefix                add prefix for all table of specified link/database tables
    -r, --removePrefix          remove specified prefix of the table, multiple prefix separated with ','
    -rf, --removeFieldPrefix    remove specified prefix of the field, multiple prefix separated with ','
    -j, --jsonCase              generated json tag case for model struct, cases are as follows:
                                | Case            | Example            |
                                |---------------- |--------------------|
                                | Camel           | AnyKindOfString    |
                                | CamelLower      | anyKindOfString    | default
                                | Snake           | any_kind_of_string |
                                | SnakeScreaming  | ANY_KIND_OF_STRING |
                                | SnakeFirstUpper | rgb_code_md5       |
                                | Kebab           | any-kind-of-string |
                                | KebabScreaming  | ANY-KIND-OF-STRING |
    -i, --importPrefix          custom import prefix for generated go files
    -t1, --tplPath              template file path for custom template
    -s, --stdTime               use time.Time from stdlib instead of gtime.Time for generated time/date fields of tables
    -w, --withTime              add created time for auto produced go files
    -n, --gJsonSupport          use gJsonSupport to use *gjson.Json instead of string for generated json fields of
                                tables
    -v, --overwrite             overwrite all template files
    -c, --descriptionTag        add comment to description tag for each field
    -k, --noJsonTag             no json tag will be added for each field
    -m, --noModelComment        no model comment will be added for each field
    -a, --clear                 delete all generated template files that do not exist in database
    -y, --typeMapping           custom local type mapping for generated struct attributes relevant to fields of table
    -fm, --fieldMapping         custom local type mapping for generated struct attributes relevant to specific fields of
                                table
    -h, --help                  more information about this command

EXAMPLE
    gf gen tpl
    gf gen tpl -l "mysql:root:12345678@tcp(127.0.0.1:3306)/test"
    gf gen tpl -p ./template -g user-center -t user,user_detail,user_login
    gf gen tpl -r user_

CONFIGURATION SUPPORT
    Options are also supported by configuration file.
    It's suggested using configuration file instead of command line arguments making producing.
    The configuration node name is "gfcli.gen.tpl", which also supports multiple databases, for example(config.yaml):
    gfcli:
      gen:
        tpl:
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
```



## 表结构信息
### 表信息
1. 表名
2. 表注释
3. 字段列表

### 字段信息
1. 字段名
2. 类型
3. 对应的 go 类型
4. 是否主键
5. 是否唯一键
6. 备注
7. 默认值
8. 是否自增

