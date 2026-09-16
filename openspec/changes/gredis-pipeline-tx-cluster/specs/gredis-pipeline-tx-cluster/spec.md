## ADDED Requirements

### Requirement: Pipeline support in AdapterOperation
The `AdapterOperation` interface SHALL provide a `Pipeline` method that returns a `Pipeliner` instance for batching multiple Redis commands into a single network round-trip.

#### Scenario: User batches multiple commands via Pipeline
- **WHEN** a user calls `Pipeline(ctx)` on a Redis client
- **THEN** the method SHALL return a `Pipeliner` that queues commands locally without sending them to the server until `Exec` is called

#### Scenario: Pipeline Exec sends all queued commands
- **WHEN** a user calls `Exec(ctx)` on a `Pipeliner` with queued commands
- **THEN** all queued commands SHALL be sent to the Redis server in a single batch
- **AND** each queued command's `Cmd` result SHALL be populated with the server's reply

#### Scenario: Pipeline Discard clears queued commands
- **WHEN** a user calls `Discard()` on a `Pipeliner` with queued commands
- **THEN** all queued commands SHALL be discarded without being sent to the server

### Requirement: Transaction support via TxPipeline
The `AdapterOperation` interface SHALL provide a `TxPipeline` method that returns a `Pipeliner` wrapping commands in a Redis MULTI/EXEC transaction.

#### Scenario: TxPipeline executes atomically
- **WHEN** a user queues commands via `TxPipeline(ctx)` and calls `Exec(ctx)`
- **THEN** all queued commands SHALL be wrapped in MULTI/EXEC and executed atomically by the Redis server

### Requirement: Optimistic locking via Watch
The `AdapterOperation` interface SHALL provide a `Watch` method that accepts a callback function and a list of keys to watch for changes.

#### Scenario: Watch detects key modification
- **WHEN** a watched key is modified by another client before the transaction executes
- **THEN** the transaction SHALL be aborted and `Watch` SHALL return a transaction-abort error

#### Scenario: Watch succeeds without modification
- **WHEN** no watched key is modified before the transaction executes
- **THEN** the callback SHALL receive a `Tx` instance and the queued commands SHALL execute successfully

### Requirement: Pipeliner group-style interface
The `Pipeliner` interface SHALL provide typed command access through pipeline group interfaces (`IPipelineGroupGeneric`, `IPipelineGroupHash`, `IPipelineGroupString`, `IPipelineGroupList`, `IPipelineGroupSet`, `IPipelineGroupSortedSet`) mirroring the existing group interfaces with `*Cmd` return types.

#### Scenario: Pipeline Hash commands
- **WHEN** a user calls `Pipeliner.PipelineGroupHash().HSet(ctx, key, fields)`
- **THEN** the command SHALL be queued and a `*Cmd` future SHALL be returned
- **AND** the `Cmd` result SHALL be populated after `Exec(ctx)` is called

### Requirement: Cluster-safe ScanAll
The `IGroupGeneric` interface SHALL provide a `ScanAll` method that returns all keys matching a pattern, transparently handling Cluster mode by iterating all master nodes.

#### Scenario: ScanAll in standalone mode
- **WHEN** `ScanAll` is called on a standalone Redis connection
- **THEN** it SHALL repeatedly call `Scan` until the cursor returns to 0 and return all accumulated keys

#### Scenario: ScanAll in Cluster mode
- **WHEN** `ScanAll` is called on a Redis Cluster connection
- **THEN** it SHALL iterate all master nodes, scan each node, and return aggregated results from all nodes

### Requirement: Cluster-safe Del
The driver implementation of `Del` SHALL automatically handle multi-key deletion in Cluster mode by performing per-key deletion to avoid CROSSSLOT errors.

#### Scenario: Del multiple keys in standalone mode
- **WHEN** `Del(ctx, key1, key2, key3)` is called on a standalone Redis connection
- **THEN** a single DEL command with all keys SHALL be sent to the server

#### Scenario: Del multiple keys in Cluster mode
- **WHEN** `Del(ctx, key1, key2, key3)` is called on a Redis Cluster connection
- **THEN** the driver SHALL perform per-key DEL operations to avoid CROSSSLOT errors
- **AND** the driver SHALL return the total number of keys deleted across all operations

### Requirement: Driver-agnostic Pipeline interface
The `Pipeliner`, `Cmd`, and `Tx` types defined in `database/gredis` SHALL NOT reference any driver-specific types (e.g., `go-redis` types). The concrete implementation in `contrib/nosql/redis` SHALL translate between the framework types and go-redis types internally.

#### Scenario: Custom adapter implements Pipeline
- **WHEN** a third-party adapter implements the `AdapterOperation` interface
- **THEN** it SHALL be able to implement `Pipeline`, `TxPipeline`, and `Watch` without depending on go-redis
