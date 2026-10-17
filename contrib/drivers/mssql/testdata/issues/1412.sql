DROP TABLE IF EXISTS [items];
CREATE TABLE [items] (
    [id] int NOT NULL,
    [name] nvarchar(255) NULL DEFAULT NULL,
    PRIMARY KEY ([id])
);
INSERT INTO [items] ([id], [name]) VALUES (1, N'金秋产品1');
INSERT INTO [items] ([id], [name]) VALUES (2, N'金秋产品2');

DROP TABLE IF EXISTS [parcels];
CREATE TABLE [parcels] (
    [id] int NOT NULL,
    [item_id] int NULL DEFAULT NULL,
    PRIMARY KEY ([id])
);
INSERT INTO [parcels] ([id], [item_id]) VALUES (1, 1);
INSERT INTO [parcels] ([id], [item_id]) VALUES (2, 2);
INSERT INTO [parcels] ([id], [item_id]) VALUES (3, 0);
