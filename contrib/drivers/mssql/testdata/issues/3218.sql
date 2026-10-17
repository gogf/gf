CREATE TABLE [issue3218_sys_config] (
    [id] int NOT NULL,
    [name] nvarchar(255) NULL DEFAULT NULL,
    [value] nvarchar(max) NULL,
    [created_at] datetime2(0) NULL DEFAULT NULL,
    [updated_at] datetime2(0) NULL DEFAULT NULL,
    PRIMARY KEY ([id]),
    UNIQUE ([name])
);
