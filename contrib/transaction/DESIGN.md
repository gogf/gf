# GF Seata 集成架构设计文档

## 一、GDB 事务功能分析

### 1.1 核心架构

GF 的数据库事务管理主要由以下组件构成:

```
┌──────────────────────────────────────────────────────────┐
│                         应用层                            │
│  db.Transaction() / Model.Transaction()                 │
└───────────────────┬──────────────────────────────────────┘
                    │
┌───────────────────▼──────────────────────────────────────┐
│                    事务传播层                             │
│  TransactionWithOptions(PropagationXXX)                 │
│  - PropagationNested (默认)                              │
│  - PropagationRequired                                  │
│  - PropagationSupports                                  │
│  - 其他4种传播类型                                        │
└───────────────────┬──────────────────────────────────────┘
                    │
┌───────────────────▼──────────────────────────────────────┐
│                    事务管理层                             │
│  TXCore (implements TX interface)                       │
│  - Begin/Commit/Rollback                                │
│  - SavePoint/RollbackTo (嵌套事务)                       │
│  - Context管理(WithTX/TXFromCtx)                        │
└───────────────────┬──────────────────────────────────────┘
                    │
┌───────────────────▼──────────────────────────────────────┐
│                   SQL执行层                               │
│  DoCommit() - 统一执行入口                                │
│  - SqlTypeBegin                                         │
│  - SqlTypeExecContext                                   │
│  - SqlTypeQueryContext                                  │
│  - SqlTypeTXCommit                                      │
│  - SqlTypeTXRollback                                    │
└───────────────────┬──────────────────────────────────────┘
                    │
┌───────────────────▼──────────────────────────────────────┐
│                  数据库驱动层                             │
│  Link Interface (dbLink / txLink)                       │
│  - QueryContext()                                       │
│  - ExecContext()                                        │
│  - PrepareContext()                                     │
└──────────────────────────────────────────────────────────┘
```

### 1.2 关键特性

#### 1.2.1 Context传递机制
- 事务对象通过 Context 传递
- 使用 `transactionKeyForContext(group)` 作为 Context key
- 支持按数据库分组隔离事务

```go
// 注入事务到 Context
ctx = context.WithValue(ctx, transactionKeyForContext(group), tx)

// 从 Context 提取事务
tx := TXFromCtx(ctx, group)
```

#### 1.2.2 嵌套事务支持
- 使用 SAVEPOINT 实现嵌套事务
- `transactionCount` 跟踪嵌套深度
- `transactionKeyForNestedPoint()` 生成保存点名称

#### 1.2.3 事务传播机制
类似 Spring 的 7 种传播类型:

| 传播类型 | 说明 | GF实现 |
|---------|------|--------|
| NESTED | 嵌套事务(默认) | SAVEPOINT |
| REQUIRED | 加入已有事务 | 共用 sql.Tx |
| SUPPORTS | 支持但不强制 | 有则用,无则普通执行 |
| REQUIRES_NEW | 挂起当前,新建事务 | WithoutTX + Begin |
| NOT_SUPPORTED | 挂起事务,非事务执行 | WithoutTX |
| MANDATORY | 必须在事务中 | 检查并报错 |
| NEVER | 不能在事务中 | 检查并报错 |

### 1.3 核心钩子点

集成 Seata 的关键钩子点:

1. **DoCommit()** - 最底层执行点,拦截所有SQL
2. **DoFilter()** - SQL执行前的过滤器
3. **DoExec()/DoQuery()** - 具体SQL执行点
4. **Model.Hook()** - ORM层钩子

## 二、Seata AT模式原理

### 2.1 AT模式核心流程

```
┌─────────────┐
│   TM (事务   │
│   管理器)    │
└──────┬──────┘
       │ 1. Begin Global TX
       ▼
┌─────────────┐
│   TC (事务   │
│   协调器)    │
└──────┬──────┘
       │ 2. Return XID
       ▼
┌─────────────┐        ┌─────────────┐        ┌─────────────┐
│  RM1 (分支   │        │  RM2 (分支   │        │  RM3 (分支   │
│  事务1)      │        │  事务2)      │        │  事务3)      │
└──────┬──────┘        └──────┬──────┘        └──────┬──────┘
       │                      │                      │
       │ 3. Register Branch   │                      │
       │ 4. Before Image      │                      │
       │ 5. Execute SQL       │                      │
       │ 6. After Image       │                      │
       │ 7. Insert UndoLog    │                      │
       │ 8. Local Commit      │                      │
       │                      │                      │
       ▼                      ▼                      ▼
    [本地DB1]              [本地DB2]              [本地DB3]
```

### 2.2 AT模式关键技术

#### 2.2.1 Before/After Image
在执行业务SQL前后,记录数据快照:

```sql
-- Before Image (UPDATE前查询)
SELECT * FROM account WHERE id = 1;
-- Result: {id: 1, balance: 1000}

-- Business SQL
UPDATE account SET balance = balance - 100 WHERE id = 1;

-- After Image (UPDATE后查询)
SELECT * FROM account WHERE id = 1;
-- Result: {id: 1, balance: 900}
```

