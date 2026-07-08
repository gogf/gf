# SQLite Save Primary Conflict

SQLite `Save` currently requires callers to specify `OnConflict(...)` explicitly, even when the target table has a primary key and the save data includes that primary key. This change aligns the SQLite and SQLiteCGO drivers with other upsert-capable drivers by automatically using table primary keys as the conflict target when `OnConflict(...)` is omitted.
