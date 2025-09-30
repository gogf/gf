# 标签配置使用指南

## 功能概述

`gf gen tpl` 现在支持灵活的标签配置,可以选择性地为生成的结构体字段添加 `omitempty` 或其他自定义标签。

## 配置选项一览

| 选项 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `jsonOmitempty` | bool | false | 为所有字段添加 omitempty |
| `jsonOmitemptyAuto` | bool | false | 仅为可空字段自动添加 omitempty |
| `withOrmTag` | bool | true | 是否添加 orm 标签 |
| `descriptionTag` | bool | false | 是否添加 description 标签 |
| `noJsonTag` | bool | false | 是否禁用 JSON 标签 |
| `fieldMapping.tags` | map | - | 字段级自定义标签 |

## 配置方式

### 1. 全局开关 - `jsonOmitempty`

为所有字段的 JSON 标签添加 `omitempty`:

```yaml
gfcli:
  gen:
    tpl:
      jsonOmitempty: true
```

**生成结果:**
```go
type User struct {
    ID       int    `json:"id,omitempty" orm:"id" description:"用户ID"`
    Name     string `json:"name,omitempty" orm:"name" description:"用户名"`
    Email    string `json:"email,omitempty" orm:"email" description:"邮箱"`
}
```

---

### 2. 智能判断 - `jsonOmitemptyAuto` (推荐)

仅为可空字段自动添加 `omitempty`:

```yaml
gfcli:
  gen:
    tpl:
      jsonOmitemptyAuto: true
```

**假设数据库表结构:**
```sql
CREATE TABLE user (
    id INT NOT NULL,
    name VARCHAR(50) NOT NULL,
    email VARCHAR(100) NULL,  -- 可空字段
    age INT NULL                -- 可空字段
);
```

**生成结果:**
```go
type User struct {
    ID    int    `json:"id" orm:"id" description:"用户ID"`
    Name  string `json:"name" orm:"name" description:"用户名"`
    Email string `json:"email,omitempty" orm:"email" description:"邮箱"`  // 自动添加
    Age   int    `json:"age,omitempty" orm:"age" description:"年龄"`     // 自动添加
}
```

---

### 3. ORM 标签控制 - `withOrmTag`

控制是否添加 orm 标签 (默认启用):

```yaml
gfcli:
  gen:
    tpl:
      withOrmTag: false  # 不生成 orm 标签
```

**生成结果:**
```go
type User struct {
    ID    int    `json:"id" description:"用户ID"`          // 没有 orm 标签
    Name  string `json:"name" description:"用户名"`        // 没有 orm 标签
    Email string `json:"email" description:"邮箱"`         // 没有 orm 标签
}
```

---

### 4. 字段级精确控制 - `fieldMapping`

针对特定字段自定义标签 (优先级最高):

```yaml
gfcli:
  gen:
    tpl:
      fieldMapping:
        user.password:
          type: string
          tags:
            json: "-"  # 不序列化

        user.email:
          type: string
          tags:
            json: "email,omitempty"
            validate: "required,email"
            binding: "required"

        user.status:
          type: int
          tags:
            json: "status,omitempty"
            validate: "oneof=0 1 2"
            example: "1"
```

**生成结果:**
```go
type User struct {
    Password string `json:"-" orm:"password" description:"密码"`
    Email    string `binding:"required" json:"email,omitempty" validate:"required,email" description:"邮箱"`
    Status   int    `example:"1" json:"status,omitempty" validate:"oneof=0 1 2" description:"状态"`
}
```

---

## 常见标签示例

### validate 标签 (gin validator)

```yaml
fieldMapping:
  user.email:
    tags:
      validate: "required,email"

  user.age:
    tags:
      validate: "gte=0,lte=150"

  user.password:
    tags:
      validate: "required,min=8,max=32"
```

### binding 标签 (gin binding)

```yaml
fieldMapping:
  user.name:
    tags:
      binding: "required"

  user.email:
    tags:
      binding: "required,email"
```

### swagger 文档标签

```yaml
fieldMapping:
  user.id:
    tags:
      example: "1"
      description: "用户唯一标识"

  user.status:
    tags:
      example: "1"
      enums: "0,1,2"
```

### 多个自定义标签组合

```yaml
fieldMapping:
  user.email:
    type: string
    tags:
      json: "email,omitempty"
      validate: "required,email"
      binding: "required"
      example: "user@example.com"
      description: "用户邮箱地址"
```

---

## 配置优先级

标签配置的优先级从高到低:

1. **fieldMapping.tags** - 字段级自定义标签 (优先级最高)
2. **jsonOmitempty** - 全局 omitempty 开关
3. **jsonOmitemptyAuto** - 智能判断可空字段
4. **默认行为** - 不添加 omitempty

---

## 完整配置示例

```yaml
gfcli:
  gen:
    tpl:
      link: "mysql:root:12345678@tcp(127.0.0.1:3306)/test"
      path: "./output"
      tplPath: "./templates"
      jsonCase: "CamelLower"
      importPrefix: "github.com/example/project"

      # 全局配置
      jsonOmitemptyAuto: true  # 可空字段自动添加 omitempty
      withOrmTag: true         # 添加 orm 标签 (默认)
      descriptionTag: true     # 添加 description 标签

      # 类型映射
      typeMapping:
        decimal:
          type: decimal.Decimal
          import: github.com/shopspring/decimal

      # 字段级配置
      fieldMapping:
        user.password:
          type: string
          tags:
            json: "-"

        user.email:
          type: string
          tags:
            json: "email,omitempty"
            validate: "required,email"
            binding: "required"

        order.total_amount:
          type: decimal.Decimal
          import: github.com/shopspring/decimal
          tags:
            json: "totalAmount,omitempty"
            validate: "gt=0"
```

---

## 命令行使用

```bash
# 使用配置文件
gf gen tpl

# 命令行参数
gf gen tpl -tp ./templates -p ./output -ja -wo
# -ja: jsonOmitemptyAuto
# -wo: withOrmTag

# 组合使用
gf gen tpl -l "mysql:root:pass@tcp(127.0.0.1:3306)/db" -tp ./tpl -ja -c -wo
```

---

## 模板中使用

如果你需要在自定义模板中使用标签功能:

```go
// entity.tpl
type {{.table.NameCaseCamel}} struct { {{range $i,$v := .table.Fields}}
    {{$v.NameCaseCamel}} {{$v.LocalType}} {{$v.BuildTags $.tagInput}} // {{$v.Comment}}{{end}}
}
```

或者分别使用单个标签方法:

```go
type {{.table.NameCaseCamel}} struct { {{range $i,$v := .table.Fields}}
    {{$v.NameCaseCamel}} {{$v.LocalType}} `json:"{{$v.JsonTag $.tagInput.JsonOmitempty $.tagInput.JsonOmitemptyAuto}}" orm:"{{$v.OrmTag}}"` // {{$v.Comment}}{{end}}
}
```

---

## 注意事项

1. **字段名格式**: `fieldMapping` 中的 key 格式为 `表名.字段名`,使用数据库中的实际字段名 (非驼峰)

2. **标签顺序**: 自定义标签会按字母顺序排列,确保生成结果一致

3. **特殊字符**: 如果标签值包含双引号,会自动转义

4. **DO 文件**: DO 文件 (model/do) 只保留 description 标签,不包含 JSON/ORM 标签

5. **兼容性**: 与现有的 `typeMapping` 和 `fieldMapping` 完全兼容

6. **默认值**: `withOrmTag` 默认为 `true`,如果不需要 orm 标签,需要显式设置为 `false`
