## Context

GoFrame's `database/gredis` package defines an abstract `Adapter` interface with group-based Redis command access (Generic, Hash, String, List, Set, SortedSet, PubSub, Script). The concrete implementation lives in `contrib/nosql/redis` and wraps `go-redis/v9`'s `redis.UniversalClient`.

The current Adapter interface (`AdapterOperation`) provides four methods: `Do`, `Conn`, `Close`, `Client`. There is no Pipeline, Transaction, or WATCH support at the interface level. Users who need batch execution must call `Client()` and type-assert to `redis.UniversalClient` to access go-redis's `Pipeline()` / `TxPipeline()` / `Watch()` — which breaks the driver-agnostic design and couples application code to a specific driver.

For Cluster mode, `Scan` only covers a single node's keyspace, and multi-key `Del` triggers `CROSSSLOT` errors. These are infrastructure-level issues that belong in the framework, not in each application.

## Goals / Non-Goals

**Goals:**
- Add Pipeline support to the `AdapterOperation` interface with full group-style type safety.
- Add Transaction (MULTI/EXEC) and optimistic locking (WATCH) support.
- Add Cluster-safe `ScanAll` to `IGroupGeneric`.
- Implement Cluster-safe `Del` in the driver (per-key fallback in Cluster mode).
- All new interfaces must be driver-agnostic (no go-redis types in `database/gredis`).
- Achieve ≥ 80% test coverage for newly added code.

**Non-Goals:**
- Batch convenience methods (BatchHGetAll, BatchHMSet, etc.) — these are application-level patterns built on top of Pipeline.
- Key-prefix namespacing, JSON serialization wrappers, index management, or CacheManager — these are application concerns, not framework primitives.
- PubSub or Script groups in the Pipeline interface — these operations do not fit the pipeline execution model.
- Configuration validation — GoFrame's `gcfg` already provides validation mechanisms.

## Decisions

### 1. Cmd as future-result container

Pipeline commands are queued locally and executed in batch. A command queued at time T1 cannot return its result until `Exec()` is called at time T2. We introduce a `Cmd` type to hold this deferred result:

```go
type Cmd struct {
    val *gvar.Var  // nil until Exec populates it
    err error
}
```

After `Exec()`, users call `cmd.Result()` to retrieve the populated value and error. This mirrors go-redis's `Cmder` pattern while staying within GoFrame's `gvar.Var` ecosystem.

### 2. Full Group-style Pipeliner interface

The `Pipeliner` interface mirrors the existing group structure but with `*Cmd` return types. Each pipeline group interface (`IPipelineGroupHash`, etc.) has the same method names as the corresponding `IGroupHash`, differing only in return type:

- Regular: `HSet(ctx, key, fields) (int64, error)`
- Pipeline: `HSet(ctx, key, fields) *Cmd`

PubSub and Script groups are excluded from the Pipeliner because they do not fit the pipeline execution model (PubSub is push-based; scripts use EVAL which is itself a single command).

### 3. Pipeline/TxPipeline return Pipeliner directly (not callback-based)

We chose `Pipeline(ctx) Pipeliner` over a callback-based `Pipeline(ctx, func(pipe) error)` because:
- The direct-return pattern is simpler to use and test.
- Users retain control over when to call `Exec()`.
- go-redis uses the same pattern (`client.Pipeline()` returns `redis.Pipeliner`).

### 4. Tx interface embeds Pipeliner

`Watch` receives a callback with a `Tx` argument. `Tx` embeds `Pipeliner`, so inside the callback the user can queue commands on the transaction. The `Tx` interface is kept thin to avoid over-abstracting go-redis's `*redis.Tx`.

### 5. ScanAll added to IGroupGeneric

`ScanAll` abstracts the common "iterate all matching keys" pattern. In Cluster mode, the driver implementation uses `ForEachMaster` to aggregate results across all master nodes. In standalone/sentinel mode, it loops `Scan` until cursor returns to 0.

### 6. Cluster-safe Del implemented in driver only

`Del`'s interface signature stays the same (`Del(ctx, keys...) (int64, error)`). The driver detects Cluster mode at runtime and switches to per-key deletion to avoid CROSSSLOT. No interface change needed — this is purely a driver implementation detail.

### 7. Files and estimated scope

| File | Action | Est. Lines |
|---|---|---|
| `database/gredis/gredis_cmd.go` | New — Cmd type | ~50 |
| `database/gredis/gredis_pipeline.go` | New — Pipeliner + pipeline group interfaces | ~250 |
| `database/gredis/gredis_adapter.go` | Modified — extend AdapterOperation | +15 |
| `database/gredis/gredis_redis.go` | Modified — add Pipeline/TxPipeline/Watch | +40 |
| `database/gredis/gredis_redis_group_generic.go` | Modified — add ScanAll | +10 |
| `contrib/nosql/redis/redis_pipeline.go` | New — Pipeliner implementation | ~500 |
| `contrib/nosql/redis/redis.go` | Modified — implement Pipeline/TxPipeline/Watch | +30 |
| `contrib/nosql/redis/redis_group_generic.go` | Modified — ScanAll + Cluster-safe Del | +80 |
| Test files (both modules) | New + modified | ~600 |
| **Total** | | **~1600** |

## Risks / Trade-offs

- **Breaking change to AdapterOperation interface** — Any custom adapter implementing `AdapterOperation` must add `Pipeline`, `TxPipeline`, and `Watch` methods. This is acceptable because v2 allows breaking changes within the major version, and these are fundamental Redis capabilities that any real adapter should support. Mitigation: provide clear migration documentation.

- **Large interface surface for Pipeline groups** — ~100 methods across 6 pipeline group interfaces. This is inherent to the "full Group style" decision. The implementation is mechanical (delegate each method to go-redis's pipeliner), and code generation could help if maintainers prefer it.

- **Cmd type is new to users** — Users must learn the Cmd pattern (queue → Exec → Result). This is standard in Redis client libraries (go-redis, jedis, redis-py) and well-understood. Documentation and examples will ease adoption.

- **Cluster-safe Del changes semantics silently** — In Cluster mode, a multi-key `Del` that previously errored with CROSSSLOT will now succeed (per-key). This is strictly better behavior but changes the error contract. Mitigation: log a debug-level message when the fallback path is taken.
