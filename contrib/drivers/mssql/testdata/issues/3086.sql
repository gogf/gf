CREATE TABLE [issue3086_user] (
    [id] int NOT NULL,
    [passport] varchar(45) NOT NULL,
    [password] varchar(45) NULL DEFAULT NULL,
    [nickname] varchar(45) NULL DEFAULT NULL,
    [create_at] datetime2(0) NULL DEFAULT NULL,
    [update_at] datetime2(0) NULL DEFAULT NULL,
    PRIMARY KEY ([id])
);
