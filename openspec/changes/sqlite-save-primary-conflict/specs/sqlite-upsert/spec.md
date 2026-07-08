# sqlite-upsert Specification

## ADDED Requirements

### Requirement: SQLite Save infers primary-key conflicts
The SQLite and SQLiteCGO drivers SHALL infer the table primary key columns as the upsert conflict target for `Save` operations when callers do not explicitly provide `OnConflict(...)`.

#### Scenario: Save data includes the primary key
- **WHEN** a caller executes `Save` against a SQLite table without `OnConflict(...)`
- **AND** the table has primary key columns present in the save data
- **THEN** the driver SHALL use those primary key columns as the conflict target
- **AND** the save SHALL insert new rows or update existing rows using SQLite upsert syntax

#### Scenario: Save data has no usable primary key
- **WHEN** a caller executes `Save` against a SQLite table without `OnConflict(...)`
- **AND** the table does not have primary key columns present in the save data
- **THEN** the driver SHALL return a missing-parameter error instructing the caller to specify `OnConflict(...)` or include primary key values

#### Scenario: Explicit conflict columns are provided
- **WHEN** a caller executes `Save` against a SQLite table with `OnConflict(...)`
- **THEN** the driver SHALL use the explicitly provided conflict columns
