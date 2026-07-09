## Why

GoFrame's `database/gredis` abstract layer currently exposes only single-command operations (Do, group-based commands like HSet/Get/etc.). It lacks support for Pipeline (batch command execution in a single network round-trip), Transaction (MULTI/EXEC atomic execution), and optimistic locking (WATCH). Users who need these capabilities must resort to type-asserting the raw `Client()` return to `redis.UniversalClient` and manually constructing pipelines with go-redis types — breaking the driver-agnostic abstraction.

Additionally, the concrete `contrib/nosql/redis` driver has no Cluster-aware SCAN (global key iteration across all master nodes) and no Cluster-safe multi-key DEL (naive multi-key DEL triggers CROSSSLOT errors). These are common infrastructure-level pitfalls that every Cluster-mode user encounters.

## What Changes

### Pipeline & Transaction Support
- Add a `Cmd` type as a future-result container for pipelined commands.
- Add a `Pipeliner` interface with full group-style command access (Generic, Hash, String, List, Set, SortedSet), mirroring existing group interfaces but returning `*Cmd` instead of `(value, error)`.
- Extend `AdapterOperation` interface with `Pipeline()`, `TxPipeline()`, and `Watch()`.
- Add `Tx` interface for transaction context (embeds `Pipeliner`).
- Implement the full `Pipeliner` in `contrib/nosql/redis` by delegating to go-redis's `redis.Pipeliner`.

### Cluster Enhancement
- Add `ScanAll` method to `IGroupGeneric` — Cluster-safe global key scan that iterates all master nodes.
- Implement Cluster-safe `Del` in the driver — auto-detects Cluster mode and falls back to per-key deletion to avoid CROSSSLOT errors.

## Capabilities

### New Capabilities
- `gredis-pipeline-tx-cluster`: Pipeline, Transaction, optimistic-locking, and Cluster-safe operation support for the Redis adapter layer.

### Modified Capabilities
- None.

## Impact

- `database/gredis/gredis_adapter.go` — extend `AdapterOperation` interface
- `database/gredis/gredis_redis.go` — add Pipeline/TxPipeline/Watch methods to `Redis`
- `database/gredis/gredis_redis_group_generic.go` — add `ScanAll` to `IGroupGeneric`
- `database/gredis/` — new files for `Cmd`, `Pipeliner`, pipeline group interfaces
- `contrib/nosql/redis/redis.go` — implement Pipeline/TxPipeline/Watch on the driver
- `contrib/nosql/redis/redis_group_generic.go` — implement ScanAll and Cluster-safe Del
- `contrib/nosql/redis/` — new file for pipeline group implementations
- New and updated test files across both modules