#### 2.2.2 Undo Log结构
```json
{
  "table_name": "account",
  "sql_type": "UPDATE",
  "before_image": {
    "id": 1,
    "balance": 1000
  },
  "after_image": {
    "id": 1,
    "balance": 900
  },
  "pk_columns": ["id"]
}
```

#### 2.2.3 全局锁机制
- 在一阶段本地事务提交前获取全局锁
- 全局锁由TC管理,防止脏写
- 本地事务提交后释放本地锁,但保留全局锁直到全局事务结束

## 三、GF + Seata 集成方案

### 3.1 整体架构

```
┌──────────────────────────────────────────────────────────┐
│                      应用业务层                           │
│  seata.GlobalTransaction(ctx, func...)                  │
└───────────────────┬──────────────────────────────────────┘
                    │
┌───────────────────▼──────────────────────────────────────┐
│                   Seata包装层                             │
│  SeataDB (包装 gdb.DB)                                   │
│  - Model() 注入Seata上下文                               │
│  - Transaction() 创建分支事务                            │
│  - DoCommit() 拦截SQL执行                                │
└───────────────────┬──────────────────────────────────────┘
                    │
┌───────────────────▼──────────────────────────────────────┐
│                  事务协调层                               │
│  SeataTX (包装 gdb.TX)                                   │
│  - registerBranch() 注册分支事务                         │
│  - Commit() 提交并报告状态                               │
│  - Rollback() 回滚并清理UndoLog                          │
└───────────────────┬──────────────────────────────────────┘
                    │
┌───────────────────▼──────────────────────────────────────┐
│                 UndoLog管理层                             │
│  UndoLogManager                                         │
│  - BuildUndoLog() 构建回滚日志                           │
│  - InsertUndoLog() 插入回滚日志                          │
│  - ExecuteUndoLog() 执行回滚                             │
└───────────────────┬──────────────────────────────────────┘
                    │
┌───────────────────▼──────────────────────────────────────┐
│                   GF原生层                                │
│  gdb.DB / gdb.TX / gdb.Core                             │
└───────────────────┬──────────────────────────────────────┘
                    │
┌───────────────────▼──────────────────────────────────────┐
│                 Seata客户端                               │
│  TM (事务管理器) + RM (资源管理器)                        │
│  - 与 TC 通信                                            │
│  - 全局事务协调                                           │
│  - 分支事务注册                                           │
└──────────────────────────────────────────────────────────┘
```

### 3.2 核心组件设计

#### 3.2.1 SeataDB
包装 gdb.DB,实现DB接口:

```go
type SeataDB struct {
    gdb.DB
    config      Config
    resourceMgr rm.ResourceManager
    undoLogMgr  *UndoLogManager
}

// 重写关键方法
func (s *SeataDB) DoCommit(ctx context.Context, in gdb.DoCommitInput) 
    (out gdb.DoCommitOutput, err error) {
    
    // 检查是否在全局事务中
    xid, inGlobalTx := XIDFromContext(ctx)
    
    if inGlobalTx {
        // 拦截SQL执行
        switch in.Type {
        case gdb.SqlTypeExecContext:
            // DML: 记录UndoLog
            return s.doCommitWithUndoLog(ctx, xid, in)
        case gdb.SqlTypeTXCommit:
            // 分支提交
            return s.doCommitBranch(ctx, xid, in)
        }
    }
    
    // 正常流程
    return s.DB.DoCommit(ctx, in)
}
```

#### 3.2.2 SeataTX
包装 gdb.TX,管理分支事务:

```go
type SeataTX struct {
    gdb.TX
    seataDB  *SeataDB
    xid      string
    branchID int64
}

func (tx *SeataTX) Commit() error {
    // 1. 提交本地事务
    if err := tx.TX.Commit(); err != nil {
        return err
    }
    
    // 2. 报告分支状态 (异步)
    tx.BranchReport(ctx, rm.BranchStatusPhaseOneDone)
    
    return nil
}
```

#### 3.2.3 UndoLogManager
管理回滚日志:

```go
type UndoLogManager struct {
    db        gdb.DB
    tableName string
}

func (m *UndoLogManager) BuildUndoLog(
    xid string, beforeImage, afterImage map[string]any, sql string,
) *UndoLog {
    // 构建UndoLog记录
}

func (m *UndoLogManager) ExecuteUndoLog(ctx context.Context, undoLog *UndoLog) error {
    // 根据SQL类型生成回滚SQL
    switch undoLog.SQLType {
    case "INSERT":
        // 生成DELETE
    case "UPDATE":
        // 生成UPDATE恢复原值
    case "DELETE":
        // 生成INSERT
    }
}
```

### 3.3 关键流程实现

#### 3.3.1 全局事务开启

