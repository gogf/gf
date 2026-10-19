CREATE TABLE [issue4034] (
    [id] int NOT NULL,
    [passport] varchar(255) NULL,
    [password] varchar(255) NULL,
    [nickname] varchar(255) NULL,
    [created_at] datetime2(0) NULL DEFAULT CURRENT_TIMESTAMP,
    [updated_at] datetime2(0) NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ([id])
);
