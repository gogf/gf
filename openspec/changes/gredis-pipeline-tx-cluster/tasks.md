## 1. Abstract Layer — Cmd Type & Pipeline Interfaces (`database/gredis/`)

- [x] 1.1 Create `gredis_cmd.go` — define `Cmd` struct with `val *gvar.Var`, `err error`, `Result()`, `Val()` methods
- [x] 1.2 Create `gredis_pipeline.go` — define `Pipeliner`, `PipelinerOperation`, `PipelinerGroup` interfaces
- [x] 1.3 Define pipeline group interfaces (`IPipelineGroupGeneric`, `IPipelineGroupHash`, `IPipelineGroupString`, `IPipelineGroupList`, `IPipelineGroupSet`, `IPipelineGroupSortedSet`) — mirror existing group method signatures with `*Cmd` return type
- [x] 1.4 Define `Tx` interface — embeds `Pipeliner`
- [x] 1.5 Extend `AdapterOperation` in `gredis_adapter.go` — add `Pipeline()`, `TxPipeline()`, `Watch()`
- [x] 1.6 Add `Pipeline()`, `TxPipeline()`, `Watch()` wrapper methods on `Redis` in `gredis_redis.go`
- [x] 1.7 Add `ScanAll` method to `IGroupGeneric` in `gredis_redis_group_generic.go`
- [x] 1.8 Verify: `go build ./database/gredis/...` compiles, `go vet` passes

## 2. Concrete Driver — Pipeline & Transaction (`contrib/nosql/redis/`)

- [x] 2.1 Create `redis_pipeline.go` — implement `Pipeliner` interface wrapping go-redis `redis.Pipeliner`
- [x] 2.2 Implement all pipeline group structs (`pipelineGroupGeneric`, `pipelineGroupHash`, etc.) delegating to `redis.Pipeliner` methods
- [x] 2.3 Implement `Cmd` population logic — translate go-redis `Cmder` results to `gredis.Cmd` after `Exec()`
- [x] 2.4 Implement `Pipeline()`, `TxPipeline()`, `Watch()` on the `*Redis` driver type in `redis.go`
- [x] 2.5 Verify: `go build ./...` in `contrib/nosql/redis/` compiles

## 3. Concrete Driver — Cluster Enhancements (`contrib/nosql/redis/`)

- [x] 3.1 Implement `ScanAll` in `redis_group_generic.go` — standalone/sentinel: loop Scan until cursor=0; cluster: `ForEachMaster` + per-node scan aggregation
- [x] 3.2 Implement Cluster-safe `Del` in `redis_group_generic.go` — detect `*redis.ClusterClient`, fall back to per-key DEL, return total deleted count
- [x] 3.3 Verify: `go build ./...` in `contrib/nosql/redis/` compiles

## 4. Unit Tests — Abstract Layer (`database/gredis/`)

- [x] 4.1 Test `Cmd` type — `Result()`, `Val()` before and after population
- [x] 4.2 Test `Redis.Pipeline()` returns non-nil `Pipeliner`
- [x] 4.3 Test `Redis.TxPipeline()` returns non-nil `Pipeliner`
- [x] 4.4 Test `ScanAll` interface exists on `IGroupGeneric`
- [x] 4.5 Verify: `go test ./database/gredis/... -count=1 -race` passes, coverage ≥ 80%

## 5. Unit Tests — Concrete Driver (`contrib/nosql/redis/`)

- [x] 5.1 Test Pipeline basic operations — queue HSet + Get, Exec, verify Cmd results populated
- [x] 5.2 Test Pipeline Discard — queue commands, Discard, verify no server interaction
- [x] 5.3 Test TxPipeline — queue multiple commands, Exec, verify atomic execution
- [x] 5.4 Test Watch — optimistic locking success and abort scenarios
- [x] 5.5 Test ScanAll in standalone mode — verify all matching keys returned
- [x] 5.6 Test Cluster-safe Del — verify per-key deletion path (may require mock/skip if no cluster)
- [x] 5.7 Verify: `go test ./... -count=1 -race` in `contrib/nosql/redis/` passes, coverage ≥ 80%

## 6. Lint & Tidy

- [x] 6.1 Run `make tidy` from repo root
- [x] 6.2 Run `make lint` and fix any issues
- [x] 6.3 Run `go build ./...` from repo root to verify no breakage