```go
func GlobalTransaction(ctx context.Context, f func(context.Context) error) error {
    // 1. 开启全局事务
    globalTx := tm.GetGlobalTransactionManager().CreateGlobalTransaction(ctx)
    globalTx.Begin(ctx, "tx-name", 30000)
    
    // 2. 获取XID并注入Context
    xid := globalTx.GetXid()
    ctx = ContextWithXID(ctx, xid)
    
    // 3. 执行业务逻辑
    err := f(ctx)
    
    // 4. 提交或回滚
    if err != nil {
        globalTx.Rollback(ctx)
        return err
    }
    
    return globalTx.Commit(ctx)
}
```

#### 3.3.2 分支事务执行

```go
func (s *SeataDB) branchTransaction(ctx context.Context, xid string, 
    f func(context.Context, gdb.TX) error) error {
    
    // 1. 开启本地事务
    tx, _ := s.DB.Begin(ctx)
    
    // 2. 创建Seata事务包装
    seataTx := &SeataTX{TX: tx, xid: xid}
    
    // 3. 注册分支事务
    branchID, _ := seataTx.registerBranch(ctx)
    seataTx.branchID = branchID
    
    // 4. 执行业务逻辑
    err := f(ctx, seataTx)
    
    // 5. 提交/回滚
    if err != nil {
        seataTx.Rollback()
        return err
    }
    
    return seataTx.Commit()
}
```

#### 3.3.3 UndoLog记录

```go
func (s *SeataDB) doCommitWithUndoLog(ctx context.Context, xid string, 
    in gdb.DoCommitInput) (out gdb.DoCommitOutput, err error) {
    
    // 1. 查询Before Image
    beforeImage := s.queryBeforeImage(ctx, in.Sql, in.Args)
    
    // 2. 执行SQL
    out, err = s.DB.DoCommit(ctx, in)
    if err != nil {
        return
    }
    
    // 3. 查询After Image
    afterImage := s.queryAfterImage(ctx, in.Sql, in.Args)
    
    // 4. 构建并插入UndoLog
    undoLog := s.undoLogMgr.BuildUndoLog(xid, beforeImage, afterImage, in.Sql)
    s.undoLogMgr.InsertUndoLog(ctx, undoLog)
    
    return
}
```

### 3.4 SQL解析器集成

为了准确提取Before/After Image,需要SQL解析器:

```go
type SQLParser interface {
    // 解析SQL获取表名
    GetTableName(sql string) string
    
    // 解析SQL获取主键条件
    GetPrimaryKeyCondition(sql string, args []any) string
    
    // 判断SQL类型
    GetSQLType(sql string) string // INSERT/UPDATE/DELETE
}

// 可选方案:
// 1. 使用 github.com/pingcap/parser (TiDB SQL Parser)
// 2. 使用正则表达式(简单场景)
// 3. 集成 Seata-Go 的 SQL Parser
```

## 四、使用示例

### 4.1 基础使用

```go
// 初始化
db := g.DB()
seataDB, _ := seata.WrapDB(db, seata.Config{
    ApplicationID:  "order-service",
    TxServiceGroup: "default_tx_group",
    ResourceID:     "order-db",
})

// 全局事务
err := seata.GlobalTransaction(ctx, func(ctx context.Context) error {
    // 创建订单
    _, err := g.DB().Model("order").Ctx(ctx).Insert(order)
    if err != nil {
        return err
    }
    
    // 调用其他服务(自动加入全局事务)
    stockService.Reduce(ctx, productID, quantity)
    accountService.Deduct(ctx, userID, amount)
    
    return nil
})
```

### 4.2 兼容原生API

```go
// Seata模式下也支持GF原生事务API
err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
    // 在全局事务中,这会自动注册为分支事务
    _, err := tx.Model("order").Insert(order)
    return err
})
```

## 五、性能优化

### 5.1 UndoLog压缩
- 配置压缩阈值(如64KB)
- 使用gzip压缩大型UndoLog

### 5.2 批量操作优化
- 批量INSERT使用单条UndoLog
- 合并多个小事务

### 5.3 异步提交
- 一阶段提交后异步报告TC
- 后台定时清理已提交的UndoLog

## 六、监控与运维

### 6.1 指标采集
- 全局事务数量
- 分支事务数量
- 事务成功率/失败率
- UndoLog大小统计

### 6.2 日志
- 全局事务日志(XID关联)
- 分支事务日志
- UndoLog操作日志

## 七、限制与注意事项

1. **数据库要求**: 必须支持本地ACID事务
2. **主键要求**: 业务表必须有主键
3. **DDL限制**: 全局事务中不支持DDL操作
4. **性能影响**: AT模式增加网络开销和日志记录
5. **隔离级别**: 建议使用READ_COMMITTED隔离级别

## 八、未来规划

1. **TCC模式支持**: 实现Try-Confirm-Cancel模式
2. **SAGA模式支持**: 长事务场景
3. **XA模式支持**: 强一致性场景
4. **性能优化**: 减少网络开销,优化UndoLog记录
5. **监控完善**: 集成Prometheus/Grafana
