# Design

The SQLite drivers will prepare `Save` insert options before delegating to the shared `gdb.Core.DoInsert` implementation. When `InsertOptionSave` is used without explicit conflict columns, each driver reads the table primary keys through `Core.GetPrimaryKeys`, verifies that the first save record includes those primary key fields, and sets `DoInsertOption.OnConflict` accordingly. Existing explicit `OnConflict(...)` behavior remains unchanged.

If no usable primary key can be found in the save data, the drivers will return a missing-parameter error with guidance to specify `OnConflict(...)` or include primary key values.
